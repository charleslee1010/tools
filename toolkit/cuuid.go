package toolkit

import (
	"fmt"
	"time"
	"sync"
)

const (
	CUUID_BASE_TS = 1556596413158
)

type CUUID struct {
	sync.Mutex
	lastTs		int64
	seq 		int
}

var cuuid *CUUID = &CUUID{}

func GetCUUID() string {
	return cuuid.Get()
}

func (this*CUUID) Get() string {
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
	return fmt.Sprintf("%x-%x", ts, myseq)
}

