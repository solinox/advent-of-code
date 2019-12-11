package main

import (
	"strconv"
	"io/ioutil"
	"strings"
	"fmt"
)

func main() {
	memory := parseInput("input.txt")

	// Part 1
	part1 := runTEST(memory, 1)
	fmt.Println("Part 1", part1)

	// Part 2
	part2 := runTEST(memory, 5)
	fmt.Println("Part 2", part2)
}

func runTEST(initialMemory []int, input int) (output int) {
	memory := make([]int, len(initialMemory))
	copy(memory, initialMemory)

	for i := 0; i < len(memory); {
		n := memory[i]
		op := n % 100
		n /= 100

		switch op {
		case 1:
			a, b, c := memory[i+1], memory[i+2], memory[i+3]
			params := getParams(memory, n, a, b)
			memory[c] = params[0] + params[1]
			i += 4
		case 2:
			a, b, c := memory[i+1], memory[i+2], memory[i+3]
			params := getParams(memory, n, a, b)
			memory[c] = params[0] * params[1]
			i += 4
		case 3:
			a := memory[i+1]
			memory[a] = input
			i += 2
		case 4:
			a := memory[i+1]
			output = memory[a]
			i += 2
		case 5:
			a, b := memory[i+1], memory[i+2]
			params := getParams(memory, n, a, b)
			if params[0] != 0 {
				i = params[1]
			} else {
				i += 3
			}
		case 6:
			a, b := memory[i+1], memory[i+2]
			params := getParams(memory, n, a, b)
			if params[0] == 0 {
				i = params[1]
			} else {
				i += 3
			}
		case 7:
			a, b, c := memory[i+1], memory[i+2], memory[i+3]
			params := getParams(memory, n, a, b)
			if params[0] < params[1] {
				memory[c] = 1
			} else {
				memory[c] = 0
			}
			i += 4
		case 8:
			a, b, c := memory[i+1], memory[i+2], memory[i+3]
			params := getParams(memory, n, a, b)
			if params[0] == params[1] {
				memory[c] = 1
			} else {
				memory[c] = 0
			}
			i += 4
		case 99:
			return output
		}
	}
	fmt.Println("Program did not halt but went out of bounds unexpectedly")
	return output
}

func getParams(memory []int, n int, params ...int) []int {
	for i := range params {
		mode := n % 10
		if mode == 0 {
			params[i] = memory[params[i]]
		} else if mode == 1 {
			params[i] = params[i]
		}
		n /= 10
	}
	return params
}

func parseInput(filename string) []int {
	data, _ := ioutil.ReadFile(filename)
	fields := strings.Split(string(data), ",")
	ints := make([]int, len(fields))
	for i := range fields {
		ints[i], _ = strconv.Atoi(fields[i])
	}
	return ints
}
