package main

import (
	"bufio"
	"log"
	"os"
)

func main() {
	l := log.New(os.Stderr, "", 0)

	if len(os.Args) < 2 {
		l.Println("filename not provided")
		os.Exit(1)
	}
	filename := os.Args[1]
	file, err := os.Open(filename)
	if err != nil {
		l.Println(err)
		os.Exit(1)
	}

	lines := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		l.Println(err)
		os.Exit(1)
	}
	err = file.Close()
	if err != nil {
		l.Println(err)
	}

	task, err := parseLines(lines)
	if err != nil {
		l.Println(err)
		os.Exit(1)
	}

	err = task.processCells()
	if err != nil {
		l.Println(err)
		os.Exit(1)
	}

	task.printCells()
}
