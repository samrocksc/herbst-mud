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
  { value: "strength", label: "Strength" },
  { value: "dexterity", label: "Dexterity" },
  { value: "constitution", label: "Constitution" },
  { value: "intelligence", label: "Intelligence" },
  { value: "wisdom", label: "Wisdom" },
];

const EMPTY_ABILITY: AbilityInput = {
  name: "",
  description: "",
  ability_type: "combat",
  requirements: "1",
  cost: 0,
  cooldown: 0,
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

  const handleEffectSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const newEffect: AbilityEffect = {
      ...currentEffect,
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
      // First create the ability
      const ability = await createAbility.mutateAsync(formData);

      // Then create each effect
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

  return (
    <PageContainer>
      <PageHeader title="Create Ability" showBack backTo="/abilities" />
      <div className="card bg-surface p-6 border border-border rounded">
        <form onSubmit={handleSubmit} className="space-y-4">
          <h3 className="text-text font-semibold mb-4">Basic Information</h3>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <FormField label="Name *" value={formData.name} onChange={(v) => set({ name: v })} required />
            <FormField label="Slug (optional)" value={formData.slug ?? ""} onChange={(v) => set({ slug: v })} placeholder="Auto-generated from name if empty" />
            <SelectField label="Ability Type" value={formData.ability_type} onChange={(v) => set({ ability_type: v })} options={ABILITY_TYPE_OPTS} />
            <SearchableSelect
              label="Required Tag (optional)"
              options={(availableTags ?? []).map((t) => ({ id: t.name, name: t.name }))}
              value={formData.required_tag || ""}
              onChange={(v) => set({ required_tag: v })}
              placeholder="Select a tag..."
            />
            <SelectField label="Ability Class" value={formData.ability_class} onChange={(v) => set({ ability_class: v })} options={ABILITY_CLASS_OPTS} />
          </div>
          <TextareaField label="Description" value={formData.description} onChange={(v) => set({ description: v })} rows={3} />

          <h3 className="text-text font-semibold mt-6 mb-4">Costs & Cooldown</h3>
          <div className="grid grid-cols-3 gap-4">
            <NumberField label="Cost" value={formData.cost} onChange={(v) => set({ cost: v })} />
            <NumberField label="Cooldown (s)" value={formData.cooldown_seconds} onChange={(v) => set({ cooldown_seconds: v })} />
            <FormField label="Unlock Tags (JSON)" value={formData.requirements} onChange={(v) => set({ requirements: v })} placeholder='{"tags":["level:5"]}' />
          </div>
          <div className="grid grid-cols-3 gap-4">
            <NumberField label="Mana Cost" value={formData.mana_cost} onChange={(v) => set({ mana_cost: v })} />
            <NumberField label="Stamina Cost" value={formData.stamina_cost} onChange={(v) => set({ stamina_cost: v })} />
            <NumberField label="HP Cost" value={formData.hp_cost} onChange={(v) => set({ hp_cost: v })} />
          </div>

          <h3 className="text-text font-semibold mt-6 mb-4">Proc Settings</h3>
          <div className="grid grid-cols-2 gap-4">
            <NumberField label="Proc Chance (0–1)" value={formData.proc_chance} onChange={(v) => set({ proc_chance: v })} step={0.01} />
            <FormField label="Proc Event" value={formData.proc_event} onChange={(v) => set({ proc_event: v })} placeholder="e.g., on_hit, on_crit" />
          </div>

          <h3 className="text-text font-semibold mt-6 mb-4">Combat Messages</h3>
          <div className="space-y-3">
            <div>
              <label className="text-text-muted text-xs block mb-1">Caster Message (optional)</label>
              <p className="text-xs text-muted mb-1">Template when actor uses the ability. Example: "{'{actor}'} casts fireball"</p>
              <FormField label="Caster Message" value={formData.caster_message} onChange={(v) => set({ caster_message: v })} placeholder={"e.g., {actor} casts fireball"} />
            </div>
            <div>
              <label className="text-text-muted text-xs block mb-1">Recipient Message (optional)</label>
              <p className="text-xs text-muted mb-1">Template when recipient is targeted. Example: "{'{actor}'} casts fireball at you"</p>
              <FormField label="Recipient Message" value={formData.recipient_message} onChange={(v) => set({ recipient_message: v })} placeholder={"e.g., {actor} casts fireball at you"} />
            </div>
          </div>

          <h3 className="text-text font-semibold mt-6 mb-4">Effects</h3>
          <p className="text-sm text-muted mb-4">Define what the ability does. Add multiple effects for complex abilities.</p>

          <form onSubmit={handleEffectSubmit} className="bg-surface/50 border border-border rounded p-4 space-y-4 mb-4">
            <h4 className="text-text font-semibold mb-2">Add Effect</h4>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <SelectField label="Effect Type" value={currentEffect.effect_type} onChange={(v) => setEffect({ effect_type: v })} options={EFFECT_TYPE_OPTS} />
              <SelectField label="Target" value={currentEffect.target} onChange={(v) => setEffect({ target: v })} options={TARGET_OPTS} />
            </div>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <SelectField label="Damage Type (if applicable)" value={currentEffect.damage_subtype} onChange={(v) => setEffect({ damage_subtype: v })} options={DAMAGE_SUBTYPE_OPTS} />
              <NumberField label="Value / Amount" value={currentEffect.value} onChange={(v) => setEffect({ value: v })} />
            </div>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <NumberField label="Duration (ticks)" value={currentEffect.duration} onChange={(v) => setEffect({ duration: v })} />
              <SelectField label="Scaling Stat (optional)" value={currentEffect.scaling_stat} onChange={(v) => setEffect({ scaling_stat: v })} options={SCALING_STAT_OPTS} />
            </div>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <NumberField label="Scaling Ratio" value={currentEffect.scaling_ratio} onChange={(v) => setEffect({ scaling_ratio: v })} step={0.1} />
            </div>
            <div className="flex justify-end">
              <Button variant="primary" type="submit" disabled={createAbilityEffect.isPending}>
                {createAbilityEffect.isPending ? "Adding..." : "Add Effect"}
              </Button>
            </div>
          </form>

          {effects.length > 0 && (
            <div className="bg-surface/50 border border-border rounded p-4">
              <h4 className="text-text font-semibold mb-3">Added Effects</h4>
              <div className="space-y-3">
                {effects.map((effect, index) => (
                  <div key={index} className="flex justify-between items-start bg-surface border border-border rounded p-3">
                    <div>
                      <div className="font-medium">{effect.effect_type}</div>
                      <div className="text-sm text-muted">
                        {effect.target} • {effect.damage_subtype} • Value: {effect.value}
                        {effect.duration > 0 && ` • Duration: ${effect.duration}`}
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

          <div className="flex gap-2 justify-end mt-6">
            <Button variant="secondary" onClick={() => navigate({ to: "/abilities" })}>Cancel</Button>
            <Button variant="primary" type="submit" disabled={createAbility.isPending}>
              {createAbility.isPending ? "Creating..." : "Create Ability"}
            </Button>
          </div>
        </form>
      </div>
    </PageContainer>
  );
}
