package main

import (
	_ "embed"
	"math"
	"strconv"
	"strings"

	"github.com/solinox/advent-of-code/2022/util"
)

//go:embed input.txt
var input string

func main() {
	snafus := strings.Split(input, "\n")
	util.RunTimed(part1, snafus)
}

func part1(snafus []string) string {
	sum := util.Sum(util.SliceFrom(snafus, fromSnafu)...)
	return toSnafu(sum)
}

var snafuToDec = map[byte]int{'=': -2, '-': -1, '0': 0, '1': 1, '2': 2}

func fromSnafu(s string) int {
	n := 0
	for i, j := len(s)-1, float64(0); i >= 0; i, j = i-1, j+1 {
		n += int(math.Pow(5, j)) * snafuToDec[s[i]]
	}
	return n
}

var decToSnafu = map[int]byte{-2: '=', -1: '-', 0: '0', 1: '1', 2: '2'}

func toSnafu(n int) string {
	s := []byte(strconv.FormatInt(int64(n), 5))
	for i := len(s) - 1; i >= 0; i-- {
		v := int(s[i] - '0')
		// carry over
		if v >= 3 && v <= 5 {
			if i == 0 {
				s = append([]byte{'0'}, s...)
				i++
			}
			s[i-1]++
			v -= 5
		}
		s[i] = decToSnafu[v]
	}
	return string(s)
}
