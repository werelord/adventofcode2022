package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

func readFile(filename string) chan string {

	var out = make(chan string, 1)
	go func() {
		file, err := os.Open(filename)
		if err != nil {
			fmt.Print("error opening file: ", err)
			close(out)
			return
		}
		defer file.Close()
		defer close(out)

		var fileScanner = bufio.NewScanner(file)

		for fileScanner.Scan() {
			out <- fileScanner.Text()
		}
		//fmt.Println("file done")
	}()
	return out
}

func currentDir() string {
	if ex, err := os.Executable(); err != nil {
		panic(err)
	} else {
		return filepath.Dir(ex)
	}

}
