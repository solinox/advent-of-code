package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"math"
	"sort"
)

type Angle1 struct {
	Quadrant int
	Val      float64
}

type Angle2 struct {
	Asteroid Point
	Quadrant int
	Val      float64
	Mag      float64
}

type Point struct {
	X, Y int
}

type Asteroids map[Point]bool

func main() {
	asteroids := parseInput("input.txt")
	
	// Part 1
	lineOfSight := getAsteroidsInLOS(asteroids)
	part1 := 0
	var bestAsteroid Point
	for asteroid, count := range lineOfSight {
		if count > part1 {
			part1 = count
			bestAsteroid = asteroid
		}
	}
	fmt.Println("Part 1", part1, bestAsteroid)

	// Part 2
	part2 := sortAsteroids(asteroids, bestAsteroid)
	fmt.Println("Part 2", part2[199].Asteroid.X * 100 + part2[199].Asteroid.Y)
}

func sortAsteroids(asteroids Asteroids, asteroid Point) []Angle2{
	sorted := make([]Angle2, 0, len(asteroids)-1)
	for other := range asteroids {
		if asteroid.X == other.X && asteroid.Y == other.Y {
			continue
		}
		_, angle := newAngles(asteroid, other)
		sorted = append(sorted, angle)
	}

	sort.Slice(sorted, func(i, j int) bool {
		a1, a2 := sorted[i], sorted[j]
		if a1.Quadrant == a2.Quadrant {
			if a1.Val == a2.Val {
				return a1.Mag < a2.Mag
			}
			if a1.Quadrant == 1 || a1.Quadrant == 3 {
				return a1.Val > a2.Val
			}
			return a1.Val < a2.Val
		}
		return a1.Quadrant < a2.Quadrant
	})

	for i := 1; i < len(sorted); {
		if sorted[i].Val == sorted[i-1].Val {
			// move to end, preserve order of the rest
			cutAngle := sorted[i]
			sorted = append(sorted[:i], append(sorted[i+1:], cutAngle)...)
		} else {
			i++
		}
	}
	return sorted
}

func getAsteroidsInLOS(asteroids Asteroids) map[Point]int {
	lineOfSight := make(map[Point]int)
	for asteroid := range asteroids {
		blockedAngles := make(map[Angle1]bool)
		lineOfSightCount := 0
		for other := range asteroids {
			if asteroid.X == other.X && asteroid.Y == other.Y {
				continue
			}
			angle, _ := newAngles(asteroid, other)
			if exists := blockedAngles[angle]; exists {
				continue
			} else {
				lineOfSightCount++
				blockedAngles[angle] = true
			}
		}

		lineOfSight[asteroid] = lineOfSightCount
	}
	return lineOfSight
}

// Lazy and tweaked this to return the type I need for either part 1 or 2
func newAngles(asteroid, other Point) (Angle1, Angle2) {
	angleVal := math.Abs(float64(other.Y - asteroid.Y) / float64(other.X - asteroid.X))
	quadrant := 0
	if other.X >= asteroid.X && other.Y <= asteroid.Y {
		quadrant = 1 // Northeast
	} else if other.X >= asteroid.X && other.Y >= asteroid.Y {
		quadrant = 2 // Southeast
	} else if other.X < asteroid.X && other.Y >= asteroid.Y {
		quadrant = 3 // Southwest
	} else {
		quadrant = 4 // Northwest
	}
	return Angle1{Quadrant: quadrant, Val: angleVal}, Angle2{Quadrant: quadrant, Val: angleVal, Asteroid: other, Mag: math.Abs(float64(other.X - asteroid.X)) + math.Abs(float64(other.Y - asteroid.Y))}
}

func parseInput(filename string) Asteroids {
	data, _ := ioutil.ReadFile(filename)
	lines := strings.Split(string(data), "\n")
	asteroids := make(Asteroids)
	for y := range lines {
		for x := range lines[y] {
			if lines[y][x] == '#' {
				pt := Point{x,y}
				asteroids[pt] = true
			}
		}
	}
	return asteroids
}
