package util

type DListNode struct {
	Pre   *DListNode
	Next  *DListNode
	owner *DList
}

type DList struct {
	head  DListNode
	tail  DListNode
}

func (this *DListNode) Remove() bool {
	if nil != this.owner {		
    	this.Pre.Next = this.Next
    	this.Next.Pre = this.Pre
		this.Pre      = nil
		this.Next     = nil
		this.owner    = nil
		return true
	}
	return false
}

func NewDList() *DList {
	d := new(DList)
	d.head.Next = &d.tail
	d.tail.Pre  = &d.head
	return d
}

func (this *DList) Empty() bool {
	return this.head.Next == &this.tail
}

func (this *DList) Begin() *DListNode {
	return this.head.Next
}

func (this *DList) End() *DListNode {
	return &this.tail
}

func (this *DList) Pop() *DListNode {
	if !this.Empty() {
		n := this.head.Next
		n.Remove()
		return n
	}
	return nil
}


func (this *DList) PushBack(n *DListNode) {
	if nil != n.owner || nil != n.Pre || nil != n.Next {
		return
	}
	this.tail.Pre.Next = n
	n.Pre = this.tail.Pre
	this.tail.Pre = n
	n.Next = &this.tail
	n.owner = this
}

func (this *DList) PushFront(n *DListNode) {
	if nil != n.owner || nil != n.Pre || nil != n.Next {
		return
	}
	next := this.head.Next
	this.head.Next = n
	n.Pre  = &this.head
	n.Next = next
	next.Pre = n
	n.owner = this	
}

func (this *DList) Move(src *DList) {
	this.head.Next = src.head.Next
	this.head.Next.Pre = &this.head
	this.tail.Pre = src.tail.Pre
	this.tail.Pre.Next = &this.tail;
	//clear src
	src.head.Next = &src.tail
	src.tail.Pre  = &src.head
}