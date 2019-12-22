package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
)

type Sequence []int

func main() {
	basePattern := []int{0, 1, 0, -1}
	sequence := parseInput("input.txt")
	fmt.Println(len(sequence))

	// Part 1
	part1Sequence := fft(sequence, basePattern, 100)
	fmt.Println("Part 1:", part1Sequence[:8])

	// Part 2
	// since offset is in the 2nd half of the sequence
	// and the 2nd half involves a triangle matrix
	// we can do summation tricks working backwards from the end of the sequence
	offset := 5971981
	part2Sequence := fft2(sequence.Repeat(10000), basePattern, 100, offset)
	fmt.Println("Part 2:", part2Sequence[:8])

}

func fft(seq Sequence, basePattern []int, numPhases int) Sequence {
	for i := 0; i < numPhases; i++ {
		newSeq := make(Sequence, len(seq))
		for j := range seq {
			newSeq[j] = seq.GenerateElementWithPattern(basePattern, j)
		}
		seq = newSeq
	}
	return seq
}

func fft2(seq Sequence, basePattern []int, numPhases, offset int) Sequence {
	seq = seq[offset:]
	for i := 0; i < numPhases; i++ {
		newSeq := make(Sequence, len(seq))
		prevSum := 0
		for j := len(newSeq) - 1; j >= 0; j-- {
			sum := prevSum + seq[j]
			newSeq[j] = abs(sum) % 10
			prevSum = sum
		}
		seq = newSeq
	}
	return seq
}

func (s Sequence) Repeat(n int) Sequence {
	newSeq := make(Sequence, len(s)*n)
	for i := range newSeq {
		newSeq[i] = s[i%len(s)]
	}
	return newSeq
}

func (s Sequence) GenerateElementWithPattern(basePattern []int, index int) int {
	ret := 0
	for i := range s {
		patternVal := basePattern[((i+1)/(index+1))%len(basePattern)]
		ret += s[i] * patternVal
	}
	return abs(ret) % 10
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func parseInput(filename string) Sequence {
	data, _ := ioutil.ReadFile(filename)
	data = bytes.TrimSpace(data)
	seq := make([]int, len(data))
	for i := range data {
		seq[i] = int(data[i] - '0')
	}
	return Sequence(seq)
}
