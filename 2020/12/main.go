package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/solinox/advent-of-code/2020/pkg/input"
	"github.com/solinox/advent-of-code/2020/pkg/intmath"
)

func main() {
	instructions := input.StringSlice(os.Stdin)

	t0 := time.Now()
	log.Println(part1(instructions), time.Since(t0))

	t0 = time.Now()
	log.Println(part2(instructions), time.Since(t0))
}

type vector struct {
	X, Y int
}

type vessel struct {
	Position  vector
	Direction vector // in part 2, the waypoint is the direction
}

func part1(instructions []string) int {
	ferry := vessel{Position: vector{0, 0}, Direction: vector{1, 0}}
	for _, instruction := range instructions {
		ferry = ferry.Do(instruction, false)
	}
	return intmath.Abs(ferry.Position.X) + intmath.Abs(ferry.Position.Y)
}

func part2(instructions []string) int {
	ferry := vessel{Position: vector{0, 0}, Direction: vector{10, 1}}
	for _, instruction := range instructions {
		ferry = ferry.Do(instruction, true)
	}
	return intmath.Abs(ferry.Position.X) + intmath.Abs(ferry.Position.Y)
}

func (v vessel) Do(instruction string, part2 bool) vessel {
	command := instruction[0]
	n, err := strconv.Atoi(instruction[1:])
	if err != nil {
		log.Fatalln(instruction, err)
	}
	// part2 changes the meaning of the N,E,S,W instructions
	switch {
	case command == 'N' && part2:
		v.Direction.Y += n
	case command == 'N' && !part2:
		v.Position.Y += n
	case command == 'E' && part2:
		v.Direction.X += n
	case command == 'E' && !part2:
		v.Position.X += n
	case command == 'S' && part2:
		v.Direction.Y -= n
	case command == 'S' && !part2:
		v.Position.Y -= n
	case command == 'W' && part2:
		v.Direction.X -= n
	case command == 'W' && !part2:
		v.Position.X -= n
	case command == 'F':
		v.Position = vector{v.Position.X + n*v.Direction.X, v.Position.Y + n*v.Direction.Y}
	case command == 'L':
		n = -n
		fallthrough
	case command == 'R':
		for n < 0 {
			n += 360
		}
		angle := n % 360
		switch angle {
		case 0:
			// do nothing
		case 90:
			v.Direction = vector{v.Direction.Y, -v.Direction.X}
		case 180:
			v.Direction = vector{-v.Direction.X, -v.Direction.Y}
		case 270:
			v.Direction = vector{-v.Direction.Y, v.Direction.X}
		default:
			log.Fatalln("unsupported turn angle", angle)
		}
	default:
		log.Fatalln("unsupported command", command)
	}
	return v
}
