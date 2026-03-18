package combat

import (
	"fmt"
	"math/rand"
	"time"
)

// VictoryResult contains the results of winning a combat
type VictoryResult struct {
	Loot          []LootItem     `json:"loot"`
	Coins         int            `json:"coins"`
	XP            int            `json:"xp"`
	SkillUps      []SkillUp      `json:"skillUps"`
	WeaponDrops   []WeaponDrop   `json:"weaponDrops,omitempty"`
	CombatLog     []string       `json:"combatLog"`
	VictoryTime   time.Time      `json:"victoryTime"`
}

// LootItem represents an item dropped by enemies
type LootItem struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
	Rarity      string `json:"rarity"` // "common", "uncommon", "rare", "legendary"
}

// SkillUp represents skill experience gained
type SkillUp struct {
	SkillName string `json:"skillName"`
	Amount    int    `json:"amount"`
	NewLevel  int    `json:"newLevel"`
}

// DefeatResult contains the consequences of losing a combat
type DefeatResult struct {
	XPLost          int            `json:"xpLost"`
	ItemDropped     *LootItem      `json:"itemDropped,omitempty"`
	RespawnPoint    int            `json:"respawnPoint"` // Room ID
	CorpseLocation  int            `json:"corpseLocation"` // Room ID
	CorpseExpiry    time.Time      `json:"corpseExpiry"`
	Consequences    []string       `json:"consequences"`
	DefeatTime      time.Time      `json:"defeatTime"`
}

// HandleVictory processes the victory and calculates rewards
func (c *Combat) HandleVictory() *VictoryResult {
	result := &VictoryResult{
		VictoryTime: time.Now(),
	}

	// Get all defeated NPCs
	for _, p := range c.Participants {
		if p.IsNPC && !p.IsAlive {
			// Generate loot from this NPC
			loot := generateNPCLoot(p)
			result.Loot = append(result.Loot, loot...)

			// Add coins
			result.Coins += rand.Intn(50) + 10 // 10-60 coins per enemy

			// Add XP
			result.XP += calculateEnemyXP(p)

			// Check for weapon drop
			if p.HasGuaranteedDrop && p.DropWeapon != "" {
				drop := WeaponDrop{
					Name:             p.DropWeapon,
					WeaponType:       getWeaponType(p.DropWeapon),
					ClassRestriction: getClassRestriction(p.DropWeapon),
					MinDamage:        4,
					MaxDamage:        8,
					Guaranteed:       true,
				}
				result.WeaponDrops = append(result.WeaponDrops, drop)
			}
		}
	}

	// Calculate skill ups for skills used in combat
	// For now, we'll give skill XP based on actions taken
	result.SkillUps = calculateSkillUps(c)

	// Copy combat log
	for _, entry := range c.Log {
		result.CombatLog = append(result.CombatLog, entry.Message)
	}

	return result
}

// HandleDefeat processes the defeat and calculates consequences
func (c *Combat) HandleDefeat(playerStartingRoom int) *DefeatResult {
	result := &DefeatResult{
		DefeatTime:     time.Now(),
		RespawnPoint:   playerStartingRoom, // Default to starting area
		CorpseExpiry:   time.Now().Add(5 * time.Minute),
	}

	// Find the player who died
	var player *Participant
	for _, p := range c.Participants {
		if p.IsPlayer && !p.IsAlive {
			player = p
			break
		}
	}

	if player == nil {
		// No dead player found (shouldn't happen, but handle gracefully)
		return result
	}

	// Calculate XP loss (10% of current level progress)
	result.XPLost = 0 // This would be calculated from player's current XP
	result.Consequences = append(result.Consequences, 
		fmt.Sprintf("Lose %d XP (10%% of current level progress)", result.XPLost))

	// Store corpse location
	result.CorpseLocation = c.RoomID
	result.Consequences = append(result.Consequences, 
		"Your equipment has been left in the combat area")

	// Add respawn info
	result.Consequences = append(result.Consequences, 
		"Return within 5 minutes to reclaim your equipment!")

	return result
}

// generateNPCLoot generates loot items for a defeated NPC
func generateNPCLoot(npc *Participant) []LootItem {
	var items []LootItem

	// Base loot - enemy-specific items would come from a definition
	// For now, generate generic salvage items
	salvageItems := []LootItem{
		{ID: 1, Name: "Rusty Pipe", Description: "A bent metal pipe, slightly corroded", Quantity: 1, Rarity: "common"},
		{ID: 2, Name: "Scrap Metal", Description: "Twisted metal fragments", Quantity: rand.Intn(3) + 1, Rarity: "common"},
		{ID: 3, Name: "Ragged Cloth", Description: "Torn fabric scraps", Quantity: rand.Intn(2) + 1, Rarity: "common"},
	}

	// Random chance for better loot
	roll := rand.Float64()
	if roll < 0.3 { // 30% chance for uncommon
		items = append(items, LootItem{
			ID:          10,
			Name:        "Polished Gear",
			Description: "A well-maintained mechanical gear",
			Quantity:     1,
			Rarity:      "uncommon",
		})
	}

	if roll < 0.1 { // 10% chance for rare
		items = append(items, LootItem{
			ID:          20,
			Name:        "Crystal Shard",
			Description: "A glowing crystal fragment",
			Quantity:     1,
			Rarity:      "rare",
		})
	}

	// Add base salvage
	items = append(items, salvageItems[rand.Intn(len(salvageItems))])

	return items
}

// calculateEnemyXP calculates XP reward for defeating an enemy
func calculateEnemyXP(npc *Participant) int {
	// Base XP based on enemy level
	baseXP := npc.Level * 15

	// Bonus for tough enemies
	if npc.MaxHP > 50 {
		baseXP += 20
	}
	if npc.MaxHP > 100 {
		baseXP += 30
	}

	return baseXP
}

// calculateSkillUps calculates skill experience gains from combat
func calculateSkillUps(combat *Combat) []SkillUp {
	// For now, return basic skill ups based on combat actions
	// In a full implementation, this would track which skills were used
	var skillUps []SkillUp

	// Check if blades was used (warrior)
	for _, entry := range combat.Log {
		if entry.Type == "action" {
			// Check action type and award appropriate skill XP
			skillUps = append(skillUps, SkillUp{
				SkillName: "blades",
				Amount:    2,
				NewLevel:   0, // Would be current + amount
			})
			break
		}
	}

	// Also give brawling XP for basic attacks
	skillUps = append(skillUps, SkillUp{
		SkillName: "brawling",
		Amount:    1,
		NewLevel:  0,
	})

	return skillUps
}

// getWeaponType returns the weapon type for a weapon name
func getWeaponType(name string) string {
	weaponTypes := map[string]string{
		"Rusty Sword":      "sword",
		"Scrap Machete":    "sword",
		"Twisted Pipe":     "pipe",
		"Chef's Knife":     "dagger",
		"Broken Bottle":    "dagger",
	}
	if t, ok := weaponTypes[name]; ok {
		return t
	}
	return "weapon"
}

// getClassRestriction returns class restrictions for a weapon
func getClassRestriction(name string) string {
	classRestrictions := map[string]string{
		"Rusty Sword":      "warrior",
		"Scrap Machete":    "warrior",
		"Twisted Pipe":     "chef",
		"Chef's Knife":     "chef",
	}
	if c, ok := classRestrictions[name]; ok {
		return c
	}
	return "" // No restriction
}