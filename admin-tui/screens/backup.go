package screens

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"herbst-mud/admin-tui/api"
	"herbst-mud/admin-tui/style"
)

// BackupModel is the backup management screen.
type BackupModel struct {
	token       string
	backups     []api.BackupManifest
	loading     bool
	errMsg      string
	successMsg  string
	selected    int
	confirmWipe bool
	width       int
}

// NewBackupScreen creates the backup screen.
func NewBackupScreen(token string) tea.Model {
	return BackupModel{token: token, loading: true, selected: 0}
}

func (m BackupModel) Init() tea.Cmd {
	return func() tea.Msg {
		backups, err := api.ListBackups()
		if err != nil {
			return BackupErrMsg{Err: err}
		}
		return BackupMsg{Backups: backups}
	}
}

func (m BackupModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		return m, nil
	case BackupMsg:
		m.backups = msg.Backups
		m.loading = false
		return m, nil
	case BackupErrMsg:
		m.errMsg = fmt.Sprintf("Failed to load backups: %v", msg.Err)
		m.loading = false
		return m, nil
	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func (m BackupModel) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		return m, func() tea.Msg { return NavigateMsg{Screen: 1} }
	case "b":
		m.successMsg = ""
		m.loading = true
		return m, func() tea.Msg {
			err := api.TriggerBackup()
			if err != nil {
				return BackupErrMsg{Err: fmt.Errorf("backup failed: %v", err)}
			}
			backups, err := api.ListBackups()
			if err != nil {
				return BackupErrMsg{Err: err}
			}
			return BackupMsg{Backups: backups}
		}
	case "r":
		m.loading = true
		return m, func() tea.Msg {
			backups, err := api.ListBackups()
			if err != nil {
				return BackupErrMsg{Err: err}
			}
			return BackupMsg{Backups: backups}
		}
	case "up", "k":
		if m.selected > 0 {
			m.selected--
		}
		return m, nil
	case "down", "j":
		if m.selected < len(m.backups)-1 {
			m.selected++
		}
		return m, nil
	}
	switch msg.String() {
	case "u":
		return m, func() tea.Msg { return NavigateMsg{Screen: 2} }
	case "c":
		return m, func() tea.Msg { return NavigateMsg{Screen: 3} }
	case "r":
		return m, func() tea.Msg { return NavigateMsg{Screen: 4} }
	case "n":
		return m, func() tea.Msg { return NavigateMsg{Screen: 5} }
	case "i":
		return m, func() tea.Msg { return NavigateMsg{Screen: 6} }
	case "q":
		return m, func() tea.Msg { return NavigateMsg{Screen: 7} }
	case "w":
		return m, func() tea.Msg { return NavigateMsg{Screen: 9} }
	}
	return m, nil
}

func (m BackupModel) View() string {
	if m.loading {
		return style.Info("Loading backups...")
	}

	lines := []string{
		style.StyleHeader.Render(fmt.Sprintf("Backups (%d)", len(m.backups))),
		style.RenderDivider(max(80, m.width-4)),
		fmt.Sprintf("%-6s %-36s %-20s %-8s",
			style.StyleTableHeader.Render("ID"),
			style.StyleTableHeader.Render("Filename"),
			style.StyleTableHeader.Render("Created"),
			style.StyleTableHeader.Render("Size"),
		),
	}

	for i, b := range m.backups {
		rowStyle := style.StyleTableRow
		if i == m.selected {
			rowStyle = lipgloss.Style{}.Foreground(style.ColorPrimary).Bold(true)
		}
		lines = append(lines, fmt.Sprintf("%-6s %-36s %-20s %-8s",
			rowStyle.Render(trunc(b.ID, 5)),
			rowStyle.Render(trunc(b.Filename, 35)),
			rowStyle.Render(trunc(b.CreatedAt, 19)),
			rowStyle.Render(humanSize(b.Size)),
		))
	}

	if len(m.backups) == 0 {
		lines = append(lines, style.StyleMuted.Render("  No backups found"))
	}

	lines = append(lines, "")
	lines = append(lines, style.StyleMuted.Render("  [B] create backup   [R] refresh list   [Esc] back"))

	if m.errMsg != "" {
		lines = append(lines, "", style.Error(m.errMsg))
	}
	if m.successMsg != "" {
		lines = append(lines, "", style.Success(m.successMsg))
	}

	return strings.Join(lines, "\n")
}

// ─── Messages ───────────────────────────────────────────────────────────────

type BackupMsg struct{ Backups []api.BackupManifest }
type BackupErrMsg struct{ Err error }

