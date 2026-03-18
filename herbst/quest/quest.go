package quest

// QuestType defines the type of quest trigger
type QuestType string

const (
	// QuestTypeSecret is a hidden quest that unlocks through specific actions
	QuestTypeSecret QuestType = "secret"
	// QuestTypeExamine is a quest triggered by examining items
	QuestTypeExamine QuestType = "examine"
)

// Quest represents a quest definition
type Quest struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Type        QuestType  `json:"type"`
	Description string     `json:"description"`
	Hidden      bool       `json:"hidden"` // Whether quest is hidden until unlocked
	Rewards     QuestRewards `json:"rewards"`

	// For examine-triggered quests
	ExamineTrigger *ExamineTrigger `json:"examine_trigger,omitempty"`
}

// ExamineTrigger defines conditions for examine-triggered quest unlock
type ExamineTrigger struct {
	Target         string `json:"target"`           // Item/NPC name to examine
	MinExamineLevel int   `json:"min_examine_level"` // Required examine skill level
	RevealText     string `json:"reveal_text"`       // Text shown when quest unlocks
	XPReward       int    `json:"xp_reward"`        // XP given to examine skill
}

// QuestRewards defines rewards for completing a quest
type QuestRewards struct {
	XP    int      `json:"xp"`
	Items []string `json:"items,omitempty"`
}

// CharacterQuest tracks a character's quest progress
type CharacterQuest struct {
	CharacterID int    `json:"character_id"`
	QuestID     string `json:"quest_id"`
	Status      string `json:"status"` // "unlocked", "in_progress", "completed"
	UnlockedAt  int64  `json:"unlocked_at"` // Unix timestamp
	CompletedAt int64  `json:"completed_at,omitempty"`
}

// QuestUnlockResult contains the result of a quest unlock attempt
type QuestUnlockResult struct {
	Unlocked    bool   `json:"unlocked"`
	QuestID     string `json:"quest_id,omitempty"`
	QuestName   string `json:"quest_name,omitempty"`
	RevealText  string `json:"reveal_text,omitempty"`
	XPGained    int    `json:"xp_gained"`
	AlreadyUnlocked bool `json:"already_unlocked"`
}

// QuestRegistry holds all defined quests
var QuestRegistry = map[string]*Quest{
	"quest_fountain_secret": {
		ID:          "quest_fountain_secret",
		Name:        "The Fountain's Secret",
		Type:        QuestTypeSecret,
		Description: "You've discovered a hidden compartment in the fountain...",
		Hidden:      true,
		Rewards: QuestRewards{
			XP: 50,
		},
		ExamineTrigger: &ExamineTrigger{
			Target:          "Stone Fountain",
			MinExamineLevel: 75,
			RevealText:      "As you examine the fountain closely, you notice faint runes glowing. Behind the water, you glimpse a hidden compartment!",
			XPReward:        5,
		},
	},
	// Add more quests as needed
}

// GetExamineQuests returns all quests that can be unlocked by examining
func GetExamineQuests() []*Quest {
	var quests []*Quest
	for _, q := range QuestRegistry {
		if q.ExamineTrigger != nil {
			quests = append(quests, q)
		}
	}
	return quests
}

// GetQuestByTarget returns quests triggered by examining a specific target
func GetQuestByTarget(targetName string) []*Quest {
	var quests []*Quest
	for _, q := range QuestRegistry {
		if q.ExamineTrigger != nil && q.ExamineTrigger.Target == targetName {
			quests = append(quests, q)
		}
	}
	return quests
}

// GetQuestByID returns a quest by its ID
func GetQuestByID(id string) *Quest {
	return QuestRegistry[id]
}