package dbinit

import (
	"context"
	"log"

	"herbst-server/db"
)

// InitSkillsAndAbilities creates default classless abilities with effects
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

	type effectDef struct {
		EffectType    string
		DamageSubtype string
		Target        string
		Value         int
		Duration      int
		ScalingStat   string
		ScalingRatio  float64
		SortOrder     int
	}

	type abilityDef struct {
		Name         string
		Desc         string
		AbilityType  string
		Cooldown     int
		ManaCost     int
		StaminaCost  int
		Effects      []effectDef
	}

	abilities := []abilityDef{
		{
			Name:        "Concentrate",
			Desc:        "Focus your mind to increase accuracy. +WIS to hit for 4 rounds.",
			AbilityType: "combat",
			Cooldown:    8,
			ManaCost:    10,
			Effects: []effectDef{
				{EffectType: "accuracy_boost", Target: "self", Duration: 4, ScalingStat: "wisdom", ScalingRatio: 0.5, SortOrder: 0},
			},
		},
		{
			Name:        "Haymaker",
			Desc:        "A powerful but reckless strike. +STR damage, -DEX to hit.",
			AbilityType: "combat",
			Cooldown:    6,
			StaminaCost: 15,
			Effects: []effectDef{
				{EffectType: "damage", Target: "enemy", Value: 15, ScalingStat: "strength", ScalingRatio: 0.5, SortOrder: 0},
				{EffectType: "debuff", Target: "self", Duration: 1, Value: 2, SortOrder: 1},
			},
		},
		{
			Name:        "Back-off",
			Desc:        "Use agility to dodge all attacks this round. Costs stamina.",
			AbilityType: "defensive",
			Cooldown:    10,
			StaminaCost: 25,
			Effects: []effectDef{
				{EffectType: "dodge_all", Target: "self", Value: 999, Duration: 1, SortOrder: 0},
			},
		},
		{
			Name:        "Scream",
			Desc:        "Release a berserker cry. -WIS/INT, +DEX/STR for 2 rounds.",
			AbilityType: "support",
			Cooldown:    12,
			ManaCost:    5,
			StaminaCost: 10,
			Effects: []effectDef{
				{EffectType: "buff", Target: "self", Duration: 2, ScalingStat: "constitution", ScalingRatio: 0.5, SortOrder: 0},
				{EffectType: "debuff", Target: "enemy", Duration: 2, ScalingStat: "constitution", ScalingRatio: 0.5, SortOrder: 1},
			},
		},
		{
			Name:        "Slap",
			Desc:        "A quick stunning strike. DEX vs CON to stun for 1 round.",
			AbilityType: "combat",
			Cooldown:    8,
			StaminaCost: 12,
			Effects: []effectDef{
				{EffectType: "stun", Target: "enemy", Duration: 1, ScalingStat: "dexterity", ScalingRatio: 0.3, SortOrder: 0},
			},
		},
	}

	for _, a := range abilities {
		ab := client.Ability.Create().
			SetName(a.Name).
			SetDescription(a.Desc).
			SetAbilityType(a.AbilityType).
			SetAbilityClass("active").
			SetCooldown(a.Cooldown).
			SetCooldownSeconds(a.Cooldown).
			SetRequirements("{}").
			SetManaCost(a.ManaCost).
			SetStaminaCost(a.StaminaCost).
			SetScalingStat(a.Effects[0].ScalingStat)

		created, err := ab.Save(ctx)
		if err != nil {
			log.Printf("Warning: failed to create ability %s: %v", a.Name, err)
			continue
		}

		for _, e := range a.Effects {
			ef := client.AbilityEffect.Create().
				SetAbilityID(created.ID).
				SetEffectType(e.EffectType).
				SetTarget(e.Target).
				SetValue(e.Value).
				SetDuration(e.Duration).
				SetSortOrder(e.SortOrder)

			if e.DamageSubtype != "" {
				ef = ef.SetDamageSubtype(e.DamageSubtype)
			}
			if e.ScalingStat != "" {
				ef = ef.SetScalingStat(e.ScalingStat)
			}
			if e.ScalingRatio > 0 {
				ef = ef.SetScalingRatio(e.ScalingRatio)
			}

			if _, err := ef.Save(ctx); err != nil {
				log.Printf("Warning: failed to create effect for %s: %v", a.Name, err)
			}
		}
	}

	log.Println("Classless abilities with effects initialized successfully")
	return nil
}