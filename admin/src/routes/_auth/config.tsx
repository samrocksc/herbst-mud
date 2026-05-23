/* eslint-disable functional/prefer-immutable-types */
import { createFileRoute, Link } from "@tanstack/react-router";
import { useState, useCallback, useEffect } from "react";
import type { ReactNode } from "react";
import { Button } from "../../components/Button";
import { DataTable } from "../../components/DataTable";
import { showToast } from "../../components/Toast";
import { PageContainer } from "../../components/PageContainer";
import { apiGet, apiDelete } from "../../utils/apiFetch";
import { humanizeKey, tryParseJSON } from "./-configUtils";
import type { GameConfig } from "./-configUtils";

export const Route = createFileRoute("/_auth/config")({ component: ConfigManagement });

function ConfigValueCell({ value }: { value: string }) {
  const parsed = tryParseJSON(value);
  const [expanded, setExpanded] = useState(false);
  if (parsed !== null) {
    const formatted = JSON.stringify(parsed, null, 2);
    const isLong = formatted.split("\n").length > 4;
    return (
      <PageContainer>
        <pre className={`font-mono text-text-secondary whitespace-pre-wrap m-0 ${!expanded ? "max-h-16 overflow-hidden" : ""}`}>{formatted}</pre>
        {isLong && <button type="button" className="text-primary text-xs mt-1 hover:underline cursor-pointer" onClick={() => setExpanded(e => !e)}>{expanded ? "Show less" : "Show more"}</button>}
      </PageContainer>
    );
  }
  return <span className="inline-block max-w-md overflow-hidden text-ellipsis whitespace-nowrap text-text-secondary text-xs">{value.length > 60 ? value.slice(0, 60) + "…" : value}</span>;
}

function ConfigManagement() {
  const [configs, setConfigs] = useState<GameConfig[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [search, setSearch] = useState("");
  const [deleteTarget, setDeleteTarget] = useState<GameConfig | null>(null);

  const fetchConfigs = useCallback(async () => {
    setLoading(true); setError(null);
    try { setConfigs(await apiGet<GameConfig[]>("/api/game-configs")); }
    catch (e: unknown) { setError(e instanceof Error ? e.message : "Unknown error"); }
    finally { setLoading(false); }
  }, []);

  useEffect(() => { fetchConfigs(); }, [fetchConfigs]);

  const handleDelete = async () => {
    if (!deleteTarget) return;
    try { await apiDelete(`/api/game-configs/${deleteTarget.key}`); showToast("Config deleted.", "success"); setDeleteTarget(null); fetchConfigs(); }
    catch (e: unknown) { showToast(`Failed to delete: ${e instanceof Error ? e.message : "Unknown error"}`); }
  };

  const filtered = configs.filter(c => c.key.toLowerCase().includes(search.toLowerCase()) || c.value.toLowerCase().includes(search.toLowerCase()));

  return (
    <div >
      <div className="page-header">
        <h2>Game Configs</h2>
        <Link to="/config/new" className="no-underline">
          <Button variant="primary">+ New Config</Button>
        </Link>
      </div>
      {error && <div className="error-banner">{error}</div>}
      <div className="flex items-center gap-3 mb-4">
        <input type="text" placeholder="Search keys or values..." value={search} onChange={e => setSearch(e.target.value)}
          className="flex-1 px-3 py-2 bg-surface-muted border-2 border-border color-text rounded" />
        <Button variant="secondary" onClick={fetchConfigs}>Refresh</Button>
      </div>
      {loading ? <div className="loading">Loading configs...</div> : (
        <DataTable columns={[
          { header: "Key", accessor: "key", render: (_, row): ReactNode => (
            <Link to="/config/$key" params={{ key: row.key }} className="no-underline">
              <div><span className="text-text text-sm font-medium">{humanizeKey(row.key)}</span><br /><code className="text-primary text-xs">{row.key}</code></div>
            </Link>
          )},
          { header: "Value", accessor: "value", render: (val) => <ConfigValueCell value={val as string} /> },
          { header: "Actions", accessor: "_actions", render: (_, row): ReactNode => (
            <div className="flex gap-2">
              <Link to="/config/$key" params={{ key: row.key }} className="no-underline">
                <Button variant="accent" size="sm">Edit</Button>
              </Link>
              <Button variant="danger" size="sm" onClick={() => setDeleteTarget(row)}>Delete</Button>
            </div>
          )},
        ]} data={filtered} getKey={row => row.id}
          emptyMessage={configs.length === 0 ? "No configs found. Create one below." : "No configs match your search."} />
      )}
      {deleteTarget && (
        <div className="modal-overlay" onClick={() => setDeleteTarget(null)}>
          <div className="modal-content max-w-md" onClick={e => e.stopPropagation()}>
            <h3>Delete Config?</h3>
            <p>Are you sure you want to delete <code>{deleteTarget.key}</code>? This cannot be undone.</p>
            <div className="flex gap-3 justify-end mt-4">
              <Button variant="secondary" onClick={() => setDeleteTarget(null)}>Cancel</Button>
              <Button variant="danger" onClick={handleDelete}>Delete</Button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
