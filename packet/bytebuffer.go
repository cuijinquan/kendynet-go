package packet

import (
	"encoding/binary"
	"fmt"
)
func IsPow2(size uint32) bool{
	return (size&(size-1)) == 0
}

func SizeofPow2(size uint32) uint32{
	if IsPow2(size){
		return size
	}
	size = size -1
	size = size-1
	size = size | (size>>1)
	size = size | (size>>2)
	size = size | (size>>4)
	size = size | (size>>8)
	size = size | (size>>16)
	return size + 1
}

func GetPow2(size uint32) uint32{
	var pow2 uint32 = 0
	if !IsPow2(size) {
		size = (size << 1)
	}
	for size > 1 {
		pow2++
	}
	return pow2
}
const (
	Max_bufsize  uint32 = 32000 
	Max_string_len  uint32 = 32000
	Max_bin_len  uint32 = 32000
)

type ByteBuffer struct {
	buffer []byte
	datasize uint32
	capacity uint32
}


var (
	ErrMaxDataSlotsExceeded     = fmt.Errorf("bytebuffer: Max Buffer Size Exceeded")
	ErrInvaildData              = fmt.Errorf("bytebuffer: Invaild Data")
)

func NewBufferByBytes(bytes []byte,datasize uint32)(*ByteBuffer){
	return &ByteBuffer{buffer:bytes,datasize:datasize,capacity:(uint32)(cap(bytes))}
}


func NewByteBuffer(size uint32)(*ByteBuffer){
	if size == 0 {
		size = 64
	}else{
		size = SizeofPow2(size)
	}
	return &ByteBuffer{buffer:make([]byte,size),datasize:0,capacity:size}
}

func (this *ByteBuffer)Clone() (*ByteBuffer){
	b := make([]byte,this.capacity)
	copy(b[0:],this.buffer[:this.capacity])
	return &ByteBuffer{buffer:b,datasize:this.datasize,capacity:this.capacity}
}

func (this *ByteBuffer)Bytes()([]byte){
	return this.buffer
}

func (this *ByteBuffer)Len()(uint32){
	return this.datasize
}

func (this *ByteBuffer)Cap()(uint32){
	return this.capacity
}

func (this *ByteBuffer)expand(newsize uint32)(error){
	newsize = SizeofPow2(newsize)
	if newsize > Max_bufsize {
		return ErrMaxDataSlotsExceeded
	}
	//allocate new buffer
	tmpbuf := make([]byte,newsize)
	//copy data
	copy(tmpbuf[0:], this.buffer[:this.datasize])
	//replace buffer
	this.buffer = tmpbuf
	this.capacity = newsize
	return nil
}

func (this *ByteBuffer)buffer_check(idx,size uint32)(error){
	if idx+size > this.capacity {
		err := this.expand(idx+size)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *ByteBuffer)SetByte(idx uint32,value byte){
	this.buffer[idx] = value
}

func (this *ByteBuffer)PutByte(idx uint32,value byte)(error){
	err := this.buffer_check(idx,1)
	if err != nil {
		return err
	}
	this.buffer[idx] = value
	this.datasize += 1
	return nil
}

func (this *ByteBuffer)SetUint16(idx uint32,value uint16){
	binary.LittleEndian.PutUint16(this.buffer[idx:idx+2],value)
}

func (this *ByteBuffer)PutUint16(idx uint32,value uint16)(error){
	err := this.buffer_check(idx,2)
	if err != nil {
		return err
	}
	binary.LittleEndian.PutUint16(this.buffer[idx:idx+2],value)
	this.datasize += 2
	return nil
}

func (this *ByteBuffer)SetUint32(idx uint32,value uint32){
	binary.LittleEndian.PutUint32(this.buffer[idx:idx+4],value)
}

func (this *ByteBuffer)PutUint32(idx uint32,value uint32)(error){
	err := this.buffer_check(idx,4)//(uint32)(unsafe.Sizeof(value)))
	if err != nil {
		return err
	}
	binary.LittleEndian.PutUint32(this.buffer[idx:idx+4],value)
	this.datasize += 4
	return nil
}


func (this *ByteBuffer)PutString(idx uint32,value string)(error){
	sizeneed := (uint32)(4)
	sizeneed += (uint32)(len(value))
	err := this.buffer_check(idx,sizeneed)
	if err != nil {
		return err
	}

	//first put string len
	this.PutUint32(idx,(uint32)(len(value)))
	
	idx += 4
	//second put string
	copy(this.buffer[idx:],value[:len(value)])
	this.datasize += (uint32)(len(value))
	return nil
}

func (this *ByteBuffer)PutBinary(idx uint32,value []byte)(error){
	sizeneed := (uint32)(4)
	sizeneed += (uint32)(len(value))
	err := this.buffer_check(idx,sizeneed)
	if err != nil {
		return err
	}

	//first put bin len
	this.PutUint32(idx,(uint32)(len(value)))
	idx += 4
	//second put bin
	copy(this.buffer[idx:],value[:len(value)])
	this.datasize += (uint32)(len(value))
	return nil
}

func (this *ByteBuffer)PutRawBinary(value []byte)(error){
	sizeneed := (uint32)(len(value))
	err := this.buffer_check(uint32(this.datasize),sizeneed)
	if err != nil {
		return err
	}
	//second put bin
	copy(this.buffer[this.datasize:],value[:len(value)])
	this.datasize += (uint32)(len(value))
	return nil
}

func (this *ByteBuffer)Uint16(idx uint32)(ret uint16,err error){
	if idx + 2 > this.datasize {
		ret = 0
		err = ErrInvaildData
		return
	}
	ret = binary.LittleEndian.Uint16(this.buffer[idx:idx+2])
	err = nil
	return
}

func (this *ByteBuffer)Uint32(idx uint32)(ret uint32,err error){
	if idx + 4 > this.datasize {
		ret = 0
		err = ErrInvaildData
		return
	}
	ret = binary.LittleEndian.Uint32(this.buffer[idx:idx+4])
	err = nil
	return
}

func (this *ByteBuffer)String(idx uint32)(ret string,err error){
	if idx + 4 > this.datasize {
		err = ErrInvaildData
		return
	}
	//read string len
	str_len,_ := this.Uint32(idx)
	idx += 4
	if idx + str_len > this.datasize {
		err = ErrInvaildData
		return
	}
	err = nil
	ret = string(this.buffer[idx:idx+str_len])
	return
}


func (this *ByteBuffer)Binary(idx uint32)(ret []byte,err error){
	if idx + 4 > this.datasize {
		err = ErrInvaildData
		return
	}
	//read bin len
	bin_len,_ := this.Uint32(idx)
	idx += 4
	if idx + bin_len > this.datasize {
		err = ErrInvaildData
		return
	}
	err = nil
	ret = this.buffer[idx:idx+bin_len]
	return
}
