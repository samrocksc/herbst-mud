/* eslint-disable functional/immutable-data, functional/no-loop-statements */
import { createFileRoute, Link, Outlet, useLocation } from "@tanstack/react-router";
import { useState, useMemo } from "react";
import { useQuery } from "@tanstack/react-query";
import { apiGet, API_BASE, buildWorldParams } from "../../utils/apiFetch";
import { useWorldStore } from "../../contexts/WorldStoreContext";
import { PageHeader } from "../../components/PageHeader";
import { DataTable, type Column } from "../../components/DataTable";
import { Button } from "../../components/Button";
import { PageContainer } from "../../components/PageContainer";

export const Route = createFileRoute("/_auth/npcs")({
  component: NPCTemplatesIndex,
});

// ─── Types ──────────────────────────────────────────────────────────────────

type NPCTemplate = Readonly<{
  id: string
  slug: string
  name: string
  level: number
  xp_value: number
  respawn_rooms: string[]
  respawn_cooldown: number
  race_id: number
}>

// ─── Component ──────────────────────────────────────────────────────────────

export function NPCTemplatesIndex() {
  const { currentWorld } = useWorldStore();
  const [searchQuery, setSearchQuery] = useState("");

  const params = buildWorldParams(currentWorld);
  const qs = params ? `?${params}` : "";

  // ── Query ────────────────────────────────────────────────────────────────

  const templatesQuery = useQuery({
    queryKey: ["npc-templates", currentWorld],
    queryFn: () => apiGet<NPCTemplate[]>(`${API_BASE}/api/npc-templates${qs}`),
  });

  // Fetch instances once to count by template
  const instancesQuery = useQuery({
    queryKey: ["npc-instances-count", currentWorld],
    queryFn: () => apiGet<Array<{ npc_template_id: string }>>(`${API_BASE}/api/npc-instances${qs}`),
  });

  const instanceCounts = useMemo(() => {
    const counts: Record<string, number> = {};
    for (const inst of instancesQuery.data ?? []) {
      const tid = inst.npc_template_id;
      if (tid) counts[tid] = (counts[tid] ?? 0) + 1;
    }
    return counts;
  }, [instancesQuery.data]);

  // ── Search filter ────────────────────────────────────────────────────────

  const filteredTemplates = (templatesQuery.data ?? []).filter((template) =>
    template.name.toLowerCase().includes(searchQuery.toLowerCase()),
  );

  // ── Render ────────────────────────────────────────────────────────────────

  const location = useLocation();
  const isList = location.pathname === "/npcs";

  const columns = useMemo<Column<NPCTemplate>[]>(
    () => [
      { header: "Slug", accessor: "slug" },
      {
        header: "Name",
        accessor: "name",
        render: (_: unknown, row: NPCTemplate) => (
          <Link
            to="/npcs/$npcId"
            params={{ npcId: row.id }}
            className="no-underline text-primary hover:underline font-bold"
          >
            {row.name}
          </Link>
        ),
      },
      { header: "Level", accessor: "level", align: "center" },
      { header: "XP Value", accessor: "xp_value", align: "right" },
      { header: "Respawn Cooldown", accessor: "respawn_cooldown", align: "center" },
      {
        header: "Instances",
        accessor: "instances",
        align: "center",
        render: (_: unknown, row: NPCTemplate) => (
          <span className="badge badge-neutral">
            {instanceCounts[row.id] ?? 0}
          </span>
        ),
      },
    ],
    [instanceCounts],
  );

  if (!isList) {
    return <Outlet />;
  }

  return (
    <PageContainer>
      <PageHeader
        title="NPC Templates"
        showBack
        backTo="/dashboard"
        actions={
          <Link to="/npcs/new">
            <Button variant="primary" size="sm">
              + Add Template
            </Button>
          </Link>
        }
      />

      {/* Search bar */}
      <div className="mb-4">
        <input
          type="text"
          placeholder="Search templates by name..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="w-full max-w-sm p-2 bg-surface border border-border rounded text-text text-sm"
        />
      </div>

      {/* Loading */}
      {templatesQuery.isLoading && (
        <div className="p-8 text-text-muted text-center text-xs">Loading NPC templates...</div>
      )}

      {/* Error */}
      {templatesQuery.isError && (
        <div className="p-4 bg-danger/10 border border-danger rounded text-danger text-xs">
          Failed to load NPC templates: {templatesQuery.error?.message ?? "Unknown error"}
        </div>
      )}

      {/* Data table */}
      {templatesQuery.isSuccess && (
        <DataTable<NPCTemplate>
          columns={columns}
          data={filteredTemplates}
          getKey={(row) => row.id}
          emptyMessage="No NPC templates found."
          variant="dark"
        />
      )}
    </PageContainer>
  );
}
