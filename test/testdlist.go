package main

import util "kendynet-go/util"
import "fmt"
import "unsafe"

type Node struct{
	util.DListNode
	Value int32
}

func (this *Node) Cast2DListNode()(*util.DListNode){
	return	(*util.DListNode)(unsafe.Pointer(this))
}

func Cast2Node(n *util.DListNode)(*Node){
	return ((*Node)(unsafe.Pointer(n)))
}


func main(){
	
	dlist := util.NewDList()

	if dlist.Empty() {
		fmt.Printf("dlist Empty\n")
	}

	dlist.PushBack((&Node{Value:100}).Cast2DListNode())
	dlist.PushBack((&Node{Value:101}).Cast2DListNode())
	dlist.PushBack((&Node{Value:104}).Cast2DListNode())
	dlist.PushBack((&Node{Value:105}).Cast2DListNode())
	dlist.PushBack((&Node{Value:106}).Cast2DListNode())

	n := dlist.Begin()
	for {
		if n == dlist.End() {
			break
		}
		fmt.Printf("%d\n",Cast2Node(n).Value)
		if Cast2Node(n).Value == 104 {
			n.Remove()
			break
		}else{
			n = n.Next
		}
	}

	fmt.Printf("------------after remove 104-------\n")
	n = dlist.Begin()
	for {
		if n == dlist.End() {
			break
		}
		fmt.Printf("%d\n",Cast2Node(n).Value)
		n = n.Next
	}	
}





