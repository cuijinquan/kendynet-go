package packet

const (
	RAWPACKET = 1
	RPACKET
	WPACKET
)

type Packet interface{
	MakeWrite()(*Packet)
	MakeRead() (*Packet)
	Clone()    (*Packet)
	PkLen()    (uint32)
	DataLen()  (uint32)
}