package screens

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"herbst-mud/admin-tui/api"
	"herbst-mud/admin-tui/style"
)

type RoomsModel struct {
	token      string
	rooms      []api.Room
	loading    bool
	errMsg     string
	selected   int
	mode       string // "list", "detail", "create"
	confirmDel bool
	width      int
	form       roomForm
	formField  int // index into roomFormFields
	createErr  string
}

type roomForm struct {
	Name         string
	Description  string
	Floor        string
	IsStarting   bool
	ExitNorth    string
	ExitSouth    string
	ExitEast     string
	ExitWest     string
	ExitUp       string
	ExitDown     string
}

var roomFormFields = []string{
	"name", "description", "floor", "isStarting",
	"exitNorth", "exitSouth", "exitEast", "exitWest", "exitUp", "exitDown",
}

func (f roomForm) reset() roomForm {
	return roomForm{}
}

func (f roomForm) toMap() map[string]any {
	exits := map[string]int{}
	dirs := map[string]string{
		"north": f.ExitNorth,
		"south": f.ExitSouth,
		"east":  f.ExitEast,
		"west":  f.ExitWest,
		"up":    f.ExitUp,
		"down":  f.ExitDown,
	}
	for dir, val := range dirs {
		if val == "" {
			continue
		}
		id, err := strconv.Atoi(val)
		if err == nil && id > 0 {
			exits[dir] = id
		}
	}
	body := map[string]any{
		"name":           f.Name,
		"description":    f.Description,
		"isStartingRoom": f.IsStarting,
		"exits":          exits,
	}
	if f.Floor != "" {
		floor, err := strconv.Atoi(f.Floor)
		if err == nil {
			body["floor"] = floor
		}
	}
	return body
}

func (f roomForm) fieldValue(field string) string {
	switch field {
	case "name":
		return f.Name
	case "description":
		return f.Description
	case "floor":
		return f.Floor
	case "isStarting":
		return fmt.Sprintf("%v", f.IsStarting)
	case "exitNorth":
		return f.ExitNorth
	case "exitSouth":
		return f.ExitSouth
	case "exitEast":
		return f.ExitEast
	case "exitWest":
		return f.ExitWest
	case "exitUp":
		return f.ExitUp
	case "exitDown":
		return f.ExitDown
	}
	return ""
}

func (f *roomForm) setFieldValue(field, val string) {
	switch field {
	case "name":
		f.Name = val
	case "description":
		f.Description = val
	case "floor":
		f.Floor = val
	case "exitNorth":
		f.ExitNorth = val
	case "exitSouth":
		f.ExitSouth = val
	case "exitEast":
		f.ExitEast = val
	case "exitWest":
		f.ExitWest = val
	case "exitUp":
		f.ExitUp = val
	case "exitDown":
		f.ExitDown = val
	}
}

func (f *roomForm) toggleBool(field string) {
	if field == "isStarting" {
		f.IsStarting = !f.IsStarting
	}
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
	// Handle form input when in create mode
	if m.mode == "create" {
		return m.handleCreateKey(msg)
	}

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
		if m.mode == "list" && len(m.rooms) > 0 {
			m.mode = "detail"
			return m, nil
		}
	case "c":
		if m.mode == "list" {
			m.mode = "create"
			m.form = roomForm.reset(m.form)
			m.formField = 0
			m.createErr = ""
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
		case "u":
			return m, func() tea.Msg { return NavigateMsg{Screen: 2} }
		case "c":
			return m, func() tea.Msg { return NavigateMsg{Screen: 3} }
		case "n":
			return m, func() tea.Msg { return NavigateMsg{Screen: 5} }
		case "i":
			return m, func() tea.Msg { return NavigateMsg{Screen: 6} }
		case "q":
			return m, func() tea.Msg { return NavigateMsg{Screen: 7} }
		case "b":
			return m, func() tea.Msg { return NavigateMsg{Screen: 8} }
		case "w":
			return m, func() tea.Msg { return NavigateMsg{Screen: 9} }
		}
	}
	return m, nil
}

