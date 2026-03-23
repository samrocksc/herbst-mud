package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ============================================================
// MODEL LIFECYCLE
// ============================================================

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		inputHeight := m.height * 20 / 100
		if inputHeight < 3 {
			inputHeight = 3
		}
		statusHeight := m.height * 10 / 100
		if statusHeight < 3 {
			statusHeight = 3
		}
		vpHeight := m.height - inputHeight - statusHeight
		if vpHeight < 5 {
			vpHeight = 5
		}
		m.viewport = viewport.New(msg.Width, vpHeight)
		if m.debugMode {
			log.Printf("DEBUG: Window size changed: %dx%d", m.width, m.height)
		}

	case spinner.TickMsg:
		if m.isLoading {
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}

	case tea.KeyMsg:
		key := msg.String()

		if key == "ctrl+c" || key == "ctrl+q" {
			return m, tea.Quit
		}

		if m.isLoading {
			return m, nil
		}

		// Viewport scrolling
		if m.screen == ScreenPlaying {
			vp, vpCmd := m.viewport.Update(msg)
			m.viewport = vp
			if cmd == nil {
				cmd = vpCmd
			}
		}

		// Message history scrolling (ctrl+p up, ctrl+n down)
		if m.screen == ScreenPlaying {
			switch key {
			case "ctrl+p":
				if !m.isScrolling {
					m.isScrolling = true
					m.historyOffset = 1
				} else {
					m.historyOffset++
				}
				maxOffset := len(m.messageHistory) - 1
				if m.historyOffset > maxOffset {
					m.historyOffset = maxOffset
				}
				return m, nil
			case "ctrl+n":
				if !m.isScrolling {
					return m, nil
				}
				m.historyOffset--
				if m.historyOffset < 0 {
					m.historyOffset = 0
					m.isScrolling = false
				}
				return m, nil
			}
		}

		// Vim-style menu navigation
		if m.screen == ScreenWelcome || m.screen == ScreenProfile {
			if key == "j" || key == "down" {
				m.menuCursor++
				if m.menuCursor >= len(m.menuItems) {
					m.menuCursor = 0
				}
				return m, nil
			}
			if key == "k" || key == "up" {
				m.menuCursor--
				if m.menuCursor < 0 {
					m.menuCursor = len(m.menuItems) - 1
				}
				return m, nil
			}
		}

		// Command history navigation (up/down arrows) - only in playing screen
		if m.screen == ScreenPlaying {
			if key == "up" {
				if len(m.commandHistory) > 0 {
					if m.historyIndex < len(m.commandHistory)-1 {
						m.historyIndex++
						m.textInput.SetValue(m.commandHistory[len(m.commandHistory)-1-m.historyIndex])
					}
				}
				return m, nil
			}
			if key == "down" {
				if m.historyIndex > 0 {
					m.historyIndex--
					if m.historyIndex == 0 {
						m.textInput.SetValue("")
					} else {
						m.textInput.SetValue(m.commandHistory[len(m.commandHistory)-1-m.historyIndex])
					}
				}
				return m, nil
			}
		}

		// Enter — process input
		if key == "enter" || key == "ctrl+m" {
			input := m.textInput.Value()
			m.textInput.SetValue("")
			// Save to command history (if not empty and not duplicate of last)
			if input != "" {
				if len(m.commandHistory) == 0 || m.commandHistory[len(m.commandHistory)-1] != input {
					m.commandHistory = append(m.commandHistory, input)
					// Keep history to last 100 commands
					if len(m.commandHistory) > 100 {
						m.commandHistory = m.commandHistory[1:]
					}
				}
			}
			m.historyIndex = 0 // Reset history position
			m.processInput(input)
			return m, nil
		}

		// Escape
		if key == "esc" {
			m.handleEscape()
			return m, nil
		}

		// Text input
		m.textInput, cmd = m.textInput.Update(msg)
		m.inputBuffer = m.textInput.Value()

		return m, cmd
	}

	return m, nil
}

func (m *model) handleEscape() {
	switch m.screen {
	case ScreenLogin, ScreenRegister:
		m.screen = ScreenWelcome
		m.textInput.SetValue("")
		m.inputBuffer = ""
		m.loginUsername = ""
		m.loginPassword = ""
		m.inputField = "username"
		m.AppendMessage("", "")
		m.menuItems = []string{"Login", "Register", "Quit"}
		m.menuCursor = 0
	case ScreenProfile, ScreenEditField:
		m.screen = ScreenPlaying
		m.textInput.SetValue("")
		m.inputBuffer = ""
		m.AppendMessage("", "")
	case ScreenPlaying:
		m.AppendMessage("Type 'quit' or press Ctrl+C to exit", "info")
	}
}

// ============================================================
// INPUT PROCESSING
// ============================================================

func (m *model) processInput(input string) {
	input = strings.TrimSpace(input)

	switch m.screen {
	case ScreenWelcome:
		m.handleWelcomeInput(input)
	case ScreenLogin:
		m.handleLoginInput(input)
	case ScreenRegister:
		m.handleRegisterInput(input)
	case ScreenProfile:
		m.handleProfileInput(input)
	case ScreenEditField:
		m.handleEditFieldInput(input)
	case ScreenPlaying:
		m.processCommand(input)
	}
}

// ============================================================
// VIEW RENDERING
// ============================================================

