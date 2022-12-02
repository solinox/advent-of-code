package main

import (
	_ "embed"
	"strings"

	"github.com/solinox/advent-of-code/2022/util"
)

//go:embed input.txt
var input string

func main() {
	in := util.ParseSections(strings.NewReader(input), util.ParseInt)
	util.RunTimed(part1, in)
	util.RunTimed(part2, in)
}

func calorieSums(elves [][]int) []int {
	return util.SliceFrom(elves, func(v []int) int { return util.Sum(v...) })
}

func part1(elves [][]int) int {
	sums := calorieSums(elves)
	return util.Max(sums...)
}

func part2(elves [][]int) int {
	sums := calorieSums(elves)
	// faster than sorting for this input
	return util.Sum(util.MaxN(3, sums...)...)
}
