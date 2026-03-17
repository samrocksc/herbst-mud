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

	// Create default skills
	defaultSkills := []*db.SkillCreate{
		client.Skill.Create().
			SetName("Slash").
			SetDescription("A basic sword attack").
			SetType("combat").
			SetCost(1).
			SetCooldown(0).
			SetPower(10),
		client.Skill.Create().
			SetName("Power Strike").
			SetDescription("A strong attack that deals extra damage").
			SetType("combat").
			SetCost(2).
			SetCooldown(3).
			SetPower(20),
		client.Skill.Create().
			SetName("Shield Block").
			SetDescription("Defensive stance that reduces damage").
			SetType("defensive").
			SetCost(1).
			SetCooldown(2).
			SetPower(5),
		client.Skill.Create().
			SetName("Quick Strike").
			SetDescription("Fast attack with low damage").
			SetType("combat").
			SetCost(1).
			SetCooldown(1).
			SetPower(5),
		client.Skill.Create().
			SetName("Heal").
			SetDescription("Restore health").
			SetType("utility").
			SetCost(2).
			SetCooldown(5).
			SetPower(-15),
		client.Skill.Create().
			SetName("Fireball").
			SetDescription("Launch a ball of fire").
			SetType("combat").
			SetCost(3).
			SetCooldown(4).
			SetPower(25),
		client.Skill.Create().
			SetName("Ice Shield").
			SetDescription("Protect against cold attacks").
			SetType("defensive").
			SetCost(2).
			SetCooldown(6).
			SetPower(10),
		client.Skill.Create().
			SetName("Sprint").
			SetDescription("Move faster for a short time").
			SetType("utility").
			SetCost(1).
			SetCooldown(10).
			SetPower(0),
	}

	for _, skill := range defaultSkills {
		if _, err := skill.Save(ctx); err != nil {
			log.Printf("Warning: failed to create skill: %v", err)
		}
	}

	// Create default talents
	defaultTalents := []*db.TalentCreate{
		client.Talent.Create().
			SetName("Warrior's Might").
			SetDescription("Increase strength by 5").
			SetRequirements(map[string]int{"level": 5, "strength": 10}),
		client.Talent.Create().
			SetName("Agile Fighter").
			SetDescription("Increase dexterity by 5").
			SetRequirements(map[string]int{"level": 5, "dexterity": 10}),
		client.Talent.Create().
			SetName("Iron Will").
			SetDescription("Increase constitution by 5").
			SetRequirements(map[string]int{"level": 5, "constitution": 10}),
		client.Talent.Create().
			SetName("Arcane Mind").
			SetDescription("Increase intelligence by 5").
			SetRequirements(map[string]int{"level": 5, "intelligence": 10}),
		client.Talent.Create().
			SetName("Leader").
			SetDescription("Increase charisma by 5").
			SetRequirements(map[string]int{"level": 5, "charisma": 10}),
		client.Talent.Create().
			SetName("Sword Mastery").
			SetDescription("+5 damage with sword weapons").
			SetRequirements(map[string]int{"level": 10, "strength": 15}),
		client.Talent.Create().
			SetName("Spell Mastery").
			SetDescription("+10% magic damage").
			SetRequirements(map[string]int{"level": 10, "intelligence": 15}),
		client.Talent.Create().
			SetName("Tank").
			SetDescription("+10% damage reduction").
			SetRequirements(map[string]int{"level": 15, "constitution": 20}),
	}

	for _, talent := range defaultTalents {
		if _, err := talent.Save(ctx); err != nil {
			log.Printf("Warning: failed to create talent: %v", err)
		}
	}

	log.Println("Skills and talents initialized successfully")
	return nil
}