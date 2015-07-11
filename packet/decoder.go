package packet

type Decoder interface{
	UnPack()(*Packet)
}