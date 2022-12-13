package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

func day13() {

	var (
		file = filepath.Join(currentDir(), "input", "day13.txt")

		pairCount        int
		correctOrderList = make([]int, 0)

		sortedPackets []packetType
	)

	if false {
		file = filepath.Join(currentDir(), "input", "day13.test.txt")
	}

	var lineChan = readLines(file)
	for {
		pairCount++
		var (
			left  = <-lineChan
			right = <-lineChan
			_     = <-lineChan // blank line
		)

		if left == "" {
			// check for channel closed
			if _, isOpen := <-lineChan; isOpen == false {
				break
			} else {
				panic("channel not closed, but left is empty")
			}
		}
		// fmt.Printf("Pair %v\n", pairCount)
		var lpacket, rpacket = decodePackets(left, right)
		if lessThan := comparePackets(lpacket, rpacket); lessThan {
			// fmt.Printf("adding to list pair #%v\n", pairCount)
			correctOrderList = append(correctOrderList, pairCount)

		}
		sortedPackets, _ = insertSorted(sortedPackets, lpacket)
		sortedPackets, _ = insertSorted(sortedPackets, rpacket)

		// fmt.Print("\n")
	}

	var lowDiv, highDiv = decodePackets("[[2]]", "[[6]]")
	var lowIndex, highIndex int
	sortedPackets, lowIndex = insertSorted(sortedPackets, lowDiv)
	_, highIndex = insertSorted(sortedPackets, highDiv)

	// fmt.Println(correctOrderList)
	var tot int
	for _, val := range correctOrderList {
		tot += val
	}

	fmt.Printf("total: %v\n", tot)

	// for i, packet := range sortedPackets {
	// 	fmt.Printf("%v: %v\n", i, packet)
	// }

	fmt.Printf("lowIndex: %v, highIndex: %v total: %v", lowIndex+1, highIndex+1, (lowIndex+1) * (highIndex+1))
}

func insertSorted(sortedList []packetType, p packetType) ([]packetType, int) {

	// fmt.Printf("inserting : %v\n", p)
	i := sort.Search(len(sortedList), func(i int) bool {
		return (p.compareTo(sortedList[i]) == LessThan)
	})

	if i == len(sortedList) {
		sortedList = append(sortedList, p)
	} else {
		// insert into i
		sortedList = append(sortedList[:i+1], sortedList[i:]...)
		sortedList[i] = p
	}
	return sortedList, i
}

func decodePackets(left, right string) (packetType, packetType) {
	var (
		leftdec  = json.NewDecoder(strings.NewReader(left))
		rightdec = json.NewDecoder(strings.NewReader(right))
	)
	leftdec.UseNumber()
	rightdec.UseNumber()

	if leftpacket, err := decode(leftdec); err != nil {
		fmt.Print(err)
		panic(err)
	} else if rightpacket, err := decode(rightdec); err != nil {
		fmt.Print(err)
		panic(err)
	} else {
		return leftpacket, rightpacket
	}
}

func comparePackets(leftpacket, rightpacket packetType) (bool) {
	// fmt.Printf("comparing %v to %v\n", leftpacket, rightpacket)
	switch cmp := leftpacket.compareTo(rightpacket); cmp {
	case LessThan:
		// fmt.Println("Left is smaller than right, returning true")
		return true
	case GreaterThan:
		// fmt.Println("right is smaller than left, returning false")
		return false
	case Equal:
		panic("Equal - this should not happen")
	default:
		panic("default should not happen")
	}
}

type compareType int

const (
	arraytype = iota
	inttype
)
const (
	LessThan compareType = iota
	Equal
	GreaterThan
)

type packetType interface {
	DataType() int
	compareTo(packetType) compareType
	String() string
}

type arrayData struct {
	data []packetType
}
type intData struct {
	data int
}

func (ad arrayData) DataType() int { return arraytype }
func (id intData) DataType() int   { return inttype }

func (ad arrayData) String() string {
	ret := "["
	for i, d := range ad.data {
		if i > 0 {
			ret += ","
		}
		ret += d.String()
	}
	ret += "]"
	return ret
}
func (id intData) String() string {
	return fmt.Sprintf("%v", id.data)
}

func (ad arrayData) compareTo(other packetType) compareType {
	// fmt.Printf("comparing %v to %v\n", ad, other)
	switch otherVal := other.(type) {
	case arrayData:

		leftSlice := ad.data
		rightSlice := otherVal.data

		for {
			if len(leftSlice) == 0 && len(rightSlice) == 0 {
				return Equal
			} else if len(leftSlice) == 0 {
				// fmt.Print("left ran out of items first, returning lessthan\n")
				return LessThan
			} else if len(rightSlice) == 0 {
				// fmt.Print("right ran out of items first, returning greaterthan\n")
				return GreaterThan
			} else {
				switch cmp := leftSlice[0].compareTo(rightSlice[0]); cmp {
				case LessThan:
					fallthrough
				case GreaterThan:
					return cmp
				case Equal:
					leftSlice = leftSlice[1:]
					rightSlice = rightSlice[1:]
				default:
					panic("unrecognized compare")

				}
			}
		}

	case intData:
		// convert other into array, compare with this
		var otherArray = arrayData{}
		otherArray.data = append(otherArray.data, other)
		return ad.compareTo(otherArray)

	default:
		panic(fmt.Sprintf("unrecognized type: %#v", otherVal))
	}
}

func (id intData) compareTo(other packetType) compareType {
	// fmt.Printf("comparing %v to %v\n", id, other)
	switch otherVal := other.(type) {
	case arrayData:
		// convert this into a list, and compare between the two
		var thisData = arrayData{}
		thisData.data = append(thisData.data, id)
		return thisData.compareTo(other)
	case intData:
		if id.data < otherVal.data {
			return LessThan
		} else if id.data == otherVal.data {
			return Equal
		} else {
			return GreaterThan
		}
	default:
		panic(fmt.Sprintf("unrecognized type: %#v", otherVal))
	}
}

func decode(dec *json.Decoder) (packetType, error) {

	if token, err := dec.Token(); err != nil {
		return nil, err
	} else {
		switch tokenVal := token.(type) {
		case json.Delim:
			// case close, open
			if tokenVal.String() == "[" {
				var arrdata = arrayData{}
				for {
					if dec.More() == false {
						// pull off closing token
						if close, err := dec.Token(); err != nil {
							return nil, err
						} else if cval, ok := close.(json.Delim); !ok || cval.String() != "]" {
							return nil, errors.New("wrong close found")
						} else {
							break
						}
					} else if newData, err := decode(dec); err != nil {
						return nil, err
					} else {
						arrdata.data = append(arrdata.data, newData)
					}
				}
				return arrdata, nil
			} else {
				return nil, fmt.Errorf("unrecognized delim: '%v'", tokenVal.String())
			}

		case json.Number:
			if val, err := tokenVal.Int64(); err != nil {
				return nil, err
			} else {
				var data = intData{data: int(val)}
				return data, nil
			}

		default:
			return nil, errors.New("unrecognized type")
		}
	}
}
