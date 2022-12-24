package main

import (
	"bytes"
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
	util.RunTimed(part2, input)
}

const nRows = 70

var empty = bytes.Repeat([]byte{'.'}, 7*nRows)

func part1(jets []byte) int {
	// return dropRocksN(rocks, jets, 2022) also works, but is a bit slower
	// due to the time spent populating the cache, moreso than the benefit for a small n (2022)
	rs, h, j := make([]byte, 7*nRows), 0, 0
	copy(rs, empty)
	for i := 0; i < 2022; i++ {
		rs, h, j = drop(rocks[i%5], rs, jets, h, j)
	}
	return h
}

func part2(jets []byte) int {
	return dropRocksN(rocks, jets, 1_000_000_000_000)
}

func dropRocksN(rocks []rock, jets []byte, n int) int {
	rs, h, j := make([]byte, 7*nRows), 0, 0
	copy(rs, empty)
	type key struct {
		rs    string
		ro, j int
	}
	cache := make(map[key][2]int)
	for i := 0; i < n; i++ {
		k := key{string(rs), i % 5, j}
		if a, ok := cache[k]; ok {
			di, dh := i-a[0], h-a[1]
			for ; i+di < n; i += di {
				h += dh
			}
		} else {
			cache[k] = [2]int{i, h}
		}
		rs, h, j = drop(rocks[i%5], rs, jets, h, j)
	}
	return h
}

func okMove(rock rock, rocks []byte, x, y, h int) bool {
	if x < 0 || x+rock.w > 7 || y < 0 {
		return false
	}
	if y > h {
		return true
	}
	for _, v := range rock.v {
		if v.Y+y <= h && rocks[(h-y-v.Y)*7+x+v.X] == '#' {
			return false
		}
	}
	return true
}

func drop(rock rock, rs []byte, jets []byte, h, j int) ([]byte, int, int) {
	x, y := 2, h+3
	for {
		dir := jets[j]
		j = (j + 1) % len(jets)
		dx := 1
		if dir == '<' {
			dx = -1
		}
		if nx := x + dx; okMove(rock, rs, nx, y, h) {
			x = nx
		}
		dy := -1
		if ny := y + dy; okMove(rock, rs, x, ny, h) {
			y = ny
		} else {
			break
		}
	}
	if dh := y + rock.h - h; dh > 0 {
		copy(rs[dh*7:], rs)
		copy(rs[:dh*7], empty)
		h += dh
	}
	for _, v := range rock.v {
		rs[(h-y-v.Y)*7+x+v.X] = '#'
	}
	return rs, h, j
}
