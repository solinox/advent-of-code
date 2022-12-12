package main

import (
	"bytes"
	_ "embed"
	"math"

	"github.com/solinox/advent-of-code/2022/util"
)

//go:embed input.txt
var input []byte

type Vector = util.Vector
type HeightMap struct {
	Grid  [][]byte
	Start Vector
	End   Vector
}

func main() {
	hm := parse(input)
	util.RunTimed(part1, hm)
	util.RunTimed(part2, hm)
}

func parse(input []byte) HeightMap {
	hm := HeightMap{Grid: make([][]byte, 0)}
	for y, l := range bytes.Split(input, []byte{'\n'}) {
		if x := bytes.Index(l, []byte{'S'}); x != -1 {
			hm.Start = Vector{Y: y, X: x}
			l[x] = 'a'
		}
		if x := bytes.Index(l, []byte{'E'}); x != -1 {
			hm.End = Vector{Y: y, X: x}
			l[x] = 'z'
		}
		hm.Grid = append(hm.Grid, l)
	}
	return hm
}

func astar(g [][]byte, start, end Vector, prevSolved map[Vector]int) int {
	openSet := map[Vector]struct{}{start: {}}
	cameFrom := map[Vector]Vector{}
	gScore := map[Vector]int{start: 0}
	fScore := map[Vector]int{start: end.Dist(start)}

	// stores min steps from each node along the path, and returns the steps from start to end
	reconstructPath := func(cameFrom map[Vector]Vector, current Vector, prevSolved map[Vector]int, v int) []Vector {
		path := []Vector{current}
		for n, ok := cameFrom[current]; ok; n, ok = cameFrom[n] {
			path = append(path, n)
			prevSolved[n] = len(path) - 1 + v
		}
		return path
	}

	min := func(openSet map[Vector]struct{}, fScore map[Vector]int) Vector {
		min := math.MaxInt
		minV := Vector{}
		for k := range openSet {
			if s := fScore[k]; s < min {
				min = s
				minV = k
			}
		}
		return minV
	}

	for len(openSet) > 0 {
		current := min(openSet, fScore)
		if current == end {
			return len(reconstructPath(cameFrom, current, prevSolved, 0)) - 1
		}
		if v, ok := prevSolved[current]; ok && v > 0 {
			// the minimal path from the current node has already been solved
			return len(reconstructPath(cameFrom, current, prevSolved, v)) - 1 + v
		}
		delete(openSet, current)
		for _, n := range []Vector{
			current.Add(Vector{1, 0}),
			current.Add(Vector{-1, 0}),
			current.Add(Vector{0, 1}),
			current.Add(Vector{0, -1}),
		} {
			if n.X < 0 || n.Y < 0 || n.Y >= len(g) || n.X >= len(g[n.Y]) {
				// out of bounds
				continue
			}
			if int(g[n.Y][n.X])-int(g[current.Y][current.X]) > 1 {
				// cannot climb that fast
				continue
			}
			tentativeG := gScore[current] + 1
			if g, ok := gScore[n]; !ok || tentativeG < g {
				gScore[n] = tentativeG
				fScore[n] = tentativeG + end.Dist(n)
				cameFrom[n] = current
				openSet[n] = struct{}{}
			}
		}
	}

	return -1
}

func part1(hm HeightMap) int {
	return astar(hm.Grid, hm.Start, hm.End, make(map[Vector]int))
}

func part2(hm HeightMap) int {
	min := math.MaxInt
	prevSolved := make(map[Vector]int)
	for y := 0; y < len(hm.Grid); y++ {
		for x := 0; x < len(hm.Grid[y]); x++ {
			if hm.Grid[y][x] == 'a' {
				start := Vector{Y: y, X: x}
				v := astar(hm.Grid, start, hm.End, prevSolved)
				if v > 0 && v < min {
					min = v
				}
			}
		}
	}
	return min
}
