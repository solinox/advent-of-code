package main

import (
	_ "embed"
	"strings"

	"github.com/solinox/advent-of-code/2022/util"
)

//go:embed input.txt
var input string

type Vector struct {
	Y, X int
}

func (p Vector) Add(o Vector) Vector {
	return Vector{Y: p.Y + o.Y, X: p.X + o.X}
}

func (p Vector) Sub(o Vector) Vector {
	return Vector{Y: p.Y - o.Y, X: p.X - o.X}
}

func (p Vector) Unit() Vector {
	v := Vector{}
	if p.Y != 0 {
		v.Y = p.Y / util.Abs(p.Y)
	}
	if p.X != 0 {
		v.X = p.X / util.Abs(p.X)
	}
	return v
}

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

func part1(moves []Move) int {
	head, tail := Vector{0, 0}, Vector{0, 0}
	visited := map[Vector]bool{tail: true}
	for _, m := range moves {
		for n := 0; n < m.N; n++ {
			head = head.Add(m.Dir)
			// fmt.Println(head, tail, m)
			if v := head.Sub(tail); util.Abs(v.Y) > 1 || util.Abs(v.X) > 1 {
				tail = tail.Add(v.Unit())
				visited[tail] = true
				// fmt.Println("moved tail", tail, v, v.Unit())
			}
		}
	}
	return len(visited)
}

func part2(moves []Move) int {
	knots := make([]Vector, 10)
	visited := map[Vector]bool{Vector{}: true}
	for _, m := range moves {
		for n := 0; n < m.N; n++ {
			knots[0] = knots[0].Add(m.Dir)
			for i := 1; i < len(knots); i++ {
				if v := knots[i-1].Sub(knots[i]); util.Abs(v.Y) > 1 || util.Abs(v.X) > 1 {
					knots[i] = knots[i].Add(v.Unit())
					if i == len(knots)-1 {
						visited[knots[len(knots)-1]] = true
					}
				}
			}
		}
	}
	return len(visited)
}
