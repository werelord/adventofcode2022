package main

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
)

// https://adventofcode.com/2022/day/7
// recursive fun!!

type Dir7 struct {
	totSize  uint64
	name     string
	parent   *Dir7
	dirList  map[string]*Dir7
	filelist map[string]*File7
}

func (cd *Dir7) AddDir(name string) {
	if _, exists := cd.dirList[name]; exists == false {
		var newDir = Dir7{
			totSize:  0,
			name:     name,
			parent:   cd,
			dirList:  make(map[string]*Dir7, 0),
			filelist: make(map[string]*File7, 0),
		}
		cd.dirList[name] = &newDir
		// fmt.Printf("%v:%v added dir\n", cd.name, name)

	} else {
		fmt.Printf("%v: dir already exists: %v\n", cd.name, name)
	}
}

func (cd *Dir7) WalkDown(name string) (*Dir7, error) {
	if dir, exists := cd.dirList[name]; exists == false {
		return nil, fmt.Errorf("%v: dir %v does not exists", cd.name, name)
	} else {
		// fmt.Printf("%v: walking down to %v\n", cd.name, dir.name)
		return dir, nil
	}
}

func (cd *Dir7) WalkUp() (*Dir7, error) {
	if cd.parent == nil {
		return nil, fmt.Errorf("%v: parent is nil (likely root?) cannot walk up", cd.name)
	} else {
		// fmt.Printf("%v: walking up to %v\n", cd.name, cd.parent.name)
		return cd.parent, nil
	}
}

func (cd *Dir7) AddFile(name string, size uint64) {
	if _, exists := cd.filelist[name]; exists == false {
		var newFile = File7{
			name: name,
			size: size,
		}
		cd.filelist[name] = &newFile
		// fmt.Printf("%v:%v (%v) added file\n", cd.name, newFile.name, newFile.size)

		// size calculation
		cd.PropegateFileSize(newFile.size)

	} else {
		fmt.Printf("%v: not adding file %v, already exists\n", cd.name, name)
	}
}

func (cd *Dir7) PropegateFileSize(size uint64) {
	// fmt.Printf("%v: adding child size: %v\n", cd.name, size)
	cd.totSize += size
	if cd.parent != nil { // if nil, we're at the root
		cd.parent.PropegateFileSize(size)
	}
}

type File7 struct {
	name string
	size uint64
}

// and a state machine!
type command int

const (
	noCommand command = iota
	moveRoot
	moveUp
	moveDown
	list
)

func day7(inp string) {
	file := filepath.Join(currentDir(), "input", inp)

	var (
		cmdCount int
		root     = Dir7{
			totSize:  0,
			name:     "/",
			parent:   nil,
			dirList:  make(map[string]*Dir7, 0),
			filelist: make(map[string]*File7, 0),
		}
		walker *Dir7
		inList bool
	)

	// test
	// var tcmd, tretstr = processCommand("$ cd dfmhjhd")
	// fmt.Printf("command %v, retstr %v", tcmd, tretstr)

	for line := range readLines(file) {
		cmdCount++

		var cmd, retstr = processCommand(line)

		switch cmd {
		case moveRoot:
			inList = false
			//fmt.Printf("moving root\n")
			walker = &root
		case moveUp:
			inList = false
			if dir, err := walker.WalkUp(); err != nil {
				panic(err)
			} else {
				walker = dir
			}
		case moveDown:
			inList = false
			if dir, err := walker.WalkDown(retstr); err != nil {
				panic(err)
			} else {
				walker = dir
			}
		case list:
			//fmt.Printf("ls (%v)\n", walker.name)
			inList = true
		case noCommand:
			// verify we're in list
			if inList == false {
				panic(fmt.Sprintf("no command, and not in list; line: %v", line))
			}
			// we're in a list; figure out file or directory
			if isDir, name := isDir(line); isDir {
				// fmt.Printf("dir found, name: %v\n", name)
				walker.AddDir(name)
			} else if isFile, name, size := isFile(line); isFile {
				// fmt.Printf("file found, name: %v size: %v\n", name, size)
				walker.AddFile(name, size)
			} else {
				panic(fmt.Sprintf("no command, not a dir or file: %v", line))
			}

		default:
			panic(fmt.Sprintf("unhandled command; line: %v", line))
		}
	}
	fmt.Printf("command count: %v\n", cmdCount)
	// recursive call figuring out dir size under limit
	var limit uint64 = 100000
	var tot = calcDirSizeUnder(&root, limit)
	fmt.Printf("total under %v: %v\n", limit, tot)

	const (
		unusedNeeded = 30000000
		totalSpace   = 70000000
	)
	var currentFree = totalSpace - root.totSize
	var amountNeeded = unusedNeeded - currentFree

	fmt.Printf("currentFree: %v, amountNeededToDelete: %v\n", currentFree, amountNeeded)
	var closest = findClosetTo(&root, amountNeeded)
	fmt.Printf("closest: %v, %v\n", closest.name, closest.totSize)

}

