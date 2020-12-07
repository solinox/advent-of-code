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
	rules := input.StringSlice(os.Stdin)

	t0 := time.Now()
	log.Println(part1(rules), time.Since(t0))

	t0 = time.Now()
	log.Println(part2(rules), time.Since(t0))
}

type bag struct {
	Name     string
	Parents  []string
	Children []innerBag
}

type innerBag struct {
	Name       string
	Quantity   int
	Multiplier int
}

var rulePattern = regexp.MustCompile(`(\w+\s+\w+) bags contain (.+)\.`)

func part1(rules []string) int {
	ruleTree := parseRules(rules)
	ancestors := make(map[string]struct{})
	stack := ruleTree["shiny gold"].Parents
	for len(stack) > 0 {
		bagName := stack[0]
		stack = stack[1:]
		ancestors[bagName] = struct{}{}
		stack = append(stack, ruleTree[bagName].Parents...)
	}
	return len(ancestors)
}

func part2(rules []string) int {
	ruleTree := parseRules(rules)
	stack := addMultiplier(ruleTree["shiny gold"].Children, 1)
	numBags := 0
	for len(stack) > 0 {
		innerBag := stack[0]
		stack = stack[1:]
		newMultiplier := innerBag.Multiplier * innerBag.Quantity
		numBags += newMultiplier
		stack = append(stack, addMultiplier(ruleTree[innerBag.Name].Children, newMultiplier)...)
	}
	return numBags
}

func addMultiplier(children []innerBag, multiplier int) []innerBag {
	for i := range children {
		children[i].Multiplier = multiplier
	}
	return children
}

func parseRules(rules []string) map[string]bag {
	ruleTree := make(map[string]bag)
	for _, rule := range rules {
		submatches := rulePattern.FindStringSubmatch(rule)
		if len(submatches) < 3 {
			log.Fatalln("regexp error")
		}
		name, children := submatches[1], strings.Split(submatches[2], ",")
		var innerBags []innerBag
		if children[0] != "no other bags" {
			for _, child := range children {
				fields := strings.Fields(child)
				if len(fields) < 3 {
					log.Fatalln("split error")
				}
				n, err := strconv.Atoi(fields[0])
				if err != nil {
					log.Fatalln(rule, child, err)
				}
				childName := fields[1] + " " + fields[2]
				if existingBag, ok := ruleTree[childName]; ok {
					existingBag.Parents = append(existingBag.Parents, name)
					ruleTree[childName] = existingBag
				} else {
					ruleTree[childName] = bag{
						Name:     childName,
						Parents:  []string{name},
						Children: nil,
					}
				}
				innerBags = append(innerBags, innerBag{Name: childName, Quantity: n})
			}
		}
		if existingBag, ok := ruleTree[name]; ok {
			existingBag.Children = innerBags
			ruleTree[name] = existingBag
		} else {
			ruleTree[name] = bag{
				Name:     name,
				Parents:  nil,
				Children: innerBags,
			}
		}
	}
	return ruleTree
}
