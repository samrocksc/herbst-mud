/* eslint-disable react-refresh/only-export-components, functional/no-mixed-types */
import { useState } from "react";
import type { Race, RaceInput } from "../hooks/useRaces";
import { TagInput } from "./TagInput";
import { Button } from "./Button";
import { SLOT_CATALOG, DEFAULT_HUMANOID_SLOTS } from "./equipConstants";
import { FormField, TextareaField, ColorField, FormError } from "./fields";
import { KeyValueEditor } from "./KeyValueEditor";

type RaceFormProps = Readonly<{
  race: Race | null
  onSubmit: (data: RaceInput) => void
  onCancel: () => void
  isLoading: boolean
  error?: string
  availableTags?: string[]
}>

const EMPTY_FORM: RaceInput = {
  name: "",
  display_name: "",
  description: "",
  stat_modifiers: "",
  stat_growth_multipliers: { hp: 1.0, mana: 1.0, stamina: 1.0 },
  skill_grants: [],
  equipment_slots: [...DEFAULT_HUMANOID_SLOTS] as unknown as ReadonlyArray<string>,
  requirement_tags: [],
  color: "",
  tags: [],
  resistances: {},
  vulnerabilities: {},
  world_id: "",
} as const;

function raceToForm(r: Race): RaceInput {
  const mult = r.stat_growth_multipliers;
  return {
    name: r.name,
    display_name: r.display_name,
    description: r.description ?? "",
    stat_modifiers: r.stat_modifiers ? JSON.stringify(r.stat_modifiers, null, 2) : "",
    stat_growth_multipliers: {
      hp: mult && typeof mult.hp === "number" ? mult.hp : 1.0,
      mana: mult && typeof mult.mana === "number" ? mult.mana : 1.0,
      stamina: mult && typeof mult.stamina === "number" ? mult.stamina : 1.0,
    },
    skill_grants: r.skill_grants ? [...r.skill_grants] as unknown as ReadonlyArray<string> : [],
    equipment_slots: r.equipment_slots ? [...r.equipment_slots] as unknown as ReadonlyArray<string> : [...DEFAULT_HUMANOID_SLOTS] as unknown as ReadonlyArray<string>,
    requirement_tags: r.requirement_tags ? [...r.requirement_tags] as unknown as ReadonlyArray<string> : [],
    color: r.color ?? "",
    tags: r.tags ? [...r.tags] as unknown as ReadonlyArray<string> : [],
    resistances: r.resistances ? { ...r.resistances } : {},
    vulnerabilities: r.vulnerabilities ? { ...r.vulnerabilities } : {},
  } as const;
}

export { EMPTY_FORM };
export { raceToForm };

