package main

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
)

// https://adventofcode.com/2022/day/5

// because rune.String() outputs code point, and can't use fmt directive to use string
func runeToString(r rune) string {
	return string(r)
}
func runeSliceToString(rs []rune) string {
	var ret string
	for _, r := range rs {
		var delim string
		if len(ret) > 0 {
			delim = " "
		}
		ret = fmt.Sprintf("%v%v'%c'", ret, delim, r)
	}
	return ret
}

func day5(inp string) {

	file := filepath.Join(currentDir(), "input", inp)

	/*	hard coding initial setup
		[T]     [Q]             [S]
		[R]     [M]             [L] [V] [G]
		[D] [V] [V]             [Q] [N] [C]
		[H] [T] [S] [C]         [V] [D] [Z]
		[Q] [J] [D] [M]     [Z] [C] [M] [F]
		[N] [B] [H] [N] [B] [W] [N] [J] [M]
		[P] [G] [R] [Z] [Z] [C] [Z] [G] [P]
		[B] [W] [N] [P] [D] [V] [G] [L] [T]
		1   2   3   4   5   6   7   8   9
	*/

	var (
		moveCommands int
		crates       = []Stack[rune]{
			NewStack([]rune("BPNQHDRT")...),
			NewStack([]rune("WGBJTV")...),
			NewStack([]rune("NRHDSVMQ")...),
			NewStack([]rune("PZNMC")...),
			NewStack([]rune("DZB")...),
			NewStack([]rune("VCWZ")...),
			NewStack([]rune("GZNCVQLS")...),
			NewStack([]rune("LGJMDNV")...),
			NewStack([]rune("TPMFZCG")...),
		}

		multicrates = []Stack[rune]{
			NewStack([]rune("BPNQHDRT")...),
			NewStack([]rune("WGBJTV")...),
			NewStack([]rune("NRHDSVMQ")...),
			NewStack([]rune("PZNMC")...),
			NewStack([]rune("DZB")...),
			NewStack([]rune("VCWZ")...),
			NewStack([]rune("GZNCVQLS")...),
			NewStack([]rune("LGJMDNV")...),
			NewStack([]rune("TPMFZCG")...),
		}

		// move 13 from 1 to 3
		r = regexp.MustCompile("move ([0-9]+) from ([0-9]+) to ([0-9]+)")
	)

	for line := range readFile(file) {

		matches := r.FindStringSubmatch(line)

		if len(matches) != 4 { // first contains entire string
			if moveCommands > 0 {
				panic(fmt.Sprintf("no match after move commands started: %v\n", line))
			}
		} else {
			moveCommands++
			var (
				moveCount int
				from, to  int
				err       error
			)

			if moveCount, err = strconv.Atoi(matches[1]); err != nil {
				panic(err)
			} else if from, err = strconv.Atoi(matches[2]); err != nil {
				panic(err)
			} else if to, err = strconv.Atoi(matches[3]); err != nil {
				panic(err)
			} else {
				// fmt.Printf("%v: count:%v from:%v to:%v\n", matches[1:], moveCount, from, to)
				// correct stack index
				from--
				to--
			}

			movePt1(&crates[from], &crates[to], moveCount)
			movePt2(&multicrates[from], &multicrates[to], moveCount)
		}
	}

	fmt.Printf("moves: %v\n", moveCommands)
	var ans1 string
	for i, stack := range crates {
		fmt.Printf("stack %v: %v\n", i+1, stack.String(runeToString))
		if r, err := stack.Peek(); err != nil {
			panic(err)
		} else {
			ans1 = fmt.Sprintf("%v%c", ans1, r)
		}
	}
	fmt.Printf("part 1 answer: %v\n", ans1)

	var ans2 string
	for i, stack := range multicrates {
		fmt.Printf("stack %v: %v\n", i+1, stack.String(runeToString))
		if r, err := stack.Peek(); err != nil {
			panic(err)
		} else {
			ans2 = fmt.Sprintf("%v%c", ans2, r)
		}
	}
	fmt.Printf("part 2 answer: %v\n", ans2)
}

func movePt1(from, to *Stack[rune], count int) {
	// part 1

	fmt.Printf("p1 before:\n\tfrom:%v\n\t  to:%v\n", from.String(runeToString), to.String(runeToString))
	if items, err := from.Pop(count); err != nil {
		panic(err)
	} else {
		fmt.Printf("p1 moved %v: %s\n", count, runeSliceToString(items))
		to.Push(items...)
	}
	fmt.Printf("p1 after:\n\tfrom:%v\n\t  to:%v\n", from.String(runeToString), to.String(runeToString))

}

func movePt2(from, to *Stack[rune], count int) {
	// part 2
	fmt.Printf("p2 before:\n\tfrom:%v\n\t  to:%v\n", from.String(runeToString), to.String(runeToString))
	if items, err := from.Pop(count); err != nil {
		panic(err)
	} else {
		// even though stack reverses it, don't want to add control into that function.. doesn't make sense..
		// so just do another reverse here
		var revItems = Reverse(items)
		fmt.Printf("p2 moved %v: %s\n", count, runeSliceToString(revItems))
		to.Push(revItems...)
	}
	fmt.Printf("p2 after:\n\tfrom:%v\n\t  to:%v\n", from.String(runeToString), to.String(runeToString))

}
