package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/solinox/advent-of-code/2020/pkg/input"
)

func main() {
	instructions := input.StringSlice(os.Stdin)

	t0 := time.Now()
	log.Println(part1(instructions), time.Since(t0))

	t0 = time.Now()
	log.Println(part2(instructions), time.Since(t0))
}

func part1(instructions []string) int {
	_, acc := boot(instructions)
	return acc
}

// brute force
func part2(instructions []string) int {
	r := strings.NewReplacer("nop", "jmp", "jmp", "nop")
	for i := range instructions {
		if instructions[i][0] == 'a' { // acc
			continue
		}
		// try replacing this instruction
		orig := instructions[i]
		instructions[i] = r.Replace(instructions[i])
		if ok, acc := boot(instructions); ok {
			return acc
		} else {
			// revert back to original
			instructions[i] = orig
		}
	}
	log.Fatalln("Out of range")
	return 0
}

// returns bool if exited successfully
// also returns the accumulator value when exiting or right before an infinite loop
func boot(instructions []string) (bool, int) {
	index, accumulator := 0, 0
	done := make(map[int]struct{})
	for {
		if index == len(instructions) {
			break
		}
		if index > len(instructions) {
			return false, accumulator
		}
		if _, ok := done[index]; ok {
			return false, accumulator
		}
		done[index] = struct{}{}
		index, accumulator = doInstruction(instructions[index], index, accumulator)
	}
	return true, accumulator
}

func doInstruction(instruction string, index, acc int) (int, int) {
	var op string
	var n int
	if m, err := fmt.Sscanf(instruction, "%s %d", &op, &n); m != 2 || err != nil {
		log.Fatalln("Invalid scan", instruction, m, err)
	}
	switch op {
	case "nop":
		return index + 1, acc
	case "acc":
		return index + 1, acc + n
	case "jmp":
		return index + n, acc
	}
	log.Fatalln("Unknown instruction", op)
	return 0, 0
}
