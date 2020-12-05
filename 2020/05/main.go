package main

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/solinox/advent-of-code/2020/pkg/input"
)

func main() {
	seats := input.StringSlice(os.Stdin)

	t0 := time.Now()
	log.Println(part1(seats), time.Since(t0))

	t0 = time.Now()
	log.Println(part2(seats), time.Since(t0))
}

func part1(seats []string) uint16 {
	max := uint16(0)
	for _, seat := range seats {
		if n := seatID(seat); n > max {
			max = n
		}
	}
	return max
}

func part2(seats []string) uint16 {
	takenLen := uint16(1024)
	taken := make([]bool, takenLen)
	for _, seat := range seats {
		taken[seatID(seat)] = true
	}
	minFound := false
	for i := uint16(0); i < takenLen; i++ {
		if taken[i] && !minFound {
			minFound = true
		}
		if minFound && !taken[i] {
			return i
		}
	}
	log.Fatalln("Not found")
	return 0
}

var replacer = strings.NewReplacer("F", "0", "B", "1", "L", "0", "R", "1")

func seatID(seat string) uint16 {
	binaryString := replacer.Replace(seat)
	n, err := strconv.ParseUint(binaryString, 2, 10)
	if err != nil {
		log.Fatalln(seat, binaryString, err)
	}
	return uint16(n)
}
