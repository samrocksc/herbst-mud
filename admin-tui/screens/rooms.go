package screens

import (
	"github.com/charmbracelet/bubbletea"
	"herbst-mud/admin-tui/style"
)

type RoomsModel struct{}

func NewRoomsScreen(token string) *RoomsModel {
	return &RoomsModel{}
}

func (m *RoomsModel) Init() tea.Cmd  { return nil }
func (m *RoomsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}
func (m *RoomsModel) View() string {
	return style.StyleHeader.Render("Rooms Screen") + "\n" + "(Coming soon)\n"
}