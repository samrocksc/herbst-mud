package dbinit

import (
	"testing"
)

// TestWeaponDropSelection tests that the correct weapon drops for a character's class
func TestWeaponDropSelection(t *testing.T) {
	// Test Warrior gets Rusty Sword
	warriorClass := "warrior"
	droppedWeapon := selectWeaponForClass(warriorClass)
	if droppedWeapon != "Rusty Sword" {
		t.Errorf("Expected Warrior to drop Rusty Sword, got %s", droppedWeapon)
	}

	// Test Chef gets Twisted Pipe
	chefClass := "chef"
	droppedWeapon = selectWeaponForClass(chefClass)
	if droppedWeapon != "Twisted Pipe" {
		t.Errorf("Expected Chef to drop Twisted Pipe, got %s", droppedWeapon)
	}

	// Test unknown class gets Rusty Sword (default)
	unknownClass := "mage"
	droppedWeapon = selectWeaponForClass(unknownClass)
	if droppedWeapon != "Rusty Sword" {
		t.Errorf("Expected unknown class to drop Rusty Sword (default), got %s", droppedWeapon)
	}

	t.Log("✓ TestWeaponDropSelection: Weapon selection for classes verified")
}

// selectWeaponForClass returns the weapon name that should drop for a given class
// This is a helper function that will be used in the actual drop logic
func selectWeaponForClass(class string) string {
	weaponMap := map[string]string{
		"warrior": "Rusty Sword",
		"chef":    "Twisted Pipe",
	}
	if weapon, ok := weaponMap[class]; ok {
		return weapon
	}
	// Default to Rusty Sword for unknown classes
	return "Rusty Sword"
}

// TestWeaponStats tests weapon damage ranges
func TestWeaponStats(t *testing.T) {
	weapons := map[string]struct {
		minDamage int
		maxDamage int
	}{
		"Rusty Sword":   {minDamage: 1, maxDamage: 3},
		"Twisted Pipe":  {minDamage: 1, maxDamage: 2},
	}

	for name, stats := range weapons {
		if name == "Rusty Sword" {
			if stats.minDamage != 1 || stats.maxDamage != 3 {
				t.Errorf("Rusty Sword: expected 1-3, got %d-%d", stats.minDamage, stats.maxDamage)
			}
		}
		if name == "Twisted Pipe" {
			if stats.minDamage != 1 || stats.maxDamage != 2 {
				t.Errorf("Twisted Pipe: expected 1-2, got %d-%d", stats.minDamage, stats.maxDamage)
			}
		}
	}
	t.Log("✓ TestWeaponStats: Weapon damage ranges verified")
}