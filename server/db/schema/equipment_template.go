package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/index"
)

// EquipmentTemplate holds the schema definition for the EquipmentTemplate entity.
type EquipmentTemplate struct {
	ent.Schema
}

// Fields of the EquipmentTemplate.
func (EquipmentTemplate) Fields() []ent.Field {
	return append(equipmentTemplateBaseFields(), equipmentTemplateCombatFields()...)
}

// Edges of the EquipmentTemplate.
func (EquipmentTemplate) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("equipment", Equipment.Type),
	}
}

// Indexes of the EquipmentTemplate.
func (EquipmentTemplate) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("slug", "world_id").Unique(),
	}
}