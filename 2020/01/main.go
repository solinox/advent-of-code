package main

import (
	"log"
	"os"

	"github.com/solinox/advent-of-code/2020/pkg/input"
)

const expectedSum = 2020

func main() {
	nums := input.IntSlice(os.Stdin)

	log.Println(part1(nums))

	log.Println(part2(nums))
}

func part1(nums []int) int {
	for i := range nums {
		for j := i + 1; j < len(nums); j++ {
			if nums[i]+nums[j] == expectedSum {
				return nums[i] * nums[j]
			}
		}
	}
	log.Fatal("Not found")
	return 0
}

func part2(nums []int) int {
	for i := range nums {
		for j := i + 1; j < len(nums); j++ {
			for k := j + 1; k < len(nums); k++ {
				if nums[i]+nums[j]+nums[k] == expectedSum {
					return nums[i] * nums[j] * nums[k]
				}
			}
		}
	}
	log.Fatal("Not found")
	return 0
}
