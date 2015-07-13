package packet
import "unsafe"

type RawPacket struct{
	buffer *ByteBuffer
	tt	    byte
}

func NewRawPacket(buffer *ByteBuffer)(*RawPacket){
	if buffer == nil {
		return nil
	}
	return &RawPacket{buffer:buffer,tt:RAWPACKET}
}

func (this RawPacket) Buffer()(*ByteBuffer){
	return this.buffer
}

func (this RawPacket) Clone() (*Packet){
	return (*Packet)(unsafe.Pointer(NewRawPacket(this.buffer)))
}


func (this RawPacket) MakeWrite()(*Packet){
	return this.Clone()
}

func (this RawPacket) MakeRead()(*Packet){
	return (*Packet)(unsafe.Pointer(NewRawPacket(this.buffer)))
}

func (this RawPacket) DataLen()(uint32){
	if this.buffer == nil {
		return 0
	}
	return this.buffer.Len()
}

func (this RawPacket) PkLen()(uint32){
	return this.DataLen()
}

func (this RawPacket) GetType()(byte){
	return this.tt
}
