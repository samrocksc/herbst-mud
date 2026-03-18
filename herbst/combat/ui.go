package combat

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// UI styles for combat display
var (
	// Colors
	red       = lipgloss.Color("196")
	green     = lipgloss.Color("46")
	yellow    = lipgloss.Color("226")
	blue      = lipgloss.Color("75")
	purple    = lipgloss.Color("141")
	white     = lipgloss.Color("15")
	gray      = lipgloss.Color("8")
	orange    = lipgloss.Color("208")
	cyan      = lipgloss.Color("51")
	darkRed   = lipgloss.Color("88")
	darkBlue  = lipgloss.Color("18")
	
	// Box styles
	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(purple)
	
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(green).
			Background(darkBlue).
			Padding(0, 1)
	
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(blue).
			Padding(0, 1)
	
	// HP bar styles
	hpBarStyle = lipgloss.NewStyle().Foreground(red)
	manaBarStyle = lipgloss.NewStyle().Foreground(blue)
	staminaBarStyle = lipgloss.NewStyle().Foreground(yellow)
	
	// Combat log styles
	damageStyle = lipgloss.NewStyle().Foreground(red)
	healStyle = lipgloss.NewStyle().Foreground(green)
	infoStyle = lipgloss.NewStyle().Foreground(cyan)
	systemStyle = lipgloss.NewStyle().Foreground(yellow)
	
	// Action bar styles
	actionReadyStyle = lipgloss.NewStyle().
				Foreground(white).
				Background(darkBlue).
				Padding(0, 2).
				Bold(true)
	
	actionCooldownStyle = lipgloss.NewStyle().
				Foreground(gray).
				Background(darkBlue).
				Padding(0, 2)
	
	// Tick timer styles
	tickActiveStyle = lipgloss.NewStyle().
				Foreground(green).
				Bold(true)
	
	tickWarningStyle = lipgloss.NewStyle().
				Foreground(yellow).
				Bold(true)
	
	tickCriticalStyle = lipgloss.NewStyle().
				Foreground(red).
				Bold(true)
)

// CombatUI renders the combat screen
type CombatUI struct {
	width  int
	height int
}

// NewCombatUI creates a new combat UI renderer
func NewCombatUI(width, height int) *CombatUI {
	return &CombatUI{
		width:  width,
		height: height,
	}
}

// Render generates the full combat screen
func (ui *CombatUI) Render(combat *Combat, stateMachine *CombatStateMachine, inputManager *InputManager, playerID int) string {
	var b strings.Builder
	
	// Title bar
	b.WriteString(ui.renderTitleBar(combat))
	b.WriteString("\n")
	
	// Enemy HP bars
	b.WriteString(ui.renderEnemyHP(combat))
	b.WriteString("\n")
	
	// Turn order
	b.WriteString(ui.renderTurnOrder(combat))
	b.WriteString("\n")
	
	// Combat log
	b.WriteString(ui.renderCombatLog(combat, 5))
	b.WriteString("\n")
	
	// Player status (HP/Mana/Stamina bars)
	b.WriteString(ui.renderPlayerStatus(combat, playerID))
	b.WriteString("\n")
	
	// Tick counter
	b.WriteString(ui.renderTickCounter(combat, stateMachine, inputManager, playerID))
	b.WriteString("\n")
	
	// Status effects
	b.WriteString(ui.renderStatusEffects(combat, playerID))
	b.WriteString("\n")
	
	// Action bar (talents 1-4)
	b.WriteString(ui.renderActionBar(combat, playerID, inputManager))
	b.WriteString("\n")
	
	// Input prompt
	b.WriteString(ui.renderInputPrompt(combat, stateMachine, playerID))
	
	return b.String()
}

// renderTitleBar renders the combat title
func (ui *CombatUI) renderTitleBar(combat *Combat) string {
	stateStr := string(combat.State)
	if combat.State == StateActive {
		roundNum := 1
		if len(combat.Participants) > 0 {
			roundNum = combat.TickNumber/len(combat.Participants) + 1
		}
		stateStr = fmt.Sprintf("Round %d", roundNum)
	}
	
	title := titleStyle.Render(fmt.Sprintf("⚔️  COMBAT - %s  ⚔️", stateStr))
	return titleStyle.Render(title)
}

