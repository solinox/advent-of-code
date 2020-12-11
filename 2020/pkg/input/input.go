package input

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

func ReadAll(r io.Reader) []byte {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}
	return data
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

func SplitString(r io.Reader, sep string) []string {
	data := string(ReadAll(r))
	return strings.Split(data, sep)
}

func IntMap(r io.Reader) map[int]int {
	scanner := bufio.NewScanner(r)
	m := make(map[int]int)
	for scanner.Scan() {
		line := scanner.Text()
		n, err := strconv.Atoi(line)
		if err != nil {
			log.Fatal(line, err)
		}
		m[n]++
	}
	return m
}

func Duplicate(r io.Reader) (io.Reader, io.Reader) {
	data := ReadAll(r)
	return bytes.NewReader(data), bytes.NewReader(data)
}
