package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

var verbose = false

type Memory map[int]int

func main() {
	initialMemory := parseInput("input.txt")

	in, out := make(chan int), make(chan int)
	go runProgram(initialMemory, in, out)
	space := [][]byte{[]byte{}}
	row := 0
	for v := range out {
		b := byte(v)
		if b == '\n' {
			space = append(space, []byte{})
			row++
		} else {
			space[row] = append(space[row], b)
		}
	}

	// Part 1
	part1 := alignmentParameters(space)
	fmt.Println("Part 1:", part1)

	// Part 2
	// See instructions.txt for manual solution of input
	initialMemory[0] = 2
	in2, out2 := make(chan int), make(chan int)
	go runProgram(initialMemory, in2, out2)
	mainMovementRoutine := "A,C,A,C,B,B,C,A,C,B"
	aRoutine := "L,12,L,10,R,8,L,12"
	bRoutine := "L,10,R,12,R,8"
	cRoutine := "R,8,R,10,R,12"
	videoFeed := "n"
	printOutput(out2)
	inputASCII(mainMovementRoutine, in2)
	printOutput(out2)
	inputASCII(aRoutine, in2)
	printOutput(out2)
	inputASCII(bRoutine, in2)
	printOutput(out2)
	inputASCII(cRoutine, in2)
	printOutput(out2)
	inputASCII(videoFeed, in2)
	printOutput(out2)
}

func inputASCII(input string, in chan int) {
	for i := range input {
		in <- int(input[i])
	}
	in <- '\n'
}

func printOutput(out chan int) {
	for {
		select {
		case b := <-out:
			if int(byte(b)) == b {
				fmt.Print(string(b))
			} else if b == 0 {
				return
			} else {
				// too big for ASCII, just print the number
				fmt.Print(b)
			}
		case <-time.After(1 * time.Second):
			fmt.Print("\n")
			return
		}
	}
}

func alignmentParameters(space [][]byte) int {
	sum := 0
	for y := range space {
		if y == 0 || y == len(space)-1 {
			continue // no intersections on borders
		}
		for x := range space[y] {
			if x == 0 || x == len(space[y])-1 {
				continue // no intersections on borders
			}
			if space[y][x] == '#' && space[y][x-1] == '#' && space[y][x+1] == '#' && space[y-1][x] == '#' && space[y+1][x] == '#' {
				sum += x * y
			}
		}
	}
	return sum
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
