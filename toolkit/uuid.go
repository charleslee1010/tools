package toolkit

import (
	"fmt"
	"time"
	"sync"
)

const (
	BASE_TS = 1556596413158
)

type UUIDStruct struct {
	sync.Mutex
	lastTs		int64
	seq 		int
	node		string
}
//
//func NewUUID() *UUIDStruct {
//	return  &UUIDStruct{}
//}

func NewUUID(node string) *UUIDStruct {
	return  &UUIDStruct{
		node :node,
	}
}

func (this*UUIDStruct) Get() string {
	ts := time.Now().UnixNano()/1e6 - BASE_TS
	
	this.Lock()
	
	if ts == this.lastTs {
		// use seq
		this.seq ++
	} else {
		// no seq
		this.lastTs = ts
		this.seq = 0
	}
	
	myseq := this.seq
	
	this.Unlock()
	
	// format ts and myseq
	return fmt.Sprintf("%v%x-%x", this.node, ts, myseq)
}

