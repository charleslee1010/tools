package toolkit

import (
	log "github.com/charles/mylog"
	"sync"
	"time"
)

type CacheItem struct {
	// timestamp of the cache item
	ts          int
	
	// data to be cached 
	data        interface{}
}

type TimedCache struct {
	
	m *IndexedQueue
	
	// expiration time in seconds for each item
	exp int

	l sync.RWMutex
}

func NewTimedCache(exp int) *TimedCache {
	c := &TimedCache{
		m:  NewIndexedQueue(),
		exp: exp,
	}
	
	// start routine to cleanup the expired item
	go cleanup(c)
	
	return c
}

//  
// push an item into the timed map with the current timestamp
// 
func (this *TimedCache) Push(key string, ts int, data interface{}) {
	
	// if already exist, remove first and then push in
	if ts == 0 {
		ts = int(time.Now().Unix())
	}
	
	this.l.Lock()
	defer this.l.Unlock()
	
	this.m.Push(key, CacheItem{ts:ts, data:data})
}

//
// return the ts stored in fifomap, 0 if not found
//
func (this *TimedCache) Find(key string) (int, interface{}, error) {

	this.l.RLock()
	defer this.l.RUnlock()
	
	v, err := this.m.Find(key)
	if err == nil {
		if vc, ok := v.(CacheItem); ok {
			return vc.ts, vc.data, nil
		} else {
			return 0, nil, IQERR_PANIC 
		}
	} else {
		return 0, nil, err
	}
}

func (this *TimedCache) Peek() (string, int, interface{}, error) {

	this.l.RLock()
	defer this.l.RUnlock()
	
	k, v, err := this.m.Peek()
	if err == nil {
		if vc, ok := v.(CacheItem); ok {
			log.Debug("peek key=%+v", vc)
			
			return k, vc.ts, vc.data, nil
		} else {
			return "", 0, nil, IQERR_PANIC 
		}
	} else {
		return "", 0, nil, err
	}
}


func (this *TimedCache) Remove(key string) {

	this.l.Lock()
	defer this.l.Unlock()
	
	log.Debug("remove key=%s", key)
	
	this.m.Remove(key)
}


const INTERVAL = time.Duration(60)


func cleanup(c *TimedCache) {
	for true {
//		log.Debug("wait for next peek")
		time.Sleep(INTERVAL *time.Second)
		
		ct := int(time.Now().Unix())
		log.Debug("start peek ct=%d", ct)
		
		for true {
			k, ts, _, err := c.Peek()
			if err != nil {
				log.Debug("stop peek, return error:%v", err)
				break
			}
			if ct >= ts + c.exp {
				// expired
				c.Remove(k)
			} else {
				log.Debug("stop peek, not expired:ct=%d, c.exp=%d, ts=%d", ct, c.exp, ts)
				break
			}			
		}
	}
}