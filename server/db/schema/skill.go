package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Skill holds the schema definition for the Skill entity.
type Skill struct {
	ent.Schema
}

// Fields of the Skill.
func (Skill) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Unique(),
		field.String("description"),
		field.String("skill_type").
			Comment("e.g., combat, magic, utility"),
		field.Int("cost").
			Default(0).
			Comment("Points cost to learn/unlearn"),
		field.Int("cooldown").
			Default(0).
			Comment("Cooldown in ticks"),
		field.String("requirements").
			Optional().
			Comment("JSON string of prerequisites"),
		// Effect system fields
		field.String("effect_type").
			Default("").
			Comment("Handler key: concentrate, haymaker, backoff, scream, slap"),
		field.Int("effect_value").
			Default(0).
			Comment("Base damage/heal amount"),
		field.Int("effect_duration").
			Default(0).
			Comment("Duration in ticks (0 = instant)"),
		field.String("scaling_stat").
			Optional().
			Comment("Which stat the skill scales off: wisdom, strength, dexterity, constitution, intelligence"),
		field.Float("scaling_percent_per_point").
			Default(0).
			Comment("% bonus per point of the scaling stat (e.g., 0.05 = +5% per stat point)"),
		field.Int("mana_cost").
			Default(0),
		field.Int("stamina_cost").
			Default(0),
		field.Int("hp_cost").
			Default(0).
			Comment("HP sacrificed to use skill"),
		// Faction skill fields
		field.String("slug").
			Unique().
			Optional().
			Comment("Globally unique skill identifier e.g., foot_clan_power_strike"),
		field.String("required_tag").
			Optional().
			Comment("Tag required to unlock this skill beyond faction membership"),
		field.String("skill_class").
			Default("active").
			Comment("passive or active"),
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

// Edges of the Skill.
func (Skill) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("characters", CharacterSkill.Type),
		edge.To("npc_skills", NPCSkill.Type),
		edge.From("faction", Faction.Type).
			Ref("skills").
			Unique(),
	}
}