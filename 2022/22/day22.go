package main

import (
	_ "embed"
	"strings"

	"github.com/solinox/advent-of-code/2022/util"
)

//go:embed input.txt
var input string

type row struct {
	start, end int
	cells      []byte
}
type Board struct {
	rows, cols   []row
	instructions string
}

func main() {
	board := parse(input)
	util.RunTimed(part1, board)
	// util.RunTimed(part2, board)
}

func parse(in string) Board {
	sections := strings.Split(in, "\n\n")
	b := Board{instructions: sections[1]}
	parseRows := func(s string) []row {
		rows := make([]row, 0)
		start, end, cells, c, substance := 1, 1, make([]byte, 0), 0, false
		for _, v := range s {
			c++
			switch v {
			case '\n':
				rows = append(rows, row{start, end, cells})
				c = 0
				substance = false
				cells = make([]byte, 0)
				continue
			case ' ':
				substance = false
			default:
				if !substance {
					start = c
					substance = true
				}
				end = c
				cells = append(cells, byte(v))
			}
		}
		if len(cells) > 0 {
			rows = append(rows, row{start, end, cells})
		}
		return rows
	}
	transposed := func(s string) string {
		lines := strings.Split(s, "\n")
		var ns strings.Builder
		found := true
		for c := 0; found; c++ {
			found = false
			for _, l := range lines {
				if c < len(l) {
					found = true
					ns.WriteByte(l[c])
				} else {
					ns.WriteByte(' ')
				}
			}
			if found {
				ns.WriteByte('\n')
			}
		}
		return ns.String()
	}
	b.rows = parseRows(sections[0])
	b.cols = parseRows(transposed(sections[0]))
	return b
}

func part1(b Board) int {
	pY, pX, pD := 1, b.rows[0].start, 0
	for i := 0; i < len(b.instructions); {
		j := strings.IndexAny(b.instructions[i:], "RL")
		if j == -1 {
			j = len(b.instructions) - i
		}
		if j == 0 {
			dD := 1
			if b.instructions[i] == 'L' {
				dD = 3
			}
			pD = (pD + dD) % 4
			i++
		} else {
			n := util.ParseInt(b.instructions[i : i+j])
			pY, pX = b.Move(pY, pX, pD, n)
			i += j
		}
	}
	return 1000*pY + 4*pX + pD
}

func (b Board) Move(y, x, d, n int) (int, int) {
	dy, dx := 0, 0
	switch d {
	case 0:
		dx++
	case 1:
		dy++
	case 2:
		dx--
	case 3:
		dy--
	}
	if dx != 0 {
		return y, b.moveX(y, x, dx, n)
	}
	return b.moveY(y, x, dy, n), x
}

func (b Board) moveX(y, x, dx, n int) int {
	for m := 0; m < n; m++ {
		nx := x + dx
		if nx < b.rows[y-1].start {
			nx = b.rows[y-1].end
		} else if nx > b.rows[y-1].end {
			nx = b.rows[y-1].start
		}
		if b.rows[y-1].cells[nx-b.rows[y-1].start] == '#' {
			break
		}
		x = nx
	}
	return x
}

func (b Board) moveY(y, x, dy, n int) int {
	for m := 0; m < n; m++ {
		ny := y + dy
		if ny < b.cols[x-1].start {
			ny = b.cols[x-1].end
		} else if ny > b.cols[x-1].end {
			ny = b.cols[x-1].start
		}
		if b.cols[x-1].cells[ny-b.cols[x-1].start] == '#' {
			break
		}
		y = ny
	}
	return y
}
