/* eslint-disable functional/immutable-data */
import { createFileRoute } from "@tanstack/react-router";
import { useState, useEffect, useRef, useCallback } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { useEventLogs, type EventLog } from "../../hooks/useEventLogs";
import { PageHeader } from "../../components/PageHeader";
import { DataTable, type Column } from "../../components/DataTable";
import { Button } from "../../components/Button";
import { PageContainer } from "../../components/PageContainer";
import { FilterBar } from "../../components/FilterBar";

export const Route = createFileRoute("/_auth/event-logs")({
  component: EventLogsPage,
});

const EVENT_TYPES = [
  "xp.gained",
  "level.up",
  "skill.leveled_up",
  "reclass",
  "rerace",
  "npc.defeated",
  "achievement.completed",
] as const;

const TYPE_BADGE: Record<string, string> = {
  "xp.gained": "bg-emerald-900/40 text-emerald-300 border-emerald-600/30",
  "level.up": "bg-violet-900/40 text-violet-300 border-violet-600/30",
  "skill.leveled_up": "bg-sky-900/40 text-sky-300 border-sky-600/30",
  "reclass": "bg-amber-900/40 text-amber-300 border-amber-600/30",
  "rerace": "bg-amber-900/40 text-amber-300 border-amber-600/30",
  "npc.defeated": "bg-red-900/40 text-red-300 border-red-600/30",
  "achievement.completed": "bg-yellow-900/40 text-yellow-300 border-yellow-600/30",
};

function EventLogsPage() {
  const [filterType, setFilterType] = useState<string>("");
  const [characterIdInput, setCharacterIdInput] = useState<string>("");
  const [characterIdFilter, setCharacterIdFilter] = useState<number | undefined>(undefined);
  const [limit, setLimit] = useState(100);
  const [autoRefresh, setAutoRefresh] = useState(false);

  const qc = useQueryClient();

  const filters = {
    event_type: filterType || undefined,
    character_id: characterIdFilter,
    limit,
  };

  const { data, isLoading, error } = useEventLogs(filters);

  // Auto-refresh: invalidate the query every 5 seconds when enabled
  useEffect(() => {
    if (!autoRefresh) return;
    const id = setInterval(() => {
      qc.invalidateQueries({ queryKey: ["event-logs"] });
    }, 5000);
    return () => clearInterval(id);
  }, [autoRefresh, qc]);

  const handleApplyCharacterId = useCallback(() => {
    const parsed = parseInt(characterIdInput, 10);
    setCharacterIdFilter(isNaN(parsed) ? undefined : parsed);
  }, [characterIdInput]);

  const handleClearFilters = () => {
    setFilterType("");
    setCharacterIdInput("");
    setCharacterIdFilter(undefined);
  };

  const handleLoadMore = () => {
    setLimit((prev) => prev + 100);
  };

  const logs = data ?? [];

  const columns: Column<EventLog>[] = [
    {
      header: "Timestamp",
      accessor: "created_at",
      render: (val: unknown) => (
        <span className="text-xs text-text-muted whitespace-nowrap">
          {formatTimestamp(String(val ?? ""))}
        </span>
      ),
      className: "w-40",
    },
    {
      header: "Event Type",
      accessor: "event_type",
      render: (val: unknown) => {
        const typeStr = String(val ?? "");
        const badgeClass = TYPE_BADGE[typeStr] ?? "bg-slate-700/50 text-slate-300 border-slate-600/30";
        return (
          <span className={`text-xs px-2 py-0.5 rounded border font-medium ${badgeClass}`}>
            {typeStr}
          </span>
        );
      },
    },
    {
      header: "Character",
      accessor: "character_id",
      render: (val: unknown) => {
        if (val == null) return <span className="text-text-muted">—</span>;
        return <span className="text-text-muted text-sm">#{String(val)}</span>;
      },
    },
    {
      header: "Summary",
      accessor: "payload",
      render: (val: unknown) => {
        const payload = val as Record<string, unknown> | null;
        if (!payload || typeof payload !== "object") {
          return <span className="text-text-muted text-xs">—</span>;
        }
        return <span className="text-text-muted text-xs">{summarizePayload(payload)}</span>;
      },
    },
  ];

  if (isLoading) return <div className="loading p-8 text-center text-text-muted">Loading event logs...</div>;
  if (error) return <div className="error p-8 text-center text-danger">Failed to load event logs: {error.message}</div>;

  return (
    <PageContainer>
      <PageHeader
        title="Event Logs"
        backTo="/dashboard"
        actions={
          <div className="flex items-center gap-2">
            <span className="text-xs text-text-muted">{logs.length} events</span>
            <Button
              variant={autoRefresh ? "primary" : "secondary"}
              size="sm"
              onClick={() => setAutoRefresh((v) => !v)}
            >
              {autoRefresh ? "● Auto-Refresh (5s)" : "○ Auto-Refresh"}
            </Button>
          </div>
        }
      />

      <p className="text-sm text-muted mb-4">
        View game events such as XP gains, level-ups, class/race changes, NPC defeats, and achievement completions.
        Enable auto-refresh to poll for new events every 5 seconds.
      </p>

      <FilterBar
        showClear={!!filterType || characterIdFilter !== undefined}
        onClear={handleClearFilters}
      >
        <div className="flex flex-col gap-1">
          <label className="text-xs text-text-muted">Event Type:</label>
          <select
            value={filterType}
            onChange={(e) => setFilterType(e.target.value)}
            className="px-3 py-2 bg-surface border border-border rounded text-sm text-text focus:outline-none focus:border-primary"
          >
            <option value="">All Events</option>
            {EVENT_TYPES.map((t) => (
              <option key={t} value={t}>{t}</option>
            ))}
          </select>
        </div>

        <div className="flex flex-col gap-1">
          <label className="text-xs text-text-muted">Character ID:</label>
          <div className="flex gap-2">
            <input
              type="text"
              value={characterIdInput}
              onChange={(e) => setCharacterIdInput(e.target.value)}
              onKeyDown={(e) => { if (e.key === "Enter") handleApplyCharacterId(); }}
              placeholder="e.g. 42"
              className="px-3 py-2 bg-surface border border-border rounded text-sm text-text focus:outline-none focus:border-primary w-28"
            />
            <Button variant="secondary" size="sm" onClick={handleApplyCharacterId}>
              Filter
            </Button>
          </div>
        </div>
      </FilterBar>

      <DataTable
        columns={columns}
        data={logs}
        getKey={(row) => row.id}
        emptyMessage={
          filterType || characterIdFilter !== undefined
            ? "No events match your current filters."
            : "No event logs found."
        }
      />

      {logs.length >= limit && (
        <div className="text-center mt-4">
          <Button variant="secondary" onClick={handleLoadMore}>
            Load More (+100)
          </Button>
        </div>
      )}
    </PageContainer>
  );
}

