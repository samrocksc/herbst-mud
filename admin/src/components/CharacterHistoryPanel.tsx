/* eslint-disable functional/prefer-immutable-types, react-hooks/purity */
import { useState } from "react";
import {
  useClassHistory,
  useRaceHistory,
  type ClassHistoryEntry,
  type RaceHistoryEntry,
} from "../hooks/useCharacterHistory";

/** Format an ISO timestamp for display; returns "—" for null/empty. */
function formatTime(ts: string | null | undefined): string {
  if (!ts) return "—";
  try {
    return new Date(ts).toLocaleString();
  } catch {
    return ts;
  }
}

/** Badge for "current" (no left_at) vs "past" class membership. */
function ClassStatusBadge({ entry }: { entry: ClassHistoryEntry }) {
  if (!entry.left_at) {
    return <span className="badge badge-success text-xs">Current</span>;
  }
  return <span className="badge badge-neutral text-xs">Past</span>;
}

/** Collapsible class-history table. */
function ClassHistoryTable({ characterId }: Readonly<{ characterId: number }>) {
  const { data: history, isLoading, isError, error } = useClassHistory(characterId);

  if (isLoading) return <div className="text-text-muted text-xs py-1">Loading class history…</div>;
  if (isError) return <div className="text-danger text-xs py-1">Error: {error?.message ?? "Unknown"}</div>;
  if (!history || history.length === 0) {
    return <div className="text-text-muted text-xs py-1">No class history recorded.</div>;
  }

  return (
    <div className="overflow-x-auto">
      <table className="w-full text-xs">
        <thead>
          <tr className="text-text-muted border-b border-border">
            <th className="text-left py-1 pr-3">Class</th>
            <th className="text-left py-1 pr-3">Joined</th>
            <th className="text-left py-1 pr-3">Left</th>
            <th className="text-left py-1 pr-3">Reason</th>
            <th className="text-left py-1">Status</th>
          </tr>
        </thead>
        <tbody>
          {history.map((entry) => (
            <tr key={entry.id} className="border-b border-border/50">
              <td className="py-1 pr-3 text-text font-medium">{entry.faction_name}</td>
              <td className="py-1 pr-3 text-text-muted">{formatTime(entry.joined_at)}</td>
              <td className="py-1 pr-3 text-text-muted">{formatTime(entry.left_at)}</td>
              <td className="py-1 pr-3 text-text-muted">{entry.reason || "—"}</td>
              <td className="py-1"><ClassStatusBadge entry={entry} /></td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

/** Collapsible race-history table. */
function RaceHistoryTable({ characterId }: Readonly<{ characterId: number }>) {
  const { data: history, isLoading, isError, error } = useRaceHistory(characterId);

  if (isLoading) return <div className="text-text-muted text-xs py-1">Loading race history…</div>;
  if (isError) return <div className="text-danger text-xs py-1">Error: {error?.message ?? "Unknown"}</div>;
  if (!history || history.length === 0) {
    return <div className="text-text-muted text-xs py-1">No race history recorded.</div>;
  }

  return (
    <div className="overflow-x-auto">
      <table className="w-full text-xs">
        <thead>
          <tr className="text-text-muted border-b border-border">
            <th className="text-left py-1 pr-3">Race</th>
            <th className="text-left py-1 pr-3">Changed At</th>
            <th className="text-left py-1">Reason</th>
          </tr>
        </thead>
        <tbody>
          {history.map((entry) => (
            <tr key={entry.id} className="border-b border-border/50">
              <td className="py-1 pr-3 text-text font-medium">{entry.race_name}</td>
              <td className="py-1 pr-3 text-text-muted">{formatTime(entry.changed_at)}</td>
              <td className="py-1 text-text-muted">{entry.reason || "—"}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

/**
 * Collapsible panel showing a character's class and race history.
 * Both sections start collapsed; click the header to expand.
 */
export function CharacterHistoryPanel({ characterId }: Readonly<{ characterId: number }>) {
  const [showClass, setShowClass] = useState(false);
  const [showRace, setShowRace] = useState(false);

  return (
    <div className="space-y-3">
      {/* Class history */}
      <div className="border border-border rounded-md">
        <button
          onClick={() => setShowClass(!showClass)}
          className="w-full flex items-center justify-between px-3 py-2 text-sm text-text hover:bg-surface-muted"
        >
          <span className="font-semibold">Class History</span>
          <span className="text-text-muted text-xs">{showClass ? "▼" : "▶"}</span>
        </button>
        {showClass && (
          <div className="px-3 pb-3">
            <ClassHistoryTable characterId={characterId} />
          </div>
        )}
      </div>

      {/* Race history */}
      <div className="border border-border rounded-md">
        <button
          onClick={() => setShowRace(!showRace)}
          className="w-full flex items-center justify-between px-3 py-2 text-sm text-text hover:bg-surface-muted"
        >
          <span className="font-semibold">Race History</span>
          <span className="text-text-muted text-xs">{showRace ? "▼" : "▶"}</span>
        </button>
        {showRace && (
          <div className="px-3 pb-3">
            <RaceHistoryTable characterId={characterId} />
          </div>
        )}
      </div>
    </div>
  );
}