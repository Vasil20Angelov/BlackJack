package main

import (
	"image/color"
	"math/rand"
	"strconv"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

const DeckSize int = 52
const minBet int = 10

var index int = 4   // index in the deck; starts from 4 because the first 4 cards are automatically drawn each game
var money int = 200 // starting balance
var bet int = 10    // default bet
var deck []Card     // the playing deck

var playersPoints int = 0
var dealersPoints int = 0
var splitPoints int = 0

// the elements of these slices are indexes (in the deck)
var playersHand []int
var dealersHand []int
var splitHand []int

type Card struct {
	value int // 1,...,13, Where 1 = Ace, 11 = Jack, 12 = Queen, 13 = King
	suit  int // 1 = clubs, 2 = diamond, 3 = heart, 4 = spade
}

type Labels struct {
	balanceLabel canvas.Text
	balance      canvas.Text
	betLabel     canvas.Text
	currBet      canvas.Text

	pLabel canvas.Text
	dLabel canvas.Text
	sLabel canvas.Text

	pCards canvas.Text
	dCards canvas.Text
	sCards canvas.Text

	PntLab1 canvas.Text
	PntLab2 canvas.Text
	PntLab3 canvas.Text

	pPoints canvas.Text
	dPoints canvas.Text
	sPoints canvas.Text
}

type Table struct {
	pHand      string
	dHand      string
	sHand      string
	pPointsStr string
	dPointsStr string
	sPointsStr string
}

// Converts to string the card that is at index i in the deck
func ShowCard(i int) string {
	var output string

	switch deck[i].value {
	case 1:
		output = "Ace"
	case 11:
		output = "Jack"
	case 12:
		output = "Queen"
	case 13:
		output = "King"
	default:
		output = strconv.Itoa(deck[i].value)
	}

	switch deck[i].suit {
	case 1:
		output += " Clubs"
	case 2:
		output += " Diamond"
	case 3:
		output += " Heart"
	case 4:
		output += " Spade"
	}

	return output
}

// Converts to string all cards in a given hand
func ShowHand(hand []int) string {
	var output string
	for i := 0; i < len(hand)-1; i++ {
		output += ShowCard(hand[i]) + ", "
	}
	output += ShowCard(hand[len(hand)-1])

	return output
}

// Initializes the deck
func SetDeck() {
	deck = make([]Card, DeckSize)

	var num int = 0
	for i := 1; i <= 13; i++ {
		for j := 1; j <= 4; j++ {
			deck[num].value = i
			deck[num].suit = j
			num++
		}
	}
}

// Shuffles the deck by swapping cards at random indexes
func Shuffle() {
	rand.Seed(time.Now().UnixNano())
	var randNum int
	var temp Card
	for i := 0; i < DeckSize; i++ {
		randNum = rand.Intn(DeckSize)
		temp = deck[i]

		deck[i] = deck[randNum]
		deck[randNum] = temp
	}
}

// Checks if a given hand has a "Black Jack" - it means that the hand has 21 points and contains 2 cards - an ace and a ten-points card
func HasBlackJack(hand []int) bool {
	if len(hand) == 2 && ((deck[hand[0]].value == 1 && deck[hand[1]].value >= 10) || (deck[hand[1]].value == 1 && deck[hand[0]].value >= 10)) {
		return true
	}

	return false
}

// Calculates how many points are in the given hand. The variable "player" means: 0 - player, 1 - splitted hand, 2 - dealer
func CalculatePoints(player int, hand []int) {
	var points int = 0
	var acesCount int = 0
	for _, element := range hand {
		if deck[element].value == 1 {
			acesCount++
		} else if deck[element].value > 9 {
			points += 10
		} else {
			points += deck[element].value
		}
	}

	// Each ace gives 1 or 11 points; always choosing the better variant
	var acesAsEleven int = 0
	for i := acesCount; i > 0; i-- {
		if points+11 < 22 {
			points += 11
			acesAsEleven++
		} else if points+1 < 22 {
			points++
		} else if acesAsEleven > 0 {
			points -= 10 // Change the given points by 1 ace from 11 to 1
			points++     // Add the current ace as 1 point
			acesAsEleven--
		} else {
			points++
		}
	}

	if player == 0 {
		playersPoints = points
	} else if player == 1 {
		splitPoints = points
	} else {
		dealersPoints = points
	}
}

// Draws the next card from the deck and adds it a hand, according to the "player" value (as explained on the previous function)
func Draw(player int) {

	var points int
	if player == 0 {
		points = playersPoints
	} else if player == 1 {
		points = splitPoints
	} else {
		points = dealersPoints
	}

	// Calculate points, by adding the next card
	if deck[index].value == 1 && points+11 < 22 {
		points += 11
	} else if deck[index].value == 1 {
		points++
	} else if deck[index].value > 9 {
		points += 10
	} else {
		points += deck[index].value
	}

	// If the calculated points are over 21, then recalculate the whole hand
	if player == 0 {
		playersHand = append(playersHand, index)
		playersPoints = points
		if playersPoints > 22 {
			CalculatePoints(0, playersHand)
		}
	} else if player == 1 {
		splitHand = append(splitHand, index)
		splitPoints = points
		if splitPoints > 22 {
			CalculatePoints(1, splitHand)
		}
	} else {
		dealersHand = append(dealersHand, index)
		dealersPoints = points
		if dealersPoints > 22 {
			CalculatePoints(2, dealersHand)
		}
	}

	index++
}

// The player draws 1 card
func Hit(splitHandTurn bool, table *Table, labels *Labels, w fyne.Window) {
	// Draw a card to the hand that is being played at the moment
	if !splitHandTurn {
		Draw(0)
		table.pHand += ", " + ShowCard(index-1)
		table.pPointsStr = strconv.Itoa(playersPoints)
		labels.pCards.Text = table.pHand
		labels.pCards.Refresh()
		labels.pPoints.Text = table.pPointsStr
		labels.pPoints.Refresh()
	} else {
		Draw(1)
		table.sHand += ", " + ShowCard(index-1)
		table.sPointsStr = strconv.Itoa(splitPoints)
		labels.sCards.Text = table.sHand
		labels.sCards.Refresh()
		labels.sPoints.Text = table.sPointsStr
		labels.sPoints.Refresh()
	}
	time.Sleep(1 * time.Second)

	// Having more than 21 points means lose (busted), no matter the dealer's hand
	if (!splitHandTurn && playersPoints > 21) || (splitHandTurn && splitPoints > 21) {
		winner := widget.NewLabel("Busted! Dealer wins!")
		popUP := widget.NewModalPopUp(winner, w.Canvas())
		time.Sleep(4 * time.Second)
		popUP.Hide()
	}
}

// The dealer draws cards until he reaches atleast 17 points to finish the current game
func Stand(labels *Labels, table *Table) {
	// Reveal the 2nd dealer's card
	table.dHand += ", " + ShowCard(3)
	labels.dCards.Text = table.dHand
	labels.dCards.Refresh()
	CalculatePoints(2, dealersHand)
	table.dPointsStr = strconv.Itoa(dealersPoints)
	labels.dPoints.Text = table.dPointsStr
	labels.dPoints.Refresh()

	// Draw cards until the dealer has less than 17 points
	for dealersPoints < 17 {
		time.Sleep(1 * time.Second)
		Draw(2)
		table.dHand += ", " + ShowCard(index-1)
		labels.dCards.Text = table.dHand
		labels.dCards.Refresh()
		table.dPointsStr = strconv.Itoa(dealersPoints)
		labels.dPoints.Text = table.dPointsStr
		labels.dPoints.Refresh()
	}
}

// When button split is clicked - update and show labels for the 2nd hand
func UpdateLabelsOnSplit(table *Table, labels *Labels, splitted bool) {

	// The method has only visual GUI purpose as in the main func it is called first to show the main card and then again to show the 2nd hand
	if !splitted { // The main hand
		table.pHand = ShowHand(playersHand)
		labels.pCards.Text = table.pHand
		labels.pCards.Refresh()
		table.pPointsStr = strconv.Itoa(playersPoints)
		labels.pPoints.Text = table.pPointsStr
		labels.pPoints.Refresh()
	} else { // The 2nd hand
		table.sHand = ShowHand(splitHand)
		labels.sCards.Text = table.sHand
		labels.sCards.Refresh()
		table.sPointsStr = strconv.Itoa(splitPoints)
		labels.sPoints.Text = table.sPointsStr
		labels.sPoints.Refresh()
	}
}

// Determines the winner and the earnings from the current game
func GameResult(pp int, ph []int, insurance, doubled bool, labels *Labels, w fyne.Window) {
	time.Sleep(1 * time.Second)
	var winner *widget.Label
	var currBet int
	var playerBJ bool = HasBlackJack(ph)
	var dealerBJ bool = HasBlackJack(dealersHand)

	// If the dealer has BJ - pay insurance bet if placed
	if dealerBJ && insurance {
		// If the player doubles that should not affect the original insurance bet
		var originalBet int = bet
		if doubled {
			originalBet /= 2
		}

		// Show to the user
		winner = widget.NewLabel("+ " + strconv.Itoa(originalBet) + " from insurance")
		popUP := widget.NewModalPopUp(winner, w.Canvas())
		time.Sleep(2 * time.Second)
		popUP.Hide()

		// Update money
		money += originalBet
		labels.balance.Text = strconv.Itoa(money)
		labels.balance.Refresh()
	}

	if doubled {
		currBet = 2 * bet
	} else {
		currBet = bet
	}

	// If the player has Black Jack and the dealer doesn't - the player wins a total of 1.5 times his original bet
	if playerBJ && !dealerBJ {
		winner = widget.NewLabel("Player wins!\n +" + strconv.Itoa(2*currBet+currBet/2))
		money += 2*currBet + currBet/2
	} else if dealersPoints > 21 || dealersPoints < pp {
		winner = widget.NewLabel("Player wins!\n +" + strconv.Itoa(2*currBet))
		money += 2 * currBet
	} else if dealersPoints == pp && playerBJ == dealerBJ {
		winner = widget.NewLabel("Tie!")
		money += currBet
	} else {
		winner = widget.NewLabel("Dealer wins!")
	}

	// Show the winner to the user
	popUP := widget.NewModalPopUp(winner, w.Canvas())
	time.Sleep(4 * time.Second)
	popUP.Hide()

	// Update the money
	labels.balance.Text = strconv.Itoa(money)
	labels.balance.Refresh()
}

// Reset the game (draws new cards and resets the labes)
func NewGame(table *Table) {
	Shuffle() // Shuffle the deck

	// The player receives 2 cards at the beginning (the first and the third card from the deck)
	playersHand = make([]int, 2, 11)
	playersHand[0] = 0
	playersHand[1] = 2

	// Update the player's labels
	playersPoints = 0
	table.pHand = ShowHand(playersHand)
	CalculatePoints(0, playersHand)
	table.pPointsStr = strconv.Itoa(playersPoints)

	// The dealer also should get 2 cards and show the first
	// First draw 1 card
	dealersHand = make([]int, 1, 11)
	dealersHand[0] = 1

	// Update the dealer's labels and show the first card
	dealersPoints = 0
	table.dHand = ShowHand(dealersHand)
	CalculatePoints(2, dealersHand)
	table.dPointsStr = strconv.Itoa(dealersPoints)
	dealersHand = append(dealersHand, 3) // The dealer receives the 2nd card (4th in the deck)

	splitPoints = 0
	table.sHand = ""
	table.sPointsStr = ""
}

// Sets all labels' (positions, colours, texts) at the beginning
func SetLabels(table Table, labels *Labels) {
	// Balance & bet amount
	labels.balanceLabel.Text = "Balance: "
	labels.balanceLabel.TextSize = 17
	labels.balanceLabel.Color = color.White
	labels.balanceLabel.Move(fyne.NewPos(0, 5))

	labels.balance.Text = strconv.Itoa(money)
	labels.balance.TextSize = 17
	labels.balance.Color = color.CMYK{0, 100, 100, 0}
	labels.balance.Move(fyne.NewPos(80, labels.balanceLabel.Position().Y))

	labels.betLabel.Text = "Current bet: "
	labels.betLabel.TextSize = 17
	labels.betLabel.Color = color.White
	labels.betLabel.Move(fyne.NewPos(0, labels.balanceLabel.Position().Y+20))

	labels.currBet.Text = strconv.Itoa(bet)
	labels.currBet.TextSize = 17
	labels.currBet.Color = color.CMYK{0, 100, 100, 0}
	labels.currBet.Move(fyne.NewPos(labels.betLabel.Position().X+110, labels.betLabel.Position().Y))

	// DEALER'S LABELS:
	labels.dLabel.Text = "Dealer's hand: "
	labels.dLabel.TextSize = 16
	labels.dLabel.Color = color.White
	labels.dLabel.TextStyle.Bold = true
	labels.dLabel.Move(fyne.NewPos(0, labels.currBet.Position().Y+70))

	labels.dCards.Text = table.dHand
	labels.dCards.TextSize = 15
	labels.dCards.Color = color.CMYK{0, 17, 100, 0}
	labels.dCards.Move(fyne.NewPos(labels.dLabel.Position().X+140, labels.dLabel.Position().Y))

	labels.PntLab2.Text = "Points: "
	labels.PntLab2.TextSize = 16
	labels.PntLab2.Color = color.White
	labels.PntLab2.TextStyle.Bold = true
	labels.PntLab2.Move(fyne.NewPos(0, labels.dLabel.Position().Y+20))

	labels.dPoints.Text = table.dPointsStr
	labels.dPoints.TextSize = 15
	labels.dPoints.Color = color.CMYK{0, 17, 100, 0}
	labels.dPoints.Move(fyne.NewPos(labels.PntLab2.Position().X+70, labels.PntLab2.Position().Y))

	// PLAYER'S LABELS:
	labels.pLabel.Text = "Player's hand: "
	labels.pLabel.TextSize = 16
	labels.pLabel.Color = color.White
	labels.pLabel.TextStyle.Bold = true
	labels.pLabel.Move(fyne.NewPos(0, labels.dPoints.Position().Y+40))

	labels.pCards.Text = table.pHand
	labels.pCards.TextSize = 15
	labels.pCards.Color = color.CMYK{0, 17, 100, 0}
	labels.pCards.Move(fyne.NewPos(labels.pLabel.Position().X+140, labels.pLabel.Position().Y))

	labels.PntLab1.Text = "Points: "
	labels.PntLab1.TextSize = 16
	labels.PntLab1.Color = color.White
	labels.PntLab1.TextStyle.Bold = true
	labels.PntLab1.Move(fyne.NewPos(0, labels.pLabel.Position().Y+20))

	labels.pPoints.Text = table.pPointsStr
	labels.pPoints.TextSize = 15
	labels.pPoints.Color = color.CMYK{0, 17, 100, 0}
	labels.pPoints.Move(fyne.NewPos(labels.PntLab1.Position().X+70, labels.PntLab1.Position().Y))

	// SPLITTED HAND:
	labels.sLabel.Text = "Splitted hand: "
	labels.sLabel.TextSize = 16
	labels.sLabel.Color = color.White
	labels.sLabel.TextStyle.Bold = true
	labels.sLabel.Move(fyne.NewPos(0, labels.pPoints.Position().Y+40))
	labels.sLabel.Hide()

	labels.sCards.Text = ""
	labels.sCards.TextSize = 15
	labels.sCards.Color = color.CMYK{0, 17, 100, 0}
	labels.sCards.Move(fyne.NewPos(labels.sLabel.Position().X+140, labels.sLabel.Position().Y))
	labels.sCards.Hide()

	labels.PntLab3.Text = "Points: "
	labels.PntLab3.TextSize = 16
	labels.PntLab3.Color = color.White
	labels.PntLab3.TextStyle.Bold = true
	labels.PntLab3.Move(fyne.NewPos(0, labels.sLabel.Position().Y+20))
	labels.PntLab3.Hide()

	labels.sPoints.Text = "0"
	labels.sPoints.TextSize = 15
	labels.sPoints.Color = color.CMYK{0, 17, 100, 0}
	labels.sPoints.Move(fyne.NewPos(labels.PntLab3.Position().X+70, labels.PntLab3.Position().Y))
	labels.sPoints.Hide()
}

// Enables/Disables betting buttons according to the current bet and balance
func UpdateBetButtons(labels *Labels, incrBetBtn, decrBetBtn *widget.Button) {

	// If the player losses as much money as they remain less than his bet - automatically change his bet to the current amount of money
	if bet > money {
		bet = money - (money % minBet) // The bet should be a number divisible by the minimum bet
		labels.currBet.Text = strconv.Itoa(bet)
		labels.currBet.Refresh()
	}

	// If the current bet can be increased - enable the increase bet button
	if money >= bet+minBet {
		incrBetBtn.Enable()
	}

	// If the current bet can be decreased - enable the decrease bet button
	if bet >= minBet*2 {
		decrBetBtn.Enable()
	}
}

func main() {

	var labels Labels
	var table Table
	var hitBtn, standBtn, newGameBtn, doubleBtn, splitBtn, incrBetBtn, decrBetBtn, insuranceBtn *widget.Button
	var doubled bool = false
	var insurance bool = false
	var splitted bool = false
	var splitHandTurn bool = false

	// Create the main window
	a := app.New()
	a.Settings().SetTheme(theme.DarkTheme())
	w := a.NewWindow("Black Jack")
	w.Resize(fyne.NewSize(900, 400))
	w.SetFixedSize(true)

	// Set the playing deck
	SetDeck()

	SetLabels(table, &labels)

	// Hit -> The player draws 1 card
	hitBtn = widget.NewButton("Hit", func() {
		splitBtn.Disable()
		hitBtn.Disable()
		doubleBtn.Disable()
		insuranceBtn.Disable()

		// Perform the action
		Hit(splitHandTurn, &table, &labels, w)
		hitBtn.Enable()

		// If playing on splitted hand and atleast 1 wasn't busted go to the dealer's turn
		if (playersPoints > 21 && splitted && !splitHandTurn) || (splitHandTurn && splitPoints > 21 && playersPoints < 22) {
			standBtn.OnTapped()
			return
		}

		// If the player is busted or both splitted hands are busted - disable the buttons
		if playersPoints > 21 && (!splitted || (splitHandTurn && splitPoints > 21)) {
			standBtn.Disable()
			hitBtn.Disable()

			// End the game if there aren't enough money
			if money < minBet {
				winner := widget.NewLabel("Out of money!\nBetter luck next time!")
				popUP := widget.NewModalPopUp(winner, w.Canvas())
				time.Sleep(4 * time.Second)
				popUP.Hide()
			} else {
				newGameBtn.Enable()
				UpdateBetButtons(&labels, incrBetBtn, decrBetBtn)
			}
		}
	})
	hitBtn.Resize(fyne.NewSize(90, 50))
	hitBtn.Move(fyne.NewPos(60, labels.sPoints.Position().Y+60))

	// Stand -> The player finishes his turn and the games goes to the dealer's turn
	standBtn = widget.NewButton("Stand", func() {
		splitBtn.Disable()
		doubleBtn.Disable()
		insuranceBtn.Disable()

		// If playing with splitted hand give turn to the 2nd hand if not already done
		if splitted && !splitHandTurn {
			splitHandTurn = true
			labels.pLabel.Color = color.White
			labels.pLabel.Refresh()
			labels.sLabel.Color = color.CMYK{100, 0, 100, 0}
			labels.sLabel.Refresh()
			return
		}

		standBtn.Disable()
		hitBtn.Disable()
		time.Sleep(1 * time.Second)

		// Play the dealer's turn
		Stand(&labels, &table)

		// Determine winner. The case where the player has more than 21 points is caught in hitBtn
		if playersPoints < 22 {
			GameResult(playersPoints, playersHand, insurance, doubled, &labels, w)
		}
		if splitted && splitPoints < 22 {
			GameResult(splitPoints, splitHand, false, doubled, &labels, w)
		}

		// End the game if there aren't enough money
		if money < minBet {
			winner := widget.NewLabel("Out of money!\nBetter luck next time!")
			popUP := widget.NewModalPopUp(winner, w.Canvas())
			time.Sleep(4 * time.Second)
			popUP.Hide()
		} else {
			newGameBtn.Enable()
			UpdateBetButtons(&labels, incrBetBtn, decrBetBtn)
		}
	})
	standBtn.Resize(fyne.NewSize(90, 50))
	standBtn.Move(fyne.NewPos(hitBtn.Position().X+100, hitBtn.Position().Y))

	// Double -> Double the current bet and draw ONLY 1 card. Then it's the dealer's turn
	doubleBtn = widget.NewButton("Double", func() {
		doubleBtn.Disable()
		hitBtn.Disable()
		standBtn.Disable()
		splitBtn.Disable()
		insuranceBtn.Disable()

		// Double the bet
		money -= bet
		labels.balance.Text = strconv.Itoa(money)
		labels.balance.Refresh()
		doubled = true

		time.Sleep(1 * time.Second)
		hitBtn.OnTapped() // Draw 1 card

		// If not busted go to the dealer's turn
		if playersPoints < 22 {
			standBtn.OnTapped()
		}

		// Update the bet buttons
		if money >= minBet {
			UpdateBetButtons(&labels, incrBetBtn, decrBetBtn)
		}
	})
	doubleBtn.Resize(fyne.NewSize(90, 50))
	doubleBtn.Move(fyne.NewPos(standBtn.Position().X+100, standBtn.Position().Y))

	// Split -> The player splits his current hand on two and automatically receives 1 card to each hand. Then he plays with both hands
	// Also the current bet is be placed again but for the 2nd hand
	splitBtn = widget.NewButton("Split", func() {
		hitBtn.Disable()
		standBtn.Disable()
		doubleBtn.Disable()
		splitBtn.Disable()
		insuranceBtn.Disable()
		splitted = true

		// Place the same bet for the 2nd hand
		money -= bet
		labels.balance.Text = strconv.Itoa(money)
		labels.balance.Refresh()

		// Move the 2nd card of the first hand to the second hand
		splitHand = make([]int, 1)
		splitHand[0] = playersHand[1]
		CalculatePoints(1, splitHand)
		playersHand = playersHand[:len(playersHand)-1]
		CalculatePoints(0, playersHand)

		UpdateLabelsOnSplit(&table, &labels, false)
		UpdateLabelsOnSplit(&table, &labels, true)
		labels.sLabel.Color = color.White
		labels.sLabel.Refresh()

		// Show the hidden labels for the 2nd hand
		labels.sLabel.Show()
		labels.sCards.Show()
		labels.sPoints.Show()
		labels.PntLab3.Show()

		time.Sleep(1 * time.Second)

		Draw(0)                                     // Draw 1 card for the main hand
		UpdateLabelsOnSplit(&table, &labels, false) // Update the labels
		// Check if the player has BlackJack on the main hand
		var firstHandBJ bool = HasBlackJack(playersHand)
		if firstHandBJ {
			bj := widget.NewLabel("Black Jack!")
			popUP := widget.NewModalPopUp(bj, w.Canvas())
			time.Sleep(2 * time.Second)
			popUP.Hide()

			splitHandTurn = true                             // If yes, the game options will be executed for the 2nd hand
			labels.sLabel.Color = color.CMYK{100, 0, 100, 0} // Change the colour on the 2nd hand's label
			labels.sLabel.Refresh()
		} else { // If there isn't BlackJack then colour the labels for the 1st hand to show who plays the current turn
			labels.pLabel.Color = color.CMYK{100, 0, 100, 0}
			labels.pLabel.Refresh()
		}

		time.Sleep(1 * time.Second)

		Draw(1)                                    // Draw 1 card for the 2nd hand
		UpdateLabelsOnSplit(&table, &labels, true) // Update it's labels
		if HasBlackJack(splitHand) {               // Check if the 2nd hand has BlackJack
			bj := widget.NewLabel("Black Jack!")
			popUP := widget.NewModalPopUp(bj, w.Canvas())
			time.Sleep(2 * time.Second)
			popUP.Hide()

			// If both hands have Black Jack go to the dealer's turn
			if firstHandBJ {
				time.Sleep(1 * time.Second)
				standBtn.OnTapped()
				return
			}
		}

		time.Sleep(1 * time.Second)
		hitBtn.Enable()
		standBtn.Enable()
	})
	splitBtn.Resize(fyne.NewSize(90, 50))
	splitBtn.Move(fyne.NewPos(doubleBtn.Position().X+100, standBtn.Position().Y))

	// Insurance -> if the dealer's first card is an Ace then the player can place half of the current bet.
	// In the end if the dealer has Black Jack then the player wins 2:1 the insurance bet
	insuranceBtn = widget.NewButton("Insurance", func() {
		insuranceBtn.Disable()

		// Show to the user that insurance bet has been placed
		messageBox := widget.NewLabel("Insurance bet of " + strconv.Itoa(bet/2) + " is placed!")
		popUP := widget.NewModalPopUp(messageBox, w.Canvas())
		time.Sleep(2 * time.Second)
		popUP.Hide()

		insurance = true
		money -= bet / 2
		labels.balance.Text = strconv.Itoa(money)
		labels.balance.Refresh()

		if money < bet {
			doubleBtn.Disable()
		}
	})
	insuranceBtn.Resize(fyne.NewSize(100, 50))
	insuranceBtn.Move(fyne.NewPos(splitBtn.Position().X+100, splitBtn.Position().Y))

	// Increase Bet -> Increases the current bet with the minimum bet
	incrBetBtn = widget.NewButton("Increase bet", func() {
		incrBetBtn.Disable() // Disabling the button at the beginning to handle the current action
		bet += minBet
		labels.currBet.Text = strconv.Itoa(bet)
		labels.currBet.Refresh()

		// If enough money enable the button again
		if money >= bet+minBet {
			incrBetBtn.Enable()
		}

		// If the bet has been increased then for sure it can be decreased atleast once
		decrBetBtn.Enable()
	})
	incrBetBtn.Resize(fyne.NewSize(150, 40))
	incrBetBtn.Move(fyne.NewPos(labels.balance.Position().X+190, 5))

	// Decrease Bet -> Decreases the current bet with the minimum bet
	decrBetBtn = widget.NewButton("Decrease bet", func() {
		decrBetBtn.Disable()
		bet -= minBet
		labels.currBet.Text = strconv.Itoa(bet)
		labels.currBet.Refresh()

		// If the bet is high enough then enable the decrease button
		if bet >= 2*minBet {
			decrBetBtn.Enable()
		}

		// If the bet has been decreased then for sure it can be increased atleast once
		incrBetBtn.Enable()
	})
	decrBetBtn.Resize(fyne.NewSize(150, 40))
	decrBetBtn.Move(fyne.NewPos(incrBetBtn.Position().X+160, 5))

	// New Game -> Places the current bet and begins the next game
	newGameBtn = widget.NewButton("Next game", func() {
		newGameBtn.Disable()
		incrBetBtn.Disable()
		decrBetBtn.Disable()

		NewGame(&table)
		index = 4
		doubled = false
		insurance = false

		// Refresh all labels with the new values
		money -= bet // Entrance bet
		labels.balance.Text = strconv.Itoa(money)
		labels.balance.Refresh()

		labels.pCards.Text = table.pHand
		labels.pCards.Refresh()
		labels.pPoints.Text = table.pPointsStr
		labels.pPoints.Refresh()

		labels.dCards.Text = table.dHand
		labels.dCards.Refresh()
		labels.dPoints.Text = table.dPointsStr
		labels.dPoints.Refresh()

		time.Sleep(1 * time.Second)

		// If the player has Black Jack then go directly to the dealer's turn
		if HasBlackJack(playersHand) {
			bj := widget.NewLabel("Black Jack!")
			popUP := widget.NewModalPopUp(bj, w.Canvas())
			time.Sleep(2 * time.Second)
			popUP.Hide()
			standBtn.OnTapped()
			return
		}

		// If the previous game was played on splitted hand - hide the labels from the 2nd hand
		if splitted {
			splitted = false
			splitHandTurn = false
			labels.sLabel.Color = color.White
			labels.sLabel.Hide()
			labels.sCards.Hide()
			labels.sPoints.Hide()
			labels.PntLab3.Hide()
		}

		hitBtn.Enable()
		standBtn.Enable()

		// If the player has enough money - enable Double option
		if money >= bet {
			doubleBtn.Enable()
		}

		// If the player has 2 same cards (except colour) and enough money - enable the split option
		if deck[0].value == deck[2].value && money >= bet {
			splitBtn.Enable()
		}

		// If the shown card of the dealer is an ace and the player has enough money - enable insurance option
		if deck[dealersHand[0]].value == 1 && money > bet/2 {
			insuranceBtn.Enable()
		}
	})
	newGameBtn.Resize(fyne.NewSize(100, 50))
	newGameBtn.Move(fyne.NewPos(insuranceBtn.Position().X+108, insuranceBtn.Position().Y))

	// Quit -> Exists the program
	quitBtn := widget.NewButton("Quit", a.Quit)
	quitBtn.Resize(fyne.NewSize(90, 50))
	quitBtn.Move(fyne.NewPos(newGameBtn.Position().X+170, newGameBtn.Position().Y))

	// When the app is started for a first time dissable the buttons below. Only New Game, Increase bet and Quit should be enabled
	hitBtn.Disable()
	standBtn.Disable()
	splitBtn.Disable()
	doubleBtn.Disable()
	insuranceBtn.Disable()
	decrBetBtn.Disable()

	// Defines the container with the GUI objects
	f := fyne.NewContainer(
		&labels.balanceLabel, &labels.balance, &labels.betLabel, &labels.currBet,
		&labels.dLabel, &labels.pLabel, &labels.sLabel,
		&labels.dCards, &labels.pCards, &labels.sCards,
		&labels.dPoints, &labels.pPoints, &labels.sPoints,
		&labels.PntLab1, &labels.PntLab2, &labels.PntLab3,
		hitBtn, splitBtn, doubleBtn, newGameBtn, standBtn, insuranceBtn, incrBetBtn, decrBetBtn, quitBtn)

	w.SetContent(f)
	w.ShowAndRun()
}
