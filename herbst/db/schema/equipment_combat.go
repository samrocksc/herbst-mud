package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

/** Combat-related fields for Equipment (weapon/armor stats). */
func equipmentCombatFields() []ent.Field {
	return []ent.Field{
		field.Int("minDamage").
			Default(0).
			Comment("Minimum damage for weapons"),
		field.Int("maxDamage").
			Default(0).
			Comment("Maximum damage for weapons"),
		field.String("weaponType").
			Default("").
			Comment("Type of weapon: sword, dagger, staff, etc."),
		field.String("classRestriction").
			Default("").
			Comment("Class that can use this weapon"),
		field.Bool("guaranteedDrop").
			Default(false).
			Comment("Always drops from certain NPCs"),
		field.Int("armor_rating").
			Default(0).
			Comment("Flat AC bonus when equipped"),
		field.String("armor_type").
			Default("").
			Comment("light|cloth|heavy or empty"),
		field.JSON("stats", map[string]int{}).
			Default(map[string]int{}).
			Comment("Passive stat bonuses when equipped"),
		field.String("rarity").
			Default("common").
			Comment("common|uncommon|rare|epic|legendary"),
		field.String("skill_requirement").
			Default("").
			Comment("Which skill governs this item"),
		field.Int("skill_requirement_level").
			Default(0).
			Comment("Min skill level for full effect"),
		field.Int("damage_dice_count").
			Default(0).
			Comment("Number of dice (0 = not a weapon)"),
		field.Int("damage_dice_sides").
			Default(0).
			Comment("Sides per die"),
		field.Int("damage_bonus").
			Default(0).
			Comment("Flat damage modifier"),
		field.String("damage_type").
			Default("").
			Comment("slashing|piercing|bludgeoning|fire|etc."),
		field.Bool("is_two_handed").
			Default(false).
			Comment("Occupies both hand slots"),
	}
}