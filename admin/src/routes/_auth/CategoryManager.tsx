import { useState } from "react";
import { Button } from "../../components/Button";
import { DataTable, type Column } from "../../components/DataTable";
import { FormField, TextareaField } from "../../components/FormFields";
import { FormError } from "../../components/fields/FormError";
import { showToast } from "../../components/Toast";
import { apiPost, apiPut, apiDelete } from "../../utils/apiFetch";
import { DeleteConfirmation } from "../../components/DeleteConfirmation";
import type { FactionCategory } from "./-factionTypes";

export function CategoryManager({ categories, onRefresh, onEdit }: Readonly<{ categories: FactionCategory[]; onRefresh: () => void; onEdit?: (cat: FactionCategory) => void }>) {
  const [deleteId, setDeleteId] = useState<number | null>(null);
  const [editingCategory, setEditingCategory] = useState<FactionCategory | null>(null);
  const [deleting, setDeleting] = useState(false);

  const handleDelete = async () => {
    if (deleteId == null) return;
    setDeleting(true);
    try {
      await apiDelete(`/api/faction-categories/${deleteId}`);
      showToast("Category deleted", "success");
      setDeleteId(null);
      onRefresh();
    } catch (err) {
      showToast(err instanceof Error ? err.message : "Failed to delete category", "error");
    } finally {
      setDeleting(false);
    }
  };

  const handleEdit = (cat: FactionCategory) => {
    if (onEdit) onEdit(cat);
    else setEditingCategory(cat);
  };

  return (
    <div className="max-w-[600px] mx-auto">
      <h2 className="mt-0 mb-4 text-text">Faction Categories</h2>
      {editingCategory != null && !onEdit && (
        <CategoryEditForm category={editingCategory} onDone={() => setEditingCategory(null)} />
      )}
      <DataTable columns={categoryColumns(setDeleteId, handleEdit)} data={categories} getKey={(c) => c.id} emptyMessage="No categories yet" />
      {deleteId != null && (
        <DeleteConfirmation
          open={deleteId != null}
          title="Delete Faction Category"
          message="Are you sure? This will unlink any factions in this category."
          onConfirm={handleDelete}
          onCancel={() => setDeleteId(null)}
          isLoading={deleting}
        />
      )}
    </div>
  );
}

export function CategoryEditForm({ category, onDone }: Readonly<{ category: FactionCategory; onDone: () => void }>) {
  const [name, setName] = useState(category.display_name || category.name);
  const [desc, setDesc] = useState(category.description || "");
  const [maxMemberships, setMaxMemberships] = useState(category.max_memberships ?? 1);
  const [autoJoin, setAutoJoin] = useState(category.auto_join ?? false);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");

  const handleSave = async () => {
    if (!name) { setError("Name is required"); return; }
    setSaving(true);
    setError("");
    try {
      await apiPut(`/api/faction-categories/${category.id}`, {
        name: category.name,
        display_name: name,
        description: desc,
        max_memberships: maxMemberships,
        auto_join: autoJoin,
      });
      showToast("Category updated", "success");
      onDone();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to update");
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="bg-surface-muted rounded-lg p-4 border border-border space-y-3 mb-4">
      <h4 className="mt-0 mb-2 text-text text-sm font-semibold">Edit: {category.name}</h4>
      {error && <FormError message={error} />}
      <FormField label="Display Name" value={name} onChange={setName} />
      <TextareaField label="Description" value={desc} onChange={setDesc} rows={2} />
      <FormField label="Max Memberships" value={String(maxMemberships)} onChange={(v) => setMaxMemberships(Number(v) || 1)} type="number" />
      <label className="flex items-center gap-2 text-sm text-text cursor-pointer">
        <input type="checkbox" checked={autoJoin} onChange={(e) => setAutoJoin(e.target.checked)} className="accent-primary" />
        Auto Join
      </label>
      <div className="flex gap-2">
        <Button variant="primary" size="sm" onClick={handleSave} disabled={saving}>{saving ? "Saving..." : "Save"}</Button>
        <Button variant="secondary" size="sm" onClick={onDone}>Cancel</Button>
      </div>
    </div>
  );
}

function categoryColumns(setDeleteId: (id: number | null) => void, onEdit: (cat: FactionCategory) => void): Column<FactionCategory>[] {
  return [
    { header: "ID", accessor: "id" },
    { header: "Name", accessor: "name", className: "font-bold" },
    { header: "Display Name", accessor: "display_name" },
    { header: "Memberships", accessor: "max_memberships", render: (v) => String(v ?? 1) },
    { header: "Description", accessor: "description" },
    {
      header: "",
      accessor: "_actions",
      align: "right",
      render: (_, row) => (
        <div className="flex gap-2 justify-end">
          <Button variant="ghost" size="sm" onClick={(e) => { e.stopPropagation(); onEdit(row); }}>Edit</Button>
          <Button variant="danger" size="sm" onClick={(e) => { e.stopPropagation(); setDeleteId(row.id); }}>Delete</Button>
        </div>
      ),
    },
  ];
}

export function CreateCategoryForm({ onDone }: Readonly<{ onDone: () => void }>) {
  const [name, setName] = useState("");
  const [desc, setDesc] = useState("");
  const [maxMemberships, setMaxMemberships] = useState(1);
  const [autoJoin, setAutoJoin] = useState(false);
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