package main

import(
	"net"
	tcpsession "kendynet-go/tcpsession"
	packet "kendynet-go/packet"
	"fmt"
)


func process_client(session *tcpsession.Tcpsession,rpk packet.Packet){
	if rpk == nil{
		session.Close()
		return
	}
	session.Send(rpk)
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
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		session := tcpsession.NewTcpSession(conn)
		fmt.Printf("a client comming\n")
		go tcpsession.ProcessSession(session,process_client,packet.NewRawDecoder())
	}
}


