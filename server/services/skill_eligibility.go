package services

import (
	"context"

	"herbst-server/db"
	"herbst-server/db/character"
	"herbst-server/db/characterfaction"
	"herbst-server/db/charactertag"
	"herbst-server/db/skill"
)

// SkillEligibility holds the eligibility status of a skill for a character.
type SkillEligibility struct {
	Eligible bool   `json:"eligible"`
	Reason   string `json:"reason,omitempty"`
}

// SkillEligibilityService checks whether a character is eligible for skills
// based on faction membership and tag requirements per RFC-FACTION-SKILLS §Eligibility Check Logic.
type SkillEligibilityService struct {
	client *db.Client
}

// NewSkillEligibilityService creates a new SkillEligibilityService.
func NewSkillEligibilityService(client *db.Client) *SkillEligibilityService {
	return &SkillEligibilityService{client: client}
}

// characterActiveFactionIDs returns the set of faction IDs where the character
// holds an "active" membership.
func (s *SkillEligibilityService) characterActiveFactionIDs(ctx context.Context, charID int) (map[int]bool, error) {
	memberships, err := s.client.CharacterFaction.Query().
		Where(
			characterfaction.HasCharacterWith(character.ID(charID)),
			characterfaction.StatusEQ("active"),
		).
		All(ctx)
	if err != nil {
		return nil, err
	}
	result := make(map[int]bool, len(memberships))
	for _, m := range memberships {
		if m.Edges.Faction != nil {
			result[m.Edges.Faction.ID] = true
		}
	}
	return result, nil
}

// characterTagSet returns the set of tag strings the character possesses.
func (s *SkillEligibilityService) characterTagSet(ctx context.Context, charID int) (map[string]bool, error) {
	tags, err := s.client.CharacterTag.Query().
		Where(charactertag.HasCharacterWith(character.ID(charID))).
		All(ctx)
	if err != nil {
		return nil, err
	}
	result := make(map[string]bool, len(tags))
	for _, t := range tags {
		result[t.Tag] = true
	}
	return result, nil
}

// CheckEligibility returns the eligibility of a single skill for a character.
// It requires the skill's faction edge to be loaded if faction_id is set.
func (s *SkillEligibilityService) CheckEligibility(ctx context.Context, charID int, sk *db.Skill, activeFactionIDs map[int]bool, tagSet map[string]bool) SkillEligibility {
	// 1. Check faction membership: if skill has a faction, character must be an active member
	if sk.Edges.Faction != nil {
		if !activeFactionIDs[sk.Edges.Faction.ID] {
			return SkillEligibility{
				Eligible: false,
				Reason:   "not_active_member_of_faction",
			}
		}
	}

	// 2. Check required tag: if skill has a required_tag, character must have it
	if sk.RequiredTag != "" {
		if !tagSet[sk.RequiredTag] {
			return SkillEligibility{
				Eligible: false,
				Reason:   "missing_required_tag:" + sk.RequiredTag,
			}
		}
	}

	return SkillEligibility{Eligible: true}
}

// CheckEligibilityForCharacter loads everything needed and returns eligibility for all skills.
// Returns a map of skill ID -> SkillEligibility.
func (s *SkillEligibilityService) CheckEligibilityForCharacter(ctx context.Context, charID int) (map[int]SkillEligibility, error) {
	// Load character's active faction memberships
	activeFactionIDs, err := s.characterActiveFactionIDs(ctx, charID)
	if err != nil {
		return nil, err
	}

	// Load character's tags
	tagSet, err := s.characterTagSet(ctx, charID)
	if err != nil {
		return nil, err
	}

	// Load all skills with faction edge
	skills, err := s.client.Skill.Query().
		WithFaction().
		All(ctx)
	if err != nil {
		return nil, err
	}

	result := make(map[int]SkillEligibility, len(skills))
	for _, sk := range skills {
		result[sk.ID] = s.CheckEligibility(ctx, charID, sk, activeFactionIDs, tagSet)
	}

	return result, nil
}

// GetEligiblePassiveSkillsForEvent returns all eligible passive skills for a character
// that match the given proc event (on_hit, on_hit_received, on_crit, on_kill).
// Used by the proc system to determine which skills can proc.
func (s *SkillEligibilityService) GetEligiblePassiveSkillsForEvent(ctx context.Context, charID int, procEvent string) ([]*db.Skill, error) {
	// Load character's active faction memberships
	activeFactionIDs, err := s.characterActiveFactionIDs(ctx, charID)
	if err != nil {
		return nil, err
	}

	// Load character's tags
	tagSet, err := s.characterTagSet(ctx, charID)
	if err != nil {
		return nil, err
	}

	// Query passive skills matching the proc event
	skills, err := s.client.Skill.Query().
		Where(
			skill.SkillClass("passive"),
			skill.ProcEvent(procEvent),
		).
		WithFaction().
		All(ctx)
	if err != nil {
		return nil, err
	}

	// Filter by eligibility
	var eligible []*db.Skill
	for _, sk := range skills {
		el := s.CheckEligibility(ctx, charID, sk, activeFactionIDs, tagSet)
		if el.Eligible {
			eligible = append(eligible, sk)
		}
	}

	return eligible, nil
}

// SkillsForCharacterWithEligibility returns all skills with eligibility info for a character.
// Each entry includes the skill's existing data plus eligibility fields.
func (s *SkillEligibilityService) SkillsForCharacterWithEligibility(ctx context.Context, charID int) ([]SkillWithEligibility, error) {
	// Load character's active faction memberships
	activeFactionIDs, err := s.characterActiveFactionIDs(ctx, charID)
	if err != nil {
		return nil, err
	}

	// Load character's tags
	tagSet, err := s.characterTagSet(ctx, charID)
	if err != nil {
		return nil, err
	}

	// Load all skills with faction edge
	skills, err := s.client.Skill.Query().
		WithFaction().
		Order(skill.ByName()).
		All(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]SkillWithEligibility, len(skills))
	for i, sk := range skills {
		el := s.CheckEligibility(ctx, charID, sk, activeFactionIDs, tagSet)
		result[i] = SkillWithEligibility{
			Skill:       sk,
			Eligibility: el,
		}
	}

	return result, nil
}

// SkillWithEligibility pairs a skill with its eligibility status.
type SkillWithEligibility struct {
	Skill       *db.Skill
	Eligibility SkillEligibility
}