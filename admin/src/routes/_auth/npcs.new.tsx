import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { apiPost, API_BASE } from "../../utils/apiFetch";
import { useWorldStore } from "../../contexts/WorldStoreContext";
import { useRaces } from "../../hooks/useRaces";
import { PageHeader } from "../../components/PageHeader";
import { Button } from "../../components/Button";
import { PageContainer } from "../../components/PageContainer";
import { SearchableSelect } from "../../components/SearchableSelect";

export const Route = createFileRoute("/_auth/npcs/new")({
  component: CreateNPCPage,
});

type NPCForm = {
  name: string
  description: string
  race_id: number
  disposition: string
  level: number
  xp_value: number
  respawn_cooldown: number
  respawn_rooms: string
  greeting: string
  skills: string
  trades_with: string
  roam_pattern: string
  roam_interval_seconds: number
  roam_pause_min_seconds: number
  roam_pause_max_seconds: number
  roam_zone_ids: string
  notify_on_enter: boolean
};

const EMPTY_FORM: NPCForm = {
  name: "",
  description: "",
  race_id: 0,
  disposition: "neutral",
  level: 1,
  xp_value: 0,
  respawn_cooldown: 60,
  respawn_rooms: "",
  greeting: "",
  skills: "",
  trades_with: "",
  roam_pattern: "static",
  roam_interval_seconds: 60,
  roam_pause_min_seconds: 15,
  roam_pause_max_seconds: 120,
  roam_zone_ids: "",
  notify_on_enter: true,
};

