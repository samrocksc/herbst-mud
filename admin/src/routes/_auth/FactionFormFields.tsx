import {
  FormField,
  TextareaField,
  SelectField,
} from "../../components/FormFields";
import { SearchableSelect } from "../../components/SearchableSelect";
import { useTags } from "../../hooks/useTags";
import type { FactionCategory, FactionForm } from "./-factionTypes";

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
  const catOptions = [
    { value: "", label: "None" },
    ...categories.map((c) => ({ value: String(c.id), label: c.display_name || c.name })),
  ];

  return (
    <div className="space-y-3">
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
        <SearchableSelect
          options={(tags ?? []).map((t) => ({ id: t.name, name: t.name }))}
          value=""
          onChange={(v) => {
            if (!form.member_tags.includes(v)) {
              set({ member_tags: [...form.member_tags, v] });
            }
          }}
          placeholder="Add existing tag..."
        />
      </div>
    </div>
  );
}