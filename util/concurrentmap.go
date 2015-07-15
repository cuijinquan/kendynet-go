//基于channel同步的一个concurrent map
package util

type op interface {
	Do(map[interface{}]interface{})
}

type opSet struct {
	key     interface{}
	val     interface{}
}

func (this *opSet) Do(m map[interface{}]interface{}) {
	m[this.key] = this.val
}

type opGet struct {
	key     interface{}
	valOk   bool
	ret		chan interface{}
}

func (this *opGet) Do(m map[interface{}]interface{}) {
	var val interface{}
	val,this.valOk = m[this.key] 
	this.ret <- val
}

type opDel struct {
	key     interface{}
}

func (this *opDel) Do(m map[interface{}]interface{}) {
	delete(m,this.key)
}

type opAll struct {
	ret		chan interface{}
}

func (this *opAll) Do(m map[interface{}]interface{}) {
	_m := make(map[interface{}]interface{})
	for k,v := range(m) {
		_m[k] = v
	}
	this.ret <- _m
}

type opKeys struct {
	ret chan interface{}
}

func (this *opKeys) Do(m map[interface{}]interface{}) {
	keys := make([]interface{},len(m))
	i := 0
	for k,_ := range(m) {
		keys[i] = k
		i++ 
	}
	this.ret <- keys
}

type opVals struct {
	ret chan interface{}
}

func (this *opVals) Do(m map[interface{}]interface{}) {
	vals := make([]interface{},len(m))
	i := 0
	for _,v := range(m) {
		vals[i] = v
		i++
	}
	this.ret <- vals
}


type ConnCurrMap struct {
	opChan  chan op
}


func NewConnCurrMap() * ConnCurrMap {
	m := new(ConnCurrMap)
	m.opChan = make(chan op,1024)
	_map    := make(map[interface{}]interface{})
	go func () {
		for{
			_op,ok := <- m.opChan
			if !ok {
				break
			}
			_op.Do(_map)
		}
	}()
	return m
}

func (this *ConnCurrMap) Set(key interface{},val interface{}) {
	o := new(opSet)
	o.key = key
	o.val = val
	this.opChan <- o
}

func (this *ConnCurrMap) Del(key interface{}) {
	o := new(opDel)
	o.key = key
	this.opChan <- o
}

func (this *ConnCurrMap) Get(key interface{}) (interface{},bool) {
	o := new(opGet)
	o.key = key
	o.ret = make(chan interface{})
	this.opChan <- o
	ret,ok := <- o.ret
	close(o.ret)
	if !ok {
		return nil,false
	}
	return ret,o.valOk
}

func (this *ConnCurrMap) All() map[interface{}]interface{} {
	o := new(opAll)
	o.ret = make(chan interface{})
	this.opChan <- o
	ret,ok := <- o.ret
	close(o.ret)
	if !ok {
		return nil
	}
	return ret.(map[interface{}]interface{})	
}

func (this *ConnCurrMap) Keys() []interface{} {
	o := new(opKeys)
	o.ret = make(chan interface{})
	this.opChan <- o
	ret,ok := <- o.ret
	close(o.ret)
	if !ok {
		return nil
	}
	return ret.([]interface{})		
} 


func (this *ConnCurrMap) Vals() []interface{} {
	o := new(opVals)
	o.ret = make(chan interface{})
	this.opChan <- o
	ret,ok := <- o.ret
	close(o.ret)
	if !ok {
		return nil
	}
	return ret.([]interface{})		
}   