func (m RoomsModel) handleCreateKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	field := roomFormFields[m.formField]

	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		m.mode = "list"
		m.createErr = ""
		return m, nil
	case "tab":
		m.formField = (m.formField + 1) % len(roomFormFields)
		return m, nil
	case "shift+tab":
		m.formField = (m.formField - 1 + len(roomFormFields)) % len(roomFormFields)
		return m, nil
	case "enter":
		return m.handleCreateSubmit()
	case "backspace":
		cur := m.form.fieldValue(field)
		if len(cur) > 0 && field != "isStarting" {
			m.form.setFieldValue(field, cur[:len(cur)-1])
		}
		return m, nil
	}

	// Toggle boolean field
	if field == "isStarting" {
		if msg.String() == "y" || msg.String() == "t" {
			m.form.IsStarting = true
		} else if msg.String() == "n" || msg.String() == "f" {
			m.form.IsStarting = false
		}
		return m, nil
	}

	// Regular text input
	if len(msg.String()) == 1 {
		cur := m.form.fieldValue(field)
		m.form.setFieldValue(field, cur+msg.String())
	}
	return m, nil
}

func (m RoomsModel) handleCreateSubmit() (tea.Model, tea.Cmd) {
	if m.form.Name == "" {
		m.createErr = "Name is required"
		return m, nil
	}

	created, err := api.CreateRoom(m.form.toMap())
	if err != nil {
		m.createErr = fmt.Sprintf("Create failed: %v", err)
		return m, nil
	}

	m.rooms = append(m.rooms, created)
	m.selected = len(m.rooms) - 1
	m.mode = "detail"
	m.createErr = ""
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
	switch m.mode {
	case "detail":
		return m.viewDetail()
	case "create":
		return m.viewCreate()
	default:
		return m.viewList()
	}
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
	lines = append(lines, style.StyleMuted.Render("  [Enter] view detail   [C] create   [D] delete   [R] refresh   [Esc] back"))
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
			lines = append(lines, fmt.Sprintf("    %s  →  Room %s",
				style.StyleValue.Render(dir),
				style.StyleValue.Render(fmt.Sprintf("%d", dest)),
			))
		}
	} else {
		lines = append(lines, "    (no exits)")
	}
	lines = append(lines, "", style.StyleMuted.Render(trunc(r.Description, 80)))
	lines = append(lines, "", style.StyleMuted.Render("  [C] create   [D] delete   [Esc] back to list"))

	if m.confirmDel {
		lines = append(lines, "", style.StyleDanger.Render("  Confirm DELETE? Press [D] again to confirm"))
	}
	return strings.Join(lines, "\n")
}

func (m RoomsModel) viewCreate() string {
	fieldLabels := []struct {
		label string
		field string
	}{
		{"Name:", "name"},
		{"Description:", "description"},
		{"Floor:", "floor"},
		{"Is Starting:", "isStarting"},
		{"Exit North:", "exitNorth"},
		{"Exit South:", "exitSouth"},
		{"Exit East:", "exitEast"},
		{"Exit West:", "exitWest"},
		{"Exit Up:", "exitUp"},
		{"Exit Down:", "exitDown"},
	}

	lines := []string{
		style.StyleHeader.Render("Create Room"),
		style.RenderDivider(max(80, m.width-4)),
		style.StyleMuted.Render("  [Tab] next field   [Shift+Tab] prev   [Enter] submit   [Esc] cancel"),
		"",
	}

	for i, fl := range fieldLabels {
		val := m.form.fieldValue(fl.field)
		cursor := "  "
		if i == m.formField {
			cursor = "▸ "
		}
		lines = append(lines, fmt.Sprintf("%s%-14s %s", cursor, style.StyleLabel.Render(fl.label), style.StyleValue.Render(val)))
	}

	if m.createErr != "" {
		lines = append(lines, "", style.Error(m.createErr))
	}

	return strings.Join(lines, "\n")
}

type RoomsMsg struct{ Rooms []api.Room }
type RoomsErrMsg struct{ Err error }