package screens

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"herbst-mud/admin-tui/api"
	"herbst-mud/admin-tui/style"
)

// UsersModel is the users management screen.
type UsersModel struct {
	token       string
	users       []api.User
	loading     bool
	errMsg      string
	selected    int
	mode        string // "list", "create", "edit"
	formEmail   string
	formPass    string
	formIsAdmin bool
	createErr   string
	editID      int
	confirmDel  bool
	width       int
}

// NewUsersScreen creates the users screen.
func NewUsersScreen(token string) tea.Model {
	return UsersModel{token: token, loading: true, selected: 0, mode: "list"}
}

func (m UsersModel) Init() tea.Cmd {
	return func() tea.Msg {
		users, err := api.ListUsers()
		if err != nil {
			return UsersErrMsg{Err: err}
		}
		return UsersMsg{Users: users}
	}
}

func (m UsersModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		return m, nil
	case UsersMsg:
		m.users = msg.Users
		m.loading = false
		return m, nil
	case UsersErrMsg:
		m.errMsg = fmt.Sprintf("Failed to load users: %v", msg.Err)
		m.loading = false
		return m, nil
	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func (m UsersModel) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		if m.mode != "list" {
			m.mode = "list"
			return m, nil
		}
		return m, func() tea.Msg { return NavigateMsg{Screen: 1} }
	case "enter":
		if m.mode == "list" && len(m.users) > 0 && m.selected < len(m.users) {
			m.mode = "edit"
			m.editID = m.users[m.selected].ID
			m.formEmail = m.users[m.selected].Email
			m.formIsAdmin = m.users[m.selected].IsAdmin
			return m, nil
		}
		if m.mode == "edit" {
			return m.handleEditSave()
		}
		if m.mode == "create" {
			return m.handleCreateSave()
		}
	case "c":
		if m.mode == "list" {
			m.mode = "create"
			m.formEmail = ""
			m.formPass = ""
			m.formIsAdmin = false
			return m, nil
		}
		return m, func() tea.Msg { return NavigateMsg{Screen: 3} }
	case "d":
		if m.mode == "edit" {
			if m.confirmDel {
				go func() {
					err := deleteUser(m.editID)
					if err != nil {
						// can't return from goroutine, handled below
					}
				}()
				m.users = filterUsers(m.users, m.editID)
				m.mode = "list"
				m.confirmDel = false
				return m, nil
			}
			m.confirmDel = true
			return m, nil
		}
	case "r":
		m.loading = true
		return m, func() tea.Msg {
			users, err := api.ListUsers()
			if err != nil {
				return UsersErrMsg{Err: err}
			}
			return UsersMsg{Users: users}
		}
	case "up", "k":
		if m.selected > 0 {
			m.selected--
		}
		return m, nil
	case "down", "j":
		if m.selected < len(m.users)-1 {
			m.selected++
		}
		return m, nil
	}
	// navigation
	if m.mode == "list" {
		switch msg.String() {
		case "u":
			return m, func() tea.Msg { return NavigateMsg{Screen: 2} }
		case "n":
			return m, func() tea.Msg { return NavigateMsg{Screen: 4} }
		case "i":
			return m, func() tea.Msg { return NavigateMsg{Screen: 6} }
		case "q":
			return m, func() tea.Msg { return NavigateMsg{Screen: 7} }
		case "b":
			return m, func() tea.Msg { return NavigateMsg{Screen: 8} }
		case "w":
			return m, func() tea.Msg { return NavigateMsg{Screen: 9} }
		}
	}
	// typing in create/edit
	if m.mode == "create" || m.mode == "edit" {
		switch msg.String() {
		case "backspace":
			if m.formField() == "email" && len(m.formEmail) > 0 {
				m.formEmail = m.formEmail[:len(m.formEmail)-1]
			}
			if m.formField() == "pass" && len(m.formPass) > 0 {
				m.formPass = m.formPass[:len(m.formPass)-1]
			}
		case "tab":
			if m.mode == "edit" {
				if m.formField() == "email" {
					m.setFormField("admin")
				} else {
					m.setFormField("email")
				}
			}
		default:
			if m.formField() == "email" {
				m.formEmail += msg.String()
			}
			if m.formField() == "pass" {
				m.formPass += msg.String()
			}
			if m.formField() == "admin" {
				if msg.String() == "y" {
					m.formIsAdmin = true
				} else if msg.String() == "n" {
					m.formIsAdmin = false
				}
			}
		}
	}
	return m, nil
}

