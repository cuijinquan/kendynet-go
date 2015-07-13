package packet

import(
	"net"
	"encoding/binary"
	"fmt"
	"io"
)

var (
	ErrPacketTooLarge     = fmt.Errorf("Packet too Large")
	ErrEOF                = fmt.Errorf("Eof")
)

type Decoder interface{
	DoRecv(Conn net.Conn)(Packet,error)
}

type RPacketDecoder struct{
	maxpacket uint32
}

func NewRPacketDecoder(maxpacket uint32)(RPacketDecoder){
	return RPacketDecoder{maxpacket:maxpacket}
}

func (this RPacketDecoder) DoRecv(Conn net.Conn)(Packet,error){
	header := []byte{0,0,0,0}
	n, err := io.ReadFull(Conn, header)
	if n == 0 && err == io.EOF {
		return nil,ErrEOF
	}else if err != nil {
		return nil,err
	}
	size := binary.LittleEndian.Uint32(header)
	if size > this.maxpacket {
		return nil,ErrPacketTooLarge
	}
	buf := make([]byte,size+4)
	copy(buf[:],header[:])
	n, err = io.ReadFull(Conn,buf[4:])
	if n == 0 && err == io.EOF {
		return nil,ErrEOF
	}else if err != nil {
		return nil,err
	}
	return NewRPacket(NewBufferByBytes(buf,(uint32)(len(buf)))),nil
}

type RawDecoder struct{
	buffsize uint32
}

func NewRawDecoder(buffsize uint32)(RawDecoder){
	return RawDecoder{buffsize:buffsize}
}

func (this RawDecoder) DoRecv(Conn net.Conn)(Packet,error){
	buff  := make([]byte,this.buffsize)
	n,err := Conn.Read(buff)
	if n == 0 && err == io.EOF {
		return nil,ErrEOF
	}else if err != nil {
		return nil,err
	}
	return NewRawPacket(NewBufferByBytes(buff,(uint32)(n))),nil		
}

