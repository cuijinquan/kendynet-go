package util
import "time"
	
type DurationTick struct{
	stop	    bool
	start       bool
}

func DurationTicker() (*DurationTick) {
	return &DurationTick{stop:false,start:false}
}

func (this *DurationTick) Start(ms uint64,callback func (time.Time))(bool) {
	if this.start {
		return false
	}
	this.start = true
	this.stop  = false
	tickchan := time.Tick(time.Millisecond * time.Duration(ms))
	go func(){
		var t time.Time
		for{ 
			t = <- tickchan
			if this.stop {
				break
			}
			callback(t)
		}
	}()
	return true
}

func (this *DurationTick) Stop() {
	this.stop = true
}
