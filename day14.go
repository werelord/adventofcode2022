package main

import (
	"fmt"
	"image"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	// tm "github.com/buger/goterm"
)

var (
	outFile *os.File
)

const test = false

func day14() {

	setupOutFile()
	defer outFile.Close()

	var file = filepath.Join(currentDir(), "input", "day14.txt")
	if test {
		file = filepath.Join(currentDir(), "input", "day14.test.txt")
	}

	var (
		sandmap = loadMap(file)
	)
	sandmap.fillSand()

	fmt.Printf("finished, first bottom: %v, source blocked: %v\n", sandmap.firstAtBottomCount, sandmap.sourceBlockedCount)

}

func setupOutFile() {
	filename := filepath.Join(currentDir(), "out.txt")
	if file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666); err != nil {
		panic(err)
	} else {
		outFile = file
	}
}

func toFile(str string) {
	outFile.Seek(0, io.SeekStart)
	if _, err := outFile.WriteString(str); err != nil {
		fmt.Printf("Error writing file: %v", err)
	}
}

func loadMap(file string) sandmap {

	var sandmap = sandmap{
		sandStart:  image.Pt(500, 0),
		rowMax:     0,
		colMin:     99999,
		colMax:     0,
		blockedMap: make(map[string]string, 0),
	}

	for line := range readLines(file) {

		var wall = parsePoints(line)

		// got a list of points, min 2..
		if len(wall) < 2 {
			panic("wall not long enough")
		}

		// because we're pulling two off, accomodate
		for idx := 0; idx < len(wall)-1; idx++ {
			var start, end = wall[idx], wall[idx+1]

			var blockedPoints = expandPoints(start, end)

			for _, pt := range blockedPoints {

				if pt.Y+1 > sandmap.rowMax {
					sandmap.rowMax = pt.Y + 1
				}
				sandmap.resetColMinMax(pt)

				sandmap.blockedMap[pt.String()] = "#"
			}
		}
	}

	toFile(sandmap.printMap(nil))

	return sandmap
}

func parsePoints(line string) []image.Point {
	var wall = make([]image.Point, 0)

	for _, ptStr := range strings.Split(line, " -> ") {
		var (
			crList = strings.Split(ptStr, ",")
			point  image.Point
		)
		if len(crList) != 2 {
			panic("bad parsing")
		} else if col, err := strconv.Atoi(crList[0]); err != nil {
			panic(err)
		} else if row, err := strconv.Atoi(crList[1]); err != nil {
			panic(err)
		} else {
			point.X = col
			point.Y = row
		}
		wall = append(wall, point)
	}

	return wall
}

func expandPoints(start, end image.Point) []image.Point {

	// sanity check
	if (start.X != end.X) && start.Y != end.Y {
		panic(fmt.Sprintf("shit got real: %v, %v", start, end))
	}

	var (
		list = make([]image.Point, 0)
		diff = end.Sub(start)
	)

	if diff.X != 0 {
		var stepX = diff.X / abs(diff.X)
		for x := start.X; ; x += stepX {
			list = append(list, image.Pt(x, start.Y))
			if x == end.X {
				break
			}
		}
	} else if diff.Y != 0 {
		var stepY = diff.Y / abs(diff.Y)
		for y := start.Y; ; y += stepY {
			list = append(list, image.Pt(start.X, y))
			if y == end.Y {
				break
			}
		}
	} else {
		panic(fmt.Sprintf("shit more real: start: %v end:%v diff:%v", start, end, diff))
	}

	// fmt.Printf("start:%v, end:%v\n%v\n", start, end, list)

	return list
}

type sandmap struct {
	iteration int

	firstAtBottomCount int
	sourceBlockedCount int

	sandStart  image.Point
	rowMax     int
	colMin     int
	colMax     int
	blockedMap map[string]string
}

func (sm *sandmap) resetColMinMax(pt image.Point) {
	if pt.X-1 < sm.colMin {
		sm.colMin = pt.X - 1
	}
	if pt.X+1 > sm.colMax {
		sm.colMax = pt.X + 1
	}

}

