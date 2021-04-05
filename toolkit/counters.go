package toolkit

import (
	"sync"
)

type Counter struct {
	V int
}

type Counters struct {
	l sync.Mutex
	m map[string]*Counter
}

func (this *Counters) Add(key string, v int) {
	this.l.Lock()
	defer this.l.Unlock()

	c, _ := this.m[key]
	if c == nil {
		// add a new item
		c = &Counter{}
		this.m[key] = c
	}

	// now we can add
	c.V += v
}

func (this *Counters) CopyAndReset() map[string]*Counter {
	this.l.Lock()
	defer this.l.Unlock()

	c := this.m
	this.m = make(map[string]*Counter)

	return c
}
