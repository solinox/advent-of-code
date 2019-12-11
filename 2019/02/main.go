package main

import (
	"strconv"
	"io/ioutil"
	"strings"
	"fmt"
)

func main() {
	ints := parseInput("input.txt")

	// Part 1
	part1 := runProgram(ints, 12, 2)
	fmt.Println("Part 1", part1) // 3654868

	// Part 2
	desiredOutput := 19690720
	for noun := 0; noun < 100; noun++ {
		for verb := 0; verb < 100; verb++ {
			part2 := runProgram(ints, noun, verb)
			if desiredOutput == part2 {
				fmt.Println("Part 2", 100 * noun + verb)
				return
			}
		}
	}
	fmt.Println("No solution found for part 2")
}

func runProgram(memory []int, noun, verb int) int {
	ints := make([]int, len(memory))
	copy(ints, memory)
	ints[1] = noun
	ints[2] = verb
	for i := 0; i < len(ints); i += 4 {
		// fmt.Printf("Position %d = %d. ", i, ints[i])
		if ints[i] != 1 && ints[i] != 2 {
			break
		}
		a, b, c := ints[i+1], ints[i+2], ints[i+3]
		if ints[i] == 1 {
			ints[c] = ints[a] + ints[b]
			// fmt.Printf(" %d + %d = %d, stored in Position %d\n", ints[a], ints[b], ints[c], c)
		}
		if ints[i] == 2 {
			ints[c] = ints[a] * ints[b]
			// fmt.Printf(" %d * %d = %d, stored in Position %d\n", ints[a], ints[b], ints[c], c)
		}
	}
	return ints[0]
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
