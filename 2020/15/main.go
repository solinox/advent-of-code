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

func play(start []int, finalTurn int) (int, map[int][2]int) {
	mem := make(map[int][2]int)
	for i, v := range start {
		mem[v] = [2]int{0, i + 1}
	}
	prevN := start[len(start)-1]
	for turn := len(start) + 1; turn <= finalTurn; turn++ {
		nums := mem[prevN]
		next := 0
		if nums[0] != 0 {
			next = nums[1] - nums[0]
		}
		nums = mem[next]
		nums[0], nums[1] = nums[1], turn
		mem[next] = nums
		prevN = next
	}
	return prevN, mem
}
