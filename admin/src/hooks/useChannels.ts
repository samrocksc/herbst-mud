import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { apiGet, apiPut } from '../utils/apiFetch';

const API = `${window.location.origin}/api`;

export type ChannelConfig = Readonly<{
  name: string
  description: string
  color: string
  default_enabled: boolean
  cooldown_seconds: number
  admin_only: boolean
}>

export type ChannelInput = Readonly<{
  name: string
  description: string
  color: string
  default_enabled: boolean
  cooldown_seconds: number
  admin_only: boolean
}>

export function useChannelConfigs() {
  return useQuery({
    queryKey: ['channels'],
    queryFn: () => apiGet<ChannelConfig[]>(`${API}/channels`),
  });
}

export function useUpdateChannel() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ name, input }: { name: string; input: ChannelInput }) =>
      apiPut<ChannelConfig>(`${API}/channels/${name}`, input),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['channels'] }),
  });
}
