package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Ability holds the schema definition for the Ability entity.
// Abilities are active or passive actions (Concentrate, Haymaker, Fireball, etc.)
// that characters can use in combat. Formerly named "Skill".
type Ability struct {
	ent.Schema
}

// Fields of the Ability.
func (Ability) Fields() []ent.Field {
	return []ent.Field{
		field.String("world_id").
			Default("1").
			Comment("World this ability belongs to (for multi-world support)"),
		field.String("name").
			Comment("Name of the ability, unique within each world"),
		field.String("description"),
		field.String("ability_type").
			Comment("e.g., combat, magic, utility, defensive"),
		field.Int("cost").
			Default(0).
			Comment("Points cost to learn/unlearn"),
		field.Int("cooldown").
			Default(0).
			Comment("Cooldown in ticks (legacy, use cooldown_seconds for actives)"),
		field.String("requirements").
			Optional().
			Comment("JSON string of prerequisites"),
		field.Int("mana_cost").
			Default(0),
		field.Int("stamina_cost").
			Default(0),
		field.Int("hp_cost").
			Default(0).
			Comment("HP sacrificed to use ability"),
		field.String("slug").
			Optional().
			Comment("World-unique ability identifier e.g., foot_clan_power_strike"),
		field.String("required_tag").
			Optional().
			Comment("Tag required to unlock beyond faction membership"),
		field.String("ability_class").
			Default("active").
			Comment("active or passive (includes former talents)"),
		field.Float("proc_chance").
			Default(0).
			Comment("For passives: % chance to proc (0.15 = 15%)"),
		field.String("proc_event").
			Optional().
			Comment("What triggers the proc: on_hit, on_hit_received, on_crit, on_kill"),
		field.Int("cooldown_seconds").
			Default(0).
			Comment("For actives: cooldown in seconds"),
		field.Int("required_skill_id").
			Optional().
			Nillable().
			Comment("FK to skills table — which skill must be leveled to unlock this ability"),
		field.Int("required_skill_level").
			Default(0).
			Comment("Minimum skill level required to unlock this ability (0 = no requirement)"),
	}
}

// Edges of the Ability.
func (Ability) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("characters", CharacterAbility.Type),
		edge.To("npc_abilities", NPCAbility.Type),
		edge.To("effects", AbilityEffect.Type),
		edge.From("faction", Faction.Type).
			Ref("abilities").
			Unique(),
		edge.From("required_skill", Skill.Type).
			Ref("abilities").
			Field("required_skill_id").
			Unique(),
	}
}