export type GameConfig = Readonly<{ id: number; key: string; value: string }>

const KEY_LABELS: Readonly<Record<string, string>> = {
  xp_thresholds: "XP Thresholds", death_penalty_percent: "Death Penalty (%)", starting_room_id: "Starting Room ID",
  max_level: "Max Level", campaign_type: "Campaign Type", corpse_rot_minutes: "Corpse Rot Minutes",
  xp_per_kill: "XP Per Kill", respawn_room_id: "Respawn Room ID", max_inventory_size: "Max Inventory Size",
  pvp_enabled: "PvP Enabled", death_penalty_type: "Death Penalty Type", starting_hp: "Starting HP",
  starting_mana: "Starting Mana", starting_gold: "Starting Gold", regen_tick_seconds: "Regen Tick (seconds)",
  fountain_room_id: "Fountain Room ID",
};

export function humanizeKey(raw: string): string {
  if (KEY_LABELS[raw]) return KEY_LABELS[raw];
  return raw.split("_").map(s => s.charAt(0).toUpperCase() + s.slice(1)).join(" ");
}

export function tryParseJSON(value: string): unknown | null {
  try { const p = JSON.parse(value); return typeof p === "object" && p !== null ? p : null; }
  catch { return null; }
}

export const ROOM_ID_KEYS = ["fountain_room_id", "starting_room_id", "respawn_room_id"] as const;
export function isRoomIdKey(key: string): boolean {
  return (ROOM_ID_KEYS as ReadonlyArray<string>).includes(key);
}

export const PRESETS = [
  { label: "XP Thresholds", key: "xp_thresholds", value: "{\"1\":100,\"2\":300,\"3\":600,\"4\":1000,\"5\":1500}" },
  { label: "Death Penalty %", key: "death_penalty_percent", value: "10" },
  { label: "Corpse Rot Minutes", key: "corpse_rot_minutes", value: "5" },
  { label: "XP Per Kill", key: "xp_per_kill", value: "50" },
  { label: "Max Level", key: "max_level", value: "100" },
  { label: "Starting Room ID", key: "starting_room_id", value: "1" },
];
