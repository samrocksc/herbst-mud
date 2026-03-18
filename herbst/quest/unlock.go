package quest

import (
	"sync"
)

// CharacterQuestStore tracks character quest unlocks
// In production, this would be backed by a database
type CharacterQuestStore struct {
	mu     sync.RWMutex
	quests map[int]map[string]*CharacterQuest // characterID -> questID -> CharacterQuest
}

// Global store for character quests
var GlobalQuestStore = NewCharacterQuestStore()

// NewCharacterQuestStore creates a new quest store
func NewCharacterQuestStore() *CharacterQuestStore {
	return &CharacterQuestStore{
		quests: make(map[int]map[string]*CharacterQuest),
	}
}

// GetCharacterQuests returns all quests for a character
func (s *CharacterQuestStore) GetCharacterQuests(characterID int) []*CharacterQuest {
	s.mu.RLock()
	defer s.mu.RUnlock()

	quests := s.quests[characterID]
	if quests == nil {
		return nil
	}

	result := make([]*CharacterQuest, 0, len(quests))
	for _, q := range quests {
		result = append(result, q)
	}
	return result
}

// GetCharacterQuest returns a specific quest for a character
func (s *CharacterQuestStore) GetCharacterQuest(characterID int, questID string) *CharacterQuest {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if quests := s.quests[characterID]; quests != nil {
		return quests[questID]
	}
	return nil
}

// HasQuestUnlocked checks if a character has unlocked a quest
func (s *CharacterQuestStore) HasQuestUnlocked(characterID int, questID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if quests := s.quests[characterID]; quests != nil {
		q, exists := quests[questID]
		return exists && q.Status != ""
	}
	return false
}

// UnlockQuest unlocks a quest for a character
func (s *CharacterQuestStore) UnlockQuest(characterID int, questID string) *CharacterQuest {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.quests[characterID] == nil {
		s.quests[characterID] = make(map[string]*CharacterQuest)
	}

	// Check if already unlocked
	if existing := s.quests[characterID][questID]; existing != nil {
		return existing
	}

	quest := &CharacterQuest{
		CharacterID: characterID,
		QuestID:     questID,
		Status:       "unlocked",
		UnlockedAt:  0, // Would use time.Now().Unix() in production
	}

	s.quests[characterID][questID] = quest
	return quest
}

// CheckExamineQuestUnlock checks if examining an item should unlock a quest
// Returns the unlock result if a quest was unlocked or already unlocked
func CheckExamineQuestUnlock(characterID int, targetName string, examineLevel int) *QuestUnlockResult {
	// Find quests triggered by examining this target
	quests := GetQuestByTarget(targetName)
	if len(quests) == 0 {
		return nil
	}

	for _, quest := range quests {
		if quest.ExamineTrigger == nil {
			continue
		}

		// Check if examine level meets requirement
		if examineLevel < quest.ExamineTrigger.MinExamineLevel {
			continue
		}

		// Check if already unlocked
		if GlobalQuestStore.HasQuestUnlocked(characterID, quest.ID) {
			return &QuestUnlockResult{
				Unlocked:        false,
				QuestID:         quest.ID,
				QuestName:       quest.Name,
				AlreadyUnlocked: true,
			}
		}

		// Unlock the quest
		GlobalQuestStore.UnlockQuest(characterID, quest.ID)

		return &QuestUnlockResult{
			Unlocked:    true,
			QuestID:     quest.ID,
			QuestName:   quest.Name,
			RevealText:  quest.ExamineTrigger.RevealText,
			XPGained:    quest.ExamineTrigger.XPReward,
		}
	}

	return nil
}

// GetUnlockedQuests returns all unlocked quests for a character
func GetUnlockedQuests(characterID int) []*Quest {
	characterQuests := GlobalQuestStore.GetCharacterQuests(characterID)
	if characterQuests == nil {
		return nil
	}

	quests := make([]*Quest, 0, len(characterQuests))
	for _, cq := range characterQuests {
		if q := GetQuestByID(cq.QuestID); q != nil {
			quests = append(quests, q)
		}
	}
	return quests
}

// GetVisibleQuests returns quests visible to a character (unlocked + non-hidden)
func GetVisibleQuests(characterID int) []*Quest {
	var visible []*Quest

	// Add all unlocked quests
	unlocked := GetUnlockedQuests(characterID)
	visible = append(visible, unlocked...)

	// Add non-hidden quests (public quests)
	for _, q := range QuestRegistry {
		if !q.Hidden {
			// Check if not already in list
			found := false
			for _, v := range visible {
				if v.ID == q.ID {
					found = true
					break
				}
			}
			if !found {
				visible = append(visible, q)
			}
		}
	}

	return visible
}