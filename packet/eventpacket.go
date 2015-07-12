package packet

type EventPacket struct{
	Type	byte
	errno	error
}

func NewEventPacket(errno error)(EventPacket){
	return EventPacket{errno:errno,Type:EPACKET}
}

func (this EventPacket) Buffer()(*ByteBuffer){
	return nil
}

func (this EventPacket)Clone() (*Packet){
	return nil
}

func (this EventPacket)MakeWrite()(*Packet){
	return nil
}

func (this EventPacket)MakeRead()(*Packet){
	return nil
}

func (this EventPacket) DataLen()(uint32){
	return 0
}

func (this EventPacket) PkLen()(uint32){
	return 0
}

func (this EventPacket) GetType()(byte){
	return this.Type
}

func (this EventPacket) GetError()(error){
	return this.errno
}
