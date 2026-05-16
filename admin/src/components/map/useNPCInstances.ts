/* eslint-disable functional/prefer-immutable-types, functional/immutable-data, functional/no-return-void */
import { useState, useCallback } from "react";
import { apiGet, apiPut, apiDelete } from "../../utils/apiFetch";
import { showToast } from "../Toast";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import type { NPCInstanceView, EditFormData } from "./NPCInstanceManager";

export function useNPCInstances(roomId: number) {
  const queryClient = useQueryClient();
  const [editingId, setEditingId] = useState<number | null>(null);
  const [confirmDeleteId, setConfirmDeleteId] = useState<number | null>(null);
  const [editForm, setEditForm] = useState<EditFormData>({ level: 0, hitpoints: 0, room_id: roomId, starting_room_id: roomId });

  const instancesQuery = useQuery({
    queryKey: ["npc-instances", roomId],
    queryFn: async (): Promise<NPCInstanceView[]> =>
      apiGet<NPCInstanceView[]>(`${window.location.origin}/api/npc-instances?roomId=${roomId}`),
  });

  const updateMutation = useMutation({
    mutationFn: async (args: { id: number; update: Record<string, unknown> }): Promise<NPCInstanceView> =>
      apiPut<NPCInstanceView>(`${window.location.origin}/api/npc-instances/${args.id}`, args.update),
    onSuccess: () => { queryClient.invalidateQueries({ queryKey: ["npc-instances"] }); },
  });

  const deleteMutation = useMutation({
    mutationFn: async (id: number): Promise<void> => { await apiDelete(`${window.location.origin}/api/npc-instances/${id}`); },
    onSuccess: () => { queryClient.invalidateQueries({ queryKey: ["npc-instances"] }); showToast("NPC instance deleted", "success"); },
  });

  const handleUpdate = useCallback(async () => {
    if (editingId === null) return;
    try {
      const update: Record<string, unknown> = { room_id: editForm.room_id, starting_room_id: editForm.starting_room_id };
      if (editForm.level > 0) update.level = editForm.level;
      if (editForm.hitpoints > 0) update.hitpoints = editForm.hitpoints;
      await updateMutation.mutateAsync({ id: editingId, update });
      showToast("NPC updated", "success");
      setEditingId(null);
      setEditForm({ level: 0, hitpoints: 0, room_id: roomId, starting_room_id: roomId });
    } catch (err) { showToast(`Update failed: ${(err as Error)?.message ?? "Unknown error"}`, "error"); }
  }, [editingId, editForm, updateMutation, roomId]);

  const handleDelete = useCallback(async (id: number) => {
    try { await deleteMutation.mutateAsync(id); setConfirmDeleteId(null); }
    catch (err) { showToast(`Delete failed: ${(err as Error)?.message ?? "Unknown error"}`, "error"); }
  }, [deleteMutation]);

  const startEdit = useCallback((inst: NPCInstanceView) => {
    setEditingId(inst.id); setConfirmDeleteId(null);
    setEditForm({ level: inst.level, hitpoints: inst.hitpoints, room_id: inst.room_id, starting_room_id: inst.room_id });
  }, []);

  return {
    instancesQuery, updateMutation, deleteMutation,
    editingId, setEditingId, confirmDeleteId, setConfirmDeleteId,
    editForm, setEditForm,
    handleUpdate, handleDelete, startEdit,
  };
}
