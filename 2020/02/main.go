package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/solinox/advent-of-code/2020/pkg/input"
)

type passwordEntry struct {
	Policy   passwordPolicy
	Password string
}

type passwordPolicy struct {
	N1, N2 int
	Letter byte
}

func main() {
	lines := input.StringSlice(os.Stdin)
	entries := make([]passwordEntry, 0, len(lines))
	for _, line := range lines {
		entry := passwordEntry{}
		n, err := fmt.Sscanf(line, "%d-%d %c: %s", &entry.Policy.N1, &entry.Policy.N2, &entry.Policy.Letter, &entry.Password)
		if err != nil {
			log.Fatalln(line, n, err)
		}
		entries = append(entries, entry)
	}

	log.Println(part1(entries))

	log.Println(part2(entries))
}

func part1(entries []passwordEntry) int {
	validCount := 0
	for _, entry := range entries {
		n := strings.Count(entry.Password, string([]byte{entry.Policy.Letter}))
		if n >= entry.Policy.N1 && n <= entry.Policy.N2 {
			validCount++
		}
	}
	return validCount
}

func part2(entries []passwordEntry) int {
	validCount := 0
	for _, entry := range entries {
		if entry.Policy.N2 > len(entry.Password) {
			continue
		}
		c1, c2 := entry.Password[entry.Policy.N1-1], entry.Password[entry.Policy.N2-1]
		if c1 != c2 && (c1 == entry.Policy.Letter || c2 == entry.Policy.Letter) {
			validCount++
		}
	}
	return validCount
}
