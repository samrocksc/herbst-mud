import { useQuery } from '@tanstack/react-query'
import { apiGet } from '../utils/apiFetch'
import type { NPC } from '../components/map/types'

const API = `${window.location.origin}`

export function useNPCs() {
  return useQuery<NPC[]>({
    queryKey: ['npcs'],
    queryFn: async () => {
      const data = await apiGet<{ npcs: NPC[] }>(`${API}/npcs`)
      return data.npcs || (data as unknown as NPC[])
    },
  })
}