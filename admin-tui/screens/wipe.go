package screens

import (
	"github.com/charmbracelet/bubbletea"
	"herbst-mud/admin-tui/style"
)

type WipeModel struct{}

func NewWipeScreen(token string) *WipeModel {
	return &WipeModel{}
}

func (m *WipeModel) Init() tea.Cmd  { return nil }
func (m *WipeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}
func (m *WipeModel) View() string {
	return style.StyleHeader.Render("Wipe World Screen") + "\n" + "(Coming soon)\n"
}