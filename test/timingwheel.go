package main

import(
	"fmt"
	"time"
	util   "kendynet-go/util"	
)

func main(){
	timingwheel := util.TimingWheel()
	timingwheel.Register(100,0,func (_ int64) int64 {
		fmt.Printf("timeout\n")
		return 0
	})
	ticker := util.DurationTicker()
	ticker.Start(1,func (_ time.Time){
		timingwheel.Tick(util.SystemMs())
	})
	for{
		time.Sleep(10000000)
	}
}
