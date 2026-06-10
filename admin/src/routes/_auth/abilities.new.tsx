import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import {
  useCreateAbility,
  useCreateAbilityEffect,
  type AbilityInput,
  type AbilityEffect,
  type CreateAbilityEffectInput,
} from "../../hooks/useAbilities";
import { useTags } from "../../hooks/useTags";
import { PageHeader } from "../../components/PageHeader";
import { Button } from "../../components/Button";
import { SearchableSelect } from "../../components/SearchableSelect";
import {
  FormField,
  NumberField,
  TextareaField,
  SelectField,
} from "../../components/FormFields";
import { showToast } from "../../components/Toast";
import { PageContainer } from "../../components/PageContainer";

export const Route = createFileRoute("/_auth/abilities/new")({
  component: CreateAbilityPage,
});

const ABILITY_TYPE_OPTS = [
  { value: "combat", label: "Combat" },
  { value: "magic", label: "Magic" },
  { value: "utility", label: "Utility" },
  { value: "healing", label: "Healing" },
  { value: "support", label: "Support" },
  { value: "defensive", label: "Defensive" },
];

const ABILITY_CLASS_OPTS = [
  { value: "active", label: "Active" },
  { value: "passive", label: "Passive" },
  { value: "toggle", label: "Toggle" },
];

const EFFECT_TYPE_OPTS = [
  { value: "damage", label: "Damage" },
  { value: "heal", label: "Heal" },
  { value: "buff", label: "Buff" },
  { value: "debuff", label: "Debuff" },
  { value: "dot", label: "Damage Over Time" },
  { value: "hot", label: "Heal Over Time" },
  { value: "stun", label: "Stun" },
  { value: "accuracy_boost", label: "Accuracy Boost" },
  { value: "dodge_all", label: "Dodge All" },
];

const DAMAGE_SUBTYPE_OPTS = [
  { value: "slashing", label: "Slashing" },
  { value: "piercing", label: "Piercing" },
  { value: "bludgeoning", label: "Bludgeoning" },
  { value: "fire", label: "Fire" },
  { value: "cold", label: "Cold" },
  { value: "lightning", label: "Lightning" },
  { value: "poison", label: "Poison" },
  { value: "psychic", label: "Psychic" },
];

const TARGET_OPTS = [
  { value: "self", label: "Self" },
  { value: "enemy", label: "Enemy" },
  { value: "ally", label: "Ally" },
  { value: "area", label: "Area" },
  { value: "random_enemy", label: "Random Enemy" },
];

const SCALING_STAT_OPTS = [
  { value: "", label: "— None —" },
  { value: "strength", label: "Strength" },
  { value: "dexterity", label: "Dexterity" },
  { value: "constitution", label: "Constitution" },
  { value: "intelligence", label: "Intelligence" },
  { value: "wisdom", label: "Wisdom" },
];

// Refinement #1 + #2: removed legacy `cost` and `cooldown` (ticks) fields.
// Refinement #12: only fields actually used by combat code remain.
const EMPTY_ABILITY: AbilityInput = {
  name: "",
  description: "",
  ability_type: "combat",
  requirements: "",
  cooldown_seconds: 0,
  mana_cost: 0,
  stamina_cost: 0,
  hp_cost: 0,
  proc_chance: 0,
  proc_event: "",
  ability_class: "active",
  required_tag: "",
  slug: "",
  caster_message: "",
  recipient_message: "",
};

const EMPTY_EFFECT: AbilityEffect = {
  id: 0,
  effect_type: "damage",
  damage_subtype: "slashing",
  target: "enemy",
  value: 0,
  duration: 0,
  scaling_stat: "",
  scaling_ratio: 0,
  sort_order: 0,
  effect_message: "",
};

