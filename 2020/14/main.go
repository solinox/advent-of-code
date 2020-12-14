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
	program := input.StringSlice(os.Stdin)

	t0 := time.Now()
	log.Println(part1(program), time.Since(t0))

	t0 = time.Now()
	log.Println(part2(program), time.Since(t0))
}

func part1(program []string) uint64 {
	mem := run1(program, make(map[int]uint64))
	sum := uint64(0)
	for _, v := range mem {
		sum += v
	}
	return sum
}

func part2(program []string) uint64 {
	mem := run2(program, make(map[int]uint64))
	sum := uint64(0)
	for _, v := range mem {
		sum += v
	}
	return sum
}

func run1(program []string, mem map[int]uint64) map[int]uint64 {
	mask1s, mask0s := uint64(0), uint64(0)
	for _, line := range program {
		f := strings.Fields(line)
		if len(f) < 3 {
			log.Fatalln("invalid line")
		}
		switch {
		case f[0] == "mask":
			mask := f[2]
			ones, err := strconv.ParseUint(strings.ReplaceAll(mask, "X", "0"), 2, 36)
			if err != nil {
				log.Fatalln(err)
			}
			zeroes, err := strconv.ParseUint(strings.ReplaceAll(mask, "X", "1"), 2, 36)
			if err != nil {
				log.Fatalln(err)
			}
			mask1s, mask0s = ones, zeroes
		case strings.HasPrefix(f[0], "mem"):
			n, err := strconv.ParseUint(f[2], 10, 64)
			if err != nil {
				log.Fatalln(err)
			}
			address, err := strconv.ParseInt(f[0][strings.IndexByte(f[0], '[')+1:strings.IndexByte(f[0], ']')], 10, 64)
			if err != nil {
				log.Fatalln(err)
			}
			mem[int(address)] = n&mask0s | mask1s
		}
	}
	return mem
}

func run2(program []string, mem map[int]uint64) map[int]uint64 {
	var mask string
	for _, line := range program {
		f := strings.Fields(line)
		if len(f) < 3 {
			log.Fatalln("invalid line")
		}
		switch {
		case f[0] == "mask":
			mask = f[2]
		case strings.HasPrefix(f[0], "mem"):
			n, err := strconv.ParseUint(f[2], 10, 64)
			if err != nil {
				log.Fatalln(err)
			}
			address, err := strconv.ParseUint(f[0][strings.IndexByte(f[0], '[')+1:strings.IndexByte(f[0], ']')], 10, 36)
			if err != nil {
				log.Fatalln(err)
			}
			addressB := []byte(strconv.FormatUint(address, 2))
			addressB = append(make([]byte, 36-len(addressB)), addressB...)
			for i := 0; i < len(mask); i++ {
				if addressB[i] == 0 {
					addressB[i] = '0'
				}
				if mask[i] == '0' {
					continue
				}
				addressB[i] = mask[i]
			}
			indexX := make([]int, 0)
			for i := range addressB {
				if addressB[i] == 'X' {
					indexX = append(indexX, i)
				}
			}
			permutations := uint64(1) << len(indexX)
			for i := uint64(0); i < permutations; i++ {
				b := strconv.FormatUint(i, 2)
				if len(b) < len(indexX) {
					b = strings.Repeat("0", len(indexX)-len(b)) + b
				}
				for j, x := range indexX {
					addressB[x] = b[j]
				}
				add, err := strconv.ParseUint(string(addressB), 2, 36)
				if err != nil {
					log.Fatalln(err)
				}
				mem[int(add)] = n
			}
		}
	}
	return mem
}
