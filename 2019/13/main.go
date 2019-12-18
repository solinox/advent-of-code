package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

type Memory map[int]int

type Arcade struct {
	grid     map[Point]int
	data     [][]int
	ball     Point
	joystick Point
	score    int
}
type Point struct {
	X, Y int
}

func (a *Arcade) Set(tile []int) *Arcade {
	if tile[0] == -1 && tile[1] == 0 {
		a.score = tile[2]
		return a
	}
	pt := Point{tile[0], tile[1]}
	a.grid[pt] = tile[2]
	a.data[pt.Y][pt.X] = tile[2]
	if tile[2] == 3 {
		a.joystick = pt
	}
	if tile[2] == 4 {
		a.ball = pt
	}
	return a
}
func (a *Arcade) String() string {
	str := "Score: " + strconv.Itoa(a.score) + "\n"
	for y := range a.data {
		row := make([]byte, len(a.data[y]))
		for x := range a.data[y] {
			switch a.data[y][x] {
			case 0:
				row[x] = ' '
			case 1:
				row[x] = 'W'
			case 2:
				row[x] = 'B'
			case 3:
				row[x] = '_'
			case 4:
				row[x] = 'o'
			default:
				row[x] = ' '
			}
		}
		str += string(row) + "\n"
	}
	return str
}

func NewArcade(outputs []int) *Arcade {
	arcade := &Arcade{grid: make(map[Point]int)}
	for i := 0; i < len(outputs); i += 3 {
		pt := Point{outputs[i], outputs[i+1]}
		arcade.grid[pt] = outputs[i+2]
	}
	minX, maxX, minY, maxY := 0, 0, 0, 0
	for pt := range arcade.grid {
		if pt.X < minX {
			minX = pt.X
		}
		if pt.X > maxX {
			maxX = pt.X
		}
		if pt.Y < minY {
			minY = pt.Y
		}
		if pt.Y > maxY {
			maxY = pt.Y
		}
	}
	arcade.data = make([][]int, maxY-minY+1)
	for i := range arcade.data {
		arcade.data[i] = make([]int, maxX-minX+1)
	}
	for pt, v := range arcade.grid {
		arcade = arcade.Set([]int{pt.X, pt.Y, v})
	}
	return arcade
}

var verbose = false

func main() {
	initialMemory := parseInput("input.txt")

	// Part 1
	in, out := make(chan int), make(chan int)
	go runProgram(initialMemory, in, out, nil)
	outputs := make([]int, 0)
	for output := range out {
		outputs = append(outputs, output)
	}
	part1 := 0
	for i := 2; i < len(outputs); i += 3 {
		if outputs[i] == 2 {
			part1++
		}
	}
	fmt.Println("Part 1", part1)

	// Part 2
	in2, out2 := make(chan int), make(chan int)
	initialMemory[0] = 2
	arcade := NewArcade(outputs)
	go runProgram(initialMemory, in2, out2, arcade)
	outputs2 := make([]int, 0, 3)
	for output := range out2 {
		outputs2 = append(outputs2, output)
		if len(outputs2) == 3 {
			arcade = arcade.Set(outputs2)
			outputs2 = make([]int, 0, 3)
		}
	}
	fmt.Println("Part 2", arcade.score)
}

func move(arcade *Arcade, input chan int) {
	if arcade.ball.X < arcade.joystick.X {
		input <- -1
	} else if arcade.ball.X > arcade.joystick.X {
		input <- 1
	} else {
		input <- 0
	}
}

