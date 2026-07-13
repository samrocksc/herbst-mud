/* eslint-disable functional/prefer-immutable-types */
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiGet, apiPost } from "../utils/apiFetch";

const API = `${window.location.origin}`;

/** A single class-history record for a character. */
export type ClassHistoryEntry = Readonly<{
  id: number
  character_id: number
  faction_id: number
  faction_name: string
  joined_at: string
  left_at: string | null
  reason: string
}>

/** A single race-history record for a character. */
export type RaceHistoryEntry = Readonly<{
  id: number
  character_id: number
  race_id: number | null
  race_name: string
  changed_at: string
  reason: string
}>

/** Result returned by reclass/rerace endpoints. */
export type ReclassReraceResult = Readonly<{
  character_id: number
  new_class?: string
  new_race?: string
  previous_class?: string
  previous_race?: string
  history_recorded: boolean
}>

/**
 * Fetch a character's class (faction) history.
 * GET /characters/:id/class-history
 */
export function useClassHistory(characterId: number) {
  return useQuery({
    queryKey: ["character-class-history", characterId],
    queryFn: async (): Promise<ClassHistoryEntry[]> => {
      const data = await apiGet<ClassHistoryEntry[]>(
        `${API}/characters/${characterId}/class-history`,
      );
      return Array.isArray(data) ? data : [];
    },
    enabled: !!characterId,
  });
}

/**
 * Fetch a character's race history.
 * GET /characters/:id/race-history
 */
export function useRaceHistory(characterId: number) {
  return useQuery({
    queryKey: ["character-race-history", characterId],
    queryFn: async (): Promise<RaceHistoryEntry[]> => {
      const data = await apiGet<RaceHistoryEntry[]>(
        `${API}/characters/${characterId}/race-history`,
      );
      return Array.isArray(data) ? data : [];
    },
    enabled: !!characterId,
  });
}

/**
 * Reclass a character — changes their class (faction).
 * POST /characters/:id/reclass  body: { faction_id: number, reason?: string }
 */
export function useReclass() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({
      characterId,
      factionId,
      reason,
    }: {
      characterId: number
      factionId: number
      reason?: string
    }) =>
      apiPost<ReclassReraceResult>(
        `${API}/characters/${characterId}/reclass`,
        { faction_id: factionId, reason },
      ),
    onSuccess: (_data, vars) => {
      qc.invalidateQueries({ queryKey: ["character", vars.characterId] });
      qc.invalidateQueries({ queryKey: ["characters"] });
      qc.invalidateQueries({ queryKey: ["character-class-history", vars.characterId] });
    },
  });
}

/**
 * Rerace a character — changes their race.
 * POST /characters/:id/rerace  body: { race_id: number, reason?: string }
 */
export function useRerace() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({
      characterId,
      raceId,
      reason,
    }: {
      characterId: number
      raceId: number
      reason?: string
    }) =>
      apiPost<ReclassReraceResult>(
        `${API}/characters/${characterId}/rerace`,
        { race_id: raceId, reason },
      ),
    onSuccess: (_data, vars) => {
      qc.invalidateQueries({ queryKey: ["character", vars.characterId] });
      qc.invalidateQueries({ queryKey: ["characters"] });
      qc.invalidateQueries({ queryKey: ["character-race-history", vars.characterId] });
    },
  });
}