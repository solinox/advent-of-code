package main

import (
	"bufio"
	"fmt"
	"os"
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
	Current Point
	Keys    map[Point]*Key
	Doors   map[Point]*Door
}

func main() {
	maze := parseInput("input.txt")

	// Part 1
	part1 := dijkstra(maze)
	fmt.Println("Part 1:", part1)
}

type State struct {
	Steps         int
	Position      Point
	KeysCollected []*Key
}

func dijkstra(maze *Maze) int {
	var tiles sync.Map
	start := State{Steps: 0, Position: maze.Current, KeysCollected: make([]*Key, 0)}

	tiles.Store(start.Position, start)
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
		availableKeys := maze.ClosestKeys(current)
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
		// time.Sleep(1 * time.Second)
		for key, dist := range availableKeys {
			newSteps := current.Steps + dist
			newKeysCollected := make([]*Key, len(current.KeysCollected)+1)
			copy(newKeysCollected, current.KeysCollected)
			newKeysCollected[len(current.KeysCollected)] = key
			newState := State{Steps: newSteps, Position: key.Position, KeysCollected: newKeysCollected}
			// kill this branch if we've already previously found a better path
			if states, ok := tiles.Load(key.Position); ok {
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
				tiles.Store(key.Position, _states)
			} else {
				tiles.Store(key.Position, []State{newState})
			}
			newWg.Add(1)
			go walk(newState, &newWg, depth+1)
		}
		newWg.Wait()
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go walk(start, &wg, 0)
	wg.Wait()
	return minSteps
}

func (maze *Maze) ClosestKeys(state State) map[*Key]int {
	keys := make(map[*Key]int)
	// Find all keys within accessible area
	// dijkstra again
	tiles := make(map[Point]State)
	start := State{Steps: 0, Position: state.Position}
	tiles[start.Position] = start
	var walk func(current State, from Point)
	walk = func(current State, from Point) {
		// if on a key
		if key, ok := maze.Keys[current.Position]; ok {
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
			Point{current.Position.X, current.Position.Y + 1}, // up
			Point{current.Position.X, current.Position.Y - 1}, // down
			Point{current.Position.X - 1, current.Position.Y}, // left
			Point{current.Position.X + 1, current.Position.Y}, // right
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
			newState := State{Steps: newSteps, Position: dir}
			tiles[newState.Position] = newState

			walk(newState, current.Position)
		}
	}
	walk(start, start.Position)
	return keys
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

func parseInput(filename string) *Maze {
	file, _ := os.Open(filename)
	scanner := bufio.NewScanner(file)
	maze := &Maze{Grid: make(map[Point]byte), Keys: make(map[Point]*Key), Doors: make(map[Point]*Door)}
	for y := 0; scanner.Scan(); y++ {
		line := scanner.Text()
		for x := range line {
			pt := Point{x, y}
			b := line[x]
			maze.Grid[pt] = b
			if b == '#' {
			} else if b == '.' {
			} else if b == '@' {
				maze.Current = pt
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
