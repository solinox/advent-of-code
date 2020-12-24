package main

import (
	"log"
	"os"
	"time"

	"github.com/solinox/advent-of-code/2020/pkg/input"
)

func main() {
	tiles := input.StringSlice(os.Stdin)

	t0 := time.Now()
	log.Println(part1(tiles), time.Since(t0))

	t0 = time.Now()
	log.Println(part2(tiles), time.Since(t0))
}

func part1(tiles []string) int {
	tileGrid := newTileGrid(tiles)

	return blackTiles(tileGrid)
}

func part2(tiles []string) int {
	tileGrid := newTileGrid(tiles)

	tileGrid = dayFlips(tileGrid, 100)

	return blackTiles(tileGrid)
}

func dayFlips(tileGrid map[vector]bool, days int) map[vector]bool {
	for i := 0; i < days; i++ {
		tilesToFlip := make([]vector, 0)

		tileGrid = trim(tileGrid)
		tileGrid = addNeighbors(tileGrid)

		for offset, black := range tileGrid {
			numBlackNeighbors := blackNeighbors(tileGrid, offset)
			if black && (numBlackNeighbors == 0 || numBlackNeighbors > 2) {
				tilesToFlip = append(tilesToFlip, offset)
			} else if !black && numBlackNeighbors == 2 {
				tilesToFlip = append(tilesToFlip, offset)
			}
		}

		for _, offset := range tilesToFlip {
			tileGrid[offset] = !tileGrid[offset]
		}
	}
	return tileGrid
}

func trim(tileGrid map[vector]bool) map[vector]bool {
	for offset := range tileGrid {
		if !tileGrid[offset] {
			delete(tileGrid, offset)
		}
	}
	return tileGrid
}

func addNeighbors(tileGrid map[vector]bool) map[vector]bool {
	for offset := range tileGrid {
		for _, neighbor := range neighbors(offset) {
			if _, ok := tileGrid[neighbor]; !ok {
				tileGrid[neighbor] = false
			}
		}
	}
	return tileGrid
}

func neighbors(offset vector) []vector {
	if offset.Y%2 == 0 {
		return []vector{
			{offset.X - 1, offset.Y},
			{offset.X + 1, offset.Y},
			{offset.X - 1, offset.Y - 1},
			{offset.X, offset.Y - 1},
			{offset.X - 1, offset.Y + 1},
			{offset.X, offset.Y + 1},
		}
	}
	return []vector{
		{offset.X - 1, offset.Y},
		{offset.X + 1, offset.Y},
		{offset.X, offset.Y - 1},
		{offset.X + 1, offset.Y - 1},
		{offset.X, offset.Y + 1},
		{offset.X + 1, offset.Y + 1},
	}
}

func blackNeighbors(tileGrid map[vector]bool, offset vector) int {
	numBlackNeighbors := 0
	for _, neighbor := range neighbors(offset) {
		if tileGrid[neighbor] {
			numBlackNeighbors++
		}
	}
	return numBlackNeighbors
}

func blackTiles(tileGrid map[vector]bool) int {
	numBlackTiles := 0
	for _, black := range tileGrid {
		if black {
			numBlackTiles++
		}
	}
	return numBlackTiles
}

type vector struct {
	X, Y int
}

func newTileGrid(tiles []string) map[vector]bool {
	grid := make(map[vector]bool)
	for _, tile := range tiles {
		offset := findOffset(tile)
		grid[offset] = !grid[offset]
	}
	return grid
}

// uses offset coordinates where each odd row is shifted right
// e.g. going se from the reference (0, 0) tile will become (-1, 0)
// but going se from that tile will become (-2, 1).
func findOffset(tile string) vector {
	offset := vector{0, 0}
	halfStep := false
	for i := range tile {
		switch tile[i] {
		case 'n':
			offset.Y++
			halfStep = true
		case 's':
			offset.Y--
			halfStep = true
		case 'e':
			if halfStep {
				if offset.Y%2 == 0 {
					offset.X++
				}
			} else {
				offset.X++
			}
			halfStep = false
		case 'w':
			if halfStep {
				if offset.Y%2 != 0 {
					offset.X--
				}
			} else {
				offset.X--
			}
			halfStep = false
		}
	}
	return offset
}
