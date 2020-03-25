package lib

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewWorkers(t *testing.T) {

	var i int

	workers := NewWorkers(5, func() interface{} {
		i++
		return i
	})

	assert.Len(t, workers.count, 5)

	worker := workers.Get()
	assert.NotNil(t, worker)

	assert.True(t, worker.(int) < 5)

	assert.Len(t, workers.count, 4)

	workers.Get()
	workers.Get()
	workers.Get()
	workers.Get()
	assert.Len(t, workers.count, 0)

	workers.Put(0)
	assert.Len(t, workers.count, 1)
	workers.Put(0)
	assert.Len(t, workers.count, 2)
	workers.Put(0)
	assert.Len(t, workers.count, 3)
}

func TestWorkers_Wait(t *testing.T) {

	var i int

	workers := NewWorkers(5, func() interface{} {
		i++
		return i
	})

	workers.Wait()

	workers.Get()

	var c chan int
	c = make(chan int, 0)
	f := func() {
		workers.Wait()
		fmt.Println("must be not reach")
		c <- 0
	}

	go f()
	select {
	case <-c:
		panic("Must be not reach")
	case <-time.After(time.Millisecond * 500):

	}
}
