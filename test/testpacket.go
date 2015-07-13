package main

import packet "kendynet-go/packet"
import "fmt"
import "unsafe"

func main(){
	
	wpk := packet.NewWPacket(packet.NewByteBuffer(64))
	wpk.PutString("中国")
	wpk.PutUint32(100)
	wpk.PutString("韶关")
	
	rpk1 := (*packet.RPacket)(unsafe.Pointer(wpk.MakeRead()))

	fmt.Printf("----------------rpk1-----------------\n")

	str,_ := rpk1.String()
	fmt.Printf("%s\n",str)
	
	rpk2 := (*packet.RPacket)(unsafe.Pointer(rpk1.Clone()))

	v,_ := rpk1.Uint32()
	fmt.Printf("%d\n",v)
	
	str,_ = rpk1.String()
	fmt.Printf("%s\n",str)

	fmt.Printf("----------------rpk2-----------------\n")

	v,_ = rpk2.Uint32()
	fmt.Printf("%d\n",v)

	str,_ = rpk2.String()
	fmt.Printf("%s\n",str)

}






