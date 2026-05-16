import { createFileRoute } from "@tanstack/react-router";
import { useState } from "react";
import { useRaces, useCreateRace, useUpdateRace, useDeleteRace, useApplyRaceTags, type Race, type RaceInput } from "../../hooks/useRaces";
import { useTags } from "../../hooks/useTags";
import { PageHeader } from "../../components/PageHeader";
import { DataTable, type Column } from "../../components/DataTable";
import { Button } from "../../components/Button";
import { RaceForm } from "../../components/RaceForm";
import { FormError } from "../../components/fields/FormError";

export const Route = createFileRoute("/_auth/races")({
  component: RacesManagement,
});

const COLUMNS: Column<Race>[] = [
  { header: "Name", accessor: "name", render: (_: unknown, row: Race) => <strong>{row.display_name || row.name}</strong> },
  { header: "Tags", accessor: "tags", render: (_: unknown, row: Race) => (
    <span className="flex flex-wrap gap-1">
      {(row.tags ?? []).map(t => <span key={t} className="badge badge-neutral">{t}</span>)}
    </span>
  )},
  { header: "Slots", accessor: "equipment_slots", render: (_: unknown, row: Race) => <span className="text-xs">{(row.equipment_slots ?? []).join(", ") || "—"}</span> },
  { header: "Playable", accessor: "is_playable" },
];

function RacesManagement() {
  const { data: races, isLoading, error } = useRaces();
  const { data: tags } = useTags();
  const createMutation = useCreateRace();
  const updateMutation = useUpdateRace();
  const deleteMutation = useDeleteRace();
  const applyTagsMutation = useApplyRaceTags();
  const [showForm, setShowForm] = useState(false);
  const [editingRace, setEditingRace] = useState<Race | null>(null);
  const [formError, setFormError] = useState("");
  const [deletingRace, setDeletingRace] = useState<Race | null>(null);
  const [applyingRace, setApplyingRace] = useState<Race | null>(null);

  const availableTags = (tags ?? []).map(t => t.name);

  const handleCreate = async (input: RaceInput) => {
    setFormError("");
    try { await createMutation.mutateAsync(input); setShowForm(false); }
    catch (e) { setFormError(e instanceof Error ? e.message : "Failed to create race"); }
  };

  const handleUpdate = async (input: RaceInput) => {
    if (!editingRace) return;
    setFormError("");
    try { await updateMutation.mutateAsync({ id: editingRace.id, input }); setEditingRace(null); }
    catch (e) { setFormError(e instanceof Error ? e.message : "Failed to update race"); }
  };

  const handleDelete = async () => {
    if (!deletingRace) return;
    try { await deleteMutation.mutateAsync(deletingRace.id); setDeletingRace(null); }
    catch { /* toasted by global handler */ }
  };

  const handleApplyTags = async () => {
    if (!applyingRace) return;
    try {
      const result = await applyTagsMutation.mutateAsync(applyingRace.id);
      alert(`Applied ${(result.tags_applied ?? []).join(", ") || "no tags"} to ${result.characters_updated} characters`);
      setApplyingRace(null);
    } catch (e) {
      alert(e instanceof Error ? e.message : "Failed to apply tags");
    }
  };

  const columns: Column<Race>[] = [
    ...COLUMNS,
    { header: "Actions", accessor: "_actions", render: (_: unknown, row: Race) => (
      <span className="inline-flex gap-2">
        <Button variant="accent" size="sm" onClick={() => { setEditingRace(row); setShowForm(false); }}>Edit</Button>
        <Button variant="secondary" size="sm" onClick={() => setApplyingRace(row)}>Apply Tags</Button>
        <Button variant="danger" size="sm" className="ml-2" onClick={() => setDeletingRace(row)}>Delete</Button>
      </span>
    )},
  ];

  const isSaving = createMutation.isPending || updateMutation.isPending;

  return (
    <div className="management-page">
      <PageHeader title="Races" backTo="/dashboard" actions={<Button variant="primary" onClick={() => { setShowForm(true); setEditingRace(null); }}>+ Add Race</Button>} />
      {error && <div className="error-banner">{error instanceof Error ? error.message : "Failed to load races"}</div>}
      {formError && <FormError message={formError} />}
      {showForm && !editingRace && <RaceForm race={null} onSubmit={handleCreate} onCancel={() => setShowForm(false)} isLoading={isSaving} error={formError} availableTags={availableTags} />}
      {editingRace && <RaceForm race={editingRace} onSubmit={handleUpdate} onCancel={() => setEditingRace(null)} isLoading={isSaving} error={formError} availableTags={availableTags} />}
      {isLoading ? <div className="loading">Loading races...</div> : (
        <DataTable columns={columns} data={races ?? []} getKey={(row: Race) => row.id} emptyMessage="No races found. Add your first race!" />
      )}
      {deletingRace && <DeleteConfirmation race={deletingRace} onConfirm={handleDelete} onCancel={() => setDeletingRace(null)} isLoading={deleteMutation.isPending} />}
      {applyingRace && <ApplyTagsConfirmation race={applyingRace} onConfirm={handleApplyTags} onCancel={() => setApplyingRace(null)} isLoading={applyTagsMutation.isPending} />}
    </div>
  );
}

function DeleteConfirmation({ race, onConfirm, onCancel, isLoading }: Readonly<{ race: Race; onConfirm: () => void; onCancel: () => void; isLoading: boolean }>) {
  return (
    <div className="modal-overlay" onClick={onCancel}>
      <div className="modal-content modal-sm" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header"><h3>Delete Race</h3><Button variant="ghost" size="sm" onClick={onCancel} aria-label="Close">×</Button></div>
        <div className="modal-body"><p>Are you sure you want to delete <strong>{race.display_name || race.name}</strong>?</p><p className="text-muted">This action cannot be undone.</p></div>
        <div className="modal-footer">
          <Button variant="danger" onClick={onConfirm} disabled={isLoading}>{isLoading ? "Deleting..." : "Delete"}</Button>
          <Button variant="secondary" onClick={onCancel}>Cancel</Button>
        </div>
      </div>
    </div>
  );
}

function ApplyTagsConfirmation({ race, onConfirm, onCancel, isLoading }: Readonly<{ race: Race; onConfirm: () => void; onCancel: () => void; isLoading: boolean }>) {
  return (
    <div className="modal-overlay" onClick={onCancel}>
      <div className="modal-content modal-sm" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header"><h3>Apply Race Tags</h3><Button variant="ghost" size="sm" onClick={onCancel} aria-label="Close">×</Button></div>
        <div className="modal-body">
          <p>Apply tags for <strong>{race.display_name || race.name}</strong> to all characters of this race?</p>
          {(race.tags ?? []).length > 0 && (
            <p className="text-muted mt-2">Tags: {(race.tags ?? []).join(", ")}</p>
          )}
          {(race.tags ?? []).length === 0 && (
            <p className="text-muted mt-2">This race has no tags configured yet. Edit the race first to add tags.</p>
          )}
        </div>
        <div className="modal-footer">
          <Button variant="primary" onClick={onConfirm} disabled={isLoading || (race.tags ?? []).length === 0}>{isLoading ? "Applying..." : "Apply Tags"}</Button>
          <Button variant="secondary" onClick={onCancel}>Cancel</Button>
        </div>
      </div>
    </div>
  );
}