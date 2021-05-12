package blackjack

import deck "github.com/KubaiDoLove/gophercises_deck_of_cards"

// Score will take in a hand of cards and return the best blackjack score
// possible with the hand.
func Score(hand ...deck.Card) int {
	min := minScore(hand...)
	if min > 11 {
		return min
	}

	for _, c := range hand {
		if c.Rank == deck.Ace {
			// ace is currently worth 1, and we are changing it to be worth 11
			return min + 10
		}
	}

	return min
}

// Blackjack returns true if a hand is a blackjack
func Blackjack(hand ...deck.Card) bool {
	return len(hand) == 2 && Score(hand...) == 21
}

// Soft return true if the score of a hand is a soft score - that is if an ace
// is being counted as 11 points.
func Soft(hand ...deck.Card) bool {
	min := minScore(hand...)
	score := Score(hand...)
	return min != score
}

func minScore(hand ...deck.Card) int {
	score := 0
	for _, c := range hand {
		score += min(int(c.Rank), 10)
	}
	return score
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
