package main

import (
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/solinox/advent-of-code/2020/pkg/input"
)

type rule struct {
	Name     string
	Inner    []rule
	Min, Max int
}

func (r rule) Contains(n int) bool {
	if len(r.Inner) > 0 {
		for _, inner := range r.Inner {
			if inner.Contains(n) {
				return true
			}
		}
		return false
	}
	return n >= r.Min && n <= r.Max
}

func (r rule) String() string {
	return r.Name
}

var rulesRegex = regexp.MustCompile(`^(.+): (\d+)-(\d+) or (\d+)-(\d+)$`)

func main() {
	sections := input.SplitString(os.Stdin, "\n\n")
	rulesInput := strings.Split(sections[0], "\n")
	rules := make([]rule, len(rulesInput))
	for i := range rulesInput {
		matches := rulesRegex.FindAllStringSubmatch(rulesInput[i], -1)
		rules[i] = rule{
			Name: matches[0][1],
			Inner: []rule{
				{Name: matches[0][1], Min: input.MustAtoi(matches[0][2]), Max: input.MustAtoi(matches[0][3])},
				{Name: matches[0][1], Min: input.MustAtoi(matches[0][4]), Max: input.MustAtoi(matches[0][5])},
			},
		}
	}
	myTicket := input.SplitInt(strings.NewReader(strings.Split(sections[1], "\n")[1]), ",")

	otherTicketInput := strings.Split(sections[2], "\n")[1:]
	otherTickets := make([][]int, len(otherTicketInput))
	for i := range otherTicketInput {
		otherTickets[i] = input.SplitInt(strings.NewReader(otherTicketInput[i]), ",")
	}

	t0 := time.Now()
	log.Println(part1(rules, otherTickets), time.Since(t0))

	t0 = time.Now()
	log.Println(part2(rules, myTicket, otherTickets), time.Since(t0))
}

func part1(rules []rule, tickets [][]int) int {
	errorRate, _ := filterTickets(rules, tickets)
	return errorRate
}

func part2(rules []rule, myTicket []int, otherTickets [][]int) int {
	_, otherTickets = filterTickets(rules, otherTickets)
	fields := determineFields(rules, append(otherTickets, myTicket))
	n := 1
	for i, field := range fields {
		if strings.HasPrefix(field.Name, "departure") {
			n *= myTicket[i]
		}
	}
	return n
}

func filterTickets(rules []rule, tickets [][]int) (int, [][]int) {
	errorRate := 0
	validTickets := make([][]int, 0, len(tickets))
	for _, ticket := range tickets {
		valid := true
	TICKET:
		for _, n := range ticket {
			for _, rule := range rules {
				if rule.Contains(n) {
					continue TICKET
				}
			}
			valid = false
			errorRate += n
		}
		if valid {
			validTickets = append(validTickets, ticket)
		}
	}
	return errorRate, validTickets
}

func determineFields(rules []rule, tickets [][]int) []rule {
	// transpose tickets matrix
	ticketFields := make([][]int, len(tickets[0]))
	for i := 0; i < len(tickets[0]); i++ {
		for _, ticket := range tickets {
			ticketFields[i] = append(ticketFields[i], ticket[i])
		}
	}

	possibleRules := make([][]rule, len(ticketFields))

	for i := range ticketFields {
		possibleRules[i] = getValidRules(rules, ticketFields[i])
	}

	// brute force
	return getOrderOfRules(possibleRules)
}

func getValidRules(rules []rule, vals []int) []rule {
	validRules := make([]rule, 0, len(rules))
RULES:
	for _, rule := range rules {
		for _, n := range vals {
			if !rule.Contains(n) {
				continue RULES
			}
		}
		validRules = append(validRules, rule)
	}
	return validRules
}

func getOrderOfRules(possibleRules [][]rule) []rule {
	var tryRules func(possibleRules [][]rule, rules []rule) []rule
	tryRules = func(possibleRules [][]rule, rules []rule) []rule {
		if len(rules) == len(possibleRules) {
			return rules
		}
		for _, rule := range possibleRules[len(rules)] {
			if contains(rules, rule) {
				continue
			}
			newRules := tryRules(possibleRules, append(rules, rule))
			if len(newRules) == len(possibleRules) {
				return newRules
			}
		}
		return rules

	}
	return tryRules(possibleRules, nil)
}

func contains(rules []rule, rule rule) bool {
	for _, r := range rules {
		if r.Name == rule.Name {
			return true
		}
	}
	return false
}
