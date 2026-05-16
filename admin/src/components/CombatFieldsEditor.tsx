 
/* eslint-disable functional/no-mixed-types */

import { NumberField, SelectField, CheckboxField, TextareaField } from "./FormFields";
import { WEAPON_TYPES, DAMAGE_TYPES, ARMOR_TYPES, RARITY_OPTIONS, isWeaponSlot, isArmorSlot } from "./equipConstants";

/** Combat form state shared by Equipment and EquipmentTemplate. */
export type CombatFields = Readonly<{
  damage_dice_count: number
  damage_dice_sides: number
  damage_bonus: number
  weapon_type: string
  damage_type: string
  is_two_handed: boolean
  armor_rating: number
  armor_type: string
  rarity: string
  skill_requirement: string
  skill_requirement_level: number
  stats: string
}>

type CombatFieldsEditorProps = Readonly<{
  form: CombatFields
  onChange: (update: Partial<CombatFields>) => void
  slot: string
}>

/** Renders weapon/armor combat fields conditionally based on slot. */
export function CombatFieldsEditor({ form, onChange, slot }: CombatFieldsEditorProps) {
  const showWeapon = isWeaponSlot(slot);
  const showArmor = isArmorSlot(slot);

  return (
    <div className="space-y-4">
      <h3 className="text-text text-sm font-semibold border-b border-border pb-1">
        Combat Stats
      </h3>

      <SelectField label="Rarity" value={form.rarity} onChange={(v) => onChange({ rarity: v })}
        options={[...RARITY_OPTIONS]} />

      {showWeapon && <WeaponFields form={form} onChange={onChange} />}
      {showArmor && <ArmorFields form={form} onChange={onChange} />}

      <SharedFields form={form} onChange={onChange} />
    </div>
  );
}

function WeaponFields({ form, onChange }: Readonly<{ form: CombatFields; onChange: (u: Partial<CombatFields>) => void }>) {
  return (
    <>
      <div className="grid grid-cols-3 gap-4">
        <NumberField label="Damage Dice Count" value={form.damage_dice_count} onChange={(v) => onChange({ damage_dice_count: v })} min={0} />
        <NumberField label="Damage Dice Sides" value={form.damage_dice_sides} onChange={(v) => onChange({ damage_dice_sides: v })} min={0} />
        <NumberField label="Damage Bonus" value={form.damage_bonus} onChange={(v) => onChange({ damage_bonus: v })} />
      </div>
      <div className="grid grid-cols-2 gap-4">
        <SelectField label="Weapon Type" value={form.weapon_type} onChange={(v) => onChange({ weapon_type: v })} options={[...WEAPON_TYPES]} />
        <SelectField label="Damage Type" value={form.damage_type} onChange={(v) => onChange({ damage_type: v })} options={[...DAMAGE_TYPES]} />
      </div>
      <CheckboxField label="Two-Handed" checked={form.is_two_handed} onChange={(v) => onChange({ is_two_handed: v })} />
    </>
  );
}

function ArmorFields({ form, onChange }: Readonly<{ form: CombatFields; onChange: (u: Partial<CombatFields>) => void }>) {
  return (
    <div className="grid grid-cols-2 gap-4">
      <NumberField label="Armor Rating" value={form.armor_rating} onChange={(v) => onChange({ armor_rating: v })} min={0} />
      <SelectField label="Armor Type" value={form.armor_type} onChange={(v) => onChange({ armor_type: v })} options={[...ARMOR_TYPES]} />
    </div>
  );
}

function SharedFields({ form, onChange }: Readonly<{ form: CombatFields; onChange: (u: Partial<CombatFields>) => void }>) {
  return (
    <>
      <div className="grid grid-cols-2 gap-4">
        <div>
          <label className="text-text-muted text-xs block mb-1">Skill Requirement</label>
          <input type="text" value={form.skill_requirement} onChange={(e) => onChange({ skill_requirement: e.target.value })}
            placeholder="e.g. skill_blades" className="w-full p-2 bg-surface border border-border rounded text-text text-sm" />
        </div>
        <NumberField label="Skill Requirement Level" value={form.skill_requirement_level} onChange={(v) => onChange({ skill_requirement_level: v })} min={0} />
      </div>
      <TextareaField label="Stats (JSON)" value={form.stats} onChange={(v) => onChange({ stats: v })} rows={3}
        placeholder='e.g. {"strength": 2, "dexterity": 1}' />
    </>
  );
}