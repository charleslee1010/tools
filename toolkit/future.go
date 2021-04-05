package toolkit

import (
    "sync"
    "time"
)


type Future struct {
    isfinished bool
    result     interface{}
    resultchan chan interface{}
    l          sync.Mutex
}

func (f *Future) GetResult() interface{} {
    f.l.Lock()
    defer f.l.Unlock()
    if f.isfinished {
        return f.result
    }

    select {
    // timeout
    case <-time.Tick(time.Second * 60):
        f.isfinished = true
        f.result = nil
        return nil
    case f.result = <-f.resultchan:
        f.isfinished = true
        return f.result
    }
}

func (f *Future) SetResult(result interface{}) {
    if f.isfinished == true {
        return
    }
    f.resultchan <- result
    close(f.resultchan)
}

func NewFuture() *Future {
    return &Future{
        isfinished: false,
        result:     nil,
        resultchan: make(chan interface{}, 1),
    }
}
