package screens

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"herbst-mud/admin-tui/api"
	"herbst-mud/admin-tui/style"
)

type NPCsModel struct {
	token       string
	npcs        []api.Character
	loading     bool
	errMsg      string
	selected    int
	mode        string
	confirmDel  bool
	width       int
}

func NewNPCsScreen(token string) tea.Model {
	return NPCsModel{token: token, loading: true, selected: 0, mode: "list"}
}

func (m NPCsModel) Init() tea.Cmd {
	return func() tea.Msg {
		chars, err := api.ListNPCs()
		if err != nil {
			return NPCsErrMsg{Err: err}
		}
		return NPCsMsg{Characters: chars}
	}
}

func (m NPCsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		return m, nil
	case NPCsMsg:
		m.npcs = msg.Characters
		m.loading = false
		return m, nil
	case NPCsErrMsg:
		m.errMsg = fmt.Sprintf("Failed to load NPCs: %v", msg.Err)
		m.loading = false
		return m, nil
	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func (m NPCsModel) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		if m.mode != "list" {
			m.mode = "list"
			return m, nil
		}
		return m, func() tea.Msg { return NavigateMsg{Screen: 1} }
	case "enter":
		m.mode = "edit"
		return m, nil
	case "d":
		if m.mode == "edit" {
			if m.confirmDel {
				go func() { api.DeleteNPC(m.npcs[m.selected].ID) }()
				m.npcs = filterNPCs(m.npcs, m.npcs[m.selected].ID)
				m.mode = "list"
				m.confirmDel = false
				return m, nil
			}
			m.confirmDel = true
			return m, nil
		}
	case "r":
		m.loading = true
		return m, func() tea.Msg {
			chars, err := api.ListNPCs()
			if err != nil {
				return NPCsErrMsg{Err: err}
			}
			return NPCsMsg{Characters: chars}
		}
	case "up", "k":
		if m.selected > 0 {
			m.selected--
		}
		return m, nil
	case "down", "j":
		if m.selected < len(m.npcs)-1 {
			m.selected++
		}
		return m, nil
	}
	if m.mode == "list" {
		switch msg.String() {
		case "u": return m, func() tea.Msg { return NavigateMsg{Screen: 2} }
		case "c": return m, func() tea.Msg { return NavigateMsg{Screen: 3} }
		case "r": return m, func() tea.Msg { return NavigateMsg{Screen: 4} }
		case "i": return m, func() tea.Msg { return NavigateMsg{Screen: 6} }
		case "q": return m, func() tea.Msg { return NavigateMsg{Screen: 7} }
		case "b": return m, func() tea.Msg { return NavigateMsg{Screen: 8} }
		case "w": return m, func() tea.Msg { return NavigateMsg{Screen: 9} }
		}
	}
	return m, nil
}

func filterNPCs(chars []api.Character, id int) []api.Character {
	result := make([]api.Character, 0, len(chars))
	for _, c := range chars {
		if c.ID != id {
			result = append(result, c)
		}
	}
	return result
}

func (m NPCsModel) View() string {
	if m.loading {
		return style.Info("Loading NPCs...")
	}
	if m.mode == "edit" && m.selected < len(m.npcs) {
		return m.viewDetail()
	}
	return m.viewList()
}

func (m NPCsModel) viewList() string {
	lines := []string{
		style.StyleHeader.Render(fmt.Sprintf("NPCs (%d)", len(m.npcs))),
		style.RenderDivider(max(90, m.width-4)),
		fmt.Sprintf("%-4s %-20s %-10s %-4s %-6s %-6s %-10s %-10s",
			style.StyleTableHeader.Render("ID"),
			style.StyleTableHeader.Render("Name"),
			style.StyleTableHeader.Render("Class"),
			style.StyleTableHeader.Render("Lvl"),
			style.StyleTableHeader.Render("HP"),
			style.StyleTableHeader.Render("Room"),
			style.StyleTableHeader.Render("Behavior"),
			style.StyleTableHeader.Render("Aggression"),
		),
	}
	for i, n := range m.npcs {
		rowStyle := style.StyleTableRow
		if i == m.selected {
			rowStyle = lipgloss.Style{}.Foreground(style.ColorPrimary).Bold(true)
		}
		lines = append(lines, fmt.Sprintf("%-4s %-20s %-10s %-4s %-6s %-6s %-10s %-10s",
			rowStyle.Render(fmt.Sprintf("%d", n.ID)),
			rowStyle.Render(trunc(n.Name, 19)),
			rowStyle.Render(n.Class),
			rowStyle.Render(fmt.Sprintf("%d", n.Level)),
			rowStyle.Render(fmt.Sprintf("%d/%d", n.HP, n.MaxHP)),
			rowStyle.Render(fmt.Sprintf("%d", n.RoomID)),
			rowStyle.Render(trunc(n.Behavior, 9)),
			rowStyle.Render(trunc(n.Aggression, 9)),
		))
	}
	if len(m.npcs) == 0 {
		lines = append(lines, style.StyleMuted.Render("  No NPCs found"))
	}
	lines = append(lines, "")
	lines = append(lines, style.StyleMuted.Render("  [Enter] view/edit   [D] delete   [R] refresh   [Esc] back"))
	if m.errMsg != "" {
		lines = append(lines, "", style.Error(m.errMsg))
	}
	return strings.Join(lines, "\n")
}

func (m NPCsModel) viewDetail() string {
	n := m.npcs[m.selected]
	lines := []string{
		style.StyleHeader.Render(fmt.Sprintf("NPC #%d — %s", n.ID, n.Name)),
		style.RenderDivider(max(90, m.width-4)),
		fmt.Sprintf("  %-12s %s", style.StyleLabel.Render("Name:"), style.StyleValue.Render(n.Name)),
		fmt.Sprintf("  %-12s %s", style.StyleLabel.Render("Class:"), style.StyleValue.Render(n.Class)),
		fmt.Sprintf("  %-12s %s", style.StyleLabel.Render("Level:"), style.StyleValue.Render(fmt.Sprintf("%d", n.Level))),
		fmt.Sprintf("  %-12s %s", style.StyleLabel.Render("HP:"), style.StyleValue.Render(fmt.Sprintf("%d/%d", n.HP, n.MaxHP))),
		fmt.Sprintf("  %-12s %s", style.StyleLabel.Render("Room:"), style.StyleValue.Render(fmt.Sprintf("%d", n.RoomID))),
		fmt.Sprintf("  %-12s %s", style.StyleLabel.Render("Behavior:"), style.StyleValue.Render(n.Behavior)),
		fmt.Sprintf("  %-12s %s", style.StyleLabel.Render("Aggression:"), style.StyleValue.Render(n.Aggression)),
		fmt.Sprintf("  %-12s %s", style.StyleLabel.Render("Description:"), style.StyleValue.Render(trunc(n.Description, 60))),
		"",
		style.StyleMuted.Render("  [D] delete   [Esc] back to list"),
	}
	if m.confirmDel {
		lines = append(lines, "", style.StyleDanger.Render("  Confirm DELETE? Press [D] again to confirm"))
	}
	return strings.Join(lines, "\n")
}

type NPCsMsg struct{ Characters []api.Character }
type NPCsErrMsg struct{ Err error }
