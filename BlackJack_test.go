package main

import (
	"testing"

	"fyne.io/fyne/app"
)

// 0-3   Aces
// 4-7   Twoes
// 8-11  Threes
// 12-15 Fours
// 16-19 Fives
// 20-23 Sixes
// 24-27 Sevens
// 28-31 Eights
// 32-35 Nines
// 36-39 Tens
// 40-43 Jacks
// 44-47 Queens
// 48-51 Kings

// Tests if the deck is correctly created
func TestSetDeck(t *testing.T) {
	SetDeck()
	for i := 0; i < 51; i++ {
		for j := i + 1; j < 52; j++ {
			if deck[i] == deck[j] {
				t.Error("The deck contains same cards!")
			}
		}
	}
}

// Tests if the output is correct when showing a card
func TestShowCard(t *testing.T) {
	SetDeck()
	tDeck := make([]int, 3)
	tDeck[0] = 0
	tDeck[1] = 1
	tDeck[2] = 51

	if ShowCard(tDeck[0]) != "Ace Clubs" || ShowCard(tDeck[1]) != "Ace Diamond" || ShowCard(tDeck[2]) != "King Spade" {
		t.Error("Wrong output for showing cards")
	}
}

// Test correct output for showing a whole hand
func TestShowHand(t *testing.T) {
	SetDeck()
	tHand := make([]int, 3)
	tHand[0] = 2
	tHand[1] = 4
	tHand[2] = 50

	if ShowHand(tHand) != "Ace Heart, 2 Clubs, King Heart" {
		t.Error("Wrong output for showing a whole hand")
	}
}

func TestCalculatePointsNormalCase(t *testing.T) {
	SetDeck()
	playersHand := make([]int, 3)
	playersHand[0] = 50 // King
	playersHand[1] = 4  // 2
	playersHand[2] = 8  // 3

	CalculatePoints(0, playersHand)
	if playersPoints != 15 {
		t.Error("Wrong points calculation")
	}
}
func TestCalculatePointsOneAceNotBusting(t *testing.T) {
	SetDeck()
	playersHand := make([]int, 3)
	playersHand[0] = 2 // Ace
	playersHand[1] = 4 // 2
	playersHand[2] = 8 // 3

	CalculatePoints(0, playersHand)
	if playersPoints != 16 {
		t.Error("Wrong points calculation")
	}
}

func TestCalculatePointsReduceAcePoints(t *testing.T) {
	SetDeck()
	playersHand := make([]int, 3)
	playersHand[0] = 2  // Ace
	playersHand[1] = 4  // 2
	playersHand[2] = 51 // King

	CalculatePoints(0, playersHand)
	if playersPoints != 13 {
		t.Error("Wrong aces points calculation")
	}
}

func TestCalculatePointsTwoAces(t *testing.T) {
	SetDeck()
	playersHand := make([]int, 2)
	playersHand[0] = 0 // Ace
	playersHand[1] = 1 // Ace

	CalculatePoints(0, playersHand)
	if playersPoints != 12 {
		t.Error("Wrong aces points calculation")
	}
}
func TestCalculatePointsMultipleAces(t *testing.T) {
	SetDeck()
	playersHand := make([]int, 5)
	playersHand[0] = 2  // Ace
	playersHand[1] = 4  // 2
	playersHand[2] = 51 // King
	playersHand[3] = 1  // Ace
	playersHand[4] = 3  // Ace

	CalculatePoints(0, playersHand)
	if playersPoints != 15 {
		t.Error("Wrong aces points calculation")
	}
}

func TestHasBlackJack(t *testing.T) {
	SetDeck()
	playersHand := make([]int, 2)

	for i := 0; i < 4; i++ {
		playersHand[0] = i // Ace
		for j := 36; j < 52; j++ {
			playersHand[1] = j // Any 10 points card
			if !HasBlackJack(playersHand) {
				t.Error("Not detecting black jack")
			}
		}
	}
}

func TestHasBlackJackFalse(t *testing.T) {
	SetDeck()
	playersHand := make([]int, 2)

	for i := 0; i < 4; i++ {
		playersHand[0] = i // Ace
		for j := i; j < 36; j++ {
			playersHand[1] = j // Any 10 points card
			if HasBlackJack(playersHand) {
				t.Error("Wrong Black Jack detection")
			}
		}
	}

	for i := 4; i < 52; i++ {
		playersHand[0] = i // Any non-Ace card
		for j := 4; j < 52; j++ {
			playersHand[1] = j // Any non-Ace card
			if HasBlackJack(playersHand) {
				t.Error("Wrong Black Jack detection")
			}
		}
	}

	// First 2 cards form BlackJack but the third ruins it
	playersHand[0] = 0
	playersHand[1] = 50
	playersHand = append(playersHand, 0)
	for j := 0; j < 52; j++ {
		playersHand[2] = j // Any card
		if HasBlackJack(playersHand) {
			t.Error("Wrong Black Jack detection")
		}
	}
}

func TestDraw(t *testing.T) {
	SetDeck()
	index = 0

	for i := 0; i < 8; i++ { // Draw 8 cards
		Draw(2)
	}

	for i := 0; i < 4; i++ { // The first 4 cards should be aces
		if deck[dealersHand[i]].value != 1 {
			t.Error("Wrong cards drawing")
		}
	}

	for i := 4; i < 8; i++ { // The second 4 cards should be 2s
		if deck[dealersHand[i]].value != 2 {
			t.Error("Wrong cards drawing")
		}
	}

	if dealersPoints != 22 {
		t.Error("Wrong points calculation in the Draw method")
	}
}

