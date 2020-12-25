package main

import (
	"log"
	"os"
	"time"

	"github.com/solinox/advent-of-code/2020/pkg/input"
)

func main() {
	publicKeys := input.IntSlice(os.Stdin)

	t0 := time.Now()
	log.Println(part1(publicKeys), time.Since(t0))
}

func part1(publicKeys []int) int {
	cardPublicKey, doorPublicKey := publicKeys[0], publicKeys[1]
	doorLoopSize := getLoopSize(doorPublicKey, 7)
	return encrypt(cardPublicKey, doorLoopSize)
}

func getLoopSize(publicKey, subjectNum int) int {
	val := 1
	loopSize := 0
	for ; val != publicKey; loopSize++ {
		val *= subjectNum
		val %= 20201227
	}
	return loopSize
}

func encrypt(subjectNum, loopSize int) int {
	val := 1
	for i := 0; i < loopSize; i++ {
		val *= subjectNum
		val %= 20201227
	}
	return val
}
