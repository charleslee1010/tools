package toolkit

import (
	"sync"
)

type SyncMap struct {
	sync.RWMutex
	M map[string]interface{}
}

func NewSyncMap() (*SyncMap) {
	return &SyncMap {
		M: make(map[string]interface{}),
	}
}

func (sm*SyncMap)Clear() {
	sm.Lock()
	sm.M = make(map[string]interface{})
	sm.Unlock()
}
func (sm*SyncMap)Load(key string ) (v interface{}, pres bool) {
	sm.RLock()
	v, pres = sm.M[key]
	sm.RUnlock()
	return
}

func (sm*SyncMap)Store(key string , value interface{}) {
	sm.Lock()
	sm.M[key] = value
	sm.Unlock()
}

func (sm*SyncMap)Set(s map[string]interface{}) {
	sm.Lock()
	sm.M = s
	sm.Unlock()
}

func (sm*SyncMap)Range(f (func (key string, value interface{}) bool)) {
	sm.RLock()
	for k, v := range sm.M {
		if f(k,v) == false {
			break
		}
	}
	
	sm.RUnlock()
}
