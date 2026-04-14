package screens

import (
	"github.com/charmbracelet/bubbletea"
	"herbst-mud/admin-tui/style"
)

type UsersModel struct{}

func NewUsersScreen(token string) *UsersModel {
	return &UsersModel{}
}

func (m *UsersModel) Init() tea.Cmd  { return nil }
func (m *UsersModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}
func (m *UsersModel) View() string {
	return style.StyleHeader.Render("Users Screen") + "\n" + "(Coming soon)\n"
}