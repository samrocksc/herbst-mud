package screens

import (
	"github.com/charmbracelet/bubbletea"
	"herbst-mud/admin-tui/style"
)

type BackupModel struct{}

func NewBackupScreen(token string) *BackupModel {
	return &BackupModel{}
}

func (m *BackupModel) Init() tea.Cmd  { return nil }
func (m *BackupModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}
func (m *BackupModel) View() string {
	return style.StyleHeader.Render("Backup Screen") + "\n" + "(Coming soon)\n"
}