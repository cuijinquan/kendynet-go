/*
 * 文件传送服务器 
 * 
 */
  
package main

import(
	"net"
	tcpsession "kendynet-go/tcpsession"
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

func (this *transfer_session)send_file(session *tcpsession.Tcpsession){
	remain := len(this.filecontent) - this.ridx
	sendsize := 0
	if remain >= 16000 {
		sendsize = 16000
	}else{
		sendsize = remain
	}
	wpk := packet.NewWpacket(packet.NewByteBuffer(uint32(sendsize)),false)
	wpk.PutUint16(transfering)
	wpk.PutBinary(this.filecontent[this.ridx:this.ridx+sendsize])
	session.Send(wpk,send_finish)
	this.ridx += sendsize
}

func (this *transfer_session)check_finish()(bool){
	if this.ridx >= len(this.filecontent) {
		return true
	}
	return false
}



func send_finish (s interface{},wpk *packet.Wpacket){
	session := s.(*tcpsession.Tcpsession)
	tsession := session.Ud().(*transfer_session)
	if tsession.check_finish(){
		session.Close()
		return
	}
	tsession.send_file(session)
}

func process_client(session *tcpsession.Tcpsession,rpk *packet.Rpacket){
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
				
				wpk := packet.NewWpacket(packet.NewByteBuffer(64),false)
				wpk.PutUint16(file_size)
				wpk.PutUint32(uint32(len(filecontent)))
				session.Send(wpk,nil)
				tsession.send_file(session)
			}	
		}
	}else{
		fmt.Printf("cmd error,%d\n",cmd)
		session.Close()
	}
}

func session_close(session *tcpsession.Tcpsession){
	fmt.Printf("client disconnect\n")
}

func loadfile(){
	//从配置导入文件
	F,err := os.Open("./config.txt")
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
			line = line[0:len(line)-1]//drop '\n'
			fileconfig := strings.Split(line,"=")
			if len(fileconfig) == 2 {
				buf, err := ioutil.ReadFile(fileconfig[0])
				if err != nil {
					fmt.Printf("%s load error\n",fileconfig[0])
				}else{	
					filemap[fileconfig[1]] = buf
					fmt.Printf("%s load success,key %s\n",fileconfig[0],fileconfig[1])
				}
			}
		}
	}
	
	if filemap["golang"] == nil {
		fmt.Printf("golang not found\n")
	}
	
	fmt.Printf("loadfile finish\n")	
}

func main(){	
	service := ":8010"
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
		session := tcpsession.NewTcpSession(conn,false)
		fmt.Printf("a client comming\n")
		go tcpsession.ProcessSession(session,process_client,session_close)
	}
}


