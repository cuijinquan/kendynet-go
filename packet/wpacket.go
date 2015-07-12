package packet
import "unsafe"
	
type WPacket struct{
	buffer     *ByteBuffer
	Type	    byte
	writeidx    uint32
	CopyCreate  byte
	Fn_sendfinish func(interface{},*WPacket)
}

func NewWPacket(buffer *ByteBuffer)(*WPacket){
	if buffer == nil {
		return nil
	}
	buffer.PutUint32(0,0)
	return &WPacket{writeidx:4,buffer:buffer,Type:WPACKET,CopyCreate:0}
}

func (this WPacket)Buffer()(*ByteBuffer){
	return this.buffer
}

func (this WPacket)Clone() (*Packet){
	wpk := NewWPacket(this.buffer)
	wpk.CopyCreate = 1
	return (*Packet)(unsafe.Pointer(wpk))
}


func (this WPacket)MakeWrite()(*Packet){
	return this.Clone()
}

func (this WPacket)MakeRead()(*Packet){
	return (*Packet)(unsafe.Pointer(NewRPacket(this.buffer)))
}

func (this WPacket)copyOnWrite(){
	if this.CopyCreate == 1 {
		this.buffer = this.buffer.Clone()
		this.CopyCreate = 0
	}
}

func (this *WPacket)PutUint16(value uint16)(error){
	if this.buffer == nil {
		return ErrInvaildData
	}
	this.copyOnWrite()
	size,err := this.buffer.Uint32(0)
	if err != nil {
		return err
	}
	err = this.buffer.PutUint16(this.writeidx,value)
	if err != nil{
		return err
	}
	size += 2
	this.writeidx += 2
	this.buffer.SetUint32(0,size)
	return nil
}

func (this *WPacket)PutUint32(value uint32)(error){
	this.copyOnWrite()
	size,err := this.buffer.Uint32(0)
	if err != nil {
		return err
	}
	err = this.buffer.PutUint32(this.writeidx,value)
	if err != nil{
		return err
	}
	size += 4
	this.writeidx += 4
	this.buffer.SetUint32(0,size)
	return nil
}

func (this *WPacket)PutString(value string)(error){
	this.copyOnWrite()
	size,err := this.buffer.Uint32(0)
	if err != nil {
		return err
	}
	err = this.buffer.PutString(this.writeidx,value)
	if err != nil{
		return err
	}
	size += (4+(uint32)(len(value)))
	this.writeidx += (4+(uint32)(len(value)))
	this.buffer.SetUint32(0,size)
	return nil
}

func (this *WPacket)PutBinary(value []byte)(error){
	this.copyOnWrite()
	size,err := this.buffer.Uint32(0)
	if err != nil {
		return err
	}
	err = this.buffer.PutBinary(this.writeidx,value)
	if err != nil{
		return err
	}
	size += (4+(uint32)(len(value)))
	this.writeidx += (4+(uint32)(len(value)))
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
	return this.Type
}