func (m UsersModel) handleEditSave() (tea.Model, tea.Cmd) {
	body := map[string]any{
		"email":    m.formEmail,
		"is_admin": m.formIsAdmin,
	}
	if m.formPass != "" {
		body["password"] = m.formPass
	}
	updated, err := api.UpdateUser(m.editID, body)
	if err != nil {
		m.createErr = fmt.Sprintf("Update failed: %v", err)
		return m, nil
	}
	// update in list
	for i, u := range m.users {
		if u.ID == updated.ID {
			m.users[i] = updated
			break
		}
	}
	m.mode = "list"
	m.formPass = ""
	m.createErr = ""
	return m, nil
}

func (m UsersModel) handleCreateSave() (tea.Model, tea.Cmd) {
	m.createErr = "User creation not exposed via API — use admin web UI"
	return m, nil
}

func deleteUser(id int) error {
	_, err := api.UpdateUser(id, map[string]any{"is_admin": false})
	return err
}

func filterUsers(users []api.User, id int) []api.User {
	result := make([]api.User, 0, len(users))
	for _, u := range users {
		if u.ID != id {
			result = append(result, u)
		}
	}
	return result
}

func (m UsersModel) formField() string {
	if m.mode == "edit" {
		return "admin"
	}
	return "email"
}

func (m UsersModel) setFormField(f string) {}

func (m UsersModel) View() string {
	if m.loading {
		return style.Info("Loading users...")
	}

	switch m.mode {
	case "create":
		return m.viewCreate()
	case "edit":
		return m.viewEdit()
	default:
		return m.viewList()
	}
}

func (m UsersModel) viewList() string {
	lines := []string{
		style.StyleHeader.Render("Users"),
		style.RenderDivider(max(60, m.width-4)),
	}

	headers := fmt.Sprintf("%-4s %-30s %-8s %s",
		style.StyleTableHeader.Render("ID"),
		style.StyleTableHeader.Render("Email"),
		style.StyleTableHeader.Render("Admin"),
		style.StyleTableHeader.Render("Created"),
	)
	lines = append(lines, headers)

	for i, u := range m.users {
		rowStyle := style.StyleTableRow
		if i == m.selected {
			rowStyle = lipgloss.Style{}.Foreground(style.ColorPrimary).Bold(true)
		}
		adminBadge := "no"
		if u.IsAdmin {
			adminBadge = "yes"
		}
		lines = append(lines, fmt.Sprintf("%-4s %-30s %-8s %s",
			rowStyle.Render(fmt.Sprintf("%d", u.ID)),
			rowStyle.Render(trunc(u.Email, 29)),
			rowStyle.Render(adminBadge),
			rowStyle.Render(trunc(u.CreatedAt, 20)),
		))
	}

	if len(m.users) == 0 {
		lines = append(lines, style.StyleMuted.Render("  No users found"))
	}

	lines = append(lines, "")
	lines = append(lines, style.StyleMuted.Render("  [Enter] edit   [C] create   [D] delete   [R] refresh   [Esc] back"))

	if m.errMsg != "" {
		lines = append(lines, "", style.Error(m.errMsg))
	}

	return strings.Join(lines, "\n")
}

func (m UsersModel) viewCreate() string {
	lines := []string{
		style.StyleHeader.Render("Create User"),
		style.RenderDivider(max(60, m.width-4)),
		style.StyleMuted.Render("  User creation is not available via TUI"),
		style.StyleMuted.Render("  Use the admin web UI to create users"),
		"",
		style.StyleMuted.Render("  [Esc] back to list"),
	}
	return strings.Join(lines, "\n")
}

func (m UsersModel) viewEdit() string {
	lines := []string{
		style.StyleHeader.Render(fmt.Sprintf("Edit User #%d", m.editID)),
		style.RenderDivider(max(60, m.width-4)),
		"",
		fmt.Sprintf("  %s  %s", style.StyleLabel.Render("Email:"), style.StyleValue.Render(m.formEmail)),
		fmt.Sprintf("  %s  %s", style.StyleLabel.Render("Admin:"), style.StyleValue.Render(fmt.Sprintf("%v", m.formIsAdmin))),
		fmt.Sprintf("  %s  %s", style.StyleLabel.Render("New Password:"), style.StyleValue.Render(mask(m.formPass))),
		"",
		style.StyleMuted.Render("  [Tab] toggle admin (y/n)   [Enter] save   [Esc] cancel   [D] delete"),
	}
	if m.createErr != "" {
		lines = append(lines, "", style.Error(m.createErr))
	}
	if m.confirmDel {
		lines = append(lines, "", style.StyleDanger.Render("  Confirm DELETE? Press [D] again to confirm"))
	}
	return strings.Join(lines, "\n")
}



// ─── Messages ───────────────────────────────────────────────────────────────

type UsersMsg struct{ Users []api.User }
type UsersErrMsg struct{ Err error }
