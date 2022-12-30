package main

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/solinox/advent-of-code/2022/util"
)

//go:embed input.txt
var input string

// 0 = Ore, 1 = Clay, 2 = Obsidian, 3 = Geode
type vec [4]int

type Blueprint struct {
	ID     int
	Robots [4]vec
}

func main() {
	inFormat := `Blueprint %d: Each ore robot costs %d ore. Each clay robot costs %d ore. Each obsidian robot costs %d ore and %d clay. Each geode robot costs %d ore and %d obsidian.`
	bps := util.ParseLines(strings.NewReader(input), func(s string) Blueprint {
		var bp Blueprint
		fmt.Sscanf(s, inFormat, &bp.ID, &bp.Robots[0][0], &bp.Robots[1][0], &bp.Robots[2][0], &bp.Robots[2][1], &bp.Robots[3][0], &bp.Robots[3][2])
		return bp
	})
	util.RunTimed(part1, bps)
	util.RunTimed(part2, bps)
}

func part1(bps []Blueprint) int {
	res := 0
	for i := range bps {
		res += (bps[i].ID * maxGeodes(bps[i], 24))
	}
	return res
}

func part2(bps []Blueprint) int {
	res := 1
	for i := 0; i < util.Min(len(bps), 3); i++ {
		res *= maxGeodes(bps[i], 32)
	}
	return res
}

func (v vec) Add(u vec) vec  { return vec{v[0] + u[0], v[1] + u[1], v[2] + u[2], v[3] + u[3]} }
func (v vec) Sub(u vec) vec  { return vec{v[0] - u[0], v[1] - u[1], v[2] - u[2], v[3] - u[3]} }
func (v vec) GTE(u vec) bool { return v[0] >= u[0] && v[1] >= u[1] && v[2] >= u[2] && v[3] >= u[3] }

// checks all possible states, with some exceptions
// skips branches where # bots for a specific material is greater than how many of that mat can be spent per minute
// skips branches where it could have made the bot last turn and chose not to
// skips branches where best possible outcome (creating geode bot for every remaining minute) is less than current best
// skips branches where it could have built a geode bot this minute and chose not to
func maxGeodes(bp Blueprint, maxT int) int {
	mats := vec{0, 0, 0, 0}
	bots := vec{1, 0, 0, 0}
	max := vec{0, 0, 0, 0}
	for i := range bp.Robots {
		for j := range max {
			if bp.Robots[i][j] > max[j] {
				max[j] = bp.Robots[i][j]
			}
		}
	}
	type state struct {
		mats, bots vec
	}
	type item struct {
		state state
		prev  int
	}
	var best int
	prev := make(map[state]struct{})
	stack := map[item]struct{}{{state: state{mats, bots}, prev: -2}: {}}
	for t := 1; t < maxT; t++ {
		// fmt.Printf("t=%d\n", t)
		newStack := make(map[item]struct{})
		for s := range stack {
			for i := 3; i >= 0; i-- {
				if !s.state.mats.GTE(bp.Robots[i]) {
					continue
				}
				if s.prev == -1 && s.state.mats.Sub(s.state.bots).GTE(bp.Robots[i]) {
					continue
				}
				mats := s.state.mats.Sub(bp.Robots[i])
				bots := s.state.bots
				mats = mats.Add(bots)
				if mats[3] > best {
					best = mats[3]
				}
				if n := maxT - t; mats[3]+n*(n+1)/2 < best {
					continue
				}
				bots[i]++
				if i < 3 && bots[i] > max[i] {
					continue
				}
				item := item{state: state{mats, bots}, prev: i}
				if _, ok := prev[item.state]; !ok {
					newStack[item] = struct{}{}
					prev[item.state] = struct{}{}
				}
				if i == 3 {
					break // building geode is guaranteed best
				}
			}
			mats := s.state.mats.Add(s.state.bots)
			if mats[3] > best {
				best = mats[3]
			}
			if n := maxT - t; mats[3]+n*(n+1)/2 < best {
				continue
			}
			item := item{state: state{mats, s.state.bots}, prev: -1}
			if _, ok := prev[item.state]; !ok {
				newStack[item] = struct{}{}
				prev[item.state] = struct{}{}
			}
			delete(stack, s)
		}
		stack = newStack
	}
	// last minute
	for s := range stack {
		mats := s.state.mats.Add(s.state.bots)
		if mats[3] > best {
			best = mats[3]
		}
	}
	return best
}
