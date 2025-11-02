package rooms

import (
	"github.com/sam/makeathing/internal/items"
	"github.com/sam/makeathing/internal/characters"
)

// Room represents a room in the game
type Room struct {
	ID          string
	Description string
	Exits       map[Direction]string // Direction to Room ID
	ImmovableObjects []items.Item
	MovableObjects   []items.Item
	Smells           string
	NPCs             []characters.Character
}

// Direction represents cardinal directions
type Direction string

const (
	North     Direction = "north"
	South     Direction = "south"
	East      Direction = "east"
	West      Direction = "west"
	Northeast Direction = "northeast"
	Northwest Direction = "northwest"
	Southeast Direction = "southeast"
	Southwest Direction = "southwest"
	Up        Direction = "up"
	Down      Direction = "down"
)