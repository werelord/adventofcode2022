package main

import (
	"fmt"
	"path/filepath"
	"strconv"
)

// https://adventofcode.com/2022/day/1
func day1() {

	file := filepath.Join(currentDir(), "input", "day1.txt")

	var (
		// maxElf   int
		// maxTotal int64

		currentTotal int64 = 0

		calList = make([]int64, 0)
	)

	for line := range readFile(file) {

		if line == "" {
			//fmt.Printf("%v has %v\n", len(calList)+1, currentTotal)
			// check max
			// if currentTotal > maxTotal {
			// 	fmt.Printf("setting new max to %v, %v\n", currentElf, currentTotal)
			// 	// new elf, record
			// 	maxTotal = currentTotal
			// 	maxElf = currentElf
			// }

			// reset, current elf and total
			// store the total, via index
			calList = append(calList, currentTotal)
			currentTotal = 0

		} else {
			var cal, err = strconv.ParseInt(line, 10, 64)
			if err != nil {
				panic(err)
			}
			currentTotal += cal
			//fmt.Printf("elf %v: %v new current %v\n", currentElf, cal, currentTotal)
		}
	}

	// there's probably an algorithm for this, but meh
	var max1idx, max2idx, max3idx int

	for idx, curTotal := range calList {
		if (max1idx == 0) || curTotal > calList[max1idx] {
			// new max found, push others down
			max3idx = max2idx
			max2idx = max1idx
			max1idx = idx

		} else if (max2idx == 0) || curTotal > calList[max2idx] {
			max3idx = max2idx
			max2idx = idx

		} else if (max3idx == 0) || curTotal > calList[max3idx] {
			max3idx = idx
		}
	}

	fmt.Printf("max elf: %v, total %v\n", max1idx+1, calList[max1idx])
	fmt.Printf("max elf: %v, total %v\n", max2idx+1, calList[max2idx])
	fmt.Printf("max elf: %v, total %v\n", max3idx+1, calList[max3idx])
	fmt.Printf("total cals: %v", calList[max1idx]+calList[max2idx]+calList[max3idx])
}
