package main

import (
	"errors"
	"fmt"
)

// スタック構造体
type Stack struct {
	elements []int
}

func main() {
	stack := &Stack{}

	//Push
	stack.Push(10)
	stack.Push(20)
	stack.Push(30)

	fmt.Println("Peek")
	fmt.Println(stack.Peek())

	fmt.Println("Pop")
	value, _ := stack.Pop()
	fmt.Println(value)

	fmt.Println("Pop")
	value, _ = stack.Pop()
	fmt.Println(value)

	fmt.Println("Pop")
	value, _ = stack.Pop()
	fmt.Println(value)
}

// Push(スタックに要素追加)
func (s *Stack) Push(value int) {
	s.elements = append(s.elements, value)
}

// Pop(スタックから要素を取り出す)
func (s *Stack) Pop() (int, error) {
	if len(s.elements) == 0 {
		return 0, errors.New("stack is empty")
	}
	value := s.elements[len(s.elements)-1]
	s.elements = s.elements[:len(s.elements)-1]
	return value, nil
}

// Peek(スタックのトップ確認)
func (s *Stack) Peek() (int, error) {
	if len(s.elements) == 0 {
		return 0, errors.New("stack is empty")
	}
	return s.elements[len(s.elements)-1], nil
}
