package toolkit

import (
	"fmt"
	"testing"
	"time"
)


func TestMain(m *testing.M) {
	
    m.Run()
    time.Sleep(2*time.Second)
}


func TestFifoMap(t *testing.T) {
	
	fm := NewIndexedQueue()
	
	
	fm.Push("a", "2020-01-20 12:21:12")
	fm.Push("a1", "2020-01-20 12:22:12")
	fm.Push("a2", "2020-01-20 12:22:12")
	fm.Push("a3", "2020-01-20 12:23:12")
	fm.Push("a4", "2020-01-20 12:24:12")

	if fm.Check() == false {
		t.Error("invalid struct")
	}

	// find	
	if d, e := fm.Find("a2"); e != nil {
		t.Error("not found")
	} else {
		if s, ok := d.(string); ok && s == "2020-01-20 12:22:12" {
			
		} else {
			t.Error("invalid value")
		}
	}


	// peek
	if k, d, e := fm.Peek(); e != nil || k != "a" {
		t.Error("peek not found or invalid element")
	} else {
		if s, ok := d.(string); ok && s == "2020-01-20 12:21:12" {
			
		} else {
			t.Error("invalid value of a")
		}
	}
	
	
	fm.Dump()
	
}

