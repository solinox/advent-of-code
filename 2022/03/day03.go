package main

import (
	_ "embed"
	"strings"

	"github.com/solinox/advent-of-code/2022/util"
)

//go:embed input.txt
var input string

func main() {
	in := util.ParseLines(strings.NewReader(input), func(s string) []byte { return []byte(s) })
	util.RunTimed(part1, in)
	util.RunTimed(part2, in)
}

func getItem(rucksack []byte) byte {
	for _, item1 := range rucksack[:len(rucksack)/2] {
		for _, item2 := range rucksack[len(rucksack)/2:] {
			if item1 == item2 {
				return item1
			}
		}
	}
	return 0
}

func priority(item byte) int {
	if item < 'a' {
		item += 58
	}
	return int(item + 1 - 'a')
}

func part1(rucksacks [][]byte) int {
	return util.Sum(util.SliceFrom(rucksacks, func(rs []byte) int { return priority(getItem(rs)) })...)
}

func getBadge(rucksacks [][]byte) byte {
	for _, item1 := range rucksacks[0] {
		for _, item2 := range rucksacks[1] {
			if item1 == item2 {
				for _, item3 := range rucksacks[2] {
					if item2 == item3 {
						return item1
					}
				}
			}
		}
	}
	return 0
}

func part2(rucksacks [][]byte) int {
	sum := 0
	for i := 0; i < len(rucksacks); i += 3 {
		sum += priority(getBadge(rucksacks[i : i+3]))
	}
	return sum
}
