package main

import (
	_ "embed"

	"github.com/solinox/advent-of-code/2022/util"
)

//go:embed input.txt
var input []byte

func main() {
	util.RunTimed(part1, input)
	util.RunTimed(part2, input)
}

func isValid(marker []byte) bool {
	for i, v := range marker {
		for _, vv := range marker[i+1:] {
			if v == vv {
				return false
			}
		}
	}
	return true
}

func findPacket(in []byte, n int) int {
	for i, j := 0, n; j <= len(in); i, j = i+1, j+1 {
		marker := in[i:j]
		if isValid(marker) {
			return j
		}
	}
	return -1
}

func part1(in []byte) int {
	return findPacket(in, 4)
}

func part2(in []byte) int {
	return findPacket(in, 14)
}
