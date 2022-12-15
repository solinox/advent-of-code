package main

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/solinox/advent-of-code/2022/util"
)

//go:embed input.txt
var input string

type Vector = util.Vector
type Points map[Vector]bool

func main() {
	rocks := parse(input)
	util.RunTimed(part1, rocks)
	util.RunTimed(part2, rocks)
}

func parse(in string) Points {
	return util.ParseLinesReduce(
		strings.NewReader(input),
		func(agg Points, s string) Points {
			rockSections := strings.Split(s, " -> ")
			for i := 1; i < len(rockSections); i++ {
				var v1, v2 Vector
				fmt.Sscanf(rockSections[i-1], "%d,%d", &v1.X, &v1.Y)
				fmt.Sscanf(rockSections[i], "%d,%d", &v2.X, &v2.Y)
				agg[v1], agg[v2] = true, true
				dv := v2.Sub(v1).Unit()
				for v := v1.Add(dv); v != v2; v = v.Add(dv) {
					agg[v] = true
				}
			}
			return agg
		},
		make(Points),
	)
}

func part1(rocks Points) int {
	return produceSand(rocks, Vector{X: 500, Y: 0}, false)
}

func part2(rocks Points) int {
	return produceSand(rocks, Vector{X: 500, Y: 0}, true) + 1 // the sand at the source is not counted in produceSand
}

func produceSand(rocks Points, sandSource Vector, hasFloor bool) int {
	obstructions := make(Points)
	for k, v := range rocks {
		obstructions[k] = v
	}

	maxY := maxY(obstructions)
	sand := 0
	for s, ok := dropSand(obstructions, sandSource, maxY, hasFloor); ok && s != sandSource; s, ok = dropSand(obstructions, sandSource, maxY, hasFloor) {
		sand++
		obstructions[s] = true
	}
	return sand
}

func dropSand(obstructions Points, start Vector, maxY int, hasFloor bool) (Vector, bool) {
	v := Vector{Y: start.Y + 1, X: start.X}
	if !hasFloor && v.Y > maxY {
		return start, false
	} else if hasFloor {
		// add enough of the infinite floor to obstructions in order to block it
		vv := Vector{Y: maxY + 2, X: v.X}
		obstructions[vv] = true
		vv.X--
		obstructions[vv] = true
		vv.X += 2
		obstructions[vv] = true
	}
	if !obstructions[v] {
		return dropSand(obstructions, v, maxY, hasFloor)
	}
	v.X--
	if !obstructions[v] {
		return dropSand(obstructions, v, maxY, hasFloor)
	}
	v.X += 2
	if !obstructions[v] {
		return dropSand(obstructions, v, maxY, hasFloor)
	}
	return start, true
}

func maxY(p Points) int {
	maxY := 0
	for v := range p {
		if v.Y > maxY {
			maxY = v.Y
		}
	}
	return maxY
}
