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

func isVisible(grid [][]byte, i, j int) bool {
	if i == 0 || i == len(grid)-1 || j == 0 || j == len(grid[i])-1 {
		// edge of forest
		return true
	}
	v := grid[i][j]
	// check row
	if util.Max(grid[i][:j]...) < v || util.Max(grid[i][j+1:]...) < v {
		return true
	}
	// check column
	fromTop, fromBottom := true, true
	for r := 0; r < i; r++ {
		if grid[r][j] >= v {
			fromTop = false
			break
		}
	}
	if fromTop {
		return true
	}
	for r := i + 1; r < len(grid); r++ {
		if grid[r][j] >= v {
			fromBottom = false
			break
		}
	}
	return fromBottom
}

func scenicScore(grid [][]byte, i, j int) int {
	v := grid[i][j]
	left, right, up, down := 0, 0, 0, 0
	for n := i - 1; n >= 0; n-- {
		up++
		if grid[n][j] >= v {
			break
		}
	}
	for n := i + 1; n < len(grid); n++ {
		down++
		if grid[n][j] >= v {
			break
		}
	}
	for n := j - 1; n >= 0; n-- {
		left++
		if grid[i][n] >= v {
			break
		}
	}
	for n := j + 1; n < len(grid[i]); n++ {
		right++
		if grid[i][n] >= v {
			break
		}
	}
	return left * right * up * down
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
