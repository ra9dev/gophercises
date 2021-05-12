package blackjack

import (
	deck "github.com/KubaiDoLove/gophercises_deck_of_cards"
)

// AI common interface
type AI interface {
	Bet(shuffled bool) int
	Play([]deck.Card, deck.Card) Move
	Results([][]deck.Card, []deck.Card)
}
