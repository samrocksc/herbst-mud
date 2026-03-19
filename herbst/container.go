package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Container represents an item that can hold other items
type Container struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Capacity    int      `json:"capacity"`
	Items       []string `json:"items"`
	IsOpen      bool     `json:"is_open"`
	IsLocked    bool     `json:"is_locked"`
	KeyID       string   `json:"key_id,omitempty"`
}

// IsContainer checks if an item name suggests it's a container
func IsContainer(name string) bool {
	containerKeywords := []string{"chest", "crate", "box", "bag", "sack", "barrel", "trunk", "coffer", "casket", "locker"}
	nameLower := strings.ToLower(name)
	for _, kw := range containerKeywords {
		if strings.Contains(nameLower, kw) {
			return true
		}
	}
	return false
}

// NewContainer creates a new container with default settings
func NewContainer(id, name, description string) *Container {
	return &Container{
		ID:          id,
		Name:        name,
		Description: description,
		Capacity:    10,
		Items:       []string{},
		IsOpen:      true,
		IsLocked:    false,
		KeyID:       "",
	}
}

// AddItem adds an item to the container
func (c *Container) AddItem(itemID string) error {
	if len(c.Items) >= c.Capacity {
		return fmt.Errorf("container is full")
	}
	c.Items = append(c.Items, itemID)
	return nil
}

// RemoveItem removes an item from the container
func (c *Container) RemoveItem(itemID string) bool {
	for i, item := range c.Items {
		if item == itemID {
			c.Items = append(c.Items[:i], c.Items[i+1:]...)
			return true
		}
	}
	return false
}

// GetContents returns the items in the container
func (c *Container) GetContents() []string {
	return c.Items
}

// Open opens the container
func (c *Container) Open() error {
	if c.IsLocked {
		return fmt.Errorf("the %s is locked", c.Name)
	}
	c.IsOpen = true
	return nil
}

// Close closes the container
func (c *Container) Close() {
	c.IsOpen = false
}

// Lock locks the container with a key
func (c *Container) Lock(keyID string) {
	c.IsLocked = true
	c.KeyID = keyID
	c.IsOpen = false
}

// Unlock unlocks the container with a key
func (c *Container) Unlock(keyID string) error {
	if c.KeyID != "" && c.KeyID != keyID {
		return fmt.Errorf("you don't have the right key")
	}
	c.IsLocked = false
	return nil
}

// FormatContents returns a formatted string of container contents
func (c *Container) FormatContents(itemDetails map[string]string) string {
	if !c.IsOpen {
		return fmt.Sprintf("The %s is closed.", c.Name)
	}

	if len(c.Items) == 0 {
		return fmt.Sprintf("The %s is empty.", c.Name)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("The %s contains:\n", c.Name))

	for i, itemID := range c.Items {
		detail := itemDetails[itemID]
		if detail == "" {
			detail = itemID
		}
		sb.WriteString(fmt.Sprintf("  [%d] %s\n", i+1, detail))
	}

	return sb.String()
}

// ContainerFromJSON creates a Container from JSON string (for DB storage)
func ContainerFromJSON(data string) (*Container, error) {
	if data == "" {
		return nil, nil
	}
	var c Container
	err := json.Unmarshal([]byte(data), &c)
	return &c, err
}

// ContainerToJSON converts a Container to JSON string
func ContainerToJSON(c *Container) (string, error) {
	if c == nil {
		return "", nil
	}
	data, err := json.Marshal(c)
	return string(data), err
}

// FindItemInContainer finds an item in container by name (partial match)
func (c *Container) FindItemInContainer(name string) string {
	nameLower := strings.ToLower(name)
	for _, itemID := range c.Items {
		itemLower := strings.ToLower(itemID)
		if strings.Contains(itemLower, nameLower) {
			return itemID
		}
	}
	return ""
}