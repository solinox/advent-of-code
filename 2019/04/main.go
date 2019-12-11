package main

import (
	"fmt"
)

func main() {
	min, max := 125730, 579381

	// Part 1
	passwords := getPossiblePasswords(min, max)
	fmt.Println("Part 1", len(passwords))

	// Part 2
	part2 := 0
	for _, pass := range passwords {
		if isPart2Valid(pass) {
			part2++
		}
	}
	fmt.Println("Part 2", part2)
}

func getPossiblePasswords(min, max int) []int {
	passwords := make([]int, 0)	
	for n := next(min); n <= max; n = next(n) {
		passwords = append(passwords, n)
	}
	return passwords
}

func next(n int) int {
	n++
	for !isValid(n) {
		n++;
	}
	return n
}

func isValid(n int) bool {
	equalAdjacent := false
	prev := 10
	for n > 0 {
		digit := n % 10
		if digit > prev {
			return false
		}
		if digit == prev {
			equalAdjacent = true
		}
		prev = digit
		n /= 10
	}
	return equalAdjacent
}

func isPart2Valid(n int) bool {
	prev := 10
	adjacent := 0
	for n > 0 {
		digit := n % 10
		if digit != prev && adjacent == 1 {
			return true
		}
		if digit == prev {
			adjacent++
		} else {
			adjacent = 0
		} 
		prev = digit
		n /= 10
	}
	return adjacent == 1
}