/**
 * Extract a human-readable summary from an event payload.
 * Different event types store different keys.
 */
function summarizePayload(payload: Record<string, unknown>): string {
  const parts: string[] = [];

  const amount = payload["amount"];
  const xpAmount = payload["xp_amount"];
  const skillName = payload["skill_name"];
  const skillLevel = payload["skill_level"];
  const oldLevel = payload["old_level"];
  const newLevel = payload["new_level"];
  const oldClass = payload["old_class"];
  const newClass = payload["new_class"];
  const oldRace = payload["old_race"];
  const newRace = payload["new_race"];
  const npcName = payload["npc_name"];
  const npcId = payload["npc_id"];
  const achievementName = payload["achievement_name"];
  const reason = payload["reason"];

  if (amount != null) parts.push(`amount: ${amount}`);
  if (xpAmount != null) parts.push(`xp: ${xpAmount}`);
  if (skillName != null) parts.push(`skill: ${skillName}`);
  if (skillLevel != null) parts.push(`level: ${skillLevel}`);
  if (oldLevel != null && newLevel != null) parts.push(`${oldLevel} → ${newLevel}`);
  if (oldClass != null && newClass != null) parts.push(`${oldClass} → ${newClass}`);
  if (oldRace != null && newRace != null) parts.push(`${oldRace} → ${newRace}`);
  if (npcName != null) parts.push(`npc: ${npcName}`);
  if (npcId != null && npcName == null) parts.push(`npc_id: ${npcId}`);
  if (achievementName != null) parts.push(`achievement: ${achievementName}`);
  if (reason != null) parts.push(`reason: ${reason}`);

  if (parts.length === 0) {
    const entries = Object.entries(payload).slice(0, 3);
    return entries.map(([k, v]) => `${k}: ${typeof v === "object" ? JSON.stringify(v) : String(v)}`).join(", ");
  }

  return parts.join(" · ");
}

function formatTimestamp(t: string): string {
  if (!t) return "—";
  try {
    const d = new Date(t);
    const now = new Date();
    const diffMs = now.getTime() - d.getTime();
    const diffMin = Math.floor(diffMs / 60000);
    if (diffMin < 1) return "just now";
    if (diffMin < 60) return `${diffMin}m ago`;
    const isToday = d.toDateString() === now.toDateString();
    if (isToday) return d.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });
    const isThisYear = d.getFullYear() === now.getFullYear();
    if (isThisYear) return d.toLocaleDateString([], { month: "short", day: "numeric" }) + " " + d.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });
    return d.toLocaleDateString([], { year: "numeric", month: "short", day: "numeric" }) + " " + d.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });
  } catch {
    return t;
  }
}