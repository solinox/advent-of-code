package main

import (
	_ "embed"
	"strconv"
	"strings"

	"github.com/solinox/advent-of-code/2022/util"
)

//go:embed input.txt
var input string

type Monkey struct {
	Items         []int
	Op            string
	TestDivisible int
	TestTrue      int
	TestFalse     int
	Inspections   int
}
type Monkeys []Monkey

func main() {
	// cannot reuse input because it is modified within each part
	util.RunTimed(part1, parse(input))
	util.RunTimed(part2, parse(input))
}

func part1(monkeys Monkeys) int {
	monkeys.Rounds(20, true)
	return monkeys.Business()
}

func part2(monkeys Monkeys) int {
	monkeys.Rounds(10000, false)
	return monkeys.Business()
}

func parse(input string) Monkeys {
	m := make(Monkeys, 0)
	for _, sect := range strings.Split(input, "\n\n") {
		lines := strings.Split(sect, "\n")
		m = append(m, Monkey{
			Items:         util.ParseDelimited(strings.TrimPrefix(lines[1], "  Starting items: "), ", ", util.ParseInt),
			Op:            strings.TrimPrefix(lines[2], "  Operation: new = old "),
			TestDivisible: util.ParseInt(strings.TrimPrefix(lines[3], "  Test: divisible by ")),
			TestTrue:      util.ParseInt(strings.TrimPrefix(lines[4], "    If true: throw to monkey ")),
			TestFalse:     util.ParseInt(strings.TrimPrefix(lines[5], "    If false: throw to monkey ")),
		})
	}
	return m
}

func (m Monkeys) Business() int {
	inspections := util.SliceFrom(m, func(m Monkey) int { return m.Inspections })
	topTwo := util.MaxN(2, inspections...)
	return topTwo[0] * topTwo[1]
}

func (m Monkeys) SafeModulo() int {
	p := 1
	for _, mo := range m {
		p *= mo.TestDivisible
	}
	return p
}

func (m Monkeys) Rounds(n int, calmAfterInspection bool) {
	for i := 0; i < n; i++ {
		for j := range m {
			m.Turn(j, calmAfterInspection, m.SafeModulo())
		}
	}
}

func (m Monkeys) Turn(i int, calmAfterInspection bool, safeMod int) {
	for _, item := range m[i].Items {
		// inspect
		m[i].Inspections++
		op := strings.Replace(m[i].Op, "old", strconv.Itoa(item), -1)
		n := util.ParseInt(op[2:])
		if op[0] == '*' {
			item *= n
		} else if op[0] == '+' {
			item += n
		} else {
			panic("unsupported input" + op)
		}
		if calmAfterInspection {
			// after inspection calming down
			item /= 3
		}
		// manage worry level without affecting results of divisible test
		item %= safeMod
		// test worry level
		result := item%m[i].TestDivisible == 0
		throwTo := m[i].TestTrue
		if !result {
			throwTo = m[i].TestFalse
		}
		// throw item to the next monkey
		m[throwTo].Items = append(m[throwTo].Items, item)
	}
	// current monkey has thrown all their items elsewhere
	m[i].Items = nil
}
