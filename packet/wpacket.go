package packet

	
type WPacket struct{
	buffer     *ByteBuffer
	tt	        byte
	writeIdx    uint32
	copyCreate  byte
}

func NewWPacket(buffer *ByteBuffer)(*WPacket){
	if buffer == nil {
		return nil
	}
	if buffer.Len() != 0 {
		return &WPacket{writeIdx:buffer.Len(),buffer:buffer,tt:WPACKET,copyCreate:1}
	}else{
		buffer.PutUint32(0,0)
		return &WPacket{writeIdx:4,buffer:buffer,tt:WPACKET,copyCreate:0}
	}
}

func (this WPacket) Buffer()(*ByteBuffer) {
	return this.buffer
}

func (this WPacket) Clone() (Packet) {
	return *NewWPacket(this.buffer)
}

func (this WPacket) MakeWrite() (Packet) {
	return this.Clone()
}

func (this WPacket) MakeRead() (Packet) {
	return *NewRPacket(this.buffer)
}

func (this WPacket) copyOnWrite(){
	if this.copyCreate == 1 {
		this.buffer = this.buffer.Clone()
		this.copyCreate = 0
	}
}

func (this *WPacket) PutUint16(value uint16)(error) {
	if this.buffer == nil {
		return ErrInvaildData
	}
	this.copyOnWrite()
	size,err := this.buffer.Uint32(0)
	if err != nil {
		return err
	}
	err = this.buffer.PutUint16(this.writeIdx,value)
	if err != nil{
		return err
	}
	size += 2
	this.writeIdx += 2
	this.buffer.SetUint32(0,size)
	return nil
}

func (this *WPacket) PutUint32(value uint32)(error){
	this.copyOnWrite()
	size,err := this.buffer.Uint32(0)
	if err != nil {
		return err
	}
	err = this.buffer.PutUint32(this.writeIdx,value)
	if err != nil{
		return err
	}
	size += 4
	this.writeIdx += 4
	this.buffer.SetUint32(0,size)
	return nil
}

func (this *WPacket) PutString(value string)(error){
	this.copyOnWrite()
	size,err := this.buffer.Uint32(0)
	if err != nil {
		return err
	}
	err = this.buffer.PutString(this.writeIdx,value)
	if err != nil{
		return err
	}
	size += (4+(uint32)(len(value)))
	this.writeIdx += (4+(uint32)(len(value)))
	this.buffer.SetUint32(0,size)
	return nil
}

func (this *WPacket) PutBinary(value []byte)(error){
	this.copyOnWrite()
	size,err := this.buffer.Uint32(0)
	if err != nil {
		return err
	}
	err = this.buffer.PutBinary(this.writeIdx,value)
	if err != nil{
		return err
	}
	size += (4+(uint32)(len(value)))
	this.writeIdx += (4+(uint32)(len(value)))
	this.buffer.SetUint32(0,size)
	return nil
}

func (this WPacket) DataLen()(uint32){
	if this.buffer == nil {
		return 0
	}
	len,err := this.buffer.Uint32(0)
	if err != nil {
		return 0
	}
	return len
}

func (this WPacket) PkLen()(uint32){
	if this.buffer == nil {
		return 0
	}
	return this.DataLen() + 4
}

func (this WPacket) GetType()(byte){
	return this.tt
}
