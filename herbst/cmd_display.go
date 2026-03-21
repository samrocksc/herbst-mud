package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

// displayItemDetails shows detailed info about an item
func (m *model) displayItemDetails(item RoomItem) {
	var details strings.Builder

	details.WriteString(fmt.Sprintf("[%s]\n", item.Name))

	desc := item.ExamineDesc
	if desc == "" {
		desc = item.Description
	}
	details.WriteString(desc + "\n")

	if item.ItemType == "weapon" || item.ItemType == "armor" {
		details.WriteString("\n--- Stats ---\n")
		if item.Weight > 0 {
			details.WriteString(fmt.Sprintf("  Weight: %d\n", item.Weight))
		}
		if item.ItemDamage > 0 {
			details.WriteString(fmt.Sprintf("  Damage: %d\n", item.ItemDamage))
		}
		if item.ItemDurability > 0 {
			details.WriteString(fmt.Sprintf("  Durability: %d\n", item.ItemDurability))
		}
		details.WriteString(fmt.Sprintf("  Type: %s\n", item.ItemType))
	}

	if len(item.HiddenDetails) > 0 {
		details.WriteString("\n--- You Notice ---\n")
		for _, hd := range item.HiddenDetails {
			details.WriteString(fmt.Sprintf("  %s\n", hd.Text))
		}
	}

	m.AppendMessage(details.String(), "info")
}

// fuzzyWordMatch returns true if all words in target appear as substrings in name.
// "grand man" matches "Grand Ol' Man". Case-insensitive.
func fuzzyWordMatch(name, target string) bool {
	nameLower := strings.ToLower(name)
	for _, word := range strings.Fields(strings.ToLower(target)) {
		if !strings.Contains(nameLower, word) {
			return false
		}
	}
	return true
}

// httpGet is a helper for making GET requests
func httpGet(url string) (*http.Response, error) {
	return http.Get(url)
}

// httpPost is a helper for making POST requests with a body
func httpPost(url, body string) (*http.Response, error) {
	return http.Post(url, "application/json", strings.NewReader(body))
}

// ioReadAll is exported for use by other files
func ioReadAll(r io.Reader) ([]byte, error) {
	return io.ReadAll(r)
}