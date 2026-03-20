package main

// ============================================================
// DATA TYPES - All model struct types
// ============================================================

// RoomItem represents an item in a room for display
type RoomItem struct {
	ID              int            `json:"id"`
	Name            string         `json:"name"`
	Description     string         `json:"description"`
	ExamineDesc     string         `json:"examineDesc"`
	HiddenDetails   []HiddenDetail `json:"hiddenDetails"`
	HiddenThreshold int            `json:"hiddenThreshold"`
	IsImmovable     bool           `json:"isImmovable"`
	Color           string         `json:"color"`
	IsVisible       bool           `json:"isVisible"`
	ItemType        string         `json:"itemType"`
	Weight          int            `json:"weight"`
	ItemDamage      int            `json:"itemDamage"`
	ItemDurability  int            `json:"itemDurability"`
	RevealCondition map[string]any `json:"revealCondition"`
}

// HiddenDetail represents hidden information about an item
type HiddenDetail struct {
	Text      string `json:"text"`
	Skill     string `json:"skill"`
	Threshold int    `json:"threshold"`
}

// roomCharacter represents a character (NPC or player) in a room for display
type roomCharacter struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	IsNPC    bool   `json:"isNPC"`
	Level    int    `json:"level"`
	Class    string `json:"class"`
	Race     string `json:"race"`
	UserID   int    `json:"userId"`
}

// inventoryItem represents an item in the player's inventory
type inventoryItem struct {
	ID          int
	Name        string
	Description string
	ItemType    string
	IsEquipped  bool
	Rarity      string
}
