/* eslint-disable functional/prefer-immutable-types */
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useState, useMemo } from "react";
import { useQuery } from "@tanstack/react-query";
import { Button } from "../../components/Button";
import { FormField, TextareaField, FormError } from "../../components/FormFields";
import { FieldLabel } from "../../components/fields/FieldLabel";
import { PageHeader } from "../../components/PageHeader";
import { apiGet, apiPost } from "../../utils/apiFetch";
import { showToast } from "../../components/Toast";
import { tryParseJSON, PRESETS, isRoomIdKey } from "./-configUtils";

type RoomOption = Readonly<{ id: number; name: string }>;

export const Route = createFileRoute("/_auth/config/new")({
  component: CreateConfigPage,
});

 
function CollapsibleJSONPreview({ value }: { value: string }) {
  const parsed = tryParseJSON(value);
  const [expanded, setExpanded] = useState(false);
  if (parsed === null) return null;
  const formatted = JSON.stringify(parsed, null, 2);
  return (
    <div className="mb-3">
      <button type="button" className="text-xs text-primary hover:underline cursor-pointer flex items-center gap-1 mb-1"
        onClick={() => setExpanded(e => !e)}>
        <span className={`inline-block transition-transform ${expanded ? "rotate-90" : ""}`}>&#9654;</span>
        {expanded ? "Collapse" : "Expand"} JSON preview
      </button>
      {expanded && <pre className="bg-surface-muted border-2 border-border rounded p-3 text-xs font-mono whitespace-pre-wrap overflow-auto max-h-64">{formatted}</pre>}
    </div>
  );
}

function CreateConfigPage() {
  const navigate = useNavigate();
  const [form, setForm] = useState({ key: "", value: "" });
  const [saving, setSaving] = useState(false);
  const [formError, setFormError] = useState("");

  const roomsQuery = useQuery<RoomOption[]>({
    queryKey: ["rooms-list"],
    queryFn: () => apiGet<RoomOption[]>("/api/rooms"),
    enabled: isRoomIdKey(form.key),
  });

  const roomOptions = useMemo(() => {
    if (!roomsQuery.data) return [];
    return roomsQuery.data.map((r) => ({ id: r.id, name: r.name }));
  }, [roomsQuery.data]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSaving(true);
    setFormError("");
    try {
      await apiPost("/api/game-configs", form);
      showToast("Config created.", "success");
      navigate({ to: "/config" });
    } catch (err) {
      setFormError(err instanceof Error ? err.message : "Unknown error");
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="p-6 max-w-[600px] mx-auto">
      <PageHeader title="New Game Config" showBack backTo="/config" />
      <div className="bg-surface p-6 border border-border rounded">
        {formError && <FormError message={formError} />}
        <form onSubmit={handleSubmit}>
          <div className="space-y-4">
            <div>
              <label className="block text-text-muted text-sm mb-1">Key</label>
              <FormField label="" value={form.key} onChange={v => setForm(f => ({ ...f, key: v }))} placeholder="e.g. xp_thresholds" required />
            </div>
            <div>
              <label className="block text-text-muted text-sm mb-1">Value {isRoomIdKey(form.key) ? "(Room ID)" : "(JSON or plain)"}</label>
              <CollapsibleJSONPreview value={form.value} />
              {isRoomIdKey(form.key) ? (
                <div>
                  <FieldLabel>Room</FieldLabel>
                  <select
                    value={form.value}
                    onChange={(e) => setForm(f => ({ ...f, value: e.target.value }))}
                    className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
                  >
                    <option value="">-- Select a room --</option>
                    {roomOptions.map((r) => (
                      <option key={r.id} value={String(r.id)}>
                        {r.name} (#{r.id})
                      </option>
                    ))}
                  </select>
                  {roomsQuery.isLoading && (
                    <div className="text-xs text-text-muted mt-1">Loading rooms...</div>
                  )}
                </div>
              ) : (
                <TextareaField label="" value={form.value} onChange={v => setForm(f => ({ ...f, value: v }))} rows={6} placeholder='{"key": "value"} or plain text' required />
              )}
            </div>
            <div>
              <label className="block text-text-muted text-sm mb-1">Presets</label>
              <div className="flex flex-wrap gap-2">
                {PRESETS.map(p => (
                  <Button type="button" key={p.key} variant="ghost" size="sm"
                    onClick={() => setForm({ key: p.key, value: p.value })}>
                    {p.label}
                  </Button>
                ))}
              </div>
            </div>
            <div className="flex gap-3 justify-end pt-2">
              <Button type="button" variant="secondary" onClick={() => navigate({ to: "/config" })}>Cancel</Button>
              <Button type="submit" variant="primary" disabled={saving}>{saving ? "Saving..." : "Create"}</Button>
            </div>
          </div>
        </form>
      </div>
    </div>
  );
}
