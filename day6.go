package main

import (
	"fmt"
	"path/filepath"
)

// https://adventofcode.com/2022/day/6

func day6(inp string) {
	file := filepath.Join(currentDir(), "input", inp)

	var (
		runeCount   int
		packetQueue = NewQueue[rune]()
		messQueue   = NewQueue[rune]()

		packetFound, messageFound int
	)

	const (
		packetLimit  = 4
		messageLimit = 14
	)

	for runeStr := range readRunes(file) {
		if len(runeStr) != 1 {
			panic("rune string isn't one character")
		}
		runeCount++
		var r = ([]rune(runeStr))[0]
		// fmt.Println(string(r))

		// this one is a simple map/reduce function.. simple utf8.runeCount() should do it
		if packetFound == 0 {
			// always append
			packetQueue.Enqueue(r)

			if packetQueue.Len() > packetLimit {
				// dequeue oldest and check
				packetQueue.Dequeue()

				var qrc = stupidCountUnique(packetQueue.data)

				if qrc == packetLimit { // four unique runes
					fmt.Printf("found on rune %c, runecount: %v, queue: %v\n", r, runeCount, string(packetQueue.data))
					packetFound = runeCount
				}
			}
		}

		if messageFound == 0 {
			// always append
			messQueue.Enqueue(r)
			if messQueue.Len() > messageLimit {
				// dequeue oldest and check
				messQueue.Dequeue()

				var qrc = stupidCountUnique(messQueue.data)

				if qrc == messageLimit { // fourteen unique runes
					fmt.Printf("found on rune %c, runecount: %v, queue: %v\n", r, runeCount, string(messQueue.data))
					messageFound = runeCount
				}
			}
		}

		if packetFound > 0 && messageFound > 0 {
			// found everything
			break
		}

	}

	fmt.Printf("packet found, runecount: %v, queue:%v\n", packetFound, string(packetQueue.data))
	fmt.Printf("message found, runecount: %v, queue:%v\n", messageFound, string(messQueue.data))

}

func stupidCountUnique(inp []rune) int {
	var runeCount = make(map[rune]int, 0)

	for _, r := range inp {
		if count, exists := runeCount[r]; exists {
			runeCount[r] = count + 1
		} else {
			runeCount[r] = 1
		}
	}

	return len(runeCount)
}
