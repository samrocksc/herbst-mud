/* eslint-disable functional/prefer-immutable-types, functional/immutable-data, functional/no-return-void */
import { useState, useCallback } from "react";
import { apiGet, apiPut, apiDelete } from "../../utils/apiFetch";
import { showToast } from "../Toast";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import type { ItemInstanceView, EditFormData } from "./types";

const API = `${window.location.origin}/api`;

export function useItemInstances(roomId: number) {
  const queryClient = useQueryClient();
  const instancesQuery = useQuery({
    queryKey: ["item-instances", roomId],
    queryFn: () => apiGet<ItemInstanceView[]>(`${API}/item-instances?roomId=${roomId}`),
  });
  const updateMutation = useMutation({
    mutationFn: (args: { id: number; update: Record<string, unknown> }) =>
      apiPut<ItemInstanceView>(`${API}/item-instances/${args.id}`, args.update),
    onSuccess: () => { queryClient.invalidateQueries({ queryKey: ["item-instances"] }); },
  });
  const deleteMutation = useMutation({
    mutationFn: (id: number) => apiDelete(`${API}/item-instances/${id}`),
    onSuccess: () => { queryClient.invalidateQueries({ queryKey: ["item-instances"] }); },
  });

  const [editingId, setEditingId] = useState<number | null>(null);
  const [confirmDeleteId, setConfirmDeleteId] = useState<number | null>(null);
  const [editForm, setEditForm] = useState<Partial<EditFormData>>({});

  const handleUpdate = useCallback(async () => {
    if (editingId === null) return;
    try {
      const u: Record<string, unknown> = {};
      if (editForm.name !== undefined) u.name = editForm.name;
      if (editForm.description !== undefined) u.description = editForm.description;
      if (editForm.slot !== undefined) u.slot = editForm.slot;
      if (editForm.level !== undefined) u.level = editForm.level;
      if (editForm.weight !== undefined) u.weight = editForm.weight;
      if (editForm.color !== undefined) u.color = editForm.color;
      await updateMutation.mutateAsync({ id: editingId, update: u });
      setEditingId(null); setEditForm({});
    } catch (err) { showToast(`Update failed: ${(err as Error)?.message ?? "Unknown error"}`); }
  }, [editingId, editForm, updateMutation]);

  const handleDelete = useCallback(async (id: number) => {
    try {
      await deleteMutation.mutateAsync(id);
      setConfirmDeleteId(null);
      if (editingId === id) { setEditingId(null); setEditForm({}); }
    } catch (err) { showToast(`Delete failed: ${(err as Error)?.message ?? "Unknown error"}`); }
  }, [deleteMutation, editingId]);

  return { instancesQuery, updateMutation, deleteMutation, editingId, setEditingId, confirmDeleteId, setConfirmDeleteId, editForm, setEditForm, handleUpdate, handleDelete };
}
