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

// simplified some pieces to optimize the solution so this isn't a true A* algo anymore, sort of a bfs+A* hybrid
func astar(g [][]byte, start Vector, check func(v Vector) bool) int {
	openSet := map[Vector]int{start: 0}
	visited := map[Vector]bool{start: true}

	min := func(openSet map[Vector]int) (Vector, int) {
		min := math.MaxInt
		minV := Vector{}
		for k, v := range openSet {
			if v < min {
				min = v
				minV = k
			}
		}
		return minV, min
	}

	for len(openSet) > 0 {
		current, cost := min(openSet)
		if check(current) {
			return cost
		}
		visited[current] = true
		delete(openSet, current)
		for _, n := range []Vector{
			{Y: current.Y + 1, X: current.X},
			{Y: current.Y - 1, X: current.X},
			{Y: current.Y, X: current.X + 1},
			{Y: current.Y, X: current.X - 1},
		} {
			if n.X < 0 || n.Y < 0 || n.Y >= len(g) || n.X >= len(g[n.Y]) {
				// out of bounds
				continue
			}
			if visited[n] {
				// already visited
				continue
			}
			// going downhill instead of uphill
			if int(g[current.Y][current.X])-int(g[n.Y][n.X]) > 1 {
				// cannot climb that fast
				continue
			}
			cost := cost + 1
			openSet[n] = cost
		}
	}

	return -1
}

func part1(hm HeightMap) int {
	return astar(hm.Grid, hm.End, func(v Vector) bool { return v == hm.Start })
}

func part2(hm HeightMap) int {
	return astar(hm.Grid, hm.End, func(v Vector) bool { return hm.Grid[v.Y][v.X] == 'a' })
}
