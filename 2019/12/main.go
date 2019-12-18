package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type Vector struct {
	X, Y, Z int
}

func (v Vector) Energy() int {
	return abs(v.X) + abs(v.Y) + abs(v.Z)
}

func (v Vector) Add(other Vector) Vector {
	return Vector{v.X + other.X, v.Y + other.Y, v.Z + other.Z}
}

type Moon struct {
	Position Vector
	Velocity Vector
}

func (m Moon) Energy() int {
	return m.Position.Energy() * m.Velocity.Energy()
}

func main() {
	moons1 := parseInput("input.txt")
	moons2 := make([]Moon, len(moons1))
	copy(moons2, moons1)

	// Part 1
	part1Moons := step(moons1, 1000)
	part1 := 0
	for i := range part1Moons {
		part1 += part1Moons[i].Energy()
	}
	fmt.Println("Part 1", part1)

	// Part 2
	xMap := make(map[[8]int]int)
	yMap := make(map[[8]int]int)
	zMap := make(map[[8]int]int)
	part2x, part2y, part2z := 0, 0, 0
	for i := 0; part2x == 0 || part2y == 0 || part2z == 0; i++ {
		xs := [8]int{moons2[0].Position.X, moons2[1].Position.X, moons2[2].Position.X, moons2[3].Position.X, moons2[0].Velocity.X, moons2[1].Velocity.X, moons2[2].Velocity.X, moons2[3].Velocity.X}
		ys := [8]int{moons2[0].Position.Y, moons2[1].Position.Y, moons2[2].Position.Y, moons2[3].Position.Y, moons2[0].Velocity.Y, moons2[1].Velocity.Y, moons2[2].Velocity.Y, moons2[3].Velocity.Y}
		zs := [8]int{moons2[0].Position.Z, moons2[1].Position.Z, moons2[2].Position.Z, moons2[3].Position.Z, moons2[0].Velocity.Z, moons2[1].Velocity.Z, moons2[2].Velocity.Z, moons2[3].Velocity.Z}
		if x, ok := xMap[xs]; ok && part2x == 0 {
			part2x = i - x
			xMap = nil
		} else if xMap != nil {
			xMap[xs] = i
		}
		if y, ok := yMap[ys]; ok && part2y == 0 {
			part2y = i - y
			yMap = nil
		} else if yMap != nil {
			yMap[ys] = i
		}
		if z, ok := zMap[zs]; ok && part2z == 0 {
			part2z = i - z
			zMap = nil
		} else if zMap != nil {
			zMap[zs] = i
		}
		moons2 = step(moons2, 1)
	}
	part2 := lcm(lcm(part2x, part2y), part2z)
	fmt.Println("Part 2", part2)
}

func lcm(a, b int) int {
	return abs(a*b) / gcd(a, b)
}

func gcd(a, b int) int {
	if a == b {
		return a
	} else if a > b {
		return gcd(a-b, b)
	}
	return gcd(a, b-a)
}

func step(moons []Moon, steps int) []Moon {
	for step := 0; step < steps; step++ {
		// First apply gravity
		for i := range moons {
			for j := range moons {
				if i <= j {
					continue
				}
				dx1, dx2 := gravity(moons[i].Position.X, moons[j].Position.X)
				moons[i].Velocity.X += dx1
				moons[j].Velocity.X += dx2

				dy1, dy2 := gravity(moons[i].Position.Y, moons[j].Position.Y)
				moons[i].Velocity.Y += dy1
				moons[j].Velocity.Y += dy2

				dz1, dz2 := gravity(moons[i].Position.Z, moons[j].Position.Z)
				moons[i].Velocity.Z += dz1
				moons[j].Velocity.Z += dz2
			}
		}

		// Then apply velocity
		for i := range moons {
			moons[i].Position = moons[i].Position.Add(moons[i].Velocity)
		}
	}
	return moons
}

func gravity(p1, p2 int) (int, int) {
	if p1 < p2 {
		return 1, -1
	} else if p1 > p2 {
		return -1, 1
	}
	return 0, 0
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func parseInput(filename string) []Moon {
	data, _ := ioutil.ReadFile(filename)
	lines := strings.Split(string(data), "\n")
	moons := make([]Moon, 0, len(lines))
	for _, line := range lines {
		var x, y, z int
		n, _ := fmt.Sscanf(line, "<x=%d, y=%d, z=%d>", &x, &y, &z)
		if n == 3 {
			moon := Moon{Position: Vector{x, y, z}, Velocity: Vector{0, 0, 0}}
			moons = append(moons, moon)
		}
	}
	return moons
}
