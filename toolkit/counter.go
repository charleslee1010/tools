package toolkit

import (
"bytes"
	"fmt"
)



type StringCounter struct {
	m map[string]int
}

func NewStringCounter() *StringCounter {
	return &StringCounter{
		m: make(map[string]int),
	}
}

func (c *StringCounter) Add(key string, val int) {
	if v, pres := c.m[key]; pres {
		c.m[key] = v + val
	} else {
		c.m[key] = val
	}
}

func (c *StringCounter) Get(key string) (v int, pres bool){
	v, pres = c.m[key]
	return
}

func (c *StringCounter) String() string {
	var b bytes.Buffer
	for k, v := range c.m {
		b.WriteString(fmt.Sprintln(k, v))
	}
	return b.String()
}

func (c *StringCounter) Traverse(a func (string, int) bool) bool {
	
	for k, v := range c.m {
		if !a(k, v) {
			return false
		} 
	}
	return true
}

type IntegerCounter struct {
	m map[int]int
}

func NewIntegerCounter() *IntegerCounter {
	return &IntegerCounter{
		m: make(map[int]int),
	}
}

func (c *IntegerCounter) Add(key int, val int) {
	if v, pres := c.m[key]; pres {
		c.m[key] = v + val
	} else {
		c.m[key] = val
	}
}



func (c *IntegerCounter) String() string {
	var b bytes.Buffer
	for k, v := range c.m {
		b.WriteString(fmt.Sprintln(k, v))
	}
	return b.String()
}

type MuliCountStruct struct {
	V1 	int
	V2	int
	V3 	int
}

type MultiCounter struct {
	m map[string]*MuliCountStruct
}

func NewMultiCounter() *MultiCounter {
	return &MultiCounter{
		m: make(map[string]*MuliCountStruct),
	}
}

func (c *MultiCounter) Add(key string, v1,v2,v3 int) {
	if v, pres := c.m[key]; pres {
		//c.m[key] = v + val
		v.V1 += v1
		v.V2 += v2
		v.V3 += v3
		
	} else {
		
		c.m[key] = &MuliCountStruct{V1:v1, V2:v2, V3:v3,}
	}
}
//
//func (c *MultiCounter) Get(key string) (v int, pres bool){
//	v, pres = c.m[key]
//	return
//}
//
//func (c *MultiCounter) String() string {
//	var b bytes.Buffer
//	for k, v := range c.m {
//		b.WriteString(fmt.Sprintln(k, v))
//	}
//	return b.String()
//}

func (c *MultiCounter) Traverse(a func (string, int, int, int) bool) bool {
	
	for k, v := range c.m {
		if !a(k, v.V1, v.V2, v.V3) {
			return false
		} 
	}
	return true
}
