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
	queue := NewQueue[*Temp]()
	//queue := NewQueue()
	//queue := Queue{}

	temp := &Temp{}
	now := time.Now()
	for j := 0; j < 5000; j++ {
		queue.Add(temp)
		//queue.Push(i)
		//arrays = append(arrays, i)
	}
	t.Log(time.Now().Sub(now).Milliseconds())
	now = time.Now()
	for j := 0; j < 5000; j++ {
		queue.Remove()
		//queue.Push(i)
		//arrays = append(arrays, i)
	}
	t.Log(time.Now().Sub(now).Milliseconds())
}
