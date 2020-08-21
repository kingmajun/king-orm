package king_orm

import (
	"container/list"
	"sync"
)

type stack struct {
	list *list.List
	mu sync.Mutex
}

func NewStack() *stack {
	list := list.New()
	return &stack{list: list,}
}

func (s *stack) Push(t interface{}){
	s.mu.Lock()
	defer s.mu.Unlock()
	s.list.PushFront(t)
}


func  (s *stack) Pop() interface{} {
	s.mu.Lock()
	defer s.mu.Unlock()
	ele := s.list.Front()
	if nil != ele {
		s.list.Remove(ele)
		return ele.Value
	}

	return nil
}


func (s *stack) Len() int {
	return s.list.Len()
}


