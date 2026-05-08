import {
  FormField,
  NumberField,
  TextareaField,
  SelectField,
  CheckboxField,
} from '../../components/FormFields'
import type { FactionCategory, FactionForm } from './factionTypes'

export function FactionFormFields({
  form,
  setForm,
  categories,
}: Readonly<{
  form: FactionForm
  setForm: (f: FactionForm) => void
  categories: FactionCategory[]
}>) {
  const set = (patch: Partial<FactionForm>) => setForm({ ...form, ...patch })
  const catOptions = [
    { value: '', label: 'None' },
    ...categories.map((c) => ({ value: String(c.id), label: c.display_name || c.name })),
  ]

  return (
    <div className="bg-surface-muted rounded-lg p-4 border border-border space-y-3">
      <FormField
        label="Name"
        value={form.name}
        onChange={(v) => set({ name: v })}
        placeholder="Faction name"
      />
      <FormField
        label="Display Name"
        value={form.display_name}
        onChange={(v) => set({ display_name: v })}
        placeholder="Shown to players (defaults to name)"
      />
      <TextareaField
        label="Description"
        value={form.description}
        onChange={(v) => set({ description: v })}
        rows={3}
      />
      <div className="grid grid-cols-2 gap-3">
        <SelectField
          label="Category"
          value={String(form.category_id)}
          onChange={(v) => set({ category_id: v ? Number(v) : '' })}
          options={catOptions}
        />
        <NumberField
          label="Standing"
          value={form.standing}
          onChange={(v) => set({ standing: v })}
        />
      </div>
      <CheckboxField
        label="Universal (applies to all characters)"
        checked={form.is_universal}
        onChange={(v) => set({ is_universal: v })}
      />
    </div>
  )
}