package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
)

// Equipment holds the schema definition for the Equipment entity.
type Equipment struct {
	ent.Schema
}

// Fields of the Equipment.
func (Equipment) Fields() []ent.Field {
	return append(equipmentBaseFields(), equipmentCombatFields()...)
}

// Edges of the Equipment.
func (Equipment) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("room", Room.Type).
			Ref("equipment").
			Unique(),
		edge.From("equipmentTemplate", EquipmentTemplate.Type).
			Ref("equipment").
			Unique().
			Field("equipment_template_id"),
	}
}