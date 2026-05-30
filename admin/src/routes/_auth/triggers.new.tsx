import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import {
  useCreateTrigger,
  type TriggerInput,
} from "../../hooks/useTriggers";
import { PageHeader } from "../../components/PageHeader";
import { Button } from "../../components/Button";
import {
  FormField,
  NumberField,
  TextareaField,
  SelectField,
  CheckboxField,
} from "../../components/FormFields";
import { showToast } from "../../components/Toast";
import { PageContainer } from "../../components/PageContainer";

export const Route = createFileRoute("/_auth/triggers/new")({
  component: CreateTriggerPage,
});

const TRIGGER_TYPE_OPTS = [
  { value: "use", label: "Use" },
  { value: "touch", label: "Touch" },
  { value: "press", label: "Press" },
  { value: "enter", label: "Enter Room" },
  { value: "examine", label: "Examine" },
];

const TARGET_TYPE_OPTS = [
  { value: "recipe", label: "Recipe" },
  { value: "effect", label: "Effect" },
  { value: "dialog_node", label: "Dialog Node" },
];

const EMPTY_TRIGGER: TriggerInput = {
  name: "",
  world_id: "1",
  trigger_type: "use",
  target_type: "recipe",
  target_id: 0,
  room_id: null,
  equipment_id: null,
  condition: "",
  enabled: true,
};

export function CreateTriggerPage() {
  const navigate = useNavigate();
  const createTrigger = useCreateTrigger();
  const [formData, setFormData] = useState<TriggerInput>(EMPTY_TRIGGER);

  const set = (patch: Partial<TriggerInput>) => setFormData((prev) => ({ ...prev, ...patch }));

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await createTrigger.mutateAsync(formData);
      showToast("Trigger created", "success");
      navigate({ to: "/triggers" });
    } catch {
      // Error is toasted by global onError handler
    }
  };

  return (
    <PageContainer>
      <PageHeader title="Create Trigger" showBack backTo="/triggers" />
      <div className="card bg-surface p-6 border border-border rounded">
        <form onSubmit={handleSubmit} className="space-y-4">
          <h3 className="text-text font-semibold mb-4">Basic Information</h3>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <FormField label="Name *" value={formData.name} onChange={(v) => set({ name: v })} required />
            <FormField label="World ID *" value={formData.world_id} onChange={(v) => set({ world_id: v })} />
            <SelectField label="Trigger Type" value={formData.trigger_type} onChange={(v) => set({ trigger_type: v })} options={TRIGGER_TYPE_OPTS} />
            <SelectField label="Target Type" value={formData.target_type} onChange={(v) => set({ target_type: v })} options={TARGET_TYPE_OPTS} />
            <NumberField label="Target ID *" value={formData.target_id} onChange={(v) => set({ target_id: v })} />
          </div>

          <h3 className="text-text font-semibold mt-6 mb-4">Target Object</h3>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <NumberField label="Room ID (optional)" value={formData.room_id ?? 0} onChange={(v) => set({ room_id: v > 0 ? v : null })} placeholder="0 = none" />
            <NumberField label="Equipment ID (optional)" value={formData.equipment_id ?? 0} onChange={(v) => set({ equipment_id: v > 0 ? v : null })} placeholder="0 = none" />
          </div>

          <h3 className="text-text font-semibold mt-6 mb-4">Conditions & Settings</h3>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <TextareaField label="Condition (SPICE expression, optional)" value={formData.condition} onChange={(v) => set({ condition: v })} rows={3} placeholder="e.g., player_level >= 10" />
            <div className="flex items-end">
              <CheckboxField label="Enabled" checked={formData.enabled} onChange={(v) => set({ enabled: v })} />
            </div>
          </div>

          <div className="flex gap-2 justify-end mt-6">
            <Button variant="secondary" onClick={() => navigate({ to: "/triggers" })}>Cancel</Button>
            <Button variant="primary" type="submit" disabled={createTrigger.isPending}>
              {createTrigger.isPending ? "Creating..." : "Create Trigger"}
            </Button>
          </div>
        </form>
      </div>
    </PageContainer>
  );
}
