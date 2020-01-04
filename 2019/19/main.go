package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var verbose = false

type Memory map[int]int

func main() {
	initialMemory := parseInput("input.txt")

	// Part 1
	var part1View strings.Builder
	part1 := 0
	for y := 0; y < 50; y++ {
		for x := 0; x < 50; x++ {
			in, out := make(chan int), make(chan int)
			go runProgram(initialMemory, in, out)
			in <- x
			in <- y
			isPulled := <-out == 1
			if isPulled {
				part1++
				part1View.WriteByte('#')
			} else {
				part1View.WriteByte('.')
			}
		}
		part1View.WriteByte('\n')
	}
	fmt.Println(part1View.String())
	fmt.Println("Part 1:", part1)

	// Part 2
	desiredWidth := 100
	desiredHeight := 100
	prevXStart := 0
	heights := make(map[int]int)
	part2 := 0
LOOP:
	for y := 0; ; y++ {
		notPulled := 0
		width := 0
		for x := prevXStart; notPulled < 5; x++ {
			in, out := make(chan int), make(chan int)
			go runProgram(initialMemory, in, out)
			in <- x
			in <- y
			isPulled := <-out == 1
			if isPulled {
				if width == 0 {
					prevXStart = x
				}
				width++
				heights[x]++
				// this is the bottom right corner
				if heights[x] >= desiredHeight && width >= desiredWidth && x-desiredWidth+1 >= prevXStart {
					part2 = (x-desiredWidth+1)*10000 + (y - desiredHeight + 1)
					break LOOP
				}
			} else {
				notPulled++
			}
		}
	}
	fmt.Println("Part 2:", part2)
}

func runProgram(initialMemory Memory, input chan int, output chan<- int) {
	memory := make(Memory)
	for k, v := range initialMemory {
		memory[k] = v
	}
	relBase := 0

	for i := 0; ; {
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
			if verbose {
				fmt.Println("Program halted")
			}
			close(output)
			return
		default:
			close(output)
			panic(fmt.Sprintf("Op %d is not supported", op))
		}

		if verbose {
			fmt.Println(msg)
		}
	}

}

func getParams(memory Memory, n, pointer, relBase, numParams int) []int {
	params := make([]int, numParams)
	for i := range params {
		value := pointer + 1 + i
		mode := n % 10
		if mode == 0 {
			params[i] = memory[value]
		} else if mode == 1 {
			params[i] = value
		} else if mode == 2 {
			params[i] = relBase + memory[value]
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
		ints[i], _ = strconv.Atoi(strings.TrimSpace(fields[i]))
	}
	return ints
}
