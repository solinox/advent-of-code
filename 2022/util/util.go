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
		ret = append(ret, convert(l))
	}
	return ret
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
