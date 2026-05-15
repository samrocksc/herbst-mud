import { FormError } from '../../components/fields/FormError';
import { FactionFormFields } from './FactionFormFields';
import type { FactionCategory, FactionForm } from './factionTypes';

export function CreateFactionView({ form, setForm, categories, createError, saving, onCreate, onCancel }: Readonly<{
  form: FactionForm; setForm: (f: FactionForm) => void; categories: FactionCategory[]
  createError: string; saving: boolean; onCreate: () => void; onCancel: () => void
}>) {
  return (
    <div className="max-w-[600px] mx-auto">
      <h2 className="mt-0 mb-4 text-text">Create Faction</h2>
      {createError && <FormError message={createError} />}
      <FactionFormFields form={form} setForm={setForm} categories={categories} />
      <div className="flex gap-2 mt-3">
        <button onClick={onCreate} disabled={saving}
          className="flex-1 py-2 px-4 rounded bg-primary text-white font-medium disabled:opacity-50">
          {saving ? 'Creating...' : 'Create Faction'}
        </button>
        <button onClick={onCancel}
          className="flex-1 py-2 px-4 rounded bg-surface-muted text-text border border-border font-medium">
          Cancel
        </button>
      </div>
    </div>
  );
}