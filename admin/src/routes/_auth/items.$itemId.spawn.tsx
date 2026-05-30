/* eslint-disable functional/prefer-immutable-types, functional/immutable-data */
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiGet, apiPost, API_BASE } from "../../utils/apiFetch";
import { PageHeader } from "../../components/PageHeader";
import { Button } from "../../components/Button";
import { FormError } from "../../components/fields/FormError";
import { ResourceSearchSelect } from "../../components/ResourceSearchSelect";
import { RESOURCE_ENDPOINTS } from "../../utils/resourceEndpoints";
import type { EquipmentTemplate as ItemTemplate } from "../../hooks/useEquipmentTemplates";

export const Route = createFileRoute("/_auth/items/$itemId/spawn")({
  component: ItemSpawnPage,
});

type TargetType = "room" | "character"

function ItemSpawnPage() {
  const { itemId } = Route.useParams();
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const templateQuery = useQuery({
    queryKey: ["item-template", itemId],
    queryFn: () => apiGet<ItemTemplate>(`${API_BASE}/api/equipment-templates/${itemId}`),
  });

  const spawnMutation = useMutation({
    mutationFn: ({ targetType, targetId }: { targetType: TargetType; targetId: number }) => {
      const body: Record<string, unknown> = { equipment_template_id: Number(itemId) };
      if (targetType === "room") body.room_id = targetId;
      else body.ownerId = targetId;
      return apiPost(`${API_BASE}/api/item-instances`, body);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["item-instances", "template", itemId] });
      navigate({ to: "/items/$itemId", params: { itemId } });
    },
  });

  const [targetType, setTargetType] = useState<TargetType>("room");
  const [targetId, setTargetId] = useState<number | string | null>(null);

  const canSpawn = targetId != null && Number(targetId) > 0;
  const error = spawnMutation.isError ? (spawnMutation.error as Error)?.message ?? "Failed" : null;

  if (templateQuery.isLoading) {
    return <div className="p-8"><PageHeader title="Loading..." backTo="/items" /></div>;
  }
  if (templateQuery.error || !templateQuery.data) {
    return <div className="p-8"><PageHeader title="Error" backTo="/items" /><div className="text-danger">Failed to load item template</div></div>;
  }

  return (
    <div className="p-6">
      <PageHeader title={`Spawn: ${templateQuery.data.name}`} showBack backTo="/items" />
      <div className="bg-surface p-6 border border-border rounded max-w-[600px]">
        <div className="space-y-4">
          {error && <FormError message={error} />}

          <div>
            <label className="block text-text text-sm font-medium mb-1">Assign to</label>
            <div className="flex gap-2">
              <button
                type="button"
                className={`px-3 py-1.5 rounded text-sm ${targetType === "room" ? "bg-primary text-white" : "bg-surface-muted text-text"}`}
                onClick={() => setTargetType("room")}
              >Room</button>
              <button
                type="button"
                className={`px-3 py-1.5 rounded text-sm ${targetType === "character" ? "bg-primary text-white" : "bg-surface-muted text-text"}`}
                onClick={() => setTargetType("character")}
              >Character</button>
            </div>
          </div>

          {targetType === "room" ? (
            <ResourceSearchSelect
              label="Room"
              value={targetId}
              onChange={(id) => setTargetId(id)}
              {...RESOURCE_ENDPOINTS.rooms}
              placeholder="Search rooms by name or ID..."
            />
          ) : (
            <ResourceSearchSelect
              label="Character"
              value={targetId}
              onChange={(id) => setTargetId(id)}
              {...RESOURCE_ENDPOINTS.characters}
              placeholder="Search characters by name or ID..."
            />
          )}

          <div className="flex gap-2 pt-2">
            <Button variant="primary" disabled={!canSpawn || spawnMutation.isPending}
              onClick={() => { if (canSpawn) spawnMutation.mutate({ targetType, targetId: Number(targetId) }); }}>
              {spawnMutation.isPending ? "Spawning..." : "Spawn Instance"}
            </Button>
            <Button variant="secondary" onClick={() => navigate({ to: "/items/$itemId", params: { itemId } })}>
              Cancel
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
}

