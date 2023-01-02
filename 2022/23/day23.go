package main

import (
	_ "embed"
	"math"
	"math/bits"
	"strings"

	"github.com/solinox/advent-of-code/2022/util"
)

//go:embed input.txt
var input string

type Vector = util.Vector
type Range = util.Range

type ground struct {
	elves map[Vector]bool
}

func main() {
	g := parse(input)
	util.RunTimed(part1, g)
	util.RunTimed(part2, g)
}

func parse(input string) ground {
	lines := strings.Split(input, "\n")
	g := ground{
		elves: make(map[Vector]bool),
	}
	for y := range lines {
		for x := range lines[y] {
			if lines[y][x] == '#' {
				g.elves[Vector{Y: y, X: x}] = true
			}
		}
	}
	return g
}

func part1(g ground) int {
	for i, ok := 0, true; i < 10 && ok; i++ {
		g, ok = gameOfLife(g, i)
	}
	return grounded(g)
}

func part2(g ground) int {
	i := 0
	for ok := true; ok; i++ {
		g, ok = gameOfLife(g, i)
	}
	return i
}

var (
	nswe = [4][3]Vector{
		{{-1, -1}, {-1, 0}, {-1, 1}}, // nw, n, ne
		{{1, -1}, {1, 0}, {1, 1}},    // sw, s, se
		{{-1, -1}, {0, -1}, {1, -1}}, // nw, w, sw
		{{-1, 1}, {0, 1}, {1, 1}},    // ne, e, se
	}
)

func gameOfLife(g ground, round int) (ground, bool) {
	dirs := [4][3]Vector{nswe[round%4], nswe[(round+1)%4], nswe[(round+2)%4], nswe[(round+3)%4]}
	proposed := make(map[Vector][]Vector)
	remain := make([]Vector, 0)
ELVES:
	for elf := range g.elves {
		blocked := uint(0)
		allowed := uint(0)
	DIRS:
		for i, dir := range dirs {
			for _, d := range dir {
				if g.elves[elf.Add(d)] {
					blocked |= (1 << i)
					if allowed != 0 {
						break DIRS
					}
					continue DIRS
				}
			}
			allowed |= (1 << i)
			if blocked != 0 {
				break DIRS
			}
		}
		if blocked == 0b1111 || allowed == 0b1111 {
			remain = append(remain, elf)
			continue ELVES
		}
		firstAllowed := bits.TrailingZeros(allowed)
		move := elf.Add(dirs[firstAllowed][1])
		proposed[move] = append(proposed[move], elf)
		continue ELVES
	}
	newElves := make(map[Vector]bool)
	for _, elf := range remain {
		newElves[elf] = true
	}
	for k, v := range proposed {
		if len(v) > 1 {
			for _, elf := range v {
				newElves[elf] = true
			}
			continue
		}
		newElves[k] = true
	}
	g.elves = newElves
	return g, len(remain) != len(newElves)
}

func grounded(g ground) int {
	count := 0
	ry, rx := g.bounds()
	for y := ry.Min; y <= ry.Max; y++ {
		for x := rx.Min; x <= rx.Max; x++ {
			if !g.elves[Vector{y, x}] {
				count++
			}
		}
	}
	return count
}

func (g ground) bounds() (Range, Range) {
	ry, rx := Range{Min: math.MaxInt, Max: math.MinInt}, Range{Min: math.MaxInt, Max: math.MinInt}
	for elf := range g.elves {
		if elf.Y < ry.Min {
			ry.Min = elf.Y
		}
		if elf.Y > ry.Max {
			ry.Max = elf.Y
		}
		if elf.X < rx.Min {
			rx.Min = elf.X
		}
		if elf.X > rx.Max {
			rx.Max = elf.X
		}
	}
	return ry, rx
}

func (g ground) String() string {
	var s strings.Builder
	ry, rx := g.bounds()
	for y := ry.Min; y <= ry.Max; y++ {
		for x := rx.Min; x <= rx.Max; x++ {
			if g.elves[Vector{y, x}] {
				s.WriteByte('#')
			} else {
				s.WriteByte('.')
			}
		}
		s.WriteByte('\n')
	}
	return s.String()
}
