package main

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

// https://adventofcode.com/2022/day/9

type knot struct {
	coord

	name     string
	next     *knot
	moveList map[string]int
}

type coord struct {
	x, y int
}

func (c knot) String() string {
	return fmt.Sprintf("%v(%v,%v)", c.name, c.x, c.y)
}
func (c coord) String() string {
	return fmt.Sprintf("(%v,%v)", c.x, c.y)
}

func (c knot) recordMove() {
	if c.moveList == nil {
		return
	}
	// record new position
	if _, exists := c.moveList[c.String()]; exists {
		c.moveList[c.String()]++
	} else {
		c.moveList[c.String()] = 1
	}
}

// func (c coord) isDifferent(new coord) bool {
// 	return (c.x != new.x) || (c.y != new.y)
// }

func (c coord) isDiagonalMove(new coord) bool {
	// head and tail are in completely different row & column, need to move diag to keep up
	return (c.x != new.x) && (c.y != new.y)
}

func (c coord) shouldMove(head coord) bool {
	dx, dy := c.getDiff(head)
	return abs(dx) > 1 || abs(dy) > 1
}

func (c knot) moveLinear(head coord) coord {
	// either row or column will be different, by 2.. other should be 0
	dx, dy := c.getDiff(head)
	// fmt.Print("L")
	if dx != 0 && dy != 0 {
		panic(fmt.Sprintf("move linear, but both dx & dy are not zero: dx%v dy%v", dx, dy))
	} else {
		// non-diff will be zero, no move
		// negatvie dx or dy will end up a subtration, positive an addition
		var n = coord{}
		if abs(dx) > 0 {
			n.x = c.x + dx/2
			n.y = c.y
		} else { // move y
			n.x = c.x
			n.y = c.y + dy/2

		}
		// fmt.Printf("lin: head%v, cur%v, dx(%v,%v), new%v\n", head, c, dx, dy, n)
		return n
	}
}

func (c knot) moveDiagonal(head coord) coord {
	// both x and y should be > 0
	// fmt.Print("D")
	dx, dy := c.getDiff(head)
	if dx == 0 || dy == 0 {
		panic(fmt.Sprintf("move diag, but both dx & dy are not zero: dx%v dy%v", dx, dy))
	} else {
		var n = coord{}
		// diagonal move, either abs(dx) == 2 and/or abs(dy) == 2.. the other coord should be 1

		// dividing by 2 keeps sign
		if abs(dx) == 2 {
			dx = dx / 2
		}
		if abs(dy) == 2 {
			dy = dy / 2
		}
		n.x = c.x + dx
		n.y = c.y + dy
		// fmt.Printf("dlg: head:%v, cur:%v, dx(%v,%v), new%v\n", head, c, dx, dy, n)
		return n
	}
}

func (c coord) getDiff(head coord) (dx, dy int) {
	return head.x - c.x, head.y - c.y
}

func day9(inp string) {

	file := filepath.Join(currentDir(), "input", inp)

	var (
		head        = &knot{name: "head"}
		currentKnot = head
	)

	for i := 1; i < 10; i++ {
		var newtail = knot{name: fmt.Sprintf("%v", i)}

		if i == 1 || i == 9 {
			// record movements of first (part1) and last (part2)
			newtail.moveList = make(map[string]int, 0)
			// make sure origin is in movelist
			newtail.recordMove()
		}
		currentKnot.next = &newtail
		currentKnot = &newtail
	}

	test(head)

	for line := range readLines(file) {
		handleMovement(line, head)
	}

	fmt.Printf("done, tail movements")
	currentKnot = head
	for {
		if currentKnot.moveList != nil {
			fmt.Printf(", knot%v: %v", currentKnot.name, len(currentKnot.moveList))
		}

		if currentKnot.next == nil {
			break
		} else {
			currentKnot = currentKnot.next
		}
	}
	fmt.Print("\n")
}

func handleMovement(line string, head *knot) {
	var movement = strings.Split(line, " ")
	if len(movement) != 2 {
		panic(fmt.Sprint("wtf, bad movement: ", movement))
	}

	var iter, err = strconv.Atoi(movement[1])
	if err != nil {
		panic(err)
	}

	// just handle one for now
	switch move := movement[0]; move {
	case "U":
		moveHead(head, 0, 1, iter)
	case "D":
		moveHead(head, 0, -1, iter)
	case "L":
		moveHead(head, -1, 0, iter)
	case "R":
		moveHead(head, 1, 0, iter)
	default:
		panic(fmt.Sprint("bad movement: ", move))
	}
}

func moveHead(head *knot, x, y, iter int) {
	// fmt.Printf("moving head (x%vy%v)*%v:\n", x, y, iter)
	for it := 0; it < iter; it++ {
		var prevHead = head.coord
		head.x += x
		head.y += y

		// fmt.Printf("moved:%v", head)
		// should recursively move down tail.. this is always head, will always have tail
		var tailMoved = head.next.headMoved(head.coord)

		fmt.Printf("\n")

		// sanity check for tail; should really be previous head
		if tailMoved && prevHead != head.next.coord {
			panic("wtf, we're off base")
		}

	}
}

func (c *knot) headMoved(newHead coord) bool {

	var currentMoved bool
	if currentMoved = c.shouldMove(newHead); currentMoved {
		// move tail
		var newCoord coord
		if c.isDiagonalMove(newHead) {
			newCoord = c.moveDiagonal(newHead)
		} else {
			newCoord = c.moveLinear(newHead)
		}

		// sanity check before assigning the move
		{
			dx, dy := newCoord.getDiff(newCoord)
			if abs(dx) > 1 || abs(dy) > 1 {
				panic("wtf, we're off base")
			}
		}
		c.coord = newCoord

		// record move if applicable
		c.recordMove()

		// fmt.Printf("%v", c)

		// } else {
		// fmt.Printf("no move needed, %v, %v\n", head, tail)
	}

	// we only need to propegate down if this one moved
	if currentMoved && c.next != nil {
		// move the tail
		c.next.headMoved(c.coord)

	}
	return currentMoved
}

func test(head *knot) {

	// var checkAfter = func(head *knot, after ...coord) {

	// 	cur := head
	// 	for _, c := range after {
	// 		if cur.coord != c {
	// 			fmt.Printf("%v: not good, exp:%v, got %v\n", cur.name, c, cur.coord)
	// 		}
	// 		cur = cur.next
	// 	}
	// }

	// handleMovement("R 4", head)
	// var after = []coord{{4, 0}, {3, 0}, {2, 0}, {1, 0}, {0, 0}}
	// checkAfter(head, after...)

	// handleMovement("U 4", head)
	// after = []coord{{4, 4}, {4, 3}, {4, 2}, {3, 2}, {2, 2}, {1, 1}, {0, 0}}
	// checkAfter(head, after...)

	// handleMovement("L 3", head)
	// after = []coord{{1, 4}, {2, 4}, {3, 3}, {3, 2}, {2, 2}, {1, 1}, {0, 0}}
	// checkAfter(head, after...)

	// handleMovement("D 1", head)
	// after = []coord{{1, 3}, {2, 4}, {3, 3}, {3, 2}, {2, 2}, {1, 1}, {0, 0}}
	// checkAfter(head, after...)

	// handleMovement("R 4", head)
	// after = []coord{{5, 3}, {4, 3}, {3, 3}, {3, 2}, {2, 2}, {1, 1}, {0, 0}}
	// checkAfter(head, after...)

	// os.Exit(0)

}
