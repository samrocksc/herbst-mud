package screens

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"herbst-mud/admin-tui/api"
	"herbst-mud/admin-tui/style"
)

type CharactersModel struct {
	token      string
	characters []api.Character
	loading    bool
	errMsg     string
	selected   int
	page       int
	filter     string // "all", "players", "npcs"
	mode       string // "list", "detail", "edit"
	confirmDel bool
	width      int
	// edit form fields
	formName   string
	formClass  string
	formRace   string
	formLevel  string
	formHp     string
	formMaxHp  string
	formXp     string
	formRoomID string
	formDesc   string
	formField  int
	editErr    string
}

const charsPerPage = 20

var charEditFields = []string{
	"name", "class", "race", "level", "hp", "maxHp", "xp", "roomID", "description",
}

func (m CharactersModel) formFieldValue(field string) string {
	switch field {
	case "name":
		return m.formName
	case "class":
		return m.formClass
	case "race":
		return m.formRace
	case "level":
		return m.formLevel
	case "hp":
		return m.formHp
	case "maxHp":
		return m.formMaxHp
	case "xp":
		return m.formXp
	case "roomID":
		return m.formRoomID
	case "description":
		return m.formDesc
	}
	return ""
}

func (m *CharactersModel) setFormFieldValue(field, val string) {
	switch field {
	case "name":
		m.formName = val
	case "class":
		m.formClass = val
	case "race":
		m.formRace = val
	case "level":
		m.formLevel = val
	case "hp":
		m.formHp = val
	case "maxHp":
		m.formMaxHp = val
	case "xp":
		m.formXp = val
	case "roomID":
		m.formRoomID = val
	case "description":
		m.formDesc = val
	}
}

func (m CharactersModel) formToMap() map[string]any {
	body := map[string]any{
		"name":        m.formName,
		"class":       m.formClass,
		"race":        m.formRace,
		"description": m.formDesc,
	}
	if m.formLevel != "" {
		body["level"] = m.formLevel
	}
	if m.formHp != "" {
		body["hp"] = m.formHp
	}
	if m.formMaxHp != "" {
		body["max_hp"] = m.formMaxHp
	}
	if m.formXp != "" {
		body["xp"] = m.formXp
	}
	if m.formRoomID != "" {
		body["room_id"] = m.formRoomID
	}
	return body
}

func NewCharactersScreen(token string) tea.Model {
	return CharactersModel{token: token, loading: true, selected: 0, page: 0, filter: "all", mode: "list"}
}

func (m CharactersModel) Init() tea.Cmd {
	return func() tea.Msg {
		chars, err := api.ListCharacters()
		if err != nil {
			return CharsErrMsg{Err: err}
		}
		return CharsMsg{Characters: chars}
	}
}

func (m CharactersModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		return m, nil
	case CharsMsg:
		m.characters = msg.Characters
		m.loading = false
		return m, nil
	case CharsErrMsg:
		m.errMsg = fmt.Sprintf("Failed to load characters: %v", msg.Err)
		m.loading = false
		return m, nil
	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func (m CharactersModel) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Handle form input when in edit mode
	if m.mode == "edit" {
		return m.handleEditKey(msg)
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
		if m.mode == "list" && len(m.filtered()) > 0 {
			m.mode = "detail"
			return m, nil
		}
	case "e":
		if m.mode == "detail" && len(m.filtered()) > 0 {
			m.mode = "edit"
			c := m.filtered()[m.selected]
			m.formName = c.Name
			m.formClass = c.Class
			m.formRace = c.Race
			m.formLevel = fmt.Sprintf("%d", c.Level)
			m.formHp = fmt.Sprintf("%d", c.HP)
			m.formMaxHp = fmt.Sprintf("%d", c.MaxHP)
			m.formXp = fmt.Sprintf("%d", c.Xp)
			m.formRoomID = fmt.Sprintf("%d", c.RoomID)
			m.formDesc = c.Description
			m.formField = 0
			m.editErr = ""
			return m, nil
		}
	case "d":
		if m.mode == "detail" {
			if m.confirmDel {
				go func() { api.DeleteCharacter(m.filtered()[m.selected].ID) }()
				m.characters = filterChars(m.characters, m.filtered()[m.selected].ID)
				m.mode = "list"
				m.confirmDel = false
				if m.selected >= len(m.filtered()) && m.selected > 0 {
					m.selected--
				}
				return m, nil
			}
			m.confirmDel = true
			return m, nil
		}
	case "r":
		m.loading = true
		return m, func() tea.Msg {
			chars, err := api.ListCharacters()
			if err != nil {
				return CharsErrMsg{Err: err}
			}
			return CharsMsg{Characters: chars}
		}
	case "p":
		if m.mode == "list" {
			m.filter = "players"
			m.selected = 0
			m.page = 0
			return m, nil
		}
	case "n":
		if m.mode == "list" {
			m.filter = "npcs"
			m.selected = 0
			m.page = 0
			return m, nil
		}
	case "a":
		if m.mode == "list" {
			m.filter = "all"
			m.selected = 0
			m.page = 0
			return m, nil
		}
	case "up", "k":
		if m.mode == "list" && m.selected > 0 {
			m.selected--
			if m.selected < m.page*charsPerPage {
				m.page--
			}
		}
		return m, nil
	case "down", "j":
		if m.mode == "list" && m.selected < len(m.filtered())-1 {
			m.selected++
			if m.selected >= (m.page+1)*charsPerPage {
				m.page++
			}
		}
		return m, nil
	}
	if m.mode == "list" {
		switch msg.String() {
		case "u":
			return m, func() tea.Msg { return NavigateMsg{Screen: 2} }
		case "r":
			return m, func() tea.Msg { return NavigateMsg{Screen: 4} }
		case "i":
			return m, func() tea.Msg { return NavigateMsg{Screen: 6} }
		case "q":
			return m, func() tea.Msg { return NavigateMsg{Screen: 7} }
		case "b":
			return m, func() tea.Msg { return NavigateMsg{Screen: 10} }
		case "w":
			return m, func() tea.Msg { return NavigateMsg{Screen: 11} }
		}
	}
	return m, nil
}

