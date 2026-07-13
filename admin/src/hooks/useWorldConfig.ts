/* eslint-disable functional/prefer-immutable-types, functional/immutable-data */
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiGet, apiPut } from "../utils/apiFetch";
import { useWorldStore } from "../contexts/WorldStoreContext";

const API = `${window.location.origin}/api`;

export type LevelCurve = Readonly<{
  mode: "percentage" | "hand_coded";
  base_xp: number;
  percentage: number;
  max_level: number;
  thresholds?: ReadonlyArray<number>;
}>;

export type StatGrowth = Readonly<{
  hp_per_level: number;
  mana_per_level: number;
  stamina_per_level: number;
}>;

export type SkillXpConfig = Readonly<{
  usage_diminishing_returns: boolean;
  usage_cap_per_hour: number;
  anti_grind_kill_threshold: number;
}>;

export type ReclassConfig = Readonly<{
  allowed: boolean;
  cost: number;
  min_level: number;
  cooldown_seconds: number;
  skill_retention: number;
}>;

export type ReraceConfig = Readonly<{
  allowed: boolean;
  cost: number;
}>;

export type WorldConfig = Readonly<{
  level_curve: LevelCurve;
  stat_growth: StatGrowth;
  skill_xp: SkillXpConfig;
  reclass: ReclassConfig;
  rerace: ReraceConfig;
}>;

export const DEFAULT_CONFIG: WorldConfig = {
  level_curve: {
    mode: "percentage",
    base_xp: 1000,
    percentage: 50,
    max_level: 50,
  },
  stat_growth: {
    hp_per_level: 10,
    mana_per_level: 5,
    stamina_per_level: 5,
  },
  skill_xp: {
    usage_diminishing_returns: true,
    usage_cap_per_hour: 100,
    anti_grind_kill_threshold: 20,
  },
  reclass: {
    allowed: true,
    cost: 1000,
    min_level: 10,
    cooldown_seconds: 3600,
    skill_retention: 0.5,
  },
  rerace: {
    allowed: false,
    cost: 5000,
  },
};

export function useWorldConfig() {
  const { currentWorld } = useWorldStore();
  const worldId = currentWorld && currentWorld !== "default" ? currentWorld : "1";

  return useQuery({
    queryKey: ["world-config", worldId],
    queryFn: async (): Promise<WorldConfig | null> => {
      const data = await apiGet<Record<string, unknown>>(`${API}/worlds/${worldId}`);
      if (!data) return null;
      if (data.config) return data.config as WorldConfig;
      return null;
    },
  });
}

export function useUpdateWorldConfig() {
  const qc = useQueryClient();
  const { currentWorld } = useWorldStore();
  const worldId = currentWorld && currentWorld !== "default" ? currentWorld : "1";

  return useMutation({
    mutationFn: (config: WorldConfig) =>
      apiPut(`${API}/worlds/${worldId}`, { config }),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["world-config", worldId] });
      qc.invalidateQueries({ queryKey: ["worlds"] });
      qc.invalidateQueries({ queryKey: ["world", worldId] });
    },
  });
}