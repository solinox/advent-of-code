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

func play(start []int, finalTurn int) (int, map[int][]int) {
	mem := make(map[int][]int)
	for i, v := range start {
		mem[v] = append(mem[v], i+1)
	}
	prevN := start[len(start)-1]
	for turn := len(start) + 1; turn <= finalTurn; turn++ {
		prev, ok := mem[prevN]
		if !ok || len(prev) == 1 {
			prevN = 0
			mem[prevN] = append(mem[prevN], turn)
			continue
		}
		if len(prev) < 2 {
			log.Fatalln(prev, "Not long enough")
		}
		diff := prev[len(prev)-1] - prev[len(prev)-2]
		mem[diff] = append(mem[diff], turn)
		prevN = diff
	}
	return prevN, mem
}
