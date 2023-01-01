package main

import (
	_ "embed"
	"strings"

	"github.com/solinox/advent-of-code/2022/util"
	"golang.org/x/exp/slices"
)

//go:embed input.txt
var input string

func main() {
	in := util.ParseLines(strings.NewReader(input), util.ParseInt)
	util.RunTimed(part1, in)
	util.RunTimed(part2, in)
}

func part1(in []int) int {
	out := decrypt(in, 1)
	return groveCoords(out)
}

func part2(in []int) int {
	keyedIn := make([]int, len(in))
	for i := range in {
		keyedIn[i] = in[i] * 811589153
	}
	out := decrypt(keyedIn, 10)
	return groveCoords(out)
}

/*
O(n^2) solution by storing a list of the indices of each value
but for n=5000 it still runs in <1s thankfully.

Initial arrangement:         Indices (ix):
1, 2, -3, 3, -2, 0, 4          [0 1 2 3 4 5 6]
1 moves between 2 and -3:    0th moves 1 (and vals 0>1 move -1)
2, 1, -3, 3, -2, 0, 4          [1 0 2 3 4 5 6]
2 moves between -3 and 3:    1st moves 2 (and vals 0>2 move -1)
1, -3, 2, 3, -2, 0, 4          [0 2 1 3 4 5 6]
-3 moves between -2 and 0:   2nd moves -3 (and vals 1>4 move -1)
1, 2, 3, -2, -3, 0, 4          [0 1 4 2 3 5 6]
3 moves between 0 and 4:     3rd moves 3 (and vals 2>5 move -1)
1, 2, -2, -3, 0, 3, 4          [0 1 3 5 2 4 6]
-2 moves between 4 and 1:    4th moves -2 (and vals 2<0=2>6 move -1)
1, 2, -3, 0, 3, 4, -2          [0 1 2 4 6 3 5]
0 does not move:             5th moves 0
1, 2, -3, 0, 3, 4, -2          [0 1 2 4 6 3 5]
4 moves between -3 and 0:    6th moves 4 (and vals 5<3 move +1)
1, 2, -3, 4, 0, 3, -2          [0 1 2 5 6 4 3]
*/
func decrypt(in []int, n int) []int {
	ix := make([]int, len(in))
	k := len(ix) - 1
	for i := range ix {
		ix[i] = i
	}
	for m := 0; m < n; m++ {
		for i, v := range in {
			if v == 0 {
				continue
			}
			v = v % k
			a, b := ix[i], ix[i]+v
			if b < 0 {
				b += k
			} else if b >= k {
				b -= k
			}
			// to simulate an element being inserted, we need to shift the index of other values
			// between the a and b index values
			if a < b {
				for x, y := range ix {
					if y >= a && y <= b {
						ix[x] = y - 1
					}
				}
			} else {
				for x, y := range ix {
					if y <= a && y >= b {
						ix[x] = y + 1
					}
				}
			}
			ix[i] = b
		}
	}
	return reconstruct(in, ix)
}

func groveCoords(out []int) int {
	i, k := slices.Index(out, 0), len(out)
	return out[(i+1000)%k] + out[(i+2000)%k] + out[(i+3000)%k]
}

func reconstruct(in, ix []int) []int {
	out := make([]int, len(in))
	for i, v := range in {
		out[ix[i]] = v
	}
	return out
}
