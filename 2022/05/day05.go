package main

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/solinox/advent-of-code/2022/util"
)

//go:embed input.txt
var input string

type Input struct {
	Stacks     [][]byte
	Directions []Direction
}

type Direction struct {
	N, From, To int
}

func main() {
	// cannot reuse input since input is modified within each part
	util.RunTimed(part1, parse(input))
	util.RunTimed(part2, parse(input))
}

func parse(input string) *Input {
	sections := strings.Split(input, "\n\n")
	stackData, directionData := sections[0], sections[1]

	stackLines := util.Reverse(strings.Split(stackData, "\n"))
	nStacks, stackLines := (len(stackLines[0])+1)/4, stackLines[1:]
	stacks := make([][]byte, nStacks)

	for _, st := range stackLines {
		for i := range stacks {
			if v := st[1+4*i]; v != ' ' {
				stacks[i] = append(stacks[i], v)
			}
		}
	}

	directions := make([]Direction, 0)

	for _, dir := range strings.Split(directionData, "\n") {
		var d Direction
		fmt.Sscanf(dir, "move %d from %d to %d", &d.N, &d.From, &d.To)
		d.From--
		d.To--
		directions = append(directions, d)
	}

	return &Input{Stacks: stacks, Directions: directions}
}

func (in *Input) move(dir Direction, reverse bool) {
	items := in.Stacks[dir.From][len(in.Stacks[dir.From])-dir.N:]
	if reverse {
		items = util.Reverse(items)
	}
	in.Stacks[dir.To] = append(in.Stacks[dir.To], items...)
	in.Stacks[dir.From] = in.Stacks[dir.From][:len(in.Stacks[dir.From])-dir.N]
}

func output(stacks [][]byte) string {
	out := make([]byte, 0, len(stacks))
	for _, stack := range stacks {
		out = append(out, stack[len(stack)-1])
	}
	return string(out)
}

func part1(in *Input) string {
	for _, dir := range in.Directions {
		in.move(dir, true)
	}
	return output(in.Stacks)
}

func part2(in *Input) string {
	for _, dir := range in.Directions {
		in.move(dir, false)
	}
	return output(in.Stacks)
}