func (sm sandmap) printMap(curSand *image.Point) string {
	// todo: print headers
	var (
		outStr    string
		colMinStr = fmt.Sprintf("%d", sm.colMin)
		colMaxStr = fmt.Sprintf("%d", sm.colMax)
		sandStr   = "500"

		sandOff = 500 - sm.colMin + 1
		endOff  = sm.colMax - 500 + 1
	)

	// fmt.Printf("maxRow:%v, minCol:%v, maxCol:%v\n", sm.rowMax, sm.colMin, sm.colMax)
	outStr += fmt.Sprintf("iteration: %v\n\n", sm.iteration)
	for i := 0; i < 3; i++ {
		outStr += fmt.Sprintf("%5s%*s%*s\n",
			string(colMinStr[i]), sandOff, string(sandStr[i]), endOff, string(colMaxStr[i]))
	}
	outStr += fmt.Sprintf("    %s\n", strings.Repeat("_", sm.colMax-sm.colMin+3))

	for row := 0; row < sm.rowMax+1; row++ {
		outStr += fmt.Sprintf("%3d|", row)
		for col := sm.colMin - 1; col <= sm.colMax+1; col++ {
			var pt = image.Pt(col, row)

			if curSand != nil && pt.Eq(*curSand) {
				outStr += "@"
			} else if char, exists := sm.blockedMap[pt.String()]; exists {
				outStr += char
			} else if pt.Eq(sm.sandStart) {
				outStr += "+"
			} else {
				outStr += "."
			}
		}
		outStr += "\n"
	}
	outStr += fmt.Sprintf("    %s\n", strings.Repeat("#", sm.colMax-sm.colMin+3))

	return outStr
}

func (sm *sandmap) fillSand() {

	var finalResting image.Point
	for {
		finalResting = sm.dropOne()

		if sm.sourceBlocked(finalResting) {
			sm.iteration++
			sm.sourceBlockedCount = sm.iteration
			fmt.Printf("source blocked at %v, iteration %v", finalResting, sm.sourceBlockedCount)
			break
		} else {

			if sm.atBottom(finalResting) {
				sm.resetColMinMax(finalResting)
				if sm.firstAtBottomCount == 0 {
					sm.firstAtBottomCount = sm.iteration
					fmt.Printf("fillSand at bottom: %v\n", sm.firstAtBottomCount)
				}
			}

			sm.iteration++
			sm.blockedMap[finalResting.String()] = "O"
			var brk = 1000
			if test {
				brk = 20
			}
			if sm.iteration%brk == 0 {
				toFile(sm.printMap(nil))
				fmt.Printf("\nend iteration %v, blocked at %v", sm.iteration, finalResting)
				fmt.Print("\n")
			}
		}
	}

	fmt.Printf("filled iteration %v\n", sm.iteration)
	toFile(sm.printMap(&finalResting))
}

func (sm *sandmap) dropOne() image.Point {

	// sand start position
	var sand = sm.sandStart

	for { // full sand movement
		if movement, err := sm.moveSand(sand); err != nil {
			// fmt.Printf("\nend iteration %v, %v\n", sm.iteration, err)
			break
		} else {
			sand = sand.Add(movement)
			// test individual
			if sm.atBottom(sand) {
				fmt.Printf("\ndropOne at bottom: %v\n", sand)
				break
			}
			// toFile(sm.printMap(&sand))
		}
	}

	return sand
}

func (sm sandmap) atBottom(sand image.Point) bool {
	return sand.Y >= sm.rowMax
}

func (sm sandmap) sourceBlocked(sand image.Point) bool {
	return sand.Eq(sm.sandStart)
}

var (
	down      = image.Pt(0, 1)
	downLeft  = image.Pt(-1, 1)
	downRight = image.Pt(1, 1)
)

func (sm sandmap) moveSand(sand image.Point) (image.Point, error) {
	if _, exists := sm.blockedMap[sand.Add(down).String()]; exists == false {
		// fmt.Print("D")
		return down, nil
	} else if _, exists := sm.blockedMap[sand.Add(downLeft).String()]; exists == false {
		// fmt.Print("L")
		return downLeft, nil
	} else if _, exists := sm.blockedMap[sand.Add(downRight).String()]; exists == false {
		// fmt.Print("R")
		return downRight, nil
	} else {
		return image.Pt(0, 0), fmt.Errorf("blocked at %v", sand)
	}
}
