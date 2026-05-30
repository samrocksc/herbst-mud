import {
  FormField,
  TextareaField,
  SelectField,
} from "../../components/FormFields";
import { TagInput } from "../../components/TagInput";
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
  const availableTags = (tags ?? []).map((t) => t.name);
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
      <TagInput
        label="Member Tags"
        value={form.member_tags}
        onChange={(tags) => set({ member_tags: tags })}
        availableTags={availableTags}
        placeholder="Tags auto-applied when characters join"
      />
    </div>
  );
}