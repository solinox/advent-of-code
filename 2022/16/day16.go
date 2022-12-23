package main

import (
	_ "embed"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/solinox/advent-of-code/2022/util"
)

//go:embed input.txt
var input string

type Valve struct {
	Name     string
	FlowRate int
	Tunnels  string
}

func main() {
	t0 := time.Now()
	valves := parse(input)
	valves, travelTimes := prepare(valves)
	t1 := time.Now()
	fmt.Println("parse & prep =>", t1.Sub(t0))
	p1 := part1(valves, travelTimes)
	t2 := time.Now()
	fmt.Println(p1, "=>", t2.Sub(t1))
	p2 := part2(valves, travelTimes)
	t3 := time.Now()
	fmt.Println(p2, "=>", t3.Sub(t2))
	fmt.Println("total", t3.Sub(t0))
}

func parse(input string) map[string]Valve {
	v := make(map[string]Valve)
	for _, line := range strings.Split(input, "\n") {
		var name string
		var flowRate int
		fmt.Sscanf(line, "Valve %s has flow rate=%d", &name, &flowRate)
		tunnels := strings.TrimSpace(line[strings.Index(line, "valve")+6:])
		v[name] = Valve{Name: name, FlowRate: flowRate, Tunnels: tunnels}
	}
	return v
}

func prepare(valves map[string]Valve) (map[string]Valve, map[string]int) {
	// get min path between each valve, and then remove valves with no flow rate
	travelTimes := minPathTravelTimes(valves)
	for k, v := range valves {
		if v.FlowRate <= 0 {
			delete(valves, k)
		}
	}
	return valves, travelTimes
}

func part1(valves map[string]Valve, travelTimes map[string]int) int {
	return maxPressureRelease(valves, "AA", 30, travelTimes, make(map[string]bool))
}

func part2(valves map[string]Valve, travelTimes map[string]int) int {
	valvesS := make([]Valve, 0, len(valves))
	for _, v := range valves {
		valvesS = append(valvesS, v)
	}
	cache := make(map[string]int)
	max := 0
	for _, vForMe := range util.AllCombinations(valvesS) {
		valvesForMe, valvesForThee := make(map[string]Valve), make(map[string]Valve)
		me, thee := "", ""
		for _, v := range valvesS {
			if util.Contains(vForMe, v) {
				valvesForMe[v.Name] = v
				me += v.Name
			} else {
				valvesForThee[v.Name] = v
				thee += v.Name
			}
		}
		scoreMe, ok := cache[me]
		if !ok {
			scoreMe = maxPressureRelease(valvesForMe, "AA", 26, travelTimes, make(map[string]bool))
			cache[me] = scoreMe
		}
		scoreThee, ok := cache[thee]
		if !ok {
			scoreThee = maxPressureRelease(valvesForThee, "AA", 26, travelTimes, make(map[string]bool))
			cache[thee] = scoreThee
		}
		if score := scoreMe + scoreThee; score > max {
			fmt.Println(score, me, thee)
			max = score
		}
	}
	return max
}

func minPathTravelTimes(valves map[string]Valve) map[string]int {
	// floyd warshall algorithm, using a map instead of 2d array
	// where the key is NameFrom+NameTo e.g. "AABB" for AA -> BB
	travelTimes := make(map[string]int)
	slice := make([]Valve, 0, len(valves))
	for _, v := range valves {
		slice = append(slice, v)
		travelTimes[v.Name+v.Name] = 0
		for _, t := range strings.Split(v.Tunnels, ", ") {
			travelTimes[v.Name+t] = 1
			travelTimes[t+v.Name] = 1
		}
	}
	for _, vk := range slice {
		for _, vi := range slice {
			for _, vj := range slice {
				dij, okij := travelTimes[vi.Name+vj.Name]
				if !okij {
					dij = math.MaxInt
				}
				dik, okik := travelTimes[vi.Name+vk.Name]
				dkj, okkj := travelTimes[vk.Name+vj.Name]
				if n := dik + dkj; okik && okkj && n < dij {
					travelTimes[vi.Name+vj.Name] = n
					travelTimes[vj.Name+vi.Name] = n
				}
			}
		}
	}
	return travelTimes
}

func maxPressureRelease(valves map[string]Valve, curV string, minutesLeft int, travelTimes map[string]int, open map[string]bool) int {
	if minutesLeft <= 0 {
		return 0
	}
	max := 0
	for k, v := range valves {
		if open[k] || v.FlowRate == 0 {
			continue
		}
		open2 := make(map[string]bool)
		for kk, vv := range open {
			open2[kk] = vv
		}
		open2[k] = true
		dtV := 1 + travelTimes[curV+k]
		score := (minutesLeft - dtV) * v.FlowRate
		potentialMax := score + maxPressureRelease(valves, k, minutesLeft-dtV, travelTimes, open2)
		if potentialMax > max {
			max = potentialMax
		}
	}
	return max
}
