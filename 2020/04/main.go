package main

import (
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/solinox/advent-of-code/2020/pkg/input"
)

func main() {
	r1, r2 := input.ScanDelim(os.Stdin, []byte{'\n', '\n'})

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		t0 := time.Now()
		log.Println("Part 1:", part1(r1,
			"byr:", "iyr:", "eyr:", "hgt:", "hcl:", "ecl:", "pid:"),
			time.Since(t0))
		wg.Done()
	}()

	go func() {
		t0 := time.Now()
		log.Println("Part 2:", part2(r2,
			regexp.MustCompile(`\bbyr:(19[2-9][0-9]|200[0-2])\b`),
			regexp.MustCompile(`\biyr:(2020|201[0-9])\b`),
			regexp.MustCompile(`\beyr:(2030|202[0-9])\b`),
			regexp.MustCompile(`\bhgt:((1[5-8][0-9]|19[0-3])cm|(59|6[0-9]|7[0-6])in)\b`),
			regexp.MustCompile(`\bhcl:#[0-9a-f]{6}\b`),
			regexp.MustCompile(`\becl:(amb|blu|brn|gry|grn|hzl|oth)\b`),
			regexp.MustCompile(`\bpid:\d{9}\b`),
		),
			time.Since(t0),
		)
		wg.Done()
	}()

	wg.Wait()
}

func part1(creds <-chan string, fields ...string) int {
	valid := 0
PASSPORT:
	for passport := range creds {
		for _, field := range fields {
			if !strings.Contains(passport, field) {
				continue PASSPORT
			}
		}
		valid++
	}
	return valid
}

func part2(creds <-chan string, patterns ...*regexp.Regexp) int {
	valid := 0
PASSPORT:
	for passport := range creds {
		for _, pattern := range patterns {
			if !pattern.MatchString(passport) {
				continue PASSPORT
			}
		}
		valid++
	}
	return valid
}
