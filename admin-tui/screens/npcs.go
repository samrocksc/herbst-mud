package screens

import (
	"github.com/charmbracelet/bubbletea"
	"herbst-mud/admin-tui/style"
)

type NPCsModel struct{}

func NewNPCsScreen(token string) *NPCsModel {
	return &NPCsModel{}
}

func (m *NPCsModel) Init() tea.Cmd  { return nil }
func (m *NPCsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}
func (m *NPCsModel) View() string {
	return style.StyleHeader.Render("NPCs Screen") + "\n" + "(Coming soon)\n"
}