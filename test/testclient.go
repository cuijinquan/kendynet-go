package main

import(
	"net"
	"fmt"
	socket "kendynet-go/socket"
	packet "kendynet-go/packet"		
)


func main(){
	bytes := make([]byte,65535)
	//if len(os.Args) < 3 {
	//	fmt.Printf("usage ./transferclient <filename> <savefilename\n")
	//	return
	//}
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
		wpk := packet.NewRawPacket(packet.NewBufferByBytes(bytes,uint32(len(bytes))))
		session.Send(wpk)	
		socket.ProcessSession(session,packet.NewRawDecoder(65535),
			func (session *socket.Tcpsession,rpk packet.Packet,errno error){	
				if rpk == nil{
					fmt.Printf("error:%s\n",errno)
					session.Close()
					return
				}
				session.Send(rpk)
		})
	}
}