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

func cal_remain(now int64,expire int64) int64 {
	return expire - now
}

func (this *WheelMgr) reg(t *Timer,tick int64) {
	var slot,wsize uint32
	var w *wheel
	wtype  := wheel_sec
	remain := cal_remain(tick,t.expire)
	for {
		w = this.wheels[wtype]
		wsize = wheelsize(byte(wtype))
		if wtype == wheel_day || int64(wsize) >= remain {
			slot = w.cur + uint32(remain)
			if slot >= wsize {
				slot = slot - wsize
			}
			w.items[slot].PushBack(cast2DListNode(t))

			break
		} else {
			remain--
			remain = remain / int64(wsize)
			wtype++
		}
	}
}

func (this *WheelMgr) fire(w *wheel,tick int64) {
	w.cur++
	if w.cur == wheelsize(w.tt) {
		w.cur = 0
	}
	if !w.items[w.cur].Empty() {
		tlist := NewDList()
		tlist.Move(w.items[w.cur])
		if w.tt == wheel_sec {
			for {
				t := cast2Timer(tlist.Pop())
				if t == nil {
					break
				}
				t.incb = true
				ret := t.callback(tick)
				t.incb = false
				if ret >= 0 && t.timeout > 0 {
					if ret > 0 {
						t.timeout = ret
					}
					t.expire = tick + t.timeout
					this.reg(t,tick)
				}
			}
		} else {
			for {
				t := cast2Timer(tlist.Pop())
				if t == nil {
					break
				}
				this.reg(t,tick);
			}
		}
	}

	if w.cur + 1 == wheelsize(w.tt) && w.tt < wheel_day {
		this.fire(this.wheels[w.tt + 1],tick)
	}
}


func (this *WheelMgr) Tick(now int64) {
	for{
		if this.lasttime == now {
			break
		}
		this.lasttime++
		this.fire(this.wheels[wheel_sec],this.lasttime)
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
	}else {
		t.timeout = timeout
	}
	t.expire = now + t.timeout
	if this.lasttime == 0 {
		this.lasttime = now
	}
	this.reg(t,this.lasttime)
	return t
}

func (this *Timer) UnRegister() {
	this.timeout = 0
	if !this.incb {
		cast2DListNode(this).Remove()
	}
}
