package main

import (
	_ "embed"
	"strings"

	"github.com/solinox/advent-of-code/2022/util"
)

//go:embed input.txt
var input string

func main() {
	grid := util.ParseLines(strings.NewReader(input), func(l string) []byte { return []byte(l) })
	util.RunTimed(part1, grid)
	util.RunTimed(part2, grid)
}

func visibleDirections(grid [][]byte, i, j int, includeEqual bool) (up, down, left, right int) {
	v := grid[i][j]
	for n := i - 1; n >= 0; n-- {
		up++
		if grid[n][j] >= v {
			if !includeEqual {
				up--
			}
			break
		}
	}
	for n := i + 1; n < len(grid); n++ {
		down++
		if grid[n][j] >= v {
			if !includeEqual {
				down--
			}
			break
		}
	}
	for n := j - 1; n >= 0; n-- {
		left++
		if grid[i][n] >= v {
			if !includeEqual {
				left--
			}
			break
		}
	}
	for n := j + 1; n < len(grid[i]); n++ {
		right++
		if grid[i][n] >= v {
			if !includeEqual {
				right--
			}
			break
		}
	}
	return
}

func isVisible(grid [][]byte, i, j int) bool {
	up, down, left, right := visibleDirections(grid, i, j, false)
	return up == i || down == len(grid)-i-1 || left == j || right == len(grid[i])-j-1
}

func scenicScore(grid [][]byte, i, j int) int {
	up, down, left, right := visibleDirections(grid, i, j, true)
	return up * down * left * right
}

func part1(grid [][]byte) int {
	visible := len(grid)*2 + len(grid[0])*2 - 4 // edges
	for i := 1; i < len(grid)-1; i++ {
		for j := 1; j < len(grid[i])-1; j++ {
			if isVisible(grid, i, j) {
				visible++
			}
		}
	}
	return visible
}

func part2(grid [][]byte) int {
	max := 0
	for i := 1; i < len(grid)-1; i++ {
		for j := 1; j < len(grid[i])-1; j++ {
			if v := scenicScore(grid, i, j); v > max {
				max = v
			}
		}
	}
	return max
}
