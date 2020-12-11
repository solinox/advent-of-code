package main

import (
	"log"
	"os"
	"sort"
	"time"

	"github.com/solinox/advent-of-code/2020/pkg/input"
)

func main() {
	r1, r2 := input.Duplicate(os.Stdin)
	joltAdaptersMap := input.IntMap(r1)
	joltAdaptersSlice := input.IntSlice(r2)

	t0 := time.Now()
	log.Println(part1(joltAdaptersMap), time.Since(t0))

	t0 = time.Now()
	log.Println(part2(joltAdaptersSlice), time.Since(t0))
}

// trying something other than the obvious sort method
func part1(joltAdapters map[int]int) int {
	numAdapters := len(joltAdapters)
	usedAdapters := 0
	lastAdapter := 0
	num1diff, num3diff := 0, 0
	for i := 0; usedAdapters < numAdapters; i++ {
		if _, ok := joltAdapters[i]; !ok {
			continue
		}
		usedAdapters++
		diff := i - lastAdapter
		if diff == 1 {
			num1diff++
		} else if diff == 3 {
			num3diff++
		}
		lastAdapter = i
	}
	num3diff++ // for the device itself
	return num1diff * num3diff
}

func part2(joltAdapters []int) uint {
	sort.Ints(joltAdapters)
	joltAdapters = append([]int{0}, append(joltAdapters, joltAdapters[len(joltAdapters)-1]+3)...) // outlet and device
	segments := make([][]int, 0)                                                                  // inner slices are a series of 1diff adapters, outer slice split by 3diff adapters
	for i := 0; i < len(joltAdapters); i++ {
		lastAdapter := joltAdapters[i]
		for j := i + 1; j < len(joltAdapters); j++ {
			if diff := joltAdapters[j] - lastAdapter; diff == 1 {
				lastAdapter = joltAdapters[j]
				continue
			}
			segments = append(segments, joltAdapters[i:j])
			i = j - 1
			break
		}
	}
	permutations := uint(1)
	// first and last item in segment must be used
	// key represents length of segment, value represents permutations of segment
	// my input only has segments up to 5 in length
	permMap := map[int]uint{
		1: 1,
		2: 1,
		3: 2,
		4: 4,
		5: 7,
	}
	for _, innerSegment := range segments {
		permutations *= permMap[len(innerSegment)]
	}
	return permutations
}
