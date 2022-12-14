package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/solinox/advent-of-code/2022/util"
)

//go:embed input.txt
var input string

func main() {
	packets := util.ParseLines(strings.NewReader(input), newPacket)
	util.RunTimed(part1, packets)
	util.RunTimed(part2, packets)
}

func part1(packets []any) int {
	sum := 0
	for i := 0; i < len(packets)/2; i++ {
		if packetCmp(packets[i*2], packets[i*2+1]) <= 0 {
			sum += i + 1
		}
	}
	return sum
}

func part2(packets []any) int {
	p1, p2 := "[[2]]", "[[6]]"
	packets = append(packets, newPacket(p1), newPacket(p2))

	sort.Slice(packets, func(i, j int) bool { return packetCmp(packets[i], packets[j]) < 0 })

	product := 1
	for i := range packets {
		if s := fmt.Sprint(packets[i]); s == p1 || s == p2 {
			product *= (i + 1)
		}
	}
	return product
}

func newPacket(s string) any {
	var p any
	json.Unmarshal([]byte(s), &p)
	return p
}

// negative if left < right, positive if left > right, 0 if equal
func packetCmp(left, right any) int {
	ls, lok := left.([]any)
	rs, rok := right.([]any)

	if !lok && !rok {
		return int(left.(float64) - right.(float64))
	}
	if !lok {
		ls = []any{left}
	}
	if !rok {
		rs = []any{right}
	}
	for i := 0; i < len(ls) && i < len(rs); i++ {
		if cmp := packetCmp(ls[i], rs[i]); cmp != 0 {
			return cmp
		}
	}
	return len(ls) - len(rs)
}
