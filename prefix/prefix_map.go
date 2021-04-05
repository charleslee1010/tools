package toolkit

import (
	"fmt"
	log "github.com/charles/mylog"
	"sync"
)

const (
	LEN_PREFIX_MIN = 1
	LEN_PREFIX_MAX = 12
)

// length range for a prefix
type LengthRange struct {
	min int
	max int
}

// key: number prefix
type Prefix struct {
	plist []map[string]LengthRange
}

// data format:
// key:
//    appkey + tela/telb + b/w

type PrefixMap struct {
	//	total   int
	//	ignored int
	l sync.RWMutex

	m map[string]*Prefix
}

func (this *PrefixMap) Dump() {
	this.l.RLock()
	defer this.l.RUnlock()

	log.Info("number of keys=%d", len(this.m))

	for k, v := range this.m {
		log.Info("key=%s", k)
		v.Dump()
	}
}

// check if specified list exists
func (this *PrefixMap) Exists(key string) bool {
	this.l.RLock()
	defer this.l.RUnlock()

	if _, pres := this.m[key]; pres {
		return true
	} else {
		return false
	}
}

func (this *PrefixMap) Find(key string, x string) bool {
	this.l.RLock()
	defer this.l.RUnlock()

	if p, pres := this.m[key]; pres {
		return p.Find(x)
	} else {
		return false
	}
}

func NewPrefix() *Prefix {
	return &Prefix{plist: make([]map[string]LengthRange, LEN_PREFIX_MAX+1)}
}

func (this *Prefix) Dump() {
	log.Info("dump prefix")

	for i := LEN_PREFIX_MIN; i <= LEN_PREFIX_MAX; i++ {
		if this.plist[i] != nil {
			log.Info("len=%d, %+v", i, this.plist[i])
		}

	}
}

func (this *Prefix) Find(x string) bool {

	l := len(x)

	if l == 0 || this.plist == nil || len(this.plist) == 0 {
		return false
	}

	for i := LEN_PREFIX_MIN; i <= LEN_PREFIX_MAX && len(x) >= i; i++ {
		if mseg := this.plist[i]; mseg != nil {
			if v, pres := mseg[x[0:i]]; pres {
				if v.max == 0 || v.min == 0 || (v.max >= l && v.min <= l) {
					log.Debug("Found prefix, x=%s", x)
					return true
				} else {
					log.Debug("prefix not found, x=%s", x)
					return false
				}
			}
		}
	}
	return false
}
func (this *Prefix) Add(prefix string, min, max int) error {
	l := len(prefix)

	if l < LEN_PREFIX_MIN || l > LEN_PREFIX_MAX {
		err := fmt.Errorf("Invalid prefix min:%d, max %d", min, max)
		return err
	}

	if this.plist[l] == nil {
		this.plist[l] = make(map[string]LengthRange)
	}

	this.plist[l][prefix] = LengthRange{min: min, max: max}
	return nil
}

func (this *PrefixMap) Add(key, prefix string, min, max int) error {

	this.l.Lock()
	defer this.l.Unlock()

	if this.m == nil {
		this.m = make(map[string]*Prefix)
	}

	p, pres := this.m[key]
	if !pres {
		p = NewPrefix()
		this.m[key] = p
	}

	// put prefix to p
	return p.Add(prefix, min, max)
}

func (this *PrefixMap) Copy(src *PrefixMap) {

	this.l.Lock()
	defer this.l.Unlock()
	this.m = src.m
}
