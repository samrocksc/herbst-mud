/* eslint-disable functional/immutable-data, functional/no-loop-statements */
import { createFileRoute, Link, Outlet, useLocation } from "@tanstack/react-router";
import { useState, useMemo } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiGet, apiPost } from "../../utils/apiFetch";
import { useWorldStore } from "../../contexts/WorldStoreContext";
import { PageHeader } from "../../components/PageHeader";
import { DataTable, type Column } from "../../components/DataTable";
import { Modal } from "../../components/Modal";
import { Button } from "../../components/Button";

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
}>

type NPCTemplateForm = Readonly<{
  name: string
  description: string
  race: string
  level: number
  xp_value: number
  respawn_cooldown: number
  respawn_rooms: string
}>

const API = `${window.location.origin}`;

// ─── Empty form ─────────────────────────────────────────────────────────────

const EMPTY_FORM: NPCTemplateForm = {
  name: "",
  description: "",
  race: "",
  level: 1,
  xp_value: 0,
  respawn_cooldown: 60,
  respawn_rooms: "",
};

// ─── Component ──────────────────────────────────────────────────────────────

function NPCTemplatesIndex() {
  const { currentWorld } = useWorldStore();
  const queryClient = useQueryClient();
  const [searchQuery, setSearchQuery] = useState("");
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [form, setForm] = useState<NPCTemplateForm>(EMPTY_FORM);

  // Build query string with world_id
  const params = useMemo(() => {
    const p = new URLSearchParams();
    if (currentWorld) p.append("world_id", currentWorld);
    return p.toString();
  }, [currentWorld]);
  const qs = params ? `?${params}` : "";

  // ── Query ────────────────────────────────────────────────────────────────

  const templatesQuery = useQuery({
    queryKey: ["npc-templates", currentWorld],
    queryFn: () => apiGet<NPCTemplate[]>(`${API}/api/npc-templates${qs}`),
  });

  // Fetch instances once to count by template
  const instancesQuery = useQuery({
    queryKey: ["npc-instances-count", currentWorld],
    queryFn: () => apiGet<Array<{ npc_template_id: string }>>(`${API}/api/npc-instances${qs}`),
  });

  const instanceCounts = useMemo(() => {
    const counts: Record<string, number> = {};
    for (const inst of instancesQuery.data ?? []) {
      const tid = inst.npc_template_id;
      if (tid) counts[tid] = (counts[tid] ?? 0) + 1;
    }
    return counts;
  }, [instancesQuery.data]);

  // ── Create mutation ──────────────────────────────────────────────────────

  const createMutation = useMutation({
    mutationFn: (input: NPCTemplateForm) => {
      const rooms = input.respawn_rooms
        .split(",")
        .map((s) => s.trim())
        .filter((s) => s !== "");
      return apiPost<NPCTemplate>(`${API}/api/npc-templates${qs}`, {
        name: input.name,
        description: input.description,
        race: input.race,
        level: input.level,
        xp_value: input.xp_value,
        respawn_cooldown: input.respawn_cooldown,
        respawn_rooms: rooms,
        skills: {},
        trades_with: [],
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["npc-templates"] });
      setShowCreateModal(false);
      setForm(EMPTY_FORM);
    },
  });

  // ── Search filter ────────────────────────────────────────────────────────

  const filteredTemplates = (templatesQuery.data ?? []).filter((template) =>
    template.name.toLowerCase().includes(searchQuery.toLowerCase()),
  );

  // ── Handlers ─────────────────────────────────────────────────────────────

  const handleCreate = () => {
    if (!form.name.trim()) return;
    createMutation.mutate(form);
  };

  const handleModalClose = () => {
    setShowCreateModal(false);
    setForm(EMPTY_FORM);
  };

  // ── Render ────────────────────────────────────────────────────────────────

  const location = useLocation();
  const isList = location.pathname === "/npcs";

  const columns = useMemo<Column<NPCTemplate>[]>(
    () => [
      { header: "ID", accessor: "id" },
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
    <div className="p-6 max-w-[1200px] mx-auto">
      <PageHeader
        title="NPC Templates"
        showBack
        backTo="/dashboard"
        actions={
          <Button variant="primary" size="sm" onClick={() => setShowCreateModal(true)}>
            + Add Template
          </Button>
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

      {/* Create modal */}
      <Modal isOpen={showCreateModal} onClose={handleModalClose} title="Add NPC Template">
        <div className="flex flex-col gap-4">
          {/* Name */}
          <div>
            <label className="text-text-muted text-xs block mb-1">Name *</label>
            <input
              type="text"
              value={form.name}
              onChange={(e) => setForm({ ...form, name: e.target.value })}
              placeholder="Display name"
              className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
            />
          </div>

          {/* Description */}
          <div>
            <label className="text-text-muted text-xs block mb-1">Description</label>
            <textarea
              value={form.description}
              rows={3}
              onChange={(e) => setForm({ ...form, description: e.target.value })}
              placeholder="Flavor text..."
              className="w-full p-2 bg-surface border border-border rounded text-text text-sm resize-y"
            />
          </div>

          {/* Race */}
          <div>
            <label className="text-text-muted text-xs block mb-1">Race</label>
            <input
              type="text"
              value={form.race}
              onChange={(e) => setForm({ ...form, race: e.target.value })}
              placeholder="e.g. goblin"
              className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
            />
          </div>

          {/* Level & XP Value */}
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="text-text-muted text-xs block mb-1">Level *</label>
              <input
                type="number"
                value={form.level}
                onChange={(e) => setForm({ ...form, level: parseInt(e.target.value) || 1 })}
                min={1}
                className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
              />
            </div>
            <div>
              <label className="text-text-muted text-xs block mb-1">XP Value *</label>
              <input
                type="number"
                value={form.xp_value}
                onChange={(e) => setForm({ ...form, xp_value: parseInt(e.target.value) || 0 })}
                min={0}
                className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
              />
            </div>
          </div>

          {/* Respawn Cooldown & Respawn Rooms */}
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="text-text-muted text-xs block mb-1">
                Respawn Cooldown <span className="text-text-muted">(seconds)</span>
              </label>
              <input
                type="number"
                value={form.respawn_cooldown}
                onChange={(e) => setForm({ ...form, respawn_cooldown: parseInt(e.target.value) || 0 })}
                min={0}
                className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
              />
            </div>
            <div>
              <label className="text-text-muted text-xs block mb-1">
                Respawn Rooms <span className="text-text-muted">(comma-separated)</span>
              </label>
              <input
                type="text"
                value={form.respawn_rooms}
                onChange={(e) => setForm({ ...form, respawn_rooms: e.target.value })}
                placeholder="1, 2, 3"
                className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
              />
            </div>
          </div>

          {/* Error display */}
          {createMutation.isError && (
            <div className="p-2 bg-danger/10 border border-danger rounded text-danger text-xs">
              Failed to create template: {createMutation.error?.message ?? "Unknown error"}
            </div>
          )}

          {/* Actions */}
          <div className="flex gap-2">
            <Button
              variant="primary"
              size="md"
              fullWidth
              onClick={handleCreate}
              disabled={!form.name.trim() || createMutation.isPending}
            >
              {createMutation.isPending ? "Creating..." : "Create Template"}
            </Button>
            <Button variant="secondary" size="md" fullWidth onClick={handleModalClose}>
              Cancel
            </Button>
          </div>
        </div>
      </Modal>
    </div>
  );
}
