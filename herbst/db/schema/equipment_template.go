package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
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