package main

import (
	"bytes"
	"log"
	"os"
	"strings"
	"time"

	"github.com/solinox/advent-of-code/2020/pkg/input"
)

const (
	floor    byte = '.'
	empty    byte = 'L'
	occupied byte = '#'
)

type grid [][]byte

func main() {
	grid := grid(input.BytesSlice(os.Stdin))

	t0 := time.Now()
	log.Println(part1(grid), time.Since(t0))

	t0 = time.Now()
	log.Println(part2(grid), time.Since(t0))
}

func part1(g grid) int {
	g = findEqulibrium(g, 1, 4)
	return strings.Count(g.String(), string(occupied))
}

func part2(g grid) int {
	g = findEqulibrium(g, -1, 5)
	return strings.Count(g.String(), string(occupied))
}

func findEqulibrium(g grid, maxDist, occupiedThreshold int) grid {
	old, new := g, g.Step(maxDist, occupiedThreshold)
	for ; old.String() != new.String(); old, new = new, new.Step(maxDist, occupiedThreshold) {
	}
	return old
}

func (g grid) String() string {
	return string(bytes.Join(g, []byte{'\n'}))
}

func (g grid) Step(maxDist, occupiedThreshold int) grid {
	next := make(grid, len(g))
	for y := range g {
		next[y] = make([]byte, len(g[y]))
		for x := range g[y] {
			next[y][x] = g[y][x]

			switch g[y][x] {
			case floor:
				continue
			case empty:
				if n := g.OccupiedAdjacent(y, x, maxDist); n == 0 {
					next[y][x] = occupied
				}
			case occupied:
				if n := g.OccupiedAdjacent(y, x, maxDist); n >= occupiedThreshold {
					next[y][x] = empty
				}
			}
		}
	}
	return next
}

type vector struct {
	Y, X int
}

var directions = []vector{
	{-1, -1},
	{-1, 0},
	{-1, 1},
	{0, -1},
	{0, 1},
	{1, -1},
	{1, 0},
	{1, 1},
}

func (g grid) OccupiedAdjacent(y, x, maxDist int) int {
	numOccupied := 0
	for _, dir := range directions {
		for i, yy, xx := 0, y+dir.Y, x+dir.X; i < maxDist || maxDist < 0; i, yy, xx = i+1, yy+dir.Y, xx+dir.X {
			if yy < 0 || yy >= len(g) || xx < 0 || xx >= len(g[yy]) {
				break
			}
			if g[yy][xx] == floor {
				continue
			}
			if g[yy][xx] == occupied {
				numOccupied++
			}
			break
		}
	}
	return numOccupied
}
