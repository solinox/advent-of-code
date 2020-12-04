package main

import (
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/solinox/advent-of-code/2020/pkg/input"
)

func main() {
	passports := input.SplitString(os.Stdin, "\n\n")

	log.Println(part1(passports, "byr:", "iyr:", "eyr:", "hgt:", "hcl:", "ecl:", "pid:"))

	log.Println(part2(passports,
		regexp.MustCompile(`\bbyr:(19[2-9][0-9]|200[0-2])\b`),
		regexp.MustCompile(`\biyr:(2020|201[0-9])\b`),
		regexp.MustCompile(`\beyr:(2030|202[0-9])\b`),
		regexp.MustCompile(`\bhgt:((1[5-8][0-9]|19[0-3])cm|(59|6[0-9]|7[0-6])in)\b`),
		regexp.MustCompile(`\bhcl:#[0-9a-f]{6}\b`),
		regexp.MustCompile(`\becl:(amb|blu|brn|gry|grn|hzl|oth)\b`),
		regexp.MustCompile(`\bpid:\d{9}\b`),
	))
}

func part1(passports []string, fields ...string) int {
	valid := 0
PASSPORT:
	for _, passport := range passports {
		for _, field := range fields {
			if !strings.Contains(passport, field) {
				continue PASSPORT
			}
		}
		valid++
	}
	return valid
}

func part2(passports []string, patterns ...*regexp.Regexp) int {
	valid := 0
PASSPORT:
	for _, passport := range passports {
		for _, pattern := range patterns {
			if !pattern.MatchString(passport) {
				continue PASSPORT
			}
		}
		valid++
	}
	return valid
}