func (m CharactersModel) handleEditKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	field := charEditFields[m.formField]

	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		m.mode = "detail"
		m.editErr = ""
		return m, nil
	case "tab":
		m.formField = (m.formField + 1) % len(charEditFields)
		return m, nil
	case "shift+tab":
		m.formField = (m.formField - 1 + len(charEditFields)) % len(charEditFields)
		return m, nil
	case "enter":
		return m.handleEditSubmit()
	case "backspace":
		cur := m.formFieldValue(field)
		if len(cur) > 0 {
			m.setFormFieldValue(field, cur[:len(cur)-1])
		}
		return m, nil
	}

	// Regular text input
	if len(msg.String()) == 1 {
		cur := m.formFieldValue(field)
		m.setFormFieldValue(field, cur+msg.String())
	}
	return m, nil
}

func (m CharactersModel) handleEditSubmit() (tea.Model, tea.Cmd) {
	if m.formName == "" {
		m.editErr = "Name is required"
		return m, nil
	}

	c := m.filtered()[m.selected]
	updated, err := api.UpdateCharacter(c.ID, m.formToMap())
	if err != nil {
		m.editErr = fmt.Sprintf("Update failed: %v", err)
		return m, nil
	}

	// Update in master list
	for i, ch := range m.characters {
		if ch.ID == updated.ID {
			m.characters[i] = updated
			break
		}
	}

	m.mode = "detail"
	m.editErr = ""
	return m, nil
}

func (m CharactersModel) filtered() []api.Character {
	switch m.filter {
	case "players":
		result := make([]api.Character, 0)
		for _, c := range m.characters {
			if !c.IsNPC {
				result = append(result, c)
			}
		}
		return result
	case "npcs":
		result := make([]api.Character, 0)
		for _, c := range m.characters {
			if c.IsNPC {
				result = append(result, c)
			}
		}
		return result
	}
	return m.characters
}

func filterChars(chars []api.Character, id int) []api.Character {
	result := make([]api.Character, 0, len(chars))
	for _, c := range chars {
		if c.ID != id {
			result = append(result, c)
		}
	}
	return result
}

func (m CharactersModel) View() string {
	if m.loading {
		return style.Info("Loading characters...")
	}
	switch m.mode {
	case "edit":
		return m.viewEdit()
	case "detail":
		return m.viewDetail()
	default:
		return m.viewList()
	}
}

func (m CharactersModel) viewList() string {
	filtered := m.filtered()
	totalPages := (len(filtered) + charsPerPage - 1) / charsPerPage
	start := m.page * charsPerPage
	end := start + charsPerPage
	if end > len(filtered) {
		end = len(filtered)
	}
	pageItems := filtered[start:end]

	filterLabel := "All"
	if m.filter == "players" {
		filterLabel = "Players"
	} else if m.filter == "npcs" {
		filterLabel = "NPCs"
	}

	lines := []string{
		style.StyleHeader.Render(fmt.Sprintf("Characters (%s — %d total)", filterLabel, len(filtered))),
		style.RenderDivider(max(90, m.width-4)),
		fmt.Sprintf("%-4s %-18s %-10s %-4s %-6s %-6s %-8s",
			style.StyleTableHeader.Render("ID"),
			style.StyleTableHeader.Render("Name"),
			style.StyleTableHeader.Render("Class"),
			style.StyleTableHeader.Render("Lvl"),
			style.StyleTableHeader.Render("HP"),
			style.StyleTableHeader.Render("Room"),
			style.StyleTableHeader.Render("Owner/NPC"),
		),
	}

	for i, c := range pageItems {
		globalIdx := start + i
		rowStyle := style.StyleTableRow
		if globalIdx == m.selected {
			rowStyle = lipgloss.Style{}.Foreground(style.ColorPrimary).Bold(true)
		}
		owner := "NPC"
		if !c.IsNPC {
			owner = fmt.Sprintf("P:%d", c.OwnerID)
		}
		lines = append(lines, fmt.Sprintf("%-4s %-18s %-10s %-4s %-6s %-6s %-8s",
			rowStyle.Render(fmt.Sprintf("%d", c.ID)),
			rowStyle.Render(trunc(c.Name, 17)),
			rowStyle.Render(c.Class),
			rowStyle.Render(fmt.Sprintf("%d", c.Level)),
			rowStyle.Render(fmt.Sprintf("%d/%d", c.HP, c.MaxHP)),
			rowStyle.Render(fmt.Sprintf("%d", c.RoomID)),
			rowStyle.Render(owner),
		))
	}

	if len(filtered) == 0 {
		lines = append(lines, style.StyleMuted.Render("  No characters found"))
	}

	if totalPages > 1 {
		lines = append(lines, "", style.StyleMuted.Render(fmt.Sprintf("  Page %d/%d", m.page+1, totalPages)))
	}
	lines = append(lines, "")
	lines = append(lines, style.StyleMuted.Render("  [P] players   [N] NPCs   [A] all   [Enter] view   [D] delete   [R] refresh   [Esc] back"))

	if m.errMsg != "" {
		lines = append(lines, "", style.Error(m.errMsg))
	}
	return strings.Join(lines, "\n")
}

