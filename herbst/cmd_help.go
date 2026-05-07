package main

import "fmt"

// handleHelpCommand displays the help text
func (m *model) handleHelpCommand(_ *model, args []string) {
	m.AppendMessage(`Commands:
  n/north, s/south, e/east, w/west - Move (cardinal)
  ne/northeast, se/southeast, sw/southwest, nw/northwest - Move (ordinal)
  u/up, d/down - Move (vertical)
  look/l [target] - Look around (or examine: look <target>, look at <target>)
  attack/a/kill/fight <target> - Attack a target
  loot [corpse] - Loot a corpse or all corpses in room
  ctrl+p - Scroll output up (older messages)
  ctrl+n - Scroll output down (newer messages)
  exits/x - Show exits
  peer <dir> - Peek at adjacent room
  take/get <item> - Pick up an item
  drop <item> - Drop an item
  inventory/inv/i - Show your inventory
  equip <item> - Equip an item from inventory
  equip talent <id> <slot> - Equip talent to slot 1-4
  equip potion <id> - Equip potion to R slot
  unequip <slot|item> - Unequip an item
  quests/q - Show your quest log
  whoami - Show your info
  profile/p - Edit character profile
  skills - Show your equipped combat skills
  skill slot <1-5> - Select a skill for a slot
  skill all - Show available classless skills
  skill swap <s1> <s2> - Swap skills between slots
  talents - Show your talents
  debug - Toggle debug mode
  clear/cls - Clear screen
  quit - Exit game`, "info")
}

func (m *model) handleWhoamiCommand(_ *model, args []string) {
	m.AppendMessage(fmt.Sprintf(
		"=== Character Status ===\nUser: %s (ID: %d)\nRoom: %s\n\n[Level %d - %d XP]\n%s",
		m.currentUserName, m.currentUserID, m.roomName,
		m.characterLevel, m.characterExperience,
		StatusBar(m.characterHP, m.characterMaxHP, m.characterStamina, m.characterMaxStamina, m.characterMana, m.characterMaxMana)), "info")
}

func (m *model) handleProfileCommand(_ *model, args []string) {
	m.screen = ScreenProfile
	m.menuItems = []string{"Edit Gender", "Edit Description", "Back to Game"}
	m.menuCursor = 0
	m.AppendMessage("", "")
}