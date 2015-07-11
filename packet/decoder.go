package packet

import "net"

type Decoder interface{
	DoRecv(Conn *net.Conn)(*Packet,error)
}