// renderEnemyHP renders HP bars for all enemies
func (ui *CombatUI) renderEnemyHP(combat *Combat) string {
	var b strings.Builder
	
	b.WriteString(headerStyle.Render("Enemies:"))
	b.WriteString("\n")
	
	enemies := combat.GetAliveByTeam(1) // Team 1 = enemies
	
	if len(enemies) == 0 {
		b.WriteString(infoStyle.Render("  No enemies remaining"))
		return b.String()
	}
	
	for _, enemy := range enemies {
		hpBar := ui.renderHPBar(enemy.HP, enemy.MaxHP, 20)
		status := ""
		if !enemy.IsAlive {
			status = " [DEAD]"
		}
		
		line := fmt.Sprintf("  %s %s %d/%d%s",
			enemy.Name,
			hpBar,
			enemy.HP,
			enemy.MaxHP,
			status,
		)
		b.WriteString(line)
		b.WriteString("\n")
	}
	
	return b.String()
}

// renderHPBar creates a text-based HP progress bar
func (ui *CombatUI) renderHPBar(current, max, width int) string {
	if max <= 0 {
		max = 1
	}
	if current < 0 {
		current = 0
	}
	if current > max {
		current = max
	}
	
	percent := float64(current) / float64(max)
	filled := int(percent * float64(width))
	empty := width - filled
	
	var color lipgloss.Color
	if percent > 0.6 {
		color = green
	} else if percent > 0.3 {
		color = yellow
	} else {
		color = red
	}
	
	filledStr := lipgloss.NewStyle().Foreground(color).Render(strings.Repeat("█", filled))
	emptyStr := lipgloss.NewStyle().Foreground(gray).Render(strings.Repeat("░", empty))
	
	return fmt.Sprintf("[%s%s]", filledStr, emptyStr)
}

// renderTurnOrder renders the turn order display
func (ui *CombatUI) renderTurnOrder(combat *Combat) string {
	var b strings.Builder
	
	b.WriteString(headerStyle.Render("Turn Order:"))
	b.WriteString("\n")
	
	order := combat.GetTurnOrder()
	if len(order) == 0 {
		return b.String()
	}
	
	for i, p := range order {
		indicator := "  "
		if i == 0 {
			indicator = "▶ "
		}
		
		teamLabel := "Ally"
		color := green
		if p.Team == 1 {
			teamLabel = "Enemy"
			color = red
		}
		
		status := ""
		if !p.IsAlive {
			status = " ✗"
		}
		
		nameStyle := lipgloss.NewStyle().Foreground(color)
		line := fmt.Sprintf("%s%s%s (%d) - Init: %d%s",
			indicator,
			nameStyle.Render(p.Name),
			teamLabel,
			p.TurnPosition,
			p.Initiative,
			status,
		)
		b.WriteString(line)
		b.WriteString("\n")
	}
	
	return b.String()
}

// renderCombatLog renders recent combat log entries
func (ui *CombatUI) renderCombatLog(combat *Combat, maxEntries int) string {
	var b strings.Builder
	
	b.WriteString(headerStyle.Render("Combat Log:"))
	b.WriteString("\n")
	
	// Get last N entries
	start := len(combat.Log) - maxEntries
	if start < 0 {
		start = 0
	}
	
	entries := combat.Log[start:]
	
	for _, entry := range entries {
		var line string
		
		switch entry.Type {
		case "damage":
			line = damageStyle.Render(fmt.Sprintf("  💢 %s", entry.Message))
		case "heal":
			line = healStyle.Render(fmt.Sprintf("  💚 %s", entry.Message))
		case "system":
			line = systemStyle.Render(fmt.Sprintf("  ⚙️ %s", entry.Message))
		case "info":
			line = infoStyle.Render(fmt.Sprintf("  ℹ️ %s", entry.Message))
		default:
			line = fmt.Sprintf("  %s", entry.Message)
		}
		
		b.WriteString(line)
		b.WriteString("\n")
	}
	
	return b.String()
}

