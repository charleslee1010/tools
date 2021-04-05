package toolkit

import (
	"time"
)

type TimerThread struct {
	flushInterval int64
	r		func ()	
} 

//const DEFAULT_FLUSH_INTERVAL = 300     // 300 seconds


func NewTimerThread(f int, r func ()) (*TimerThread) {
	t := &TimerThread{
//		  flushInterval : DEFAULT_FLUSH_INTERVAL,
		  r : r,
	}
		
	if f >= 0 && f <= 10 {
		if f == 0 {
			t.flushInterval = 0
		} else if 60%f == 0 {
			t.flushInterval = int64(f*60)           // minute to seconds
		} else {
			// invalid interval
			return nil
		}	
	} else {
		return nil
	}
	
	return t
}


//func (this*TimerThread) GetFlushInterval() int64 {
//	return this.flushInterval
//}

func (this*TimerThread) Run() {
	timer1 := &time.Timer{}
	
	if  this.flushInterval != 0 {
		r := this.getNext() * time.Second
	    timer1 = time.NewTimer(r)
	}
	
	for true {
		select {
		case <-timer1.C:
		    // flush stat to db
		    this.r()
		    // start next flush timer
		    timer1 = time.NewTimer(this.getNext() * time.Second)
		}
	}
}

func (this*TimerThread) getNext() time.Duration {
	// seconds to next ticker
	t := time.Now().Unix()%3600
	t = (t+this.flushInterval)/this.flushInterval*this.flushInterval - t
	if t <= 1 {
		return time.Duration(this.flushInterval+t)
	}
	return time.Duration(t)
}
//
//func (this*TimerThread) NormalizeTs(ts int64) int64 {
//	// seconds to start time of this stat slot
//		
//	t := ts%3600
//	
//	t = t/this.flushInterval*this.flushInterval
//	
//	
//	return ts/3600*3600+t
//}
//
//func (this*TimerThread) NormalizeTsStr(tstr string, fmt string) (string, error) {
//	// seconds to start time of this stat slot
//		
//	// 
//	// time.Parse("2006-01-02 15:04:05", value)	
//	tm, err := time.Parse(fmt, tstr)
//	if err != nil {
//		return "", err
//	}	
//	
//	ts := tm.Unix()	
//		
//	t := ts%3600
//	
//	t = ts/3600*3600 + t/this.flushInterval*this.flushInterval
//	
//	tf := time.Unix(t, 0)
//	return tf.Format(fmt), nil 
//}
