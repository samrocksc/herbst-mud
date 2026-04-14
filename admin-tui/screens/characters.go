package screens

import (
	"github.com/charmbracelet/bubbletea"
	"herbst-mud/admin-tui/style"
)

type CharactersModel struct{}

func NewCharactersScreen(token string) *CharactersModel {
	return &CharactersModel{}
}

func (m *CharactersModel) Init() tea.Cmd  { return nil }
func (m *CharactersModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}
func (m *CharactersModel) View() string {
	return style.StyleHeader.Render("Characters Screen") + "\n" + "(Coming soon)\n"
}