package packet

type Rpacket struct{
	readidx uint32
	buffer *ByteBuffer
}

func NewRpacket(buffer *ByteBuffer)(*Rpacket){
	if buffer == nil {
		return nil
	}
	return &Rpacket{readidx:4,buffer:buffer}
}

func (this *Rpacket) Buffer()(*ByteBuffer){
	return this.buffer
}

func (this *Rpacket) Len()(uint32){
	if this.buffer == nil {
		return 0
	}

	len,err := this.buffer.Uint32(0)
	if err != nil {
		return 0
	}
	return len
}

func (this *Rpacket) Uint16()(uint16,error){
	value,err := this.buffer.Uint16(this.readidx)
	if err != nil {
		return 0,err
	}
	this.readidx += 2
	return value,nil
}

func (this *Rpacket) Uint32()(uint32,error){
	value,err := this.buffer.Uint32(this.readidx)
	if err != nil {
		return 0,err
	}
	this.readidx += 4
	return value,nil
}

func (this *Rpacket) String()(string,error){
	value,err := this.buffer.String(this.readidx)
	if err != nil {
		return "",err
	}
	this.readidx += (4 + (uint32)(len(value)) + 1)
	return value,nil

}

func (this *Rpacket) Binary()([]byte,error){
	value,err := this.buffer.Binary(this.readidx)
	if err != nil {
		return nil,err
	}
	this.readidx += (4 + (uint32)(len(value)))
	return value,nil

}
