package screens

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"herbst-mud/admin-tui/api"
	"herbst-mud/admin-tui/style"
)

type FactionsModel struct {
	token       string
	factions    []api.Faction
	categories  []api.FactionCategory
	characters  []api.Character
	loading     bool
	errMsg      string
	selected    int
	selectedChar int
	mode        string // "list" | "detail" | "create" | "assign"
	confirmDel  bool
	width       int
	// form fields for create/edit
	formName        string
	formDesc        string
	formCategoryID  string
	formStanding    string
	formIsUniversal bool
}

func NewFactionsScreen(token string) tea.Model {
	return FactionsModel{token: token, loading: true, selected: 0, mode: "list"}
}

func (m FactionsModel) Init() tea.Cmd {
	return func() tea.Msg {
		factions, err := api.ListFactions()
		if err != nil {
			return FactionsErrMsg{Err: err}
		}
		categories, err := api.ListFactionCategories()
		if err != nil {
			return FactionsErrMsg{Err: err}
		}
		return FactionsMsg{Factions: factions, Categories: categories}
	}
}

func (m FactionsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		return m, nil
	case FactionsMsg:
		m.factions = msg.Factions
		m.categories = msg.Categories
		m.loading = false
		return m, nil
	case FactionsErrMsg:
		m.errMsg = fmt.Sprintf("Failed to load factions: %v", msg.Err)
		m.loading = false
		return m, nil
	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func (m FactionsModel) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		return m.back()
	case "enter":
		if m.mode == "list" && len(m.factions) > 0 {
			m.mode = "detail"
		}
		return m, nil
	case "d":
		return m.handleDelete()
	case "c":
		if m.mode == "list" {
			m.mode = "create"
			m.formName = ""
			m.formDesc = ""
			m.formCategoryID = ""
			m.formStanding = "0"
			m.formIsUniversal = false
		}
		return m, nil
	case "a":
		if m.mode == "detail" && len(m.factions) > 0 {
			m.mode = "assign"
			m.selectedChar = 0
		}
		return m, nil
	case "r":
		m.loading = true
		return m, func() tea.Msg {
			factions, err := api.ListFactions()
			if err != nil {
				return FactionsErrMsg{Err: err}
			}
			categories, err := api.ListFactionCategories()
			if err != nil {
				return FactionsErrMsg{Err: err}
			}
			return FactionsMsg{Factions: factions, Categories: categories}
		}
	case "up", "k":
		if m.mode == "list" && m.selected > 0 {
			m.selected--
		}
		if m.mode == "assign" && m.selectedChar > 0 {
			m.selectedChar--
		}
		return m, nil
	case "down", "j":
		if m.mode == "list" && m.selected < len(m.factions)-1 {
			m.selected++
		}
		if m.mode == "assign" && len(m.characters) > 0 && m.selectedChar < len(m.characters)-1 {
			m.selectedChar++
		}
		return m, nil
	}
	if m.mode == "list" {
		switch msg.String() {
		case "u": return m, func() tea.Msg { return NavigateMsg{Screen: 2} }
		case "n": return m, func() tea.Msg { return NavigateMsg{Screen: 4} }
		case "i": return m, func() tea.Msg { return NavigateMsg{Screen: 6} }
		case "q": return m, func() tea.Msg { return NavigateMsg{Screen: 7} }
		case "b": return m, func() tea.Msg { return NavigateMsg{Screen: 8} }
		case "w": return m, func() tea.Msg { return NavigateMsg{Screen: 9} }
		}
	}
	return m, nil
}

func (m FactionsModel) back() (tea.Model, tea.Cmd) {
	switch m.mode {
	case "detail":
		m.mode = "list"
	case "create":
		m.mode = "list"
	case "assign":
		m.mode = "detail"
	default:
		return m, func() tea.Msg { return NavigateMsg{Screen: 1} }
	}
	return m, nil
}

func (m FactionsModel) handleDelete() (tea.Model, tea.Cmd) {
	if m.mode != "detail" || len(m.factions) == 0 {
		return m, nil
	}
	if m.confirmDel {
		go func() { api.DeleteFaction(m.factions[m.selected].ID) }()
		m.factions = filterFactions(m.factions, m.factions[m.selected].ID)
		m.mode = "list"
		m.confirmDel = false
		return m, nil
	}
	m.confirmDel = true
	return m, nil
}

func filterFactions(factions []api.Faction, id int) []api.Faction {
	result := make([]api.Faction, 0, len(factions))
	for _, f := range factions {
		if f.ID != id {
			result = append(result, f)
		}
	}
	return result
}

func (m FactionsModel) View() string {
	if m.loading {
		return style.Info("Loading factions...")
	}
	switch m.mode {
	case "create":
		return m.viewCreate()
	case "detail":
		return m.viewDetail()
	case "assign":
		return m.viewAssign()
	}
	return m.viewList()
}

