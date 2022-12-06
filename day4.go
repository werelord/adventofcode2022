package main

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
)

// https://adventofcode.com/2022/day/4

type Pair struct {
	lower int
	upper int
}

func day4(inp string) {
	file := filepath.Join(currentDir(), "input", inp)

	var (
		fullyContainCount, overlapCount int
	)

	for line := range readFile(file) {

		var first, second = toPair(line)

		if first.fullyContains(second) || second.fullyContains(first) {
			fullyContainCount++
			//fmt.Printf("new count: %v\n", fullyContainCount)
		}
		if first.overlaps(second) {
			overlapCount++
			//fmt.Printf("new overlap count: %v\n", overlapCount)
		}
	}
	fmt.Printf("\ncontains count: %v\noverlap count: %v\n", fullyContainCount, overlapCount)

}

func toPair(s string) (Pair, Pair) {

	var (
		rx      = regexp.MustCompile("([0-9]*)-([0-9]*),([0-9]*)-([0-9]*)")
		matches = rx.FindStringSubmatch(s)

		pair1, pair2 Pair
		err          error
	)

	if len(matches) != 5 { // first contains the entire string
		panic(fmt.Sprintf("regex shit: %v", s))
	}

	if pair1.lower, err = strconv.Atoi(matches[1]); err != nil {
		panic(fmt.Sprintf("bad conversion: %v, err: %v", matches[1], err))
	} else if pair1.upper, err = strconv.Atoi(matches[2]); err != nil {
		panic(fmt.Sprintf("bad conversion: %v, err: %v", matches[2], err))
	} else if pair2.lower, err = strconv.Atoi(matches[3]); err != nil {
		panic(fmt.Sprintf("bad conversion: %v, err: %v", matches[3], err))
	} else if pair2.upper, err = strconv.Atoi(matches[4]); err != nil {
		panic(fmt.Sprintf("bad conversion: %v, err: %v", matches[4], err))
	}
	// fmt.Printf("parsed: %#v : %#v\n", pair1, pair2)

	return pair1, pair2
}

func (p Pair) fullyContains(other Pair) bool {

	if (other.lower >= p.lower) && (p.upper >= other.upper) {
		// fmt.Printf("%#v contains %#v\n", p, other)
		return true
	}
	return false
}

func (p Pair) overlaps(other Pair) (ol bool) {
	if p.lower <= other.lower && other.lower <= p.upper {
		ol = true
	} else if other.lower <= p.lower && p.lower <= other.upper {
		ol = true
	}

	//fmt.Printf("%v (%v):(%v)\n", ol, p, other)
	return ol
}
