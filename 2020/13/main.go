package main

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/solinox/advent-of-code/2020/pkg/input"
	"github.com/solinox/advent-of-code/2020/pkg/intmath"
)

func main() {
	shuttleInfo := input.StringSlice(os.Stdin)

	t0 := time.Now()
	log.Println(part1(shuttleInfo), time.Since(t0))

	t0 = time.Now()
	log.Println(part2(shuttleInfo[1]), time.Since(t0))
}

func part1(shuttleInfo []string) int {
	soonestTimestamp, err := strconv.Atoi(shuttleInfo[0])
	if err != nil {
		log.Fatalln(err)
		return 0
	}
	shuttleTimes := strings.Split(shuttleInfo[1], ",")
	validShuttleTimes := make([]int, 0)
	for _, shuttleTime := range shuttleTimes {
		n, err := strconv.Atoi(shuttleTime)
		if err != nil {
			continue
		}
		validShuttleTimes = append(validShuttleTimes, n)
	}

	min := soonestTimestamp
	busID := 0
	for _, shuttleTime := range validShuttleTimes {
		if next := (soonestTimestamp/shuttleTime + 1) * shuttleTime; next-soonestTimestamp < min {
			min = next - soonestTimestamp
			busID = shuttleTime
		}
	}
	return min * busID
}

type shuttle struct {
	ID     int
	Offset int
}

func part2(shuttleInfo string) int {
	shuttles := make([]shuttle, 0)
	for i, shuttleTime := range strings.Split(shuttleInfo, ",") {
		n, err := strconv.Atoi(shuttleTime)
		if err != nil {
			continue
		}
		shuttles = append(shuttles, shuttle{ID: n, Offset: i})
	}

	iter := shuttles[0].ID
	lcmInts := []int{shuttles[0].ID}
	shuttles = shuttles[1:]
	for t := 0; len(shuttles) > 0; t += iter {
		for i := 0; i < len(shuttles); i++ {
			if (t+shuttles[i].Offset)%shuttles[i].ID == 0 {
				lcmInts = append(lcmInts, shuttles[i].ID)
				shuttles = append(shuttles[:i], shuttles[i+1:]...)
				if len(shuttles) == 0 {
					return t
				}
				iter = intmath.LCM(lcmInts...)
			} else {
				break
			}
		}
	}
	return 0
}
