package main

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/solinox/advent-of-code/2022/util"
)

//go:embed input.txt
var input string

func main() {
	in := util.ParseLines(strings.NewReader(input), func(s string) [4]int {
		var n [4]int
		fmt.Sscanf(s, "%d-%d,%d-%d", &n[0], &n[1], &n[2], &n[3])
		return n
	})
	util.RunTimed(part1, in)
	util.RunTimed(part2, in)
}

func part1(pairs [][4]int) int {
	overlaps := 0
	for _, pair := range pairs {
		if (pair[0] <= pair[2] && pair[1] >= pair[3]) || (pair[0] >= pair[2] && pair[1] <= pair[3]) {
			overlaps++
		}
	}
	return overlaps
}

func part2(pairs [][4]int) int {
	overlaps := 0
	for _, pair := range pairs {
		if pair[1] < pair[2] || pair[0] > pair[3] {
			continue
		}
		overlaps++
	}
	return overlaps
}
