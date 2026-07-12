package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// CharacterSkill holds the schema definition for the CharacterSkill join entity.
type CharacterSkill struct {
	ent.Schema
}

func (CharacterSkill) Fields() []ent.Field {
	return []ent.Field{
		field.Int("character_id").
			Comment("FK to character"),
		field.Int("skill_id").
			Comment("FK to skill"),
		field.Int("level").
			Default(0).
			Comment("Current skill level"),
		field.Int("xp").
			Default(0).
			Comment("Current skill XP"),
	}
}

func (CharacterSkill) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("character", Character.Type).
			Ref("character_skills").
			Field("character_id").
			Unique().
			Required(),
		edge.From("skill", Skill.Type).
			Ref("character_skills").
			Field("skill_id").
			Unique().
			Required(),
	}
}