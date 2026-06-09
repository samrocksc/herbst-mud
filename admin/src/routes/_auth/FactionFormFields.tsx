import { useState } from "react";
import {
  FormField,
  TextareaField,
  SelectField,
} from "../../components/FormFields";
import { SearchableSelect } from "../../components/SearchableSelect";
import { useTags, useCreateTag } from "../../hooks/useTags";
import { showToast } from "../../components/Toast";
import type { FactionCategory, FactionForm } from "./-factionTypes";

const FACTION_EMOJI_OPTS = ["🍕", "🥷", "📰", "🗡️", "🛡️", "🧙", "🏹", "⚔️", "🐉", "👑", "🔮", "💀", "🎭", "🦊", "🐺", "🦅", "🌟", "🍺", "📚", "⚒️", "🎩", "🌹", "🎲", "🎨", "🎪", "🎭"];

export function FactionFormFields({
  form,
  setForm,
  categories,
}: Readonly<{
  form: FactionForm
  setForm: (f: FactionForm) => void
  categories: FactionCategory[]
}>) {
  const set = (patch: Partial<FactionForm>) => setForm({ ...form, ...patch });
  const { data: tags } = useTags();
  const createTag = useCreateTag();
  const [showTagInput, setShowTagInput] = useState(false);
  const [newTag, setNewTag] = useState("");

  const catOptions = [
    { value: "", label: "— None —" },
    ...categories.map((c) => ({ value: String(c.id), label: c.display_name || c.name })),
  ];

  const tagOptions = (tags ?? []).map((t) => ({ id: t.name, name: t.name }));
  const emoji = (form as FactionForm & { emoji?: string }).emoji;

  const handleCreateTag = async () => {
    const trimmed = newTag.trim();
    if (!trimmed) return;
    try {
      await createTag.mutateAsync({ name: trimmed, color: "#888888" });
      if (!form.member_tags.includes(trimmed)) {
        set({ member_tags: [...form.member_tags, trimmed] });
      }
      setNewTag("");
      setShowTagInput(false);
      showToast(`Tag "${trimmed}" created`, "success");
    } catch (err) {
      showToast(err instanceof Error ? err.message : "Failed to create tag", "error");
    }
  };

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
              onClick={() => set({ emoji: emoji === e ? undefined : e } as Partial<FactionForm>)}
              className={`w-9 h-9 rounded text-lg flex items-center justify-center border ${
                emoji === e
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
        onChange={(v) => {
          set({ name: v, display_name: form.display_name === form.name || !form.display_name ? v.replace(/_/g, " ").replace(/\b\w/g, c => c.toUpperCase()) : form.display_name });
        }}
        placeholder="pizza_chief"
        tooltip="Identifier used in code and content files (snake_case). Display Name auto-fills from this."
      />
      <FormField
        label="Display Name"
        value={form.display_name}
        onChange={(v) => set({ display_name: v })}
        placeholder="Pizza Chief"
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
      <div>
        <label className="text-text-muted text-xs block mb-1">Member Tags</label>
        {form.member_tags.length > 0 && (
          <div className="flex flex-wrap gap-1.5 mb-2">
            {form.member_tags.map((tag) => (
              <span
                key={tag}
                className="inline-flex items-center gap-1 px-2 py-0.5 bg-primary/10 border border-primary/30 rounded text-sm text-text"
              >
                {tag}
                <button
                  type="button"
                  onClick={() => set({ member_tags: form.member_tags.filter((t) => t !== tag) })}
                  className="text-text-muted hover:text-danger px-0.5 text-xs"
                  aria-label={`Remove ${tag}`}
                >
                  ✕
                </button>
              </span>
            ))}
          </div>
        )}
        <div className="flex gap-2 items-center">
          <SearchableSelect
            options={tagOptions}
            value=""
            onChange={(v) => {
              if (!form.member_tags.includes(v)) {
                set({ member_tags: [...form.member_tags, v] });
              }
            }}
            placeholder="Add existing tag..."
          />
          {!showTagInput ? (
            <button
              type="button"
              onClick={() => setShowTagInput(true)}
              className="text-xs text-primary hover:underline shrink-0"
            >
              + New tag
            </button>
          ) : (
            <div className="flex gap-1 items-center shrink-0">
              <input
                type="text"
                value={newTag}
                onChange={(e) => setNewTag(e.target.value)}
                onKeyDown={(e) => { if (e.key === "Enter") { e.preventDefault(); void handleCreateTag(); } }}
                placeholder="new_tag"
                className="px-2 py-1 bg-surface border border-border rounded text-sm w-24"
                autoFocus
              />
              <button
                type="button"
                onClick={handleCreateTag}
                disabled={!newTag.trim() || createTag.isPending}
                className="text-xs px-2 py-1 bg-primary text-white rounded disabled:opacity-50"
              >
                Add
              </button>
              <button
                type="button"
                onClick={() => { setShowTagInput(false); setNewTag(""); }}
                className="text-xs text-text-muted hover:text-text"
              >
                ✕
              </button>
            </div>
          )}
        </div>
        <p className="text-xs text-text-muted mt-1">Tags auto-applied to characters when they join this faction. Create new ones inline or pick from existing.</p>
      </div>
    </div>
  );
}
