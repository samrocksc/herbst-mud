import { useState } from 'react'

export type GameConfig = Readonly<{ id: number; key: string; value: string }>
export type ConfigForm = Readonly<{ key: string; value: string }>

const KEY_LABELS: Readonly<Record<string, string>> = {
  xp_thresholds: 'XP Thresholds',
  death_penalty_percent: 'Death Penalty (%)',
  starting_room_id: 'Starting Room ID',
  max_level: 'Max Level',
  campaign_type: 'Campaign Type',
  corpse_rot_minutes: 'Corpse Rot Minutes',
  xp_per_kill: 'XP Per Kill',
  respawn_room_id: 'Respawn Room ID',
  max_inventory_size: 'Max Inventory Size',
  pvp_enabled: 'PvP Enabled',
  death_penalty_type: 'Death Penalty Type',
  starting_hp: 'Starting HP',
  starting_mana: 'Starting Mana',
  starting_gold: 'Starting Gold',
  regen_tick_seconds: 'Regen Tick (seconds)',
}

export function humanizeKey(raw: string): string {
  if (KEY_LABELS[raw]) return KEY_LABELS[raw]
  return raw.split('_').map(s => s.charAt(0).toUpperCase() + s.slice(1)).join(' ')
}

export function tryParseJSON(value: string): unknown | null {
  try {
    const parsed = JSON.parse(value)
    if (typeof parsed === 'object' && parsed !== null) return parsed
    return null
  } catch { return null }
}

export const PRESETS = [
  { label: 'XP Thresholds', key: 'xp_thresholds', value: '{"1":100,"2":300,"3":600,"4":1000,"5":1500}' },
  { label: 'Death Penalty %', key: 'death_penalty_percent', value: '10' },
  { label: 'Corpse Rot Minutes', key: 'corpse_rot_minutes', value: '5' },
  { label: 'XP Per Kill', key: 'xp_per_kill', value: '50' },
  { label: 'Max Level', key: 'max_level', value: '100' },
  { label: 'Starting Room ID', key: 'starting_room_id', value: '1' },
]

/** Renders a config value cell: pretty-prints JSON, truncates long plain strings */
export function ConfigValueCell({ value }: { value: string }) {
  const parsed = tryParseJSON(value)
  const [expanded, setExpanded] = useState(false)
  if (parsed !== null) {
    const formatted = JSON.stringify(parsed, null, 2)
    const isLong = formatted.split('\n').length > 4
    return (
      <div className="text-xs">
        <pre className={`font-mono text-text-secondary whitespace-pre-wrap m-0 ${!expanded ? 'max-h-16 overflow-hidden' : ''}`}>
          {formatted}
        </pre>
        {isLong && (
          <button type="button" className="text-primary text-xs mt-1 hover:underline cursor-pointer"
            onClick={() => setExpanded(e => !e)}>
            {expanded ? 'Show less' : 'Show more'}
          </button>
        )}
      </div>
    )
  }
  return (
    <span className="inline-block max-w-md overflow-hidden text-ellipsis whitespace-nowrap text-text-secondary text-xs">
      {value.length > 60 ? value.slice(0, 60) + '…' : value}
    </span>
  )
}