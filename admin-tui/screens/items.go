package screens

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"herbst-mud/admin-tui/api"
	"herbst-mud/admin-tui/style"
)

// ItemsModel is the items/equipment management screen.
type ItemsModel struct {
	token      string
	items      []api.EquipmentItem
	loading    bool
	errMsg     string
	selected   int
	mode       string // "list", "detail", "create", "edit"
	item       *api.EquipmentItem
	form       itemForm
	formField  int
	createErr  string
	confirmDel bool
	width      int
}

type itemForm struct {
	Name        string
	Description string
	Slot        string
	ItemType    string
	Level       string
	Weight      string
	Color       string
	Healing     string
	Effect      string
}

var itemFormFields = []string{
	"name", "description", "slot", "itemType", "level",
	"weight", "color", "healing", "effect",
}

func (f *itemForm) reset() {
	*f = itemForm{}
}

func (f itemForm) slot() string { return f.Slot }
func (f itemForm) itemType() string { return f.ItemType }

func (f itemForm) fieldValue(field string) string {
	switch field {
	case "name":
		return f.Name
	case "description":
		return f.Description
	case "slot":
		return f.Slot
	case "itemType":
		return f.ItemType
	case "level":
		return f.Level
	case "weight":
		return f.Weight
	case "color":
		return f.Color
	case "healing":
		return f.Healing
	case "effect":
		return f.Effect
	}
	return ""
}

func (f *itemForm) setFieldValue(field, val string) {
	switch field {
	case "name":
		f.Name = val
	case "description":
		f.Description = val
	case "slot":
		f.Slot = val
	case "itemType":
		f.ItemType = val
	case "level":
		f.Level = val
	case "weight":
		f.Weight = val
	case "color":
		f.Color = val
	case "healing":
		f.Healing = val
	case "effect":
		f.Effect = val
	}
}

func (f itemForm) toMap() map[string]any {
	body := map[string]any{
		"name":        f.Name,
		"description": f.Description,
		"slot":        f.Slot,
		"itemType":    f.ItemType,
		"color":       f.Color,
		"effect":      f.Effect,
		"level":       f.Level,
		"weight":      f.Weight,
		"healing":     f.Healing,
	}
	return body
}

// NewItemsScreen creates the items screen.
func NewItemsScreen(token string) tea.Model {
	return ItemsModel{token: token, loading: true, selected: 0, mode: "list"}
}

func (m ItemsModel) Init() tea.Cmd {
	return func() tea.Msg {
		items, err := api.ListItems()
		if err != nil {
			return ItemsErrMsg{Err: err}
		}
		return ItemsMsg{Items: items}
	}
}

