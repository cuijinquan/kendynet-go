package tcpsession

import "net"
import packet "kendynet-go/packet"
import "encoding/binary"
import "fmt"

var (
	ErrUnPackError     = fmt.Errorf("TcpSession: UnpackError")
)

type Tcpsession struct{
	Conn net.Conn
	Packet_que chan *packet.Rpacket
	Send_que chan *packet.Wpacket
}


type tcprecver struct{
	Session *Tcpsession
}

type tcpsender struct{
	Session *Tcpsession
}


func unpack(begidx uint32,buffer []byte,packet_que chan *packet.Rpacket)(int,error){
	unpack_size := 0
	for{
		packet_size :=	binary.LittleEndian.Uint32(buffer[begidx:begidx+4])
		if packet_size > packet.Max_bufsize-4 {
			return 0,ErrUnPackError
		}
		if packet_size+4 <= (uint32)(len(buffer)){
			rpk := packet.NewRpacket(packet.NewBufferByBytes(buffer[begidx:(begidx+packet_size+4)]))
			packet_que <- rpk
			begidx += packet_size+4
			unpack_size += (int)(packet_size)+4
		}else{
			break
		}
	}
	return unpack_size,nil
}


func dorecv(recver *tcprecver){
	recvbuf := make([]byte,packet.Max_bufsize)
	unpackbuf := make([]byte,packet.Max_bufsize*2)
	unpack_idx := 0
	for{
		n,err := recver.Session.Conn.Read(recvbuf)
		if err != nil {
			close(recver.Session.Packet_que)
			return
		}
		//copy to unpackbuf
		copy(unpackbuf[len(unpackbuf):],recvbuf[:n])
		//unpack
		n,err = unpack((uint32)(unpack_idx),unpackbuf,recver.Session.Packet_que)
		if err != nil {
			close(recver.Session.Packet_que)
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

func dosend(sender *tcpsender){
	for{
		wpk,ok :=  <-sender.Session.Send_que
		if !ok {
			return
		}
		_,err := sender.Session.Conn.Write(wpk.Buffer().Bytes())
		if err != nil {
			close(sender.Session.Packet_que)
			return
		}

	}
}

func NewTcpSession(conn net.Conn)(*Tcpsession){
	session := &Tcpsession{Conn:conn,Packet_que:make(chan *packet.Rpacket,1024),Send_que:make(chan *packet.Wpacket,1024)}
	go dorecv(&tcprecver{Session:session})
	go dosend(&tcpsender{Session:session})
	return session
}

func (this *Tcpsession)Send(wpk *packet.Wpacket)(error){
	this.Send_que <- wpk
	return nil
}

func (this *Tcpsession)Process(){

}