export function RaceForm({ race, onSubmit, onCancel, isLoading, error, availableTags }: RaceFormProps) {
  const [form, setForm] = useState<RaceInput>(() => race ? raceToForm(race) : { ...EMPTY_FORM });
  const set = <K extends keyof RaceInput>(key: K, value: RaceInput[K]) =>
    setForm(prev => ({ ...prev, [key]: value }));

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!form.name.trim()) return;
    onSubmit(form);
  };

  return (
    <div className="form-card space-y-3">
      <h3 className="mt-0 mb-0 text-text text-base font-semibold">
        {race ? "Edit Race" : "Add New Race"}
      </h3>
      {error && <FormError message={error} />}
      <form onSubmit={handleSubmit} className="space-y-3">
        <FormField label="Name" value={form.name} onChange={(v) => set("name", v)} placeholder="e.g. elf, dwarf, human" required />
        <FormField label="Display Name" value={form.display_name} onChange={(v) => set("display_name", v)} placeholder="Defaults to name if blank" />
        <TextareaField label="Description" value={form.description} onChange={(v) => set("description", v)} rows={3} />
        <TextareaField label="Stat Modifiers (JSON)" value={form.stat_modifiers} onChange={(v) => set("stat_modifiers", v)} rows={4} placeholder='e.g. {"str": 2, "dex": -1}' />
        <div>
          <h4 className="text-sm font-semibold text-text mb-1">Stat Growth Multipliers</h4>
          <p className="text-xs text-text-muted mb-2">
            Racial multipliers applied to base stat growth per level. 1.0 = normal, 1.1 = +10%, 0.9 = -10%.
          </p>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <FloatField
              label="HP Multiplier"
              value={form.stat_growth_multipliers.hp}
              onChange={(v) => set("stat_growth_multipliers", { ...form.stat_growth_multipliers, hp: v })}
              tooltip="Multiplier for HP growth per level. 1.0 = default."
            />
            <FloatField
              label="Mana Multiplier"
              value={form.stat_growth_multipliers.mana}
              onChange={(v) => set("stat_growth_multipliers", { ...form.stat_growth_multipliers, mana: v })}
              tooltip="Multiplier for mana growth per level. 1.0 = default."
            />
            <FloatField
              label="Stamina Multiplier"
              value={form.stat_growth_multipliers.stamina}
              onChange={(v) => set("stat_growth_multipliers", { ...form.stat_growth_multipliers, stamina: v })}
              tooltip="Multiplier for stamina growth per level. 1.0 = default."
            />
          </div>
        </div>
        <TagInput label="Skill Grants" value={form.skill_grants} onChange={(slots) => set("skill_grants", slots)}
          availableTags={[]} placeholder="Add skill..." tooltip="Skills automatically granted to characters of this race" />
        <TagInput label="Equipment Slots" value={form.equipment_slots} onChange={(slots) => set("equipment_slots", slots)}
          availableTags={[...SLOT_CATALOG]} placeholder="Add slot..." tooltip="Slots this race can equip items into" />
        <TagInput label="Requirement Tags" value={form.requirement_tags} onChange={(tags) => set("requirement_tags", tags)}
          availableTags={[]} placeholder="Add requirement..." tooltip="Tags that must be satisfied for race to be selectable (empty = playable)" />
        <TagInput label="Race Tags" value={form.tags} onChange={(tags) => set("tags", tags)}
          availableTags={availableTags} placeholder="Add tag..." tooltip="Tags automatically granted to characters of this race" />
        <div className="pt-2 border-t border-border">
          <h4 className="text-sm font-semibold text-text mb-1">Resistances</h4>
          <p className="text-xs text-text-muted mb-2">
            Damage type → percentage. Positive values reduce damage taken (e.g. fire: 25 = 25% less fire damage).
          </p>
          <KeyValueEditor
            label="Resistances"
            value={{ ...form.resistances }}
            onChange={(v) => set("resistances", v)}
            keyPlaceholder="e.g. fire"
            tooltip="Map of damage type to resistance percentage. Positive = less damage taken."
          />
        </div>
        <div className="pt-2 border-t border-border">
          <h4 className="text-sm font-semibold text-text mb-1">Vulnerabilities</h4>
          <p className="text-xs text-text-muted mb-2">
            Damage type → percentage. Positive values increase damage taken (e.g. lightning: 25 = 25% more lightning damage).
          </p>
          <KeyValueEditor
            label="Vulnerabilities"
            value={{ ...form.vulnerabilities }}
            onChange={(v) => set("vulnerabilities", v)}
            keyPlaceholder="e.g. lightning"
            tooltip="Map of damage type to vulnerability percentage. Positive = more damage taken."
          />
        </div>
        <ColorField label="Color" value={form.color} onChange={(v) => set("color", v)} placeholder="e.g. #8b5cf6" />
        <div className="flex gap-2 pt-1">
          <Button type="submit" variant="primary" disabled={isLoading || !form.name.trim()} fullWidth>
            {isLoading ? "Saving..." : race ? "Update Race" : "Create Race"}
          </Button>
          <Button type="button" variant="secondary" onClick={onCancel} fullWidth>Cancel</Button>
        </div>
      </form>
    </div>
  );
}

/**
 * Float input field — like NumberField but for float values.
 * Used for stat growth multipliers (e.g. 1.1, 0.9).
 */
function FloatField({ label, value, onChange, tooltip }: Readonly<{ label: string; value: number; onChange: (v: number) => void; tooltip?: string }>) {
  const fieldId = label.toLowerCase().replace(/\s+/g, "-");
  return (
    <div>
      <label htmlFor={fieldId} className="text-text-muted text-xs block mb-1" title={tooltip}>
        {label}
      </label>
      <input
        id={fieldId}
        type="text"
        inputMode="decimal"
        value={Number.isFinite(value) ? String(value) : ""}
        onChange={(e) => {
          const raw = e.target.value;
          if (raw === "" || raw === "-") { onChange(0); return; }
          const parsed = parseFloat(raw);
          if (!isNaN(parsed)) onChange(parsed);
        }}
        className="w-full px-3 py-2 bg-surface border border-border rounded text-sm text-text focus:outline-none focus:border-primary"
      />
    </div>
  );
}
