import { createFileRoute } from "@tanstack/react-router";
import { useEffect, useState, useCallback } from "react";
import { apiGet, apiPost, apiPut } from "../utils/apiFetch";
import { Button } from "../components/Button";
import { Modal } from "../components/Modal";
import { PageHeader } from "../components/PageHeader";
import { FormField, NumberField, TextareaField, SelectField } from "../components/fields";
import { FormError } from "../components/fields/FormError";
import { showToast } from "../components/Toast";

type NPCTemplate = Readonly<{
  id: string
  name: string
  level: number
  xp_value: number
  respawn_rooms: string[]
  respawn_cooldown: number
}>

type Race = Readonly<{
  id: number
  name: string
  display_name: string
}>

export const Route = createFileRoute("/npc-templates")({
  component: NPCTemplatePage,
});

const DISPOSITION_OPTS = [
  { value: "neutral", label: "Neutral" },
  { value: "friendly", label: "Friendly" },
  { value: "hostile", label: "Hostile" },
];

function NPCTemplatePage() {
  const [templates, setTemplates] = useState<NPCTemplate[]>([]);
  const [races, setRaces] = useState<Race[]>([]);
  const [loading, setLoading] = useState(true);
  const [editingID, setEditingID] = useState<string | null>(null);
  const [saving, setSaving] = useState(false);
  const [roomsInput, setRoomsInput] = useState("");
  const [cooldownInput, setCooldownInput] = useState("60");
  const [showCreate, setShowCreate] = useState(false);
  const [createForm, setCreateForm] = useState({
    id: "", name: "", description: "", race_id: "", disposition: "neutral",
    level: 1, xp_value: 0, greeting: "", respawn_cooldown: 60, respawn_rooms: "",
  });
  const [createError, setCreateError] = useState("");

  const load = useCallback(async () => {
    try {
      const [data, raceData] = await Promise.all([
        apiGet<NPCTemplate[]>(`${window.location.origin}/api/npc-templates`),
        apiGet<{ races: Race[] }>(`${window.location.origin}/api/races`).then(r => r.races ?? []).catch(() => [] as Race[]),
      ]);
      setTemplates(data);
      setRaces(raceData);
    } catch { showToast("Failed to load NPC templates", "error"); }
    finally { setLoading(false); }
  }, []);

  useEffect(() => { load(); }, [load]);

  async function handleSave(id: string) {
    setSaving(true);
    try {
      const rooms = roomsInput.split(",").map((s) => s.trim().replace(/^r/i, "")).filter(Boolean);
      const cooldown = parseInt(cooldownInput, 10) || 0;
      const current = templates.find((t) => t.id === id);
      await apiPut(`${window.location.origin}/api/npc-templates/${id}`, {
        xp_value: current?.xp_value ?? 0, respawn_rooms: rooms, respawn_cooldown: cooldown,
      });
      await load(); setEditingID(null);
      showToast("NPC template saved", "success");
    } catch { showToast("Failed to save NPC template", "error"); }
    finally { setSaving(false); }
  }

  async function handleCreate() {
    setCreateError("");
    try {
      const rooms = createForm.respawn_rooms.split(",").map((s) => s.trim().replace(/^r/i, "")).filter(Boolean);
      const payload: Record<string, unknown> = {
        id: createForm.id, name: createForm.name, description: createForm.description,
        disposition: createForm.disposition, level: createForm.level,
        xp_value: createForm.xp_value, greeting: createForm.greeting,
        respawn_cooldown: createForm.respawn_cooldown, respawn_rooms: rooms, skills: {}, trades_with: [],
      };
      if (createForm.race_id !== "") payload.race_id = Number(createForm.race_id);
      await apiPost(`${window.location.origin}/api/npc-templates`, payload);
      await load(); setShowCreate(false);
      setCreateForm({ id: "", name: "", description: "", race_id: "", disposition: "neutral", level: 1, xp_value: 0, greeting: "", respawn_cooldown: 60, respawn_rooms: "" });
    } catch (e: unknown) { setCreateError(e instanceof Error ? e.message : String(e)); }
  }

  function startEdit(t: NPCTemplate) {
    setEditingID(t.id); setRoomsInput(t.respawn_rooms?.join(", ") ?? ""); setCooldownInput(String(t.respawn_cooldown ?? 60));
  }

  const raceOpts = races.map(r => ({ value: String(r.id), label: r.display_name || r.name }));

  if (loading) return <div className="min-h-screen bg-surface p-6"><PageHeader title="NPC Templates" backTo="/dashboard" /><p className="text-text-muted">Loading templates...</p></div>;

  return (
    <div className="min-h-screen bg-surface p-6">
      <PageHeader title="NPC Templates" backTo="/dashboard"
        actions={<Button variant="primary" onClick={() => setShowCreate(true)}>+ New NPC</Button>} />
      <div className="space-y-4 max-w-[900px]">
        {templates.map((t) => (
          <div key={t.id} className="bg-surface-muted border border-border rounded p-4 space-y-3">
            <div className="flex justify-between items-start">
              <div><div className="text-text font-medium">{t.name}</div><div className="text-text-muted text-xs">Level {t.level} · XP {t.xp_value} · ID: {t.id}</div></div>
              {editingID !== t.id ? (
                <Button variant="accent" size="sm" onClick={() => startEdit(t)}>Edit</Button>
              ) : (
                <div className="flex gap-2">
                  <Button variant="primary" size="sm" onClick={() => handleSave(t.id)} disabled={saving}>{saving ? "Saving..." : "Save"}</Button>
                  <Button variant="ghost" size="sm" onClick={() => setEditingID(null)}>Cancel</Button>
                </div>
              )}
            </div>
            {editingID === t.id && (
              <div className="space-y-3">
                <FormField label="Respawn Rooms" value={roomsInput} onChange={setRoomsInput} placeholder="1, 2, 3" />
                <NumberField label="Respawn Cooldown (seconds)" value={parseInt(cooldownInput) || 0} onChange={(v) => setCooldownInput(String(v))} />
              </div>
            )}
          </div>
        ))}
        {templates.length === 0 && <p className="text-text-muted">No NPC templates found.</p>}
      </div>
      <Modal isOpen={showCreate} onClose={() => setShowCreate(false)} title="Create NPC Instance">
        <div className="space-y-3">
          {createError && <FormError message={createError} />}
          <FormField label="ID" value={createForm.id} onChange={(v) => setCreateForm({...createForm, id: v})} placeholder="e.g. goblin_guard_01" />
          <FormField label="Name" value={createForm.name} onChange={(v) => setCreateForm({...createForm, name: v})} placeholder="Display name" />
          <TextareaField label="Description" value={createForm.description} onChange={(v) => setCreateForm({...createForm, description: v})} rows={2} placeholder="Flavor text..." />
          <div className="flex gap-3">
            <div className="flex-1"><SelectField label="Race" value={createForm.race_id} onChange={(v) => setCreateForm({...createForm, race_id: v})} options={raceOpts} /></div>
            <div className="flex-1"><SelectField label="Disposition" value={createForm.disposition} onChange={(v) => setCreateForm({...createForm, disposition: v})} options={DISPOSITION_OPTS} /></div>
          </div>
          <div className="flex gap-3">
            <div className="flex-1"><NumberField label="Level" value={createForm.level} onChange={(v) => setCreateForm({...createForm, level: v})} /></div>
            <div className="flex-1"><NumberField label="XP Value" value={createForm.xp_value} onChange={(v) => setCreateForm({...createForm, xp_value: v})} /></div>
          </div>
          <TextareaField label="Greeting" value={createForm.greeting} onChange={(v) => setCreateForm({...createForm, greeting: v})} rows={2} placeholder="NPC greeting message..." />
          <div className="flex gap-3">
            <div className="flex-1"><NumberField label="Respawn Cooldown (s)" value={createForm.respawn_cooldown} onChange={(v) => setCreateForm({...createForm, respawn_cooldown: v})} /></div>
            <div className="flex-1"><FormField label="Respawn Rooms" value={createForm.respawn_rooms} onChange={(v) => setCreateForm({...createForm, respawn_rooms: v})} placeholder="1, 2, 3" /></div>
          </div>
          <div className="flex gap-2 pt-2">
            <Button variant="primary" onClick={handleCreate}>Create</Button>
            <Button variant="secondary" onClick={() => setShowCreate(false)}>Cancel</Button>
          </div>
        </div>
      </Modal>
    </div>
  );
}
