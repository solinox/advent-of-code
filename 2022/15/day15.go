package main

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/solinox/advent-of-code/2022/util"
)

//go:embed input.txt
var input string

type Vector = util.Vector
type Range = util.Range
type Sensor struct {
	Loc, Beacon Vector
}

func main() {
	sensors := util.ParseLines(strings.NewReader(input), func(line string) Sensor {
		var s Sensor
		fmt.Sscanf(line, "Sensor at x=%d, y=%d: closest beacon is at x=%d, y=%d", &s.Loc.X, &s.Loc.Y, &s.Beacon.X, &s.Beacon.Y)
		return s
	})
	util.RunTimed(part1, sensors)
	util.RunTimed(part2, sensors)
}

func part1(sensors []Sensor) int {
	rs := coverageY(sensors, 2000000)
	coverage := 0
	for _, r := range rs {
		coverage += r.Dist()
	}
	return coverage
}

// current brute-force iterates over each y, finding the first row which is not covered from 0-4000000
// current takes a couple seconds to run
// optimization idea: get points just outside each sensor's perimeter that are within the bounds
// and check those points until one of them is not covered by any sensor
func part2(sensors []Sensor) int {
	const tuningMult = 4000000
	bounds := Range{Min: 0, Max: 4000000}
	for y := bounds.Min; y <= bounds.Max; y++ {
		rs := coverageY(sensors, y)
		r, ok := notCovered(rs, bounds)
		if ok {
			return r.Min*tuningMult + y
		}
	}
	return -1
}

func coverageY(sensors []Sensor, y int) []Range {
	rs := make([]Range, 0)
	for _, s := range sensors {
		if r, ok := s.coverageAtY(y); ok {
			rs = append(rs, r)
		}
	}
	return util.MergeRanges(rs)
}

func (s Sensor) coverageAtY(y int) (Range, bool) {
	d, dy := s.Loc.Dist(s.Beacon), util.Abs(y-s.Loc.Y)
	if d < dy {
		// beacon is too close, or sensor is too far away, to cover anything at y
		return Range{}, false
	}
	return Range{Min: s.Loc.X - (d - dy), Max: s.Loc.X + (d - dy)}, true
}

func notCovered(rs []Range, bounds Range) (Range, bool) {
	r := bounds
	for i := range rs {
		if rs[i].Min <= r.Min && rs[i].Max >= r.Max {
			return r, false
		}
		if rs[i].Min <= r.Min && rs[i].Max >= r.Min {
			r.Min = rs[i].Max + 1
			if r.Min > r.Max {
				return r, false
			}
		}
	}
	return r, r.Min <= r.Max
}
