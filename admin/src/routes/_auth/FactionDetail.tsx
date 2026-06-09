import { useState } from "react";
import { Button } from "../../components/Button";
import { DeleteConfirmation } from "../../components/DeleteConfirmation";
import { FormError } from "../../components/fields/FormError";
import { FormField, TextareaField } from "../../components/FormFields";
import { showToast } from "../../components/Toast";
import { apiPut, apiDelete } from "../../utils/apiFetch";
import { factionToForm, type Faction, type FactionForm, type FactionCategory } from "./-factionTypes";
import { useWorldStore } from "../../contexts/WorldStoreContext";

const FACTION_EMOJI_OPTS = ["🍕", "🥷", "📰", "🗡️", "🛡️", "🧙", "🏹", "⚔️", "🐉", "👑", "🔮", "💀", "🎭", "🦊", "🐺", "🦅", "🌟", "🍺", "📚", "⚒️"];

export function FactionDetail({
  faction,
  categories,
  onRefresh,
}: Readonly<{
  faction: Faction
  categories: FactionCategory[]
  onRefresh: () => void
}>) {
  const { currentWorld } = useWorldStore();
  const [editing, setEditing] = useState(false);
  const [form, setForm] = useState<FactionForm>(factionToForm(faction));
  const [saving, setSaving] = useState(false);
  const [confirmDelete, setConfirmDelete] = useState(false);
  const [error, setError] = useState("");

  const handleUpdate = async () => {
    if (!form.name?.trim()) {
      setError("Name is required");
      return;
    }
    setSaving(true);
    setError("");
    try {
      await apiPut(`/api/factions/${faction.id}`, {
        ...form,
        category_id: form.category_id || null,
      });
      showToast("Faction updated", "success");
      setEditing(false);
      onRefresh();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Update failed");
    } finally {
      setSaving(false);
    }
  };

  const handleDelete = async () => {
    await apiDelete(`/api/factions/${faction.id}`);
    showToast("Faction deleted", "success");
    onRefresh();
  };

  const cat = categories.find((c) => c.id === faction.category_id);

  if (editing) {
    return (
      <div>
        <h2 className="mt-0 mb-4 text-text">Edit Faction</h2>
        {error && <FormError message={error} />}
        <FactionFormFields form={form} setForm={setForm} categories={categories} />
        <p className="text-xs text-text-muted mt-3">
          Editing faction in world <code className="text-text">{currentWorld}</code>.
        </p>
        <div className="flex gap-2 mt-3">
          <Button variant="primary" size="md" fullWidth onClick={handleUpdate} disabled={saving}>
            {saving ? "Saving..." : "Save Changes"}
          </Button>
          <Button variant="secondary" size="md" fullWidth onClick={() => setEditing(false)}>
            Cancel
          </Button>
        </div>
      </div>
    );
  }

  const memberTags = faction.member_tags ?? [];
  const emoji = (faction as Faction & { emoji?: string }).emoji;

  return (
    <div>
      <h2 className="mt-0 mb-4 text-text flex items-center gap-2">
        {emoji && <span className="text-2xl">{emoji}</span>}
        {faction.display_name || faction.name}
      </h2>
      <div className="bg-surface-muted rounded-lg p-4 border border-border space-y-2">
        <DetailRow label="ID" value={String(faction.id)} />
        <DetailRow label="Name" value={faction.name} />
        <DetailRow label="Display Name" value={faction.display_name || "—"} />
        {faction.description && <DetailRow label="Description" value={faction.description} />}
        {cat && (
          <DetailRow
            label="Category"
            value={
              <span className="inline-flex items-center gap-1">
                <span className="px-2 py-0.5 rounded-full text-xs bg-accent/20 text-accent">{cat.display_name || cat.name}</span>
                {cat.max_memberships > 1 && (
                  <span className="text-xs text-text-muted">({cat.max_memberships} max)</span>
                )}
              </span>
            }
          />
        )}
        {cat?.auto_join && (
          <div className="detail-row">
            <label>&nbsp;</label>
            <span className="px-2 py-0.5 rounded-full text-xs bg-success/20 text-success">Auto-Join</span>
          </div>
        )}
        <DetailRow label="Members" value={String(faction.member_count ?? 0)} />
        {memberTags.length > 0 && (
          <div className="detail-row">
            <label>Member Tags</label>
            <div className="flex flex-wrap gap-1">
              {memberTags.map((tag) => (
                <span key={tag} className="px-2 py-0.5 bg-primary/20 text-primary text-xs rounded-full">{tag}</span>
              ))}
            </div>
          </div>
        )}
        <div className="flex gap-2 mt-3">
          <Button variant="primary" size="md" fullWidth onClick={() => setEditing(true)}>Edit</Button>
          <Button variant="danger" size="md" fullWidth onClick={() => setConfirmDelete(true)}>Delete</Button>
        </div>
      </div>
      <DeleteConfirmation
        open={confirmDelete}
        title="Delete Faction"
        message={`Are you sure you want to delete "${faction.display_name || faction.name}"? This cannot be undone.`}
        onConfirm={handleDelete}
        onCancel={() => setConfirmDelete(false)}
      />
    </div>
  );
}

type DetailRowValue = string | number | React.ReactNode;
function DetailRow({ label, value }: Readonly<{ label: string; value: DetailRowValue }>) {
  return (
    <div className="detail-row">
      <label>{label}</label>
      <span className="text-text text-sm">{value}</span>
    </div>
  );
}

/** Inline form fields shared by create + edit. */
function FactionFormFields({
  form,
  setForm,
  categories,
}: Readonly<{
  form: FactionForm
  setForm: (f: FactionForm) => void
  categories: FactionCategory[]
}>) {
  const set = (patch: Partial<FactionForm>) => setForm({ ...form, ...patch });
  const catOptions = [
    { value: "", label: "— None —" },
    ...categories.map((c) => ({ value: String(c.id), label: c.display_name || c.name })),
  ];

  return (
    <div className="space-y-3">
      {/* Emoji picker */}
      <div>
        <label className="text-text-muted text-xs block mb-1">Icon (optional)</label>
        <div className="flex flex-wrap gap-1">
          {FACTION_EMOJI_OPTS.map((e) => (
            <button
              key={e}
              type="button"
              onClick={() => set({ emoji: (form as FactionForm & { emoji?: string }).emoji === e ? undefined : e } as Partial<FactionForm>)}
              className={`w-9 h-9 rounded text-lg flex items-center justify-center border ${
                (form as FactionForm & { emoji?: string }).emoji === e
                  ? "bg-primary/30 border-primary"
                  : "bg-surface border-border hover:bg-surface-hover"
              }`}
            >
              {e}
            </button>
          ))}
        </div>
        <p className="text-xs text-text-muted mt-1">Pick an icon to help players identify this faction quickly.</p>
      </div>
      <FormField
        label="Name *"
        value={form.name}
        onChange={(v) => set({ name: v })}
        placeholder="pizza_chef"
        tooltip="Identifier used in code and content files (snake_case)."
      />
      <FormField
        label="Display Name"
        value={form.display_name}
        onChange={(v) => set({ display_name: v })}
        placeholder="Pizza Chef (auto-fills from Name)"
      />
      <TextareaField
        label="Description"
        value={form.description}
        onChange={(v) => set({ description: v })}
        rows={3}
        placeholder="What is this faction about? How do players join?"
      />
      <SelectField
        label="Category"
        value={String(form.category_id)}
        onChange={(v) => set({ category_id: v ? Number(v) : "" })}
        options={catOptions}
      />
    </div>
  );
}
