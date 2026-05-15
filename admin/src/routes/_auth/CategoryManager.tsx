import { useState } from 'react';
import { Button } from '../../components/Button';
import { DataTable, type Column } from '../../components/DataTable';
import { FormField, TextareaField } from '../../components/FormFields';
import { FormError } from '../../components/fields/FormError';
import { showToast } from '../../components/Toast';
import { apiPost } from '../../utils/apiFetch';
import type { FactionCategory } from './factionTypes';

export function CategoryManager({ categories }: Readonly<{ categories: FactionCategory[] }>) {
  return (
    <div className="max-w-[600px] mx-auto">
      <h2 className="mt-0 mb-4 text-text">Faction Categories</h2>
      <DataTable columns={categoryColumns} data={categories} getKey={(c) => c.id} emptyMessage="No categories yet" />
    </div>
  );
}

const categoryColumns: Column<FactionCategory>[] = [
  { header: 'ID', accessor: 'id' },
  { header: 'Name', accessor: 'name', className: 'font-bold' },
  { header: 'Display Name', accessor: 'display_name' },
  { header: 'Description', accessor: 'description' },
];

export function CreateCategoryForm({ onDone }: Readonly<{ onDone: () => void }>) {
  const [name, setName] = useState('');
  const [desc, setDesc] = useState('');
  const [error, setError] = useState('');
  const [saving, setSaving] = useState(false);

  const handleCreate = async () => {
    if (!name) { setError('Category name is required'); return; }
    setSaving(true);
    setError('');
    try {
      await apiPost('/api/faction-categories', { name, display_name: name, description: desc });
      showToast('Category created', 'success');
      onDone();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create category');
      showToast('Failed to create category', 'error');
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="max-w-[500px] mx-auto">
      <h2 className="mt-0 mb-4 text-text">Create Category</h2>
      {error && <FormError message={error} />}
      <div className="bg-surface-muted rounded-lg p-4 border border-border space-y-3">
        <FormField label="Name" value={name} onChange={setName} placeholder="Category name" />
        <TextareaField label="Description" value={desc} onChange={setDesc} rows={2} />
        <div className="flex gap-2">
          <Button variant="primary" size="md" fullWidth onClick={handleCreate} disabled={saving}>
            {saving ? 'Creating...' : 'Create Category'}
          </Button>
          <Button variant="secondary" size="md" fullWidth onClick={onDone}>Cancel</Button>
        </div>
      </div>
    </div>
  );
}