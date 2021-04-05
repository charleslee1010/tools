package toolkit

import (
	"container/list"
	log "github.com/charles/mylog"
//	"sync"
//	"time"
	"errors"
)

type IndexedQueueItem struct {
	k  string
	v  interface{}  // in seconds
}

type IndexedQueue struct {
	fl *list.List

	m map[string]*list.Element

}

var IQERR_NOT_FOUND error =	errors.New("not found")
var IQERR_PANIC error =	errors.New("should never happen")

func NewIndexedQueue() *IndexedQueue {
	return &IndexedQueue{
		m:  make(map[string]*list.Element),
		fl: list.New(),
	}
}

func (this *IndexedQueue) Push(key string, value interface{}) {
	log.Info("push into fifomap:k=%s,v=%v", key, value)
	
	// if already exist, remove first and then push in

	if v, pres := this.m[key]; pres {
		this.fl.Remove(v)
	}
	
	newv := this.fl.PushBack(IndexedQueueItem{k: key, v: value})
	this.m[key] = newv
}


func (this *IndexedQueue) Find(key string) (interface{}, error) {
	log.Info("search in fifomap, key=%s", key)

	if v, pres := this.m[key]; pres {
		if ci, ok := v.Value.(IndexedQueueItem); ok {
			return ci.v, nil
		}
	}
	return nil, IQERR_NOT_FOUND
}


//
// return key, value of the front element
// error is used to indicate empty queue
// 
func (this *IndexedQueue) Peek() (string, interface{}, error) {
	le := this.fl.Front()
	if le == nil {
		return "", nil, IQERR_NOT_FOUND
	}
	
	if iqi, ok := le.Value.(IndexedQueueItem); ok {
		return iqi.k, iqi.v, nil
	}
	
	return 	"", nil, IQERR_PANIC
}


func (this *IndexedQueue) Remove(key string) {

	if le, pres := this.m[key]; pres {
		this.fl.Remove(le)
		delete(this.m, key)
	}
}

func (this *IndexedQueue) Dump() {
	log.Info("dump fifo map")

	log.Info("link list")

	for v := this.fl.Front(); v != nil; v = v.Next() {
		if ci, ok := v.Value.(IndexedQueueItem); ok {
			log.Info("CI:%+v", ci)
		} else {
			log.Info("invalid link element")
		}
	}

	log.Info("element map")

	for k, v := range this.m {
		if ci, ok := v.Value.(IndexedQueueItem); ok {
			log.Info("key:%s, CI:%+v", k, ci)
		} else {
			log.Info("invalid map element")
		}
	}

}


func (this *IndexedQueue) Check() bool {
	log.Info("check fifo map")


	ok := true
	
	// check if all map item exists in queue
	for _, v := range this.m {
		found := false
		for ele := this.fl.Front(); ele != nil; ele = ele.Next() {
			if ele == v {
				found = true
				break
			}
		}
		
		if !found {
			log.Error("element no in queue")
			ok = false
			break
		}
	}
	
	if !ok {
		return ok
	}

	for ele := this.fl.Front(); ele != nil; ele = ele.Next() {	
		found := false

		for _, v := range this.m {			
			if ele == v {
				found = true
				break
			}
		}
		
		if !found {
			log.Error("element not in map")
			ok = false
			break
		}
	}

	return ok
}
