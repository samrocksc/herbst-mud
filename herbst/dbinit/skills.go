package dbinit

import (
	"context"
	"log"

	"herbst/db"
)

// InitSkillsAndAbilities creates default abilities if they don't exist
func InitSkillsAndAbilities(client *db.Client) error {
	ctx := context.Background()

	existingAbilities, err := client.Ability.Query().Limit(1).All(ctx)
	if err != nil {
		return err
	}
	if len(existingAbilities) > 0 {
		log.Println("Abilities already initialized, skipping")
		return nil
	}

	defaultAbilities := []*db.AbilityCreate{
		client.Ability.Create().
			SetName("Concentrate").
			SetDescription("Focus your mind to increase accuracy. +WIS to hit for 4 rounds.").
			SetAbilityType("combat").
			SetCost(0).
			SetCooldown(8).
			SetRequirements("{}").
			SetManaCost(10).
			SetStaminaCost(0).
			SetHpCost(0),
		client.Ability.Create().
			SetName("Haymaker").
			SetDescription("A powerful but reckless strike. +STR to damage, -DEX to hit.").
			SetAbilityType("combat").
			SetCost(0).
			SetCooldown(6).
			SetRequirements("{}").
			SetManaCost(0).
			SetStaminaCost(15).
			SetHpCost(0),
		client.Ability.Create().
			SetName("Back-off").
			SetDescription("Use agility to dodge all attacks this round. Costs stamina.").
			SetAbilityType("combat").
			SetCost(0).
			SetCooldown(10).
			SetRequirements("{}").
			SetManaCost(0).
			SetStaminaCost(25).
			SetHpCost(0),
		client.Ability.Create().
			SetName("Scream").
			SetDescription("Release a berserker cry. -WIS/INT, +DEX/STR for 2 rounds.").
			SetAbilityType("combat").
			SetCost(0).
			SetCooldown(12).
			SetRequirements("{}").
			SetManaCost(5).
			SetStaminaCost(10).
			SetHpCost(0),
		client.Ability.Create().
			SetName("Slap").
			SetDescription("A quick stunning strike. DEX vs CON to stun for 1 round.").
			SetAbilityType("combat").
			SetCost(0).
			SetCooldown(8).
			SetRequirements("{}").
			SetManaCost(0).
			SetStaminaCost(12).
			SetHpCost(0),
	}

	for _, ability := range defaultAbilities {
		if _, err := ability.Save(ctx); err != nil {
			log.Printf("Warning: failed to create ability: %v", err)
		}
	}

	log.Println("Abilities initialized successfully")
	return nil
}