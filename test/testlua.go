package main

import "golua/lua"
import "fmt"


func main() {
	L := lua.NewState()
	defer L.Close()
	L.OpenLibs()
	L.DoFile("test.lua")
	L.GetGlobal("language")
	size := int(L.ObjLen(-1))	
	for i := 1; i <= size; i++ {
		L.RawGeti(-1-i+1,i)
		fmt.Printf("%s\n",L.ToString(-1))
	}
	
	L.GetGlobal("hello")
	L.Call(0,0)	
}
