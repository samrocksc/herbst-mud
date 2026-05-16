/* eslint-disable functional/prefer-immutable-types */
import { useState } from "react";
import type { Tag, TagInput } from "../../hooks/useTags";
import { FormField } from "../../components/fields/FormField";
import { ColorField } from "../../components/fields/ColorField";
import { FormError } from "../../components/fields/FormError";
import { Button } from "../../components/Button";

const DEFAULT_COLOR = "var(--color-tag-default)";

export function TagForm({
  tag,
  onSubmit,
  onCancel,
  isLoading,
  error,
}: {
  tag: Tag | null
  onSubmit: (data: TagInput) => void
  onCancel: () => void
  isLoading: boolean
  error: string | null
}) {
  const [form, setForm] = useState<TagInput>(() => ({
    name: tag?.name ?? "",
    color: tag?.color ?? DEFAULT_COLOR,
  }));

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!form.name.trim()) return;
    onSubmit({ name: form.name.trim(), color: form.color });
  };

  return (
    <div className="form-card">
      <h3>{tag ? "Edit Tag" : "Add New Tag"}</h3>
      {error && <FormError message={error} />}
      <form onSubmit={handleSubmit}>
        <div className="form-row">
          <FormField
            label="Name"
            value={form.name}
            onChange={(name) => setForm({ ...form, name })}
            placeholder="e.g. fire, magic, warrior"
            required
          />
        </div>

        <div className="form-row">
          <ColorField
            label="Color"
            value={form.color || DEFAULT_COLOR}
            onChange={(color) => setForm({ ...form, color })}
            placeholder="CSS color / hex"
          />
        </div>

        <div className="form-actions">
          <Button type="submit" variant="primary" disabled={isLoading || !form.name.trim()}>
            {isLoading ? "Saving…" : tag ? "Update Tag" : "Create Tag"}
          </Button>
          <Button type="button" variant="secondary" onClick={onCancel}>
            Cancel
          </Button>
        </div>
      </form>
    </div>
  );
}