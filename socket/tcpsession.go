package socket

import(
	   "net"
	   packet "kendynet-go/packet"
	   "fmt"
	   "time"
	   "sync/atomic"
   )

var (
	ErrUnPackError     = fmt.Errorf("TcpSession: UnpackError")
	ErrSendClose       = fmt.Errorf("send close")
	ErrSocketClose     = fmt.Errorf("socket close")
)

var (
	SendClose          = packet.NewEventPacket(fmt.Errorf("SendClose"))
	RecvClose          = packet.NewEventPacket(fmt.Errorf("RecvClose"))
	NotifyClose        = packet.NewEventPacket(fmt.Errorf("notifyClose"))
)

type Tcpsession struct{
	Conn         net.Conn
	Packet_que   chan packet.Packet
	Send_que     chan packet.Packet
	decoder      packet.Decoder
	socket_close int32
	ud           interface{}
	recv_timeout uint64   //in ms
	send_timeout uint64   //in ms
}

func (this *Tcpsession) SetUd(ud interface{}){
	this.ud = ud
}

func (this *Tcpsession) Ud()(interface{}){
	return this.ud
}

func (this *Tcpsession) SetRecvTimeout(timeout uint64){
	this.recv_timeout = timeout
}

func (this *Tcpsession) SetSendTimeout(timeout uint64){
	this.send_timeout = timeout
}

func dorecv(session *Tcpsession){
	for{
		if session.recv_timeout > 0 {
			t := time.Now()
			session.Conn.SetReadDeadline(t.Add(time.Millisecond * time.Duration(session.recv_timeout)))
		}
		p,err := session.decoder.DoRecv(session.Conn)
		if 1 == atomic.LoadInt32(&session.socket_close) {
			break
		}
		if err != nil {
			session.Packet_que <- packet.NewEventPacket(err)
			break
		}
		session.Packet_que <- p	
	}
	session.Packet_que <- RecvClose
}

func dosend(session *Tcpsession) {
	for {
		wpk,ok := <- session.Send_que
		if !ok || wpk == NotifyClose {
			break
		}
		idx := (uint32)(0)
		for{
			buff  := wpk.Buffer().Bytes()
			end   := wpk.PkLen()
			if session.send_timeout > 0 {
				t := time.Now()
				session.Conn.SetWriteDeadline(t.Add(time.Millisecond * time.Duration(session.send_timeout)))
			}		
			n,err := session.Conn.Write(buff[idx:end])
			if err != nil {
				if 0 == atomic.LoadInt32(&session.socket_close) {
					session.Packet_que <- packet.NewEventPacket(err)
				}else{
					session.Packet_que <- SendClose
					return
				}
				break
			}
			idx += (uint32)(n)
			if idx >= (uint32)(end){
				break
			}
		}
	}
	session.Packet_que <- SendClose
}


func ProcessSession(tcpsession *Tcpsession,decoder packet.Decoder,
					process_packet func (*Tcpsession,packet.Packet,error))(error) {
	if 1 == atomic.LoadInt32(&tcpsession.socket_close) {
		return ErrSocketClose
	}
	tcpsession.decoder = decoder
	go dorecv(tcpsession)
	go dosend(tcpsession)
	cc := 0
	for{
		msg,ok := <- tcpsession.Packet_que
		if !ok {
			//log error
			break
		}
		if msg == SendClose || msg == RecvClose{
			cc += 1
		}else if packet.EPACKET == msg.GetType(){
			process_packet(tcpsession,nil,msg.(packet.EventPacket).GetError())
		}else{
			process_packet(tcpsession,msg,nil)
		}
		if 2 == cc {
			break
		}
	}
	close(tcpsession.Packet_que)
	close(tcpsession.Send_que)
	tcpsession.Conn.Close()
	return nil
}

func NewTcpSession(conn net.Conn)(*Tcpsession){
	session 			:= new(Tcpsession)
	session.Conn 		 = conn
	session.Packet_que   = make(chan packet.Packet,64)
	session.Send_que     = make(chan packet.Packet,64)
	return session
}

func (this *Tcpsession) Send(wpk packet.Packet)(error){
	if 1 == atomic.LoadInt32(&this.socket_close) {
		return ErrSocketClose
	}
	this.Send_que <- wpk
	return nil
}

func (this *Tcpsession) Close(){
	if 1 == atomic.LoadInt32(&this.socket_close) {
		return
	}
	atomic.StoreInt32(&this.socket_close,1)
	tcpconn := this.Conn.(*net.TCPConn)
	tcpconn.CloseRead()
	this.Send_que <- NotifyClose
}
