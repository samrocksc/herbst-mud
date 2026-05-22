package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type CraftingInput struct {
	EquipmentTemplateSlug string `json:"equipment_template_slug"`
	Quantity              int    `json:"quantity"`
	Consumed              bool   `json:"consumed"`
}

type CraftingOutput struct {
	EquipmentTemplateSlug string `json:"equipment_template_slug"`
	Quantity               int    `json:"quantity"`
}

type CraftingRecipe struct {
	ent.Schema
}

func (CraftingRecipe) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Unique(),
		field.String("display_name"),
		field.Text("description").Optional(),
		field.String("required_station_tag"),
		field.String("required_class").Optional(),
		field.Int("required_skill_level").Default(0),
		field.String("required_skill").Optional(),
		field.JSON("inputs", []CraftingInput{}),
		field.JSON("outputs", []CraftingOutput{}),
		field.Int("craft_time_secs").Default(3),
		field.String("world_id").Default("default"),
	}
}
