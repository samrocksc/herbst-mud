package screens

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"herbst-mud/admin-tui/api"
	"herbst-mud/admin-tui/style"
)

type RoomsModel struct {
	token       string
	rooms       []api.Room
	loading     bool
	errMsg      string
	selected    int
	mode        string
	confirmDel  bool
	width       int
}

func NewRoomsScreen(token string) tea.Model {
	return RoomsModel{token: token, loading: true, selected: 0, mode: "list"}
}

func (m RoomsModel) Init() tea.Cmd {
	return func() tea.Msg {
		rooms, err := api.ListRooms()
		if err != nil {
			return RoomsErrMsg{Err: err}
		}
		return RoomsMsg{Rooms: rooms}
	}
}

func (m RoomsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		return m, nil
	case RoomsMsg:
		m.rooms = msg.Rooms
		m.loading = false
		return m, nil
	case RoomsErrMsg:
		m.errMsg = fmt.Sprintf("Failed to load rooms: %v", msg.Err)
		m.loading = false
		return m, nil
	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func (m RoomsModel) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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
		m.mode = "detail"
		return m, nil
	case "c":
		if m.mode == "list" {
			// create not wired — show message
			return m, nil
		}
	case "d":
		if m.mode == "detail" {
			if m.confirmDel {
				go func() { api.DeleteRoom(m.rooms[m.selected].ID) }()
				m.rooms = filterRooms(m.rooms, m.rooms[m.selected].ID)
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
			rooms, err := api.ListRooms()
			if err != nil {
				return RoomsErrMsg{Err: err}
			}
			return RoomsMsg{Rooms: rooms}
		}
	case "up", "k":
		if m.selected > 0 {
			m.selected--
		}
		return m, nil
	case "down", "j":
		if m.selected < len(m.rooms)-1 {
			m.selected++
		}
		return m, nil
	}
	if m.mode == "list" {
		switch msg.String() {
		case "u": return m, func() tea.Msg { return NavigateMsg{Screen: 2} }
		case "c": return m, func() tea.Msg { return NavigateMsg{Screen: 3} }
		case "n": return m, func() tea.Msg { return NavigateMsg{Screen: 5} }
		case "i": return m, func() tea.Msg { return NavigateMsg{Screen: 6} }
		case "q": return m, func() tea.Msg { return NavigateMsg{Screen: 7} }
		case "b": return m, func() tea.Msg { return NavigateMsg{Screen: 8} }
		case "w": return m, func() tea.Msg { return NavigateMsg{Screen: 9} }
		}
	}
	return m, nil
}

func filterRooms(rooms []api.Room, id int) []api.Room {
	result := make([]api.Room, 0, len(rooms))
	for _, r := range rooms {
		if r.ID != id {
			result = append(result, r)
		}
	}
	return result
}

func (m RoomsModel) View() string {
	if m.loading {
		return style.Info("Loading rooms...")
	}
	if m.mode == "detail" {
		return m.viewDetail()
	}
	return m.viewList()
}

func (m RoomsModel) viewList() string {
	lines := []string{
		style.StyleHeader.Render(fmt.Sprintf("Rooms (%d)", len(m.rooms))),
		style.RenderDivider(max(80, m.width-4)),
		fmt.Sprintf("%-4s %-24s %-6s %-8s %s",
			style.StyleTableHeader.Render("ID"),
			style.StyleTableHeader.Render("Name"),
			style.StyleTableHeader.Render("Floor"),
			style.StyleTableHeader.Render("Exits"),
			style.StyleTableHeader.Render("Description"),
		),
	}
	for i, r := range m.rooms {
		rowStyle := style.StyleTableRow
		if i == m.selected {
			rowStyle = lipgloss.Style{}.Foreground(style.ColorPrimary).Bold(true)
		}
		exitCount := 0
		if r.Exits != nil {
			exitCount = len(r.Exits)
		}
		lines = append(lines, fmt.Sprintf("%-4s %-24s %-6s %-8s %s",
			rowStyle.Render(fmt.Sprintf("%d", r.ID)),
			rowStyle.Render(trunc(r.Name, 23)),
			rowStyle.Render(fmt.Sprintf("%d", r.Floor)),
			rowStyle.Render(fmt.Sprintf("%d", exitCount)),
			rowStyle.Render(trunc(r.Description, 30)),
		))
	}
	if len(m.rooms) == 0 {
		lines = append(lines, style.StyleMuted.Render("  No rooms found"))
	}
	lines = append(lines, "")
	lines = append(lines, style.StyleMuted.Render("  [Enter] view detail   [D] delete   [R] refresh   [Esc] back"))
	if m.errMsg != "" {
		lines = append(lines, "", style.Error(m.errMsg))
	}
	return strings.Join(lines, "\n")
}

func (m RoomsModel) viewDetail() string {
	if m.selected >= len(m.rooms) {
		m.mode = "list"
		return m.viewList()
	}
	r := m.rooms[m.selected]

	lines := []string{
		style.StyleHeader.Render(fmt.Sprintf("Room #%d — %s", r.ID, r.Name)),
		style.RenderDivider(max(80, m.width-4)),
		fmt.Sprintf("  %-12s %s", style.StyleLabel.Render("Name:"), style.StyleValue.Render(r.Name)),
		fmt.Sprintf("  %-12s %s", style.StyleLabel.Render("Floor:"), style.StyleValue.Render(fmt.Sprintf("%d", r.Floor))),
		fmt.Sprintf("  %-12s %s", style.StyleLabel.Render("Starting:"), style.StyleValue.Render(fmt.Sprintf("%v", r.IsStarting))),
		"",
		style.StyleLabel.Render("  Exits:"),
	}
	if r.Exits != nil && len(r.Exits) > 0 {
		for dir, dest := range r.Exits {
			lines = append(lines, fmt.Sprintf("    %s  →  Room %d",
				style.StyleValue.Render(dir),
				style.StyleValue.Render(fmt.Sprintf("%d", dest)),
			))
		}
	} else {
		lines = append(lines, "    (no exits)")
	}
	lines = append(lines, "", style.StyleMuted.Render(trunc(r.Description, 80)))
	lines = append(lines, "", style.StyleMuted.Render("  [D] delete   [Esc] back to list"))

	if m.confirmDel {
		lines = append(lines, "", style.StyleDanger.Render("  Confirm DELETE? Press [D] again to confirm"))
	}
	return strings.Join(lines, "\n")
}

type RoomsMsg struct{ Rooms []api.Room }
type RoomsErrMsg struct{ Err error }
