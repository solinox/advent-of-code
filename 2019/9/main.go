package main

import (
	"strconv"
	"io/ioutil"
	"strings"
	"fmt"
)

var verbose = false

type Memory map[int]int

func main() {
	memory := parseInput("input.txt")

	// Part 1
	in1, out1 := make(chan int), make(chan int)
	go runProgram(memory, in1, out1)
	in1 <- 1
	part1 := make([]int, 0)
	for output := range out1 {
		part1 = append(part1, output)
	}
	fmt.Println("Part 1", part1)

	// Part 2
	in2, out2 := make(chan int), make(chan int)
	go runProgram(memory, in2, out2)
	in2 <- 2
	part2 := make([]int, 0)
	for output := range out2 {
		part2 = append(part2, output)
	}
	fmt.Println("Part 2", part2)

}

func runProgram(initialMemory Memory, input <-chan int, output chan<- int) {
	memory := make(Memory)
	for k, v := range initialMemory {
		memory[k] = v
	}
	relBase := 0

	for i := 0;; {
		n := memory[i]
		op := n % 100
		n /= 100

		msg := fmt.Sprintf("i=%d\tOp=%d\tRel=%d\t", i, memory[i], relBase)

		switch op {
		case 1:
			p := getParams(memory, n, i, relBase, 3)
			memory[p[2]] = memory[p[0]] + memory[p[1]]
			msg += fmt.Sprintf("memory[%d] = memory[%d] + memory[%d] ==> %d", p[2], p[0], p[1], memory[p[2]])
			i += 4
		case 2:
			p := getParams(memory, n, i, relBase, 3)
			memory[p[2]] = memory[p[0]] * memory[p[1]]
			msg += fmt.Sprintf("memory[%d] = memory[%d] * memory[%d] ==> %d", p[2], p[0], p[1], memory[p[2]])
			i += 4
		case 3:
			p := getParams(memory, n, i, relBase, 1)
			memory[p[0]] = <-input
			msg += fmt.Sprintf("memory[%d] = input ==> %d", p[0], memory[p[0]])
			i += 2
		case 4:
			p := getParams(memory, n, i, relBase, 1)
			output <- memory[p[0]]
			msg += fmt.Sprintf("output = memory[%d] ==> %d", p[0], memory[p[0]])
			i += 2
		case 5:
			p := getParams(memory, n, i, relBase, 2)
			msg += fmt.Sprintf("memory[%d] != 0 ==> %d != 0. ", p[0], memory[p[0]])
			if memory[p[0]] != 0 {
				i = memory[p[1]]
				msg += fmt.Sprintf("Set i to memory[%d] = %d", p[1], memory[p[1]])
			} else {
				i += 3
				msg += "do nothing but increment pointer"
			}
		case 6:
			p := getParams(memory, n, i, relBase, 2)
			msg += fmt.Sprintf("memory[%d] == 0 ==> %d == 0. ", p[0], memory[p[0]])
			if memory[p[0]] == 0 {
				i = memory[p[1]]
				msg += fmt.Sprintf("Set i to memory[%d] = %d", p[1], memory[p[1]])
			} else {
				i += 3
				msg += "do nothing but increment pointer"
			}
		case 7:
			p := getParams(memory, n, i, relBase, 3)
			msg += fmt.Sprintf("memory[%d] < memory[%d] ==> %d < %d. ", p[0], p[1], memory[p[0]], memory[p[1]])
			if memory[p[0]] < memory[p[1]] {
				memory[p[2]] = 1
				msg += fmt.Sprintf("Set memory[%d] to 1", p[2])
			} else {
				memory[p[2]] = 0
				msg += fmt.Sprintf("Set memory[%d] to 0", p[2])
			}
			i += 4
		case 8:
			p := getParams(memory, n, i, relBase, 3)
			msg += fmt.Sprintf("memory[%d] == memory[%d] ==> %d == %d. ", p[0], p[1], memory[p[0]], memory[p[1]])
			if memory[p[0]] == memory[p[1]] {
				memory[p[2]] = 1
				msg += fmt.Sprintf("Set memory[%d] to 1", p[2])
			} else {
				memory[p[2]] = 0
				msg += fmt.Sprintf("Set memory[%d] to 0", p[2])
			}
			i += 4
		case 9:
			p := getParams(memory, n, i, relBase, 1)
			msg += fmt.Sprintf("Add memory[%d] = %d to old relBase = %d.", p[0], memory[p[0]], relBase)
			relBase += memory[p[0]]
			i += 2
		case 99:
			close(output)
			return
		default:
			panic(fmt.Sprintf("Op %d is not supported", op))
		}

		if verbose {
			fmt.Println(msg)
		}
	}
	fmt.Println("Program did not halt but went out of bounds unexpectedly")
	close(output)
	return
}

func getParams(memory Memory, n, pointer, relBase, numParams int) []int {
	params := make([]int, numParams)
	for i := range params {
		value := pointer+1+i
		mode := n % 10
		if mode == 0 {
			params[i] = memory[value]
		} else if mode == 1 {
			params[i] = value
		} else if mode == 2 {
			params[i] = relBase+memory[value]
		}
		n /= 10
	}
	return params
}

func parseInput(filename string) Memory {
	data, _ := ioutil.ReadFile(filename)
	fields := strings.Split(string(data), ",")
	ints := make(Memory)
	for i := range fields {
		ints[i], _ = strconv.Atoi(fields[i])
	}
	return ints
}
