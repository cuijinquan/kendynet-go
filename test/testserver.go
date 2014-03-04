package main

import(
	"net"
	tcpsession "kendynet-go/tcpsession"
	packet "kendynet-go/packet"
	"fmt"
)

func send_finish (s interface{},wpk *packet.Wpacket){
	session := s.(*tcpsession.Tcpsession)
	session.Close()
}

func process_client(session *tcpsession.Tcpsession,rpk *packet.Rpacket){
	session.Send(packet.NewWpacket(rpk.Buffer(),rpk.IsRaw()),send_finish)
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
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil{
		fmt.Printf("ListenTCP")
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		session := tcpsession.NewTcpSession(conn,true)
		fmt.Printf("a client comming\n")
		go tcpsession.ProcessSession(session,process_client,session_close)
	}
}


