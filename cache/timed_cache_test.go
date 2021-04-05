package toolkit

import (

	"testing"
	"time"
)





func TestTimedCache(t *testing.T) {
	
	fm := NewTimedCache(50)
	
	ts := int(time.Now().Unix())
	
	var data interface{} = "b2"
	
	fm.Push("a", ts, "b")
	fm.Push("a1", ts+1, "b1")
	fm.Push("a2", ts+2, data)
	fm.Push("a3", ts+3, "b3")
	fm.Push("a4", ts+4, "b4")
	
	if ct, d, e := fm.Find("a2"); e != nil || ct != ts + 2 || d != data {
		t.Error("invalid x2")
	} 
	
	// wait for 60 seconds for cleanup 
	time.Sleep(70*time.Second)
	
	if _, _, e := fm.Find("a2"); e == nil {
		t.Error("invalid a2, should not be found")
	} else {
		t.Log(e)
	}
	
	time.Sleep(5*time.Second)
}