func (m ItemsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		return m, nil
	case ItemsMsg:
		m.items = msg.Items
		m.loading = false
		return m, nil
	case ItemsErrMsg:
		m.errMsg = fmt.Sprintf("Failed to load items: %v", msg.Err)
		m.loading = false
		return m, nil
	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func (m ItemsModel) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Handle form input when in edit/create mode
	if m.mode == "edit" || m.mode == "create" {
		return m.handleFormKey(msg)
	}

	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		if m.mode != "list" {
			m.mode = "list"
			m.form.reset()
			return m, nil
		}
		return m, func() tea.Msg { return NavigateMsg{Screen: 1} }
	case "enter":
		if m.mode == "list" && len(m.items) > 0 && m.selected < len(m.items) {
			m.mode = "detail"
			m.item = &m.items[m.selected]
			return m, nil
		}
	case "e":
		if m.mode == "detail" && m.item != nil {
			m.mode = "edit"
			m.form = itemForm{
				Name:        m.item.Name,
				Description: m.item.Description,
				Slot:        m.item.Slot,
				ItemType:    m.item.ItemType,
				Level:       fmt.Sprintf("%d", m.item.Level),
				Weight:      fmt.Sprintf("%d", m.item.Weight),
				Color:       m.item.Color,
				Healing:     fmt.Sprintf("%d", m.item.Healing),
				Effect:      m.item.Effect,
			}
			m.formField = 0
			m.createErr = ""
			return m, nil
		}
	case "c":
		if m.mode == "list" {
			m.mode = "create"
			m.form.reset()
			m.formField = 0
			m.createErr = ""
			return m, nil
		}
	case "d":
		if m.mode == "detail" {
			if m.confirmDel {
				go func() {
					api.DeleteItem(m.item.ID)
				}()
				m.items = filterItems(m.items, m.item.ID)
				m.mode = "list"
				m.confirmDel = false
				return m, nil
			}
			m.confirmDel = true
			return m, nil
		}
	case "r":
		m.loading = true
		m.confirmDel = false
		return m, func() tea.Msg {
			items, err := api.ListItems()
			if err != nil {
				return ItemsErrMsg{Err: err}
			}
			return ItemsMsg{Items: items}
		}
	case "up", "k":
		if m.selected > 0 {
			m.selected--
		}
		return m, nil
	case "down", "j":
		if m.selected < len(m.items)-1 {
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
		case "n":
			return m, func() tea.Msg { return NavigateMsg{Screen: 5} }
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

func (m ItemsModel) handleFormKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	field := itemFormFields[m.formField]

	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		if m.mode == "edit" {
			m.mode = "detail"
		} else {
			m.mode = "list"
		}
		m.form.reset()
		m.createErr = ""
		return m, nil
	case "tab":
		m.formField = (m.formField + 1) % len(itemFormFields)
		return m, nil
	case "shift+tab":
		m.formField = (m.formField - 1 + len(itemFormFields)) % len(itemFormFields)
		return m, nil
	case "enter":
		return m.handleSave()
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

func (m ItemsModel) handleSave() (tea.Model, tea.Cmd) {
	if m.form.Name == "" {
		m.createErr = "Name is required"
		return m, nil
	}

	if m.mode == "edit" && m.item != nil {
		updated, err := api.UpdateItem(m.item.ID, m.form.toMap())
		if err != nil {
			m.createErr = fmt.Sprintf("Save failed: %v", err)
			return m, nil
		}
		for i, it := range m.items {
			if it.ID == updated.ID {
				m.items[i] = updated
				break
			}
		}
		m.item = &updated
	}
	if m.mode == "create" {
		created, err := api.CreateItem(m.form.toMap())
		if err != nil {
			m.createErr = fmt.Sprintf("Create failed: %v", err)
			return m, nil
		}
		m.items = append(m.items, created)
		m.selected = len(m.items) - 1
		m.item = &created
	}
	m.mode = "detail"
	m.form.reset()
	m.createErr = ""
	return m, nil
}

func filterItems(items []api.EquipmentItem, id int) []api.EquipmentItem {
	result := make([]api.EquipmentItem, 0, len(items))
	for _, it := range items {
		if it.ID != id {
			result = append(result, it)
		}
	}
	return result
}

func (m ItemsModel) View() string {
	if m.loading {
		return style.Info("Loading items...")
	}
	switch m.mode {
	case "detail":
		return m.viewDetail()
	case "edit", "create":
		return m.viewEdit()
	default:
		return m.viewList()
	}
}

func (m ItemsModel) viewList() string {
	lines := []string{
		style.StyleHeader.Render(fmt.Sprintf("Items (%d)", len(m.items))),
		style.RenderDivider(max(80, m.width-4)),
		fmt.Sprintf("%-4s %-22s %-10s %-10s %-4s %-6s %-8s",
			style.StyleTableHeader.Render("ID"),
			style.StyleTableHeader.Render("Name"),
			style.StyleTableHeader.Render("Slot"),
			style.StyleTableHeader.Render("Type"),
			style.StyleTableHeader.Render("Lvl"),
			style.StyleTableHeader.Render("Weight"),
			style.StyleTableHeader.Render("Room"),
		),
	}
	for i, it := range m.items {
		rowStyle := style.StyleTableRow
		if i == m.selected {
			rowStyle = lipgloss.Style{}.Foreground(style.ColorPrimary).Bold(true)
		}
		owner := ""
		if it.OwnerID != nil {
			owner = fmt.Sprintf("P:%d", *it.OwnerID)
		} else if it.RoomID > 0 {
			owner = fmt.Sprintf("R:%d", it.RoomID)
		}
		lines = append(lines, fmt.Sprintf("%-4s %-22s %-10s %-10s %-4s %-6s %-8s",
			rowStyle.Render(fmt.Sprintf("%d", it.ID)),
			rowStyle.Render(trunc(it.Name, 21)),
			rowStyle.Render(it.Slot),
			rowStyle.Render(it.ItemType),
			rowStyle.Render(fmt.Sprintf("%d", it.Level)),
			rowStyle.Render(fmt.Sprintf("%d", it.Weight)),
			rowStyle.Render(owner),
		))
	}
	if len(m.items) == 0 {
		lines = append(lines, style.StyleMuted.Render("  No items found"))
	}
	lines = append(lines, "")
	lines = append(lines, style.StyleMuted.Render("  [Enter] view   [C] create   [D] delete   [R] refresh   [Esc] back"))
	if m.errMsg != "" {
		lines = append(lines, "", style.Error(m.errMsg))
	}
	return strings.Join(lines, "\n")
}

func (m ItemsModel) viewDetail() string {
	if m.item == nil {
		return m.viewList()
	}
	lines := []string{
		style.StyleHeader.Render(fmt.Sprintf("Item #%d", m.item.ID)),
		style.RenderDivider(max(80, m.width-4)),
		fmt.Sprintf("  %-12s %s", style.StyleLabel.Render("Name:"), style.StyleValue.Render(m.item.Name)),
		fmt.Sprintf("  %-12s %s", style.StyleLabel.Render("Slot:"), style.StyleValue.Render(m.item.Slot)),
		fmt.Sprintf("  %-12s %s", style.StyleLabel.Render("Type:"), style.StyleValue.Render(m.item.ItemType)),
		fmt.Sprintf("  %-12s %s", style.StyleLabel.Render("Level:"), style.StyleValue.Render(fmt.Sprintf("%d", m.item.Level))),
		fmt.Sprintf("  %-12s %s", style.StyleLabel.Render("Weight:"), style.StyleValue.Render(fmt.Sprintf("%d", m.item.Weight))),
		fmt.Sprintf("  %-12s %s", style.StyleLabel.Render("Color:"), style.StyleValue.Render(m.item.Color)),
		fmt.Sprintf("  %-12s %s", style.StyleLabel.Render("Healing:"), style.StyleValue.Render(fmt.Sprintf("%d", m.item.Healing))),
		fmt.Sprintf("  %-12s %s", style.StyleLabel.Render("Effect:"), style.StyleValue.Render(m.item.Effect)),
		fmt.Sprintf("  %-12s %s", style.StyleLabel.Render("Equipped:"), style.StyleValue.Render(fmt.Sprintf("%v", m.item.IsEquipped))),
		fmt.Sprintf("  %-12s %s", style.StyleLabel.Render("Visible:"), style.StyleValue.Render(fmt.Sprintf("%v", m.item.IsVisible))),
		fmt.Sprintf("  %-12s %s", style.StyleLabel.Render("Room/Owner:"), style.StyleValue.Render(roomOwner(m.item))),
		"",
		style.StyleMuted.Render(fmt.Sprintf("  Description: %s", m.item.Description)),
		"",
		style.StyleMuted.Render("  [E] edit   [C] create   [D] delete   [Esc] back to list"),
	}
	if m.confirmDel {
		lines = append(lines, "", style.StyleDanger.Render("  Confirm DELETE? Press [D] again to confirm"))
	}
	return strings.Join(lines, "\n")
}

func (m ItemsModel) viewEdit() string {
	modeLabel := "Edit Item"
	if m.mode == "create" {
		modeLabel = "Create Item"
	}

	fieldLabels := []struct {
		label string
		field string
	}{
		{"Name:", "name"},
		{"Description:", "description"},
		{"Slot:", "slot"},
		{"Type:", "itemType"},
		{"Level:", "level"},
		{"Weight:", "weight"},
		{"Color:", "color"},
		{"Healing:", "healing"},
		{"Effect:", "effect"},
	}

	lines := []string{
		style.StyleHeader.Render(modeLabel),
		style.RenderDivider(max(80, m.width-4)),
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

func roomOwner(it *api.EquipmentItem) string {
	if it.OwnerID != nil {
		return fmt.Sprintf("Player %d", *it.OwnerID)
	}
	if it.RoomID > 0 {
		return fmt.Sprintf("Room %d", it.RoomID)
	}
	return "—"
}

// ─── Messages ───────────────────────────────────────────────────────────────

type ItemsMsg struct{ Items []api.EquipmentItem }
type ItemsErrMsg struct{ Err error }