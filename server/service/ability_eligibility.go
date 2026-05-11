package service

import (
	"context"

	"herbst-server/db"
	"herbst-server/db/ability"
	"herbst-server/db/character"
	"herbst-server/db/characterfaction"
	"herbst-server/db/charactertag"
)

// abilityEligibilityService checks whether a character is eligible for abilities
// based on faction membership and tag requirements per RFC-FACTION-SKILLS §Eligibility Check Logic.
type abilityEligibilityService struct {
	client *db.Client
}

// NewAbilityEligibilityService creates a new AbilityEligibilityService.
func NewAbilityEligibilityService(client *db.Client) AbilityEligibilityService {
	return &abilityEligibilityService{client: client}
}

// characterActiveFactionIDs returns the set of faction IDs where the character
// holds an "active" membership.
func (s *abilityEligibilityService) characterActiveFactionIDs(ctx context.Context, charID int) (map[int]bool, error) {
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
func (s *abilityEligibilityService) characterTagSet(ctx context.Context, charID int) (map[string]bool, error) {
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

// CheckEligibility returns the eligibility of a single ability for a character.
// It requires the ability's faction edge to be loaded if faction_id is set.
func (s *abilityEligibilityService) CheckEligibility(ctx context.Context, charID int, sk *db.Ability, activeFactionIDs map[int]bool, tagSet map[string]bool) AbilityEligibility {
	// 1. Check faction membership: if ability has a faction, character must be an active member
	if sk.Edges.Faction != nil {
		if !activeFactionIDs[sk.Edges.Faction.ID] {
			return AbilityEligibility{
				Eligible: false,
				Reason:   "not_active_member_of_faction",
			}
		}
	}

	// 2. Check required tag: if ability has a required_tag, character must have it
	if sk.RequiredTag != "" {
		if !tagSet[sk.RequiredTag] {
			return AbilityEligibility{
				Eligible: false,
				Reason:   "missing_required_tag:" + sk.RequiredTag,
			}
		}
	}

	return AbilityEligibility{Eligible: true}
}

// CheckEligibilityForCharacter loads everything needed and returns eligibility for all abilities.
// Returns a map of ability ID -> AbilityEligibility.
func (s *abilityEligibilityService) CheckEligibilityForCharacter(ctx context.Context, charID int) (map[int]AbilityEligibility, error) {
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

	// Load all abilities with faction edge
	abilities, err := s.client.Ability.Query().
		WithFaction().
		All(ctx)
	if err != nil {
		return nil, err
	}

	result := make(map[int]AbilityEligibility, len(abilities))
	for _, sk := range abilities {
		result[sk.ID] = s.CheckEligibility(ctx, charID, sk, activeFactionIDs, tagSet)
	}

	return result, nil
}

// GetEligiblePassiveAbilitiesForEvent returns all eligible passive abilities for a character
// that match the given proc event (on_hit, on_hit_received, on_crit, on_kill).
// Used by the proc system to determine which abilities can proc.
func (s *abilityEligibilityService) GetEligiblePassiveAbilitiesForEvent(ctx context.Context, charID int, procEvent string) ([]*db.Ability, error) {
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

	// Query passive abilities matching the proc event
	abilities, err := s.client.Ability.Query().
		Where(
			ability.AbilityClass("passive"),
			ability.ProcEvent(procEvent),
		).
		WithFaction().
		All(ctx)
	if err != nil {
		return nil, err
	}

	// Filter by eligibility
	var eligible []*db.Ability
	for _, sk := range abilities {
		el := s.CheckEligibility(ctx, charID, sk, activeFactionIDs, tagSet)
		if el.Eligible {
			eligible = append(eligible, sk)
		}
	}

	return eligible, nil
}

// AbilitiesForCharacterWithEligibility returns all abilities with eligibility info for a character.
// Each entry includes the ability's existing data plus eligibility fields.
func (s *abilityEligibilityService) AbilitiesForCharacterWithEligibility(ctx context.Context, charID int) ([]AbilityWithEligibility, error) {
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

	// Load all abilities with faction edge
	abilities, err := s.client.Ability.Query().
		WithFaction().
		Order(ability.ByName()).
		All(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]AbilityWithEligibility, len(abilities))
	for i, sk := range abilities {
		el := s.CheckEligibility(ctx, charID, sk, activeFactionIDs, tagSet)
		result[i] = AbilityWithEligibility{
			Ability:     sk,
			Eligibility: el,
		}
	}

	return result, nil
}