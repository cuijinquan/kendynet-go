/*
 * 文件传送客户端 
 * 
 */ 

package main

import(
	"net"
	socket "kendynet-go/socket"
	packet "kendynet-go/packet"
	"fmt"
	"io/ioutil"
	"os"
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
	filesize    int
}

func (this *transfer_session)recv_file(rpk packet.RPacket)(bool){
	content,_ := rpk.Binary()
	copy(this.filecontent[this.widx:],content[:])
	this.widx += len(content)
	if this.widx >= this.filesize {
		ioutil.WriteFile(this.filename, this.filecontent, 0x777)
		return true
	}
	return false
}

func process_client(session *socket.Tcpsession,p packet.Packet,_ error){
	rpk := p.(packet.RPacket)
	cmd,_ := rpk.Uint16()
	if cmd == file_size {
		if session.Ud() == nil {
			fmt.Printf("error\n")
			session.Close()
			return
		}
		tsession := session.Ud().(*transfer_session)
		filesize,_ := rpk.Uint32()
		fmt.Printf("file size:%d\n",filesize)
		tsession.widx = 0
		tsession.filesize = int(filesize)
		tsession.filecontent = make([]byte,filesize)
		
	}else if cmd == transfering {
		if session.Ud() == nil {
			fmt.Printf("close here\n")
			session.Close()
			return
		}
		tsession := session.Ud().(*transfer_session)
		if tsession.recv_file(rpk) {
			//传输完毕
			fmt.Printf("transfer finish\n")
			session.Close()
			return
		}
	}else{
		fmt.Printf("cmd error,%d\n",cmd)
		//session.Close()
	}
}

func main(){
	
	if len(os.Args) < 3 {
		fmt.Printf("usage ./transferclient <filename> <savefilename\n")
		return
	}
	service := "127.0.0.1:8010"
	tcpAddr,err := net.ResolveTCPAddr("tcp4", service)
	if err != nil{
		fmt.Printf("ResolveTCPAddr")
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Printf("DialTcp error,%s\n",err)
	}else{
		session := socket.NewTcpSession(conn)
		fmt.Printf("connect sucessful\n")
		//发出文件请求
		wpk := packet.NewWPacket(packet.NewByteBuffer(64))
		wpk.PutUint16(request_file)
		wpk.PutString(os.Args[1])
		session.Send(wpk)
		tsession := &transfer_session{filename:os.Args[2]}
		session.SetUd(tsession)	
		socket.ProcessSession(session,packet.NewRPacketDecoder(65535),process_client)
	}
}

