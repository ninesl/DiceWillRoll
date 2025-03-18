package main

import (
	"slices"
)

// This pkg is used to determine hand outcome from die.Roll().Value()

// HandRank is the type of hand played, Straight, Full House, Five of a kind, etc.
//
// Use .String() to get the name of the hand
type HandRank uint8

// Each HandRank is based on the actual value associated with it.
//
//	TWO_PAIR > NO_HAND
//	STRAIGHT_LARGE < FOUR_OF_A_KIND
const (
	NO_HAND HandRank = iota
	HIGH_DIE
	ONE_PAIR
	SNAKE_EYES
	TWO_PAIR
	THREE_OF_A_KIND
	// 4 consecutive
	STRAIGHT_SMALL
	// 5 consecutive
	STRAIGHT_LARGE
	FULL_HOUSE
	FOUR_OF_A_KIND
	FIVE_OF_A_KIND

	// 6 die
	// straight + pair?

	// 2 + 2 + 2
	THREE_PAIR
	// 4 + 2
	CROWDED_HOUSE
	SIX_OF_A_KIND
	// 6 consecutive
	STRAIGHT_LARGER
	// 3 + 3
	TWO_THREE_OF_A_KIND

	// 7 die
	// straight + three of a kind?

	// 4 + 3
	OVERPOPULATED_HOUSE
	// 7 consecutive
	STRAIGHT_LARGEST
	// 5 + 2
	FULLEST_HOUSE
	SEVEN_OF_A_KIND
	// 7 of a kind where all Value is 7
	SEVEN_SEVENS

	// Special, usually from modifiers
	STRAIGHT_MAX

	// other
	UNKNOWN_HAND
	NUM_TYP
)

// could be modified?
var (
	STRAIGHT_SMALL_LENGTH   = 4
	STRAIGHT_LARGE_LENGTH   = 5
	STRAIGHT_LARGER_LENGTH  = 6
	STRAIGHT_LARGEST_LENGTH = 7

	SNAKE_EYES_TARGET   = 1 // default to 1 bc obvious
	SEVEN_SEVENS_TARGET = 7
)

// first check when determining HandRank
//
// returns NO_HAND if a straight is not found, otherwise returns the associated HandRank for the straight
func checkStraight(values []int) HandRank {
	// consecutive := [MAX_PIPS]bool{}

	// for _, value := range values {
	// 	consecutive[value] = true
	// }

	inARow := 1
	var maxRow int
	straight := NO_HAND
	var lastValue, curValue int
	slices.Sort(values)

	// fmt.Println(values)

	lastValue = values[0]
	for i := 1; i < len(values); i += 1 {
		curValue = values[i]

		// TODO: modifiers for skip straight?

		if lastValue+1 == curValue {
			inARow += 1
		} else {
			maxRow = inARow
			inARow = 1
		}

		lastValue = curValue
	}

	// for value slices that are the entire sequence
	if maxRow == 0 {
		maxRow = inARow
	}

	// fmt.Println(maxRow, "in a row")

	// for i := range consecutive {
	// 	if consecutive[i] {
	// 		curValue = i

	// 		if lastValue+1 == curValue { // still works for 0 - 1
	// 			inARow += 1
	// 		} else if lastValue != curValue {
	// 			inARow = 0
	// 		}

	// 		lastValue = curValue
	// 	}
	// }

	if maxRow >= STRAIGHT_SMALL_LENGTH { // small straight
		switch maxRow {
		case STRAIGHT_SMALL_LENGTH:
			straight = STRAIGHT_SMALL
		case STRAIGHT_LARGE_LENGTH:
			straight = STRAIGHT_LARGE
		case STRAIGHT_LARGER_LENGTH:
			straight = STRAIGHT_LARGER
		case STRAIGHT_LARGEST_LENGTH:
			straight = STRAIGHT_LARGEST
		case STRAIGHT_LARGEST_LENGTH + 1: // for + 1 from modifiers??? FIXME: this shouldn't be hardcoded?
			straight = STRAIGHT_MAX // placeholder TODO: STRAIGHT_MAX impl
		}
	}

	return straight
}

