package screens

import (
	"fmt"

	"github.com/charmbracelet/bubbletea"
	"herbst-mud/admin-tui/api"
	"herbst-mud/admin-tui/style"
)

// DashboardModel is the main dashboard / home screen.
type DashboardModel struct {
	token      string
	user       UserInfo
	health     api.HealthResponse
	loading    bool
	errMsg     string
	quitting   bool
	width      int
	height     int
}

// UserInfo mirrors main.UserInfo so screens don't import main.
type UserInfo = struct {
	ID     int
	Email  string
	IsAdmin bool
}

// NewDashboardScreen creates the dashboard screen.
func NewDashboardScreen(token string, user UserInfo) tea.Model {
	m := DashboardModel{token: token, user: user, loading: true}
	return m
}

func (m DashboardModel) Init() tea.Cmd {
	return func() tea.Msg {
		h, err := api.Health()
		if err != nil {
			return HealthErrMsg{Err: err}
		}
		return HealthMsg{Health: h}
	}
}

func (m DashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case HealthMsg:
		m.health = msg.Health
		m.loading = false
		return m, nil
	case HealthErrMsg:
		m.errMsg = fmt.Sprintf("Health check failed: %v", msg.Err)
		m.loading = false
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "u":
			return m, func() tea.Msg { return NavigateMsg{Screen: 2} } // screenUsers
		case "c":
			return m, func() tea.Msg { return NavigateMsg{Screen: 3} } // screenCharacters
		case "r":
			return m, func() tea.Msg { return NavigateMsg{Screen: 4} } // screenRooms
		case "n":
			return m, func() tea.Msg { return NavigateMsg{Screen: 5} } // screenNPCs
		case "i":
			return m, func() tea.Msg { return NavigateMsg{Screen: 6} } // screenItems
		case "q":
			return m, func() tea.Msg { return NavigateMsg{Screen: 7} } // screenQuests
		case "f":
			return m, func() tea.Msg { return NavigateMsg{Screen: 8} } // screenFactions
		case "s":
			return m, func() tea.Msg { return NavigateMsg{Screen: 9} } // screenSkills
		case "b":
			return m, func() tea.Msg { return NavigateMsg{Screen: 10} } // screenBackup
		case "w":
			return m, func() tea.Msg { return NavigateMsg{Screen: 11} } // screenWipe
		}
	}
	return m, nil
}

func (m DashboardModel) View() string {
	if m.loading {
		return style.Info("Loading health status...")
	}

	healthBox := style.StyleBox.Width(max(50, m.width-4)).Render(fmt.Sprintf(
		"%s  %s\n%s  %s\n%s  %s",
		style.StyleLabel.Render("Status:"), healthBadge(m.health.Status),
		style.StyleLabel.Render("SSH:"),   m.health.SSH,
		style.StyleLabel.Render("DB:"),     m.health.DB,
	))

	errStr := ""
	if m.errMsg != "" {
		errStr = "\n" + style.Error(m.errMsg)
	}

	return fmt.Sprintf(`
%s

  Welcome back, %s

  %s

  %s

%s
`,
		style.StyleTitle.Render("🌿 herbst-mud  Dashboard"),
		style.StyleValue.Render(m.user.Email),
		healthBox,
		m.renderMenu(),
		errStr,
	)
}

func (m DashboardModel) renderMenu() string {
	items := []string{
		"[U] Users",
		"[C] Characters",
		"[R] Rooms",
		"[N] NPCs",
		"[I] Items",
		"[Q] Quests",
		"[F] Factions",
		"[S] Skills",
		"[B] Backup",
		"[W] Wipe World",
	}
	result := style.StyleLabel.Render("Navigate:") + "\n"
	for _, item := range items {
		result += "  " + style.StyleValue.Render(item) + "\n"
	}
	return result
}

func healthBadge(s string) string {
	switch s {
	case "ok", "healthy", "running":
		return style.StyleSuccess.Render(s)
	default:
		return style.StyleDanger.Render(s)
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func trunc(s string, max int) string {
	if len(s) > max {
		return s[:max-1] + "…"
	}
	return s
}

func humanSize(bytes int64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	}
	if bytes < 1024*1024 {
		return fmt.Sprintf("%.1f KB", float64(bytes)/1024)
	}
	return fmt.Sprintf("%.1f MB", float64(bytes)/(1024*1024))
}

func mask(s string) string {
	r := ""
	for range s {
		r += "•"
	}
	return r
}

// ─── Shared message types ────────────────────────────────────────────────────

type HealthMsg struct {
	Health api.HealthResponse
}

type HealthErrMsg struct {
	Err error
}

// NavigateMsg signals a screen change.
type NavigateMsg struct {
	Screen int
}
