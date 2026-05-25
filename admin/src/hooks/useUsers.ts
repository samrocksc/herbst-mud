 
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiGet, apiPost, apiDelete } from "../utils/apiFetch";

const API = `${window.location.origin}`;

export type User = Readonly<{
  id: number
  email: string
  is_admin: boolean
  created_at: string
  character_name?: string
}>

export function useUsers() {
  return useQuery({
    queryKey: ["users"],
    queryFn: () => apiGet<User[]>(`${API}/users`),
  });
}

export function useUser(id: number | null) {
  return useQuery({
    queryKey: ["user", id],
    queryFn: () => (id ? apiGet<User>(`${API}/users/${id}`) : null),
    enabled: !!id,
  });
}

export function useResetPassword() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: number) =>
      apiPost<{ message: string }>(`${API}/users/${id}/reset-password`, {}),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["users"] }),
  });
}

export function useDeleteCharacter() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: number) => apiDelete<void>(`${API}/characters/${id}`),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["characters"] });
      qc.invalidateQueries({ queryKey: ["user-characters"] });
    },
  });
}