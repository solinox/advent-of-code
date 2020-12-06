package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/solinox/advent-of-code/2020/pkg/input"
)

func main() {
	groups := input.SplitString(os.Stdin, "\n\n")

	t0 := time.Now()
	log.Println(part1(groups), time.Since(t0))

	t0 = time.Now()
	log.Println(part2(groups), time.Since(t0))
}

func part1(groups []string) int {
	sum := 0
	for _, group := range groups {
		sum += countAnswers1(group)
	}
	return sum
}

func countAnswers1(group string) int {
	answers := make(map[byte]struct{})
	for i := range group {
		if group[i] != '\n' {
			answers[group[i]] = struct{}{}
		}
	}
	return len(answers)
}

func part2(groups []string) int {
	sum := 0
	for _, group := range groups {
		sum += countAnswers2(group)
	}
	return sum
}

func countAnswers2(group string) int {
	answers := make(map[byte]int)
	people := strings.Split(group, "\n")
	if len(people) == 1 {
		return len(people[0])
	}
	for _, person := range people {
		for i := range person {
			answers[person[i]]++
		}
	}
	sum := 0
	for _, v := range answers {
		if v == len(people) {
			sum++
		}
	}
	return sum
}
