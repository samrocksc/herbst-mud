import { useQuery } from '@tanstack/react-query'
import { apiGet } from '../utils/apiFetch'

const API = `${window.location.origin}`

export type LogEntry = Readonly<{
  id: number
  level: string
  message: string
  service?: string
  character_id?: number
  room_id?: number
  template_id?: string
  metadata?: Record<string, unknown>
  created_at: string
}>

export type LogFilters = {
  level?: string
  service?: string
  character_id?: string
  room_id?: string
  template_id?: string
  limit?: number
  offset?: number
}

export function useLogs(filters?: LogFilters) {
  const params = new URLSearchParams()
  if (filters?.level) params.append('level', filters.level)
  if (filters?.service) params.append('service', filters.service)
  if (filters?.character_id) params.append('character_id', filters.character_id)
  if (filters?.room_id) params.append('room_id', filters.room_id)
  if (filters?.template_id) params.append('template_id', filters.template_id)
  if (filters?.limit) params.append('limit', String(filters.limit))
  if (filters?.offset) params.append('offset', String(filters.offset))
  const qs = params.toString()
  const url = `${API}/api/logs${qs ? '?' + qs : ''}`

  return useQuery({
    queryKey: ['logs', filters],
    queryFn: async () => {
      const data = await apiGet<{ logs: LogEntry[]; total: number }>(url)
      return { logs: data.logs ?? [], total: data.total ?? 0 }
    },
  })
}

export function useLogServices() {
  return useQuery({
    queryKey: ['log-services'],
    queryFn: async () => {
      const data = await apiGet<{ services: string[] }>(`${API}/api/logs/services`)
      return data.services ?? []
    },
  })
}