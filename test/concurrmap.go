package main

import (
	Util "kendynet-go/util"
	"fmt"
)


func main() {
	cM := Util.NewConnCurrMap()
	cM.Set("a","a")
	cM.Set("b","b")
	cM.Set("c","c")
	for k,v := range cM.All() {
		fmt.Printf("%s,%s\n",k.(string),v.(string))
	}
	_,ok := cM.Get("d")
	if !ok {
		fmt.Printf("no d\n")
	}
}