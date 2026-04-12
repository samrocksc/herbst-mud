package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Gender holds the schema definition for the Gender entity.
type Gender struct {
	ent.Schema
}

// Fields of the Gender.
func (Gender) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Unique().
			Comment("Internal ID: he_him, she_her, they_them"),
		field.String("display_name").
			Comment("Shown in UI: He/Him, She/Her, They/Them"),
		field.String("subject_pronoun").
			Comment("he, she, they"),
		field.String("object_pronoun").
			Comment("him, her, them"),
		field.String("possessive_pronoun").
			Comment("his, hers, theirs"),
	}
}

// Edges of the Gender.
func (Gender) Edges() []ent.Edge {
	return nil
}