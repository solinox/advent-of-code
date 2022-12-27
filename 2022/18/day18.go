package main

import (
	_ "embed"
	"fmt"
	"math"
	"strings"

	"github.com/solinox/advent-of-code/2022/util"
)

//go:embed input.txt
var input string

type Vector = util.Vector3
type Range = util.Range

var dirs = []Vector{{-1, 0, 0}, {1, 0, 0}, {0, -1, 0}, {0, 1, 0}, {0, 0, -1}, {0, 0, 1}}

func main() {
	cubes := util.ParseLinesReduce(strings.NewReader(input), func(cubes map[Vector]bool, s string) map[Vector]bool {
		var v Vector
		fmt.Sscanf(s, "%d,%d,%d", &v.X, &v.Y, &v.Z)
		cubes[v] = true
		return cubes
	}, make(map[Vector]bool))
	util.RunTimed(part1, cubes)
	util.RunTimed(part2Optimized, cubes)
	util.RunTimed(part2, cubes)
}

func part1(cubes map[Vector]bool) int {
	surfaceArea := 0
	for c := range cubes {
		for _, d := range dirs {
			if !cubes[c.Add(d)] {
				surfaceArea++
			}
		}
	}
	return surfaceArea
}

// get bounding cube to cover all cubes, then flood fill from one corner
// and count all the surfaces that flooding comes in contact with
func part2Optimized(cubes map[Vector]bool) int {
	surfaceArea := 0
	minX, maxX, minY, maxY, minZ, maxZ := math.MaxInt, math.MinInt, math.MaxInt, math.MinInt, math.MaxInt, math.MinInt
	for c := range cubes {
		if c.X <= minX {
			minX = c.X - 1
		}
		if c.X >= maxX {
			maxX = c.X + 1
		}
		if c.Y <= minY {
			minY = c.Y - 1
		}
		if c.Y >= maxY {
			maxY = c.Y + 1
		}
		if c.Z <= minZ {
			minZ = c.Z - 1
		}
		if c.Z >= maxZ {
			maxZ = c.Z + 1
		}
	}
	start := Vector{Z: minZ, Y: minY, X: minX}
	checked := map[Vector]bool{start: true}
	cells := []Vector{start}
	for len(cells) > 0 {
		c := cells[0]
		cells = cells[1:]
		for _, d := range dirs {
			cd := c.Add(d)
			if cubes[cd] {
				surfaceArea++
			} else if cd.Z >= minZ && cd.Z <= maxZ && cd.Y >= minY && cd.Y <= maxY && cd.X >= minX && cd.X <= maxX && !checked[cd] {
				cells = append(cells, cd)
				checked[cd] = true
			}
		}
	}
	return surfaceArea
}

// for each cube, check if cube face is an exterior face
// by getting min/max values for each direction (xy/yz/xz)
// unfortunately there are some edge cases which cause us to have to recursively search neighbors
// to confirm if an open face is exterior or not
func part2(cubes map[Vector]bool) int {
	surfaceArea := 0
	ranges := make(map[Vector]Range)
	exteriors := make(map[Vector]bool)
	interiors := make(map[Vector]bool)
	for c := range cubes {
		xy, yz, xz := Vector{X: c.X, Y: c.Y}, Vector{Y: c.Y, Z: c.Z}, Vector{X: c.X, Z: c.Z}
		rxy, ok := ranges[xy]
		if !ok || c.Z < rxy.Min {
			rxy.Min = c.Z
		}
		if !ok || c.Z > rxy.Max {
			rxy.Max = c.Z
		}
		ryz, ok := ranges[yz]
		if !ok || c.X < ryz.Min {
			ryz.Min = c.X
		}
		if !ok || c.X > ryz.Max {
			ryz.Max = c.X
		}
		rxz, ok := ranges[xz]
		if !ok || c.Y < rxz.Min {
			rxz.Min = c.Y
		}
		if !ok || c.Y > rxz.Max {
			rxz.Max = c.Y
		}
		ranges[xy], ranges[yz], ranges[xz] = rxy, ryz, rxz
	}
	for c := range cubes {
		for _, d := range dirs {
			if vd := c.Add(d); isExterior(cubes, exteriors, interiors, make(map[Vector]bool), ranges, vd) {
				exteriors[vd] = true
				surfaceArea++
			} else {
				interiors[vd] = true
			}
		}
	}
	return surfaceArea
}

func isExterior(cubes, exteriors, interiors, checked map[Vector]bool, ranges map[Vector]Range, v Vector) bool {
	if cubes[v] || checked[v] || interiors[v] {
		return false
	}
	if exteriors[v] {
		return true
	}
	xy, yz, xz := Vector{X: v.X, Y: v.Y}, Vector{Y: v.Y, Z: v.Z}, Vector{X: v.X, Z: v.Z}
	if r, ok := ranges[xy]; !ok || v.Z < r.Min || v.Z > r.Max {
		return true
	}
	if r, ok := ranges[yz]; !ok || v.X < r.Min || v.X > r.Max {
		return true
	}
	if r, ok := ranges[xz]; !ok || v.Y < r.Min || v.Y > r.Max {
		return true
	}
	checked[v] = true
	for _, d := range dirs {
		if vd := v.Add(d); isExterior(cubes, exteriors, interiors, checked, ranges, vd) {
			exteriors[vd] = true
			return true
		}
	}
	return false
}