// renderPlayerStatus renders the player's HP/Mana/Stamina bars
func (ui *CombatUI) renderPlayerStatus(combat *Combat, playerID int) string {
	player := combat.GetParticipantByID(playerID)
	if player == nil {
		return ""
	}
	
	var b strings.Builder
	
	b.WriteString(headerStyle.Render(fmt.Sprintf("%s's Status:", player.Name)))
	b.WriteString("\n")
	
	// HP Bar
	hpBar := ui.renderHPBar(player.HP, player.MaxHP, 15)
	b.WriteString(fmt.Sprintf("  ❤️ HP:    %s %d/%d\n", hpBar, player.HP, player.MaxHP))
	
	// Mana Bar
	manaPercent := float64(player.Mana) / float64(player.MaxMana)
	manaFilled := int(manaPercent * 15)
	manaBar := lipgloss.NewStyle().Foreground(blue).Render(strings.Repeat("✨", manaFilled)) +
		lipgloss.NewStyle().Foreground(gray).Render(strings.Repeat("○", 15-manaFilled))
	b.WriteString(fmt.Sprintf("  ✨ Mana:  [%s] %d/%d\n", manaBar, player.Mana, player.MaxMana))
	
	// Stamina Bar
	staminaPercent := float64(player.Stamina) / float64(player.MaxStamina)
	staminaFilled := int(staminaPercent * 15)
	staminaBar := lipgloss.NewStyle().Foreground(yellow).Render(strings.Repeat("💪", staminaFilled)) +
		lipgloss.NewStyle().Foreground(gray).Render(strings.Repeat("○", 15-staminaFilled))
	b.WriteString(fmt.Sprintf("  💪 STA:   [%s] %d/%d\n", staminaBar, player.Stamina, player.MaxStamina))
	
	return b.String()
}

// renderTickCounter renders the tick countdown timer
func (ui *CombatUI) renderTickCounter(combat *Combat, stateMachine *CombatStateMachine, inputManager *InputManager, playerID int) string {
	remaining := inputManager.GetTimeRemaining(playerID)
	
	var style lipgloss.Style
	remainingFloat := 0.0
	fmt.Sscanf(remaining, "%f", &remainingFloat)
	
	if remainingFloat > 1.0 {
		style = tickActiveStyle
	} else if remainingFloat > 0.5 {
		style = tickWarningStyle
	} else {
		style = tickCriticalStyle
	}
	
	tickStr := fmt.Sprintf("⏱️ Tick %d | Time: %ss", combat.TickNumber, remaining)
	
	// Show current actor
	currentActor := stateMachine.GetCurrentActor()
	if currentActor != nil {
		actorStyle := lipgloss.NewStyle().Foreground(cyan).Bold(true)
		tickStr += fmt.Sprintf(" | %s's turn", actorStyle.Render(currentActor.Name))
	}
	
	return style.Render(tickStr)
}

// renderStatusEffects renders active status effects for a participant
func (ui *CombatUI) renderStatusEffects(combat *Combat, playerID int) string {
	var b strings.Builder

	// Get active effects from the registry
	effects := combat.Effects.GetEffectsForParticipant(playerID)
	if len(effects) == 0 {
		return ""
	}

	b.WriteString(headerStyle.Render("Status Effects:"))
	b.WriteString("\n")

	for _, effect := range effects {
		var icon string
		switch effect.Type {
		case "buff":
			icon = "💚"
		case "debuff":
			icon = "💔"
		case "channel":
			icon = "🔄"
		case "charge":
			icon = "⏳"
		default:
			icon = "⚪"
		}

		remaining := ""
		if effect.TicksRemaining > 0 {
			remaining = fmt.Sprintf(" (%d ticks)", effect.TicksRemaining)
		}

		line := fmt.Sprintf("  %s %s%s\n", icon, effect.Name, remaining)
		b.WriteString(line)
	}

	return b.String()
}

