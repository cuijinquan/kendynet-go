package main

import util "kendynet-go/util"
import "fmt"
import "encoding/binary"
import "unsafe"

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

}





