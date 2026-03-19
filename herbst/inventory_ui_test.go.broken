package main

import (
	"testing"
)

func TestGetItemIcon(t *testing.T) {
	t.Parallel()

	t.Run("weapon returns sword emoji", func(t *testing.T) {
		result := getItemIcon("weapon")
		if result != "⚔️" {
			t.Errorf("expected ⚔️, got %s", result)
		}
	})

	tests.Run("armor returns shield emoji", func(t *testing.T) {
		result := getItemIcon("armor")
		if result != "🛡️" {
			t.Errorf("expected 🛡️, got %s", result)
		}
	})

	tests.Run("potion returns potion emoji", func(t *testing.T) {
		result := getItemIcon("potion")
		if result != "🧪" {
			t.Errorf("expected 🧪, got %s", result)
		}
	})

	tests.Run("food returns food emoji", func(t *testing.T) {
		result := getItemIcon("food")
		if result != "🍖" {
			t.Errorf("expected 🍖, got %s", result)
		}
	})

	tests.Run("scroll returns scroll emoji", func(t *testing.T) {
		result := getItemIcon("scroll")
		if result != "📜" {
			t.Errorf("expected 📜, got %s", result)
		}
	})

	tests.Run("key returns key emoji", func(t *testing.T) {
		result := getItemIcon("key")
		if result != "🔑" {
			t.Errorf("expected 🔑, got %s", result)
		}
	})

	tests.Run("treasure returns gem emoji", func(t *testing.T) {
		result := getItemIcon("treasure")
		if result != "💎" {
			t.Errorf("expected 💎, got %s", result)
		}
	})

	tests.Run("quest returns clipboard emoji", func(t *testing.T) {
		result := getItemIcon("quest")
		if result != "📋" {
			t.Errorf("expected 📋, got %s", result)
		}
	})

	tests.Run("unknown type returns box emoji", func(t *testing.T) {
		result := getItemIcon("unknown")
		if result != "📦" {
			t.Errorf("expected 📦, got %s", result)
		}
	})
}

func TestGetItemRarityColor(t *testing.T) {
	t.Parallel()

	t.Run("rare returns blue", func(t *testing.T) {
		result := getItemRarityColor("rare")
		if result != "51" {
			t.Errorf("expected 51, got %s", result)
		}
	})

	tests.Run("epic returns magenta", func(t *testing.T) {
		result := getItemRarityColor("epic")
		if result != "201" {
			t.Errorf("expected 201, got %s", result)
		}
	})

	tests.Run("legendary returns gold", func(t *testing.T) {
		result := getItemRarityColor("legendary")
		if result != "220" {
			t.Errorf("expected 220, got %s", result)
		}
	})

	tests.Run("common returns white", func(t *testing.T) {
		result := getItemRarityColor("common")
		if result != "white" {
			t.Errorf("expected white, got %s", result)
		}
	})

	tests.Run("empty returns white", func(t *testing.T) {
		result := getItemRarityColor("")
		if result != "white" {
			t.Errorf("expected white, got %s", result)
		}
	})
}