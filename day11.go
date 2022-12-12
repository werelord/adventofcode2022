package main

import (
	"fmt"
	"math"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
)

//  https://adventofcode.com/2022/day/11

// my god, overflow fun.. Recognized that right away, and noticed immediately that all the divisors were prime..
// but since its been 20+ years since college math, I didn't recognize the least common multiple hack until
// after hints on reddit..

// tricky tricky..

func day11() {
	var (
		monkeylist    []*monkeyType
		LCM           uint
		maxRounds     = 20
		decreaseWorry = true

		roundCheck = []int{1, 20, 1000, 2000, 3000, 4000, 5000, 6000, 7000, 8000, 9000, 10000}
	)

	monkeylist, LCM = generateMonkeys(false)

	if true { // part2
		maxRounds = 10000
		decreaseWorry = false
	}

	for round := 1; round <= maxRounds; round++ {

		for _, monkey := range monkeylist {
			// get the first item
			if monkey.itemQueue.Len() == 0 {
				// no items, skip
				continue
			}

			for item := range monkey.itemQueue.DequeueIter() {
				monkey.inspectCount++

				// keep the item within bounds
				newWorry := monkey.operation(item) % LCM

				if decreaseWorry {
					newWorry = uint(math.Floor(float64(newWorry) / 3.0))
				}

				// if decreaseWorry == false && newWorry < item {
				// 	fmt.Printf("op(1): %v\nop(2): %v\n", monkey.operation(1), monkey.operation(2))
				// 	fmt.Println("Overflow detected!!!")
				// }

				if monkey.test(newWorry) {
					monkeylist[monkey.testTrue].catch(newWorry)
				} else {
					monkeylist[monkey.testFalse].catch(newWorry)
				}
			}
		}

		// end of round
		if slices.Contains(roundCheck, round) {
			fmt.Printf("round %v:\n", round)
			for _, m := range monkeylist {
				fmt.Printf("  %v\n", m)
			}
		}
	}

	fmt.Println("finished:")
	var max1, max2 uint64
	for _, m := range monkeylist {
		fmt.Printf("  %v\n", m)

		if m.inspectCount > max1 {
			max2 = max1
			max1 = m.inspectCount
		} else if m.inspectCount > max2 {
			max2 = m.inspectCount
		}
	}
	fmt.Printf("%v * %v = %v\n", max1, max2, (max1 * max2))
}

type monkeyType struct {
	id           string
	inspectCount uint64
	itemQueue    Queue[uint]
	operation    func(uint) uint
	test         func(uint) bool

	testTrue, testFalse int
}

func (m *monkeyType) catch(item uint) {
	m.itemQueue.Enqueue(item)
}

func (m monkeyType) String() string {
	return fmt.Sprintf("%v:(inspect:%v) %v", m.id, m.inspectCount, m.itemQueue.data)
}

func generateMonkeys(test bool) ([]*monkeyType, uint) {
	var filename string
	if test {
		// return test11()
		filename = filepath.Join(currentDir(), "input", "day11.test.txt")
	} else {
		filename = filepath.Join(currentDir(), "input", "day11.txt")
	}

	// because I don't want to make mistakes hardcoding.. and where's the fun in that??
	var (
		ret      = make([]*monkeyType, 0)
		LCM uint = 1
	)

	fileChan := readLines(filename)
	for {

		var monkeyStr = []string{<-fileChan, <-fileChan, <-fileChan, <-fileChan,
			<-fileChan, <-fileChan, <-fileChan}

		if monkeyStr[0] == "" {
			// doublecheck that the channel is closed
			if _, isOpen := <-fileChan; isOpen == false {
				break
			} else {
				panic("channel not closed, but monkey name is empty")
			}
		}
		// parse monkey
		monkey, prime := parseMonkey(monkeyStr)
		LCM *= prime
		ret = append(ret, monkey)
	}

	return ret, LCM
}

func parseMonkey(strList []string) (*monkeyType, uint) {

	// name
	var monkey monkeyType
	var prime uint

	matchName := regexp.MustCompile("(Monkey [0-9]+):").FindStringSubmatch(strList[0])
	if len(matchName) != 2 {
		panic("no match" + strList[0])
	}
	monkey.id = matchName[1]

	// itemlist
	monkey.itemQueue = NewQueue[uint]()
	itemMatch := regexp.MustCompile("  Starting items: ([0-9 ,]*)").FindStringSubmatch(strList[1])
	if len(itemMatch) != 2 {
		panic("no match" + strList[1])
	}
	itemListStr := strings.FieldsFunc(itemMatch[1], func(r rune) bool { return r == ',' || r == ' ' })
	for _, item := range itemListStr {
		if item, err := strconv.ParseUint(item, 10, 32); err != nil {
			panic(err)
		} else {
			monkey.itemQueue.Enqueue(uint(item))
		}
	}

	// operation
	opMatch := regexp.MustCompile(`  Operation: new = old ([\+\*]) ([0-9]+|old)`).FindStringSubmatch(strList[2])
	if len(opMatch) != 3 {
		panic("no match" + strList[2])
	}
	// operation is only times or plus
	oper := opMatch[1]
	targetStr := opMatch[2]
	if targetStr == "old" {
		if oper == "+" {
			monkey.operation = func(i uint) uint { return i + i }
		} else if oper == "*" {
			monkey.operation = func(i uint) uint { return i * i }
		} else {
			panic("unrecognized operation" + strList[2])
		}
	} else {
		if target, err := strconv.ParseUint(targetStr, 10, 32); err != nil {
			panic(err)
		} else if oper == "+" {
			monkey.operation = func(i uint) uint { return i + uint(target) }
		} else if oper == "*" {
			monkey.operation = func(i uint) uint { return i * uint(target) }
		} else {
			panic("unrecognized operation" + strList[2])
		}
	}

	// test match
	testMatch := regexp.MustCompile("  Test: divisible by ([0-9]+)").FindStringSubmatch(strList[3])
	if len(testMatch) != 2 {
		panic("no match" + strList[3])
	}
	if testVal, err := strconv.ParseUint(testMatch[1], 10, 32); err != nil {
		panic(err)
	} else {
		prime = uint(testVal)
		monkey.test = func(i uint) bool { return (i % prime) == 0 }
	}

	// assign targets
	trueMatch := regexp.MustCompile("    If true: throw to monkey ([0-9]+)").FindStringSubmatch(strList[4])
	falseMatch := regexp.MustCompile("    If false: throw to monkey ([0-9]+)").FindStringSubmatch(strList[5])
	if len(trueMatch) != 2 {
		panic("no match" + strList[4])
	} else if len(falseMatch) != 2 {
		panic("no match" + strList[5])
	}
	if trueVal, err := strconv.Atoi(trueMatch[1]); err != nil {
		panic(err)
	} else if falseVal, err := strconv.Atoi(falseMatch[1]); err != nil {
		panic(err)
	} else {
		monkey.testTrue = trueVal
		monkey.testFalse = falseVal
	}

	return &monkey, prime
}
