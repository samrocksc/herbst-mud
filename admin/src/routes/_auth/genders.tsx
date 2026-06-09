/* eslint-disable react-hooks/purity */

import { createFileRoute } from "@tanstack/react-router";
import { useState } from "react";
import { useGenders, useCreateGender, useUpdateGender, useDeleteGender, type Gender, type GenderInput } from "../../hooks/useGenders";
import { useWorldStore } from "../../contexts/WorldStoreContext";
import { PageHeader } from "../../components/PageHeader";
import { DataTable, type Column } from "../../components/DataTable";
import { PageContainer } from "../../components/PageContainer";
import { Button } from "../../components/Button";
import { FormField } from "../../components/FormFields";
import { FormError } from "../../components/fields/FormError";

export const Route = createFileRoute("/_auth/genders")({
  component: GendersManagement,
});

function GendersManagement() {
  const { currentWorld } = useWorldStore();
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
    try { await deleteMutation.mutateAsync(deletingGender); setDeletingGender(null); }
    catch { /* toasted by global handler */ }
  };

  const columns: Column<Gender>[] = [
    { header: "ID", accessor: "id", align: "center" },
    {
      header: "Name",
      accessor: "display_name",
      render: (_, row) => (
        <span>
          <span className="font-medium text-text">{row.display_name || row.name}</span>
          {row.display_name !== row.name && (
            <code className="ml-2 text-xs text-text-muted">{row.name}</code>
          )}
        </span>
      ),
    },
    { header: "Subject", accessor: "subject_pronoun" },
    { header: "Object", accessor: "object_pronoun" },
    { header: "Possessive", accessor: "possessive_pronoun" },
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
      <div className="mb-4 text-sm text-text-muted">
        Managing genders for world <span className="text-text font-medium">{currentWorld}</span>. Switch worlds via the dashboard dropdown.
      </div>
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

type GenderFormProps = Readonly<{
  gender: Gender | null
  onSubmit: (input: GenderInput) => void
  onCancel: () => void
  isLoading: boolean
  error: string
}>

function GenderForm({ gender, onSubmit, onCancel, isLoading, error }: GenderFormProps) {
  const { currentWorld } = useWorldStore();
  const [form, setForm] = useState<GenderInput>({
    name: gender?.name || "",
    display_name: gender?.display_name || "",
    subject_pronoun: gender?.subject_pronoun || "",
    object_pronoun: gender?.object_pronoun || "",
    possessive_pronoun: gender?.possessive_pronoun || "",
    world_id: currentWorld || "default",
  });

  const handleField = (field: keyof GenderInput) => (value: string) => {
    setForm(prev => ({ ...prev, [field]: value }));
  };

  return (
    <div className="mt-4 p-6 bg-surface border border-border rounded">
      <h3 className="text-lg font-semibold mb-4">{gender ? "Edit Gender" : "Create New Gender"}</h3>
      <p className="text-xs text-text-muted mb-4">Define how characters of this gender are referenced in messages. World is automatically set to <code className="text-text">{currentWorld}</code>.</p>
      {error && <div className="mb-4 p-2 bg-danger/10 border border-danger rounded text-danger">{error}</div>}
      <div className="space-y-3">
        {/* Name fields */}
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <FormField
            label="Name"
            value={form.name}
            onChange={handleField("name")}
            placeholder="he_him"
            required
            tooltip="Identifier used in code and content files (snake_case). Example: he_him, she_her, they_them."
          />
          <FormField
            label="Display Name"
            value={form.display_name}
            onChange={handleField("display_name")}
            placeholder="He/Him"
            tooltip="How the gender appears in admin UI and character creation. Example: He/Him, She/Her, They/Them."
          />
        </div>
        {/* Pronoun fields */}
        <p className="text-xs text-text-muted pt-1">Pronouns — used in game messages like "You see <em>them</em>." and "It belongs to <em>him</em>."</p>
        <div className="grid grid-cols-3 gap-3">
          <FormField
            label="Subject"
            value={form.subject_pronoun}
            onChange={handleField("subject_pronoun")}
            placeholder="he"
            tooltip='Used when the character is the subject. Example: "He swings a sword."'
          />
          <FormField
            label="Object"
            value={form.object_pronoun}
            onChange={handleField("object_pronoun")}
            placeholder="him"
            tooltip='Used when the character is the object. Example: "The goblin hits him."'
          />
          <FormField
            label="Possessive"
            value={form.possessive_pronoun}
            onChange={handleField("possessive_pronoun")}
            placeholder="his"
            tooltip='Used for ownership. Example: "his sword."'
          />
        </div>
        {/* Live preview */}
        {(form.subject_pronoun || form.object_pronoun || form.possessive_pronoun) && (
          <div className="bg-surface-muted border border-border rounded p-3">
            <p className="text-xs text-text-muted mb-1 font-medium uppercase tracking-wider">Preview</p>
            <div className="space-y-1 text-sm text-text">
              {form.subject_pronoun && (
                <p><span className="text-primary font-medium">{form.subject_pronoun}</span> enters the room. <span className="text-primary font-medium">{form.subject_pronoun}</span> looks around.</p>
              )}
              {form.object_pronoun && (
                <p>The goblin attacks <span className="text-primary font-medium">{form.object_pronoun}</span>! You see <span className="text-primary font-medium">{form.object_pronoun}</span>.</p>
              )}
              {form.possessive_pronoun && (
                <p>This is <span className="text-primary font-medium">{form.possessive_pronoun}</span> sword. It belongs to <span className="text-primary font-medium">{form.object_pronoun || "them"}</span>.</p>
              )}
            </div>
          </div>
        )}
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