// renderActionBar renders the action bar with talent slots
func (ui *CombatUI) renderActionBar(combat *Combat, playerID int, inputManager *InputManager) string {
	var b strings.Builder
	
	b.WriteString(headerStyle.Render("Actions:"))
	b.WriteString("\n")
	
	// Get player's talent bindings
	bindings := inputManager.GetTalentBindings(playerID)
	
	// Default to basic actions if no bindings
	if len(bindings) == 0 {
		bindings = map[int]string{1: "attack", 2: "defend", 3: "item", 4: "wait"}
	}
	
	// Render action slots
	slots := []int{1, 2, 3, 4}
	
	for _, slot := range slots {
		actionID := bindings[slot]
		action, found := GetActionDefinition(actionID)
		if !found {
			action = BasicActions["attack"] // Fallback
		}
		
		// Check if awaiting input
		awaiting := inputManager.IsAwaitingInput(playerID)
		
		var slotStr string
		if awaiting {
			slotStr = actionReadyStyle.Render(fmt.Sprintf("[%d] %s", slot, action.Name))
		} else {
			slotStr = actionCooldownStyle.Render(fmt.Sprintf("[%d] %s", slot, action.Name))
		}
		
		b.WriteString(fmt.Sprintf("  %s", slotStr))
		
		// Add tick cost info
		if action.TickCost > 0 {
			b.WriteString(fmt.Sprintf(" (%d tick)", action.TickCost))
		}
		
		// Add resource costs
		if action.ManaCost > 0 {
			b.WriteString(fmt.Sprintf(" -%d mana", action.ManaCost))
		}
		if action.StaminaCost > 0 {
			b.WriteString(fmt.Sprintf(" -%d sta", action.StaminaCost))
		}
		
		// Add type indicator
		switch action.Type {
		case ActionCharge:
			b.WriteString(" ⏳")
		case ActionChannel:
			b.WriteString(" 🔄")
		}
		
		b.WriteString("\n")
	}
	
	return b.String()
}

// renderInputPrompt renders the input prompt for the current actor
func (ui *CombatUI) renderInputPrompt(combat *Combat, stateMachine *CombatStateMachine, playerID int) string {
	player := combat.GetParticipantByID(playerID)
	if player == nil {
		return ""
	}
	
	// Check if it's the player's turn
	currentActor := stateMachine.GetCurrentActor()
	if currentActor == nil || currentActor.ID != playerID {
		return systemStyle.Render("Waiting for enemy action...")
	}
	
	// Check if awaiting input
	if !stateMachine.IsAwaitingInput(playerID) {
		return infoStyle.Render("Processing...")
	}
	
	remaining := stateMachine.GetTickCountdown()
	
	prompt := fmt.Sprintf(">>> Your turn! Select action [1-4] (%.1fs remaining)", remaining)
	return tickActiveStyle.Render(prompt)
}

// RenderCompactHPBar renders a compact inline HP bar
func RenderCompactHPBar(current, max int) string {
	if max <= 0 {
		max = 1
	}
	percent := float64(current) / float64(max)
	
	var color lipgloss.Color
	if percent > 0.6 {
		color = green
	} else if percent > 0.3 {
		color = yellow
	} else {
		color = red
	}
	
	blocks := int(percent * 5)
	if blocks > 5 {
		blocks = 5
	}
	if blocks < 0 {
		blocks = 0
	}
	
	bar := lipgloss.NewStyle().Foreground(color).Render(strings.Repeat("█", blocks)) +
		lipgloss.NewStyle().Foreground(gray).Render(strings.Repeat("░", 5-blocks))
	
	return fmt.Sprintf("%s %d/%d", bar, current, max)
}

// FormatTimeRemaining formats time remaining for display
func FormatTimeRemaining(d time.Duration) string {
	seconds := d.Seconds()
	if seconds < 0 {
		return "0.0"
	}
	return fmt.Sprintf("%.1f", seconds)
}

// VictoryScreenStyle for victory display
var (
	victoryTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(yellow).
				Background(darkRed).
				Padding(1, 4).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(yellow)

	victoryBoxStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(green).
				Padding(1, 2)

	lootHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(yellow).
				MarginTop(1)

	lootItemStyle = lipgloss.NewStyle().
				Foreground(white)

	xpStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(green)

	skillUpStyle = lipgloss.NewStyle().
			Foreground(cyan)

	defeatTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(white).
				Background(darkRed).
				Padding(1, 4).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(red)

	defeatBoxStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(red).
				Padding(1, 2)

	consequenceStyle = lipgloss.NewStyle().
				Foreground(yellow)

	respawnStyle = lipgloss.NewStyle().
			Foreground(cyan)
)

