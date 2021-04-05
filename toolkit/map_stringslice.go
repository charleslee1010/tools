package toolkit

import (
//"bytes"
//	"fmt"
)



type MapStringSlice struct {
	m map[string][]string
}

func NewMapStringSlice() *MapStringSlice {
	return &MapStringSlice{
		m: make(map[string][]string),
	}
}

func (c *MapStringSlice) Add(key string, val string) {
	if v, pres := c.m[key]; pres {
		c.m[key] = append(v, val)
	} else {
		c.m[key] = []string{val}
	}
}

func (c *MapStringSlice) Traverse(a func (string, []string) bool) bool {
	
	for k, v := range c.m {
		if !a(k, v) {
			return false
		} 
	}
	return true
}

