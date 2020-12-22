package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/solinox/advent-of-code/2020/pkg/input"
)

func main() {
	players := input.SplitString(os.Stdin, "\n\n")
	deck1 := input.IntSlice(strings.NewReader(players[0][strings.IndexByte(players[0], ':')+2:]))
	deck2 := input.IntSlice(strings.NewReader(players[1][strings.IndexByte(players[1], ':')+2:]))

	t0 := time.Now()
	log.Println(part1(deck1, deck2), time.Since(t0))

	t0 = time.Now()
	log.Println(part2(deck1, deck2), time.Since(t0))
}

func part1(deck1, deck2 []int) int {
	for len(deck1) > 0 && len(deck2) > 0 {
		deck1, deck2 = combat(deck1, deck2)
	}
	winner := deck1
	if len(deck2) > 0 {
		winner = deck2
	}
	return score(winner)
}

func score(deck []int) int {
	score := 0
	for i := len(deck) - 1; i >= 0; i-- {
		score += deck[i] * (len(deck) - i)
	}
	return score
}

type player int

const (
	player1 player = iota
	player2
)

func part2(deck1, deck2 []int) int {
	var winner player
	state := make(map[int]struct{})
	for len(deck1) > 0 && len(deck2) > 0 {
		deck1, deck2, winner = recursiveCombat(deck1, deck2, state)
	}
	if winner == player1 {
		return score(deck1)
	}
	return score(deck2)
}

func combat(deck1, deck2 []int) ([]int, []int) {
	min := len(deck1)
	if len(deck2) < min {
		min = len(deck2)
	}

	for i := 0; i < min; i++ {
		if deck1[i] > deck2[i] {
			deck1 = append(deck1, deck1[i], deck2[i])
		} else {
			deck2 = append(deck2, deck2[i], deck1[i])
		}
	}
	return deck1[min:], deck2[min:]
}

func recursiveCombat(deck1, deck2 []int, state map[int]struct{}) ([]int, []int, player) {
	hash := hash(deck1, deck2)
	if _, ok := state[hash]; ok {
		deck1 = append(deck1[1:], deck1[0], deck2[0])
		deck2 = deck2[1:]
		return deck1, deck2, player1
	}
	state[hash] = struct{}{}
	card1, card2 := deck1[0], deck2[0]
	if len(deck1) > card1 && len(deck2) > card2 {
		subDeck1, subDeck2 := make([]int, card1), make([]int, card2)
		copy(subDeck1, deck1[1:1+card1])
		copy(subDeck2, deck2[1:1+card2])
		var winner player
		subState := make(map[int]struct{})
		for len(subDeck1) > 0 && len(subDeck2) > 0 {
			subDeck1, subDeck2, winner = recursiveCombat(subDeck1, subDeck2, subState)
		}
		if winner == player1 {
			deck1 = append(deck1[1:], card1, card2)
			deck2 = deck2[1:]
			return deck1, deck2, winner
		}
		deck1 = deck1[1:]
		deck2 = append(deck2[1:], card2, card1)
		return deck1, deck2, winner
	}
	if card1 > card2 {
		deck1 = append(deck1[1:], card1, card2)
		deck2 = deck2[1:]
		return deck1, deck2, player1
	}
	deck1 = deck1[1:]
	deck2 = append(deck2[1:], card2, card1)
	return deck1, deck2, player2
}

func hash(deck1, deck2 []int) int {
	return (score(deck1) + 1) * (score(deck2) + 2)
}
