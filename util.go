package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func readLines(filename string) chan string {
	return readFile(filename, bufio.ScanLines)
}

func readRunes(filename string) chan string {
	return readFile(filename, bufio.ScanRunes)
}

func readFile(filename string, splitFunc bufio.SplitFunc) chan string {

	var out = make(chan string, 1)
	go func() {
		file, err := os.Open(filename)
		if err != nil {
			fmt.Print("error opening file: ", err)
			close(out)
			return
		}
		defer file.Close()
		defer close(out)

		var fileScanner = bufio.NewScanner(file)
		fileScanner.Split(splitFunc)

		for fileScanner.Scan() {
			out <- fileScanner.Text()
		}
		//fmt.Println("file done")
	}()
	return out
}

// todo: make generic
func abs(i int) int {
	if i < 0  { return -i } else { return i }
}

func currentDir() string {
	if ex, err := os.Executable(); err != nil {
		panic(err)
	} else {
		return filepath.Dir(ex)
	}

}

type Stack[T any] struct {
	data []T
}

func NewStack[T any](vals ...T) Stack[T] {
	var ret Stack[T]
	ret.data = make([]T, 0, len(vals))

	if len(vals) > 0 {
		ret.data = append(ret.data, vals...)
	}
	return ret
}
func (s Stack[T]) Len() int {
	return len(s.data)
}
func (s Stack[T]) Peek() (v T, err error) {
	if len(s.data) == 0 {
		err = errors.New("stack is empty")
		return
	} else {
		return s.data[len(s.data)-1], nil
	}
}
func (s *Stack[T]) Push(v ...T) {
	s.data = append(s.data, v...)
}
func (s *Stack[T]) PopOne() (T, error) {
	v, err := s.Pop(1)
	return v[0], err
}

func (s *Stack[T]) Pop(count int) (v []T, err error) {
	v = make([]T, 0, count)

	if count < 0 {
		err = errors.New("count cannot be < 0")
		return
	} else if count == 0 {
		return
	} else if len(s.data)-count < 0 {
		err = fmt.Errorf("unable to pop %v items", count)
		return
	} else {
		v = Reverse(s.data[len(s.data)-count:])
		s.data = s.data[:len(s.data)-count]
		return v, nil
	}
}
func (s Stack[T]) String(fn func(T) string) string {
	// top to bottom
	var str string
	for i := len(s.data) - 1; i >= 0; i-- {
		var delim = ""
		if len(str) > 0 {
			delim = ":"
		} else {
			delim = "(top) "
		}
		str = fmt.Sprintf("%v%v%v", str, delim, fn(s.data[i]))

	}
	return str
}

func Reverse[T any](inp []T) []T {
	var ret = make([]T, 0, len(inp))
	for i := len(inp) - 1; i >= 0; i-- {
		ret = append(ret, inp[i])
	}
	return ret
}

type Queue[T any] struct {
	data []T
}

func NewQueue[T any](vals ...T) Queue[T] {
	var ret Queue[T]
	ret.data = make([]T, 0, len(vals))

	if len(vals) > 0 {
		ret.data = append(ret.data, vals...)
	}
	return ret
}
func (q Queue[T]) Len() int {
	return len(q.data)
}
func (q Queue[T]) Peek() (v T, err error) {
	if len(q.data) == 0 {
		err = errors.New("stack is empty")
		return
	} else {
		return q.data[0], nil
	}
}
func (q *Queue[T]) Enqueue(v ...T) {
	q.data = append(q.data, v...)
}
func (q *Queue[T]) Dequeue() (v T, err error) {
	if len(q.data) == 0 {
		err = fmt.Errorf("unable to dequeue; queue is empty")
		return
	} else {
		v = q.data[0]
		q.data = q.data[1:]
		return v, nil
	}
}
