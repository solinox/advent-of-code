package main

import (
	"log"
	"os"

	"github.com/solinox/advent-of-code/2020/pkg/input"
)

type grid [][]byte
type point struct {
	X, Y int
}

func (g grid) move(position, slope point) point {
	gridX := len(g[0])
	return point{(position.X + slope.X) % gridX, position.Y + slope.Y}
}

func main() {
	grid := grid(input.BytesSlice(os.Stdin))

	log.Println(part1(grid, point{0, 0}, point{3, 1}))

	log.Println(part2(grid, point{0, 0}, point{1, 1}, point{3, 1}, point{5, 1}, point{7, 1}, point{1, 2}))
}

func part1(grid grid, start, slope point) int {
	trees := 0
	gridY := len(grid)
	for position := start; position.Y < gridY; position = grid.move(position, slope) {
		if grid[position.Y][position.X] == '#' {
			trees++
		}
	}
	return trees
}

func part2(grid grid, start point, slopes ...point) int {
	n := 1
	for _, slope := range slopes {
		n *= part1(grid, start, slope)
	}
	return n
}
