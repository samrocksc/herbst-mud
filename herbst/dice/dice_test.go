package dice

import "testing"

func TestRoll(t *testing.T) {
	for range 100 {
		roll, total := Roll(6, 1, 0)
		if roll < 1 || roll > 6 {
			t.Errorf("Roll(6, 1, 0) = %d, want 1-6", roll)
		}
		if total != roll {
			t.Errorf("total should equal roll when modifier is 0")
		}
	}
}

func TestRollWithModifier(t *testing.T) {
	roll, total := Roll(6, 1, 5)
	if total != roll+5 {
		t.Errorf("Roll(6, 1, 5) total = %d, want %d", total, roll+5)
	}
}

func TestRollWithCrit(t *testing.T) {
	// Test many times to hit both crit and fumble
	critCount := 0
	fumbleCount := 0
	for range 1000 {
		roll, total, isCrit, isFumble := RollWithCrit(0)
		if roll < 1 || roll > 20 {
			t.Errorf("RollWithCrit roll = %d, want 1-20", roll)
		}
		if total != roll {
			t.Errorf("total should equal roll when modifier is 0")
		}
		if isCrit && roll != 20 {
			t.Errorf("isCrit should only be true when roll is 20")
		}
		if isFumble && roll != 1 {
			t.Errorf("isFumble should only be true when roll is 1")
		}
		if isCrit && isFumble {
			t.Errorf("cannot be both crit and fumble")
		}
		if isCrit {
			critCount++
		}
		if isFumble {
			fumbleCount++
		}
	}
	// Should hit at least a few crits and fumbles in 1000 rolls
	if critCount == 0 {
		t.Errorf("expected at least one critical hit in 1000 rolls")
	}
	if fumbleCount == 0 {
		t.Errorf("expected at least one fumble in 1000 rolls")
	}
}

func TestD20(t *testing.T) {
	roll, total := D20(5)
	if roll < 1 || roll > 20 {
		t.Errorf("D20(5) roll = %d, want 1-20", roll)
	}
	if total != roll+5 {
		t.Errorf("D20(5) total = %d, want %d", total, roll+5)
	}
}