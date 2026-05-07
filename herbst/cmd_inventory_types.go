package main

// inventoryItem holds equipment data for inventory display.
type inventoryItem struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Slot            string `json:"slot"`
	IsEquipped      bool   `json:"isEquipped"`
	ItemType        string `json:"itemType"`
	DamageDiceCount int    `json:"damage_dice_count"`
	DamageDiceSides int    `json:"damage_dice_sides"`
	DamageBonus     int    `json:"damage_bonus"`
	ArmorRating     int    `json:"armor_rating"`
	IsTwoHanded     bool   `json:"is_two_handed"`
	Rarity          string `json:"rarity"`
}

// raceData holds the race info needed for inventory display.
type raceData struct {
	Name           string   `json:"name"`
	DisplayName    string   `json:"display_name"`
	EquipmentSlots []string `json:"equipment_slots"`
}