// second check when determining handRank
//
//	// annoying logic for each one. comments are placed where it isn't obv after a glance
//
// returns NO_HAND if nothing is found, otherwise returns the highest HandRank that can be associated (other than straights)
func checkHandOtherThanStraight(valueCount map[int]int, values []int, numDice int) HandRank {
	uniqueValues := len(valueCount)
	switch uniqueValues {
	case 0:
		return NO_HAND
	case 1:
		switch numDice {
		case 1:
			return HIGH_DIE
		case 2:
			if valueCount[SNAKE_EYES_TARGET] == numDice {
				return SNAKE_EYES
			} else {
				return ONE_PAIR
			}
		case 3:
			return THREE_OF_A_KIND
		case 4:
			return FOUR_OF_A_KIND
		case 5:
			return FIVE_OF_A_KIND
		case 6:
			return SIX_OF_A_KIND
		case 7:
			if valueCount[SEVEN_SEVENS_TARGET] == 7 {
				return SEVEN_SEVENS
			} else {
				return SEVEN_OF_A_KIND
			}
		}
	case 2:
		switch numDice {
		case 2: // 4 2
			return HIGH_DIE
		case 3: // 2 2 1
			if valueCount[SNAKE_EYES_TARGET] == 2 {
				return SNAKE_EYES
			} else {
				return ONE_PAIR
			}
		case 4: // 1 1 2 2, 1 1 1 2
			if valueCount[values[0]] == 3 || valueCount[values[1]] == 3 {
				return THREE_OF_A_KIND
			}
			return TWO_PAIR
		case 5: // 1 2 2 2 2, 1 1 2 2 2
			if valueCount[values[0]] == 2 && valueCount[values[1]] == 3 ||
				valueCount[values[1]] == 2 && valueCount[values[0]] == 3 {
				return FULL_HOUSE
			} else {
				return FOUR_OF_A_KIND
			}
		case 6: // 3 3 3 5 5 5, 1 1 2 2 2 2, 4 3 3 3 3 3,
			if valueCount[values[0]] == valueCount[values[1]] {
				return TWO_THREE_OF_A_KIND
			} else if valueCount[values[0]] == 2 && valueCount[values[1]] == 4 ||
				valueCount[values[1]] == 2 && valueCount[values[0]] == 4 {
				return CROWDED_HOUSE
			} else {
				return FIVE_OF_A_KIND
			}
		case 7: // 1 1 1 2 2 2 2, 1 1 2 2 2 2 2, 1 2 2 2 2 2 2

			if valueCount[values[0]] == 3 && valueCount[values[1]] == 4 ||
				valueCount[values[1]] == 3 && valueCount[values[0]] == 4 {
				return OVERPOPULATED_HOUSE
			} else if valueCount[values[0]] == 2 && valueCount[values[1]] == 5 ||
				valueCount[values[1]] == 2 && valueCount[values[0]] == 5 {
				return FULLEST_HOUSE
			} else {
				return SIX_OF_A_KIND
			}
		}
	case 3:
		switch numDice {
		case 3:
			return HIGH_DIE
		case 4: // 1 1 2 3
			if valueCount[SNAKE_EYES_TARGET] == 2 {
				return SNAKE_EYES
			} else {
				return ONE_PAIR
			}
		case 5: // 1 1 1 2 3, 1 1 2 2 3
			if valueCount[values[0]] == 3 ||
				valueCount[values[1]] == 3 ||
				valueCount[values[2]] == 3 {
				return THREE_OF_A_KIND
			} else {
				return TWO_PAIR
			}
		case 6: // 1 1 1 1 2 3, 1 1 1 2 2 3, 1 1 2 2 3 3
			if valueCount[values[0]] == 4 ||
				valueCount[values[1]] == 4 ||
				valueCount[values[2]] == 4 {
				return FOUR_OF_A_KIND
			} else if valueCount[values[0]] == valueCount[values[1]] && valueCount[values[1]] == valueCount[values[2]] {
				return THREE_PAIR
			} else {
				return FULL_HOUSE
			}
		case 7: // 1 1 1 1 1 2 3, 1 1 1 1 2 2 3, 1 1 1 2 2 2 3
			if valueCount[values[0]] == 5 ||
				valueCount[values[1]] == 5 ||
				valueCount[values[2]] == 5 {
				return FIVE_OF_A_KIND
			} else if valueCount[values[0]] == valueCount[values[1]] ||
				valueCount[values[1]] == valueCount[values[2]] {
				return TWO_THREE_OF_A_KIND
			} else if (valueCount[values[0]] == 2 && valueCount[values[1]] == 4 ||
				valueCount[values[1]] == 2 && valueCount[values[0]] == 4) ||

				(valueCount[values[1]] == 2 && valueCount[values[2]] == 4 ||
					valueCount[values[2]] == 2 && valueCount[values[1]] == 4) ||

				(valueCount[values[0]] == 2 && valueCount[values[2]] == 4 ||
					valueCount[values[2]] == 2 && valueCount[values[0]] == 4) {
				return CROWDED_HOUSE
			} else {
				return FULL_HOUSE
			}
		}
	case 4:
		switch numDice {
		case 4:
			return HIGH_DIE
		case 5: // 1 1 2 3 4
			if valueCount[SNAKE_EYES_TARGET] == 2 {
				return SNAKE_EYES
			} else {
				return ONE_PAIR
			}
		case 6: // 1 1 1 2 3 4, 1 1 2 2 3 4
			if valueCount[values[0]] == 3 ||
				valueCount[values[1]] == 3 ||
				valueCount[values[2]] == 3 ||
				valueCount[values[3]] == 3 {
				return THREE_OF_A_KIND
			} else {
				return TWO_PAIR
			}
		case 7: // 1 1 1 1 2 3 4, 1 1 1 2 2 3 4, 1 1 2 2 3 3 4
			if valueCount[values[0]] == 4 ||
				valueCount[values[1]] == 4 ||
				valueCount[values[2]] == 4 ||
				valueCount[values[3]] == 4 {
				return FOUR_OF_A_KIND
			} else if (valueCount[values[0]] == valueCount[values[1]] && valueCount[values[1]] == valueCount[values[2]]) ||
				(valueCount[values[1]] == valueCount[values[2]] && valueCount[values[2]] == valueCount[values[3]]) {
				return THREE_PAIR
			} else {
				return FULL_HOUSE
			}
		}
	case 5:
		switch numDice {
		case 5:
			return HIGH_DIE
		case 6: // 1 1 2 3 4 5
			if valueCount[SNAKE_EYES_TARGET] == 2 {
				return SNAKE_EYES
			} else {
				return ONE_PAIR
			}
		case 7: // 1 1 1 2 3 4 5, 1 1 2 2 3 4 5
			if valueCount[values[0]] == 3 ||
				valueCount[values[1]] == 3 ||
				valueCount[values[2]] == 3 ||
				valueCount[values[3]] == 3 ||
				valueCount[values[4]] == 3 {
				return THREE_OF_A_KIND
			} else {
				return TWO_PAIR
			}
		}
	case 6:
		switch numDice {
		case 6:
			return HIGH_DIE
		case 7:
			if valueCount[SNAKE_EYES_TARGET] == 2 {
				return SNAKE_EYES
			} else {
				return ONE_PAIR
			}
		}
	case 7:
		return HIGH_DIE
	}

	return UNKNOWN_HAND // SHOULD NOT GET HERE!
}