function CreateNPCPage() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { currentWorld } = useWorldStore();
  const [form, setForm] = useState<NPCForm>(EMPTY_FORM);

  const qs = currentWorld ? `?world_id=${currentWorld}` : "";

  const { data: races } = useRaces();

  const createMutation = useMutation({
    mutationFn: (input: NPCForm) => {
      const rooms = input.respawn_rooms
        .split(",")
        .map((s) => s.trim())
        .filter((s) => s !== "");
      const roamZoneIds = input.roam_zone_ids
        .split(",")
        .map((s) => s.trim())
        .filter((s) => s !== "");
      const skills: Record<string, number> = {}
      for (const line of input.skills.split("\n")) {
        const trimmed = line.trim();
        if (!trimmed) continue;
        const colonIdx = trimmed.lastIndexOf(":");
        if (colonIdx > 0) {
          skills[trimmed.slice(0, colonIdx)] = parseInt(trimmed.slice(colonIdx + 1)) || 0;
        } else {
          skills[trimmed] = 0;
        }
      }
      return apiPost(`${API_BASE}/api/npc-templates${qs}`, {
        name: input.name,
        description: input.description,
        race_id: input.race_id,
        disposition: input.disposition,
        level: input.level,
        xp_value: input.xp_value,
        respawn_cooldown: input.respawn_cooldown,
        respawn_rooms: rooms,
        greeting: input.greeting,
        skills,
        trades_with: input.trades_with.split("\n").map((s) => s.trim()).filter(Boolean),
        roam_pattern: input.roam_pattern,
        roam_interval_seconds: input.roam_interval_seconds,
        roam_pause_min_seconds: input.roam_pause_min_seconds,
        roam_pause_max_seconds: input.roam_pause_max_seconds,
        roam_zone_ids: roamZoneIds,
        notify_on_enter: input.notify_on_enter,
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["npc-templates"] });
      navigate({ to: "/npcs" });
    },
    onError: (error: unknown) => {
      console.error("Failed to create NPC template:", error);
      const message = error instanceof Error ? error.message : "Failed to create NPC template";
      alert(message); // Show visible error
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!form.name.trim()) return;
    createMutation.mutate(form, {
      onError: (error: unknown) => {
        console.error("NPC creation error:", error);
        const message = error instanceof Error ? error.message : "Failed to create NPC template";
        alert(message);
      },
    });
  };

  const set = (patch: Partial<NPCForm>) => setForm((prev) => ({ ...prev, ...patch }));

  return (
    <PageContainer>
      <PageHeader title="Create NPC Template" showBack backTo="/npcs" />

      <div className="card bg-surface p-6 border border-border rounded">
        <form onSubmit={handleSubmit} className="space-y-4">
          {/* Name */}
          <div>
            <label className="text-text-muted text-xs block mb-1">Name *</label>
            <input
              type="text"
              value={form.name}
              onChange={(e) => set({ name: e.target.value })}
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
              onChange={(e) => set({ description: e.target.value })}
              placeholder="Flavor text..."
              className="w-full p-2 bg-surface border border-border rounded text-text text-sm resize-y"
            />
          </div>

          {/* Race */}
          <div>
            <label className="text-text-muted text-xs block mb-1">Race</label>
            <SearchableSelect
              options={(races ?? []).map((r) => ({ id: String(r.id), name: r.display_name || r.name }))}
              value={String(form.race_id || "")}
              onChange={(v) => set({ race_id: Number(v) || 0 })}
              placeholder="Select race..."
            />
          </div>

          {/* Level & XP Value */}
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="text-text-muted text-xs block mb-1">Level *</label>
              <input
                type="number"
                value={form.level}
                onChange={(e) => set({ level: parseInt(e.target.value) || 1 })}
                min={1}
                className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
              />
            </div>
            <div>
              <label className="text-text-muted text-xs block mb-1">XP Value *</label>
              <input
                type="number"
                value={form.xp_value}
                onChange={(e) => set({ xp_value: parseInt(e.target.value) || 0 })}
                min={0}
                className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
              />
            </div>
          </div>

          {/* Respawn Cooldown & Rooms */}
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="text-text-muted text-xs block mb-1">
                Respawn Cooldown <span className="text-text-muted">(seconds)</span>
              </label>
              <input
                type="number"
                value={form.respawn_cooldown}
                onChange={(e) => set({ respawn_cooldown: parseInt(e.target.value) || 0 })}
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
                onChange={(e) => set({ respawn_rooms: e.target.value })}
                placeholder="1, 2, 3"
                className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
              />
            </div>
          </div>

          {/* Disposition & Greeting */}
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="text-text-muted text-xs block mb-1">Disposition</label>
              <select
                value={form.disposition}
                onChange={(e) => set({ disposition: e.target.value })}
                className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
              >
                <option value="neutral">Neutral</option>
                <option value="friendly">Friendly</option>
                <option value="hostile">Hostile</option>
                <option value="shopkeeper">Shopkeeper</option>
              </select>
            </div>
            <div>
              <label className="text-text-muted text-xs block mb-1">Greeting</label>
              <input
                type="text"
                value={form.greeting}
                onChange={(e) => set({ greeting: e.target.value })}
                placeholder="Hello, traveler..."
                className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
              />
            </div>
          </div>

          {/* Skills & Trades With */}
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="text-text-muted text-xs block mb-1">
                Skills <span className="text-text-muted">(skill:value per line)</span>
              </label>
              <textarea
                value={form.skills}
                onChange={(e) => set({ skills: e.target.value })}
                rows={3}
                placeholder={"blades:30\nstaves:20"}
                className="w-full p-2 bg-surface border border-border rounded text-text text-sm resize-y font-mono"
              />
            </div>
            <div>
              <label className="text-text-muted text-xs block mb-1">
                Trades With <span className="text-text-muted">(one per line)</span>
              </label>
              <textarea
                value={form.trades_with}
                onChange={(e) => set({ trades_with: e.target.value })}
                rows={3}
                placeholder={"tag:merchant\ntag:blacksmith"}
                className="w-full p-2 bg-surface border border-border rounded text-text text-sm resize-y font-mono"
              />
            </div>
          </div>

          {/* Behavior */}
          <h3 className="mt-6 mb-2 text-text text-sm font-semibold border-b border-border pb-1">Behavior</h3>
          <div>
            <label className="text-text-muted text-xs block mb-1">Roam Pattern</label>
            <select
              value={form.roam_pattern}
              onChange={(e) => set({ roam_pattern: e.target.value })}
              className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
            >
              <option value="static">Static (never moves)</option>
              <option value="wander">Wander (random within zones)</option>
              <option value="patrol">Patrol (sequential through exits)</option>
              <option value="return_home">Return Home (back to home room)</option>
            </select>
          </div>
          <div className="grid grid-cols-3 gap-4">
            <div>
              <label className="text-text-muted text-xs block mb-1">
                Roam Interval <span className="text-text-muted">(seconds)</span>
              </label>
              <input
                type="number"
                value={form.roam_interval_seconds}
                onChange={(e) => set({ roam_interval_seconds: parseInt(e.target.value) || 0 })}
                min={0}
                className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
              />
            </div>
            <div>
              <label className="text-text-muted text-xs block mb-1">
                Roam Pause Min <span className="text-text-muted">(seconds)</span>
              </label>
              <input
                type="number"
                value={form.roam_pause_min_seconds}
                onChange={(e) => set({ roam_pause_min_seconds: parseInt(e.target.value) || 0 })}
                min={0}
                className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
              />
            </div>
            <div>
              <label className="text-text-muted text-xs block mb-1">
                Roam Pause Max <span className="text-text-muted">(seconds)</span>
              </label>
              <input
                type="number"
                value={form.roam_pause_max_seconds}
                onChange={(e) => set({ roam_pause_max_seconds: parseInt(e.target.value) || 0 })}
                min={0}
                className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
              />
            </div>
          </div>
          <div>
            <label className="text-text-muted text-xs block mb-1">
              Roam Zone IDs <span className="text-text-muted">(comma-separated)</span>
            </label>
            <input
              type="text"
              value={form.roam_zone_ids}
              onChange={(e) => set({ roam_zone_ids: e.target.value })}
              placeholder="zone-1, zone-2"
              className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
            />
          </div>
          <div>
            <label className="flex items-center gap-2 text-text-muted text-xs cursor-pointer">
              <input
                type="checkbox"
                checked={form.notify_on_enter}
                onChange={(e) => set({ notify_on_enter: e.target.checked })}
              />
              Notify on Enter
            </label>
          </div>

          {/* Error display */}
          {createMutation.isError && (
            <div className="p-2 bg-danger/10 border border-danger rounded text-danger text-xs">
              Failed to create template: {createMutation.error?.message ?? "Unknown error"}
            </div>
          )}

          {/* Actions */}
          <div className="flex gap-2 justify-end mt-6">
            <Button variant="secondary" onClick={() => navigate({ to: "/npcs" })}>
              Cancel
            </Button>
            <Button
              variant="primary"
              type="submit"
              disabled={!form.name.trim() || createMutation.isPending}
            >
              {createMutation.isPending ? "Creating..." : "Create Template"}
            </Button>
          </div>
        </form>
      </div>
    </PageContainer>
  );
}
