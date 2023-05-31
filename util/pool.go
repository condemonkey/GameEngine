package util

type Poolable interface {
	Node() *PoolNode
	SetNode(node *PoolNode)
}

type PoolNode struct {
	next *PoolNode
	item Poolable
}

func newNode(item Poolable) *PoolNode {
	node := &PoolNode{
		item: item,
	}
	item.SetNode(node)
	return node
}

type Pool[T Poolable] struct {
	head     *PoolNode
	factory  func() T
	capacity int
}

func NewPool[T Poolable](capacity int, factory func() T) *Pool[T] {
	pool := &Pool[T]{
		head:     nil,
		capacity: capacity,
		factory:  factory,
	}

	pool.Resize()
	return pool
}

func (s *Pool[T]) Resize() {
	for i := 0; i < s.capacity; i++ {
		head := newNode(s.factory())
		head.next = s.head
		s.head = head
	}
}

func (s *Pool[T]) Push(item T) {
	node := item.Node()
	node.next = s.head
	s.head = node
}

func (s *Pool[T]) Pop() T {
	if s.head == nil {
		s.Resize()
	}
	tmpHead := s.head
	s.head = tmpHead.next
	return tmpHead.item.(T)
}

// IsEmpty returns true if the stack is empty, one the other hand, it returns false if it is not empty
func (s *Pool[T]) IsEmpty() bool {
	return s.head == nil
}

// Cap returns current capacity of stack
func (s *Pool[T]) Cap() (cnt int) {
	tmpHead := s.head
	for tmpHead != nil {
		cnt++
		tmpHead = tmpHead.next
	}
	return cnt
}
