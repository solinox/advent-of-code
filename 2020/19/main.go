package main

import (
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/solinox/advent-of-code/2020/pkg/input"
)

func main() {
	sections := input.SplitString(os.Stdin, "\n\n")
	rules, messages := strings.Split(sections[0], "\n"), strings.Split(sections[1], "\n")

	t0 := time.Now()
	log.Println(part1(rules, messages), time.Since(t0))

	t0 = time.Now()
	log.Println(part2(rules, messages), time.Since(t0))
}

func part1(rulesString, messages []string) int {
	rules := buildRules(rulesString, "0", make(map[string]string))
	pattern := regexp.MustCompile("^" + rules["0"] + "$")
	valid := 0
	for _, msg := range messages {
		if pattern.MatchString(msg) {
			valid++
		}
	}
	return valid
}

func part2(rulesString, messages []string) int {
	// Instead of making the buildRules method (developed in part1) support recursive rules
	// I noticed that rule `8: 42 | 42 8` is equivalent regexp to (<42>)+
	// and that rule `11: 42 31 | 42 11 31` is equivalent regexp to (<42>){n}(<31>){n}, where n >= 1
	// therefore, rule `0: 8 11` is equivalent regexp to (<42>)+(<42>){n}(<31>){n}, where n >= 1
	// which simplifies to (<42>){m,}(<31>){n}, where n >= 1 and m = n+1
	rules := buildRules(rulesString, "42", make(map[string]string))
	rules = buildRules(rulesString, "31", rules)

	rule0Template := "((" + rules["42"] + "){m,})((" + rules["31"] + "){n})"

	// try out different values of n until number of matches reduces to 0
	// aggregate the sum of the matches
	sum := 0
	for n := 1; ; n++ {
		rules["0"] = strings.ReplaceAll(strings.ReplaceAll(rule0Template, "n", strconv.Itoa(n)), "m", strconv.Itoa(n+1))

		pattern := regexp.MustCompile("^" + rules["0"] + "$")
		valid := 0
		for _, msg := range messages {
			if pattern.MatchString(msg) {
				valid++
			}
		}
		if valid == 0 {
			break
		}
		sum += valid
	}
	return sum
}

func buildRules(rulesString []string, key string, memo map[string]string) map[string]string {
	rulesMap := make(map[string]string)
	for _, rulesStr := range rulesString {
		i := strings.IndexByte(rulesStr, ':')
		ruleNum := rulesStr[:i]
		rule := rulesStr[i+1:]
		rulesMap[ruleNum] = rule
	}

	var parseRule func(string, map[string]string) string
	parseRule = func(rule string, memo map[string]string) string {
		rule = strings.TrimSpace(rule)
		if rule[0] == '"' {
			return rule[1:strings.LastIndexByte(rule, '"')]
		}

		fields := strings.Fields(rule)
		numParens := 1
		pattern := ""
		for _, f := range fields {
			if f == "|" {
				pattern += strings.Repeat(")", numParens) + "|" + strings.Repeat("(", numParens)
				numParens++
				continue
			}
			if val, ok := memo[f]; ok {
				pattern += val
				continue
			}
			subRule := rulesMap[f]
			subPattern := parseRule(subRule, memo)
			memo[f] = subPattern
			pattern += subPattern
		}
		return strings.Repeat("(", numParens) + pattern + strings.Repeat(")", numParens)
	}

	pattern := parseRule(rulesMap[key], memo)
	memo[key] = pattern
	return memo
}
