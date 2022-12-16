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
type Line struct {
	V1, V2 Vector
}

func main() {
	sensors := util.ParseLines(strings.NewReader(input), func(line string) Sensor {
		var s Sensor
		fmt.Sscanf(line, "Sensor at x=%d, y=%d: closest beacon is at x=%d, y=%d", &s.Loc.X, &s.Loc.Y, &s.Beacon.X, &s.Beacon.Y)
		return s
	})
	util.RunTimed(part1, sensors)
	util.RunTimed(part2Optimized, sensors)
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

// brute-force solution takes 1.5-3s, iterates over each row within the bounds
// and checks if the range is completely covered by sensors
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

// checks intersections of the lines just outside each sensor perimeter
// runs in <1ms
func part2Optimized(sensors []Sensor) int {
	const tuningMult = 4000000
	bounds := Range{Min: 0, Max: 4000000}
	positiveSlopeLines := make([]Line, 0, len(sensors)*2)
	negativeSlopeLines := make([]Line, 0, len(sensors)*2)
	for _, s := range sensors {
		perimeters := s.BeyondPerimeter()
		positiveSlopeLines = append(positiveSlopeLines, perimeters[0], perimeters[3])
		negativeSlopeLines = append(negativeSlopeLines, perimeters[1], perimeters[2])
	}
	// the problem guarantees there is 1 point in the bounds that is not covered
	// so at least 2 sensor coverages are 1 range short
	// extending the range by one, the perimeters must intersect on the missing beacon
INTS:
	for _, intersection := range intersections(positiveSlopeLines, negativeSlopeLines) {
		if intersection.X >= bounds.Min && intersection.X <= bounds.Max && intersection.Y >= bounds.Min && intersection.Y <= bounds.Max {
			for _, s := range sensors {
				if s.Covers(intersection) {
					continue INTS
				}
			}
			return tuningMult*intersection.X + intersection.Y
		}
	}
	return -1
}

func intersections(pSlopes, nSlopes []Line) []Vector {
	ints := make([]Vector, 0)
	// y = mx + b
	// positive slopes: m=1
	// negative slopes: m=-1
	// (-x + b1 = y) == (y = x + b2)
	// x = (b1-b2)/2
	// y = b2+(b1-b2)/2
	// b1 = y+x
	// b2 = y-x
	for _, pSlope := range pSlopes {
		for _, nSlope := range nSlopes {
			// pSlope and nSlope will always intersect
			b1, b2 := nSlope.V1.Y+nSlope.V1.X, pSlope.V1.Y-pSlope.V1.X
			x, y := (b1-b2)/2, b2+(b1-b2)/2
			ints = append(ints, Vector{Y: y, X: x})
		}
	}
	return ints
}

func (s Sensor) BeyondPerimeter() []Line {
	d := s.Loc.Dist(s.Beacon) + 1
	up, down := Vector{Y: s.Loc.Y - d, X: s.Loc.X}, Vector{Y: s.Loc.Y + d, X: s.Loc.X}
	left, right := Vector{Y: s.Loc.Y, X: s.Loc.X - d}, Vector{Y: s.Loc.Y, X: s.Loc.X + d}
	lines := []Line{
		{V1: left, V2: up},
		{V1: up, V2: right},
		{V1: left, V2: down},
		{V1: down, V2: right},
	}
	return lines
}

func (s Sensor) Covers(v Vector) bool {
	d := s.Loc.Dist(s.Beacon)
	return s.Loc.Dist(v) <= d
}
