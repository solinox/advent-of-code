package main

import (
	"log"
	"os"
	"time"

	"github.com/solinox/advent-of-code/2020/pkg/input"
)

func main() {
	start := input.SplitInt(os.Stdin, ",")

	t0 := time.Now()
	log.Println(part1(start), time.Since(t0))

	t0 = time.Now()
	log.Println(part2(start), time.Since(t0))
}

func part1(start []int) int {
	n, _ := play(start, 2020)
	return n
}

func part2(start []int) int {
	n, _ := play(start, 30000000)
	return n
}

func play(start []int, finalTurn int) (int, map[int]int) {
	mem := make(map[int]int)
	for i, v := range start[:len(start)-1] {
		mem[v] = i + 1
	}
	prevN := start[len(start)-1]
	for turn := len(start) + 1; turn <= finalTurn; turn++ {
		last := mem[prevN]
		mem[prevN] = turn - 1
		next := 0
		if last != 0 {
			next = turn - 1 - last
		}
		prevN = next
	}
	return prevN, mem
}
