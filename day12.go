package main

import (
	"fmt"
	"image"
	"path/filepath"
	"strings"

	"github.com/fzipp/astar"
)

// https://adventofcode.com/2022/day/12

// thought 1: never implemented A* algorithm before, should be interesting..
// thought 2: well shit, A* is just breath first search with a priority queue?? wonder if I should implement my own.. meh

// nope: do a naive implementation (Dijkstra) rather than a*.. compare diff, because thats more fun

func day12() {

	var elevationMap, start, end = buildMap(false)

	path := astar.FindPath[image.Point](elevationMap, start, end, nodeDist, nodeDist)

	fmt.Printf("path part 1 found: len:%v\n", len(path)-1) // subtract start

	// part two, find the shortest path from any a
	var pathMap = make(map[image.Point]int, 0)

	pathMap[start] = len(path) - 1
	var minPt = start

	for row, str := range elevationMap {
		for col, val := range []rune(str) {
			if val == 'a' {
				var newStart = image.Pt(row, col)
				if _, exists := pathMap[newStart]; exists == false {
					newPath := astar.FindPath[image.Point](elevationMap, newStart, end, nodeDist, nodeDist)
					if len(newPath) > 0 {
						// fmt.Printf("checking point %v, pathlen: %v\n", newStart, len(newPath) - 1)
						pathMap[newStart] = len(newPath) - 1

						if pathMap[newStart] < pathMap[minPt] {
							minPt = newStart
						}
					}
				}
			}
		}
	}

	fmt.Printf("minimum point found: %v len: %v", minPt, pathMap[minPt])
}

func nodeDist(p, q image.Point) float64 {
	d := q.Sub(p)
	return float64(abs(d.X) + abs(d.Y))
}

func buildMap(test bool) (graph, image.Point, image.Point) {

	var (
		rowIdx     = 0
		file       = filepath.Join(currentDir(), "input", "day12.txt")
		m          = make([]string, 0)
		start, end image.Point
	)
	if test {
		file = filepath.Join(currentDir(), "input", "day12.test.txt")
	}
	for line := range readLines(file) {

		if idy := strings.IndexRune(line, 'S'); idy >= 0 {
			start = image.Pt(rowIdx, idy)
			// replace start with its default elevation 'a'
			line = line[:idy] + "a" + line[idy+1:]
		}

		if idy := strings.IndexRune(line, 'E'); idy >= 0 {
			end = image.Pt(rowIdx, idy)
			// replace end with its default elevation 'z'
			line = line[:idy] + "z" + line[idy+1:]
		}

		m = append(m, line)

		rowIdx++
	}

	return m, start, end
}

type graph []string

var (
	// with 0,0 being top left: north, south, east west
	// X is row, Y is column
	ptOffsets    = []image.Point{image.Pt(-1, 0), image.Pt(1, 0), image.Pt(0, 1), image.Pt(0, -1)}
	inspectPoint = image.Pt(23, 151)
)

// Neighbours implements the astar.Graph[Node] interface (with Node = image.Point).
func (g graph) Neighbours(currentPt image.Point) []image.Point {

	res := make([]image.Point, 0, 4)
	for _, offset := range ptOffsets {
		newPt := currentPt.Add(offset)
		if currentPt.Eq(inspectPoint) {
			fmt.Print()
		}
		if g.canReach(currentPt, newPt) {
			res = append(res, newPt)
		}
	}

	return res
}

var runemap = make(map[rune]rune, 0)

func (g graph) canReach(currentPt, newPt image.Point) bool {
	// check in bounds
	if newPt.X < 0 || newPt.Y < 0 {
		return false
	} else if newPt.X >= len(g) || newPt.Y >= len(g[0]) {
		return false
	}

	var (
		curRune = rune(g[currentPt.X][currentPt.Y])
		newRune = rune(g[newPt.X][newPt.Y])
	)
	if newRune <= curRune { // can drop down infinitely
		return true
	} else if (newRune - curRune) == 1 { // can move up by one
		return true
	} else {
		return false
	}
}
