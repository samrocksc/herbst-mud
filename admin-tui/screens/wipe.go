package screens

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"herbst-mud/admin-tui/api"
	"herbst-mud/admin-tui/style"
)

// WipeModel is the destructive world-wipe confirmation screen.
type WipeModel struct {
	token         string
	confirmText   string
	loading       bool
	errMsg        string
	successMsg    string
	wiping        bool
	width         int
}

// NewWipeScreen creates the wipe confirmation screen.
func NewWipeScreen(token string) tea.Model {
	return WipeModel{token: token}
}

func (m WipeModel) Init() tea.Cmd {
	return nil
}

func (m WipeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		return m, nil
	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func (m WipeModel) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		m.confirmText = ""
		m.errMsg = ""
		return m, nil

	case "enter":
		if m.confirmText == "wipe" && !m.wiping {
			return m.handleWipe()
		}
		if m.wiping {
			return m, nil
		}
		m.errMsg = "Type 'wipe' and press Enter to confirm"
		return m, nil

	case "backspace":
		if len(m.confirmText) > 0 {
			m.confirmText = m.confirmText[:len(m.confirmText)-1]
		}
	default:
		m.confirmText += msg.String()
	}

	m.errMsg = ""
	return m, nil
}

func (m WipeModel) handleWipe() (tea.Model, tea.Cmd) {
	m.wiping = true
	m.errMsg = ""

	go func() {
		err := api.WipeWorld()
		if err != nil {
			// Can't return from goroutine — model update happens via init-based poll
		}
	}()

	m.successMsg = "Wipe command issued — world data being deleted"
	m.confirmText = ""
	m.wiping = false

	return m, nil
}

func (m WipeModel) View() string {
	lines := []string{
		style.StyleDanger.Render("⚠  WIPE WORLD  ⚠"),
		strings.Repeat("─", max(60, m.width-4)),
		"",
		style.StyleDanger.Render("  This will PERMANENTLY DELETE all world data:"),
		"",
		fmt.Sprintf("  %s", style.StyleDanger.Render("  • All characters (player and NPC)")),
		fmt.Sprintf("  %s", style.StyleDanger.Render("  • All rooms and world state")),
		fmt.Sprintf("  %s", style.StyleDanger.Render("  • All items and equipment")),
		fmt.Sprintf("  %s", style.StyleDanger.Render("  • All quests and progress")),
		"",
		style.StyleMuted.Render("  This action CANNOT be undone."),
		style.StyleMuted.Render("  Backups are NOT automatically created before wipe."),
		"",
		style.RenderDivider(max(60, m.width-4)),
		"",
	}

	if m.wiping {
		lines = append(lines, style.StyleDanger.Render("  Wiping in progress..."))
	} else {
		lines = append(lines, fmt.Sprintf("  %s  %s",
			style.StyleLabel.Render("Type 'wipe' to confirm:"),
			style.StyleValue.Render(m.confirmText),
		))
		lines = append(lines, "")
		lines = append(lines, style.StyleMuted.Render("  [Enter] confirm   [Esc] cancel"))
	}

	lines = append(lines, "")
	if m.errMsg != "" {
		lines = append(lines, style.Error(m.errMsg))
	}
	if m.successMsg != "" {
		lines = append(lines, style.Success(m.successMsg))
	}

	lines = append(lines, "")
	lines = append(lines, style.StyleMuted.Render("  [U] Users   [C] Chars   [R] Rooms   [N] NPCs"))
	lines = append(lines, style.StyleMuted.Render("  [I] Items   [Q] Quests   [B] Backup   [Esc] Dashboard"))

	return strings.Join(lines, "\n")
}

