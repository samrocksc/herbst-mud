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

// SkillsModel is the skill management screen.
type SkillsModel struct {
	token    string
	skills   []api.SkillRecord
	loading  bool
	errMsg   string
	selected int
	mode     string // "list" | "edit"
	width    int
	// edit form fields
	editField       int // index into skillEditFields
	editSkillClass  string
	editRequiredTag string
	editProcChance  string
	editProcEvent   string
	editCooldownSec string
	editErr         string
}

var skillEditFields = []string{
	"skill_class", "required_tag", "proc_chance", "proc_event", "cooldown_seconds",
}

func (m SkillsModel) editFieldValue(field string) string {
	switch field {
	case "skill_class":
		return m.editSkillClass
	case "required_tag":
		return m.editRequiredTag
	case "proc_chance":
		return m.editProcChance
	case "proc_event":
		return m.editProcEvent
	case "cooldown_seconds":
		return m.editCooldownSec
	}
	return ""
}

func (m *SkillsModel) setEditFieldValue(field, val string) {
	switch field {
	case "skill_class":
		m.editSkillClass = val
	case "required_tag":
		m.editRequiredTag = val
	case "proc_chance":
		m.editProcChance = val
	case "proc_event":
		m.editProcEvent = val
	case "cooldown_seconds":
		m.editCooldownSec = val
	}
}

func NewSkillsScreen(token string) tea.Model {
	return SkillsModel{token: token, loading: true, selected: 0, mode: "list"}
}

func (m SkillsModel) Init() tea.Cmd {
	return func() tea.Msg {
		skills, err := api.ListSkills()
		if err != nil {
			return SkillsErrMsg{Err: err}
		}
		return SkillsMsg{Skills: skills}
	}
}

func (m SkillsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		return m, nil
	case SkillsMsg:
		m.skills = msg.Skills
		m.loading = false
		return m, nil
	case SkillsErrMsg:
		m.errMsg = fmt.Sprintf("Failed to load skills: %v", msg.Err)
		m.loading = false
		return m, nil
	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func (m SkillsModel) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.mode == "edit" {
		return m.handleEditKey(msg)
	}

	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		return m, func() tea.Msg { return NavigateMsg{Screen: 1} }
	case "up", "k":
		if m.selected > 0 {
			m.selected--
		}
		return m, nil
	case "down", "j":
		if m.selected < len(m.skills)-1 {
			m.selected++
		}
		return m, nil
	case "enter":
		if len(m.skills) > 0 {
			m.mode = "detail"
		}
		return m, nil
	case "e":
		if len(m.skills) > 0 {
			return m.startEdit()
		}
		return m, nil
	case "r":
		m.loading = true
		return m, func() tea.Msg {
			skills, err := api.ListSkills()
			if err != nil {
				return SkillsErrMsg{Err: err}
			}
			return SkillsMsg{Skills: skills}
		}
	}
	return m, nil
}

func (m SkillsModel) startEdit() (tea.Model, tea.Cmd) {
	if m.selected >= len(m.skills) {
		return m, nil
	}
	s := m.skills[m.selected]
	m.editSkillClass = s.SkillClass
	m.editRequiredTag = s.RequiredTag
	m.editProcChance = fmt.Sprintf("%.2f", s.ProcChance)
	m.editProcEvent = s.ProcEvent
	m.editCooldownSec = strconv.Itoa(s.CooldownSeconds)
	m.editField = 0
	m.editErr = ""
	m.mode = "edit"
	return m, nil
}

func (m SkillsModel) handleEditKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	field := skillEditFields[m.editField]

	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		m.mode = "list"
		m.editErr = ""
		return m, nil
	case "tab":
		m.editField = (m.editField + 1) % len(skillEditFields)
		return m, nil
	case "shift+tab":
		m.editField = (m.editField - 1 + len(skillEditFields)) % len(skillEditFields)
		return m, nil
	case "enter":
		return m.handleEditSubmit()
	case "backspace":
		cur := m.editFieldValue(field)
		if len(cur) > 0 {
			m.setEditFieldValue(field, cur[:len(cur)-1])
		}
		return m, nil
	}

	// Regular text input
	if len(msg.String()) == 1 {
		cur := m.editFieldValue(field)
		m.setEditFieldValue(field, cur+msg.String())
	}
	return m, nil
}

