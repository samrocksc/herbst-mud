import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { apiGet, apiDelete } from '../utils/apiFetch';

const API = `${window.location.origin}/api`;

export type TellQueueEntry = Readonly<{
  id: number
  sender_id: number
  sender_name: string
  recipient_id: number
  recipient_name: string
  message: string
  sent_at: string
  delivered_at: string | null
  expires_at: string
}>

export type TellQueueFilters = {
  recipient_id?: number
  undelivered?: boolean
  limit?: number
}

export function useTellQueue(filters?: TellQueueFilters) {
  const params = new URLSearchParams();
  if (filters?.recipient_id) params.append('recipient_id', String(filters.recipient_id));
  if (filters?.undelivered) params.append('undelivered', 'true');
  if (filters?.limit) params.append('limit', String(filters.limit));
  const qs = params.toString();
  return useQuery({
    queryKey: ['tell-queue', filters],
    queryFn: () => apiGet<TellQueueEntry[]>(`${API}/tell-queue${qs ? '?' + qs : ''}`),
  });
}

export function useDeleteTell() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: number) => apiDelete(`${API}/tell-queue/${id}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['tell-queue'] }),
  });
}
