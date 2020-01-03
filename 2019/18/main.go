package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"sync"
)

type Point struct {
	X, Y int
}
type Key struct {
	Position Point
	Letter   byte
	Door     *Door
	Acquired bool
}

func (k *Key) String() string {
	return string(k.Letter)
}

type Door struct {
	Position Point
	Letter   byte
	Key      *Key
	Open     bool
}
type Maze struct {
	Steps   int
	Grid    map[Point]byte
	Current []Point
	Keys    map[Point]*Key
	Doors   map[Point]*Door
}

func main() {
	// Part 1
	maze1 := parseInput("part1.txt")
	part1 := dijkstra(maze1)
	fmt.Println("Part 1:", part1)

	// Part 2
	maze2 := parseInput("part2.txt")
	part2 := dijkstra(maze2)
	fmt.Println("Part 2:", part2)
}

type State struct {
	Steps         int
	Position      []Point
	KeysCollected []*Key
}

func dijkstra(maze *Maze) int {
	var tiles sync.Map
	start := State{Steps: 0, Position: maze.Current, KeysCollected: make([]*Key, 0)}

	tiles.Store(hash(start.Position), start)
	minSteps := -1
	var mu sync.Mutex
	var walk func(current State, wg *sync.WaitGroup, depth int)
	walk = func(current State, wg *sync.WaitGroup, depth int) {
		defer wg.Done()
		// if already a worse option than a found solution, stop
		mu.Lock()
		if minSteps > 0 && current.Steps > minSteps {
			mu.Unlock()
			return
		}
		mu.Unlock()
		// find closest available keys
		allAvailableKeys := maze.ClosestKeys(current)
		for i, availableKeys := range allAvailableKeys {
			// fmt.Println(strings.Repeat("  ", depth), current, availableKeys)
			// if found all the keys, check if distance is minimum seen
			if len(availableKeys) == 0 && len(current.KeysCollected) == len(maze.Keys) {
				mu.Lock()
				// fmt.Println("Solved:", current.Steps, current.KeysCollected)
				if minSteps < 0 || current.Steps < minSteps {
					minSteps = current.Steps
				}
				mu.Unlock()
				return
			}
			var newWg sync.WaitGroup
			for key, dist := range availableKeys {
				newSteps := current.Steps + dist
				newPosition := make([]Point, len(current.Position))
				copy(newPosition, current.Position)
				newPosition[i] = key.Position
				newKeysCollected := make([]*Key, len(current.KeysCollected)+1)
				copy(newKeysCollected, current.KeysCollected)
				newKeysCollected[len(current.KeysCollected)] = key

				newState := State{Steps: newSteps, Position: newPosition, KeysCollected: newKeysCollected}
				// kill this branch if we've already previously found a better path
				if states, ok := tiles.Load(hash(newPosition)); ok {
					skip := false
					_states := states.([]State)
					for _, state := range _states {
						if state.Steps <= newSteps && key == state.KeysCollected[len(state.KeysCollected)-1] && containsAllKeys(newKeysCollected, state.KeysCollected) {
							// fmt.Println(strings.Repeat("  ", depth), "skipping", state, newState)
							skip = true
							break
						}
					}
					if skip {
						continue
					}
					_states = append(_states, newState)
					tiles.Store(hash(newPosition), _states)
				} else {
					tiles.Store(hash(newPosition), []State{newState})
				}
				newWg.Add(1)
				go walk(newState, &newWg, depth+1)
			}
			newWg.Wait()
		}
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go walk(start, &wg, 0)
	wg.Wait()
	return minSteps
}

func (maze *Maze) ClosestKeys(state State) []map[*Key]int {
	allKeys := make([]map[*Key]int, len(state.Position))
	tiles := make(map[Point]State)
	// Find all keys within accessible area
	// dijkstra again
	var walk func(current State, from Point, keys map[*Key]int)
	walk = func(current State, from Point, keys map[*Key]int) {
		// if on a key
		if key, ok := maze.Keys[current.Position[0]]; ok {
			// if we haven't collected the key
			collected := false
			for _, collectedKey := range state.KeysCollected {
				if key == collectedKey {
					collected = true
					break
				}
			}
			if !collected {
				keys[key] = current.Steps
				return
			}
		}
		// check all directions
		dirs := []Point{
			Point{current.Position[0].X, current.Position[0].Y + 1}, // up
			Point{current.Position[0].X, current.Position[0].Y - 1}, // down
			Point{current.Position[0].X - 1, current.Position[0].Y}, // left
			Point{current.Position[0].X + 1, current.Position[0].Y}, // right
		}
		for _, dir := range dirs {
			// if just came from this direction
			if from.X == dir.X && from.Y == dir.Y {
				continue
			}
			newSteps := current.Steps + 1
			// if direction not a valid point in the maze or if the direction is a wall
			if b, ok := maze.Grid[dir]; !ok || b == '#' {
				continue
				// if direction is a door that we do not have a key foor
			} else if b >= 'A' && b <= 'Z' && !haveKeyForDoor(b, state.KeysCollected) {
				continue
			}
			// if we've already previously found a better path
			if next, ok := tiles[dir]; ok && next.Steps <= newSteps {
				continue
			}
			newState := State{Steps: newSteps, Position: []Point{dir}}
			tiles[newState.Position[0]] = newState

			walk(newState, current.Position[0], keys)
		}
	}

	for i := range allKeys {
		start := State{Steps: 0, Position: state.Position[i : i+1]}
		tiles[start.Position[0]] = start
		keys := make(map[*Key]int)
		walk(start, start.Position[0], keys)
		allKeys[i] = keys
	}

	return allKeys
}

func haveKeyForDoor(doorLetter byte, keysCollected []*Key) bool {
	for _, key := range keysCollected {
		if key.Letter+'A'-'a' == doorLetter {
			return true
		}
	}
	return false
}

func containsAllKeys(keys, otherKeys []*Key) bool {
	for _, key := range keys {
		found := false
		for _, other := range otherKeys {
			if key == other {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func hash(pts []Point) string {
	str := ""
	for i := range pts {
		str += strconv.Itoa(pts[i].X) + strconv.Itoa(pts[i].Y)
	}
	return str
}

func parseInput(filename string) *Maze {
	file, _ := os.Open(filename)
	scanner := bufio.NewScanner(file)
	maze := &Maze{Grid: make(map[Point]byte), Keys: make(map[Point]*Key), Doors: make(map[Point]*Door), Current: make([]Point, 0)}
	for y := 0; scanner.Scan(); y++ {
		line := scanner.Text()
		for x := range line {
			pt := Point{x, y}
			b := line[x]
			maze.Grid[pt] = b
			if b == '#' {
			} else if b == '.' {
			} else if b == '@' {
				maze.Current = append(maze.Current, pt)
				continue
			} else if b >= 'a' && b <= 'z' {
				maze.Keys[pt] = &Key{Position: pt, Letter: b}
			} else if b >= 'A' && b <= 'Z' {
				maze.Doors[pt] = &Door{Position: pt, Letter: b, Open: false}
			}
		}
	}
	for _, key := range maze.Keys {
		doorLetter := key.Letter + 'A' - 'a'
		for _, door := range maze.Doors {
			if door.Letter == doorLetter {
				key.Door = door
				door.Key = key
				break
			}
		}
	}
	return maze
}