func TestHitWithMainHand(t *testing.T) {
	SetDeck()
	var labels Labels
	var table Table
	a := app.New()
	w := a.NewWindow("Test")

	index = 28 // Eight
	Hit(false, &table, &labels, w)
	if splitPoints != 0 && playersPoints != 8 {
		t.Error("Draw fails from hit option")
	}

	index = 44 // Jack
	Hit(false, &table, &labels, w)
	if splitPoints != 0 && playersPoints != 18 {
		t.Error("Draw fails from hit option")
	}

	index = 16 // Three
	Hit(false, &table, &labels, w)
	if splitPoints != 0 && playersPoints != 23 {
		t.Error("Draw fails from hit option")
	}
}

func TestHitWithSplittedHand(t *testing.T) {
	SetDeck()
	var labels Labels
	var table Table
	a := app.New()
	w := a.NewWindow("Test")

	index = 32 // Nine
	Hit(true, &table, &labels, w)
	if splitPoints != 9 && playersPoints != 0 {
		t.Error("Draw fails from hit option")
	}

	index = 48 // Queen
	Hit(true, &table, &labels, w)
	if splitPoints != 19 && playersPoints != 0 {
		t.Error("Draw fails from hit option")
	}

	index = 8 // Three
	Hit(true, &table, &labels, w)
	if splitPoints != 22 && playersPoints != 0 {
		t.Error("Draw fails from hit option")
	}
}

func TestStand(t *testing.T) {
	SetDeck()
	var labels Labels
	var table Table

	Shuffle()
	index = 0
	dealersPoints = 0
	dealersHand = nil
	Stand(&labels, &table)
	if dealersPoints < 17 || dealersPoints > 26 {
		t.Error("Stand option fails")
	}
}

func TestGameResultPlayerWins(t *testing.T) {
	SetDeck()
	var labels Labels
	var table Table
	a := app.New()
	w := a.NewWindow("Test")

	playersPoints = 0
	playersHand = nil

	index = 0 // ace
	Hit(false, &table, &labels, w)
	index = 20 // Six
	Hit(false, &table, &labels, w)
	index = 21 // Six
	Hit(false, &table, &labels, w)
	index = 28 // Eight
	Hit(false, &table, &labels, w)

	if playersPoints != 21 {
		t.Error("Wrong cards drawing with hit option")
	}

	index = 3
	dealersPoints = 0
	dealersHand = nil
	Stand(&labels, &table) // Ace, Two, Two, Two
	if dealersPoints != 17 {
		t.Error("Stand option fails")
	}

	bet = 10
	money = 200
	GameResult(playersPoints, playersHand, false, false, &labels, w)
	if money != 220 {
		t.Error("Wrong game result or earnings calculations")
	}
}

func TestGameResultTie(t *testing.T) {
	SetDeck()
	var labels Labels
	var table Table
	a := app.New()
	w := a.NewWindow("Test")

	playersPoints = 0
	playersHand = nil

	index = 0 // ace
	Hit(false, &table, &labels, w)
	index = 20 // Six
	Hit(false, &table, &labels, w)
	index = 21 // Six
	Hit(false, &table, &labels, w)
	index = 28 // Eight
	Hit(false, &table, &labels, w)

	if playersPoints != 21 {
		t.Error("Wrong cards drawing with hit option")
	}

	dealersHand = make([]int, 3)
	dealersHand[0] = 1  // Ace
	dealersHand[1] = 16 // Five
	dealersHand[2] = 17 // Five
	CalculatePoints(2, dealersHand)
	if dealersPoints != 21 {
		t.Error("Wrong points calculation")
	}

	bet = 10
	money = 200
	GameResult(playersPoints, playersHand, false, false, &labels, w)
	if money != 210 { // The player should receive back his bet
		t.Error("Wrong game result or earnings calculations")
	}
}

func TestGameResultDealersBJ(t *testing.T) {
	SetDeck()
	var labels Labels
	var table Table
	a := app.New()
	w := a.NewWindow("Test")

	playersPoints = 0
	playersHand = nil

	index = 0 // ace
	Hit(false, &table, &labels, w)
	index = 20 // Six
	Hit(false, &table, &labels, w)
	index = 21 // Six
	Hit(false, &table, &labels, w)
	index = 28 // Eight
	Hit(false, &table, &labels, w)

	if playersPoints != 21 {
		t.Error("Wrong cards drawing with hit option")
	}

	dealersHand = make([]int, 2)
	dealersHand[0] = 1  // Ace
	dealersHand[1] = 50 // King
	CalculatePoints(2, dealersHand)
	if !HasBlackJack(dealersHand) {
		t.Error("Not detecting Black Jack")
	}

	money = 200
	GameResult(playersPoints, playersHand, false, false, &labels, w)
	if money != 200 { // The betting is on earlier stage of the game so after lossing the game the player mustn't lose any money
		t.Error("Wrong game result or earnings calculations")
	}
}
func TestGameResultWinWithBJ(t *testing.T) {
	SetDeck()
	var labels Labels
	var table Table
	a := app.New()
	w := a.NewWindow("Test")

	playersPoints = 0
	playersHand = nil

	index = 0 // ace
	Hit(false, &table, &labels, w)
	index = 51 // King
	Hit(false, &table, &labels, w)

	if playersPoints != 21 {
		t.Error("Wrong cards drawing with hit option")
	}

	dealersHand = make([]int, 3)
	dealersHand[0] = 1  // Ace
	dealersHand[1] = 16 // Five
	dealersHand[2] = 17 // Five
	CalculatePoints(2, dealersHand)
	if dealersPoints != 21 {
		t.Error("Wrong points calculation")
	}

	bet = 10
	money = 200
	GameResult(playersPoints, playersHand, false, false, &labels, w)
	if money != 225 {
		t.Error("Wrong game result or earnings calculations")
	}
}
