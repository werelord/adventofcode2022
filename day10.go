package main

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"

	"golang.org/x/exp/slices"
)

// https://adventofcode.com/2022/day/10

type operation struct {
	val   int
	cycle int
}

func day10() {

	// file := filepath.Join(currentDir(), "input", "day10.test.txt")
	file := filepath.Join(currentDir(), "input", "day10.txt")

	var (
		curCycle       int
		regX           = 1 // default starting point for register
		cycleInterrupt = []int{20, 60, 100, 140, 180, 220}

		// screen = []string{"", "", "", "", "", ""}
		screen string

		signalSum int
	)

	for line := range readLines(file) {

		var curOp = getOperation(line, curCycle)

		for curCycle < curOp.cycle {
			curCycle++
			if slices.Contains(cycleInterrupt, curCycle) {

				var strength = curCycle * regX
				signalSum += strength

				fmt.Printf("cycle %v, regX: %v, strength: %v, totSum: %v\n", curCycle, regX, strength, signalSum)
			}

			// var bufIndex = (curCycle - 1) / 40
			// drawScreen(curCycle, regX, &(screen[bufIndex]))
			drawScreen(curCycle, regX, &screen)

		}

		// after cycle, do operation
		// fmt.Printf("cycle %v: op:%v+x:%v = x%v\n", curCycle, curOp.val, regX, regX+curOp.val)
		regX += curOp.val

	}

	fmt.Printf("total cycles: %v, totSigStrength:%v\nScreen:\n", curCycle, signalSum)

	// for _, buf := range screen {
	// 	fmt.Printf("%v\n", buf)
	// }
	for x := 0; x < 6; x++ {
		fmt.Printf("%v\n", screen[x*40:(x+1)*40])
	}
}

func getOperation(line string, curCycle int) operation {
	var rx = regexp.MustCompile("(addx) ([0-9-]*)")
	var curOp operation
	if line == "noop" {
		curOp = operation{
			cycle: curCycle + 1,
		}
	} else {
		match := rx.FindStringSubmatch(line)
		if len(match) != 3 {
			panic("unable to find operation")
		}
		if match[1] == "addx" {
			if amt, err := strconv.Atoi(match[2]); err != nil {
				panic(err)
			} else {
				curOp = operation{
					cycle: curCycle + 2,
					val:   amt,
				}
			}
		}
	}
	return curOp
}

func drawScreen(cycle, regX int, screenBuf *string) {
	// if (regX <= 0) || (regX >= 40) {
	// fmt.Print("regx beyond screen!\n")
	// }

	var screenPos = (cycle - 1) % 40
	if (regX-1 == screenPos) || (regX == screenPos) || (regX+1 == screenPos) {
		*screenBuf += "#"
	} else {
		*screenBuf += "."
	}
}
