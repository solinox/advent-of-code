package main

import (
	_ "embed"
	"strings"

	"github.com/solinox/advent-of-code/2022/util"
)

//go:embed input.txt
var input string

type Vector = util.Vector

type Move struct {
	Dir Vector
	N   int
}

func main() {
	moves := util.ParseLines(strings.NewReader(input), func(l string) Move {
		n := util.ParseInt(l[2:])
		d := Vector{}
		switch l[0] {
		case 'R':
			d.X++
		case 'L':
			d.X--
		case 'U':
			d.Y++
		case 'D':
			d.Y--
		}
		return Move{Dir: d, N: n}
	})
	util.RunTimed(part1, moves)
	util.RunTimed(part2, moves)
}

func trackTail(knots []Vector, moves []Move) int {
	visited := map[Vector]bool{knots[len(knots)-1]: true}
	for _, m := range moves {
		for n := 0; n < m.N; n++ {
			// move the head
			knots[0] = knots[0].Add(m.Dir)
			// check if following knots need to be moved
			for i := 1; i < len(knots); i++ {
				if v := knots[i-1].Sub(knots[i]); util.Abs(v.Y) > 1 || util.Abs(v.X) > 1 {
					knots[i] = knots[i].Add(v.Unit())
					// track the tail knot
					if i == len(knots)-1 {
						visited[knots[len(knots)-1]] = true
					}
				}
			}
		}
	}
	return len(visited)
}

func part1(moves []Move) int {
	return trackTail(make([]Vector, 2), moves)
}

func part2(moves []Move) int {
	return trackTail(make([]Vector, 10), moves)
}
