package util

import (
	"testing"
	"time"
)

type Temp struct {
	va  int
	va2 int
	va3 int
}

func TestQueueSimple(t *testing.T) {
	//for i := 0; i < MinQueueLen; i++ {
	//queue := NewQueue[*Temp]()
	//queue := NewQueue()
	//queue := Queue{}
	var queue []*Temp

	temp := &Temp{}
	now := time.Now()
	for j := 0; j < 5000; j++ {
		//queue.Push(temp)
		//queue.Push(i)
		//arrays = append(arrays, i)
	}
	t.Log(time.Now().Sub(now).Milliseconds())
	now = time.Now()
	for j := 0; j < 5000; j++ {
		queue.Pop()
		//queue.Push(i)
		//arrays = append(arrays, i)
	}
	t.Log(time.Now().Sub(now).Milliseconds())
}
