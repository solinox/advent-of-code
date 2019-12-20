package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reactions := parseInput("input.txt")

	// Part 1
	part1, overflow1 := findNumInput(Chemical{"FUEL", 1}, "ORE", make(map[string]int), reactions)
	fmt.Println("Part 1:", part1, overflow1)

	// Part 2
	part2 := numFuelFromOre(1e12, reactions)
	fmt.Println("Part 2:", part2)
}

func numFuelFromOre(ore int, reactions []Rx) int {
	fuel := 1 << 14
	min, max := 0, -1
	for max != min {
		x, _ := findNumInput(Chemical{"FUEL", fuel}, "ORE", make(map[string]int), reactions)
		if x < ore {
			min = fuel + 1
		} else if x > ore {
			max = fuel - 1
		} else {
			break
		}
		// keep doubling fuel until we have a max
		if max < 0 {
			fuel *= 2
		} else {
			fuel = (max + min) / 2
		}
	}
	return fuel
}

func findNumInput(desiredOutput Chemical, base string, overflow map[string]int, reactions []Rx) (int, map[string]int) {
	if desiredOutput.Name == base || desiredOutput.Quantity == 0 {
		return desiredOutput.Quantity, overflow
	}

	// Find the rx which gives output
	// Instructions guarantee there is only one
	var reaction Rx
	for _, rx := range reactions {
		if rx.Outputs[0].Name == desiredOutput.Name {
			reaction = rx
			break
		}
	}
	// Figure out number of times to do reaction
	// In order to get desired output quantity
	times := desiredOutput.Quantity / reaction.Outputs[0].Quantity
	if desiredOutput.Quantity%reaction.Outputs[0].Quantity != 0 {
		times++ // round up for uneven ratios
	}
	actualOutput := Chemical{reaction.Outputs[0].Name, reaction.Outputs[0].Quantity * times}
	inputQty := 0
	for _, chem := range reaction.Inputs {
		chem.Quantity *= times
		if v, ok := overflow[chem.Name]; ok && v > 0 {
			if chem.Quantity > v {
				chem.Quantity -= v
				overflow[chem.Name] = 0
			} else {
				overflow[chem.Name] -= chem.Quantity
				chem.Quantity = 0
			}
		}
		inputForChem, extra := findNumInput(chem, base, overflow, reactions)
		overflow = extra
		inputQty += inputForChem
	}
	if over := actualOutput.Quantity - desiredOutput.Quantity; over > 0 {
		overflow[actualOutput.Name] += over
	}
	return inputQty, overflow
}

type Chemical struct {
	Name     string
	Quantity int
}

type Rx struct {
	Inputs  []Chemical
	Outputs []Chemical
}

func parseInput(filename string) []Rx {
	file, _ := os.Open(filename)
	scanner := bufio.NewScanner(file)
	rxs := make([]Rx, 0)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			equation := strings.Split(line, "=>")
			in, out := equation[0], equation[1]
			rx := Rx{Inputs: parseChemicals(in), Outputs: parseChemicals(out)}
			rxs = append(rxs, rx)
		}
	}
	return rxs
}

func parseChemicals(s string) []Chemical {
	chems := make([]Chemical, 0)
	items := strings.Split(s, ",")
	for _, item := range items {
		chem := strings.TrimSpace(item)
		var quantity int
		var name string
		fmt.Sscanf(chem, "%d %s", &quantity, &name)
		chems = append(chems, Chemical{Quantity: quantity, Name: name})
	}
	return chems
}
