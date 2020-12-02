package input

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
)

func File(filename string) *os.File {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	return file
}

func IntSlice(r io.Reader) []int {
	scanner := bufio.NewScanner(r)
	slice := make([]int, 0)
	for scanner.Scan() {
		line := scanner.Text()
		n, err := strconv.Atoi(line)
		if err != nil {
			log.Fatal(line, err)
		}
		slice = append(slice, n)
	}
	return slice
}

func StringSlice(r io.Reader) []string {
	scanner := bufio.NewScanner(r)
	slice := make([]string, 0)
	for scanner.Scan() {
		slice = append(slice, scanner.Text())
	}
	return slice
}
