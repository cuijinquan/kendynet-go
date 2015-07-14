package main

import(
	"net"
	"fmt"
	"time"
	socket "kendynet-go/socket"
	packet "kendynet-go/packet"		
)

func main(){
	var i uint32
	for i = 1; i <= 100; i++ {
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
			wpk := packet.NewWPacket(packet.NewByteBuffer(uint32(64)))
			wpk.PutUint32(uint32(i))
			session.Send(wpk)
			idx := i	
			go socket.ProcessSession(session,packet.NewRPacketDecoder(1024),
				func (session *socket.Tcpsession,rpk packet.Packet,errno error){	
					if rpk == nil{
						fmt.Printf("error:%s\n",errno)
						session.Close()
						return
					}
					r  := rpk.(*packet.RPacket)
					id,_:= r.Uint32()
					if id == idx {
						session.Send(r.MakeWrite())
					}
			})
		}
	}

	for{
		time.Sleep(10000000)
	}
}