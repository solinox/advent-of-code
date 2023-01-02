package main

import (
	_ "embed"
	"fmt"
	"math/bits"
	"strings"

	"github.com/solinox/advent-of-code/2022/util"
)

//go:embed input.txt
var input []byte

type Vector = util.Vector

type state uint

const (
	// use bits in order to preserve state of overlapping blizzard storms
	clear state = 0
	wall  state = 0b1
	up    state = 0b10
	down  state = 0b100
	left  state = 0b1000
	right state = 0b10000
)

type blizzards struct {
	w, h   int
	valley []state
}

func main() {
	bz := parse(input)
	util.RunTimed(part1, bz)
	util.RunTimed(part2, bz)
}

func parse(input []byte) blizzards {
	bz := blizzards{h: 1, valley: make([]state, 0)}
	for i := range input {
		switch input[i] {
		case '\n':
			bz.h++
		case '.':
			bz.valley = append(bz.valley, clear)
		case '#':
			bz.valley = append(bz.valley, wall)
		case '^':
			bz.valley = append(bz.valley, up)
		case 'v':
			bz.valley = append(bz.valley, down)
		case '<':
			bz.valley = append(bz.valley, left)
		case '>':
			bz.valley = append(bz.valley, right)
		}
	}
	bz.w = len(bz.valley) / bz.h
	return bz
}

func part1(bz blizzards) int {
	start, end := 1, len(bz.valley)-2
	_, t := bfsBlizzardWalk(bz, start, end)
	return t
}

func part2(bz blizzards) int {
	start, end := 1, len(bz.valley)-2
	bz, t1 := bfsBlizzardWalk(bz, start, end)
	bz, t2 := bfsBlizzardWalk(bz, end, start)
	bz, t3 := bfsBlizzardWalk(bz, start, end)
	return t1 + t2 + t3
}

// using BFS, we can calculate the next blizzards state once per minute
// and keep track of all the possible current positions
func bfsBlizzardWalk(bz blizzards, start, end int) (blizzards, int) {
	steps := []int{0, -bz.w, bz.w, -1, 1} // wait, up, down, left, right
	stack := map[int]struct{}{start: {}}
	for t := 0; ; t++ {
		bz = bz.next()
		newStack := make(map[int]struct{})
		for p := range stack {
			for _, step := range steps {
				np := p + step
				if np >= 0 && np < len(bz.valley) && bz.valley[np] == clear {
					if np == end {
						return bz, t + 1
					}
					newStack[np] = struct{}{}
				}
			}
		}
		stack = newStack
	}
}

func (bz blizzards) next() blizzards {
	nz := bz.cleanCopy()
	for i := range bz.valley {
		v := bz.valley[i]
		if v == clear || v == wall {
			continue
		}
		if v&up == up {
			ni := i - bz.w
			if bz.valley[ni] == wall {
				ni += (bz.w * (bz.h - 2))
			}
			nz.valley[ni] |= up
		}
		if v&down == down {
			ni := i + bz.w
			if bz.valley[ni] == wall {
				ni -= (bz.w * (bz.h - 2))
			}
			nz.valley[ni] |= down
		}
		if v&left == left {
			ni := i - 1
			if bz.valley[ni] == wall {
				ni += bz.w - 2
			}
			nz.valley[ni] |= left
		}
		if v&right == right {
			ni := i + 1
			if bz.valley[ni] == wall {
				ni -= bz.w - 2
			}
			nz.valley[ni] |= right
		}
	}
	return nz
}

func (bz blizzards) cleanCopy() blizzards {
	cp := blizzards{w: bz.w, h: bz.h, valley: make([]state, len(bz.valley))}
	// only copy walls (clear spaces are the default 0 value)
	for i := range bz.valley {
		if bz.valley[i] == wall {
			cp.valley[i] = wall
		}
	}
	return cp
}

func (bz blizzards) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintf("h=%d w=%d\n", bz.h, bz.w))
LOOP:
	for y := 0; y < bz.h; y++ {
		for x := 0; x < bz.w; x++ {
			i := y*bz.w + x
			if i >= len(bz.valley) {
				break LOOP
			}
			v := bz.valley[i]
			switch {
			case v == clear:
				s.WriteByte('.')
			case v == wall:
				s.WriteByte('#')
			case bits.OnesCount(uint(v)) > 1:
				s.WriteByte('0' + byte(bits.OnesCount(uint(v))))
			case v == up:
				s.WriteByte('^')
			case v == down:
				s.WriteByte('v')
			case v == left:
				s.WriteByte('<')
			case v == right:
				s.WriteByte('>')
			default:
				s.WriteByte(' ')
			}
		}
		s.WriteByte('\n')
	}
	return s.String()
}
