import { createFileRoute, useNavigate } from '@tanstack/react-router';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useState } from 'react';
import { apiGet, apiPut, apiDelete } from '../../utils/apiFetch';
import { PageHeader } from '../../components/PageHeader';
import { Button } from '../../components/Button';
import { DeleteConfirmation } from '../../components/DeleteConfirmation';

// ─── Types ──────────────────────────────────────────────────────────────────

type NPCInstance = Readonly<{
  id: number
  name: string
  instance_number: number
  room_id: number
  starting_room_id: number
  hitpoints: number
  max_hitpoints: number
  level: number
  race: string
  isNPC: boolean
  is_instance: boolean
}>

type EditForm = {
  room_id: number
  starting_room_id: number
  instance_number: number
  hitpoints: number
}

const API = `${window.location.origin}`;

// ─── Route ─────────────────────────────────────────────────────────────────

export const Route = createFileRoute('/_auth/npcs/$npcId/instances/$instanceId')({
  component: NpcInstanceDetail,
});

// ─── Component ──────────────────────────────────────────────────────────────

function NpcInstanceDetail() {
  const { npcId, instanceId } = Route.useParams();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [isEditing, setIsEditing] = useState(false);
  const [showDeleteModal, setShowDeleteModal] = useState(false);

  const [form, setForm] = useState<EditForm | null>(null);

  // ── Queries ────────────────────────────────────────────────────────────────

  const { data: instance, isLoading, error } = useQuery<NPCInstance>({
    queryKey: ['npc-instances', instanceId],
    queryFn: () => apiGet<NPCInstance>(`${API}/api/npc-instances/${instanceId}`),
  });

  // ── Mutations ──────────────────────────────────────────────────────────────

  const updateMutation = useMutation({
    mutationFn: (body: Record<string, unknown>) =>
      apiPut(`${API}/api/npc-instances/${instanceId}`, body),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['npc-instances', instanceId] });
      queryClient.invalidateQueries({ queryKey: ['npc-instances'] });
      setIsEditing(false);
    },
  });

  const deleteMutation = useMutation({
    mutationFn: () => apiDelete(`${API}/api/npc-instances/${instanceId}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['npc-instances'] });
      navigate({ to: '/npcs/$npcId', params: { npcId } });
    },
  });

  // ── Handlers ───────────────────────────────────────────────────────────────

  const startEditing = () => {
    if (!instance) return;
    setForm({
      room_id: instance.room_id,
      starting_room_id: instance.starting_room_id,
      instance_number: instance.instance_number,
      hitpoints: instance.hitpoints,
    });
    setIsEditing(true);
  };

  const handleSave = () => {
    if (!form) return;
    updateMutation.mutate({
      room_id: form.room_id,
      starting_room_id: form.starting_room_id,
      instance_number: form.instance_number,
      hitpoints: form.hitpoints,
    });
  };

  const handleDelete = () => {
    deleteMutation.mutate();
  };

  // ── Loading / Error ────────────────────────────────────────────────────────

  if (isLoading) {
    return (
      <div className="p-8">
        <PageHeader title="Loading..." backTo={`/npcs/${npcId}`} />
        <div className="text-text-muted">Loading NPC instance...</div>
      </div>
    );
  }

  if (error || !instance) {
    return (
      <div className="p-8">
        <PageHeader title="Error" backTo={`/npcs/${npcId}`} />
        <div className="text-danger">
          Failed to load instance: {error?.message ?? 'Unknown error'}
        </div>
      </div>
    );
  }

  // ── Render ─────────────────────────────────────────────────────────────────

  return (
    <div className="p-8">
      <PageHeader
        title={instance.name}
        backTo={`/npcs/${npcId}`}
        actions={
          <div className="flex items-center gap-2">
            {!isEditing ? (
              <Button variant="primary" size="sm" onClick={startEditing}>
                Edit
              </Button>
            ) : (
              <Button variant="secondary" size="sm" onClick={() => setIsEditing(false)}>
                Cancel
              </Button>
            )}
            <Button
              variant="danger"
              size="sm"
              onClick={() => setShowDeleteModal(true)}
              disabled={deleteMutation.isPending}
            >
              Delete
            </Button>
          </div>
        }
      />

      <div className="max-w-2xl">
        <div className="bg-surface-muted rounded-lg p-6 border border-border mb-6">
          {!isEditing ? (
            <>
              <h2 className="mt-0 mb-4 text-text text-lg font-semibold">Instance Details</h2>
              <div className="grid grid-cols-2 gap-x-6 gap-y-3">
                <DetailField label="ID" value={String(instance.id)} />
                <DetailField label="Name" value={instance.name} />
                <DetailField label="Race" value={instance.race} />
                <DetailField label="Level" value={String(instance.level)} />
                <DetailField label="Instance Number" value={String(instance.instance_number)} />
                <DetailField label="Room ID" value={String(instance.room_id)} />
                <DetailField label="Starting Room" value={String(instance.starting_room_id)} />
                <DetailField label="Hitpoints" value={`${instance.hitpoints} / ${instance.max_hitpoints}`} />
              </div>
            </>
          ) : (
            <div className="space-y-4">
              <h2 className="mt-0 mb-4 text-text text-lg font-semibold">Edit Instance</h2>
              <div className="grid grid-cols-2 gap-4">
                <div className="flex flex-col">
                  <label className="text-text-muted text-xs block mb-1">Room ID</label>
                  <input
                    type="number"
                    value={form?.room_id ?? 0}
                    onChange={(e) => setForm(p => p ? ({ ...p, room_id: parseInt(e.target.value) || 0 }) : null)}
                    min={1}
                    className="p-2 bg-surface border border-border rounded text-text text-sm"
                  />
                </div>
                <div className="flex flex-col">
                  <label className="text-text-muted text-xs block mb-1">Starting Room</label>
                  <input
                    type="number"
                    value={form?.starting_room_id ?? 0}
                    onChange={(e) => setForm(p => p ? ({ ...p, starting_room_id: parseInt(e.target.value) || 0 }) : null)}
                    min={1}
                    className="p-2 bg-surface border border-border rounded text-text text-sm"
                  />
                </div>
                <div className="flex flex-col">
                  <label className="text-text-muted text-xs block mb-1">Instance Number</label>
                  <input
                    type="number"
                    value={form?.instance_number ?? 1}
                    onChange={(e) => setForm(p => p ? ({ ...p, instance_number: parseInt(e.target.value) || 1 }) : null)}
                    min={1}
                    className="p-2 bg-surface border border-border rounded text-text text-sm"
                  />
                </div>
                <div className="flex flex-col">
                  <label className="text-text-muted text-xs block mb-1">Hitpoints</label>
                  <input
                    type="number"
                    value={form?.hitpoints ?? 0}
                    onChange={(e) => setForm(p => p ? ({ ...p, hitpoints: parseInt(e.target.value) || 0 }) : null)}
                    min={0}
                    className="p-2 bg-surface border border-border rounded text-text text-sm"
                  />
                </div>
              </div>
              
              {updateMutation.isError && (
                <div className="p-2 bg-danger/10 text-danger rounded text-xs">
                  Failed to save: {(updateMutation.error as Error)?.message}
                </div>
              )}

              <div className="flex gap-2 pt-2">
                <Button
                  variant="primary"
                  onClick={handleSave}
                  disabled={updateMutation.isPending}
                >
                  {updateMutation.isPending ? 'Saving...' : 'Save'}
                </Button>
                <Button variant="secondary" onClick={() => setIsEditing(false)}>
                  Cancel
                </Button>
              </div>
            </div>
          )}
        </div>
      </div>

      <DeleteConfirmation
        open={showDeleteModal}
        title="Delete NPC Instance"
        message="Are you sure? This will permanently delete this NPC instance."
        onConfirm={handleDelete}
        onCancel={() => setShowDeleteModal(false)}
        isLoading={deleteMutation.isPending}
      />
    </div>
  );
}

function DetailField({ label, value }: Readonly<{ label: string; value: string }>) {
  return (
    <div>
      <span className="text-text-muted text-xs block mb-0.5">{label}</span>
      <span className="text-text text-sm font-medium">{value}</span>
    </div>
  );
}
