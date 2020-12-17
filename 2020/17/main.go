package main

import (
	"log"
	"os"
	"time"

	"github.com/solinox/advent-of-code/2020/pkg/input"
)

type point struct {
	W, X, Y, Z int
}

type pocketDimension struct {
	MinW, MaxW int
	MinX, MaxX int
	MinY, MaxY int
	MinZ, MaxZ int
	Cubes      map[point]bool
}

type change struct {
	Point  point
	Active bool
}

func main() {
	data := input.BytesSlice(os.Stdin)

	t0 := time.Now()
	log.Println(part1(data), time.Since(t0))

	t0 = time.Now()
	log.Println(part2(data), time.Since(t0))
}

func part1(data [][]byte) int {
	pd := newPocketDimension(data)
	pd = pd.CycleN(6, false)
	activeCubes := 0
	for _, active := range pd.Cubes {
		if active {
			activeCubes++
		}
	}
	return activeCubes
}

func part2(data [][]byte) int {
	pd := newPocketDimension(data)
	pd = pd.CycleN(6, true)
	activeCubes := 0
	for _, active := range pd.Cubes {
		if active {
			activeCubes++
		}
	}
	return activeCubes
}

func newPocketDimension(data [][]byte) pocketDimension {
	cubes := make(map[point]bool)
	for y := range data {
		for x := range data[y] {
			cubes[point{0, x, y, 0}] = data[y][x] == '#'
		}
	}
	return pocketDimension{
		MinW: 0, MaxW: 0,
		MinX: 0, MaxX: len(data[0]) - 1,
		MinY: 0, MaxY: len(data) - 1,
		MinZ: 0, MaxZ: 0,
		Cubes: cubes,
	}
}

func (pd pocketDimension) CycleN(n int, part2 bool) pocketDimension {
	for i := 0; i < n; i++ {
		pd = pd.Cycle(part2)
	}
	return pd
}

func (pd pocketDimension) Cycle(part2 bool) pocketDimension {
	// Only check active cubes and their neighbors to see if state should change
	cubesToCheck := make(map[point]bool)
	actives := make([]point, 0)
	for pt, active := range pd.Cubes {
		if active {
			actives = append(actives, pt)
			cubesToCheck[pt] = active
			for z := pt.Z - 1; z <= pt.Z+1; z++ {
				for y := pt.Y - 1; y <= pt.Y+1; y++ {
					for x := pt.X - 1; x <= pt.X+1; x++ {
						for w := pt.W - 1; w <= pt.W+1; w++ {
							ww := w
							if !part2 {
								ww = 0
							}
							neighborPt := point{ww, x, y, z}
							if _, ok := cubesToCheck[neighborPt]; ok {
								continue
							}
							cubesToCheck[neighborPt] = false
						}
					}
				}
			}
		}
	}

	// Get list of cubes that will need to change state
	changes := make([]change, 0)
	for pt, active := range cubesToCheck {
		numActiveNeighbors := numActiveNeighbors(actives, pt)
		if active {
			if numActiveNeighbors == 2 || numActiveNeighbors == 3 {
				continue
			}
			changes = append(changes, change{Point: pt, Active: false})
		}
		if numActiveNeighbors != 3 {
			continue
		}
		changes = append(changes, change{Point: pt, Active: true})
	}

	// apply changes
	for _, change := range changes {
		pd.Cubes[change.Point] = change.Active
	}

	return pd
}

func numActiveNeighbors(actives []point, pt point) int {
	numActiveNeighbors := 0
	for i := range actives {
		dw, dx, dy, dz := actives[i].W-pt.W, actives[i].X-pt.X, actives[i].Y-pt.Y, actives[i].Z-pt.Z
		if dw >= -1 && dw <= 1 && dx >= -1 && dx <= 1 && dy >= -1 && dy <= 1 && dz >= -1 && dz <= 1 {
			if dw == 0 && dx == 0 && dy == 0 && dz == 0 {
				continue
			}
			numActiveNeighbors++
			// Optimization: return early since rules only care about up to 3 active neighbors
			if numActiveNeighbors > 3 {
				return numActiveNeighbors
			}
		}
	}
	return numActiveNeighbors
}
