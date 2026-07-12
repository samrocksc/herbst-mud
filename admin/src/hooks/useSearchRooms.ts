/* eslint-disable functional/prefer-immutable-types */
import { useQuery } from "@tanstack/react-query";
import { useWorldStore } from "../contexts/WorldStoreContext";
import { apiGet } from "../utils/apiFetch";

type Room = Readonly<{
  id: number;
  name: string;
  description: string;
  isStartingRoom: boolean;
  isRootRoom: boolean;
  exits: Record<string, number>;
  posZ: number;
  version: number;
}>;

export function useSearchRooms(query: string) {
  const { currentWorld } = useWorldStore();

  return useQuery({
    queryKey: ["search-rooms", query, currentWorld],
    queryFn: async (): Promise<Room[]> => {
      const params = new URLSearchParams({
        ...(query.length >= 2 ? { search: query } : {}),
        ...(currentWorld ? { world_id: currentWorld } : {}),
      });
      const url = `${window.location.origin}/api/rooms?${params.toString()}`;
      const data = await apiGet<Room[]>(url);
      const results = Array.isArray(data) ? data : [];
      // Return up to 4 rooms for initial guidance, or filtered results
      return query.length >= 2 ? results : results.slice(0, 4);
    },
    enabled: query.length >= 2 || query.length === 0,
  });
}
