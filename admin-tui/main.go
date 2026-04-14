package main

import (
	"log"
	"os"

	"github.com/charmbracelet/bubbletea"
	"herbst-mud/admin-tui/api"
	"herbst-mud/admin-tui/screens"
)

func main() {
	// Allow BaseURL to be overridden by environment
	if url := os.Getenv("API_BASE_URL"); url != "" {
		api.BaseURL = url
	}

	// Try auto-login with stored token
	if token := api.GetToken(); token != "" {
		// Token exists — start at dashboard
		model := NewRootModel()
		model.token = token
		model.currentScreen = screenDashboard
		model.screenModel = screens.NewDashboardScreen(token, model.currentUser)
		if _, err := tea.NewProgram(model).Run(); err != nil {
			log.Fatal(err)
		}
		return
	}

	// No token — show login screen
	if _, err := tea.NewProgram(NewRootModel()).Run(); err != nil {
		log.Fatal(err)
	}
}