export function CreateAbilityPage() {
  const navigate = useNavigate();
  const createAbility = useCreateAbility();
  const createAbilityEffect = useCreateAbilityEffect();
  const { data: availableTags } = useTags();
  const [formData, setFormData] = useState<AbilityInput>(EMPTY_ABILITY);
  const [effects, setEffects] = useState<AbilityEffect[]>([]);
  const [currentEffect, setCurrentEffect] = useState<AbilityEffect>(EMPTY_EFFECT);

  const set = (patch: Partial<AbilityInput>) => setFormData((prev) => ({ ...prev, ...patch }));
  const setEffect = (patch: Partial<AbilityEffect>) => setCurrentEffect((prev) => ({ ...prev, ...patch }));

  // Refinement #7: when effect type changes, default the damage_subtype
  // to the first sensible value (no damage for heals/buffs).
  const handleEffectTypeChange = (newType: string) => {
    const damageDefault = newType === "damage" ? "slashing" : newType === "dot" ? "fire" : "";
    setEffect({ effect_type: newType, damage_subtype: damageDefault });
  };

  const handleEffectSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const newEffect: AbilityEffect = {
      ...currentEffect,
      // Refinement #7: scrub damage_subtype when not a damage effect
      damage_subtype:
        currentEffect.effect_type === "damage" || currentEffect.effect_type === "dot"
          ? currentEffect.damage_subtype
          : "",
      sort_order: effects.length,
    };
    setEffects((prev) => [...prev, newEffect]);
    setCurrentEffect(EMPTY_EFFECT);
  };

  const removeEffect = (index: number) => {
    setEffects((prev) => prev.filter((_, i) => i !== index));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const ability = await createAbility.mutateAsync(formData);

      for (const effect of effects) {
        await createAbilityEffect.mutateAsync({
          abilityId: ability.id,
          input: {
            effect_type: effect.effect_type,
            damage_subtype: effect.damage_subtype,
            target: effect.target,
            value: effect.value,
            duration: effect.duration,
            scaling_stat: effect.scaling_stat || "",
            scaling_ratio: effect.scaling_ratio,
            sort_order: effect.sort_order,
            effect_message: effect.effect_message || "",
          },
        });
      }

      showToast("Ability created with effects", "success");
      navigate({ to: "/abilities" });
    } catch {
      // Error is toasted by global onError handler
    }
  };

  // Refinement #6: only show Proc Settings for passive abilities
  const isPassive = formData.ability_class === "passive";
  // Refinement #7: only show Damage Type for damage/dot
  const isDamageEffect =
    currentEffect.effect_type === "damage" || currentEffect.effect_type === "dot";

  return (
    <PageContainer>
      <PageHeader title="Create Ability" showBack backTo="/abilities" />
      <p className="text-sm text-muted mb-4">
        Abilities are actions characters perform in combat — fireballs, healing spells, melee
        strikes, defensive buffs. Set costs, cooldowns, and effects below. Add multiple effects for
        complex abilities (e.g. Bless adds an attack bonus AND a save bonus).
      </p>
      <div className="card bg-surface p-6 border border-border rounded">
        <form onSubmit={handleSubmit} className="space-y-6">
          {/* Refinement #12: sectioned with a left-border accent + heading */}
          <section>
            <h3 className="text-text font-semibold mb-1 border-l-4 border-primary pl-3">Basic Information</h3>
            <p className="text-xs text-muted mb-3 ml-1">Name, type, class. Required: Name.</p>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <FormField label="Name" value={formData.name} onChange={(v) => set({ name: v })} required placeholder="e.g., Fireball" />
              <FormField label="Slug (optional)" value={formData.slug ?? ""} onChange={(v) => set({ slug: v })} placeholder="Auto-generated from name if empty" />
              <SelectField label="Ability Type" value={formData.ability_type} onChange={(v) => set({ ability_type: v })} options={ABILITY_TYPE_OPTS} tooltip="Combat = melee, Magic = spells, Utility = misc actions, Healing = restore HP, Support = buff allies, Defensive = ward or shield" />
              <SearchableSelect
                label="Required Tag (optional)"
                options={(availableTags ?? []).map((t) => ({ id: t.name, name: t.name }))}
                value={formData.required_tag || ""}
                onChange={(v) => set({ required_tag: v })}
                placeholder="No tag required"
              />
              <SelectField label="Ability Class" value={formData.ability_class} onChange={(v) => set({ ability_class: v })} options={ABILITY_CLASS_OPTS} tooltip="Active = use on your turn. Passive = auto-triggers from Proc Settings. Toggle = turn on/off." />
            </div>
            <div className="mt-4">
              <TextareaField label="Description" value={formData.description} onChange={(v) => set({ description: v })} rows={3} placeholder="What the ability does — shown in the help panel and combat log." />
            </div>
          </section>

          {/* Refinement #1 + #2 + #12: removed legacy Cost/Cooldown ticks, single Cooldown (s) */}
          <section>
            <h3 className="text-text font-semibold mb-1 border-l-4 border-primary pl-3">Costs & Cooldown</h3>
            <p className="text-xs text-muted mb-3 ml-1">
              Active abilities cost mana, stamina, or HP to use. Set to 0 for abilities that don't
              consume resources.
            </p>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <NumberField
                label="Cooldown (s)"
                tooltip="Cooldown in seconds before the ability can be used again. 0 = no cooldown."
                value={formData.cooldown_seconds}
                onChange={(v) => set({ cooldown_seconds: v })}
              />
              <FormField
                label="Unlock Tags (JSON)"
                tooltip='JSON prerequisites for unlocking this ability. Example: {"tags":["level:5"]} means the character must be level 5 or higher.'
                value={formData.requirements}
                onChange={(v) => set({ requirements: v })}
                placeholder='{"tags":["level:5"]}'
              />
            </div>
            <div className="grid grid-cols-3 gap-4 mt-4">
              <NumberField label="Mana Cost" tooltip="Mana points consumed when using this ability. Use for magical abilities." value={formData.mana_cost} onChange={(v) => set({ mana_cost: v })} />
              <NumberField label="Stamina Cost" tooltip="Stamina points consumed when using this ability. Use for physical abilities." value={formData.stamina_cost} onChange={(v) => set({ stamina_cost: v })} />
              <NumberField label="HP Cost" tooltip="HP sacrificed to use the ability. Use for blood magic or self-damaging abilities." value={formData.hp_cost} onChange={(v) => set({ hp_cost: v })} />
            </div>
          </section>

          {/* Refinement #6: only show for passive abilities */}
          {isPassive && (
            <section>
              <h3 className="text-text font-semibold mb-1 border-l-4 border-warning pl-3">Proc Settings</h3>
              <p className="text-xs text-muted mb-3 ml-1">
                Passives trigger automatically from combat events. Set the chance and the event that
                fires them.
              </p>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <NumberField
                  label="Proc Chance (0–1)"
                  tooltip="Probability of triggering on each relevant event. 0.15 = 15% chance."
                  value={formData.proc_chance}
                  onChange={(v) => set({ proc_chance: v })}
                  step={0.01}
                />
                <FormField
                  label="Proc Event"
                  tooltip="Combat event that triggers this passive: on_hit, on_hit_received, on_crit, on_kill."
                  value={formData.proc_event}
                  onChange={(v) => set({ proc_event: v })}
                  placeholder="e.g., on_hit, on_kill"
                />
              </div>
            </section>
          )}

          {/* Refinement #7: collapsed to one row, no double-labels */}
          <section>
            <h3 className="text-text font-semibold mb-1 border-l-4 border-primary pl-3">Combat Messages</h3>
            <p className="text-xs text-muted mb-3 ml-1">
              Templates shown in the combat log. Use <code className="bg-muted px-1 rounded">{'{actor}'}</code>{" "}
              for the user and <code className="bg-muted px-1 rounded">{'{target}'}</code> for the recipient.
            </p>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <FormField
                label="Caster Message"
                tooltip='What the actor sees. Example: "{actor} casts fireball"'
                value={formData.caster_message}
                onChange={(v) => set({ caster_message: v })}
                placeholder="e.g., {actor} casts fireball"
              />
              <FormField
                label="Recipient Message"
                tooltip='What the target sees. Example: "{actor} casts fireball at you"'
                value={formData.recipient_message}
                onChange={(v) => set({ recipient_message: v })}
                placeholder="e.g., {actor} casts fireball at you"
              />
            </div>
          </section>

          {/* Refinement #7, #8, #9: per-effect-type-conditional fields, card preview, no double-labels */}
          <section>
            <h3 className="text-text font-semibold mb-1 border-l-4 border-primary pl-3">Effects</h3>
            <p className="text-xs text-muted mb-3 ml-1">
              Define what the ability does. Add multiple effects for complex abilities (Bless =
              AttackBonus + SaveBonus). Effects are saved when you click "Create Ability" below.
            </p>

            <div className="bg-surface-muted/50 border border-border rounded p-4 space-y-4 mb-4">
              <h4 className="text-text font-semibold mb-2 text-sm">Add Effect</h4>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <SelectField
                  label="Effect Type"
                  value={currentEffect.effect_type}
                  onChange={handleEffectTypeChange}
                  options={EFFECT_TYPE_OPTS}
                  tooltip="Damage = direct HP loss. Heal = direct HP gain. Buff/Debuff = stat changes. DoT/HoT = over time. Stun = skip turn. Accuracy/Dodge = hit/miss modifiers."
                />
                <SelectField
                  label="Target"
                  value={currentEffect.target}
                  onChange={(v) => setEffect({ target: v })}
                  options={TARGET_OPTS}
                  tooltip="Who the effect hits. Self for self-buffs, Ally for healing allies, Area for AoE."
                />
              </div>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {isDamageEffect ? (
                  <SelectField
                    label="Damage Type"
                    value={currentEffect.damage_subtype}
                    onChange={(v) => setEffect({ damage_subtype: v })}
                    options={DAMAGE_SUBTYPE_OPTS}
                    tooltip="Elemental type. Affects resistance/weakness calculations."
                  />
                ) : null}
                <NumberField
                  label="Value / Amount"
                  value={currentEffect.value}
                  onChange={(v) => setEffect({ value: v })}
                  tooltip="Numeric effect magnitude: damage amount, heal amount, buff bonus, etc."
                />
              </div>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <NumberField
                  label="Duration (ticks)"
                  value={currentEffect.duration}
                  onChange={(v) => setEffect({ duration: v })}
                  tooltip="Effect duration in combat ticks. 0 = instant. Buffs/debuffs typically 3-10."
                />
                <SelectField
                  label="Scaling Stat (optional)"
                  value={currentEffect.scaling_stat}
                  onChange={(v) => setEffect({ scaling_stat: v })}
                  options={SCALING_STAT_OPTS}
                  tooltip="Scales the value by a stat ratio. None = fixed value."
                />
              </div>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <NumberField
                  label="Scaling Ratio"
                  value={currentEffect.scaling_ratio}
                  onChange={(v) => setEffect({ scaling_ratio: v })}
                  step={0.1}
                  tooltip="Multiplier applied to the stat. Final = value + (stat × ratio)."
                />
              </div>
              <div className="flex justify-end">
                <Button variant="primary" type="submit" disabled={createAbilityEffect.isPending}>
                  {createAbilityEffect.isPending ? "Adding..." : "Add Effect"}
                </Button>
              </div>
            </div>

            {effects.length > 0 && (
              <div className="bg-surface-muted/50 border border-border rounded p-4">
                <h4 className="text-text font-semibold mb-3 text-sm">
                  Added Effects ({effects.length})
                </h4>
                <div className="space-y-2">
                  {effects.map((effect, index) => (
                    <div
                      key={index}
                      className="flex justify-between items-start bg-surface border border-border rounded p-3"
                    >
                      <div>
                        <div className="font-medium capitalize">
                          {effect.effect_type}
                          {effect.damage_subtype ? ` (${effect.damage_subtype})` : ""}
                        </div>
                        <div className="text-sm text-muted">
                          Target: {effect.target} • Value: {effect.value}
                          {effect.duration > 0 ? ` • Duration: ${effect.duration} ticks` : ""}
                          {effect.scaling_stat ? ` • Scales: ${effect.scaling_stat}` : ""}
                        </div>
                      </div>
                      <Button variant="danger" onClick={() => removeEffect(index)} size="sm">
                        Remove
                      </Button>
                    </div>
                  ))}
                </div>
              </div>
            )}
          </section>

          <div className="flex gap-2 justify-end pt-4 border-t border-border">
            <Button variant="secondary" onClick={() => navigate({ to: "/abilities" })}>
              Cancel
            </Button>
            <Button variant="primary" type="submit" disabled={createAbility.isPending}>
              {createAbility.isPending ? "Creating..." : "Create Ability"}
            </Button>
          </div>
        </form>
      </div>
    </PageContainer>
  );
}
