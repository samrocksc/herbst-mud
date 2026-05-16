/* eslint-disable functional/immutable-data */
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useState, useCallback } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiGet, apiPost } from "../utils/apiFetch";
import { PageHeader } from "../components/PageHeader";
import { Button } from "../components/Button";
import { FormField } from "../components/fields/FormField";
import { NumberField } from "../components/fields/NumberField";
import { SelectField } from "../components/fields/SelectField";
import { FormError } from "../components/fields/FormError";
import { ResourceSearchSelect } from "../components/ResourceSearchSelect";
import { RESOURCE_ENDPOINTS } from "../utils/resourceEndpoints";
import type { EquipmentTemplate } from "../components/map/types";

const API = `${window.location.origin}/api`;

export const Route = createFileRoute("/map/rooms/$roomId/items/spawn")({
  component: ItemSpawnPage,
});

function ItemSpawnPage() {
  const { roomId } = Route.useParams();
  const roomIdNum = Number(roomId);
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const { data: templates = [], isLoading: templatesLoading } = useQuery({
    queryKey: ["equipment-templates"],
    queryFn: () => apiGet<EquipmentTemplate[]>(`${API}/equipment-templates`),
  });

  const createMutation = useMutation({
    mutationFn: (input: Record<string, unknown>) => apiPost(`${API}/item-instances`, input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["item-instances"] });
      navigate({ to: "/map", search: { room: roomIdNum } });
    },
  });

  const [templateId, setTemplateId] = useState("");
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [slot, setSlot] = useState("none");
  const [level, setLevel] = useState(0);
  const [weight, setWeight] = useState(0);
  const [color, setColor] = useState("");
  const [spawnRoomId, setSpawnRoomId] = useState<number | string | null>(roomIdNum);

  const selectedTemplate = templates.find((t) => t.equipment_template_id === templateId);

  const applyTemplateDefaults = useCallback((id: string) => {
    const t = templates.find((tmpl) => tmpl.equipment_template_id === id);
    if (!t) return;
    setTemplateId(t.equipment_template_id);
    setName(t.name);
    setDescription(t.description);
    setSlot(t.slot);
    setLevel(t.level);
    setWeight(t.weight);
    setColor(t.color);
  }, [templates]);

  const templateOptions = templates.map((t) => ({
    value: t.equipment_template_id,
    label: `${t.name} (${t.slot}, lv.${t.level})`,
  }));

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!templateId) return;
    const p: Record<string, unknown> = { equipment_template_id: templateId, room_id: Number(spawnRoomId) || roomIdNum };
    if (name.trim()) p.name = name.trim();
    if (description.trim()) p.description = description.trim();
    if (slot && slot !== "none") p.slot = slot;
    if (level > 0) p.level = level;
    if (weight > 0) p.weight = weight;
    if (color.trim()) p.color = color.trim();
    createMutation.mutate(p);
  };

  return (
    <div className="p-6">
      <PageHeader title="Spawn Item Instance" showBack backTo="/map" />
      <div className="bg-surface p-6 border border-border rounded max-w-[600px]">
        <form onSubmit={handleSubmit} className="space-y-3">
          {createMutation.isError && (
            <FormError message={(createMutation.error as Error)?.message || "Failed to spawn item instance"} />
          )}
          <SelectField
            label="Equipment Template *"
            value={templateId}
            onChange={(id) => applyTemplateDefaults(id)}
            options={templateOptions}
            placeholder="-- Select template --"
            disabled={templatesLoading}
          />
          {selectedTemplate && (
            <div className="p-2 bg-surface-muted border border-border rounded text-xs text-text-muted space-y-0.5">
              <div>Slot: {selectedTemplate.slot}</div>
              <div>Level: {selectedTemplate.level} Weight: {selectedTemplate.weight}</div>
              <div>Type: {selectedTemplate.item_type}</div>
            </div>
          )}
          <FormField label="Name" value={name} onChange={setName} />
          <FormField label="Description" value={description} onChange={setDescription} />
          <div className="flex gap-2">
            <div className="flex-1">
              <FormField label="Slot" value={slot} onChange={setSlot} />
            </div>
            <div className="flex-1">
              <NumberField label="Level" value={level} onChange={setLevel} min={0} />
            </div>
          </div>
          <div className="flex gap-2">
            <div className="flex-1">
              <NumberField label="Weight" value={weight} onChange={setWeight} min={0} />
            </div>
            <div className="flex-1">
              <FormField label="Color" value={color} onChange={setColor} placeholder="#8b5cf6" />
            </div>
          </div>
          <ResourceSearchSelect
            label="Room"
            value={spawnRoomId}
            onChange={(id) => setSpawnRoomId(id)}
            {...RESOURCE_ENDPOINTS.rooms}
            placeholder="Search rooms by name or ID..."
          />
          <div className="flex gap-2 pt-2">
            <Button type="submit" variant="primary" disabled={!templateId || createMutation.isPending}>
              {createMutation.isPending ? "Spawning..." : "Spawn Instance"}
            </Button>
            <Button type="button" variant="secondary" onClick={() => navigate({ to: "/map", search: { room: roomIdNum } })}>
              Cancel
            </Button>
          </div>
        </form>
      </div>
    </div>
  );
}
