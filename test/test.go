package main

import "fmt"
import "unsafe"
import "reflect"


type STest struct {
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
}