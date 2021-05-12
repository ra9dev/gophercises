package blackjack

import (
	"errors"

	deck "github.com/KubaiDoLove/gophercises_deck_of_cards"
)

var (
	errBust = errors.New("hand score exceeded 21")
)

// Move is a func for changing game state
type Move func(*Game) error

// MoveHit for drawing a card
func MoveHit(g *Game) error {
	hand := g.currentHand()
	var card deck.Card

	card, g.deck = draw(g.deck)
	*hand = append(*hand, card)

	if Score(*hand...) > 21 {
		return errBust
	}

	return nil
}

// MoveSplit ...
func MoveSplit(g *Game) error {
	curHand := *g.currentHand()

	if len(curHand) != 2 {
		return errors.New("you can only split with 2 cards in your hand")
	}
	if curHand[0].Rank != curHand[1].Rank {
		return errors.New("both cards must have the same rank to split")
	}
	g.player = append(g.player, hand{
		cards: []deck.Card{curHand[1]},
		bet:   g.player[g.handIdx].bet,
	})
	g.player[g.handIdx].cards = curHand[:1]

	return nil
}

// MoveDouble ...
func MoveDouble(g *Game) error {
	if len(*g.currentHand()) != 2 {
		return errors.New("Can only double on a hand with 2 cards")
	}
	g.playerBet *= 2
	MoveHit(g)

	return MoveStand(g)
}

// MoveStand for updating state
func MoveStand(g *Game) error {
	if g.state == stateDealerTurn {
		g.state++
		return nil
	}
	if g.state == statePlayerTurn {
		g.handIdx++
		if g.handIdx == len(g.player) {
			g.state++
		}
		return nil
	}

	return errors.New("invalid state")
}

func draw(cards []deck.Card) (deck.Card, []deck.Card) {
	return cards[0], cards[1:]
}
