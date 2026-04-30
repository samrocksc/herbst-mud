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
	token         string
	npcs          []api.Character
	npcTemplates  []api.NPCTemplate
	loading       bool
	errMsg        string
	selected      int
	mode          string // "list", "detail", "create", "edit", "editXp"
	confirmDel    bool
	width         int
	form          npcForm
	formField     int
	createErr     string
	xpEditValue   string
	xpEditErr     string
	pendingXpID   string
	pendingXpName string
}

type npcForm struct {
	Name        string
	Description string
	Class       string
	Race        string
	Behavior    string
	Aggression  string
	Hp          string
	Level       string
	RoomID      string
}

var npcFormFields = []string{
	"name", "description", "class", "race", "behavior", "aggression", "hp", "level", "roomID",
}

func (f npcForm) reset() npcForm {
	return npcForm{}
}

func (f npcForm) fieldValue(field string) string {
	switch field {
	case "name":
		return f.Name
	case "description":
		return f.Description
	case "class":
		return f.Class
	case "race":
		return f.Race
	case "behavior":
		return f.Behavior
	case "aggression":
		return f.Aggression
	case "hp":
		return f.Hp
	case "level":
		return f.Level
	case "roomID":
		return f.RoomID
	}
	return ""
}

func (f *npcForm) setFieldValue(field, val string) {
	switch field {
	case "name":
		f.Name = val
	case "description":
		f.Description = val
	case "class":
		f.Class = val
	case "race":
		f.Race = val
	case "behavior":
		f.Behavior = val
	case "aggression":
		f.Aggression = val
	case "hp":
		f.Hp = val
	case "level":
		f.Level = val
	case "roomID":
		f.RoomID = val
	}
}

func (f npcForm) toMap() map[string]any {
	body := map[string]any{
		"name":        f.Name,
		"description": f.Description,
		"class":       f.Class,
		"race":        f.Race,
		"behavior":    f.Behavior,
		"aggression":  f.Aggression,
		"isNPC":       true,
	}
	if f.Hp != "" {
		body["hp"] = f.Hp
	}
	if f.Level != "" {
		body["level"] = f.Level
	}
	if f.RoomID != "" {
		body["currentRoomId"] = f.RoomID
	}
	return body
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
	case NPCTemplatesMsg:
		m.npcTemplates = msg.Templates
		return m, nil
	case NPCTemplateXpUpdatedMsg:
		// Refresh templates after update
		return m, func() tea.Msg {
			templates, err := api.GetNPCTemplates()
			if err != nil {
				return NPCsErrMsg{Err: err}
			}
			return NPCTemplatesMsg{Templates: templates}
		}
	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func (m NPCsModel) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Handle XP edit mode
	if m.mode == "editXp" {
		return m.handleXpEditKey(msg)
	}

	// Handle form input when in create or edit mode
	if m.mode == "create" {
		return m.handleCreateKey(msg)
	}
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
		if m.mode == "list" && len(m.npcs) > 0 && m.selected < len(m.npcs) {
			m.mode = "detail"
			return m, nil
		}
	case "e":
		if m.mode == "detail" && m.selected < len(m.npcs) {
			m.mode = "edit"
			n := m.npcs[m.selected]
			m.form = npcForm{
				Name:        n.Name,
				Description: n.Description,
				Class:       n.Class,
				Race:        n.Race,
				Behavior:    n.Behavior,
				Aggression:  n.Aggression,
				Hp:          fmt.Sprintf("%d", n.HP),
				Level:       fmt.Sprintf("%d", n.Level),
				RoomID:      fmt.Sprintf("%d", n.RoomID),
			}
			m.formField = 0
			m.createErr = ""
			return m, nil
		}
	case "x":
		if m.mode == "detail" {
			// Enter mini-edit mode for XP template value
			// Refresh NPC templates first
			cmd := func() tea.Msg {
				templates, err := api.GetNPCTemplates()
				if err != nil {
					return NPCsErrMsg{Err: err}
				}
				return NPCTemplatesMsg{Templates: templates}
			}
			m.mode = "editXp"
			m.xpEditValue = ""
			m.xpEditErr = ""
			m.pendingXpID = ""
			m.pendingXpName = ""
			return m, cmd
		}
	case "c":
		if m.mode == "list" {
			m.mode = "create"
			m.form = npcForm.reset(m.form)
			m.formField = 0
			m.createErr = ""
			return m, nil
		}
	case "d":
		if m.mode == "detail" {
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
		case "u":
			return m, func() tea.Msg { return NavigateMsg{Screen: 2} }
		case "c":
			return m, func() tea.Msg { return NavigateMsg{Screen: 3} }
		case "r":
			return m, func() tea.Msg { return NavigateMsg{Screen: 4} }
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

func (m NPCsModel) handleXpEditKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		m.mode = "detail"
		m.xpEditErr = ""
		return m, nil
	case "enter":
		return m.handleXpEditSubmit()
	case "backspace":
		if len(m.xpEditValue) > 0 {
			m.xpEditValue = m.xpEditValue[:len(m.xpEditValue)-1]
		}
		return m, nil
	}

	// Regular text input (digits only for XP)
	if len(msg.String()) == 1 && msg.String() >= "0" && msg.String() <= "9" {
		m.xpEditValue += msg.String()
	}
	return m, nil
}

