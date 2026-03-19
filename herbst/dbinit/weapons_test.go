package dbinit

import (
	"context"
	"testing"

	"herbst/db"
	"herbst/db/equipment"
)

// TestInitWeapons tests that starter weapons can be created
func TestInitWeapons(t *testing.T) {
	ctx := context.Background()

	client := db.NewClient()
	if client == nil {
		t.Skipf("Skipping test: database not available")
	}
	defer client.Close()

	// Clean up any existing weapons before test
	existingWeapons, err := client.Equipment.Query().
		Where(equipment.ItemTypeEQ("weapon")).
		All(ctx)
	if err == nil {
		for _, w := range existingWeapons {
			client.Equipment.DeleteOne(w).Exec(ctx)
		}
	}

	// Run the InitWeapons function
	err = InitWeapons(client)
	if err != nil {
		t.Fatalf("InitWeapons failed: %v", err)
	}

	// Verify weapons were created
	allEquipment, err := client.Equipment.Query().All(ctx)
	if err != nil {
		t.Fatalf("Failed to query equipment: %v", err)
	}

	// Filter for weapon types with guaranteed drop
	var weapons []*db.Equipment
	for _, e := range allEquipment {
		if e.ItemType == "weapon" && e.GuaranteedDrop {
			weapons = append(weapons, e)
		}
	}

	if len(weapons) != 2 {
		t.Errorf("Expected 2 guaranteed drop weapons, got %d", len(weapons))
	}

	// Verify Rusty Sword
	var rustySword *db.Equipment
	for _, e := range allEquipment {
		if e.Name == "Rusty Sword" {
			rustySword = e
			break
		}
	}
	if rustySword == nil {
		t.Fatal("Failed to find Rusty Sword")
	}
	if rustySword.MinDamage != 1 || rustySword.MaxDamage != 3 {
		t.Errorf("Rusty Sword damage: expected 1-3, got %d-%d", rustySword.MinDamage, rustySword.MaxDamage)
	}
	if rustySword.ClassRestriction != "warrior" {
		t.Errorf("Rusty Sword class restriction: expected warrior, got %s", rustySword.ClassRestriction)
	}
	if rustySword.WeaponType != "sword" {
		t.Errorf("Rusty Sword weapon type: expected sword, got %s", rustySword.WeaponType)
	}

	// Verify Twisted Pipe
	var twistedPipe *db.Equipment
	for _, e := range allEquipment {
		if e.Name == "Twisted Pipe" {
			twistedPipe = e
			break
		}
	}
	if twistedPipe == nil {
		t.Fatal("Failed to find Twisted Pipe")
	}
	if twistedPipe.MinDamage != 1 || twistedPipe.MaxDamage != 2 {
		t.Errorf("Twisted Pipe damage: expected 1-2, got %d-%d", twistedPipe.MinDamage, twistedPipe.MaxDamage)
	}
	if twistedPipe.ClassRestriction != "chef" {
		t.Errorf("Twisted Pipe class restriction: expected chef, got %s", twistedPipe.ClassRestriction)
	}
	if twistedPipe.WeaponType != "pipe" {
		t.Errorf("Twisted Pipe weapon type: expected pipe, got %s", twistedPipe.WeaponType)
	}

	t.Log("✓ TestInitWeapons: All starter weapons verified")
}

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