var (
	rxMoveDown = regexp.MustCompile(`\$ cd ([A-Za-z]+)`)
	rxDir      = regexp.MustCompile("dir ([A-Za-z]+)")
	rxFile     = regexp.MustCompile("([0-9]+) ([A-Za-z.]+)")
)

func processCommand(line string) (cmd command, ret string) {
	if line == "$ cd /" {
		cmd = moveRoot
	} else if line == "$ cd .." {
		cmd = moveUp
	} else if line == "$ ls" {
		cmd = list
	} else if moveDownMatch := rxMoveDown.FindStringSubmatch(line); len(moveDownMatch) == 2 {
		cmd = moveDown
		ret = moveDownMatch[1]
	} else {
		cmd = noCommand
		ret = line
	}
	return
}

func isDir(line string) (bool, string) {
	var dirMatch = rxDir.FindStringSubmatch(line)
	if len(dirMatch) != 2 {
		return false, ""
	} else {
		return true, dirMatch[1]
	}
}

func isFile(line string) (bool, string, uint64) {
	var fileMatch = rxFile.FindStringSubmatch(line)
	if len(fileMatch) != 3 {
		return false, "", 0
	} else if size, err := strconv.ParseUint(fileMatch[1], 10, 64); err != nil || size == 0 {
		panic(fmt.Sprintf("filematch found, but error or size when parsing; size:%v, error: %v", size, err))
	} else {
		return true, fileMatch[2], size
	}
}

func calcDirSizeUnder(dir *Dir7, limit uint64) (total uint64) {
	var tot uint64
	if dir.totSize < limit {
		// fmt.Printf("%v (%v) under limit %v\n", dir.name, dir.totSize, limit)
		tot = dir.totSize
	}
	for _, child := range dir.dirList {
		tot += calcDirSizeUnder(child, limit)
	}

	return tot
}
func findClosetTo(dir *Dir7, limit uint64) *Dir7 {

	// I know there's a better way of doing, this; propegating current smallest along is probably better
	// than bubbling smallest up from below; but meh

	if dir.totSize < limit {
		// won't satisfy the requirement..
		panic("wtf, this shouldn't happen")
		// return nil
	} else {
		var currentCandidate = dir
		// this one is over the limit and is a candidate.. check subdir for other candidates
		for _, subdir := range dir.dirList {
			// if child won't satisfy, just skip
			if subdir.totSize >= limit {
				if childCandidate := findClosetTo(subdir, limit); childCandidate == nil {
					panic("again, shouldn't happen")
				} else {
					// see if child is better than current
					var (
						childDiff = childCandidate.totSize - limit
						currDiff  = currentCandidate.totSize - limit
					)
					fmt.Printf("current (%v): %v (diff %v), child (%v): %v (diff %v)\n",
						currentCandidate.name, currentCandidate.totSize, currDiff,
						childCandidate.name, childCandidate.totSize, childDiff)
					if (childDiff) < (currDiff) {
						fmt.Printf("setting current to (%v): %v (diff %v)\n",
							childCandidate.name, childCandidate.totSize, childDiff)
						currentCandidate = childCandidate
					}
				}
			}
		}
		return currentCandidate
	}
}
