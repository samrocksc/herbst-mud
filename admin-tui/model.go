package main

import (
	"fmt"

	"github.com/charmbracelet/bubbletea"
	"herbst-mud/admin-tui/screens"
	"herbst-mud/admin-tui/style"
	"github.com/charmbracelet/lipgloss"
)

// screen is the active screen identifier.
type screen int

const (
	screenLogin     screen = iota
	screenDashboard screen = iota
	screenUsers
	screenCharacters
	screenRooms
	screenNPCs
	screenItems
	screenQuests
	screenBackup
	screenWipe
)

var screenNames = map[screen]string{
	screenLogin:      "Login",
	screenDashboard:  "Dashboard",
	screenUsers:      "Users",
	screenCharacters: "Characters",
	screenRooms:      "Rooms",
	screenNPCs:       "NPCs",
	screenItems:      "Items",
	screenQuests:     "Quests",
	screenBackup:     "Backup",
	screenWipe:       "Wipe World",
}

// RootModel holds shared state across all screens.
type RootModel struct {
	currentScreen screen
	token         string
	currentUser   UserInfo
	screenModel   tea.Model
	quitting      bool
	width         int
	height        int
}

// UserInfo is stored after successful login.
type UserInfo struct {
	ID     int
	Email  string
	IsAdmin bool
}

// NewRootModel creates the root model, starting at the login screen.
func NewRootModel() RootModel {
	return RootModel{
		currentScreen: screenLogin,
		screenModel:   screens.NewLoginScreen(),
	}
}

func (m RootModel) Init() tea.Cmd {
	return nil
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		// Global navigation when not on login
		if m.currentScreen != screenLogin {
			switch msg.String() {
			case "esc":
				m.currentScreen = screenDashboard
				m.screenModel = screens.NewDashboardScreen(m.token, m.currentUser)
				return m, nil
			case "ctrl+c":
				return m, tea.Quit
			}
		}

	case screens.AuthSuccessMsg:
		// Login succeeded — switch to dashboard
		m.token = msg.Token
		m.currentUser = UserInfo{ID: msg.UserID, Email: msg.Email, IsAdmin: msg.IsAdmin}
		m.currentScreen = screenDashboard
		m.screenModel = screens.NewDashboardScreen(m.token, m.currentUser)
		return m, nil

	case screens.NavigateMsg:
		m.currentScreen = screen(msg.Screen)
		switch screen(msg.Screen) {
		case screenDashboard:
			m.screenModel = screens.NewDashboardScreen(m.token, m.currentUser)
		case screenUsers:
			m.screenModel = screens.NewUsersScreen(m.token)
		case screenCharacters:
			m.screenModel = screens.NewCharactersScreen(m.token)
		case screenRooms:
			m.screenModel = screens.NewRoomsScreen(m.token)
		case screenNPCs:
			m.screenModel = screens.NewNPCsScreen(m.token)
		case screenItems:
			m.screenModel = screens.NewItemsScreen(m.token)
		case screenQuests:
			m.screenModel = screens.NewQuestsScreen(m.token)
		case screenBackup:
			m.screenModel = screens.NewBackupScreen(m.token)
		case screenWipe:
			m.screenModel = screens.NewWipeScreen(m.token)
		}
		return m, nil
	}

	// Pass updates to the active screen
	updated, cmd := m.screenModel.Update(msg)
	m.screenModel = updated
	return m, cmd
}

func (m RootModel) View() string {
	nav := m.renderNav()
	content := m.screenModel.View()
	return nav + "\n" + content
}

func (m RootModel) renderNav() string {
	if m.currentScreen == screenLogin {
		return ""
	}

	items := []string{
		"[D]ashboard",
		"[U]sers",
		"[C]hars",
		"[R]ooms",
		"[N]PCs",
		"[I]tems",
		"[Q]uests",
		"[B]ackup",
		"[W]ipe",
		"[Esc]Back",
	}

	width := m.width
	if width < 10 {
		width = 80
	}
	div := style.RenderDivider(width)

	navLine := ""
	for _, item := range items {
		itemStyle := lipgloss.Style{}.Foreground(style.ColorMuted)
		if m.currentScreen == screenFromKey(item) {
			itemStyle = lipgloss.Style{}.Foreground(style.ColorPrimary).Bold(true)
		}
		navLine += itemStyle.Render(item) + "  "
	}

	return fmt.Sprintf("%s\n%s  Logged in: %s\n%s",
		div,
		style.StyleHeader.Render("🌿 herbst-mud  admin"),
		style.StyleValue.Render(m.currentUser.Email),
		navLine,
	)
}

func screenFromKey(item string) screen {
	switch item[1] {
	case 'D':
		return screenDashboard
	case 'U':
		return screenUsers
	case 'C':
		return screenCharacters
	case 'R':
		return screenRooms
	case 'N':
		return screenNPCs
	case 'I':
		return screenItems
	case 'Q':
		return screenQuests
	case 'B':
		return screenBackup
	case 'W':
		return screenWipe
	}
	return screenDashboard
}
