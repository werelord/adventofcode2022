package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

// https://adventofcode.com/2022/day/2

type Shape int
type Result int

const (
	Rock     Shape = 1
	Paper    Shape = 2
	Scissors Shape = 3

	Lose Result = 0
	Draw Result = 3
	Win  Result = 6
)

func (s Shape) String() string {
	// fucking go generate doesn't work
	var str string
	switch s {
	case Rock:
		str = "rock"
	case Scissors:
		str = "scissors"
	case Paper:
		str = "paper"
	default:
		panic("unrecognized shape")
	}

	return fmt.Sprintf("%v(%v)", str, int(s))

}

func (r Result) String() string {
	var str string
	switch r {
	case Win:
		str = "win"
	case Lose:
		str = "lose"
	case Draw:
		str = "draw"
	default:
		panic("unrecognized result")
	}

	return fmt.Sprintf("%v(%v)", str, int(r))
}

func day2() {
	file := filepath.Join(currentDir(), "input", "day2.txt")

	var part1score = 0
	var part2score = 0

	for line := range readFile(file) {
		s := strings.Fields(line)
		if len(s) != 2 {
			panic(fmt.Errorf("bad parse: %q", line))
		}

		newpt1 := calcPart1(s...)

		part1score += newpt1
		fmt.Printf("pt 1, score = %v, total = %v\n", newpt1, part1score)

		newpt2 := calcPart2(s...)
		part2score += newpt2
		fmt.Printf("pt 2, score = %v, total = %v\n", newpt2, part2score)

	}

	fmt.Printf("pt 1 score: %v\npt 2 score: %v\n", part1score, part2score)

}

func calcPart1(s ...string) int {

	var (
		opp, you = toShape(s[0]), toShape(s[1])
	)

	return  calcScore(you, calcWLD(opp, you))
}
func calcPart2(s ...string) int {
	var (
		opp = toShape(s[0])
		res = toResult(s[1])
		you = calcShape(opp, res)
	)
	return calcScore(you, res)
}

func toShape(str string) Shape {
	if str == "A" || str == "X" {
		return Rock
	} else if str == "B" || str == "Y" {
		return Paper
	} else if str == "C" || str == "Z" {
		return Scissors
	} else {
		panic("unrecognized shape")
	}
}

func toResult(str string) Result {
	switch {
	case str == "X":
		return Lose
	case str == "Y":
		return Draw
	case str == "Z":
		return Win
	default:
		panic("unrecogized result")
	}
}

var (
	/*
		WLD matrix:
			you	R	P	S
			_____________
		opp R |	D	W	L
			P |	L	D	W
			S |	W	L	D


		Result matrix 
			res	L	D	W
			_____________
		opp R |	S	R	P
			P |	R	P	S
			S |	P	S	R
		

	*/
	wldMatrix [][]Result = [][]Result{
		{Draw, Win, Lose},
		{Lose, Draw, Win},
		{Win, Lose, Draw}}

	resMatrix [][]Shape = [][]Shape{
		{Scissors, Rock, Paper},
		{Rock, Paper, Scissors},
		{Paper, Scissors, Rock},
	}
)

func calcShape(opp Shape, resNeeded Result) Shape {
	// opp shape index is 1 index; we can figure out shape index by shape/3
	var (
		oppIndex = int(opp) - 1
		resIndex = int(int(resNeeded)/3)
		shape = resMatrix[oppIndex][resIndex]
	)
	//fmt.Printf("opp:%v, result:%v, shape:%v\n", opp, resNeeded, shape)
	return shape

}

func calcWLD(opp Shape, you Shape) Result {
	var (
		// shapes are 1-indexed, convert to 0 indexed
		oppIndex = int(opp) - 1
		youIndex = int(you) - 1
		res      = wldMatrix[oppIndex][youIndex]
	)
	//fmt.Printf("opp:%v, you,%v, result:%v\n", opp, you, res)
	return res
}
func calcScore(s Shape, wld Result) int {
	//fmt.Printf("shape %v, result %v\n", s, wld)
	return int(s) + int(wld)
}
