package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// NPCTemplate holds the schema definition for the NPC Template entity.
type NPCTemplate struct {
	ent.Schema
}

// Fields of the NPCTemplate.
func (NPCTemplate) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Unique(),
		field.String("name"),
		field.Text("description"),
		field.String("race"),
		field.Enum("disposition").
			Values("hostile", "friendly", "neutral").
			Default("neutral"),
		field.Int("level").
			Default(1),
		field.JSON("skills", map[string]int{}),
		field.Strings("trades_with"),
		field.Text("greeting"),
	}
}