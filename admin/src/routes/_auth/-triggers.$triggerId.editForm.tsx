import { useState } from "react";
import {
  useUpdateTrigger,
  useDeleteTrigger,
  type Trigger,
  type TriggerInput,
} from "../../hooks/useTriggers";
import { useNavigate } from "@tanstack/react-router";
import { Button } from "../../components/Button";
import { DeleteConfirmation } from "../../components/DeleteConfirmation";
import {
  FormField,
  TextareaField,
  SelectField,
  CheckboxField,
  NumberField,
} from "../../components/FormFields";
import { ResourceIdField } from "../../components/ResourceIdField";
import { showToast } from "../../components/Toast";
import { useWorldStore } from "../../contexts/WorldStoreContext";

const TRIGGER_TYPE_OPTS = [
  { value: "use",     label: "Use — items, switches, doors" },
  { value: "touch",   label: "Touch — surface interaction" },
  { value: "press",   label: "Press — buttons, levers" },
  { value: "enter",   label: "Enter Room — on room entry" },
  { value: "examine", label: "Examine — detailed inspection" },
  { value: "talk",    label: "Talk — NPC conversation" },
];

const TARGET_TYPE_OPTS = [
  { value: "recipe",      label: "Recipe" },
  { value: "effect",      label: "Effect" },
  { value: "dialog_node", label: "Dialog Node" },
] as const;

const TRIGGER_HELP: Record<string, string> = {
  use: "Fires when a player uses an item or interacts with an object in a room.",
  touch: "Fires when a player touches a surface or object.",
  press: "Fires when a player presses a button or lever.",
  enter: "Fires automatically when a player enters the specified room.",
  examine: "Fires when a player examines an object. Use examine_weight to control reveal order.",
  talk: "Fires when a player starts a conversation with the linked NPC.",
};

export function TriggerEditForm({
  trigger,
  triggerId,
  onDone,
}: Readonly<{
  trigger: Trigger
  triggerId: number
  onDone: () => void
}>) {
  const navigate = useNavigate();
  const updateTrigger = useUpdateTrigger();
  const deleteTrigger = useDeleteTrigger();
  const { currentWorld } = useWorldStore();
  const [showDeleteModal, setShowDeleteModal] = useState(false);

  const [formData, setFormData] = useState<TriggerInput>({
    name: trigger.name,
    world_id: trigger.world_id,
    trigger_type: trigger.trigger_type,
    target_type: trigger.target_type,
    target_id: trigger.target_id || "",
    room_id: trigger.room_id,
    equipment_id: trigger.equipment_id,
    condition: trigger.condition,
    enabled: trigger.enabled,
    examine_weight: trigger.examine_weight ?? 0,
  });

  const set = (patch: Partial<TriggerInput>) => setFormData((prev) => ({ ...prev, ...patch }));

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await updateTrigger.mutateAsync({ id: triggerId, input: formData });
      showToast("Trigger updated", "success");
      onDone();
    } catch {
      // Error is toasted by global onError handler
    }
  };

  const handleDelete = async () => {
    try {
      await deleteTrigger.mutateAsync(triggerId);
      showToast("Trigger deleted", "success");
      navigate({ to: "/triggers" });
    } catch {
      // Error is toasted by global onError handler
    }
  };

  return (
    <div className="bg-surface-muted rounded-lg p-6 border border-border mb-6">
      <h2 className="mt-0 mb-4 text-text text-lg font-semibold">Edit Trigger</h2>
      <form onSubmit={handleSubmit} className="space-y-4">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <FormField label="Name" value={formData.name} onChange={(v) => set({ name: v })} />
          <SelectField label="Trigger Type" value={formData.trigger_type} onChange={(v) => set({ trigger_type: v })} options={TRIGGER_TYPE_OPTS} />
          <SelectField label="Target Type" value={formData.target_type} onChange={(v) => set({ target_type: v })} options={TARGET_TYPE_OPTS} />
          <ResourceIdField
            label="Target ID"
            value={formData.target_id ?? ""}
            onChange={(v) => set({ target_id: v || "" })}
            resourceType="targets"
            apiBase=""
            worldId={currentWorld}
          />
        </div>

        {/* Help for selected trigger type */}
        <div className="bg-surface-muted border border-border rounded p-3">
          <p className="text-xs text-text-muted leading-relaxed">{TRIGGER_HELP[formData.trigger_type] || "No description available for this trigger type."}</p>
        </div>

        <h3 className="text-text font-semibold mt-6 mb-4">Target Object</h3>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <ResourceIdField
            label="Room ID (optional)"
            value={formData.room_id ?? ""}
            onChange={(v) => set({ room_id: v === null || v === "" ? null : Number(v) > 0 ? Number(v) : null })}
            resourceType="rooms"
            apiBase={window.location.origin}
          />
          <ResourceIdField
            label="Equipment ID (optional)"
            value={formData.equipment_id ?? ""}
            onChange={(v) => set({ equipment_id: v === null || v === "" ? null : Number(v) > 0 ? Number(v) : null })}
            resourceType="equipment-templates"
            apiBase={window.location.origin}
          />
        </div>

        <h3 className="text-text font-semibold mt-6 mb-4">Conditions & Settings</h3>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <NumberField
              label="Examine Weight"
              value={formData.examine_weight ?? 0}
              onChange={(v) => set({ examine_weight: v })}
              placeholder="0 = always fires"
            />
            <div className="mt-1">
              <a href="/docs/triggers" target="_blank" rel="noopener noreferrer" className="text-xs text-primary hover:underline">
                Learn more about conditions
              </a>
            </div>
          </div>
          <TextareaField label="Condition (SPICE expression, optional)" value={formData.condition} onChange={(v) => set({ condition: v })} rows={3} placeholder="e.g., player_level >= 10" />
          <div className="flex items-end">
            <CheckboxField label="Enabled" checked={formData.enabled} onChange={(v) => set({ enabled: v })} />
          </div>
        </div>

        <div className="flex gap-2 pt-2">
          <Button type="submit" variant="primary" disabled={updateTrigger.isPending}>
            {updateTrigger.isPending ? "Saving..." : "Save Changes"}
          </Button>
          <Button variant="secondary" onClick={onDone} type="button">Cancel</Button>
          <Button variant="danger" onClick={() => setShowDeleteModal(true)} type="button">Delete Trigger</Button>
        </div>
      </form>

      <DeleteConfirmation
        open={showDeleteModal}
        title="Delete Trigger"
        message="Are you sure you want to delete this trigger? This action cannot be undone."
        onConfirm={handleDelete}
        onCancel={() => setShowDeleteModal(false)}
        isLoading={deleteTrigger.isPending}
      />
    </div>
  );
}
