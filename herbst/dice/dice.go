package dice

import (
	"math/rand"
	"time"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// Roll performs a dice roll and returns the raw roll and total with modifier
func Roll(sides, count, modifier int) (roll int, total int) {
	roll = 0
	for range count {
		roll += rng.Intn(sides) + 1
	}
	total = roll + modifier
	return
}

// D20 is a convenience function for d20 rolls
func D20(modifier int) (roll, total int) {
	roll = rng.Intn(20) + 1
	total = roll + modifier
	return
}

// RollWithCrit returns roll result with critical hit/fumble detection
func RollWithCrit(modifier int) (roll int, total int, isCrit bool, isFumble bool) {
	roll = rng.Intn(20) + 1
	total = roll + modifier
	isCrit = (roll == 20)
	isFumble = (roll == 1)
	return
}

// RollDamage rolls damage dice and returns the result with a description
func RollDamage(sides, count, modifier int) (total int, rollStr string) {
	roll := 0
	for range count {
		roll += rng.Intn(sides) + 1
	}
	total = roll + modifier

	// Build description string
	if modifier == 0 {
		rollStr = ""
	} else if modifier > 0 {
		rollStr = ""
	} else {
		rollStr = ""
	}
	return
}