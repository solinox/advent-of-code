package main

import (
	_ "embed"
	"encoding/json"
	"reflect"
	"sort"
	"strings"

	"github.com/solinox/advent-of-code/2022/util"
)

//go:embed input.txt
var input string

type order int

const (
	unknownOrder = iota
	inOrder
	outOfOrder
)

type Packet []any

func (p Packet) String() string {
	data, _ := json.Marshal(p)
	return string(data)
}

type ByPacket []Packet

func (s ByPacket) Len() int           { return len(s) }
func (s ByPacket) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ByPacket) Less(i, j int) bool { return slicesOrder(s[i], s[j]) != outOfOrder }

func main() {
	packets := parse(input)
	util.RunTimed(part1, pairs(packets))
	util.RunTimed(part2, packets)
}

func part1(pairs [][]Packet) int {
	sum := 0
	for i, pair := range pairs {
		if slicesOrder([]any(pair[0]), []any(pair[1])) == inOrder {
			sum += i + 1
		}
	}
	return sum
}

func part2(packets []Packet) int {
	p1, p2 := Packet{[]any{2.0}}, Packet{[]any{6.0}}
	d1, d2 := p1.String(), p2.String()
	packets = append(packets, p1, p2)

	sort.Sort(ByPacket(packets))

	product := 1
	for i := range packets {
		if s := packets[i].String(); s == d1 || s == d2 {
			product *= (i + 1)
		}
	}
	return product
}

func parse(input string) []Packet {
	return util.ParseLines(strings.NewReader(input), func(line string) Packet {
		var p Packet
		json.Unmarshal([]byte(line), &p)
		return p
	})
}

func pairs(packets []Packet) [][]Packet {
	p := make([][]Packet, 0, len(packets)/2)
	for i := 0; i < len(packets); i += 2 {
		p = append(p, []Packet{packets[i], packets[i+1]})
	}
	return p
}

func slicesOrder(left, right []any) order {
	n := util.Min(len(left), len(right))
	for i := 0; i < n; i++ {
		leftSlice, rightSlice := reflect.TypeOf(left[i]).Kind() == reflect.Slice, reflect.TypeOf(right[i]).Kind() == reflect.Slice
		if !leftSlice && !rightSlice {
			l := left[i].(float64)
			r := right[i].(float64)
			if l < r {
				return inOrder
			} else if l > r {
				return outOfOrder
			}
			continue
		}
		var ls []any
		if !leftSlice {
			ls = []any{left[i]}
		} else {
			ls = left[i].([]any)
		}
		var rs []any
		if !rightSlice {
			rs = []any{right[i]}
		} else {
			rs = right[i].([]any)
		}
		innerOrder := slicesOrder(ls, rs)
		if innerOrder != unknownOrder {
			return innerOrder
		}
	}
	if len(left) < len(right) {
		return inOrder
	}
	if len(right) < len(left) {
		return outOfOrder
	}
	return unknownOrder
}
