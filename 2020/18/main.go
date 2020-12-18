package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/solinox/advent-of-code/2020/pkg/input"
)

func main() {
	problems := input.StringSlice(os.Stdin)

	t0 := time.Now()
	log.Println(part1(problems), time.Since(t0))

	t0 = time.Now()
	log.Println(part2(problems), time.Since(t0))
}

func part1(problems []string) int {
	sum := 0
	for _, problem := range problems {
		answer, _ := solve(problem, false)
		sum += answer
	}
	return sum
}

func part2(problems []string) int {
	sum := 0
	for _, problem := range problems {
		answer, _ := solve(problem, true)
		sum += answer
	}
	return sum
}

type opFunc func(int, int) int

var (
	add opFunc = func(x, y int) int {
		return x + y
	}
	multiply opFunc = func(x, y int) int {
		return x * y
	}
)

// second return value is how many bytes were read
func solve(problem string, part2 bool) (int, int) {
	op := add
	bufferedMults := make([]int, 0)
	expr, next := 0, 0
	cursor := 0
EXPRLOOP:
	for cursor < len(problem) {
		switch problem[cursor] {
		case ' ':
			cursor++
			continue
		case '+':
			cursor++
			op = add
		case '*':
			cursor++
			op = multiply
			if part2 {
				bufferedMults = append(bufferedMults, expr)
				expr = 0
				op = add
			}
		case ')':
			cursor++
			break EXPRLOOP
		case '(':
			cursor++
			innerExpr, n := solve(problem[cursor:], part2)
			cursor += n
			next = innerExpr
			expr = op(expr, next)
		default:
			n, err := fmt.Sscanf(problem[cursor:], "%d", &next)
			if err != nil {
				panic(err)
			}
			cursor += n
			expr = op(expr, next)
		}
	}
	// will run for part2 only
	for _, buf := range bufferedMults {
		expr = multiply(expr, buf)
	}
	return expr, cursor
}
