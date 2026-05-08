package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Ability holds the schema definition for the Ability entity.
// Abilities are active or passive actions (Concentrate, Haymaker, Fireball, etc.)
// that characters can use in combat. Formerly named "Skill".
type Ability struct {
	ent.Schema
}

// Fields of the Ability.
func (Ability) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Unique(),
		field.String("description"),
		field.String("ability_type").
			Comment("e.g., combat, magic, utility, defensive"),
		field.Int("cost").
			Default(0).
			Comment("Points cost to learn/unlearn"),
		field.Int("cooldown").
			Default(0).
			Comment("Cooldown in ticks (legacy, use cooldown_seconds for actives)"),
		field.String("requirements").
			Optional().
			Comment("JSON string of prerequisites"),
		// Effect system fields (flat, will be migrated to AbilityEffect entity)
		field.String("effect_type").
			Default("").
			Comment("Generic type: damage, heal, buff, debuff, dot, hot, stun, accuracy_boost, dodge_all"),
		field.Int("effect_value").
			Default(0).
			Comment("Base damage/heal amount"),
		field.Int("effect_duration").
			Default(0).
			Comment("Duration in ticks (0 = instant)"),
		field.String("scaling_stat").
			Optional().
			Comment("Which stat scales: wisdom, strength, dexterity, constitution, intelligence"),
		field.Float("scaling_percent_per_point").
			Default(0).
			Comment("% bonus per point of the scaling stat (e.g., 0.05 = +5% per stat point)"),
		field.Int("mana_cost").
			Default(0),
		field.Int("stamina_cost").
			Default(0),
		field.Int("hp_cost").
			Default(0).
			Comment("HP sacrificed to use ability"),
		// Faction and classification fields
		field.String("slug").
			Unique().
			Optional().
			Comment("Globally unique ability identifier e.g., foot_clan_power_strike"),
		field.String("required_tag").
			Optional().
			Comment("Tag required to unlock beyond faction membership"),
		field.String("ability_class").
			Default("active").
			Comment("active or passive (includes former talents)"),
		field.Float("proc_chance").
			Default(0).
			Comment("For passives: % chance to proc (0.15 = 15%)"),
		field.String("proc_event").
			Optional().
			Comment("What triggers the proc: on_hit, on_hit_received, on_crit, on_kill"),
		field.Int("cooldown_seconds").
			Default(0).
			Comment("For actives: cooldown in seconds"),
	}
}

// Edges of the Ability.
func (Ability) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("characters", CharacterAbility.Type),
		edge.To("npc_abilities", NPCAbility.Type),
		edge.To("effects", AbilityEffect.Type),
		edge.From("faction", Faction.Type).
			Ref("abilities").
			Unique(),
	}
}