package main

import (
	_ "embed"
	"strings"

	"github.com/solinox/advent-of-code/2022/util"
)

//go:embed input.txt
var input string

func main() {
	in := util.ParseLines(strings.NewReader(input), func(l string) string { return l })
	util.RunTimed(part1, in)
	util.RunTimed(part2, in)
}

func doOp(op string, reg int) ([]int, int) {
	if op == "noop" {
		return []int{reg}, reg
	} else if op[:4] == "addx" {
		v := util.ParseInt(op[5:])
		return []int{reg, reg}, reg + v
	}
	return nil, reg
}

func calcCycles(ops []string) []int {
	cycles := make([]int, 0, 240)
	reg := 1
	for _, op := range ops {
		newCycles, nextReg := doOp(op, reg)
		reg = nextReg
		cycles = append(cycles, newCycles...)
	}
	return cycles
}

func part1(ops []string) int {
	cycles := calcCycles(ops)
	return cycles[19]*20 + cycles[59]*60 + cycles[99]*100 + cycles[139]*140 + cycles[179]*180 + cycles[219]*220
}

func part2(ops []string) string {
	cycles := calcCycles(ops)
	var s strings.Builder
	for i := 0; i < len(cycles); i++ {
		sprite, pix := cycles[i], i%40
		if sprite+1 >= pix && pix >= sprite-1 {
			s.WriteByte('#')
		} else {
			s.WriteByte('.')
		}
		if (i+1)%40 == 0 {
			s.WriteByte('\n')
		}
	}
	return s.String()
}
