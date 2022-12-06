package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"unicode"
)

// https://adventofcode.com/2022/day/3

func day3() {
	file := filepath.Join(currentDir(), "input", "day3.txt")

	var (
		prioritySum = 0 // 7850

		badgeSum = 0

		badgelist = make([]string, 0, 3)
	)

	for line := range readLines(file) {

		var (
			packSize = len(line) / 2
			pack1    = line[:packSize]
			pack2    = line[packSize:]

			foundRune rune
		)

		// badge calculation
		badgelist = append(badgelist, line)
		if len(badgelist) == 3 {
			// do calculation

			var badgeRune rune

			for _, r := range badgelist[0] {
				if strings.ContainsRune(badgelist[1], r) && strings.ContainsRune(badgelist[2], r) {
					badgeRune = r
					break
				}
			}

			if badgeRune == 0 {
				fmt.Printf("line: '%v'\npack1: '%v'\npack2: '%v'\n", line, pack1, pack2)
				panic("common rune not found")
			}

			var bpri = calcRunePriority(badgeRune)
			badgeSum += bpri
			// fmt.Printf("'%v':'%v':'%v'\nfound: '%c'(%v), new total: %v\n",
			// 	badgelist[0], badgelist[1], badgelist[2], badgeRune, bpri, badgeSum)

			// reset for next three
			badgelist = make([]string, 0, 3)
		}

		//fmt.Printf("line: '%v'\npack1: '%v'\npack2: '%v'\n", line, pack1, pack2)

		for _, r := range pack1 {
			if strings.ContainsRune(pack2, r) {
				foundRune = r
				break
			}
		}

		if foundRune == 0 {
			fmt.Printf("line: '%v'\npack1: '%v'\npack2: '%v'\n", line, pack1, pack2)
			panic("common rune not found")
		}

		var pri = calcRunePriority(foundRune)
		prioritySum += pri
		//fmt.Printf("rune: '%c'(%v), sum: %v\n", foundRune, pri, prioritySum)

	}

	fmt.Printf("total priority sum: %v\ntotal badge sum: %v\n ", prioritySum, badgeSum)

}

func calcRunePriority(r rune) int {
	// abcdefghi	jklmnopqrs	tuvwxyz
	// 123456789	0123456789	0123456
	// ABC	DEFGHIJKLM	NOPQRSTUVW	XYZ
	// 789	0123456789	0123456789	012
	var pri int
	if unicode.IsUpper(r) {
		pri = int(r) - int('A') + 27
	} else { // isLower
		pri = int(r) - int('a') + 1
	}
	return pri
}
