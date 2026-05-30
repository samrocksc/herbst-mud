/* eslint-disable react-refresh/only-export-components, functional/no-mixed-types */
import { useState } from "react";
import type { Race, RaceInput } from "../hooks/useRaces";
import { TagInput } from "./TagInput";
import { Button } from "./Button";
import { SLOT_CATALOG, DEFAULT_HUMANOID_SLOTS } from "./equipConstants";
import { FormField, TextareaField, ColorField, FormError } from "./fields";

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
  skill_grants: [],
  equipment_slots: [...DEFAULT_HUMANOID_SLOTS] as unknown as ReadonlyArray<string>,
  requirement_tags: [],
  color: "",
  tags: [],
} as const;

function raceToForm(r: Race): RaceInput {
  return {
    name: r.name,
    display_name: r.display_name,
    description: r.description ?? "",
    stat_modifiers: r.stat_modifiers ? JSON.stringify(r.stat_modifiers, null, 2) : "",
    skill_grants: r.skill_grants ? [...r.skill_grants] as unknown as ReadonlyArray<string> : [],
    equipment_slots: r.equipment_slots ? [...r.equipment_slots] as unknown as ReadonlyArray<string> : [...DEFAULT_HUMANOID_SLOTS] as unknown as ReadonlyArray<string>,
    requirement_tags: r.requirement_tags ? [...r.requirement_tags] as unknown as ReadonlyArray<string> : [],
    color: r.color ?? "",
    tags: r.tags ? [...r.tags] as unknown as ReadonlyArray<string> : [],
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
        <TagInput label="Skill Grants" value={form.skill_grants} onChange={(slots) => set("skill_grants", slots)}
          availableTags={[]} placeholder="Add skill..." tooltip="Skills automatically granted to characters of this race" />
        <TagInput label="Equipment Slots" value={form.equipment_slots} onChange={(slots) => set("equipment_slots", slots)}
          availableTags={[...SLOT_CATALOG]} placeholder="Add slot..." tooltip="Slots this race can equip items into" />
        <TagInput label="Requirement Tags" value={form.requirement_tags} onChange={(tags) => set("requirement_tags", tags)}
          availableTags={[]} placeholder="Add requirement..." tooltip="Tags that must be satisfied for race to be selectable (empty = playable)" />
        <TagInput label="Race Tags" value={form.tags} onChange={(tags) => set("tags", tags)}
          availableTags={availableTags} placeholder="Add tag..." tooltip="Tags automatically granted to characters of this race" />
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
