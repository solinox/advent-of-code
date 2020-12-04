package input

import (
	"bufio"
	"bytes"
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

func BytesSlice(r io.Reader) [][]byte {
	scanner := bufio.NewScanner(r)
	slice := make([][]byte, 0)
	for scanner.Scan() {
		slice = append(slice, []byte(scanner.Text()))
	}
	return slice
}

func ScanDelim(r io.Reader, delim []byte) (<-chan string, <-chan string) {
	ch1, ch2 := make(chan string), make(chan string)
	scanner := bufio.NewScanner(r)
	// copied and tweaked from bufio.ScanLines
	split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		if i := bytes.Index(data, delim); i >= 0 {
			// We have a full newline-terminated line.
			return i + 1, dropCR(data[0:i]), nil
		}
		// If we're at EOF, we have a final, non-terminated line. Return it.
		if atEOF {
			return len(data), dropCR(data), nil
		}
		// Request more data.
		return 0, nil, nil
	}
	scanner.Split(split)
	go func() {
		for scanner.Scan() {
			ch1 <- scanner.Text()
			ch2 <- scanner.Text()
		}
		close(ch1)
		close(ch2)
	}()
	return ch1, ch2
}

// dropCR drops a terminal \r from the data.
func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}
