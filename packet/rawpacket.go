package packet


type RawPacket struct{
	buffer *ByteBuffer
	tt	    byte
}

func NewRawPacket(buffer *ByteBuffer)(*RawPacket) {
	if buffer == nil {
		buffer = NewByteBuffer(64)
	}
	return &RawPacket{buffer:buffer,tt:RAWPACKET}
}

func (this *RawPacket) PutBinary(value []byte)(error){
	err := this.buffer.PutBinary(this.buffer.Len(),value)
	if err != nil{
		return err
	}
	return nil
}

func (this RawPacket) Buffer()(*ByteBuffer) {
	return this.buffer
}

func (this RawPacket) Clone() (Packet) {
	return *NewRawPacket(this.buffer)
}


func (this RawPacket) MakeWrite()(Packet) {
	return this.Clone()
}

func (this RawPacket) MakeRead()(Packet) {
	return *NewRawPacket(this.buffer)
}

func (this RawPacket) DataLen()(uint32) {
	if this.buffer == nil {
		return 0
	}
	return this.buffer.Len()
}

func (this RawPacket) PkLen()(uint32) {
	return this.DataLen()
}

func (this RawPacket) GetType()(byte) {
	return this.tt
}
