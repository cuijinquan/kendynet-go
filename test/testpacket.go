package main

import packet "kendynet-go/packet"
import "fmt"
//import "encoding/binary"

func main(){
	
	wpk := packet.NewWpacket(packet.NewByteBuffer(64),false)
	wpk.PutString("中国")
	wpk.PutUint32(100)
	wpk.PutString("韶关")
	rpk := packet.NewRpacket(wpk.Buffer(),false)
	str,_ := rpk.String()
	fmt.Printf("%s\n",str)
	
	v,_ := rpk.Uint32()
	fmt.Printf("%d\n",v)
	
	str,_ = rpk.String()
	fmt.Printf("%s\n",str)
	if str != "韶关" {
		fmt.Printf("bad\n")
	}
	
	
	//fmt.Printf("buf len %d\n",uint32(wpk.Buffer().Len()))
	//packet_size :=	binary.LittleEndian.Uint32(wpk.Buffer().Bytes()[0:4])
	//fmt.Printf("packet_size : %d\n",packet_size)
	
/*
	rpk := packet.NewRpacket(wpk.Buffer(),false)
	value,_ := rpk.Uint32()
	fmt.Printf("%d\n",value)
	value,_ = rpk.Uint32()
	fmt.Printf("%d\n",value)
	value,_ = rpk.Uint32()
	fmt.Printf("%d\n",value)

	str,_ := rpk.String()
	fmt.Printf("%s\n",str)*/
}

/*
type Node struct{
	util.ListNode
	Value int32
}

func (this *Node) Cast2ListNode()(*util.ListNode){
	return	(*util.ListNode)(unsafe.Pointer(this))
}

func Cast2Node(n *util.ListNode)(*Node){
	return ((*Node)(unsafe.Pointer(n)))
}


func main(){

	wpos := 0
	rpos := 0

	buffer := make([]byte,100)
	binary.LittleEndian.PutUint32(buffer[wpos:wpos+4], uint32(100))
	wpos += 4
	binary.LittleEndian.PutUint32(buffer[wpos:wpos+4], uint32(101))
	wpos += 4

	fmt.Printf("%d\n",binary.LittleEndian.Uint32(buffer[rpos:rpos+4]))
	rpos += 4

	fmt.Printf("%d\n",binary.LittleEndian.Uint32(buffer[rpos:rpos+4]))
	rpos += 4

}*/





