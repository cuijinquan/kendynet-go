package main

import(
	"net"
	"fmt"
	"time"
	"sync/atomic"
	util   "kendynet-go/util"
	socket "kendynet-go/socket"
	packet "kendynet-go/packet"		
)

func main(){

	clientcount := int32(0)
	packetcount := int32(0)

	clients     := util.NewConnCurrMap()

	service := ":8010"
	tcpAddr,err := net.ResolveTCPAddr("tcp4", service)
	if err != nil{
		fmt.Printf("ResolveTCPAddr")
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil{
		fmt.Printf("ListenTCP")
	}

	ticker := util.DurationTicker()
	ticker.Start(1000,func (_ time.Time){
		tmp := 	atomic.LoadInt32(&packetcount)
		atomic.StoreInt32(&packetcount,0)
		fmt.Printf("clientcount:%d,packetcount:%d/s\n",clientcount,tmp)
		})

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		session := socket.NewTcpSession(conn)
		clients.Set(session,session)
		atomic.AddInt32(&clientcount,1)
		go socket.ProcessSession(session,packet.NewRPacketDecoder(1024),
		   func (session *socket.Tcpsession,rpk packet.Packet,errno error){	
			if rpk == nil {
				session.Close()
				atomic.AddInt32(&clientcount,-1)
				fmt.Printf("error:%s\n",errno)
				clients.Del(session)
				return
			}
			for _,v := range(clients.Vals()) {
				atomic.AddInt32(&packetcount,int32(1))
				wpk := rpk.MakeWrite()
				v.(*socket.Tcpsession).Send(wpk)
			}
		})
	}
}


