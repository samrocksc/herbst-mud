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
		field.String("effect_type"),              // damage|heal|buff|debuff|dot|hot|stun|accuracy_boost|dodge_all
		field.Int("effect_value"),                // base damage/heal amount
		field.Int("effect_duration"),             // duration in ticks, 0 = instant
		field.String("scaling_stat").Optional(),  // wisdom|strength|dexterity|constitution|intelligence
		field.Float("scaling_percent_per_point"), // e.g. 0.05 = +5% per stat point
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