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
	bytescount  := int32(0)

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
		tmp := 	atomic.LoadInt32(&bytescount)
		atomic.StoreInt32(&bytescount,0)
		fmt.Printf("clientcount:%d,transrfer:%d mb/s\n",clientcount,tmp/1024/1024)
		})

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		session := socket.NewTcpSession(conn)
		atomic.AddInt32(&clientcount,1)
		session.SetRecvTimeout(5000)
		go socket.ProcessSession(session,packet.NewRawDecoder(65535),
		   func (session *socket.Tcpsession,rpk packet.Packet,errno error){	
			if rpk == nil{
				atomic.AddInt32(&clientcount,-1)
				fmt.Printf("error:%s\n",errno)
				session.Close()
				return
			}
			atomic.AddInt32(&bytescount,int32(rpk.PkLen()))
			session.Send(rpk)
		   })
	}
}