func (m FactionsModel) viewList() string {
	lines := []string{
		style.StyleHeader.Render(fmt.Sprintf("Factions (%d)", len(m.factions))),
		style.RenderDivider(max(90, m.width-4)),
		fmt.Sprintf("%-4s %-22s %-10s %-8s %-12s",
			style.StyleTableHeader.Render("ID"),
			style.StyleTableHeader.Render("Name"),
			style.StyleTableHeader.Render("Category"),
			style.StyleTableHeader.Render("Standing"),
			style.StyleTableHeader.Render("Universal"),
		),
	}
	for i, f := range m.factions {
		rowStyle := style.StyleTableRow
		if i == m.selected {
			rowStyle = lipgloss.Style{}.Foreground(style.ColorPrimary).Bold(true)
		}
		catName := ""
		for _, c := range m.categories {
			if c.ID == f.CategoryID {
				catName = c.Name
				break
			}
		}
		universal := "no"
		if f.IsUniversal {
			universal = "yes"
		}
		lines = append(lines, fmt.Sprintf("%-4s %-22s %-10s %-8s %-12s",
			rowStyle.Render(fmt.Sprintf("%d", f.ID)),
			rowStyle.Render(trunc(f.Name, 21)),
			rowStyle.Render(trunc(catName, 9)),
			rowStyle.Render(fmt.Sprintf("%d", f.Standing)),
			rowStyle.Render(universal),
		))
	}
	if len(m.factions) == 0 {
		lines = append(lines, style.StyleMuted.Render("  No factions found"))
	}
	lines = append(lines, "")
	lines = append(lines, style.StyleMuted.Render("  [Enter] view   [C] create   [R] refresh   [Esc] back"))
	if m.errMsg != "" {
		lines = append(lines, "", style.Error(m.errMsg))
	}
	return strings.Join(lines, "\n")
}

func (m FactionsModel) viewDetail() string {
	if m.selected >= len(m.factions) {
		return m.viewList()
	}
	f := m.factions[m.selected]
	catName := ""
	for _, c := range m.categories {
		if c.ID == f.CategoryID {
			catName = c.Name
			break
		}
	}
	lines := []string{
		style.StyleHeader.Render(fmt.Sprintf("Faction #%d — %s", f.ID, f.Name)),
		style.RenderDivider(max(90, m.width-4)),
		fmt.Sprintf("  %-14s %s", style.StyleLabel.Render("Name:"), style.StyleValue.Render(f.Name)),
		fmt.Sprintf("  %-14s %s", style.StyleLabel.Render("Description:"), style.StyleValue.Render(trunc(f.Description, 60))),
		fmt.Sprintf("  %-14s %s", style.StyleLabel.Render("Category:"), style.StyleValue.Render(catName)),
		fmt.Sprintf("  %-14s %d", style.StyleLabel.Render("Standing:"), style.StyleValue.Render(fmt.Sprintf("%d", f.Standing))),
		fmt.Sprintf("  %-14s %t", style.StyleLabel.Render("Universal:"), style.StyleValue.Render(fmt.Sprintf("%t", f.IsUniversal))),
		fmt.Sprintf("  %-14s %v", style.StyleLabel.Render("Member IDs:"), style.StyleValue.Render(fmt.Sprintf("%v", f.Members))),
		"",
		style.StyleMuted.Render("  [A] assign character   [D] delete   [Esc] back to list"),
	}
	if m.confirmDel {
		lines = append(lines, "", style.StyleDanger.Render("  Confirm DELETE? Press [D] again to confirm"))
	}
	return strings.Join(lines, "\n")
}

func (m FactionsModel) viewCreate() string {
	lines := []string{
		style.StyleHeader.Render("Create Faction"),
		style.RenderDivider(max(90, m.width-4)),
		style.StyleMuted.Render("  [Tab] switch field   [Enter] submit   [Esc] cancel"),
		"",
		fmt.Sprintf("  %-14s [%s]", style.StyleLabel.Render("Name:"), m.formName),
		fmt.Sprintf("  %-14s [%s]", style.StyleLabel.Render("Description:"), m.formDesc),
		fmt.Sprintf("  %-14s [%s]", style.StyleLabel.Render("Category ID:"), m.formCategoryID),
		fmt.Sprintf("  %-14s [%s]", style.StyleLabel.Render("Standing:"), m.formStanding),
		fmt.Sprintf("  %-14s [%t]", style.StyleLabel.Render("Universal:"), m.formIsUniversal),
	}
	return strings.Join(lines, "\n")
}

func (m FactionsModel) viewAssign() string {
	// Load characters list for assignment
	if len(m.characters) == 0 {
		chars, err := api.ListCharacters()
		if err != nil {
			return style.Error(fmt.Sprintf("Failed to load characters: %v", err))
		}
		m.characters = chars
	}
	f := m.factions[m.selected]
	lines := []string{
		style.StyleHeader.Render(fmt.Sprintf("Assign to Faction #%d — %s", f.ID, f.Name)),
		style.RenderDivider(max(90, m.width-4)),
		style.StyleMuted.Render("  Use [Up/Down] to select, [A] assign, [R] remove, [Esc] back"),
		"",
		fmt.Sprintf("%-4s %-22s %-10s %-6s",
			style.StyleTableHeader.Render("ID"),
			style.StyleTableHeader.Render("Name"),
			style.StyleTableHeader.Render("Class"),
			style.StyleTableHeader.Render("Level"),
		),
	}
	for i, c := range m.characters {
		rowStyle := style.StyleTableRow
		if i == m.selectedChar {
			rowStyle = lipgloss.Style{}.Foreground(style.ColorPrimary).Bold(true)
		}
		lines = append(lines, fmt.Sprintf("%-4s %-22s %-10s %-6s",
			rowStyle.Render(fmt.Sprintf("%d", c.ID)),
			rowStyle.Render(trunc(c.Name, 21)),
			rowStyle.Render(c.Class),
			rowStyle.Render(fmt.Sprintf("%d", c.Level)),
		))
	}
	if len(m.characters) == 0 {
		lines = append(lines, style.StyleMuted.Render("  No characters found"))
	}
	return strings.Join(lines, "\n")
}

type FactionsMsg struct {
	Factions   []api.Faction
	Categories []api.FactionCategory
}
type FactionsErrMsg struct{ Err error }
