/* eslint-disable functional/prefer-immutable-types */
import { useState } from "react";
import { useCharacterBiography, type BiographyEntry } from "../hooks/useCharacters";

const EVENT_ICONS: Record<string, string> = {
  boss_killed: "⚔️",
  npc_defeated: "🗡️",
  character_died: "💀",
  character_respawned: "🔄",
  level_up: "⭐",
  reclass: "🎭",
  rerace: "🧬",
  quest_complete: "📜",
  skill_learned: "📚",
  "skill.leveled_up": "📈",
  "skill.ability_unlocked": "🔓",
  "xp.gained": "✨",
  "skill_xp.gained": "📖",
};

const FILTER_OPTIONS = [
  { label: "All", value: "" },
  { label: "⚔️ Boss Kills", value: "boss_killed" },
  { label: "🗡️ Kills", value: "npc_defeated" },
  { label: "💀 Deaths", value: "character_died" },
  { label: "⭐ Level-Ups", value: "level_up" },
  { label: "🎭 Class Changes", value: "reclass" },
  { label: "📜 Quests", value: "quest_complete" },
];

export function CharacterBiography({ characterId }: Readonly<{ characterId: number }>) {
  const [filter, setFilter] = useState("");
  const { data, isLoading, isError, error } = useCharacterBiography(characterId, filter || undefined);

  if (isLoading) return <div className="text-text-muted text-sm">Loading biography...</div>;
  if (isError) return <div className="text-danger text-sm">Error: {error?.message ?? "Failed to load"}</div>;
  if (!data) return null;

  const summary = data.summary;

  return (
    <div className="space-y-3">
      {/* Summary Stats */}
      <div className="flex flex-wrap gap-3 text-xs">
        <span className="badge badge-primary">📊 {summary.total_events} events</span>
        {summary.boss_kills > 0 && (
          <span className="badge badge-warning">⚔️ {summary.boss_kills} boss kills</span>
        )}
        {summary.npc_kills > 0 && (
          <span className="badge badge-info">🗡️ {summary.npc_kills} kills</span>
        )}
        {summary.deaths > 0 && (
          <span className="badge badge-danger">💀 {summary.deaths} deaths</span>
        )}
        {summary.level_ups > 0 && (
          <span className="badge badge-success">⭐ {summary.level_ups} level-ups</span>
        )}
        {summary.class_changes > 0 && (
          <span className="badge badge-neutral">🎭 {summary.class_changes} class changes</span>
        )}
        {summary.quest_completions > 0 && (
          <span className="badge badge-info">📜 {summary.quest_completions} quests</span>
        )}
        {summary.skills_learned > 0 && (
          <span className="badge badge-neutral">📚 {summary.skills_learned} skills</span>
        )}
      </div>

      {/* Filter Buttons */}
      <div className="flex flex-wrap gap-1">
        {FILTER_OPTIONS.map((opt) => (
          <button
            key={opt.value}
            onClick={() => setFilter(opt.value)}
            className={`px-2 py-1 rounded text-xs border transition-colors ${
              filter === opt.value
                ? "bg-primary text-white border-primary"
                : "bg-surface text-text-muted border-border hover:bg-surface-muted"
            }`}
          >
            {opt.label}
          </button>
        ))}
      </div>

      {/* Timeline */}
      <div className="space-y-1 max-h-[400px] overflow-y-auto">
        {data.entries.length === 0 ? (
          <p className="text-text-muted text-sm italic">No biography entries found.</p>
        ) : (
          data.entries.map((entry: BiographyEntry) => {
            const icon = EVENT_ICONS[entry.type] ?? "📋";
            const date = new Date(entry.timestamp);
            const dateStr = date.toLocaleDateString("en-US", { month: "short", day: "numeric" });
            const timeStr = date.toLocaleTimeString("en-US", { hour: "2-digit", minute: "2-digit" });

            return (
              <div
                key={entry.id}
                className="flex items-start gap-2 p-2 rounded border border-border bg-surface hover:bg-surface-muted transition-colors"
              >
                <span className="text-lg shrink-0">{icon}</span>
                <div className="flex-1 min-w-0">
                  <p className="text-sm text-text m-0">{entry.description}</p>
                  <p className="text-xs text-text-muted m-0 mt-0.5">
                    {dateStr} {timeStr} · <span className="font-mono">{entry.type}</span>
                  </p>
                </div>
              </div>
            );
          })
        )}
      </div>
    </div>
  );
}