package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Skill holds the schema definition for the Skill entity.
type Skill struct {
	ent.Schema
}

func (Skill) Fields() []ent.Field {
	return []ent.Field{
		field.Int("world_id").
			Comment("World this skill belongs to"),
		field.String("name").
			Comment("Machine name e.g., 'blades', 'heavy_armor'"),
		field.String("display_name").
			Comment("Human-readable name e.g., 'Blades', 'Heavy Armor'"),
		field.Text("description").
			Optional().
			Comment("Description of the skill"),
		field.String("category").
			Default("weapon").
			Comment("Category: weapon, armor, craft, magic, etc."),
		field.Int("parent_skill_id").
			Optional().
			Nillable().
			Comment("Parent skill for tree structure (nullable)"),
		field.Int("max_level").
			Default(100).
			Comment("Maximum level cap for this skill"),
		field.String("xp_curve_mode").
			Default("percentage").
			Comment("XP curve mode: 'percentage' or 'hand_coded'"),
		field.JSON("xp_curve_data", map[string]interface{}{}).
			Optional().
			Comment("XP curve config: {percentage: 50} or {thresholds: [0, 1000, 2500, ...]}"),
	}
}

func (Skill) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("world", World.Type).
			Ref("skills").
			Field("world_id").
			Unique().
			Required(),
		edge.To("character_skills", CharacterSkill.Type),
		edge.From("parent", Skill.Type).
			Ref("children").
			Field("parent_skill_id").
			Unique(),
		edge.To("children", Skill.Type),
		edge.To("abilities", Ability.Type),
	}
}

func (Skill) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name", "world_id").Unique(),
	}
}