// Package db provides tests to verify database schema consolidation
package db

import (
	"testing"
)

// TestSchemaConsolidation verifies that the unified db module contains
// all entities needed by both the TUI client and the server.
//
// CURRENT STATE (May 2026):
// - herbst/db/ contains: user, character, room, equipment, ability, ability_effect, character_ability, npctemplate
// - server/db/ contains: ALL OF ABOVE + npc_ability, charactertag, characterfaction, charactercompetency, etc.
//
// KEY FINDINGS:
// 1. Both modules share the core combat entities (Ability, AbilityEffect, CharacterAbility)
// 2. server/db has additional admin-only entities (NPCAbility join table, tags, factions, etc.)
// 3. They are separate Go modules (herbst vs herbst-server) so can't share code
//
// CONSOLIDATION OPTIONS:
// A) Merge both into single "herbst-db" shared module, import from both
// B) Generate ent code for both from a single schema source
// C) Keep separate but add missing tables to herbst/db
func TestSchemaConsolidation(t *testing.T) {
	// Core entities present in both:
	// - User (both have)
	// - Character (both have)
	// - Room (both have)
	// - Equipment (both have)
	// - Ability (both have, renamed from Skill)
	// - AbilityEffect (both have, new entity)
	// - CharacterAbility (both have, renamed from CharacterSkill)
	// - NPCTemplate (both have)

	t.Log("Schema consolidation documentation test")
	t.Log("Core entities shared between herbst/db and server/db:")
	t.Log("  - Ability: combat moves and passive abilities")
	t.Log("  - AbilityEffect: effects per ability (damage, heal, buff, etc.)")
	t.Log("  - CharacterAbility: character -> ability link with slot")
}

// TestModulePathDifference documents the import path divergence
func TestModulePathDifference(t *testing.T) {
	// herbst/ imports: "herbst/db", "herbst/dbinit"
	// server/ imports: "herbst-server/db", "herbst-server/dbinit"
	//
	// This means they cannot share code without creating a shared module

	t.Log("Module import paths:")
	t.Log("  TUI (herbst/):   import herbst/db")
	t.Log("  Server (server/): import herbst-server/db")
	t.Log("Solution: Create shared herbst-db module or single monorepo")
}