func (m CharactersModel) viewDetail() string {
	filtered := m.filtered()
	if m.selected >= len(filtered) {
		m.mode = "list"
		return m.viewList()
	}
	c := filtered[m.selected]
	isNPC := c.IsNPC

	lines := []string{
		style.StyleHeader.Render(fmt.Sprintf("Character #%d — %s", c.ID, c.Name)),
		style.RenderDivider(max(90, m.width-4)),
		fmt.Sprintf("  %-14s %s", style.StyleLabel.Render("Name:"), style.StyleValue.Render(c.Name)),
		fmt.Sprintf("  %-14s %s", style.StyleLabel.Render("Class:"), style.StyleValue.Render(c.Class)),
		fmt.Sprintf("  %-14s %s", style.StyleLabel.Render("Race:"), style.StyleValue.Render(c.Race)),
		fmt.Sprintf("  %-14s %s", style.StyleLabel.Render("Level:"), style.StyleValue.Render(fmt.Sprintf("%d", c.Level))),
		fmt.Sprintf("  %-14s %s", style.StyleLabel.Render("XP:"), style.StyleValue.Render(fmt.Sprintf("%d", c.Xp))),
		fmt.Sprintf("  %-14s %s", style.StyleLabel.Render("HP:"), style.StyleValue.Render(fmt.Sprintf("%d/%d", c.HP, c.MaxHP))),
		fmt.Sprintf("  %-14s %s", style.StyleLabel.Render("Room:"), style.StyleValue.Render(fmt.Sprintf("%d", c.RoomID))),
	}
	if isNPC {
		lines = append(lines,
			fmt.Sprintf("  %-14s %s", style.StyleLabel.Render("Behavior:"), style.StyleValue.Render(c.Behavior)),
			fmt.Sprintf("  %-14s %s", style.StyleLabel.Render("Aggression:"), style.StyleValue.Render(c.Aggression)),
		)
	} else {
		lines = append(lines,
			fmt.Sprintf("  %-14s %s", style.StyleLabel.Render("Owner ID:"), style.StyleValue.Render(fmt.Sprintf("%d", c.OwnerID))),
		)
	}
	lines = append(lines,
		fmt.Sprintf("  %-14s %s", style.StyleLabel.Render("Description:"), style.StyleValue.Render(trunc(c.Description, 60))),
		"",
		style.StyleMuted.Render("  [E] edit   [D] delete   [Esc] back to list"),
	)
	if m.confirmDel {
		lines = append(lines, "", style.StyleDanger.Render("  Confirm DELETE? Press [D] again to confirm"))
	}
	return strings.Join(lines, "\n")
}

func (m CharactersModel) viewEdit() string {
	fieldLabels := []struct {
		label string
		field string
	}{
		{"Name:", "name"},
		{"Class:", "class"},
		{"Race:", "race"},
		{"Level:", "level"},
		{"HP:", "hp"},
		{"Max HP:", "maxHp"},
		{"XP:", "xp"},
		{"Room ID:", "roomID"},
		{"Description:", "description"},
	}

	lines := []string{
		style.StyleHeader.Render("Edit Character"),
		style.RenderDivider(max(90, m.width-4)),
		style.StyleMuted.Render("  [Tab] next field   [Enter] submit   [Esc] cancel"),
		"",
	}

	for i, fl := range fieldLabels {
		val := m.formFieldValue(fl.field)
		cursor := "  "
		if i == m.formField {
			cursor = "▸ "
		}
		lines = append(lines, fmt.Sprintf("%s%-14s %s", cursor, style.StyleLabel.Render(fl.label), style.StyleValue.Render(val)))
	}

	if m.editErr != "" {
		lines = append(lines, "", style.Error(m.editErr))
	}

	return strings.Join(lines, "\n")
}

type CharsMsg struct{ Characters []api.Character }
type CharsErrMsg struct{ Err error }