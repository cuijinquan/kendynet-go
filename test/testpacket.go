package main

import packet "kendynet-go/packet"
import "fmt"
import "unsafe"

func main(){
	
	wpk := packet.NewWPacket(packet.NewByteBuffer(64))
	wpk.PutString("中国")
	wpk.PutUint32(100)
	wpk.PutString("韶关")
	
	rpk := (*packet.RPacket)(unsafe.Pointer(wpk.MakeRead()))
	str,_ := rpk.String()
	fmt.Printf("%s\n",str)
	
	v,_ := rpk.Uint32()
	fmt.Printf("%d\n",v)
	
	str,_ = rpk.String()
	fmt.Printf("%s\n",str)
	if str != "韶关" {
		fmt.Printf("bad\n")
	}
}






