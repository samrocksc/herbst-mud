import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useWorldStore } from "../../contexts/WorldStoreContext";
import { useState } from "react";
import { useCreateTemplate } from "../../hooks/useEquipmentTemplates";
import { PageHeader } from "../../components/PageHeader";
import { Button } from "../../components/Button";
import { FormField, NumberField, SelectField, CheckboxField, TextareaField } from "../../components/FormFields";
import { ResourceSearchSelect } from "../../components/ResourceSearchSelect";
import { RESOURCE_ENDPOINTS } from "../../utils/resourceEndpoints";
import { CombatFieldsEditor, type CombatFields } from "../../components/CombatFieldsEditor";
import { showToast } from "../../components/Toast";
import { SLOT_OPTIONS, ITEM_TYPE_OPTIONS } from "../../components/itemConstants";
import { PageContainer } from "../../components/PageContainer";

/** Map item type to default slot for common equipment patterns. */
const getDefaultSlotForItemType = (itemType: string): string => {
  const armorSlots = ["head", "chest", "hands", "legs", "feet"];
  const weaponSlots = ["main_hand", "off_hand"];

  switch (itemType) {
    case "armor":
      return armorSlots[0]; // head for armor
    case "weapon":
      return weaponSlots[0]; // main_hand for weapons
    case "consumable":
    case "potion":
      return "neck"; // consumables often go in neck slot
    case "container":
      return "back"; // containers on back
    default:
      return ""; // empty for misc/quest items
  }
};

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

export function CreateItemPage() {
  const navigate = useNavigate();
  const { currentWorld } = useWorldStore();
  const { mutate: createTemplate, isPending } = useCreateTemplate();

  const [form, setForm] = useState<{
    name: string; description: string; slot: string; level: number; weight: number
    item_type: string; color: string; is_visible: boolean; is_immovable: boolean
    effect_type: string; effect_value: number; effect_duration: number
    is_container: boolean; container_capacity: number; is_locked: boolean
    key_item_id: string; reveal_condition: string
    world_id: string
  } & CombatFields>({
    name: "",
    description: "",
    slot: getDefaultSlotForItemType("armor"),
    item_type: "armor",
    level: 1,
    weight: 0,
    color: "",
    is_visible: true,
    is_immovable: false,
    effect_type: "",
    effect_value: 0,
    effect_duration: 0,
    is_container: false,
    container_capacity: 0,
    is_locked: false,
    key_item_id: "",
    reveal_condition: "",
    armor_rating: 0,
    armor_type: "",
    rarity: "",
    skill_requirement: "",
    skill_requirement_level: 0,
    damage_dice_count: 0,
    damage_dice_sides: 0,
    damage_bonus: 0,
    damage_type: "",
    weapon_type: "",
    is_two_handed: false,
    stats: "{}",
    world_id: currentWorld || "default",
  });

  const [conflict, setConflict] = useState<{ existing: Record<string, unknown> } | null>(null);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!form.name.trim()) return;
    setConflict(null);
    const body: Record<string, unknown> = { ...form, world_id: currentWorld || "default" };
    try { body.stats = JSON.parse(form.stats); } catch { body.stats = {}; }
    createTemplate(body, {
      onSuccess: () => {
        showToast("Item template created", "success");
        navigate({ to: "/items" });
      },
      onError: (err: unknown) => {
        console.error("Item template creation error:", err);
        const apiErr = err as { response?: { status?: number; data?: { existing?: Record<string, unknown> } } };
        if (apiErr?.response?.status === 409 && apiErr.response.data?.existing) {
          setConflict({ existing: apiErr.response.data.existing });
        } else {
          const message = err instanceof Error ? err.message : "Failed to create item template";
          showToast(message, "error");
        }
      },
    });
  };

  return (
    <PageContainer>
      <PageHeader title="Create Item Template" showBack backTo="/items" />

      <div className="card bg-surface p-6 border border-border rounded">
        <form onSubmit={handleSubmit} className="space-y-4">
          <h3 className="text-text font-semibold mb-4">Basic Information</h3>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <FormField label="Name *" value={form.name} onChange={(v) => setForm({ ...form, name: v })} required />
            <SelectField label="Slot" value={form.slot} onChange={(v) => setForm({ ...form, slot: v })} options={[...SLOT_OPTIONS]} />
            <SelectField label="Item Type" value={form.item_type} onChange={(v) => {
              setForm({ ...form, item_type: v, slot: getDefaultSlotForItemType(v) });
            }} options={[...ITEM_TYPE_OPTIONS]} />
            <NumberField label="Level" value={form.level} onChange={(v) => setForm({ ...form, level: v })} min={1} />
            <NumberField label="Weight" value={form.weight} onChange={(v) => setForm({ ...form, weight: v })} min={0} />
            <TextareaField label="Description" value={form.description} onChange={(v) => setForm({ ...form, description: v })} rows={2} />
          </div>

          <h3 className="text-text font-semibold mt-6 mb-4">Properties</h3>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <CheckboxField label="Visible" checked={form.is_visible} onChange={(v) => setForm({ ...form, is_visible: v })} />
            <CheckboxField label="Immovable" checked={form.is_immovable} onChange={(v) => setForm({ ...form, is_immovable: v })} />
            <CheckboxField label="Container" checked={form.is_container} onChange={(v) => setForm({ ...form, is_container: v })} />
          </div>
          {form.is_container && (
            <div className="grid grid-cols-3 gap-4 mt-2">
              <NumberField label="Container Capacity" value={form.container_capacity} onChange={(v) => setForm({ ...form, container_capacity: v })} min={0} />
              <CheckboxField label="Locked" checked={form.is_locked} onChange={(v) => setForm({ ...form, is_locked: v })} />
              <ResourceSearchSelect
                label="Key Item"
                value={form.key_item_id || null}
                onChange={(id) => setForm({ ...form, key_item_id: String(id ?? "") })}
                placeholder="Search key item template..."
                {...RESOURCE_ENDPOINTS.equipmentTemplates}
              />
            </div>
          )}
          <div className="mt-2">
            <FormField label="Reveal Condition" value={form.reveal_condition} onChange={(v) => setForm({ ...form, reveal_condition: v })} placeholder='e.g. {"type":"examine","minLevel":3}' tooltip="JSON condition for revealing hidden details" />
          </div>

          <h3 className="text-text font-semibold mt-6 mb-4">Effect</h3>
          <div className="grid grid-cols-3 gap-4">
            <SelectField label="Effect Type" value={form.effect_type} onChange={(v) => setForm({ ...form, effect_type: v })} options={EFFECT_TYPE_OPTS} />
            <NumberField label="Effect Value" value={form.effect_value} onChange={(v) => setForm({ ...form, effect_value: v })} min={0} />
            <NumberField label="Duration (ticks)" value={form.effect_duration} onChange={(v) => setForm({ ...form, effect_duration: v })} min={0} tooltip="0 = instant" />
          </div>

          <div className="mt-4 pt-4 border-t border-border">
            <CombatFieldsEditor form={form} onChange={(u) => setForm(prev => ({ ...prev, ...u }))} slot={form.slot} />
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
    </PageContainer>
  );
}
