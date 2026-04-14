package screens

import (
	"github.com/charmbracelet/bubbletea"
	"herbst-mud/admin-tui/style"
)

type QuestsModel struct{}

func NewQuestsScreen(token string) *QuestsModel {
	return &QuestsModel{}
}

func (m *QuestsModel) Init() tea.Cmd  { return nil }
func (m *QuestsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}
func (m *QuestsModel) View() string {
	return style.StyleHeader.Render("Quests Screen") + "\n" + "(Coming soon)\n"
}