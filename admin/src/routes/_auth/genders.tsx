/* eslint-disable react-hooks/purity */

import { createFileRoute } from "@tanstack/react-router";
import { useState } from "react";
import { useGenders, useCreateGender, useUpdateGender, useDeleteGender, type Gender, type GenderInput } from "../../hooks/useGenders";
import { PageHeader } from "../../components/PageHeader";
import { DataTable, type Column } from "../../components/DataTable";
import { PageContainer } from "../../components/PageContainer";
import { Button } from "../../components/Button";
import { FormError } from "../../components/fields/FormError";

export const Route = createFileRoute("/_auth/genders")({
  component: GendersManagement,
});

function GendersManagement() {
  const { data: genders, isLoading, error } = useGenders();
  const createMutation = useCreateGender();
  const updateMutation = useUpdateGender();
  const deleteMutation = useDeleteGender();
  const [showForm, setShowForm] = useState(false);
  const [editingGender, setEditingGender] = useState<Gender | null>(null);
  const [formError, setFormError] = useState("");
  const [deletingGender, setDeletingGender] = useState<Gender | null>(null);

  const handleCreate = async (input: GenderInput) => {
    setFormError("");
    try { await createMutation.mutateAsync(input); setShowForm(false); }
    catch (e) { setFormError(e instanceof Error ? e.message : "Failed to create gender"); }
  };

  const handleUpdate = async (input: GenderInput) => {
    if (!editingGender) return;
    setFormError("");
    try { await updateMutation.mutateAsync({ id: editingGender.id, input }); setEditingGender(null); }
    catch (e) { setFormError(e instanceof Error ? e.message : "Failed to update gender"); }
  };

  const handleDelete = async () => {
    if (!deletingGender) return;
    try { await deleteMutation.mutateAsync(deletingGender.id); setDeletingGender(null); }
    catch { /* toasted by global handler */ }
  };

  const columns: Column<Gender>[] = [
    { header: "ID", accessor: "id", align: "center" },
    { header: "Name", accessor: "name" },
    { header: "Display Name", accessor: "display_name" },
    { header: "Subject Pronoun", accessor: "subject_pronoun" },
    { header: "Object Pronoun", accessor: "object_pronoun" },
    { header: "Possessive Pronoun", accessor: "possessive_pronoun" },
    { header: "World", accessor: "world_id", align: "center" },
    {
      header: "Actions",
      accessor: "_actions",
      render: (_: unknown, row: Gender) => (
        <span className="inline-flex gap-2">
          <Button variant="accent" size="sm" onClick={() => { setEditingGender(row); setShowForm(false); }}>Edit</Button>
          <Button variant="danger" size="sm" onClick={() => setDeletingGender(row)}>Delete</Button>
        </span>
      ),
    },
  ];

  const isSaving = createMutation.isPending || updateMutation.isPending;

  return (
    <PageContainer>
      <PageHeader title="Genders" backTo="/dashboard" actions={<Button variant="primary" onClick={() => { setShowForm(true); setEditingGender(null); }}>+ Add Gender</Button>} />
      {error && <div className="error-banner">{error instanceof Error ? error.message : "Failed to load genders"}</div>}
      {formError && <FormError message={formError} />}
      {showForm && !editingGender && <GenderForm gender={null} onSubmit={handleCreate} onCancel={() => setShowForm(false)} isLoading={isSaving} error={formError} />}
      {editingGender && <GenderForm gender={editingGender} onSubmit={handleUpdate} onCancel={() => setEditingGender(null)} isLoading={isSaving} error={formError} />}
      {isLoading ? <div className="loading">Loading genders...</div> : (
        <DataTable columns={columns} data={genders ?? []} getKey={(row: Gender) => row.id} emptyMessage="No genders found. Add your first gender!" />
      )}
      {deletingGender && <DeleteConfirmation gender={deletingGender} onConfirm={handleDelete} onCancel={() => setDeletingGender(null)} isLoading={deleteMutation.isPending} />}
    </PageContainer>
  );
}

