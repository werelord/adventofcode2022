package main

import (
	"fmt"
	"path/filepath"
)

// https://adventofcode.com/2022/day/8

// this one's a PITA; read the problem carefully

type Visibility uint8

const (
	NotVisible  = 0
	VisibleLeft = 1 << iota
	VisibleRight
	VisibleUp
	VisibleDown
)

type tree struct {
	height int
	vis    Visibility
	score  int
}

func (t *tree) String() string {
	return fmt.Sprintf("%d", t.height)
}

func day8(inp string) {
	file := filepath.Join(currentDir(), "input", inp)

	var (
		forest = make([][]*tree, 0) // rows
	)

	// build the forest
	for line := range readLines(file) {

		var row = make([]*tree, 0, len(line))

		for _, r := range line {
			var tree = tree{
				height: int(r - '0'),
				vis:    NotVisible,
			}
			row = append(row, &tree)
		}

		forest = append(forest, row)
	}

	// I don't care about efficency of cycles, just get it done.. as I don't care about memory efficiency either

	redo(forest)
}

// func part1(forest [][]*tree) {
// 	// scan from left
// 	for row := 0; row < 99; row++ {
// 		var localMax = -1
// 		for col := 0; col < 99; col++ {
// 			var tree = forest[row][col]
// 			if tree.height > localMax {
// 				localMax = tree.height
// 				tree.vis = tree.vis | VisibleLeft
// 			}
// 		}
// 	}

// 	// scan from right
// 	for row := 0; row < 99; row++ {
// 		var localMax = -1
// 		for col := 98; col >= 0; col-- {
// 			var tree = forest[row][col]
// 			if tree.height > localMax {
// 				localMax = tree.height
// 				tree.vis = tree.vis | VisibleRight
// 			}
// 		}
// 	}

// 	// scan from up
// 	for col := 0; col < 99; col++ {
// 		var localMax = -1
// 		for row := 0; row < 99; row++ {
// 			var tree = forest[row][col]
// 			if tree.height > localMax {
// 				localMax = tree.height
// 				tree.vis = tree.vis | VisibleUp
// 			}
// 		}
// 	}

// 	// scan from down
// 	for col := 0; col < 99; col++ {
// 		var localMax = -1
// 		for row := 98; row >= 0; row-- {
// 			var tree = forest[row][col]
// 			if tree.height > localMax {
// 				localMax = tree.height
// 				tree.vis = tree.vis | VisibleDown
// 			}
// 		}
// 	}

// 	var (
// 		visbileCount, hiddenCount int
// 	)

// 	// finally, count visible
// 	for x := 0; x < 99; x++ {
// 		for y := 0; y < 99; y++ {
// 			var tree = forest[x][y]
// 			if tree.vis > 0 {
// 				visbileCount++
// 			} else {
// 				hiddenCount++
// 			}
// 		}
// 	}

// 	fmt.Printf("hidden: %v, Visible: %v\n", hiddenCount, visbileCount)
// 	fmt.Println("done")
// }

func redo(forest [][]*tree) {

	var (
		cols = make([][]*tree, 0, 99)
		rows = make([][]*tree, 0, 99)
	)

	for x := 0; x < len(forest); x++ {
		rows = append(rows, forest[x])
	}

	for y := 0; y < len(forest); y++ {
		var newcol = make([]*tree, 0, len(forest))
		for x := 0; x < len(forest); x++ {
			newcol = append(newcol, forest[x][y])
		}
		cols = append(cols, newcol)
	}

	var (
		visbileCount, hiddenCount int
		maxScore                  = 0
	)

	for x := 0; x < len(forest); x++ {
		for y := 0; y < len(forest); y++ {
			var (
				tree        = rows[x][y]
				left, right = splittrees(rows[x], y)
				up, down    = splittrees(cols[y], x)
			)
			if visbileFrom(tree.height, left) {
				tree.vis |= VisibleLeft
			}
			if visbileFrom(tree.height, right) {
				tree.vis |= VisibleRight
			}
			if visbileFrom(tree.height, up) {
				tree.vis |= VisibleUp
			}
			if visbileFrom(tree.height, down) {
				tree.vis |= VisibleDown
			}

			if tree.vis > 0 {
				visbileCount++
			} else {
				hiddenCount++
			}

			// score trees
			if x == 0 || y == 0 || x == len(forest)-1 || y == len(forest)-1 {
				tree.score = 0
			} else {

				var sc = 1
				// left/up are reversed, right/down are not
				sc *= score(tree.height, left, true)
				sc *= score(tree.height, right, false)
				sc *= score(tree.height, up, true)
				sc *= score(tree.height, down, false)

				// fmt.Println("totscore: ", sc)

				tree.score = sc
			}

			if tree.score > maxScore {
				maxScore = tree.score
			}
		}
	}

	fmt.Printf("hidden: %v, Visible: %v\nmaxscore:%v\n", hiddenCount, visbileCount, maxScore)

	fmt.Println("done")

}

func splittrees(list []*tree, index int) ([]*tree, []*tree) {
	return list[:index], list[index+1:]
}

func visbileFrom(height int, list []*tree) bool {

	for _, tree := range list {
		if tree.height >= height {
			//fmt.Printf("not vis, h:%d, trees:%s\n", height, printTrees(list))
			return false
		}
	}
	// fmt.Printf("vis, h:%d, trees:%s\n", height, printTrees(list))
	return true
}

func printTrees(list []*tree) string {
	var ret string
	for _, tree := range list {
		ret += tree.String()
	}
	return ret
}

func score(curHeight int, list []*tree, reverse bool) int {
	if len(list) == 0 {
		return 0
	}

	var maxDistance = 0

	var doreverse = func() int {
		var dist int
		// var treestr = fmt.Sprintf("h%d:", curHeight)
		// var carstr = "   "
		for x := len(list) - 1; x >= 0; x-- {
			var tree = list[x]
			// treestr += fmt.Sprintf("%d", tree.height)
			if dist == 0 && tree.height >= curHeight {
				// carstr += "^"
				// because fucking reverse
				dist = len(list) - x
				// } else {
				// carstr += " "
			}
		}

		if dist == 0 {
			dist = len(list)
		}
		// fmt.Printf("%v\n%v score:%v\n", treestr, carstr, dist)

		return dist
	}
	var donorm = func() int {
		var dist int
		// var treestr = fmt.Sprintf("h%d:", curHeight)
		// var carstr = "   "
		for x, tree := range list {
			// treestr += fmt.Sprintf("%d", tree.height)
			if dist == 0 && tree.height >= curHeight {
				// carstr += "^"
				dist = x + 1
				// } else {
				// carstr += " "
			}
		}
		if dist == 0 {
			dist = len(list)
		}
		// fmt.Printf("%v\n%v score:%v\n", treestr, carstr, dist)

		return dist
	}

	if reverse {
		maxDistance = doreverse()
	} else {
		maxDistance = donorm()
	}

	return maxDistance
}
