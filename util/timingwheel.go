package util

import "unsafe"

const (
	wheel_sec   = 0  
	wheel_hour  = 1     
	wheel_day   = 2
	MAX_TIMEOUT = (1000*3600*24-1)
)

func wheelsize(t byte) uint32 {
	if t == wheel_sec {
		return 1000
	}else if t == wheel_hour {
		return 3600
	}else if t == wheel_day {
		return 24
	}else {
		return 0
	}
}

func precision(t byte) uint64 {
	if t == wheel_sec {
		return 1
	}else if t == wheel_hour {
		return 1000
	}else if t == wheel_day {
		return 3600
	}else {
		return 0
	}	
}


type wheel struct {
	tt    byte
	cur   uint32
	items []*DList
}

type Timer struct {
	DListNode
	timeout  int64
	expire   int64
	incb     bool
	callback func (int64) int64 //返回0,则timeout不变,大于0,则timeout为返回值,小于0从定时器移除
}


func cast2DListNode(t *Timer)(*DListNode){
	return	(*DListNode)(unsafe.Pointer(t))
}

func cast2Timer(n *DListNode)(*Timer){
	return ((*Timer)(unsafe.Pointer(n)))
}

func (this *Timer) UnRegister() {
	if !this.incb {
		cast2DListNode(this).Remove()
	}
}
 
type WheelMgr struct {
	wheels 	 [wheel_day+1]*wheel;
	lasttime int64
}

func newWheel(tt byte) *wheel {
	w      := new(wheel)
	w.tt    = tt
	w.cur   = 0
	w.items = make([]*DList,wheelsize(tt))
	for i:= uint32(0); i < wheelsize(tt); i++ {
		w.items[i] = NewDList()
	}
	return w
}

func TimingWheel() *WheelMgr {
	t := new(WheelMgr)
	for i:= uint32(0);i < wheel_day + 1; i++ {
		t.wheels[i] = newWheel(byte(i))
	}
	return t
}


func (this *WheelMgr) add2Wheel(w *wheel,t *Timer,remain int64) {
	var i uint32
	slots := wheelsize(w.tt) - w.cur
	if w.tt == wheel_day || slots > uint32(remain) {
		i = (w.cur + uint32(remain))%(wheelsize(w.tt))
		w.items[i].PushBack(cast2DListNode(t))			
	}else {
		remain -= int64(slots)
		remain /= int64(wheelsize(w.tt))
		this.add2Wheel(this.wheels[w.tt+1],t,remain)		
	}
}

func (this *WheelMgr) reg(t *Timer,w *wheel,tick int64) {
	if t.expire > tick {
		if w == nil {
			this.add2Wheel(this.wheels[wheel_sec],t,t.expire - tick)
		}else{
			this.add2Wheel(w,t,t.expire - tick)
		}
	}
}

//将本级超时的定时器推到下级时间轮中
func (this *WheelMgr) down(t *Timer,w *wheel,tick int64) {
	var remain int64
	if t.expire >= tick {
		remain = (t.expire - tick) - int64(wheelsize(w.tt-1))
		remain /= int64(precision(w.tt))
		w.items[w.cur + uint32(remain)].PushBack(cast2DListNode(t))		
	}	
}

//处理上一级时间轮
func (this *WheelMgr) tickup(w *wheel,tick int64) {
	var t *Timer
	items := w.items[w.cur]
	for{
		t = cast2Timer(items.Pop())
		if t == nil {
			break
		}
		this.down(t,this.wheels[w.tt-1],tick)
	}
	w.cur = (w.cur+1)%wheelsize(w.tt)
	if w.cur == 0 && w.tt != wheel_day {
		this.tickup(this.wheels[w.tt+1],tick)
	}	
}

func (this *WheelMgr) fire(tick int64) {
	w := this.wheels[wheel_sec]
	w.cur = (w.cur+1)%wheelsize(wheel_sec)
	if w.cur == 0 {
		this.tickup(this.wheels[wheel_hour],tick)
	}
	items := w.items[w.cur]
	for {
		t := cast2Timer(items.Pop())
		if t == nil {
			break
		}
		t.incb = true
		ret := t.callback(tick)
		t.incb = false
		if ret >= 0 {
			if ret > 0 {
				t.timeout = ret
			}
			t.expire = tick + t.timeout
			this.reg(t,nil,tick)
		}
	}
}

func (this *WheelMgr) Tick(now int64) {
	for{
		if this.lasttime == now {
			break
		}
		this.fire(this.lasttime)
		this.lasttime++
	}
} 

func (this *WheelMgr) Register(timeout int64,now int64,
							   callback func (int64) int64) *Timer {
	if timeout == 0 || callback == nil {
		return nil
	}
	if now == 0 {
		now = SystemMs()
	}
	t := new(Timer)
	t.callback = callback
	if timeout > MAX_TIMEOUT {
		t.timeout = MAX_TIMEOUT
	}else{
		t.timeout = t.timeout
	}
	t.expire = now + t.timeout
	if this.lasttime == 0 {
		this.lasttime = now
	}
	this.reg(t,nil,this.lasttime)
	return t
}
