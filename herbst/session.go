package main

// ============================================================
// SESSION MANAGEMENT - Screen transitions and state management
// ============================================================

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
		// Re-initialize menu items for welcome screen
		m.menuItems = []string{"Login", "Register", "Quit"}
		m.menuCursor = 0
	case ScreenProfile, ScreenEditField:
		m.screen = ScreenPlaying
		m.textInput.SetValue("")
		m.inputBuffer = ""
		m.AppendMessage("", "")
	case ScreenPlaying:
		// Could add a "really quit?" confirmation
		m.AppendMessage("Type 'quit' or press Ctrl+C to exit", "info")
	}
}
