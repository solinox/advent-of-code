package main

import (
	_ "embed"
	"strings"

	"github.com/solinox/advent-of-code/2022/util"
)

//go:embed input.txt
var input string

func main() {
	in := util.ParseLines(strings.NewReader(input), func(s string) string { return s })
	util.RunTimed(part1, in)
	util.RunTimed(part2, in)
}

func part1(games []string) int {
	scores := map[string]int{
		"A X": 3 + 1,
		"A Y": 6 + 2,
		"A Z": 0 + 3,
		"B X": 0 + 1,
		"B Y": 3 + 2,
		"B Z": 6 + 3,
		"C X": 6 + 1,
		"C Y": 0 + 2,
		"C Z": 3 + 3,
	}
	results := util.SliceFrom(games, func(game string) int { return scores[game] })
	return util.Sum(results...)
}

func part2(games []string) int {
	scores := map[string]int{
		"A X": 0 + 3,
		"A Y": 3 + 1,
		"A Z": 6 + 2,
		"B X": 0 + 1,
		"B Y": 3 + 2,
		"B Z": 6 + 3,
		"C X": 0 + 2,
		"C Y": 3 + 3,
		"C Z": 6 + 1,
	}
	results := util.SliceFrom(games, func(game string) int { return scores[game] })
	return util.Sum(results...)
}
