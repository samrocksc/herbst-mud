package screens

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"herbst-mud/admin-tui/api"
	"herbst-mud/admin-tui/style"
)

// QuestsModel is the quests read-only management screen.
type QuestsModel struct {
	token    string
	quests   []api.Quest
	loading  bool
	errMsg   string
	selected int
	width    int
}

// NewQuestsScreen creates the quests screen.
func NewQuestsScreen(token string) tea.Model {
	return QuestsModel{token: token, loading: true, selected: 0}
}

func (m QuestsModel) Init() tea.Cmd {
	return func() tea.Msg {
		quests, err := api.ListQuests()
		if err != nil {
			return QuestsErrMsg{Err: err}
		}
		return QuestsMsg{Quests: quests}
	}
}

func (m QuestsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		return m, nil
	case QuestsMsg:
		m.quests = msg.Quests
		m.loading = false
		return m, nil
	case QuestsErrMsg:
		m.errMsg = fmt.Sprintf("Failed to load quests: %v", msg.Err)
		m.loading = false
		return m, nil
	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func (m QuestsModel) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		return m, func() tea.Msg { return NavigateMsg{Screen: 1} }
	case "r":
		m.loading = true
		return m, func() tea.Msg {
			quests, err := api.ListQuests()
			if err != nil {
				return QuestsErrMsg{Err: err}
			}
			return QuestsMsg{Quests: quests}
		}
	case "up", "k":
		if m.selected > 0 {
			m.selected--
		}
		return m, nil
	case "down", "j":
		if m.selected < len(m.quests)-1 {
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
	case "b":
		return m, func() tea.Msg { return NavigateMsg{Screen: 8} }
	case "w":
		return m, func() tea.Msg { return NavigateMsg{Screen: 9} }
	}
	return m, nil
}

func (m QuestsModel) View() string {
	if m.loading {
		return style.Info("Loading quests...")
	}

	lines := []string{
		style.StyleHeader.Render(fmt.Sprintf("Quests (%d)", len(m.quests))),
		style.RenderDivider(max(90, m.width-4)),
		fmt.Sprintf("%-6s %-22s %-12s %-18s %s",
			style.StyleTableHeader.Render("ID"),
			style.StyleTableHeader.Render("Name"),
			style.StyleTableHeader.Render("Type"),
			style.StyleTableHeader.Render("Giver"),
			style.StyleTableHeader.Render("Description"),
		),
	}

	for i, q := range m.quests {
		rowStyle := style.StyleTableRow
		if i == m.selected {
			rowStyle = lipgloss.Style{}.Foreground(style.ColorPrimary).Bold(true)
		}
		lines = append(lines, fmt.Sprintf("%-6s %-22s %-12s %-18s %s",
			rowStyle.Render(trunc(q.ID, 5)),
			rowStyle.Render(trunc(q.Name, 21)),
			rowStyle.Render(q.Type),
			rowStyle.Render(trunc(q.Giver, 17)),
			rowStyle.Render(trunc(q.Description, 40)),
		))
	}

	if len(m.quests) == 0 {
		lines = append(lines, style.StyleMuted.Render("  No quests found — quests are defined in content YAML"))
	}

	lines = append(lines, "")
	lines = append(lines, style.StyleMuted.Render("  Quests are read-only (defined in content YAML)   [R] refresh   [Esc] back"))

	if m.errMsg != "" {
		lines = append(lines, "", style.Error(m.errMsg))
	}

	return strings.Join(lines, "\n")
}

// ─── Messages ───────────────────────────────────────────────────────────────

type QuestsMsg struct{ Quests []api.Quest }
type QuestsErrMsg struct{ Err error }
