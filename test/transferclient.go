package main

import(
	"net"
	tcpsession "kendynet-go/tcpsession"
	packet "kendynet-go/packet"
	"fmt"
	"io/ioutil"
)

const (
	request_file = 1
	file_size = 2
	transfering = 3
)

type transfer_session struct{
	filecontent []byte
	widx        int
	filename    string
}

func (this *transfer_session)recv_file(rpk *packet.Rpacket)(bool){
	content,_ := rpk.Binary()
	copy(content[:],this.filecontent[this.widx:])
	this.widx += len(content)
	if this.widx >= len(this.filecontent) {
		ioutil.WriteFile(this.filename, this.filecontent, 0x644)
		return true
	}
	return false
}

func process_client(session *tcpsession.Tcpsession,rpk *packet.Rpacket){
	cmd,_ := rpk.Uint16()
	if cmd == file_size {
		if session.Ud() == nil {
			session.Close()
			return
		}
		tsession := session.Ud().(transfer_session)
		filesize,_ := rpk.Uint32()
		tsession.widx = 0
		tsession.filecontent = make([]byte,filesize)
		
	}else if cmd == transfering {
		if session.Ud() == nil {
			session.Close()
			return
		}
		tsession := session.Ud().(transfer_session)
		if tsession.recv_file(rpk) {
			//传输完毕
			session.Close()
			return
		}
	}
}

func session_close(session *tcpsession.Tcpsession){
	fmt.Printf("client disconnect\n")
}

func main(){
	service := ":8010"
	tcpAddr,err := net.ResolveTCPAddr("tcp4", service)
	if err != nil{
		fmt.Printf("ResolveTCPAddr")
	}
	
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	
	if err != nil {
		fmt.Printf("DialTcp error\n")
	}else{
		session := tcpsession.NewTcpSession(conn,true)
		fmt.Printf("connect sucessful\n")
		//发出文件请求
		wpk := packet.NewWpacket(packet.NewByteBuffer(64),false)
		wpk.PutUint16(request_file)
		wpk.PutString("learnyouhaskell.pdf")
		session.Send(wpk,nil)
		tsession := &transfer_session{filename:"learnyouhaskell.pdf"}
		session.SetUd(tsession)	
		tcpsession.ProcessSession(session,process_client,session_close)
	}
}

