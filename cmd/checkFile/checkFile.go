package main

import (
	"fmt"
	"os"

	"github.com/d24wang/go-file-checker/util"
)

var badLines = []*badLine{}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("%s file", os.Args[0])
		os.Exit(1)
	}

	filename := os.Args[1]

	file, err := os.Open(filename)
	checkError(err)

	breaker := util.NewBytesLineBreaker(file)
	defer breaker.Close()
	defer printResult()

	breakerErr := breaker.ErrorChan()
	lines := breaker.LinesChan()
	lineNum := 1

	for {
		select {
		case err = <-breakerErr:
			checkError(err)
		case line, ok := <-lines:
			if good, pos := util.ValidLine(line); !good {
				badLines = append(badLines, &badLine{
					Line: lineNum,
					Pos:  pos,
				})
			}
			if !ok {
				return
			}
			lineNum++
		}
	}
}

func printResult() {
	if len(badLines) == 0 {
		fmt.Println("All good!")
		return
	}

	for _, line := range badLines {
		fmt.Println(line)
	}
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

type badLine struct {
	Line int
	Pos  int
}

func (l *badLine) String() string {
	return fmt.Sprintf("Possible non-unicode character at Ln %v, Col %v", l.Line, l.Pos)
}
