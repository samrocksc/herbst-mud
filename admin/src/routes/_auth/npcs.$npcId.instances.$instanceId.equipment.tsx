import { createFileRoute } from "@tanstack/react-router";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";
import { apiGet, apiPost, apiDelete, apiPut, API_BASE } from "../../utils/apiFetch";
import { PageHeader } from "../../components/PageHeader";
import { Button } from "../../components/Button";
import { DataTable, type Column } from "../../components/DataTable";
import { ResourceSearchSelect } from "../../components/ResourceSearchSelect";
import { RESOURCE_ENDPOINTS } from "../../utils/resourceEndpoints";

export const Route = createFileRoute("/_auth/npcs/$npcId/instances/$instanceId/equipment")({
  component: NPCInstanceEquipmentPage,
});

type EquipmentItem = Readonly<{
  id: number;
  name: string;
  slot: string;
  item_type: string;
  is_equipped: boolean;
}>;

function NPCInstanceEquipmentPage() {
  const { npcId, instanceId } = Route.useParams();
  const queryClient = useQueryClient();
  const [templateId, setTemplateId] = useState<number | null>(null);

  const { data, isLoading, error } = useQuery<EquipmentItem[]>({
    queryKey: ["npc-instance-equipment", instanceId],
    queryFn: () => apiGet<EquipmentItem[]>(`${API_BASE}/api/npc-instances/${instanceId}/equipment`),
  });

  const addMutation = useMutation({
    mutationFn: () => apiPost(`${API_BASE}/api/npc-instances/${instanceId}/equipment`, { equipment_template_id: templateId }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["npc-instance-equipment", instanceId] });
      setTemplateId(null);
    },
  });

  const toggleMutation = useMutation({
    mutationFn: ({ id, isEquipped }: { id: number; isEquipped: boolean }) =>
      apiPut(`${API_BASE}/api/npc-instances/${instanceId}/equipment/${id}`, { is_equipped: isEquipped }),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["npc-instance-equipment", instanceId] }),
  });

  const deleteMutation = useMutation({
    mutationFn: (id: number) => apiDelete(`${API_BASE}/api/npc-instances/${instanceId}/equipment/${id}`),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["npc-instance-equipment", instanceId] }),
  });

  if (isLoading) return <LoadingState npcId={npcId} instanceId={instanceId} />;
  if (error || !data) return <ErrorState npcId={npcId} instanceId={instanceId} error={error?.message} />;

  const columns: Column<EquipmentItem>[] = [
    { header: "ID", accessor: "id", align: "center" },
    { header: "Name", accessor: "name" },
    { header: "Slot", accessor: "slot" },
    { header: "Equipped", accessor: "is_equipped", align: "center" },
    {
      header: "",
      accessor: "_actions",
      render: (_value, row) => (
        <div className="flex items-center justify-end gap-2">
          <Button
            variant={row.is_equipped ? "secondary" : "success"}
            size="sm"
            onClick={() => toggleMutation.mutate({ id: row.id, isEquipped: !row.is_equipped })}
            disabled={toggleMutation.isPending}
          >
            {row.is_equipped ? "Unequip" : "Equip"}
          </Button>
          <Button
            variant="danger"
            size="sm"
            onClick={() => deleteMutation.mutate(row.id)}
            disabled={deleteMutation.isPending}
          >
            Remove
          </Button>
        </div>
      ),
    },
  ];

  return (
    <div className="p-8">
      <PageHeader title="Equipment" backTo={`/npcs/${npcId}/instances/${instanceId}`} />

      <div className="bg-surface-muted rounded-lg p-6 border border-border mb-6">
        <h2 className="mt-0 mb-4 text-text text-lg font-semibold">Add Equipment</h2>
        <div className="flex gap-3 items-end max-w-xl">
          <div className="flex-1">
            <ResourceSearchSelect
              label="Equipment Template"
              value={templateId}
              onChange={(id) => setTemplateId(id ? Number(id) : null)}
              placeholder="Search equipment template..."
              {...RESOURCE_ENDPOINTS.equipmentTemplates}
            />
          </div>
          <Button
            variant="primary"
            size="sm"
            onClick={() => addMutation.mutate()}
            disabled={!templateId || addMutation.isPending}
          >
            {addMutation.isPending ? "Adding..." : "Add"}
          </Button>
        </div>
        {addMutation.isError && (
          <div className="mt-3 p-2 bg-danger/10 border border-danger rounded text-danger text-xs">
            Failed to add: {(addMutation.error as Error)?.message}
          </div>
        )}
      </div>

      <DataTable<EquipmentItem>
        columns={columns}
        data={data}
        getKey={(row) => row.id}
        emptyMessage="No equipment on this NPC instance."
        variant="dark"
      />
    </div>
  );
}

function LoadingState({ npcId, instanceId }: Readonly<{ npcId: string; instanceId: string }>) {
  return (
    <div className="p-8">
      <PageHeader title="Loading..." backTo={`/npcs/${npcId}/instances/${instanceId}`} />
      <div className="text-text-muted">Loading equipment...</div>
    </div>
  );
}

function ErrorState({ npcId, instanceId, error }: Readonly<{ npcId: string; instanceId: string; error?: string }>) {
  return (
    <div className="p-8">
      <PageHeader title="Error" backTo={`/npcs/${npcId}/instances/${instanceId}`} />
      <div className="text-danger">Failed to load equipment: {error ?? "Unknown error"}</div>
    </div>
  );
}
