/* eslint-disable functional/prefer-immutable-types */
import { useQuery } from "@tanstack/react-query";
import { useWorldStore } from "../contexts/WorldStoreContext";
import { apiGet } from "../utils/apiFetch";

type Equipment = Readonly<{
  id: number;
  name: string;
  description: string;
  slot: string;
  level: number;
  weight: number;
  isEquipped: boolean;
  isImmovable: boolean;
  color: string;
  isVisible: boolean;
  itemType: string;
  ownerID: number | null;
  equipmentTemplateID: number | null;
  healing: number;
  effect: string;
  armorRating: number;
  armorType: string;
  stats: Record<string, number>;
  rarity: string;
  skillRequirement: string;
  skillRequirementLevel: number;
  damageDiceCount: number;
  damageDiceSides: number;
  damageBonus: number;
  damageType: string;
  weaponType: string;
  isTwoHanded: boolean;
}>;

export function useSearchEquipment(query: string) {
  const { currentWorld } = useWorldStore();

  return useQuery({
    queryKey: ["search-equipment", query, currentWorld],
    queryFn: async (): Promise<Equipment[]> => {
      if (!query || query.length < 2) return [];
      const params = new URLSearchParams({
        search: query,
        ...(currentWorld ? { world_id: currentWorld } : {}),
      });
      const url = `${window.location.origin}/api/equipment/search?${params.toString()}`;
      const data = await apiGet<Equipment[]>(url);
      return Array.isArray(data) ? data : [];
    },
    enabled: query.length >= 2,
  });
}