func runProgram(initialMemory Memory, input chan int, output chan<- int, arcade *Arcade) {
	memory := make(Memory)
	for k, v := range initialMemory {
		memory[k] = v
	}
	relBase := 0

	for i := 0; ; {
		n := memory[i]
		op := n % 100
		n /= 100

		msg := fmt.Sprintf("i=%d\tOp=%d\tRel=%d\t", i, memory[i], relBase)

		switch op {
		case 1:
			p := getParams(memory, n, i, relBase, 3)
			memory[p[2]] = memory[p[0]] + memory[p[1]]
			msg += fmt.Sprintf("memory[%d] = memory[%d] + memory[%d] ==> %d", p[2], p[0], p[1], memory[p[2]])
			i += 4
		case 2:
			p := getParams(memory, n, i, relBase, 3)
			memory[p[2]] = memory[p[0]] * memory[p[1]]
			msg += fmt.Sprintf("memory[%d] = memory[%d] * memory[%d] ==> %d", p[2], p[0], p[1], memory[p[2]])
			i += 4
		case 3:
			p := getParams(memory, n, i, relBase, 1)
			fmt.Println(arcade)
			// There's a nasty race condition where the position of the ball isn't updated before the joystick tries to move
			// Which causes the joystick to end up missing the ball and ending the game early
			// Add a sleep here to "fix" the race condition
			// Has a bonus of creating a nice visual
			// May need to adjust sleep time for your computer to prevent the race condition
			time.Sleep(10 * time.Millisecond)
			go move(arcade, input)
			memory[p[0]] = <-input
			msg += fmt.Sprintf("memory[%d] = input ==> %d", p[0], memory[p[0]])
			i += 2
		case 4:
			p := getParams(memory, n, i, relBase, 1)
			output <- memory[p[0]]
			msg += fmt.Sprintf("output = memory[%d] ==> %d", p[0], memory[p[0]])
			i += 2
		case 5:
			p := getParams(memory, n, i, relBase, 2)
			msg += fmt.Sprintf("memory[%d] != 0 ==> %d != 0. ", p[0], memory[p[0]])
			if memory[p[0]] != 0 {
				i = memory[p[1]]
				msg += fmt.Sprintf("Set i to memory[%d] = %d", p[1], memory[p[1]])
			} else {
				i += 3
				msg += "do nothing but increment pointer"
			}
		case 6:
			p := getParams(memory, n, i, relBase, 2)
			msg += fmt.Sprintf("memory[%d] == 0 ==> %d == 0. ", p[0], memory[p[0]])
			if memory[p[0]] == 0 {
				i = memory[p[1]]
				msg += fmt.Sprintf("Set i to memory[%d] = %d", p[1], memory[p[1]])
			} else {
				i += 3
				msg += "do nothing but increment pointer"
			}
		case 7:
			p := getParams(memory, n, i, relBase, 3)
			msg += fmt.Sprintf("memory[%d] < memory[%d] ==> %d < %d. ", p[0], p[1], memory[p[0]], memory[p[1]])
			if memory[p[0]] < memory[p[1]] {
				memory[p[2]] = 1
				msg += fmt.Sprintf("Set memory[%d] to 1", p[2])
			} else {
				memory[p[2]] = 0
				msg += fmt.Sprintf("Set memory[%d] to 0", p[2])
			}
			i += 4
		case 8:
			p := getParams(memory, n, i, relBase, 3)
			msg += fmt.Sprintf("memory[%d] == memory[%d] ==> %d == %d. ", p[0], p[1], memory[p[0]], memory[p[1]])
			if memory[p[0]] == memory[p[1]] {
				memory[p[2]] = 1
				msg += fmt.Sprintf("Set memory[%d] to 1", p[2])
			} else {
				memory[p[2]] = 0
				msg += fmt.Sprintf("Set memory[%d] to 0", p[2])
			}
			i += 4
		case 9:
			p := getParams(memory, n, i, relBase, 1)
			msg += fmt.Sprintf("Add memory[%d] = %d to old relBase = %d.", p[0], memory[p[0]], relBase)
			relBase += memory[p[0]]
			i += 2
		case 99:
			close(output)
			return
		default:
			close(output)
			panic(fmt.Sprintf("Op %d is not supported", op))
		}

		if verbose {
			fmt.Println(msg)
		}
	}

}

func getParams(memory Memory, n, pointer, relBase, numParams int) []int {
	params := make([]int, numParams)
	for i := range params {
		value := pointer + 1 + i
		mode := n % 10
		if mode == 0 {
			params[i] = memory[value]
		} else if mode == 1 {
			params[i] = value
		} else if mode == 2 {
			params[i] = relBase + memory[value]
		}
		n /= 10
	}
	return params
}

func parseInput(filename string) Memory {
	data, _ := ioutil.ReadFile(filename)
	fields := strings.Split(string(data), ",")
	ints := make(Memory)
	for i := range fields {
		ints[i], _ = strconv.Atoi(strings.TrimSpace(fields[i]))
	}
	return ints
}
