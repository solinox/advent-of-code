// util is a small, miscellaneous package for advent-of-code
// includes common functions used by multiple days for manipulating input, printing output, etc
package util

import (
	"bufio"
	"fmt"
	"io"
	"reflect"
	"runtime"
	"strconv"
	"time"

	"golang.org/x/exp/constraints"
)

func RunTimed[T, U any](fn func(T) U, in T) {
	t0 := time.Now()
	out := fn(in)
	dt := time.Since(t0)
	name := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	fmt.Printf("%v\n%s in %s\n", out, name, dt)
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
