package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Talent holds the schema definition for the Talent entity.
type Talent struct {
	ent.Schema
}

// Fields of the Talent.
func (Talent) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Unique(),
		field.String("description"),
		field.String("requirements").
			Optional().
			Comment("JSON string of prerequisites (skills, levels, etc.)"),
		// Effect system fields
		field.String("effect_type").
			Default("").
			Comment("heal|damage|dot|buff_armor|buff_dodge|buff_crit|debuff"),
		field.Int("effect_value").
			Default(0).
			Comment("Amount: HP healed, damage dealt, armor bonus, etc."),
		field.Int("effect_duration").
			Default(0).
			Comment("Duration in ticks (0 = instant)"),
		field.Int("cooldown").
			Default(0).
			Comment("Cooldown in ticks before can use again"),
		field.Int("mana_cost").
			Default(0).
			Comment("Mana cost to use"),
		field.Int("stamina_cost").
			Default(0).
			Comment("Stamina cost to use"),
	}
}

// Edges of the Talent.
func (Talent) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("characters", CharacterTalent.Type),
		edge.To("available_to_characters", AvailableTalent.Type),
	}
}