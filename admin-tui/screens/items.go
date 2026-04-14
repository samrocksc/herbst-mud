package screens

import (
	"github.com/charmbracelet/bubbletea"
	"herbst-mud/admin-tui/style"
)

type ItemsModel struct{}

func NewItemsScreen(token string) *ItemsModel {
	return &ItemsModel{}
}

func (m *ItemsModel) Init() tea.Cmd  { return nil }
func (m *ItemsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}
func (m *ItemsModel) View() string {
	return style.StyleHeader.Render("Items Screen") + "\n" + "(Coming soon)\n"
}