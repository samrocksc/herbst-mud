package routes

import (
	"herbst-server/db/schema"
)

// strPtr returns a pointer to s, or nil if s is empty
func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// derefInt returns the value of p, or 0 if p is nil
func derefInt(p *int) int {
	if p == nil {
		return 0
	}
	return *p
}

// derefFloat64 returns the value of p, or 0 if p is nil
func derefFloat64(p *float64) float64 {
	if p == nil {
		return 0
	}
	return *p
}

// derefBool returns the value of p, or false if p is nil
func derefBool(p *bool) bool {
	if p == nil {
		return false
	}
	return *p
}

func inputsFromInterface(in []map[string]any) []schema.CraftingInput {
	result := make([]schema.CraftingInput, len(in))
	for i, m := range in {
		slug, _ := m["equipment_template_slug"].(string)
		q, _ := m["quantity"].(float64)
		consumed, _ := m["consumed"].(bool)
		result[i] = schema.CraftingInput{
			EquipmentTemplateSlug: slug,
			Quantity:              int(q),
			Consumed:              consumed,
		}
	}
	return result
}

func outputsFromInterface(in []map[string]any) []schema.CraftingOutput {
	result := make([]schema.CraftingOutput, len(in))
	for i, m := range in {
		slug, _ := m["equipment_template_slug"].(string)
		q, _ := m["quantity"].(float64)
		result[i] = schema.CraftingOutput{
			EquipmentTemplateSlug: slug,
			Quantity:              int(q),
		}
	}
	return result
}
