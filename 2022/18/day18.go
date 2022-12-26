package main

import (
	_ "embed"
	"fmt"
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

func part2(cubes map[Vector]bool) int {
	surfaceArea := 0
	ranges := make(map[Vector]Range)
	exteriors := make(map[Vector]bool)
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
			if vd := c.Add(d); isExterior(cubes, exteriors, make(map[Vector]bool), ranges, vd) {
				exteriors[vd] = true
				surfaceArea++
			}
		}
	}
	return surfaceArea
}

func isExterior(cubes, exteriors, checked map[Vector]bool, ranges map[Vector]Range, v Vector) bool {
	if cubes[v] || checked[v] {
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
		if vd := v.Add(d); isExterior(cubes, exteriors, checked, ranges, vd) {
			exteriors[vd] = true
			return true
		}
	}
	return false
}
