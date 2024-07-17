package main

import (
	"fmt"
)

type Stack struct {
	buff []any
}

func NewStack() *Stack {
	return &Stack{
		buff: make([]any, 0),
	}
}

func (s *Stack) IsEmpty() bool {
	if len(s.buff) == 0 {
		return true
	}
	return false
}

func (s *Stack) Push(elem any) {
	s.buff = append(s.buff, elem)
}

func (s *Stack) Top() (any, error) {
	if len(s.buff) == 0 {
		return -1, fmt.Errorf("Stack is empty")
	}
	return s.buff[len(s.buff)-1], nil
}

func (s *Stack) Pop() (any, error) {
	e, err := s.Top()
	if err != nil {
		return -1, nil
	}
	s.buff = s.buff[:len(s.buff)-1]
	return e, nil
}
func main() {
	stack := NewStack()
	fmt.Println(stack.IsEmpty())
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)
	fmt.Printf("%+v\n", stack)
	t, err := stack.Top()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Top is %d\n", t)

	e, err := stack.Pop()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Popped element is %d\n", e)
	t, err = stack.Top()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Top is %d\n", t)
	fmt.Println(stack.IsEmpty())
	stack.Pop()
	stack.Pop()
	fmt.Println(stack.IsEmpty())
}
