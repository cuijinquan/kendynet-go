package main

import(
	"net"
	tcpsession "kendynet-go/tcpsession"
	packet "kendynet-go/packet"
	"fmt"
)

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
		session.SetRecvTimeout(5000)
		go tcpsession.ProcessSession(session,packet.NewRawDecoder(),
		   func (session *tcpsession.Tcpsession,rpk packet.Packet,errno error){	
			if rpk == nil{
				fmt.Printf("error:%s\n",errno)
				session.Close()
				return
			}
			session.Send(rpk)
		   })
	}
}


