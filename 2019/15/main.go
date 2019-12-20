package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var verbose = false

type Memory map[int]int

func main() {
	initialMemory := parseInput("input.txt")

	in, out := make(chan int), make(chan int)
	go runProgram(initialMemory, in, out)
	part1, part2 := make(chan int), make(chan int)
	// Part1 and Part2 are in dijkstra because the map exists there already
	go dijkstra(in, out, part1, part2)
	fmt.Println("Part 1:", <-part1)
	fmt.Println("Part 2:", <-part2)
}

type Point struct {
	X, Y int
}
type Tile int

const (
	Wall        Tile = 0
	Empty       Tile = 1
	Destination Tile = 2
)

type Vertex struct {
	Prev  Point
	Point Point
	Dist  int
}

func dijkstra(in, out, part1, part2 chan int) {
	tiles := make(map[Point]Vertex)
	walls := make(map[Point]bool)
	var dest Point

	// Part 1
	start := Vertex{Prev: Point{0, 0}, Point: Point{0, 0}, Dist: 0}
	tiles[start.Point] = start
	inputs := []int{
		1, // up
		2, // down
		3, // left
		4, // right
	}
	reverse := []int{
		2,
		1,
		4,
		3,
	}
	var walk func(current Vertex)
	walk = func(current Vertex) {
		// check all directions
		dirs := []Point{
			Point{current.Point.X, current.Point.Y + 1}, // up
			Point{current.Point.X, current.Point.Y - 1}, // down
			Point{current.Point.X - 1, current.Point.Y}, // left
			Point{current.Point.X + 1, current.Point.Y}, // right
		}
		for i, dir := range dirs {
			if walls[dir] {
				continue
			}
			v, ok := tiles[dir]
			if ok && v.Dist <= current.Dist {
				continue
			}
			// direction is unexplored. Let's try it and see!
			in <- inputs[i]
			result := Tile(<-out)
			if result == Wall {
				walls[dir] = true
				continue
			}
			if result == Destination {
				dest = dir
			}
			newV := Vertex{Prev: current.Point, Point: dir, Dist: current.Dist + 1}
			if newV.Dist < v.Dist || !ok {
				tiles[dir] = newV
			}
			// if we actually moved, check out the new tile
			walk(newV)
			// move back (we know it will be an empty tile, no need to see output)
			in <- reverse[i]
			<-out
		}
	}
	walk(start)
	part1 <- tiles[dest].Dist

	// Part 2
	allOxygen := map[Point]bool{dest: true}
	oxygenEdges := map[Point]bool{dest: true}
	minutes := 0
	for len(allOxygen) != len(tiles) {
		minutes++
		newEdges := make(map[Point]bool)
		for edge := range oxygenEdges {
			dirs := []Point{
				Point{edge.X, edge.Y + 1},
				Point{edge.X, edge.Y - 1},
				Point{edge.X - 1, edge.Y},
				Point{edge.X + 1, edge.Y},
			}
			for _, dir := range dirs {
				if allOxygen[dir] || walls[dir] {
					continue
				}
				if _, ok := tiles[dir]; ok {
					allOxygen[dir] = true
					newEdges[dir] = true
				}
			}
		}
		oxygenEdges = newEdges
	}
	part2 <- minutes
}

func runProgram(initialMemory Memory, input chan int, output chan<- int) {
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
