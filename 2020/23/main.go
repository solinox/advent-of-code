package main

import (
	"container/ring"
	"log"
	"os"
	"time"

	"github.com/solinox/advent-of-code/2020/pkg/input"
)

func main() {
	cups := input.SplitInt(os.Stdin, "")

	t0 := time.Now()
	log.Println(part1(cups), time.Since(t0))

	t0 = time.Now()
	log.Println(part2(cups), time.Since(t0))
}

func part1(cups []int) int {
	lenCups := len(cups)
	r := ring.New(lenCups)
	for _, cup := range cups {
		r.Value = cup
		r = r.Next()
	}

	r, rMap := move(r, 100)

	label := 0
	r = rMap[1]
	r.Do(func(x interface{}) {
		v := x.(int)
		if v == 1 {
			return
		}
		label = label*10 + v
	})

	return label
}

func part2(cups []int) int {
	lenCups := 1000000
	r := ring.New(lenCups)
	for _, cup := range cups {
		r.Value = cup
		r = r.Next()
	}
	for i := len(cups); i < lenCups; i++ {
		cup := i + 1
		r.Value = cup
		r = r.Next()
	}

	r, rMap := move(r, 10000000)

	star1 := rMap[1].Next()
	star2 := star1.Next()

	return star1.Value.(int) * star2.Value.(int)
}

func move(r *ring.Ring, n int) (*ring.Ring, []*ring.Ring) {
	lenCups := r.Len()
	rMap := make([]*ring.Ring, lenCups+1)
	for i := 0; i < lenCups; i++ {
		rMap[r.Value.(int)] = r
		r = r.Next()
	}

	for i := 0; i < n; i++ {
		current := r
		cup := current.Value.(int)
		u := r.Unlink(3)
		dest := find(r, rMap, u, lenCups, cup-1)
		dest.Link(u)
		r = current.Next()
	}
	return r, rMap
}

func find(r *ring.Ring, rMap []*ring.Ring, removed *ring.Ring, lenCups, destCup int) *ring.Ring {
	if destCup <= 0 {
		destCup = lenCups
	}
	dest := rMap[destCup]
	if contains(removed, dest) {
		return find(r, rMap, removed, lenCups, destCup-1)
	}
	return dest
}

func contains(r *ring.Ring, v *ring.Ring) bool {
	if r == nil || v == nil {
		return false
	}
	if r == v {
		return true
	}
	for p := r.Next(); p != r; p = p.Next() {
		if p == v {
			return true
		}
	}
	return false
}
