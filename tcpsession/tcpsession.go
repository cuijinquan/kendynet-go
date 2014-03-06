package tcpsession

import(
	   "net"
	   packet "kendynet-go/packet"
	   "encoding/binary"
	   "fmt"
   )

var (
	ErrUnPackError     = fmt.Errorf("TcpSession: UnpackError")
)

type Tcpsession struct{
	Conn net.Conn
	Packet_que chan interface{}
	Send_que chan *packet.Wpacket
	raw bool
	send_close bool
	ud interface{}
}

func (this *Tcpsession) IsRaw()(bool){
	return this.raw
}

func (this *Tcpsession) SetUd(ud interface{}){
	this.ud = ud
}

func (this *Tcpsession) Ud()(interface{}){
	return this.ud
}

func unpack(begidx uint32,buffer []byte,packet_que chan interface{})(int,error){
	unpack_size := 0
	for{
		packet_size :=	binary.LittleEndian.Uint32(buffer[begidx:begidx+4])
		if packet_size > packet.Max_bufsize-4 {
			return 0,ErrUnPackError
		}
		if packet_size+4 <= (uint32)(len(buffer)){
			rpk := packet.NewRpacket(packet.NewBufferByBytes(buffer[begidx:(begidx+packet_size+4)]),false)
			packet_que <- rpk
			begidx += packet_size+4
			unpack_size += (int)(packet_size)+4
		}else{
			break
		}
	}
	return unpack_size,nil
}


func dorecv(session *Tcpsession){
	recvbuf := make([]byte,packet.Max_bufsize)
	unpackbuf := make([]byte,packet.Max_bufsize*2)
	unpack_idx := 0
	for{
		n,err := session.Conn.Read(recvbuf)
		if err != nil {
			session.Packet_que <- "rclose"
			return
		}
		//copy to unpackbuf
		copy(unpackbuf[len(unpackbuf):],recvbuf[:n])
		//unpack
		n,err = unpack((uint32)(unpack_idx),unpackbuf,session.Packet_que)
		if err != nil {
			close(session.Packet_que)
			return
		}
		unpack_idx += n
		if cap(unpackbuf) - len(unpackbuf) < (int)(packet.Max_bufsize) {
			tmpbuf := make([]byte,packet.Max_bufsize*2)
			n = len(unpackbuf) - unpack_idx
			if n > 0 {
				copy(tmpbuf[0:],unpackbuf[unpack_idx:unpack_idx+n])
			}
			unpackbuf = tmpbuf
			unpack_idx = 0
		}
	}
}

func dorecv_raw(session *Tcpsession){
	for{
		recvbuf := make([]byte,packet.Max_bufsize)
		_,err := session.Conn.Read(recvbuf)
		if err != nil {
			session.Packet_que <- "rclose"
			return
		}
		rpk := packet.NewRpacket(packet.NewBufferByBytes(recvbuf),true)
		session.Packet_que <- rpk
	}
}

func dosend(session *Tcpsession){
	for{
		wpk,ok :=  <-session.Send_que
		if !ok {
			return
		}
		_,err := session.Conn.Write(wpk.Buffer().Bytes())
		if err != nil {
			session.send_close = true
			return
		}
		if wpk.Fn_sendfinish != nil{
			wpk.Fn_sendfinish(session,wpk)
		}
	}
}


func ProcessSession(tcpsession *Tcpsession,process_packet func (*Tcpsession,*packet.Rpacket),session_close func (*Tcpsession)){
	for{
		msg,ok := <- tcpsession.Packet_que
		if !ok {
			fmt.Printf("client disconnect\n")
			return
		}
		switch msg.(type){
			case * packet.Rpacket:
				rpk := msg.(*packet.Rpacket)
				process_packet(tcpsession,rpk)
			case string:
				str := msg.(string)
				if str == "rclose"{
					session_close(tcpsession)
					close(tcpsession.Packet_que)
					close(tcpsession.Send_que)
					tcpsession.Conn.Close()
					return
				}
		}
	}
}

func NewTcpSession(conn net.Conn,raw bool)(*Tcpsession){
	session := &Tcpsession{Conn:conn,Packet_que:make(chan interface{},1024),Send_que:make(chan *packet.Wpacket,1024),raw:raw,send_close:false}
	if raw{
		go dorecv_raw(session)
	}else{
		go dorecv(session)
	}
	go dosend(session)
	return session
}

func (this *Tcpsession)Send(wpk *packet.Wpacket,send_finish func(interface{},*packet.Wpacket))(error){
	if !this.send_close{
		wpk.Fn_sendfinish = send_finish
		this.Send_que <- wpk
	}
	return nil
}

func (this *Tcpsession)Close(){
	this.Conn.Close()
}
