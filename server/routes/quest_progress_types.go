package routes

import (
	"herbst-server/db"
)

// questProgressView is the JSON response shape for a QuestProgress entity.
type questProgressView struct {
	ID               int                    `json:"id"`
	CharacterID      int                    `json:"character_id"`
	QuestID          int                    `json:"quest_id"`
	Status           string                 `json:"status"`
	StartedAt        string                 `json:"started_at"`
	CompletedAt      *string                `json:"completed_at,omitempty"`
	CurrentStep      int                    `json:"current_step"`
	ObjectiveCounts  map[string]int         `json:"objective_counts"`
	QuestName        string                 `json:"quest_name,omitempty"`
	RewardsApplied   map[string]interface{} `json:"rewards_applied,omitempty"`
}

// questProgressToView converts a db.QuestProgress entity to a questProgressView.
func questProgressToView(p *db.QuestProgress) questProgressView {
	statusStr := string(p.Status)
	startedAt := p.StartedAt.Format("2006-01-02T15:04:05Z")
	var completedAt *string
	if p.CompletedAt != nil {
		f := p.CompletedAt.Format("2006-01-02T15:04:05Z")
		completedAt = &f
	}
	counts := p.ObjectiveCounts
	if counts == nil {
		counts = map[string]int{}
	}
	view := questProgressView{
		ID:              p.ID,
		Status:          statusStr,
		StartedAt:       startedAt,
		CompletedAt:     completedAt,
		CurrentStep:     p.CurrentStep,
		ObjectiveCounts: counts,
	}
	// Populate IDs from edges if loaded
	if p.Edges.Character != nil {
		view.CharacterID = p.Edges.Character.ID
	}
	if p.Edges.Quest != nil {
		view.QuestID = p.Edges.Quest.ID
		view.QuestName = p.Edges.Quest.Name
	}
	// If edges weren't loaded, we can't populate IDs — caller must use WithCharacter/WithQuest
	return view
}

// questAcceptInput is the JSON request for accepting a quest.
type questAcceptInput struct {
	QuestID int `json:"quest_id" binding:"required"`
}

// questCheckInput is the JSON request body for checking quest progress.
type questCheckInput struct {
	ObjectiveKey string `json:"objective_key" binding:"required"`
	Count        int    `json:"count"`
}