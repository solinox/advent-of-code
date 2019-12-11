package main

import (
	"strconv"
	"io/ioutil"
	"strings"
	"fmt"
)

var verbose = false

func main() {
	memory := parseInput("input.txt")

	// Part 1
	part1 := 0
	ch1 := make(chan []int)
	go generatePermutations([]int{0,1,2,3,4}, ch1)
	for phaseSettings := range ch1 {
		outA := make(chan int, 2)
		outB := make(chan int, 2)
		outC := make(chan int, 2)
		outD := make(chan int, 2)
		outE := make(chan int, 2)

		inA := make(chan int, 2)
		inA <- phaseSettings[0]
		outA <- phaseSettings[1]
		outB <- phaseSettings[2]
		outC <- phaseSettings[3]
		outD <- phaseSettings[4]

		inA <- 0
		runProgram(memory, inA, outA, nil)
		runProgram(memory, outA, outB, nil)
		runProgram(memory, outB, outC, nil)
		runProgram(memory, outC, outD, nil)
		runProgram(memory, outD, outE, nil)

		output := <-outE
		
		if output > part1 {
			part1 = output
		}
	}
	fmt.Println("Part 1", part1)

	// Part 2
	part2 := 0
	ch2 := make(chan []int)
	go generatePermutations([]int{5,6,7,8,9}, ch2)
	for phaseSettings := range ch2 {
		outA := make(chan int, 2)
		outB := make(chan int, 2)
		outC := make(chan int, 2)
		outD := make(chan int, 2)
		outE := make(chan int, 2)

		outE <- phaseSettings[0]
		outA <- phaseSettings[1]
		outB <- phaseSettings[2]
		outC <- phaseSettings[3]
		outD <- phaseSettings[4]

		outE <- 0
		done := make(chan struct{})
		go runProgram(memory, outE, outA, nil)
		go runProgram(memory, outA, outB, nil)
		go runProgram(memory, outB, outC, nil)
		go runProgram(memory, outC, outD, nil)
		go runProgram(memory, outD, outE, done)

		<- done
		output := <-outE
		
		if output > part2 {
			part2 = output
		}
	}
	fmt.Println("Part 2", part2)

}

func generatePermutations(slice []int, ch chan []int) {
	// heap's algorithm from stackoverflow
	var helper func([]int, int)

    helper = func(slice []int, n int){
        if n == 1 {
            tmp := make([]int, len(slice))
            copy(tmp, slice)
            ch <- tmp
        } else {
            for i := 0; i < n; i++{
                helper(slice, n - 1)
                if n % 2 == 1{
                    tmp := slice[i]
                    slice[i] = slice[n - 1]
                    slice[n - 1] = tmp
                } else {
                    tmp := slice[0]
                    slice[0] = slice[n - 1]
                    slice[n - 1] = tmp
                }
            }
        }
    }
    helper(slice, len(slice))
    close(ch)
} 

func runProgram(initialMemory []int, input <-chan int, output chan<- int, done chan<- struct{}) (outputValue int) {
	memory := make([]int, len(initialMemory))
	copy(memory, initialMemory)

	for i := 0; i < len(memory); {
		n := memory[i]
		op := n % 100
		n /= 100
		msg := ""

		switch op {
		case 1:
			a, b, c := memory[i+1], memory[i+2], memory[i+3]
			params := getParams(memory, n, a, b)
			memory[c] = params[0] + params[1]
			msg = fmt.Sprintf("i=%d\tOp 1\tmemory[%d] = %d\t%v", i, c, memory[c], memory[i:i+4])
			i += 4
		case 2:
			a, b, c := memory[i+1], memory[i+2], memory[i+3]
			params := getParams(memory, n, a, b)
			memory[c] = params[0] * params[1]
			msg = fmt.Sprintf("i=%d\tOp 2\tmemory[%d] = %d\t%v", i, c, memory[c], memory[i:i+4])
			i += 4
		case 3:
			a := memory[i+1]
			memory[a] = <-input
			msg = fmt.Sprintf("i=%d\tOp 3\tInput %d to memory[%d]", i, input, a)
			i += 2
		case 4:
			a := memory[i+1]
			output <- memory[a]
			outputValue = memory[a]
			msg = fmt.Sprintf("i=%d\tOp 4\tOutput %d", i, memory[a])
			i += 2
		case 5:
			a, b := memory[i+1], memory[i+2]
			params := getParams(memory, n, a, b)
			msg = fmt.Sprintf("i=%d\tOp 5\tSet i to %d if %d != 0", i, params[1], params[0])
			if params[0] != 0 {
				i = params[1]
			} else {
				i += 3
			}
		case 6:
			a, b := memory[i+1], memory[i+2]
			params := getParams(memory, n, a, b)
			msg = fmt.Sprintf("i=%d\tOp 6\tSet i to %d if %d == 0", i, params[1], params[0])
			if params[0] == 0 {
				i = params[1]
			} else {
				i += 3
			}
		case 7:
			a, b, c := memory[i+1], memory[i+2], memory[i+3]
			params := getParams(memory, n, a, b)
			if params[0] < params[1] {
				memory[c] = 1
			} else {
				memory[c] = 0
			}
			i += 4
		case 8:
			a, b, c := memory[i+1], memory[i+2], memory[i+3]
			params := getParams(memory, n, a, b)
			if params[0] == params[1] {
				memory[c] = 1
			} else {
				memory[c] = 0
			}
			i += 4
		case 99:
			if done != nil {
				done <- struct{}{}
			}
			return outputValue
		default:
			panic(fmt.Sprintf("Op %d is not supported", op))
		}

		if msg != "" && verbose {
			fmt.Println(msg)
		}
	}
	fmt.Println("Program did not halt but went out of bounds unexpectedly")
	if done != nil {
		done <- struct{}{}
	}
	return outputValue
}

func getParams(memory []int, n int, params ...int) []int {
	for i := range params {
		mode := n % 10
		if mode == 0 {
			params[i] = memory[params[i]]
		} else if mode == 1 {
			params[i] = params[i]
		}
		n /= 10
	}
	return params
}

func parseInput(filename string) []int {
	data, _ := ioutil.ReadFile(filename)
	fields := strings.Split(string(data), ",")
	ints := make([]int, len(fields))
	for i := range fields {
		ints[i], _ = strconv.Atoi(fields[i])
	}
	return ints
}
