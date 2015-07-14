package packet



type RPacket struct{
	buffer *ByteBuffer
	tt	    byte
	readIdx uint32
}

func NewRPacket(buffer *ByteBuffer)(*RPacket){
	if buffer == nil {
		return nil
	}
	return &RPacket{readIdx:4,buffer:buffer,tt:RPACKET}
}

func (this RPacket) Buffer()(*ByteBuffer){
	return this.buffer
}

func (this RPacket) Clone() (Packet){
	rpk := &RPacket{readIdx:this.readIdx,buffer:this.buffer,tt:RPACKET}
	return *rpk
}


func (this RPacket) MakeWrite()(Packet){
	return *NewWPacket(this.buffer)
}

func (this RPacket) MakeRead()(Packet){
	return this.Clone()
}

func (this *RPacket) Uint16()(uint16,error){
	value,err := this.buffer.Uint16(this.readIdx)
	if err != nil {
		return 0,err
	}
	this.readIdx += 2
	return value,nil
}

func (this *RPacket) Uint32()(uint32,error){
	value,err := this.buffer.Uint32(this.readIdx)
	if err != nil {
		return 0,err
	}
	this.readIdx += 4
	return value,nil
}

func (this *RPacket) String()(string,error){
	value,err := this.buffer.String(this.readIdx)
	if err != nil {
		return "",err
	}
	this.readIdx += (4 + (uint32)(len(value)))
	return value,nil

}

func (this *RPacket) Binary()([]byte,error){
	value,err := this.buffer.Binary(this.readIdx)
	if err != nil {
		return nil,err
	}
	this.readIdx += (4 + (uint32)(len(value)))
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
	return this.tt
}
