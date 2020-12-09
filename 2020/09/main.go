package main

import (
	"log"
	"os"
	"time"

	"github.com/solinox/advent-of-code/2020/pkg/input"
	"github.com/solinox/advent-of-code/2020/pkg/sliceutils"
)

func main() {
	xmas := input.IntSlice(os.Stdin)

	t0 := time.Now()
	log.Println(part1(xmas), time.Since(t0))

	t0 = time.Now()
	log.Println(part2(xmas), time.Since(t0))
}

func part1(xmas []int) int {
	preambleSums := preamble(xmas, 25)
	for i := 25; i < len(xmas); i++ {
		if preambleSums[xmas[i]] == 0 {
			return xmas[i]
		}
		preambleSums = next(xmas, 25, i, preambleSums)
	}
	log.Fatalln("Not found")
	return 0
}

func preamble(xmas []int, preambleLen int) map[int]int {
	preambleSums := make(map[int]int)
	if preambleLen >= len(xmas) {
		log.Fatalln("Out of range")
		return preambleSums
	}
	for i := 0; i < preambleLen; i++ {
		for j := i + 1; j < preambleLen; j++ {
			sum := xmas[i] + xmas[j]
			preambleSums[sum]++
		}
	}
	return preambleSums
}

func next(xmas []int, preambleLen, index int, preambleSums map[int]int) map[int]int {
	if index >= len(xmas) || index-preambleLen < 0 {
		log.Fatalln("Out of range")
		return preambleSums
	}

	toRemove, toAdd := xmas[index-preambleLen], xmas[index]
	for i := index - preambleLen + 1; i < index; i++ {
		preambleSums[xmas[i]+toRemove]--
		preambleSums[xmas[i]+toAdd]++
	}
	return preambleSums
}

func part2(xmas []int) int {
	n := part1(xmas)
OUTER:
	for i := 0; i < len(xmas); i++ {
		for j := i + 2; j <= len(xmas); j++ {
			if m := sliceutils.Sum(xmas[i:j]); m > n {
				continue OUTER
			} else if m == n {
				return sliceutils.Min(xmas[i:j]) + sliceutils.Max(xmas[i:j])
			}
		}
	}
	log.Fatalln("Not found")
	return 0
}
