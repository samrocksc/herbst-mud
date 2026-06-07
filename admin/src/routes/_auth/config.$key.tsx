/* eslint-disable functional/prefer-immutable-types */
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useState, useEffect, useMemo } from "react";
import { Button } from "../../components/Button";
import { TextareaField, FormError } from "../../components/FormFields";
import { FieldLabel } from "../../components/fields/FieldLabel";
import { PageHeader } from "../../components/PageHeader";
import { apiGet, apiPut, apiDelete } from "../../utils/apiFetch";
import { showToast } from "../../components/Toast";
import { humanizeKey, tryParseJSON, isRoomIdKey } from "./-configUtils";
import type { GameConfig } from "./-configUtils";

type RoomOption = Readonly<{ id: number; name: string }>;

export const Route = createFileRoute("/_auth/config/$key")({
  component: ConfigDetailPage,
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

function ConfigDetailPage() {
  const { key } = Route.useParams();
  const navigate = useNavigate();
  const [config, setConfig] = useState<GameConfig | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [editing, setEditing] = useState(false);
  const [value, setValue] = useState("");
  const [saving, setSaving] = useState(false);
  const [formError, setFormError] = useState<string | null>(null);

  useEffect(() => {
    setLoading(true);
    apiGet<GameConfig>(`/api/game-configs/${key}`)
      .then((config) => {
        setConfig(config);
        setValue(config.value);
      })
      .catch((e) => {
        const msg = e instanceof Error ? e.message : "Failed to load config";
        setError(`Config "${key}" not found: ${msg}`);
      })
      .finally(() => setLoading(false));
  }, [key]);

  const roomsQuery = useQuery<RoomOption[]>({
    queryKey: ["rooms-list"],
    queryFn: () => apiGet<RoomOption[]>("/api/rooms"),
    enabled: isRoomIdKey(key),
  });

  const roomOptions = useMemo(() => {
    if (!roomsQuery.data) return [];
    return roomsQuery.data.map((r) => ({ id: r.id, name: r.name }));
  }, [roomsQuery.data]);

  const handleSave = async () => {
    if (!config) return;
    setSaving(true);
    setFormError(null);
    try {
      const updated = await apiPut<GameConfig>(`/api/game-configs/${config.key}`, { value });
      showToast("Config updated.", "success");
      setConfig(updated);
      setEditing(false);
    } catch (e) {
      setFormError(e instanceof Error ? e.message : "Unknown error");
    } finally {
      setSaving(false);
    }
  };

  const handleDelete = async () => {
    if (!config) return;
    try {
      await apiDelete(`/api/game-configs/${config.key}`);
      showToast("Config deleted.", "success");
      navigate({ to: "/config" });
    } catch (e) {
      setFormError(e instanceof Error ? e.message : "Unknown error");
    }
  };

  if (loading) return <div className="p-8"><PageHeader title="Loading..." backTo="/config" /><div className="text-text-muted">Loading config...</div></div>;
  if (error || !config) return <div className="p-8"><PageHeader title="Error" backTo="/config" /><div className="text-danger">{error ?? "Config not found"}</div></div>;

  const parsed = tryParseJSON(config.value);

  return (
    <div className="p-6 max-w-[600px] mx-auto">
      <PageHeader title={humanizeKey(config.key)} showBack backTo="/config" actions={
        <div className="flex gap-2">
          <Button variant="danger" size="sm" onClick={handleDelete}>Delete</Button>
          <Button variant="primary" size="sm" onClick={() => setEditing(!editing)}>
            {editing ? "Cancel" : "Edit"}
          </Button>
        </div>
      } />
      <div className="bg-surface p-6 border border-border rounded">
        {formError && <FormError message={formError} />}
        <div className="space-y-4">
          <div>
            <label className="text-text-muted text-sm block mb-1">Key</label>
            <code className="block p-2 bg-surface-muted rounded text-text text-sm">{config.key}</code>
          </div>
          {editing ? (
            <div>
              <label className="text-text-muted text-sm block mb-1">Value</label>
              <CollapsibleJSONPreview value={value} />
              {isRoomIdKey(config.key) ? (
                <div>
                  <FieldLabel>Room</FieldLabel>
                  <select
                    value={value}
                    onChange={(e) => setValue(e.target.value)}
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
                <TextareaField label="" value={value} onChange={setValue} rows={6} />
              )}
              <div className="flex gap-2 mt-3">
                <Button variant="primary" size="sm" onClick={handleSave} disabled={saving}>
                  {saving ? "Saving..." : "Save"}
                </Button>
                <Button variant="secondary" size="sm" onClick={() => { setEditing(false); setValue(config.value); }}>
                  Cancel
                </Button>
              </div>
            </div>
          ) : (
            <div>
              <label className="text-text-muted text-sm block mb-1">Value</label>
              {parsed !== null ? (
                <pre className="bg-surface-muted border border-border rounded p-3 text-xs font-mono whitespace-pre-wrap overflow-auto max-h-64 text-text">{JSON.stringify(parsed, null, 2)}</pre>
              ) : (
                <div className="p-2 bg-surface-muted rounded text-text text-sm">{config.value}</div>
              )}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
