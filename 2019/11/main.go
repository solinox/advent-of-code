package main

import (
	"strconv"
	"io/ioutil"
	"strings"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

var verbose = false

type Memory map[int]int

type Color int
const (
	black Color = 0
	white = 1
)

type Point struct {
	X, Y int
}

type Direction int
const (
	up Direction = iota
	left
	down
	right
)

type Position struct {
	Point
	Direction
}

func main() {
	memory := parseInput("input.txt")

	// Part 1
	part1 := paintPanels(memory, black, up)
	fmt.Println("Part 1", len(part1))

	// Part 2
	part2 := paintPanels(memory, white, up)
	part2File := "part2.png"
	part2String := createImage(part2, part2File)

	// The outputs are mirrored off the horizontal axis
	// Since image coordinates assume the minimum values are the top/leftmost pixels
	fmt.Println("Part 2 image saved to", part2File)
	fmt.Println(part2String)
}

func createImage(colors map[Point]Color, filename string) string {
	minX, maxX, minY, maxY := getBounds(colors)
	fmt.Println(minX, maxX, minY, maxY)
	m := image.NewRGBA(image.Rect(minX, minY, maxX+1, maxY+1))
	file, _ := os.Create(filename)

	imgBytes := make([][]byte, maxY - minY + 1)
	for y := minY; y <= maxY; y++ {
		imgBytes[y-minY] = make([]byte, maxX - minX + 1)
		for x := minX; x <= maxX; x++ {
			imgBytes[y-minY][x-minX] = ' '
		}
	}

	for pt, colour := range colors {
		if colour == black {
			m.Set(pt.X, pt.Y, color.Black)
			imgBytes[pt.Y-minY][pt.X-minX] = ' '
		} else if colour == white {
			m.Set(pt.X, pt.Y, color.White)
			imgBytes[pt.Y-minY][pt.X-minX] = '#'
		}
	}

	png.Encode(file, m)
	imgString := ""
	for i := range imgBytes {
		imgString += string(imgBytes[i]) + "\n"
	}

	return imgString
}

func getBounds(colors map[Point]Color) (int, int, int, int) {
	minX, maxX, minY, maxY := 0, 0, 0, 0
	for pt := range colors {
		if pt.X < minX {
			minX = pt.X
		} else if pt.X > maxX {
			maxX = pt.X
		}
		if pt.Y < minY {
			minY = pt.Y
		} else if pt.Y > maxY {
			maxY = pt.Y
		}
	}
	return minX, maxX, minY, maxY
}

func paintPanels(initialMemory Memory, initialColor Color, initialDirection Direction) map[Point]Color {
	in, out := make(chan int, 2), make(chan int)
	position := Position{Point: Point{0, 0}, Direction: up}
	paintedPanels := make(map[Point]Color)
	go runProgram(initialMemory, in, out)
	in <- int(initialColor)
	outputs := make([]int, 0, 2)
	for output := range out {
		outputs = append(outputs, output)
		if len(outputs) == 2 {
			color, turn := Color(outputs[0]), outputs[1]
			// paint panel
			paintedPanels[position.Point] = color

			// turn
			switch turn {
			case 0: // left 90 degrees
				position.Direction = (position.Direction + 1) % 4
			case 1: // right 90 degrees
				position.Direction = (position.Direction + 3) % 4
			}

			// move 1 forward
			switch position.Direction {
			case up:
				position.Point = Point{position.Point.X, position.Point.Y+1}
			case down:
				position.Point = Point{position.Point.X, position.Point.Y-1}
			case right:
				position.Point = Point{position.Point.X+1, position.Point.Y}
			case left:
				position.Point = Point{position.Point.X-1, position.Point.Y}
			}
			outputs = make([]int, 0, 2)
			in <- int(paintedPanels[position.Point])
		}
	}
	return paintedPanels
}

func runProgram(initialMemory Memory, input <-chan int, output chan<- int) {
	memory := make(Memory)
	for k, v := range initialMemory {
		memory[k] = v
	}
	relBase := 0

	for i := 0;; {
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
			panic(fmt.Sprintf("Op %d is not supported", op))
		}

		if verbose {
			fmt.Println(msg)
		}
	}
	fmt.Println("Program did not halt but went out of bounds unexpectedly")
	close(output)
	return
}

func getParams(memory Memory, n, pointer, relBase, numParams int) []int {
	params := make([]int, numParams)
	for i := range params {
		value := pointer+1+i
		mode := n % 10
		if mode == 0 {
			params[i] = memory[value]
		} else if mode == 1 {
			params[i] = value
		} else if mode == 2 {
			params[i] = relBase+memory[value]
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
		ints[i], _ = strconv.Atoi(fields[i])
	}
	return ints
}
