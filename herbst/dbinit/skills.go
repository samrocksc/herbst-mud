package dbinit

import (
	"context"
	"log"

	"herbst/db"
)

// InitSkillsAndTalents creates default skills and talents if they don't exist
func InitSkillsAndTalents(client *db.Client) error {
	ctx := context.Background()

	// Check if skills already exist
	existingSkills, err := client.Skill.Query().Limit(1).All(ctx)
	if err != nil {
		return err
	}
	if len(existingSkills) > 0 {
		log.Println("Skills already initialized, skipping")
		return nil
	}

	// Create the 5 classless combat skills
	defaultSkills := []*db.SkillCreate{
		client.Skill.Create().
			SetName("Concentrate").
			SetDescription("Focus your mind to increase accuracy. +WIS to hit for 4 rounds.").
			SetSkillType("combat").
			SetCost(0).
			SetCooldown(8).
			SetRequirements("{}").
			SetEffectType("concentrate").
			SetEffectValue(10).
			SetEffectDuration(4).
			SetScalingStat("wisdom").
			SetScalingPercentPerPoint(0.05).
			SetManaCost(10).
			SetStaminaCost(0).
			SetHpCost(0),
		client.Skill.Create().
			SetName("Haymaker").
			SetDescription("A powerful but reckless strike. +STR to damage, -DEX to hit.").
			SetSkillType("combat").
			SetCost(0).
			SetCooldown(6).
			SetRequirements("{}").
			SetEffectType("haymaker").
			SetEffectValue(12).
			SetEffectDuration(1).
			SetScalingStat("strength").
			SetScalingPercentPerPoint(0.05).
			SetManaCost(0).
			SetStaminaCost(15).
			SetHpCost(0),
		client.Skill.Create().
			SetName("Back-off").
			SetDescription("Use agility to dodge all attacks this round. Costs stamina.").
			SetSkillType("combat").
			SetCost(0).
			SetCooldown(10).
			SetRequirements("{}").
			SetEffectType("backoff").
			SetEffectValue(0).
			SetEffectDuration(1).
			SetScalingStat("").
			SetScalingPercentPerPoint(0).
			SetManaCost(0).
			SetStaminaCost(25).
			SetHpCost(0),
		client.Skill.Create().
			SetName("Scream").
			SetDescription("Release a berserker cry. -WIS/INT, +DEX/STR for 2 rounds.").
			SetSkillType("combat").
			SetCost(0).
			SetCooldown(12).
			SetRequirements("{}").
			SetEffectType("scream").
			SetEffectValue(0).
			SetEffectDuration(2).
			SetScalingStat("").
			SetScalingPercentPerPoint(0).
			SetManaCost(5).
			SetStaminaCost(10).
			SetHpCost(0),
		client.Skill.Create().
			SetName("Slap").
			SetDescription("A quick stunning strike. DEX vs CON to stun for 1 round.").
			SetSkillType("combat").
			SetCost(0).
			SetCooldown(8).
			SetRequirements("{}").
			SetEffectType("slap").
			SetEffectValue(8).
			SetEffectDuration(1).
			SetScalingStat("").
			SetScalingPercentPerPoint(0).
			SetManaCost(0).
			SetStaminaCost(12).
			SetHpCost(0),
	}

	for _, skill := range defaultSkills {
		if _, err := skill.Save(ctx); err != nil {
			log.Printf("Warning: failed to create skill: %v", err)
		}
	}

	log.Println("Skills and talents initialized successfully")
	return nil
}