func (m *model) View() string {
	var s strings.Builder

	if m.isLoading {
		s.WriteString(m.spinner.View())
		s.WriteString(" " + m.loadingMessage)
		m.message = ""
		m.messageType = ""
		return s.String()
	}

	switch m.screen {
	case ScreenWelcome:
		var inputContent strings.Builder
		inputContent.WriteString(promptStyle.Render("> "))
		inputContent.WriteString(m.textInput.View())
		inputContent.WriteString("\n\n")
		inputContent.WriteString(lipgloss.NewStyle().Foreground(gray).Render("Press 1/2/3 or type login/register/quit"))
		return welcomeScreen(m.width, m.height, inputContent.String())

	case ScreenLogin:
		promptText := "> "
		if m.inputField == "username" {
			promptText = "Username: "
		} else if m.inputField == "password" {
			promptText = "Password: "
		}
		inputContent := promptStyle.Render(promptText) + m.textInput.View()
		return loginScreen(m.width, m.height, m.message, m.messageType, inputContent)

	case ScreenRegister:
		promptText := "> "
		if m.inputField == "username" {
			promptText = "Username: "
		} else if m.inputField == "password" {
			promptText = "Password: "
		} else if m.inputField == "confirm_password" {
			promptText = "Confirm: "
		} else if m.inputField == "email" {
			promptText = "Email: "
		}
		inputContent := promptStyle.Render(promptText) + m.textInput.View()
		return registerScreen(m.width, m.height, m.message, m.messageType, inputContent)

	case ScreenProfile:
		s.WriteString("=== CHARACTER PROFILE ===\n\n")
		s.WriteString(fmt.Sprintf("Name: %s\n", lipgloss.NewStyle().Bold(true).Render(m.currentCharacterName)))
		s.WriteString(fmt.Sprintf("Gender: %s\n", m.characterGender))
		s.WriteString(fmt.Sprintf("Description: %s\n\n", m.characterDescription))
		s.WriteString("Stats:\n")
		s.WriteString(StatusBar(m.characterHP, m.characterMaxHP, m.characterStamina, m.characterMaxStamina, m.characterMana, m.characterMaxMana))
		s.WriteString("\n\n")

		for i, item := range m.menuItems {
			cursor := "  "
			if i == m.menuCursor {
				cursor = "▶ "
				s.WriteString(menuSelectedStyle.Render(cursor + item))
			} else {
				s.WriteString(menuNormalStyle.Render(cursor + item))
			}
			s.WriteString("\n")
		}

		s.WriteString("\n")
		if m.message != "" {
			s.WriteString(m.styledMessage(m.message))
			s.WriteString("\n")
		}
		s.WriteString(promptStyle.Render("> "))
		s.WriteString(m.textInput.View())

	case ScreenEditField:
		if m.editField == "gender" {
			s.WriteString("Enter your gender (e.g., he/him, she/her, they/them):\n\n")
		} else {
			s.WriteString("Enter your description (what people see when they look at you):\n\n")
		}
		if m.message != "" {
			s.WriteString(m.styledMessage(m.message))
			s.WriteString("\n\n")
		}
		s.WriteString(promptStyle.Render("> "))
		s.WriteString(m.textInput.View())

	case ScreenPlaying:
		width := m.width
		height := m.height
		if width < 40 {
			width = 80
		}
		if height < 10 {
			height = 24
		}

		inputHeight := height * 20 / 100
		if inputHeight < 3 {
			inputHeight = 3
		}
		statusHeight := height * 10 / 100
		if statusHeight < 3 {
			statusHeight = 3
		}

		roomInfo := fmt.Sprintf("[%s]\n%s\n\nExits: %s",
			lipgloss.NewStyle().Bold(true).Foreground(green).Render(m.roomName),
			m.roomDesc,
			m.formatExitsWithColor())

		if len(m.messageHistory) > 0 {
			msgContent := m.buildOutputContent()
			if msgContent != "" {
				roomInfo += "\n\n" + msgContent
			}
		}

		if m.viewport.Width != width {
			m.viewport.Width = width
		}
		if m.viewport.Height != height-statusHeight-inputHeight {
			m.viewport.Height = height - statusHeight - inputHeight
		}
		m.viewport.SetContent(roomInfo)

		viewportStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(pink).
			Width(width)
		s.WriteString(viewportStyle.Render(m.viewport.View()))
		s.WriteString("\n")

		statsLine := MiniStatusBar(m.characterHP, m.characterMaxHP, m.characterStamina, m.characterMaxStamina, m.characterMana, m.characterMaxMana)
		debugInfo := ""
		if m.debugMode {
			debugInfo = " " + lipgloss.NewStyle().Foreground(yellow).Bold(true).Render(fmt.Sprintf("[Room: %d]", m.currentRoom))
		}
		statusBarStyle := lipgloss.NewStyle().
			Foreground(pink).
			Background(lipgloss.Color("235")).
			Bold(true).
			Width(width).
			Padding(0, 1)
		s.WriteString(statusBarStyle.Render(statsLine + debugInfo))
		s.WriteString("\n")

		inputStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(pink).
			Padding(0, 1).
			Width(width).
			Height(inputHeight - 2)
		s.WriteString(inputStyle.Render(promptStyle.Render("> ") + m.textInput.View()))

		return s.String()
	}

	// Center in terminal
	if m.width > 0 && m.height > 0 && m.width > 60 {
		lines := strings.Split(s.String(), "\n")
		var centered []string
		for _, line := range lines {
			visualWidth := lipgloss.Width(line)
			padding := (m.width - visualWidth) / 2
			if padding > 0 && visualWidth < m.width-10 {
				centered = append(centered, fmt.Sprintf("%*s%s", padding, "", line))
			} else {
				centered = append(centered, line)
			}
		}
		return strings.Join(centered, "\n")
	}

	return s.String()
}
