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
	//var stack [5000]*Temp
	//stack[0] = b.root
	//var top int = 1

	temp := &Temp{}
	//queue := NewQueue[*Temp]()

	//var stack [5000]*Temp
	//stack[0] = b.root
	//var top int = 1

	for i := 0; i < 100; i++ {
		now := time.Now()
		var stack []*Temp
		//var cnt int
		//queue := NewQueue[*Temp]()
		for j := 0; j < 10000; j++ {
			stack = append(stack, temp)
			//queue.Push(temp)
			//queue.Pop()
			//queue.Push(i)
			//arrays = append(arrays, i)
		}
		t.Log(time.Now().Sub(now).Milliseconds())
	}

	//now = time.Now()
	//for j := 0; j < 5000; j++ {
	//	queue.Pop()
	//	//queue.Push(i)
	//	//arrays = append(arrays, i)
	//}
	//t.Log(time.Now().Sub(now).Milliseconds())
}
