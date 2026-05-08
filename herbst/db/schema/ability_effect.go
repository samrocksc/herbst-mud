package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// AbilityEffect holds the schema definition for the AbilityEffect entity.
// Effects are the building blocks of abilities — damage, heals, buffs, debuffs, etc.
type AbilityEffect struct {
	ent.Schema
}

// Fields of the AbilityEffect.
func (AbilityEffect) Fields() []ent.Field {
	return []ent.Field{
		field.String("effect_type").
			Comment("damage|heal|buff|debuff|dot|hot|stun|accuracy_boost|dodge_all"),
		field.String("damage_subtype").
			Default("").
			Comment("slashing|piercing|bludgeoning|fire|cold|lightning|poison|psychic"),
		field.String("target").
			Default("enemy").
			Comment("self|enemy|ally|area|random_enemy"),
		field.Int("value").
			Default(0).
			Comment("Base magnitude of the effect"),
		field.Int("duration").
			Default(0).
			Comment("Duration in ticks (0 = instant)"),
		field.String("scaling_stat").
			Optional().
			Comment("strength|dexterity|constitution|intelligence|wisdom"),
		field.Float("scaling_ratio").
			Default(0).
			Comment("Multiplier per point of scaling stat"),
		field.Int("sort_order").
			Default(0).
			Comment("Order of effects within an ability"),
	}
}

// Edges of the AbilityEffect.
func (AbilityEffect) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("ability", Ability.Type).
			Ref("effects").
			Unique(),
	}
}