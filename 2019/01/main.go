package main

import (
	"strconv"
	"bufio"
	"os"
	"fmt"
)

func main() {
	masses := parseInput("input.txt")

	// Part 1
	part1 := 0
	for _, mass := range masses {
		part1 += (mass / 3 - 2)
	}
	fmt.Println("Part 1", part1)

	// Part 2
	part2 := 0
	for _, mass := range masses {
		part2 += getFuelRequirement(mass)
	}
	fmt.Println("Part 2", part2)
}

func getFuelRequirement(mass int) int {
	fuel := mass / 3 - 2
	if fuel > 0 {
		return fuel + getFuelRequirement(fuel)
	}
	return 0
}

func parseInput(filename string) []int {
	file, _ := os.Open(filename)
	masses := make([]int, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		mass, _ := strconv.Atoi(line)
		masses = append(masses, mass)
	}
	return masses
}
