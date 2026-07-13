import { useQuery } from "@tanstack/react-query";
import { apiGet } from "../utils/apiFetch";

const API = `${window.location.origin}/api/event-logs`;

export type EventLog = Readonly<{
  id: number
  event_type: string
  character_id: number | null
  world_id: string | null
  payload: Record<string, unknown>
  created_at: string
}>

export type EventLogFilters = Readonly<{
  event_type?: string
  character_id?: number
  limit?: number
  offset?: number
}>

/**
 * Fetch game event logs with optional filters.
 * GET /api/event-logs?event_type=X&character_id=Y&limit=Z
 */
export function useEventLogs(filters?: EventLogFilters) {
  return useQuery({
    queryKey: ["event-logs", filters],
    queryFn: async (): Promise<EventLog[]> => {
      const params = new URLSearchParams();
      if (filters?.event_type) params.append("event_type", filters.event_type);
      if (filters?.character_id) params.append("character_id", String(filters.character_id));
      if (filters?.limit) params.append("limit", String(filters.limit));
      if (filters?.offset) params.append("offset", String(filters.offset));
      const qs = params.toString();
      const url = `${API}${qs ? "?" + qs : ""}`;
      const data = await apiGet<EventLog[] | { event_logs: EventLog[] }>(url);
      // Handle both array and wrapped-object response shapes
      if (Array.isArray(data)) return data;
      if (data && typeof data === "object" && "event_logs" in data) {
        return (data as { event_logs: EventLog[] }).event_logs ?? [];
      }
      return [];
    },
  });
}