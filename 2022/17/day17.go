package main

import (
	_ "embed"

	"github.com/solinox/advent-of-code/2022/util"
)

//go:embed input.txt
var input []byte

type Vector = util.Vector

type rock struct {
	w, h int
	v    []Vector
}

var rocks = []rock{
	/*
		####
	*/
	{w: 4, h: 1, v: []Vector{{0, 0}, {0, 1}, {0, 2}, {0, 3}}},
	/*
		.#.
		###
		.#.
	*/
	{w: 3, h: 3, v: []Vector{{0, 1}, {1, 0}, {1, 1}, {1, 2}, {2, 1}}},
	/*
		..#
		..#
		###
	*/
	{w: 3, h: 3, v: []Vector{{0, 0}, {0, 1}, {0, 2}, {1, 2}, {2, 2}}},
	/*
		#
		#
		#
		#
	*/
	{w: 1, h: 4, v: []Vector{{0, 0}, {1, 0}, {2, 0}, {3, 0}}},
	/*
		##
		##
	*/
	{w: 2, h: 2, v: []Vector{{0, 0}, {0, 1}, {1, 0}, {1, 1}}},
}

func main() {
	util.RunTimed(part1, input)
	// util.RunTimed(part2, input)
}

func part1(jets []byte) int {
	dropped, h, j := make(map[Vector]bool), 0, 0
	for i := 0; i < 2022; i++ {
		dropped, h, j = drop(rocks[i%5], dropped, jets, h, j)
	}
	return h
}

func part2(jets []byte) int {
	dropped, h, j := make(map[Vector]bool), 0, 0
	// this doesn't work (takes too long, OOMKilled)
	// but there must be a pattern to the rocks falling that lets us increase the height dh for every n rocks
	for i := 0; i < 1_000_000_000_000; i++ {
		dropped, h, j = drop(rocks[i%5], dropped, jets, h, j)
	}
	return h
}

func okMove(rock rock, dropped map[Vector]bool, x, y int) bool {
	if x < 0 || x+rock.w > 7 || y < 0 {
		return false
	}
	for _, v := range rock.v {
		v.X += x
		v.Y += y
		if dropped[v] {
			return false
		}
	}
	return true
}

func drop(rock rock, dropped map[Vector]bool, jets []byte, h, j int) (map[Vector]bool, int, int) {
	x, y := 2, h+3
	for {
		dir := jets[j]
		j = (j + 1) % len(jets)
		dx := 1
		if dir == '<' {
			dx = -1
		}
		if nx := x + dx; okMove(rock, dropped, nx, y) {
			x = nx
		}
		dy := -1
		if ny := y + dy; okMove(rock, dropped, x, ny) {
			y = ny
		} else {
			// at rest
			break
		}
	}
	for _, v := range rock.v {
		v.X += x
		v.Y += y
		dropped[v] = true
	}
	if nh := y + rock.h; nh > h {
		h = nh
	}
	return dropped, h, j
}
