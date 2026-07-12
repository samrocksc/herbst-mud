// useZones — TanStack Query hooks for the zone admin API.
//
// Zones are world-scoped. The current world is read from WorldStoreContext;
// components do not pass world_id explicitly. This matches the standing
// herbst-mud rule: forms never expose world_id as a dropdown.
//
// Uses the shared apiFetch helpers from admin/src/utils/apiFetch.ts which
// auto-inject the Authorization Bearer token from localStorage. The
// helpers also strip known response envelopes (e.g. { "zones": [...] }).

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useWorldStore } from "../contexts/WorldStoreContext";
import { API_BASE, apiGet, apiPost, apiPut, apiDelete } from "../utils/apiFetch";

export type Zone = Readonly<{
  id: string;
  world_id: string;
  name: string;
  description?: string;
  min_level: number;
  parent_zone_id?: string;
  color?: string;
  room_ids?: number[];
}>;

export type ZoneRoom = Readonly<{
  id: number;
  name: string;
  exists: boolean;
  message?: string;
}>;

export type ZoneInput = Readonly<{
  id: string;
  name: string;
  description?: string;
  min_level?: number;
  parent_zone_id?: string;
  color?: string;
  room_ids?: number[];
}>;

export type ZoneUpdate = Partial<Omit<ZoneInput, "id">>;

function zoneUrl(path: string): string {
  return API_BASE + path;
}

function zoneRoomsUrl(zoneId: string): string {
  return API_BASE + "/api/zones/" + encodeURIComponent(zoneId) + "/rooms";
}

function zoneUrlById(zoneId: string): string {
  return API_BASE + "/api/zones/" + encodeURIComponent(zoneId);
}

export function useZones() {
  const { currentWorld } = useWorldStore();
  const queryClient = useQueryClient();

  const zonesQuery = useQuery<Zone[]>({
    queryKey: ["zones", currentWorld],
    queryFn: async () => {
      const qs = currentWorld ? "?world_id=" + currentWorld : "";
      // apiGet auto-unwraps { "zones": [...] } envelope.
      return apiGet<Zone[]>(zoneUrl("/api/zones" + qs));
    },
    enabled: !!currentWorld,
  });

  const createMutation = useMutation({
    mutationFn: (input: ZoneInput) =>
      apiPost<Zone>(zoneUrl("/api/zones"), {
        ...input,
        world_id: currentWorld,
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["zones", currentWorld] });
    },
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, update }: { id: string; update: ZoneUpdate }) =>
      apiPut<Zone>(zoneUrlById(id), update),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["zones", currentWorld] });
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (id: string) =>
      apiDelete<{ message: string }>(zoneUrlById(id)),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["zones", currentWorld] });
    },
  });

  return {
    zones: zonesQuery.data ?? [],
    isLoading: zonesQuery.isLoading,
    createZone: createMutation.mutate,
    createZoneAsync: createMutation.mutateAsync,
    updateZone: updateMutation.mutate,
    updateZoneAsync: updateMutation.mutateAsync,
    deleteZone: deleteMutation.mutate,
    isCreating: createMutation.isPending,
    isUpdating: updateMutation.isPending,
    isDeleting: deleteMutation.isPending,
  };
}

// useZoneRooms fetches the room list for a single zone, including the
// `exists: false` ghost entries for rooms that have been removed from the
// world. The list is sorted by ID server-side.
export function useZoneRooms(zoneId: string | null) {
  const queryClient = useQueryClient();
  const roomsQuery = useQuery<{ zone_id: string; rooms: ZoneRoom[] }>({
    queryKey: ["zones", zoneId, "rooms"],
    queryFn: () => apiGet<{ zone_id: string; rooms: ZoneRoom[] }>(zoneRoomsUrl(zoneId ?? "")),
    enabled: !!zoneId,
  });

  const addRooms = useMutation({
    mutationFn: async (roomIds: number[]) => {
      // Fetch current, then PUT the union.
      const current = roomsQuery.data?.rooms ?? [];
      const existing = current.map((r) => r.id);
      const seen = new Set<number>(existing);
      const merged = [...existing];
      for (const rid of roomIds) {
        if (!seen.has(rid)) {
          merged.push(rid);
          seen.add(rid);
        }
      }
      return apiPut<Zone>(zoneUrlById(zoneId ?? ""), {
        room_ids: merged,
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["zones", zoneId, "rooms"] });
      queryClient.invalidateQueries({ queryKey: ["zones"] });
    },
  });

  const removeRoom = useMutation({
    mutationFn: async (roomId: number) => {
      const current = roomsQuery.data?.rooms ?? [];
      const filtered = current.map((r) => r.id).filter((id) => id !== roomId);
      return apiPut<Zone>(zoneUrlById(zoneId ?? ""), {
        room_ids: filtered,
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["zones", zoneId, "rooms"] });
      queryClient.invalidateQueries({ queryKey: ["zones"] });
    },
  });

  return {
    rooms: roomsQuery.data?.rooms ?? [],
    isLoading: roomsQuery.isLoading,
    addRooms: addRooms.mutateAsync,
    isAdding: addRooms.isPending,
    removeRoom: removeRoom.mutateAsync,
    isRemoving: removeRoom.isPending,
  };
}
