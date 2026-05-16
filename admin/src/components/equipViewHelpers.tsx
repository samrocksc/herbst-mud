/* eslint-disable react-refresh/only-export-components */
/** Shared helper functions for equipment view components. */

/** Format a slot name like "main_hand" into "Main Hand". */
export function formatSlotName(slot: string): string {
  return slot.replace(/_/g, " ").replace(/\b\w/g, (c) => c.toUpperCase());
}

/** Rarity badge with color coding. */
export function RarityBadge({ rarity }: Readonly<{ rarity: string }>) {
  const colors: Record<string, string> = {
    common: "bg-gray-700/30 text-gray-300 border-gray-600/40",
    uncommon: "bg-green-900/30 text-green-400 border-green-700/40",
    rare: "bg-blue-900/30 text-blue-400 border-blue-700/40",
    epic: "bg-purple-900/30 text-purple-400 border-purple-700/40",
    legendary: "bg-orange-900/30 text-orange-400 border-orange-700/40",
  };
  const colorClass = colors[rarity] || colors.common;
  return (
    <span className={`inline-block px-1.5 py-0.5 rounded text-xs font-medium border ${colorClass}`}>
      {rarity}
    </span>
  );
}
