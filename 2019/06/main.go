package main

import (
	"os"
	"bufio"
	"strings"
	"fmt"
)

func main() {
	orbits := parseInput("input.txt")

	// Part 1
	part1 := 0
	for _, parent := range orbits {
		part1++ // direct orbit
		part1 += countIndirectOrbits(parent, orbits)
	}
	fmt.Println("Part 1", part1)

	// Part 2
	part2 := degreesOfSeparation("YOU", "SAN", orbits)
	fmt.Println("Part 2", part2)

}

func degreesOfSeparation(system1, system2 string, orbits map[string]string) int {
	path1 := walk(system1, orbits, make([]string, 0))
	path2 := walk(system2, orbits, make([]string, 0))
	path1, path2 = reverse(path1), reverse(path2)
	for i := range path1 {
		if path1[i] != path2[i] {
			// remove parts that are similar between the two paths
			path1, path2 = path1[i:], path2[i:]
			break
		}
	}
	return len(path1) + len(path2)
}

func reverse(slice []string) []string {
	for i, j := 0, len(slice)-1; i < len(slice)/2; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
	return slice
}

func walk(system string, orbits map[string]string, path []string) []string {
	if parent, exists := orbits[system]; exists {
		path = append(path, parent)
		return walk(parent, orbits, path)
	}
	return path
}

func countIndirectOrbits(system string, orbits map[string]string) int {
	if parent, exists := orbits[system]; exists {
		return 1 + countIndirectOrbits(parent, orbits)
	}
	return 0
}

func parseInput(filename string) map[string]string {
	orbits := make(map[string]string)
	file, _ := os.Open(filename)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		systems := strings.Split(line, ")")
		parent, child := systems[0], systems[1]
		orbits[child] = parent
	}
	return orbits
}