func (m NPCsModel) handleXpEditSubmit() (tea.Model, tea.Cmd) {
	if m.pendingXpID == "" {
		// User hasn't selected a template yet — this shouldn't happen
		// since we auto-select the first matching template, but handle gracefully
		m.xpEditErr = "No template selected"
		return m, nil
	}
	if m.xpEditValue == "" {
		m.xpEditErr = "XP value is required"
		return m, nil
	}

	// Parse the XP value (already validated as digits only)
	var xpVal int
	fmt.Sscanf(m.xpEditValue, "%d", &xpVal)

	_, err := api.UpdateNPCTemplate(m.pendingXpID, xpVal)
	if err != nil {
		m.xpEditErr = fmt.Sprintf("Update failed: %v", err)
		return m, nil
	}

	m.mode = "detail"
	m.xpEditErr = ""
	return m, func() tea.Msg {
		templates, err := api.GetNPCTemplates()
		if err != nil {
			return NPCsErrMsg{Err: err}
		}
		return NPCTemplatesMsg{Templates: templates}
	}
}

func (m NPCsModel) handleCreateKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	field := npcFormFields[m.formField]

	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		m.mode = "list"
		m.createErr = ""
		return m, nil
	case "tab":
		m.formField = (m.formField + 1) % len(npcFormFields)
		return m, nil
	case "shift+tab":
		m.formField = (m.formField - 1 + len(npcFormFields)) % len(npcFormFields)
		return m, nil
	case "enter":
		return m.handleCreateSubmit()
	case "backspace":
		cur := m.form.fieldValue(field)
		if len(cur) > 0 {
			m.form.setFieldValue(field, cur[:len(cur)-1])
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

func (m NPCsModel) handleCreateSubmit() (tea.Model, tea.Cmd) {
	if m.form.Name == "" {
		m.createErr = "Name is required"
		return m, nil
	}

	created, err := api.CreateNPC(m.form.toMap())
	if err != nil {
		m.createErr = fmt.Sprintf("Create failed: %v", err)
		return m, nil
	}

	m.npcs = append(m.npcs, created)
	m.selected = len(m.npcs) - 1
	m.mode = "detail"
	m.createErr = ""
	return m, nil
}

func (m NPCsModel) handleEditKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	field := npcFormFields[m.formField]

	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		m.mode = "detail"
		m.createErr = ""
		return m, nil
	case "tab":
		m.formField = (m.formField + 1) % len(npcFormFields)
		return m, nil
	case "shift+tab":
		m.formField = (m.formField - 1 + len(npcFormFields)) % len(npcFormFields)
		return m, nil
	case "enter":
		return m.handleEditSubmit()
	case "backspace":
		cur := m.form.fieldValue(field)
		if len(cur) > 0 {
			m.form.setFieldValue(field, cur[:len(cur)-1])
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

func (m NPCsModel) handleEditSubmit() (tea.Model, tea.Cmd) {
	if m.form.Name == "" {
		m.createErr = "Name is required"
		return m, nil
	}

	n := m.npcs[m.selected]
	updated, err := api.UpdateNPC(n.ID, m.form.toMap())
	if err != nil {
		m.createErr = fmt.Sprintf("Update failed: %v", err)
		return m, nil
	}

	// Update in master list
	for i, c := range m.npcs {
		if c.ID == updated.ID {
			m.npcs[i] = updated
			break
		}
	}

	m.mode = "detail"
	m.createErr = ""
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
	switch m.mode {
	case "editXp":
		return m.viewEditXp()
	case "edit":
		return m.viewEdit()
	case "detail":
		if m.selected < len(m.npcs) {
			return m.viewDetail()
		}
		m.mode = "list"
		return m.viewList()
	case "create":
		return m.viewCreate()
	default:
		return m.viewList()
	}
}

func (m NPCsModel) viewList() string {
	lines := []string{
		style.StyleHeader.Render(fmt.Sprintf("NPCs (%d)", len(m.npcs))),
		style.RenderDivider(max(90, m.width-4)),
		fmt.Sprintf("%-4s %-20s %-10s %-4s %-6s %-6s %-8s %-10s %-10s",
			style.StyleTableHeader.Render("ID"),
			style.StyleTableHeader.Render("Name"),
			style.StyleTableHeader.Render("Class"),
			style.StyleTableHeader.Render("Lvl"),
			style.StyleTableHeader.Render("XP"),
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
		// Show xp from the Character struct
		xpStr := "-"
		if n.Xp > 0 {
			xpStr = fmt.Sprintf("%d", n.Xp)
		}
		lines = append(lines, fmt.Sprintf("%-4s %-20s %-10s %-4s %-6s %-6s %-8s %-10s %-10s",
			rowStyle.Render(fmt.Sprintf("%d", n.ID)),
			rowStyle.Render(trunc(n.Name, 19)),
			rowStyle.Render(n.Class),
			rowStyle.Render(fmt.Sprintf("%d", n.Level)),
			rowStyle.Render(xpStr),
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
	lines = append(lines, style.StyleMuted.Render("  [Enter] view detail   [C] create   [D] delete   [R] refresh   [Esc] back"))
	if m.errMsg != "" {
		lines = append(lines, "", style.Error(m.errMsg))
	}
	return strings.Join(lines, "\n")
}

func (m NPCsModel) viewDetail() string {
	n := m.npcs[m.selected]

	// Find matching NPC template to show xp_value
	var templateXp string
	for _, t := range m.npcTemplates {
		if t.Name == n.Name || t.ID == n.Name {
			templateXp = fmt.Sprintf("%d", t.XpValue)
			break
		}
	}
	if templateXp == "" {
		templateXp = "—"
	}

	lines := []string{
		style.StyleHeader.Render(fmt.Sprintf("NPC #%d — %s", n.ID, n.Name)),
		style.RenderDivider(max(90, m.width-4)),
		fmt.Sprintf("  %-12s %s", style.StyleLabel.Render("Name:"), style.StyleValue.Render(n.Name)),
		fmt.Sprintf("  %-12s %s", style.StyleLabel.Render("Class:"), style.StyleValue.Render(n.Class)),
		fmt.Sprintf("  %-12s %s", style.StyleLabel.Render("Race:"), style.StyleValue.Render(n.Race)),
		fmt.Sprintf("  %-12s %s", style.StyleLabel.Render("Level:"), style.StyleValue.Render(fmt.Sprintf("%d", n.Level))),
		fmt.Sprintf("  %-12s %s", style.StyleLabel.Render("XP:"), style.StyleValue.Render(fmt.Sprintf("%d", n.Xp))),
		fmt.Sprintf("  %-12s %s", style.StyleLabel.Render("HP:"), style.StyleValue.Render(fmt.Sprintf("%d/%d", n.HP, n.MaxHP))),
		fmt.Sprintf("  %-12s %s", style.StyleLabel.Render("Room:"), style.StyleValue.Render(fmt.Sprintf("%d", n.RoomID))),
		fmt.Sprintf("  %-12s %s", style.StyleLabel.Render("Behavior:"), style.StyleValue.Render(n.Behavior)),
		fmt.Sprintf("  %-12s %s", style.StyleLabel.Render("Aggression:"), style.StyleValue.Render(n.Aggression)),
		fmt.Sprintf("  %-12s %s", style.StyleLabel.Render("Description:"), style.StyleValue.Render(trunc(n.Description, 60))),
	}

	// Show template XP info if available
	if templateXp != "—" {
		lines = append(lines, "")
		lines = append(lines, style.StyleLabel.Render("  Template XP Value: ")+style.StyleValue.Render(templateXp))
	}

	lines = append(lines, "")
	lines = append(lines, style.StyleMuted.Render("  [E] edit   [X] edit template XP   [C] create   [D] delete   [Esc] back to list"))
	if m.confirmDel {
		lines = append(lines, "", style.StyleDanger.Render("  Confirm DELETE? Press [D] again to confirm"))
	}
	return strings.Join(lines, "\n")
}

func (m NPCsModel) viewEditXp() string {
	n := m.npcs[m.selected]

	// Find matching template for this NPC
	var matchedTemplate *api.NPCTemplate
	for i, t := range m.npcTemplates {
		if t.Name == n.Name || t.ID == n.Name {
			matchedTemplate = &m.npcTemplates[i]
			break
		}
	}

	// If no match by name, show template picker
	if matchedTemplate == nil && len(m.npcTemplates) > 0 {
		// Auto-select first template if only one, or let user pick
		if len(m.npcTemplates) == 1 {
			matchedTemplate = &m.npcTemplates[0]
		}
	}

	// Set pending XP ID if we found a template
	if matchedTemplate != nil && m.pendingXpID == "" {
		m.pendingXpID = matchedTemplate.ID
		m.pendingXpName = matchedTemplate.Name
		if m.xpEditValue == "" {
			m.xpEditValue = fmt.Sprintf("%d", matchedTemplate.XpValue)
		}
	}

	lines := []string{
		style.StyleHeader.Render(fmt.Sprintf("Edit Template XP — %s", n.Name)),
		style.RenderDivider(max(90, m.width-4)),
	}

	if matchedTemplate != nil {
		lines = append(lines, fmt.Sprintf("  Template: %s (ID: %s)", style.StyleValue.Render(matchedTemplate.Name), style.StyleMuted.Render(matchedTemplate.ID)))
		lines = append(lines, fmt.Sprintf("  Current XP Value: %s", style.StyleValue.Render(fmt.Sprintf("%d", matchedTemplate.XpValue))))
	} else if len(m.npcTemplates) == 0 {
		lines = append(lines, style.StyleMuted.Render("  No NPC templates found. Press [Esc] to go back."))
		lines = append(lines, "")
		lines = append(lines, style.StyleMuted.Render("  [Esc] cancel"))
		return strings.Join(lines, "\n")
	} else {
		// No matching template found — show available templates
		lines = append(lines, style.StyleMuted.Render("  No matching template found. Available templates:"))
		for _, t := range m.npcTemplates {
			lines = append(lines, fmt.Sprintf("    • %s (XP: %d)", t.Name, t.XpValue))
		}
	}

	if matchedTemplate != nil {
		lines = append(lines, "")
		lines = append(lines, fmt.Sprintf("  New XP Value: %s█", style.StyleValue.Render(m.xpEditValue)))
		lines = append(lines, "")
		lines = append(lines, style.StyleMuted.Render("  [Enter] save   [Esc] cancel"))
	}

	if m.xpEditErr != "" {
		lines = append(lines, "", style.Error(m.xpEditErr))
	}

	return strings.Join(lines, "\n")
}

func (m NPCsModel) viewCreate() string {
	fieldLabels := []struct {
		label string
		field string
	}{
		{"Name:", "name"},
		{"Description:", "description"},
		{"Class:", "class"},
		{"Race:", "race"},
		{"Behavior:", "behavior"},
		{"Aggression:", "aggression"},
		{"HP:", "hp"},
		{"Level:", "level"},
		{"Room ID:", "roomID"},
	}

	lines := []string{
		style.StyleHeader.Render("Create NPC"),
		style.RenderDivider(max(90, m.width-4)),
		style.StyleMuted.Render("  [Tab] next field   [Enter] submit   [Esc] cancel"),
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

func (m NPCsModel) viewEdit() string {
	fieldLabels := []struct {
		label string
		field string
	}{
		{"Name:", "name"},
		{"Description:", "description"},
		{"Class:", "class"},
		{"Race:", "race"},
		{"Behavior:", "behavior"},
		{"Aggression:", "aggression"},
		{"HP:", "hp"},
		{"Level:", "level"},
		{"Room ID:", "roomID"},
	}

	lines := []string{
		style.StyleHeader.Render("Edit NPC"),
		style.RenderDivider(max(90, m.width-4)),
		style.StyleMuted.Render("  [Tab] next field   [Enter] submit   [Esc] cancel"),
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

type NPCsMsg struct{ Characters []api.Character }
type NPCsErrMsg struct{ Err error }
type NPCTemplatesMsg struct{ Templates []api.NPCTemplate }
type NPCTemplateXpUpdatedMsg struct{ ID string; XpValue int }