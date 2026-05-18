import { useState } from "react";
import { Button } from "../../components/Button";
import { DataTable, type Column } from "../../components/DataTable";
import { FormField, TextareaField } from "../../components/FormFields";
import { FormError } from "../../components/fields/FormError";
import { showToast } from "../../components/Toast";
import { apiPost } from "../../utils/apiFetch";
import type { FactionCategory } from "./factionTypes";

export function CategoryManager({ categories }: Readonly<{ categories: FactionCategory[] }>) {
  return (
    <div className="max-w-[600px] mx-auto">
      <h2 className="mt-0 mb-4 text-text">Faction Categories</h2>
      <DataTable columns={categoryColumns} data={categories} getKey={(c) => c.id} emptyMessage="No categories yet" />
    </div>
  );
}

const categoryColumns: Column<FactionCategory>[] = [
  { header: "ID", accessor: "id" },
  { header: "Name", accessor: "name", className: "font-bold" },
  { header: "Display Name", accessor: "display_name" },
  { header: "Memberships", accessor: "max_memberships", render: (v) => String(v ?? 1) },
  { header: "Wizard", accessor: "initial_config", render: (v) => v ? "Yes" : "No" },
  { header: "Description", accessor: "description" },
];

export function CreateCategoryForm({ onDone }: Readonly<{ onDone: () => void }>) {
  const [name, setName] = useState("");
  const [desc, setDesc] = useState("");
  const [maxMemberships, setMaxMemberships] = useState(1);
  const [autoJoin, setAutoJoin] = useState(false);
  const [initialConfig, setInitialConfig] = useState(false);
  const [error, setError] = useState("");
  const [saving, setSaving] = useState(false);

  const handleCreate = async () => {
    if (!name) { setError("Category name is required"); return; }
    setSaving(true);
    setError("");
    try {
      await apiPost("/api/faction-categories", {
        name,
        display_name: name,
        description: desc,
        max_memberships: maxMemberships,
        auto_join: autoJoin,
        initial_config: initialConfig,
      });
      showToast("Category created", "success");
      onDone();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to create category");
      showToast("Failed to create category", "error");
    } finally {
      setSaving(false);
    }
  };

  const initialConfigDisabled = maxMemberships > 1;

  return (
    <div className="max-w-[500px] mx-auto">
      <h2 className="mt-0 mb-4 text-text">Create Category</h2>
      {error && <FormError message={error} />}
      <div className="bg-surface-muted rounded-lg p-4 border border-border space-y-3">
        <FormField label="Name" value={name} onChange={setName} placeholder="Category name" />
        <TextareaField label="Description" value={desc} onChange={setDesc} rows={2} />
        <FormField label="Max Memberships" value={String(maxMemberships)} onChange={(v) => setMaxMemberships(Number(v) || 1)} type="number" />
        <label className="flex items-center gap-2 text-sm text-text cursor-pointer">
          <input
            type="checkbox"
            checked={autoJoin}
            onChange={(e) => setAutoJoin(e.target.checked)}
            className="accent-primary"
          />
          Auto Join (earning required tag auto-joins faction)
        </label>
        <label className="flex items-center gap-2 text-sm text-text cursor-pointer">
          <input
            type="checkbox"
            checked={initialConfig}
            onChange={(e) => setInitialConfig(e.target.checked)}
            disabled={initialConfigDisabled}
            className="accent-primary"
          />
          Initial Config (appears in character creation wizard)
          {initialConfigDisabled && (
            <span className="text-text-muted text-xs italic">— Only available for max_memberships = 1</span>
          )}
        </label>
        <div className="flex gap-2">
          <Button variant="primary" size="md" fullWidth onClick={handleCreate} disabled={saving}>
            {saving ? "Creating..." : "Create Category"}
          </Button>
          <Button variant="secondary" size="md" fullWidth onClick={onDone}>Cancel</Button>
        </div>
      </div>
    </div>
  );
}