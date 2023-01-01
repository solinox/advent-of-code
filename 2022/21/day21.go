package main

import (
	_ "embed"
	"strconv"
	"strings"

	"github.com/solinox/advent-of-code/2022/util"
)

//go:embed input.txt
var input string

func main() {
	monkeys := util.ParseLinesReduce(strings.NewReader(input), func(agg map[string]string, s string) map[string]string {
		kv := strings.Split(s, ": ")
		agg[kv[0]] = kv[1]
		return agg
	}, make(map[string]string))
	util.RunTimed(part1, monkeys)
	util.RunTimed(part2, monkeys)
}

func part1(monkeys map[string]string) int {
	return number(monkeys, "root")
}

func part2(monkeys map[string]string) int {
	job := strings.Split(monkeys["root"], " ")
	l, r := job[0], job[2]
	// in both test and my input, "humn" only affects the left number
	// and I don't feel like adding extra validation to this right now
	// also, changing humn has an approximate stepped-linear change on left
	// eg. no change for x steps, then it changes by y
	// so figure out the dx and dy here, then interpolate to get a rough estimate of the desired humn number
	// repeat as needed to get closer and closer, until it is exactly matching
	left, right, humn := number(monkeys, l), number(monkeys, r), util.ParseInt(monkeys["humn"])
	for left != right {
		l1, l2 := left, left
		dl1 := -1
		// step humn in both directions (-/+) to get the full dx between changes in left
		for l1 == left {
			l1 = numberWithHumn(monkeys, l, humn+dl1)
			if l1 != left {
				break
			}
			dl1--
		}
		dl2 := 1
		for l2 == left {
			l2 = numberWithHumn(monkeys, l, humn+dl2)
			if l2 != left {
				break
			}
			dl2++
		}
		dx, dy := dl2-dl1-1, left-l1
		// for every humn changes dx, left changes dy
		humn = humn + (right-left)/dy*dx
		left = numberWithHumn(monkeys, l, humn)
	}
	return humn
}

func numberWithHumn(monkeys map[string]string, name string, humn int) int {
	monkeys["humn"] = strconv.Itoa(humn)
	return number(monkeys, name)
}

func number(monkeys map[string]string, name string) int {
	v := monkeys[name]
	if n, err := strconv.Atoi(v); err == nil {
		return n
	}
	job := strings.Split(v, " ")
	n1, op, n2 := number(monkeys, job[0]), job[1], number(monkeys, job[2])
	switch op {
	case "+":
		return n1 + n2
	case "-":
		return n1 - n2
	case "*":
		return n1 * n2
	case "/":
		return n1 / n2
	}
	panic(op)
}
