 
import { useQuery } from "@tanstack/react-query";
import { apiGet } from "../utils/apiFetch";
import type { Equipment } from "../components/map/types";

const API = `${window.location.origin}`;

export function useRoomEquipment(roomId: number | null) {
  return useQuery<Equipment[]>({
    queryKey: ["room-equipment", roomId],
    enabled: !!roomId,
    queryFn: () => apiGet<Equipment[]>(`${API}/rooms/${roomId}/equipment`),
  });
}