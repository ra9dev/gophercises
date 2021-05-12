//go:generate stringer -type=Suit,Rank

package deck

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

// Suit type
type Suit uint8

// Rank type
type Rank uint8

// Spade and others is a Suit types represented by incremented integers
const (
	Spade Suit = iota
	Diamond
	Club
	Heart
	Joker // special case
)

var suits = [...]Suit{Spade, Diamond, Club, Heart}

// Ace and others is a Rank types represented by incremented integers
// Underscore is for 1-based increment
const (
	_ Rank = iota
	Ace
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

const (
	minRank = Ace
	maxRank = King
)

// Card type
type Card struct {
	Suit
	Rank
}

func (c Card) String() string {
	if c.Suit == Joker {
		return c.Suit.String()
	}

	return fmt.Sprintf("%s of %ss", c.Rank.String(), c.Suit.String())
}

// New deck of cards
func New(opts ...func([]Card) []Card) []Card {
	var cards []Card

	for _, suit := range suits {
		for rank := minRank; rank <= maxRank; rank++ {
			cards = append(cards, Card{Suit: suit, Rank: rank})
		}
	}

	for _, opt := range opts {
		cards = opt(cards)
	}

	return cards
}

// DefaultSort for our deck of cards
func DefaultSort(cards []Card) []Card {
	sort.Slice(cards, Less(cards))
	return cards
}

// Sort creates appropriate custom sorting option for a deck of cards
func Sort(less func(cards []Card) func(i, j int) bool) func([]Card) []Card {
	return func(cards []Card) []Card {
		sort.Slice(cards, less(cards))
		return cards
	}
}

// Less func helper for sorting
func Less(cards []Card) func(i, j int) bool {
	return func(i, j int) bool {
		return absRank(cards[i]) < absRank(cards[j])
	}
}

func absRank(c Card) int {
	return int(c.Suit)*int(maxRank) + int(c.Rank)
}

var shuffleRand = rand.New(rand.NewSource(time.Now().Unix()))

// Shuffle deck of card
func Shuffle(cards []Card) []Card {
	ret := make([]Card, len(cards))
	perm := shuffleRand.Perm(len(cards))

	for i, j := range perm {
		ret[i] = cards[j]
	}

	return ret
}

// Jokers add n Jokers to your deck
func Jokers(n int) func([]Card) []Card {
	return func(cards []Card) []Card {
		for i := 0; i < n; i++ {
			cards = append(cards, Card{
				Suit: Joker,
				Rank: Rank(i),
			})
		}

		return cards
	}
}

// Filter cards with custom filter func
func Filter(f func(card Card) bool) func([]Card) []Card {
	return func(cards []Card) []Card {
		var ret []Card

		for _, c := range cards {
			if f(c) {
				ret = append(ret, c)
			}
		}

		return ret
	}
}

// Deck copy n-times our deck and returns multiple decks in a big slice of cards
func Deck(n int) func([]Card) []Card {
	return func(cards []Card) []Card {
		var ret []Card

		for i := 0; i < n; i++ {
			ret = append(ret, cards...)
		}

		return ret
	}
}
