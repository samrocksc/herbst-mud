import { useState } from 'react'
import type { Race, RaceInput } from '../hooks/useRaces'
import { TagInput } from './TagInput'
import { Button } from './Button'
import { SLOT_CATALOG, DEFAULT_HUMANOID_SLOTS } from './equipConstants'

type RaceFormProps = Readonly<{
  race: Race | null
  onSubmit: (data: RaceInput) => void
  onCancel: () => void
  isLoading: boolean
}>

const EMPTY_FORM: RaceInput = {
  name: '',
  display_name: '',
  description: '',
  stat_modifiers: '',
  equipment_slots: [...DEFAULT_HUMANOID_SLOTS],
  is_playable: true,
  color: '',
}

function raceToForm(r: Race): RaceInput {
  return {
    name: r.name,
    display_name: r.display_name,
    description: r.description ?? '',
    stat_modifiers: r.stat_modifiers ? JSON.stringify(r.stat_modifiers, null, 2) : '',
    equipment_slots: r.equipment_slots ? [...r.equipment_slots] : [...DEFAULT_HUMANOID_SLOTS],
    is_playable: r.is_playable,
    color: r.color ?? '',
  }
}

export { EMPTY_FORM, raceToForm }

export function RaceForm({ race, onSubmit, onCancel, isLoading }: RaceFormProps) {
  const [form, setForm] = useState<RaceInput>(() => race ? raceToForm(race) : { ...EMPTY_FORM })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!form.name.trim()) return
    onSubmit(form)
  }

  const set = <K extends keyof RaceInput>(key: K, value: RaceInput[K]) => {
    setForm(prev => ({ ...prev, [key]: value }))
  }

  return (
    <div className="form-card">
      <h3>{race ? 'Edit Race' : 'Add New Race'}</h3>
      <form onSubmit={handleSubmit}>
        <div className="form-row"><label>Name:</label>
          <input type="text" value={form.name} onChange={(e) => set('name', e.target.value)} placeholder="e.g. elf, dwarf, human" required /></div>
        <div className="form-row"><label>Display Name:</label>
          <input type="text" value={form.display_name} onChange={(e) => set('display_name', e.target.value)} placeholder="Defaults to name if blank" /></div>
        <div className="form-row"><label>Description:</label>
          <textarea value={form.description} onChange={(e) => set('description', e.target.value)} rows={3} /></div>
        <div className="form-row"><label>Stat Modifiers (JSON):</label>
          <textarea value={form.stat_modifiers} onChange={(e) => set('stat_modifiers', e.target.value)} rows={4} placeholder='e.g. {"str": 2, "dex": -1}' /></div>
        <TagInput label="Equipment Slots" value={form.equipment_slots} onChange={(slots) => set('equipment_slots', slots)}
          availableTags={[...SLOT_CATALOG]} placeholder="Add slot..." tooltip="Slots this race can equip items into" />
        <div className="form-row"><label>Playable:</label>
          <input type="checkbox" checked={form.is_playable} onChange={(e) => set('is_playable', e.target.checked)} /></div>
        <div className="form-row"><label>Color:</label>
          <input type="text" value={form.color} onChange={(e) => set('color', e.target.value)} placeholder="e.g. #8b5cf6" /></div>
        <div className="form-actions">
          <Button type="submit" variant="primary" disabled={isLoading || !form.name.trim()}>
            {isLoading ? 'Saving...' : race ? 'Update Race' : 'Create Race'}
          </Button>
          <Button type="button" variant="secondary" onClick={onCancel}>Cancel</Button>
        </div>
      </form>
    </div>
  )
}