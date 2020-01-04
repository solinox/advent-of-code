package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
)

func main() {
	maze := parseInput("input.txt")

	// Part 1
	part1 := dijkstra(maze, false)
	fmt.Println("Part 1:", part1)

	// Part 2
	part2 := dijkstra(maze, true)
	fmt.Println("Part 2:", part2)
}

type Maze struct {
	Grid           map[Point]Tile
	PortalsByName  map[string][]Point
	PortalsByPoint map[Point]string
}

type Tile struct {
	IsPortal      bool
	IsInnerPortal bool
	Name          string
	Position      Point
}

type Point struct {
	X, Y, Depth int
}

// Pt returns a Point with a zero depth
func (pt Point) Pt() Point {
	return Point{pt.X, pt.Y, 0}
}

type State struct {
	Steps int
	Tile  Tile
}

func dijkstra(maze *Maze, isPart2 bool) int {
	tiles := make(map[Point]State)

	start := State{Steps: 0, Tile: maze.Grid[maze.PortalsByName["AA"][0].Pt()]}
	destPt := maze.PortalsByName["ZZ"][0].Pt()

	var walk func(current State, from Point)
	walk = func(current State, from Point) {
		// Some paths will go endlessly and cause an overflow
		// Set an upper limit to the amount of depth one can go
		// find limit through trial and error
		if current.Tile.Position.Depth > 50 {
			return
		}
		if current.Tile.Position.Y == destPt.Y && current.Tile.Position.X == destPt.X && current.Tile.Position.Depth == destPt.Depth {
			return
		}
		// check all directions
		possibleNeighbors := []Point{
			Point{current.Tile.Position.X, current.Tile.Position.Y + 1, current.Tile.Position.Depth}, // up
			Point{current.Tile.Position.X, current.Tile.Position.Y - 1, current.Tile.Position.Depth}, // down
			Point{current.Tile.Position.X - 1, current.Tile.Position.Y, current.Tile.Position.Depth}, // left
			Point{current.Tile.Position.X + 1, current.Tile.Position.Y, current.Tile.Position.Depth}, // right
		}
		if current.Tile.IsPortal {
			if !isPart2 || current.Tile.Position.Depth > 0 || (current.Tile.Position.Depth == 0 && current.Tile.IsInnerPortal) {
				for _, pt := range maze.PortalsByName[current.Tile.Name] {
					if pt.X != current.Tile.Position.X || pt.Y != current.Tile.Position.Y {
						pt.Depth = current.Tile.Position.Depth
						if isPart2 {
							ddepth := 1
							if !current.Tile.IsInnerPortal {
								ddepth = -1
							}
							pt.Depth += ddepth
						}
						possibleNeighbors = append(possibleNeighbors, pt)
					}
				}
			}
		}
		for _, neighbor := range possibleNeighbors {
			// if just came from this direction
			if from.X == neighbor.X && from.Y == neighbor.Y {
				continue
			}
			// if direction not a valid point in the maze
			if _, ok := maze.Grid[neighbor.Pt()]; !ok {
				continue
			}

			newSteps := current.Steps + 1
			newTile := maze.Grid[neighbor.Pt()]
			newTile.Position.Depth = neighbor.Depth
			newState := State{Steps: newSteps, Tile: newTile}
			// if we've already previously found a better path
			if prevState, ok := tiles[neighbor]; ok && prevState.Steps <= newSteps {
				continue
			}
			tiles[newState.Tile.Position] = newState
			walk(newState, current.Tile.Position)
		}
	}

	tiles[start.Tile.Position] = start
	walk(start, start.Tile.Position)
	return tiles[destPt].Steps
}

func copyMap(m map[string]bool) map[string]bool {
	o := make(map[string]bool)
	for k, v := range m {
		o[k] = v
	}
	return o
}

func parseInput(filename string) *Maze {
	data, _ := ioutil.ReadFile(filename)
	lines := bytes.Split(data, []byte{'\n'})
	maze := &Maze{Grid: make(map[Point]Tile), PortalsByName: make(map[string][]Point), PortalsByPoint: make(map[Point]string)}
	for y := 2; y < len(lines)-2; y++ {
		for x := 2; x < len(lines[y]); x++ {
			pt := Point{x, y, 0}
			b := lines[y][x]
			if b == '.' {
				tile := Tile{Position: pt}
				// check if next to a portal label
				up, down, left, right := Point{x, y - 1, 0}, Point{x, y + 1, 0}, Point{x - 1, y, 0}, Point{x + 1, y, 0}
				if c := lines[up.Y][up.X]; c >= 'A' && c <= 'Z' {
					tile.IsPortal = true
					tile.Name = string(lines[up.Y-1][up.X]) + string(c)
					tile.IsInnerPortal = y > 5
				}
				if c := lines[down.Y][down.X]; c >= 'A' && c <= 'Z' {
					tile.IsPortal = true
					tile.Name = string(c) + string(lines[down.Y+1][down.X])
					tile.IsInnerPortal = y < len(lines)-5
				}
				if c := lines[left.Y][left.X]; c >= 'A' && c <= 'Z' {
					tile.IsPortal = true
					tile.Name = string(lines[left.Y][left.X-1]) + string(c)
					tile.IsInnerPortal = x > 5
				}
				if c := lines[right.Y][right.X]; c >= 'A' && c <= 'Z' {
					tile.IsPortal = true
					tile.Name = string(c) + string(lines[right.Y][right.X+1])
					tile.IsInnerPortal = x < len(lines[y])-5
				}
				if tile.IsPortal {
					if pts, ok := maze.PortalsByName[tile.Name]; ok {
						maze.PortalsByName[tile.Name] = append(pts, pt)
					} else {
						maze.PortalsByName[tile.Name] = []Point{pt}
					}
					maze.PortalsByPoint[pt] = tile.Name
				}
				maze.Grid[pt] = tile
			}
		}
	}
	return maze
}
