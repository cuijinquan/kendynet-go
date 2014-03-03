package packet

type rpacket struct{
	readidx uint32
	buffer *bytebuffer
}

func NewRpacket(buffer *bytebuffer)(*rpacket){
	if buffer == nil {
		return nil
	}
	return &rpacket{readidx:4,buffer:buffer}
}

func (this *rpacket) Buffer()(*bytebuffer){
	return this.buffer
}

func (this *rpacket) Len()(uint32){
	if this.buffer == nil {
		return 0
	}

	len,err := this.buffer.Uint32(0)
	if err != nil {
		return 0
	}
	return len
}

func (this *rpacket) Uint16()(uint16,error){
	value,err := this.buffer.Uint16(this.readidx)
	if err != nil {
		return 0,err
	}
	this.readidx += 2
	return value,nil
}

func (this *rpacket) Uint32()(uint32,error){
	value,err := this.buffer.Uint32(this.readidx)
	if err != nil {
		return 0,err
	}
	this.readidx += 4
	return value,nil
}

func (this *rpacket) String()(string,error){
	value,err := this.buffer.String(this.readidx)
	if err != nil {
		return "",err
	}
	this.readidx += (4 + (uint32)(len(value)) + 1)
	return value,nil

}

func (this *rpacket) Binary()([]byte,error){
	value,err := this.buffer.Binary(this.readidx)
	if err != nil {
		return nil,err
	}
	this.readidx += (4 + (uint32)(len(value)))
	return value,nil

}
