package main

import (
	"io/ioutil"
	"strings"
	"fmt"
	"math/bits"
)

// I used complex numbers as a 2d grid because I wanted to try it out
type grid map[complex128]point

type point struct {
	Wires uint
	Steps int
}

func main() {
	wirePaths := parseInput("input.txt")
	grid := buildGrid(wirePaths)
	intersections := getIntersections(grid)
	if len(intersections) == 0 {
		fmt.Println("No intersections found")
		return
	}

	// Part 1	
	part1 := dist(intersections[0])
	for _, intersection := range intersections[1:] {
		if d := dist(intersection); d < part1 {
			part1 = d
		}
	}
	fmt.Println("Part 1", part1)

	// Part 2
	stepsPerIntersection := stepsPerIntersection(intersections, grid)
	part2 := stepsPerIntersection[0]
	for _, steps := range stepsPerIntersection[1:] {
		if steps < part2 {
			part2 = steps
		}
	}
	fmt.Println("Part 2", part2)

}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func dist(intersection complex128) int {
	return abs(int(real(intersection))) + abs(int(imag(intersection)))
}

func buildGrid(wirePaths [][]string) grid {
	grid := make(grid)
	for wire, path := range wirePaths {
		wirePosition := 0 + 0i
		steps := 0
		for _, instruction := range path {
			var direction byte
			var units int
			var step complex128
			fmt.Sscanf(instruction, "%c%d", &direction, &units)
			switch direction {
			case 'U':
				step = 0+1i
			case 'D':
				step = 0-1i
			case 'R':
				step = 1+0i
			case 'L':
				step = -1+0i
				
			}
			for n := 0; n < units; n++ {
				steps++
				wirePosition = wirePosition + step
				point := grid[wirePosition]
				if point.Wires & (1 << wire) == 0 {
					point.Steps += steps
				}
				point.Wires |= (1 << wire)
				grid[wirePosition] = point
			}
		}
	}
	return grid
}

func getIntersections(grid grid) []complex128 {
	intersections := make([]complex128, 0)
	for pos, wireCrossings := range grid {
		if bits.OnesCount(wireCrossings.Wires) > 1 {
			intersections = append(intersections, pos)
		}
	}
	return intersections
}

func stepsPerIntersection(intersections []complex128, grid grid) []int {
	steps := make([]int, len(intersections))
	for i := range intersections {
		steps[i] = grid[intersections[i]].Steps
	}
	return steps
}

func parseInput(filename string) [][]string {
	data, _ := ioutil.ReadFile(filename)
	wires := strings.Split(string(data), "\n")
	wirePaths := make([][]string, len(wires))
	for i := range wires {
		wirePaths[i] = strings.Split(wires[i], ",")
	}
	return wirePaths
}
