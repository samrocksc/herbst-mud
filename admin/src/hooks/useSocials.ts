/* eslint-disable functional/prefer-immutable-types */
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiGet, apiPost, apiPut, apiDelete } from "../utils/apiFetch";

const API = `${window.location.origin}/api`;

export type SocialCommand = Readonly<{
  id: number
  name: string
  self_text: string
  room_text: string
  target_self_text: string
  target_text: string
  target_room_text: string
}>

export type SocialInput = Readonly<{
  id?: number
  name: string
  self_text: string
  room_text: string
  target_self_text: string
  target_text: string
  target_room_text: string
}>

export function useSocials() {
  return useQuery({
    queryKey: ["socials"],
    queryFn: () => apiGet<SocialCommand[]>(`${API}/socials`),
  });
}

export function useSocial(id: number | null) {
  return useQuery({
    queryKey: ["social", id],
    queryFn: () => apiGet<SocialCommand>(`${API}/socials/${id}`),
    enabled: !!id,
  });
}

export function useCreateSocial() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (input: SocialInput) => apiPost<SocialCommand>(`${API}/socials`, input),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["socials"] }),
  });
}

export function useUpdateSocial() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ id, input }: { id: number; input: SocialInput }) =>
      apiPut<SocialCommand>(`${API}/socials/${id}`, input),
    onSuccess: (_, { id }) => {
      qc.invalidateQueries({ queryKey: ["socials"] });
      qc.invalidateQueries({ queryKey: ["social", id] });
    },
  });
}

export function useDeleteSocial() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: number) => apiDelete(`${API}/socials/${id}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["socials"] }),
  });
}
