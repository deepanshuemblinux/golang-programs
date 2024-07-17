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
func IsValidParenthesis(str string) bool {
        runes := []rune(str)
        stack := NewStack()
        for _, val := range runes  {
                switch val {
                        case '{', '[', '(' :
                                 stack.Push(val)
                        case '}':
                                e,_ := stack.Top()
                                if e != '{' {
                                        return false
                                }
                                stack.Pop()
                        case ')':
                                e,_ := stack.Top()
                                if e != '(' {
                                        return false
                                }
                                stack.Pop()
                        case ']':
                                e,_ := stack.Top()
                                if e != '[' {
                                        return false
                                }
                                stack.Pop()

                }
        }
        return stack.IsEmpty()
}
func main() {
        str := "()[{}]"
	fmt.Println(IsValidParenthesis(str))
        str = "([{}])"
	fmt.Println(IsValidParenthesis(str))
        str = "()[]"
	fmt.Println(IsValidParenthesis(str))
        str = "([)]"
	fmt.Println(IsValidParenthesis(str))
        str = "(()"
	fmt.Println(IsValidParenthesis(str))
        str = "["
	fmt.Println(IsValidParenthesis(str))

}

