package screens

import (
	"fmt"

	"github.com/charmbracelet/bubbletea"
	"herbst-mud/admin-tui/api"
	"herbst-mud/admin-tui/style"
)

// LoginModel is the login screen state.
type LoginModel struct {
	email    string
	password string
	focused  int // 0=email, 1=password
	errMsg   string
	loading  bool
}

// NewLoginScreen creates a fresh login screen.
func NewLoginScreen() tea.Model {
	return LoginModel{focused: 0}
}

func (m LoginModel) Init() tea.Cmd {
	return nil
}

func (m LoginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			m.focused = (m.focused + 1) % 2
			return m, nil
		case "enter":
			if m.focused == 0 {
				m.focused = 1
				return m, nil
			}
			// Attempt login
			m.loading = true
			resp, err := api.Login(m.email, m.password)
			if err != nil {
				m.errMsg = fmt.Sprintf("Login failed: %v", err)
				m.loading = false
				return m, nil
			}
			// Login succeeded — signal main to switch screen
			return m, func() tea.Msg {
				return AuthSuccessMsg{Token: resp.Token, UserID: resp.UserID, Email: resp.Email, IsAdmin: resp.IsAdmin}
			}
		case "ctrl+c", "esc":
			return m, tea.Quit
		}

		// Typing
		if m.focused == 0 {
			switch msg.String() {
			case "backspace":
				if len(m.email) > 0 {
					m.email = m.email[:len(m.email)-1]
				}
			default:
				m.email += msg.String()
			}
		} else {
			switch msg.String() {
			case "backspace":
				if len(m.password) > 0 {
					m.password = m.password[:len(m.password)-1]
				}
			default:
				m.password += msg.String()
			}
		}
	}
	return m, nil
}

func (m LoginModel) View() string {
	emailStyle := style.StyleInput
	passwordStyle := style.StyleInput
	if m.focused == 0 {
		emailStyle = style.StyleInputFocused
	}
	if m.focused == 1 {
		passwordStyle = style.StyleInputFocused
	}

	errStr := ""
	if m.errMsg != "" {
		errStr = "\n" + style.Error(m.errMsg)
	}

	loadingStr := ""
	if m.loading {
		loadingStr = "\n" + style.Info("Logging in...")
	}

	return fmt.Sprintf(`
%s

  %s  %s
  %s  %s

  %s  %s
  %s  %s

%s%s
`,
		style.StyleTitle.Render("🌿 herbst-mud  Admin Login"),
		style.StyleLabel.Render("Email:"),
		emailStyle.Width(40).Render(m.email),
		style.StyleLabel.Render("Password:"),
		passwordStyle.Width(40).Render(mask(m.password)),
		style.StyleMuted.Render("  [Tab] toggle field"),
		style.StyleMuted.Render("[Enter] login"),
		style.StyleMuted.Render("  [Esc] quit"),
		errStr, loadingStr,
	)
}

func mask(s string) string {
	result := ""
	for range s {
		result += "•"
	}
	return result
}

// AuthSuccessMsg signals login success to the root model.
type AuthSuccessMsg struct {
	Token   string
	UserID  int
	Email   string
	IsAdmin bool
}
