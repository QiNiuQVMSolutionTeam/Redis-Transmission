package lib

import (
	"sync"
)

type Workers struct {
	count chan interface{}
	wg    sync.WaitGroup
}

func NewWorkers(n int, initFunc func() interface{}) *Workers {
	ws := Workers{}
	ws.count = make(chan interface{}, n)
	for i := 0; i < n; i++ {
		ws.count <- initFunc()
	}

	return &ws
}

func (ws *Workers) Get() interface{} {
	ws.wg.Add(1)
	return <-ws.count
}

func (ws *Workers) Put(obj interface{}) {
	ws.count <- obj
	ws.wg.Done()
}

func (ws *Workers) Wait() {

	ws.wg.Wait()
}

func (ws *Workers) IdleCount() int {

	return len(ws.count)
}