function DeleteConfirmation({ gender, onConfirm, onCancel, isLoading }: Readonly<{ gender: Gender; onConfirm: () => void; onCancel: () => void; isLoading: boolean }>) {
  return (
    <div className="modal-overlay" onClick={onCancel}>
      <div className="modal-content modal-sm" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header"><h3>Delete Gender</h3><Button variant="ghost" size="sm" onClick={onCancel} aria-label="Close">×</Button></div>
        <div className="modal-body"><p>Are you sure you want to delete <strong>{gender.display_name || gender.name}</strong>?</p><p className="text-muted">This action cannot be undone.</p></div>
        <div className="modal-footer">
          <Button variant="danger" onClick={onConfirm} disabled={isLoading}>{isLoading ? "Deleting..." : "Delete"}</Button>
          <Button variant="secondary" onClick={onCancel}>Cancel</Button>
        </div>
      </div>
    </div>
  );
}

function GenderForm({ gender, onSubmit, onCancel, isLoading, error }: Readonly<{ gender: Gender | null; onSubmit: (input: GenderInput) => void; onCancel: () => void; isLoading: boolean; error: string }>) {
  const [form, setForm] = useState<GenderInput>({
    name: gender?.name || "",
    display_name: gender?.display_name || "",
    subject_pronoun: gender?.subject_pronoun || "",
    object_pronoun: gender?.object_pronoun || "",
    possessive_pronoun: gender?.possessive_pronoun || "",
    world_id: gender?.world_id || "1",
  });

  return (
    <div className="mt-4 p-4 bg-surface border border-border rounded">
      <h3 className="text-lg font-semibold mb-4">{gender ? "Edit Gender" : "Create New Gender"}</h3>
      {error && <div className="mb-4 p-2 bg-danger/10 border border-danger rounded text-danger">{error}</div>}
      <div className="space-y-3">
        <div className="flex gap-2">
          <input
            type="text"
            placeholder="Name (e.g., he_him)"
            value={form.name}
            onChange={(e) => setForm({ ...form, name: e.target.value })}
            className="flex-1 p-2 bg-surface border border-border rounded text-text text-sm"
            required
          />
          <input
            type="text"
            placeholder="Display Name (e.g., He/Him)"
            value={form.display_name}
            onChange={(e) => setForm({ ...form, display_name: e.target.value })}
            className="flex-1 p-2 bg-surface border border-border rounded text-text text-sm"
          />
        </div>
        <div className="flex gap-2">
          <input
            type="text"
            placeholder="Subject Pronoun (e.g., he)"
            value={form.subject_pronoun}
            onChange={(e) => setForm({ ...form, subject_pronoun: e.target.value })}
            className="flex-1 p-2 bg-surface border border-border rounded text-text text-sm"
          />
          <input
            type="text"
            placeholder="Object Pronoun (e.g., him)"
            value={form.object_pronoun}
            onChange={(e) => setForm({ ...form, object_pronoun: e.target.value })}
            className="flex-1 p-2 bg-surface border border-border rounded text-text text-sm"
          />
          <input
            type="text"
            placeholder="Possessive Pronoun (e.g., his)"
            value={form.possessive_pronoun}
            onChange={(e) => setForm({ ...form, possessive_pronoun: e.target.value })}
            className="flex-1 p-2 bg-surface border border-border rounded text-text text-sm"
          />
        </div>
        <div className="flex gap-2">
          <input
            type="text"
            placeholder="World ID (default: 1)"
            value={form.world_id}
            onChange={(e) => setForm({ ...form, world_id: e.target.value })}
            className="w-24 p-2 bg-surface border border-border rounded text-text text-sm"
          />
        </div>
        <div className="flex gap-2 pt-2">
          <Button
            variant="primary"
            onClick={() => onSubmit(form)}
            disabled={isLoading}
          >
            {isLoading ? "Saving..." : gender ? "Save Changes" : "Create Gender"}
          </Button>
          <Button variant="secondary" onClick={onCancel}>
            Cancel
          </Button>
        </div>
      </div>
    </div>
  );
}