// RenderVictoryScreen renders the victory screen UI
func (ui *CombatUI) RenderVictoryScreen(result *VictoryResult, width int) string {
	var b strings.Builder

	// Title
	title := victoryTitleStyle.Render("⚔️  VICTORY!  ⚔️")
	b.WriteString(title)
	b.WriteString("\n\n")

	// Victory message
	b.WriteString(victoryBoxStyle.Render("All enemies have been defeated!"))
	b.WriteString("\n\n")

	// Loot section
	if len(result.Loot) > 0 {
		b.WriteString(lootHeaderStyle.Render("LOOT:"))
		b.WriteString("\n")
		for _, item := range result.Loot {
			rarityColor := white
			switch item.Rarity {
			case "uncommon":
				rarityColor = green
			case "rare":
				rarityColor = purple
			case "legendary":
				rarityColor = yellow
			}

			itemStyle := lipgloss.NewStyle().Foreground(rarityColor)
			b.WriteString(fmt.Sprintf("  • %s", itemStyle.Render(item.Name)))
			if item.Quantity > 1 {
				b.WriteString(fmt.Sprintf(" x%d", item.Quantity))
			}
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	// Coins
	if result.Coins > 0 {
		b.WriteString(xpStyle.Render(fmt.Sprintf("  💰 %d coins", result.Coins)))
		b.WriteString("\n")
	}

	// XP
	if result.XP > 0 {
		b.WriteString(xpStyle.Render(fmt.Sprintf("  ⭐ %d XP", result.XP)))
		b.WriteString("\n\n")
	}

	// Weapon drops
	if len(result.WeaponDrops) > 0 {
		b.WriteString(lootHeaderStyle.Render("WEAPONS:"))
		b.WriteString("\n")
		for _, drop := range result.WeaponDrops {
			classText := ""
			if drop.ClassRestriction != "" {
				classText = fmt.Sprintf(" (%s only)", drop.ClassRestriction)
			}
			b.WriteString(fmt.Sprintf("  ⚔️  %s%s\n", drop.Name, classText))
			b.WriteString(fmt.Sprintf("      Damage: %d-%d\n", drop.MinDamage, drop.MaxDamage))
		}
		b.WriteString("\n")
	}

	// Skill ups
	if len(result.SkillUps) > 0 {
		b.WriteString(skillUpStyle.Render("SKILL UP:"))
		b.WriteString("\n")
		for _, su := range result.SkillUps {
			b.WriteString(fmt.Sprintf("  • %s +%d\n", su.SkillName, su.Amount))
		}
	}

	return b.String()
}

// RenderDefeatScreen renders the defeat screen UI
func (ui *CombatUI) RenderDefeatScreen(result *DefeatResult, enemyName string, width int) string {
	var b strings.Builder

	// Title
	title := defeatTitleStyle.Render("💀  DEFEAT  💀")
	b.WriteString(title)
	b.WriteString("\n\n")

	// Defeat message
	defeatMsg := fmt.Sprintf("You have fallen to %s...", enemyName)
	b.WriteString(defeatBoxStyle.Render(defeatMsg))
	b.WriteString("\n\n")

	// Consequences
	b.WriteString(consequenceStyle.Render("CONSEQUENCES:"))
	b.WriteString("\n")
	for _, c := range result.Consequences {
		b.WriteString(fmt.Sprintf("  • %s\n", c))
	}
	b.WriteString("\n")

	// Respawn info
	b.WriteString(respawnStyle.Render("You will respawn at the last save point."))
	b.WriteString("\n")

	// Corpse timer
	if !result.CorpseExpiry.IsZero() {
		timeRemaining := time.Until(result.CorpseExpiry)
		if timeRemaining > 0 {
			b.WriteString(respawnStyle.Render(fmt.Sprintf("Equipment reclaimable for: %.0f seconds", timeRemaining.Seconds())))
			b.WriteString("\n")
		}
	}

	return b.String()
}

// RenderCombatEnd renders either victory or defeat screen based on result
func RenderCombatEnd(combat *Combat, playerID, respawnRoom int, width int) string {
	ui := NewCombatUI(width, 20)

	// Check if player won or lost
	if combat.AllEnemiesDefeated() {
		result := combat.HandleVictory()
		return ui.RenderVictoryScreen(result, width)
	}

	if combat.AllPlayersDefeated() {
		// Find the enemy that delivered the killing blow
		enemyName := "an enemy"
		for _, p := range combat.Participants {
			if p.IsNPC && p.IsAlive && p.Team == 1 {
				enemyName = p.Name
				break
			}
		}

		result := combat.HandleDefeat(respawnRoom)
		return ui.RenderDefeatScreen(result, enemyName, width)
	}

	// Combat still ongoing
	return ""
}

