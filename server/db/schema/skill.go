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
			Comment("Cooldown in seconds"),
		field.String("requirements").
			Optional().
			Comment("JSON string of prerequisites"),
		// Effect system fields
		field.String("effect_type").
			Default("").
			Comment("heal|damage|dot|buff_armor|buff_dodge|buff_crit|passive"),
		field.Int("effect_value").
			Default(0).
			Comment("Amount: HP healed, damage dealt, armor bonus, etc."),
		field.Int("effect_duration").
			Default(0).
			Comment("Duration in ticks (0 = instant)"),
		field.Int("mana_cost").
			Default(0).
			Comment("Mana cost to use"),
		field.Int("stamina_cost").
			Default(0).
			Comment("Stamina cost to use"),
	}
}

// Edges of the Skill.
func (Skill) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("characters", CharacterSkill.Type),
	}
}