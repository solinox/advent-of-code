// util is a small, miscellaneous package for advent-of-code
// includes common functions used by multiple days for manipulating input, printing output, etc
package util

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/constraints"
)

func RunTimed[T, U any](fn func(T) U, in T) {
	t0 := time.Now()
	out := fn(in)
	dt := time.Since(t0)
	fmt.Printf("%v\n  => in %s\n", out, dt)
}

func ParseLines[T any](r io.Reader, convert func(line string) T) []T {
	s := bufio.NewScanner(r)
	ret := make([]T, 0)
	for s.Scan() {
		l := s.Text()
		if l == "" {
			continue
		}
		ret = append(ret, convert(l))
	}
	return ret
}

func ParseLinesReduce[T any](r io.Reader, fn func(agg T, line string) T, val T) T {
	s := bufio.NewScanner(r)
	for s.Scan() {
		l := s.Text()
		if l == "" {
			continue
		}
		val = fn(val, l)
	}
	return val
}

func ParseSections[T any](r io.Reader, convert func(line string) T) [][]T {
	s := bufio.NewScanner(r)
	ret := make([][]T, 0)
	section := make([]T, 0)
	for s.Scan() {
		l := s.Text()
		if l == "" {
			ret = append(ret, section)
			section = make([]T, 0)
			continue
		}
		section = append(section, convert(l))
	}
	if len(section) > 0 {
		ret = append(ret, section)
	}
	return ret
}

func ParseDelimited[T any](line, delim string, convert func(s string) T) []T {
	fields := strings.Split(line, delim)
	return SliceFrom(fields, convert)
}

func ParseInt(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}

func SliceFrom[T, U any](in []U, fn func(u U) T) []T {
	out := make([]T, 0, len(in))
	for _, v := range in {
		out = append(out, fn(v))
	}
	return out
}

func Max[T constraints.Ordered](n ...T) T {
	var max T
	if len(n) == 0 {
		return max
	}
	max = n[0]
	for _, m := range n[1:] {
		if m > max {
			max = m
		}
	}
	return max
}

// MaxN gets the n max numbers, running in O(n*m) instead of O(n*logn) since it does not sort
// may or may not be faster than sorting, depending on input n value
func MaxN[T constraints.Ordered](n int, nn ...T) []T {
	if n <= 0 || len(nn) == 0 {
		return nil
	}
	max := make([]T, n)
	max[0] = nn[0]
	for _, m := range nn[1:] {
		for i, mm := range max {
			if m > mm {
				copy(max[i+1:], max[i:])
				max[i] = m
				break
			}
		}
	}
	return max
}

func Min[T constraints.Ordered](n ...T) T {
	var min T
	if len(n) == 0 {
		return min
	}
	min = n[0]
	for _, m := range n[1:] {
		if m < min {
			min = m
		}
	}
	return min
}

func Sum[T constraints.Ordered](n ...T) T {
	var sum T
	for _, m := range n {
		sum += m
	}
	return sum
}

func Reverse[T any](s []T) []T {
	for i, j := 0, len(s)-1; i < len(s)/2; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

type Number interface {
	constraints.Signed | constraints.Float
}

func Abs[T Number](n T) T {
	if n < 0 {
		return -n
	}
	return n
}

type Vector struct {
	Y, X int
}

func (p Vector) Add(o Vector) Vector {
	return Vector{Y: p.Y + o.Y, X: p.X + o.X}
}

func (p Vector) Sub(o Vector) Vector {
	return Vector{Y: p.Y - o.Y, X: p.X - o.X}
}

func (p Vector) Unit() Vector {
	v := Vector{}
	if p.Y != 0 {
		v.Y = p.Y / Abs(p.Y)
	}
	if p.X != 0 {
		v.X = p.X / Abs(p.X)
	}
	return v
}

func (p Vector) Dist(o Vector) int {
	return Abs(p.Y-o.Y) + Abs(p.X-o.X)
}

type Vector3 struct {
	Z, Y, X int
}

func (p Vector3) Add(o Vector3) Vector3 {
	return Vector3{Z: p.Z + o.Z, Y: p.Y + o.Y, X: p.X + o.X}
}

type Range struct {
	Min, Max int
}

func (r Range) Dist() int {
	d := Abs(r.Max-r.Min) + 1
	if r.Max > 0 && r.Min < 0 {
		d--
	}
	return d
}

func MergeRanges(rs []Range) []Range {
	anyMerged := true
	for anyMerged {
		anyMerged = false
		for i := len(rs) - 1; i >= 1; i-- {
			for j := i - 1; j >= 0; j-- {
				if rs[i].Max < rs[j].Min || rs[i].Min > rs[i].Max || rs[j].Max < rs[i].Min || rs[j].Min > rs[i].Max {
					continue
				}
				anyMerged = true
				rs[j] = Range{Min: Min(rs[i].Min, rs[j].Min), Max: Max(rs[i].Max, rs[j].Max)}
				if i < len(rs)-1 {
					rs = append(rs[:i], rs[i+1:]...)
				} else {
					rs = rs[:i]
				}
				break
			}
		}
	}
	return rs
}

// All returns all combinations for a given string array.
// This is essentially a powerset of the given set except that the empty set is disregarded.
// ABCDE combinations
// A, AB, AC, AD, AE, ABC, ABD, ABE, ACD, ACE, ADE, ABCD, ACDE, ABCDE
// B, BC, BD, BE, BCD, BCE, BCDE
// C, CD, CE, CDE
// D, DE
// E
func AllCombinations[T any](set []T) (subsets [][]T) {
	length := uint(len(set))

	// Go through all possible combinations of objects
	// from 1 (only first object in subset) to 2^length (all objects in subset)
	for subsetBits := 1; subsetBits < (1 << length); subsetBits++ {
		var subset []T

		for object := uint(0); object < length; object++ {
			// checks if object is contained in subset
			// by checking if bit 'object' is set in subsetBits
			if (subsetBits>>object)&1 == 1 {
				// add object to subset
				subset = append(subset, set[object])
			}
		}
		// add subset to subsets
		subsets = append(subsets, subset)
	}
	return subsets
}

func Contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}
