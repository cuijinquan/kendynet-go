package util

type ListNode struct{
	next * ListNode
}


type List struct{
	head * ListNode
	tail * ListNode
	size int32
}
func NewList()(*List){
	return &List{head:nil,tail:nil,size:0}
}

func (this *List) Push(element *ListNode){
	element.next = nil
	if this.head == nil && this.tail == nil {
		this.head = element
		this.tail = element
	}else{
		this.tail.next = element
		this.tail = element
	}
	this.size += 1
}

func (this *List) Pop() (*ListNode){
	if this.head == nil {
		return nil
	}else
	{
		front := this.head
		this.head = this.head.next
		if this.head == nil {
			this.tail = nil
		}
		this.size -= 1
		return front
	}
}

func (this *List) Len()(int32){
	return this.size
}









