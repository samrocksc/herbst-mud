package items

// Item represents an item in the game
type Item struct {
	Name        string
	Description string
	Type        ItemType
	Stats       ItemStats
	IsMagical   bool
}

// ItemType represents the type of item
type ItemType string

const (
	Weapon    ItemType = "weapon"
	Wearable  ItemType = "wearable"
	Movable   ItemType = "movable"
	Immovable ItemType = "immovable"
)

// ItemStats represents stats for an item
type ItemStats struct {
	Strength     int
	Intelligence int
	Dexterity    int
}