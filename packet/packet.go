package packet

const (
	RAWPACKET = 1
	RPACKET   = 2
	WPACKET   = 3
	EPACKET   = 4
)

type Packet interface{
	MakeWrite()(Packet)
	MakeRead() (Packet)
	Clone()    (Packet)
	PkLen()    (uint32)
	DataLen()  (uint32)
	Buffer()   (*ByteBuffer)
	GetType()  (byte)
}