// TODO: unit test
//
// Calculates given slice of dice's active face's values.
//
// returns a HandRank that corresponds to the input dice
func DetermineHandRank(dice []Die) HandRank {
	var (
		numDice       = len(dice)
		valueCount    = map[int]int{}
		values        = []int{} // tracks the # of values' occurences from valueCount for comparisons. A slice of the values of valueCount
		foundStraight HandRank
		handFound     HandRank
	)

	// gather occurences of unique values
	for i := range numDice {
		value := dice[i].ActiveFace().Value()
		// starts from 0 bc of nil int
		valueCount[value] = valueCount[value] + 1
	}

	// used in determining hand logic
	for value := range valueCount {
		values = append(values, value)
	}

	// if a straight exists, assign check variable
	if numDice >= STRAIGHT_SMALL_LENGTH { // has to be at least 4. STRAIGHT_SMALL
		foundStraight = checkStraight(values)
	}

	// determine type of hand based off # of unique values
	handFound = checkHandOtherThanStraight(valueCount, values, numDice)

	// found straight vs handfound. Edge case for full houses and straights in the same hand
	if foundStraight > handFound {
		return foundStraight
	}
	return handFound
}

var (
	handRankStringMap = map[HandRank]string{
		UNKNOWN_HAND:        "UNKNOWN HAND (shouldn't see this)",
		NO_HAND:             "No Hand",
		HIGH_DIE:            "High Die",
		ONE_PAIR:            "One Pair",
		SNAKE_EYES:          "Snake Eyes",
		TWO_PAIR:            "Two Pair",
		THREE_OF_A_KIND:     "Three of a Kind",
		STRAIGHT_SMALL:      "Small Straight",
		STRAIGHT_LARGE:      "Large Straight",
		FULL_HOUSE:          "Full House",
		FOUR_OF_A_KIND:      "Four of a Kind",
		FIVE_OF_A_KIND:      "Five of a Kind",
		THREE_PAIR:          "Three Pair",
		CROWDED_HOUSE:       "Crowded House",
		SIX_OF_A_KIND:       "Six of a Kind",
		STRAIGHT_LARGER:     "Large-r Straight",
		TWO_THREE_OF_A_KIND: "Three's a Crowd",
		OVERPOPULATED_HOUSE: "Overpopulated House",
		STRAIGHT_LARGEST:    "Ultra Straight",
		FULLEST_HOUSE:       "Fire Code Violation",
		SEVEN_OF_A_KIND:     "Seven of a Kind",
		SEVEN_SEVENS:        "Lucky Sevens",
		STRAIGHT_MAX:        "MEGA Straight",
	}
)

func (h *HandRank) String() string {
	return handRankStringMap[*h]
}