func (m SkillsModel) handleEditSubmit() (tea.Model, tea.Cmd) {
	if m.selected >= len(m.skills) {
		return m, nil
	}
	s := m.skills[m.selected]

	body := map[string]any{
		"skill_class":  m.editSkillClass,
		"required_tag":  m.editRequiredTag,
		"proc_event":    m.editProcEvent,
	}

	if pc, err := strconv.ParseFloat(m.editProcChance, 64); err == nil {
		body["proc_chance"] = pc
	}
	if cs, err := strconv.Atoi(m.editCooldownSec); err == nil {
		body["cooldown_seconds"] = cs
	}

	updated, err := api.UpdateSkill(s.ID, body)
	if err != nil {
		m.editErr = fmt.Sprintf("Update failed: %v", err)
		return m, nil
	}

	// Update local state
	m.skills[m.selected] = updated
	m.mode = "list"
	m.editErr = ""
	return m, nil
}

func (m SkillsModel) View() string {
	if m.loading {
		return style.Info("Loading skills...")
	}
	switch m.mode {
	case "edit":
		return m.viewEdit()
	}
	return m.viewList()
}

func (m SkillsModel) viewList() string {
	lines := []string{
		style.StyleHeader.Render(fmt.Sprintf("Skills (%d)", len(m.skills))),
		style.RenderDivider(max(100, m.width-4)),
		fmt.Sprintf("%-4s %-20s %-16s %-10s %-12s %-10s %-12s",
			style.StyleTableHeader.Render("ID"),
			style.StyleTableHeader.Render("Name"),
			style.StyleTableHeader.Render("Slug"),
			style.StyleTableHeader.Render("Class"),
			style.StyleTableHeader.Render("ReqTag"),
			style.StyleTableHeader.Render("Proc%"),
			style.StyleTableHeader.Render("Cooldown"),
		),
	}
	for i, s := range m.skills {
		rowStyle := style.StyleTableRow
		if i == m.selected {
			rowStyle = lipgloss.Style{}.Foreground(style.ColorPrimary).Bold(true)
		}
		procStr := fmt.Sprintf("%.0f%%", s.ProcChance*100)
		cdStr := fmt.Sprintf("%ds", s.CooldownSeconds)
		lines = append(lines, fmt.Sprintf("%-4s %-20s %-16s %-10s %-12s %-10s %-12s",
			rowStyle.Render(fmt.Sprintf("%d", s.ID)),
			rowStyle.Render(trunc(s.Name, 19)),
			rowStyle.Render(trunc(s.Slug, 15)),
			rowStyle.Render(trunc(s.SkillClass, 9)),
			rowStyle.Render(trunc(s.RequiredTag, 11)),
			rowStyle.Render(procStr),
			rowStyle.Render(cdStr),
		))
	}
	if len(m.skills) == 0 {
		lines = append(lines, style.StyleMuted.Render("  No skills found"))
	}
	lines = append(lines, "")
	lines = append(lines, style.StyleMuted.Render("  [Enter] view   [E] edit faction fields   [R] refresh   [Esc] back"))
	if m.errMsg != "" {
		lines = append(lines, "", style.Error(m.errMsg))
	}
	return strings.Join(lines, "\n")
}

func (m SkillsModel) viewEdit() string {
	if m.selected >= len(m.skills) {
		return m.viewList()
	}
	s := m.skills[m.selected]

	fieldLabels := []struct {
		label string
		field string
	}{
		{"Skill Class:", "skill_class"},
		{"Required Tag:", "required_tag"},
		{"Proc Chance:", "proc_chance"},
		{"Proc Event:", "proc_event"},
		{"Cooldown (s):", "cooldown_seconds"},
	}

	lines := []string{
		style.StyleHeader.Render(fmt.Sprintf("Edit Skill #%d — %s", s.ID, s.Name)),
		style.RenderDivider(max(100, m.width-4)),
		style.StyleMuted.Render("  [Tab] next field   [Enter] submit   [Esc] cancel"),
		"",
	}

	for i, fl := range fieldLabels {
		val := m.editFieldValue(fl.field)
		cursor := "  "
		if i == m.editField {
			cursor = "▸ "
		}
		lines = append(lines, fmt.Sprintf("%s%-16s %s", cursor, style.StyleLabel.Render(fl.label), style.StyleValue.Render(val)))
	}

	if m.editErr != "" {
		lines = append(lines, "", style.Error(m.editErr))
	}

	return strings.Join(lines, "\n")
}

// ─── Skills message types ──────────────────────────────────────────────────

type SkillsMsg struct {
	Skills []api.SkillRecord
}

type SkillsErrMsg struct {
	Err error
}