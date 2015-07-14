/*
 * 文件传送服务器 
 * 
 */
  
package main

import(
	"net"
	socket "kendynet-go/socket"
	packet "kendynet-go/packet"
	"fmt"
	"strings"
	"os"
	"io/ioutil"
	"bufio"
	"io"
)

const (
	request_file = 1
	file_size = 2
	transfering = 3
)

var filemap map[string][]byte

type transfer_session struct{
	filecontent []byte
	ridx        int
}

func (this *transfer_session)send_file(session *socket.Tcpsession){
	for{
		remain := len(this.filecontent) - this.ridx
		sendsize := 0
		if remain >= 16000 {
			sendsize = 16000
		}else{
			sendsize = remain
		}
		wpk := packet.NewWPacket(packet.NewByteBuffer(uint32(sendsize)))
		wpk.PutUint16(transfering)
		wpk.PutBinary(this.filecontent[this.ridx:this.ridx+sendsize])
		if nil == session.Send(wpk){
			this.ridx += sendsize
			if this.ridx >= len(this.filecontent){
				break
			}
		}else{
			break
		}
	}
}

func process_client(session *socket.Tcpsession,p packet.Packet,_ error){
	rpk := p.(*packet.RPacket)
	cmd,_ := rpk.Uint16()
	if cmd == request_file {
		if session.Ud() != nil {
			fmt.Printf("already in transfer session\n")
		}else
		{
			filename,_ := rpk.String()
			filecontent := filemap[filename]
			if filecontent == nil {
				fmt.Printf("%s not found\n",filename)
				session.Close()
			}else{
				fmt.Printf("request file %s\n",filename)
				tsession := &transfer_session{filecontent:filecontent,ridx:0}
				session.SetUd(tsession)
				wpk := packet.NewWPacket(packet.NewByteBuffer(64))
				wpk.PutUint16(file_size)
				wpk.PutUint32(uint32(len(filecontent)))
				if nil == session.Send(wpk){
					tsession.send_file(session)
				}
			}	
		}
	}
	session.Close()
}

func drop_linebreak(input string)(string){
	size := len(input)
	if size > 2 && input[size-2] == '\r' && input[size-1] == '\n' {
		return input[0:size-2]
	}else if size > 1 && input[size-1] == '\n' {
		return input[0:size-1]
	}
	return input
}

func loadfile(){
	//从配置导入文件
	F,err := os.Open("./test/config.txt")
	if err != nil {
		fmt.Printf("config.txt open failed\n")
		return
	}
	filemap = make(map[string][]byte)
	bufferReader := bufio.NewReader(F)
	eof := false
	for !eof {
		line,err := bufferReader.ReadString('\n')
		if err == io.EOF{
			err = nil
			eof = true
		}else if err != nil{
			fmt.Printf("parse file error\n")
			return
		}
		if len(line) > 1 {
			line = drop_linebreak(line)//去掉行尾换行符
			fileconfig := strings.Split(line,"=")
			if len(fileconfig) == 2 {
				buf, err := ioutil.ReadFile(fileconfig[0])
				if err != nil || buf == nil{
					fmt.Printf("%s load error\n",fileconfig[0])
				}else{	
					filemap[fileconfig[1]] = buf
					fmt.Printf("%s load success,key %s\n",fileconfig[0],fileconfig[1])
				}
			}
		}
	}
	fmt.Printf("loadfile finish\n")	
}

func main(){	
	service := "127.0.0.1:8010"
	tcpAddr,err := net.ResolveTCPAddr("tcp4", service)
	if err != nil{
		fmt.Printf("ResolveTCPAddr")
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil{
		fmt.Printf("ListenTCP")
	}
	loadfile()
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		session := socket.NewTcpSession(conn)
		fmt.Printf("a client comming\n")
		go socket.ProcessSession(session,packet.NewRPacketDecoder(4096),process_client)
	}
}


