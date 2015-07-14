package util

const (
	SET = 1
	GET = 2
	DEL = 3
	ALL = 4
)

type op struct {
	tt      byte
	key     interface{}
	val     interface{}
	valOk   bool
	retChan chan *op
	all     map[interface{}]interface{}
}


type ConnCurrMap struct {
	opChan  chan *op
	m       map[interface{}]interface{}
}


func NewConnCurrMap() * ConnCurrMap {
	m := new(ConnCurrMap)
	m.opChan = make(chan *op,1024)
	_map    := make(map[interface{}]interface{})
	go func () {
		for{
			_op,ok := <- m.opChan
			if !ok {
				break
			}
			if _op.tt == GET {
				_op.val,_op.valOk = _map[_op.key]
				_op.retChan <- _op
			}else if _op.tt == ALL {
				_op.all = make(map[interface{}]interface{})
				for k,v := range _map {
					_op.all[k] = v
				}
				_op.retChan <- _op
			}else if _op.tt == DEL{
				delete(_map,_op.key)
			}else{
				_map[_op.key] = _op.val
			}
		}
	}()
	return m
}

func (this *ConnCurrMap) Set(key interface{},val interface{}) {
	o := new(op)
	o.tt  = SET
	o.key = key
	o.val = val
	this.opChan <- o
}

func (this *ConnCurrMap) Del(key interface{}) {
	o := new(op)
	o.tt  = DEL
	o.key = key
	this.opChan <- o
}

func (this *ConnCurrMap) Get(key interface{}) (interface{},bool) {
	o := new(op)
	o.tt  = GET
	o.key = key
	o.retChan = make(chan *op)
	this.opChan <- o
	ret,ok := <- o.retChan
	close(o.retChan)
	if !ok {
		return nil,false
	}
	return ret.val,ret.valOk
}

func (this *ConnCurrMap) All() map[interface{}]interface{} {
	o := new(op)
	o.tt  = ALL
	o.retChan = make(chan *op)
	this.opChan <- o
	ret,ok := <- o.retChan
	close(o.retChan)
	if !ok {
		return nil
	}
	return ret.all	
} 


