package packet
import "unsafe"


type RPacket struct{
	buffer *ByteBuffer
	Type	byte
	readidx uint32
}

func NewRPacket(buffer *ByteBuffer)(*RPacket){
	if buffer == nil {
		return nil
	}
	return &RPacket{readidx:4,buffer:buffer,Type:RPACKET}
}

func (this RPacket) Buffer()(*ByteBuffer){
	return this.buffer
}

func (this RPacket)Clone() (*Packet){
	rpk := &RPacket{readidx:this.readidx,buffer:this.buffer,Type:RPACKET}
	return (*Packet)(unsafe.Pointer(rpk))
}


func (this RPacket)MakeWrite()(*Packet){
	return this.Clone()
}

func (this RPacket)MakeRead()(*Packet){
	return (*Packet)(unsafe.Pointer(NewWPacket(this.buffer)))
}

func (this *RPacket) Uint16()(uint16,error){
	value,err := this.buffer.Uint16(this.readidx)
	if err != nil {
		return 0,err
	}
	this.readidx += 2
	return value,nil
}

func (this *RPacket) Uint32()(uint32,error){
	value,err := this.buffer.Uint32(this.readidx)
	if err != nil {
		return 0,err
	}
	this.readidx += 4
	return value,nil
}

func (this *RPacket) String()(string,error){
	value,err := this.buffer.String(this.readidx)
	if err != nil {
		return "",err
	}
	this.readidx += (4 + (uint32)(len(value)))
	return value,nil

}

func (this *RPacket) Binary()([]byte,error){
	value,err := this.buffer.Binary(this.readidx)
	if err != nil {
		return nil,err
	}
	this.readidx += (4 + (uint32)(len(value)))
	return value,nil
}

func (this RPacket) DataLen()(uint32){
	if this.buffer == nil {
		return 0
	}
	len,err := this.buffer.Uint32(0)
	if err != nil {
		return 0
	}
	return len
}

func (this RPacket) PkLen()(uint32){
	if this.buffer == nil {
		return 0
	}
	return this.DataLen() + 4
}

func (this RPacket) GetType()(byte){
	return this.Type
}
