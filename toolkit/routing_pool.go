package toolkit

import (
	//	"fmt"
	"sync/atomic"
	"time"
)

type RoutinePool struct {
	tasks     chan interface{}
	bufSize   int   // channel buffer size
	workerNum int   // configured worker number
	running   int32 // current worker number
	//atomic.LoadInt32(&p.release) atomic.StoreInt32(&p.capacity, int32(size)),atomic.AddUint32(&ui32
	f func(int, interface{})
}

func NewRoutinePool() *RoutinePool {
	return &RoutinePool{}
}

func (this *RoutinePool) RoutineWorker(id int) {

	defer func() {
		if p := recover(); p != nil {
			atomic.AddInt32(&this.running, -1)
			//			fmt.Println("worker panic")
		}
	}()

	//	fmt.Println("worker started, id=%d", id)
	atomic.AddInt32(&this.running, 1)

	for true {
		task := <-this.tasks
		this.f(id, task)
	}
}

const (
	MAX_WORKER_NUM = 1024
	MAX_BUF_SIZE   = 10000
)

func (this *RoutinePool) Init(num int, bufSize int, f func(int, interface{})) bool {

	if num > 0 && num <= MAX_WORKER_NUM && bufSize > 0 && bufSize <= MAX_BUF_SIZE && f != nil {
		this.workerNum = num
		this.bufSize = bufSize
		this.f = f
		this.tasks = make(chan interface{}, bufSize)

		// start go routine
		for i := 0; i < this.workerNum; i++ {
			go this.RoutineWorker(i)
		}

		go this.monitoring()

		return true
	}

	return false
}

func (this *RoutinePool) monitoring() {

	ticker := time.NewTicker(1000 * time.Millisecond)
	for true {
		<-ticker.C

		if c := atomic.LoadInt32(&this.running); c < int32(this.workerNum) {
			//			fmt.Println("current nm:", c)
			go this.RoutineWorker(-1)
		}

	}
}

func (this *RoutinePool) AddTask(v interface{}) bool {

	if len(this.tasks) >= this.bufSize {
		// it is going to be blocked
		return false
	} else {
		// put the task into channel
		this.tasks <- v
		return true
	}
}
