package main

import "fmt"
import util "kendynet-go/util"
//import "unsafe"
//import "reflect"


type a interface {
	funca ()
}

type b struct {}

func (this *b) funca(){
	fmt.Printf("funca\n")
}

type c struct {
	b
}

func main() {
	var aa a
	aa = new(c)
	aa.funca()
	fmt.Printf("%ld\n",util.SystemMs())
}


/*type STest struct {
	A byte
	B byte
	C byte
}

func ByteToStruct(b []byte) uintptr{
	pbytes := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	return pbytes.Data
}

func main() {
	b := []byte{'a','b','c'}
	var s STest = *(*STest)(unsafe.Pointer(ByteToStruct(b)))
	fmt.Println(s)
}*/