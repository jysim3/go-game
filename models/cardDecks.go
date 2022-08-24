package models

import (
	"math/rand"
)

type CardsDeck struct {
	cards     []Card
	discarded []Card
}

type Card struct {
	Number int `json:"number"`
	Suit   int `json:"suit"`
}

func (c *CardsDeck) DealCards(n int) []Card {
	var cards []Card
	for i := 0; i < n; i++ {
		cards = append(cards, c.Deal())
	}
	return cards
}

func (c *CardsDeck) Deal() Card {
	randInt := rand.Intn(len(c.cards))
	card := c.cards[randInt]
	c.discarded = append(c.discarded, card)
	c.cards[randInt], c.cards = c.cards[len(c.cards)-1], c.cards[:len(c.cards)-1]
	return card
}

func (c *CardsDeck) NewDeck(nDeck int) bool {
	if nDeck < 1 {
		return false
	}
	c.cards = make([]Card, 4*13*nDeck)
	i := 0
	for Deck := 0; Deck < nDeck; Deck++ {
		for Suit := 0; Suit < 4; Suit++ {
			for Number := 1; Number < 14; Number++ {
				c.cards[i] = Card{Suit: Suit, Number: Number}
				i += 1
			}
		}
	}
	return true
}

func (c *CardsDeck) RefreshDeck() {
}
