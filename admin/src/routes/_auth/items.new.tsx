import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useWorldStore } from "../../contexts/WorldStoreContext";
import { useState } from "react";
import { useCreateTemplate } from "../../hooks/useEquipmentTemplates";
import { PageHeader } from "../../components/PageHeader";
import { Button } from "../../components/Button";
import { FormField, NumberField, SelectField, CheckboxField, TextareaField } from "../../components/FormFields";
import { showToast } from "../../components/Toast";
import { DeleteConfirmation } from "../../components/DeleteConfirmation";
import { SLOT_OPTIONS, ITEM_TYPE_OPTIONS } from "../../components/itemConstants";

const EFFECT_TYPE_OPTS = [
  { value: "", label: "— None —" },
  { value: "heal", label: "Heal" },
  { value: "damage", label: "Damage" },
  { value: "dot", label: "DoT (Damage over Time)" },
  { value: "hot", label: "HoT (Heal over Time)" },
  { value: "buff", label: "Buff" },
  { value: "debuff", label: "Debuff" },
  { value: "stun", label: "Stun" },
  { value: "buff_armor", label: "Buff Armor" },
  { value: "buff_dodge", label: "Buff Dodge" },
  { value: "buff_crit", label: "Buff Crit" },
];

export const Route = createFileRoute("/_auth/items/new")({
  component: CreateItemPage,
});

function CreateItemPage() {
  const navigate = useNavigate();
  const { currentWorld } = useWorldStore();
  const { mutate: createTemplate, isPending } = useCreateTemplate();

  const [form, setForm] = useState({
    name: "",
    description: "",
    slot: "",
    item_type: "misc",
    level: 1,
    weight: 0,
    color: "",
    is_visible: true,
    is_immovable: false,
    effect_type: "",
    effect_value: 0,
    effect_duration: 0,
    world_id: currentWorld || "default",
  });

  const [conflict, setConflict] = useState<{ existing: Record<string, unknown> } | null>(null);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!form.name.trim()) return;
    setConflict(null);
    createTemplate({ ...form, world_id: currentWorld || "default" }, {
      onSuccess: () => {
        showToast("Item template created", "success");
        navigate({ to: "/items" });
      },
      onError: (err: unknown) => {
        const apiErr = err as { response?: { status?: number; data?: { existing?: Record<string, unknown> } } };
        if (apiErr?.response?.status === 409 && apiErr.response.data?.existing) {
          setConflict({ existing: apiErr.response.data.existing });
        } else {
          showToast("Failed to create item template", "error");
        }
      },
    });
  };

  return (
    <div className="p-6 max-w-[1200px] mx-auto">
      <PageHeader title="Create Item Template" showBack backTo="/items" />

      <div className="card bg-surface p-6 border border-border rounded">
        <form onSubmit={handleSubmit} className="space-y-4">
          <h3 className="text-text font-semibold mb-4">Basic Information</h3>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <FormField label="Name *" value={form.name} onChange={(v) => setForm({ ...form, name: v })} required />
            <SelectField label="Slot" value={form.slot} onChange={(v) => setForm({ ...form, slot: v })} options={[...SLOT_OPTIONS]} />
            <SelectField label="Item Type" value={form.item_type} onChange={(v) => setForm({ ...form, item_type: v })} options={[...ITEM_TYPE_OPTIONS]} />
            <NumberField label="Level" value={form.level} onChange={(v) => setForm({ ...form, level: v })} min={1} />
            <NumberField label="Weight" value={form.weight} onChange={(v) => setForm({ ...form, weight: v })} min={0} />
            <TextareaField label="Description" value={form.description} onChange={(v) => setForm({ ...form, description: v })} rows={2} />
          </div>

          <h3 className="text-text font-semibold mt-6 mb-4">Properties</h3>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <CheckboxField label="Visible" checked={form.is_visible} onChange={(v) => setForm({ ...form, is_visible: v })} />
            <CheckboxField label="Immovable" checked={form.is_immovable} onChange={(v) => setForm({ ...form, is_immovable: v })} />
          </div>

          <h3 className="text-text font-semibold mt-6 mb-4">Effect</h3>
          <div className="grid grid-cols-3 gap-4">
            <SelectField label="Effect Type" value={form.effect_type} onChange={(v) => setForm({ ...form, effect_type: v })} options={EFFECT_TYPE_OPTS} />
            <NumberField label="Effect Value" value={form.effect_value} onChange={(v) => setForm({ ...form, effect_value: v })} min={0} />
            <NumberField label="Duration (ticks)" value={form.effect_duration} onChange={(v) => setForm({ ...form, effect_duration: v })} min={0} tooltip="0 = instant" />
          </div>

          <div className="flex gap-2 justify-end mt-6">
            <Button variant="secondary" onClick={() => navigate({ to: "/items" })}>Cancel</Button>
            <Button variant="primary" type="submit" disabled={isPending || !form.name.trim()}>
              {isPending ? "Creating…" : "Create"}
            </Button>
          </div>
        </form>
      </div>

      {conflict && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
             onClick={() => setConflict(null)}>
          <div className="bg-surface border border-border rounded-lg p-6 max-w-lg w-full mx-4 shadow-xl"
               onClick={(e) => e.stopPropagation()}>
            <h3 className="text-text font-semibold mb-2">Name Already Exists</h3>
            <p className="text-text-muted text-sm mb-4">
              "{form.name}" is already taken (slug: <code className="bg-surface-dark px-1 rounded">{form.name.toLowerCase().replace(/\s+/g, '_')}</code>).
              You can:
            </p>
            <ul className="text-sm text-text-muted space-y-2 mb-4">
              <li>• Choose a different name</li>
              <li>• Edit the existing template instead</li>
              <li>• Add a suffix like "_2" or " (new)"</li>
            </ul>
            <div className="flex gap-2 justify-end">
              <Button variant="secondary" onClick={() => setConflict(null)}>Keep Editing</Button>
              <Button variant="ghost" onClick={() => { setConflict(null); navigate({ to: "/items" }); }}>
                Back to Items
              </Button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
