package stack

type Stack[T any] struct {
	data []T
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
		data: make([]T, 0),
	}
}

func (s *Stack[T]) Push(item T) {
	s.data = append(s.data, item)
}

func (s *Stack[T]) Pop() {
	if s.Empty() {
		return
	}

	s.data = s.data[:len(s.data)-1]
}

func (s *Stack[T]) Top() T {
	return s.data[len(s.data)-1]
}

func (s *Stack[T]) Empty() bool {
	return len(s.data) == 0
}
