package main

import util "kendynet-go/util"
import "fmt"
import "unsafe"
import "encoding/binary"

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
	list := util.NewList()

	list.Push((&Node{Value:100}).Cast2ListNode())
	list.Push((&Node{Value:101}).Cast2ListNode())

	fmt.Printf("%d\n",list.Len())

	fmt.Printf("%d\n",Cast2Node(list.Pop()).Value)
	fmt.Printf("%d\n",Cast2Node(list.Pop()).Value)

	fmt.Printf("%d\n",list.Len())


	header := uint32(0)
	ptr    := ([]byte)(unsafe.Pointer(&header))
	binary.LittleEndian.PutUint16(ptr[:],100)

	fmt.Printf("%d\n",header)

}





