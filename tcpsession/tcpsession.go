package tcpsession

import(
	   "net"
	   packet "kendynet-go/packet"
	   "fmt"
   )

var (
	ErrUnPackError     = fmt.Errorf("TcpSession: UnpackError")
	ErrSendClose       = fmt.Errorf("send close")
	ErrSocketClose     = fmt.Errorf("socket close")
)


type Tcpsession struct{
	Conn         net.Conn
	Packet_que   chan packet.Packet
	decoder      packet.Decoder
	socket_close bool
	ud           interface{}
}

func (this *Tcpsession) SetUd(ud interface{}){
	this.ud = ud
}

func (this *Tcpsession) Ud()(interface{}){
	return this.ud
}

func dorecv(session *Tcpsession){
	for{
		p,err := session.decoder.DoRecv(session.Conn)
		if session.socket_close{
			break
		}
		if err != nil {
			session.Packet_que <- packet.NewEventPacket(err)
			break
		}
		session.Packet_que <- p	
	}
	close(session.Packet_que)
}


func ProcessSession(tcpsession *Tcpsession,
					process_packet func (*Tcpsession,packet.Packet,error),
					decoder packet.Decoder)(error){
	if tcpsession.socket_close{
		return ErrSocketClose
	}
	tcpsession.decoder = decoder
	go dorecv(tcpsession)
	for{
		msg,ok := <- tcpsession.Packet_que
		if !ok {
			//log error
			return nil
		}
		if packet.EPACKET == msg.GetType(){
			process_packet(tcpsession,nil,msg.(packet.EventPacket).GetError())
		}else{
			process_packet(tcpsession,msg,nil)
		}
		if tcpsession.socket_close{
			return nil
		}
	}
}

func NewTcpSession(conn net.Conn)(*Tcpsession){
	session := new(Tcpsession)
	session.Conn = conn
	session.Packet_que   = make(chan packet.Packet,1024)
	session.socket_close = false
	return session
}

func (this *Tcpsession)Send(wpk packet.Packet)(error){
	if this.socket_close{
		return ErrSocketClose
	}
	idx := (uint32)(0)
	for{
		buff  := wpk.Buffer().Bytes()
		end   := wpk.PkLen()
		n,err := this.Conn.Write(buff[idx:end])
		if err != nil || n < 0 {
			return ErrSendClose
		}
		idx += (uint32)(n)
		if idx >= (uint32)(end){
			break
		}
	}
	return nil
}

func (this *Tcpsession)Close(){
	if this.socket_close{
		return
	}
	this.socket_close = true
	this.Conn.Close()
}
