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
		field.String("name"),
		field.String("description"),
		field.String("ability_type"), // "combat", "magic", "utility", "defensive"
		field.Int("cost"),            // points required to learn
		field.Int("cooldown"),        // in ticks (combat) / room movements (outside)
		field.String("requirements").Optional(), // JSON: {"race": [], "class": []}
		field.Int("mana_cost"),
		field.Int("stamina_cost"),
		field.Int("hp_cost"), // HP sacrificed to use ability
	}
}

// Edges of the Ability.
func (Ability) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("characters", Character.Type),
		edge.To("effects", AbilityEffect.Type),
	}
}