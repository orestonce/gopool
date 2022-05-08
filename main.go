package gopool

import "sync"

type ThreadPool struct {
	ch             chan func()
	wg             sync.WaitGroup
	isClosed       bool
	isClosedLocker sync.Mutex
}

func NewThreadPool(maxSize int) *ThreadPool {
	if maxSize <= 0 {
		panic("gopool.NewThreadPool maxSize <= 0")
	}
	this := &ThreadPool{
		ch:             make(chan func()),
		wg:             sync.WaitGroup{},
		isClosed:       false,
		isClosedLocker: sync.Mutex{},
	}
	this.wg.Add(maxSize)
	for idx := 0; idx < maxSize; idx++ {
		go this.workThreadRun()
	}
	return this
}

func (this *ThreadPool) workThreadRun() {
	defer this.wg.Done()

	for job := range this.ch {
		job()
	}
}

func (this *ThreadPool) AddJob(fn func()) {
	this.isClosedLocker.Lock()
	if this.isClosed {
		this.isClosedLocker.Unlock()
		panic("gopool.ThreadPool: AddJob after close!")
	}
	this.ch <- fn
	this.isClosedLocker.Unlock()
}

func (this *ThreadPool) CloseAndWait() {
	this.isClosedLocker.Lock()
	if this.isClosed {
		this.isClosedLocker.Unlock()
		return
	}
	this.isClosed = true
	close(this.ch)
	this.wg.Wait()
	this.isClosedLocker.Unlock()
}
