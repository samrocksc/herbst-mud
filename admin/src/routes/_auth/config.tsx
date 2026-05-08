import { createFileRoute } from '@tanstack/react-router'
import { useState, useCallback, useEffect } from 'react'
import type { ReactNode } from 'react'
import { Button } from '../../components/Button'
import { DataTable } from '../../components/DataTable'
import { showToast } from '../../components/Toast'
import { apiGet, apiDelete } from '../../utils/apiFetch'
import { ConfigForm } from './ConfigForm'
import { ConfigValueCell, DeleteConfigModal, humanizeKey } from './ConfigHelpers'
import type { GameConfig } from './ConfigHelpers'

export const Route = createFileRoute('/_auth/config')({ component: ConfigManagement })

function ConfigManagement() {
  const [configs, setConfigs] = useState<GameConfig[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [search, setSearch] = useState('')
  const [activeForm, setActiveForm] = useState<'create' | 'edit' | null>(null)
  const [editing, setEditing] = useState<GameConfig | null>(null)
  const [deleteTarget, setDeleteTarget] = useState<GameConfig | null>(null)

  const fetchConfigs = useCallback(async () => {
    setLoading(true); setError(null)
    try { setConfigs(await apiGet<GameConfig[]>('/api/game-configs')) }
    catch (e: unknown) { setError(e instanceof Error ? e.message : 'Unknown error') }
    finally { setLoading(false) }
  }, [])

  useEffect(() => { fetchConfigs() }, [fetchConfigs])

  const handleDelete = async () => {
    if (!deleteTarget) return
    try {
      await apiDelete(`/api/game-configs/${deleteTarget.key}`)
      showToast('Config deleted.', 'success')
      setDeleteTarget(null); fetchConfigs()
    } catch (e: unknown) {
      showToast(`Failed to delete: ${e instanceof Error ? e.message : 'Unknown error'}`)
    }
  }

  const filtered = configs.filter(c =>
    c.key.toLowerCase().includes(search.toLowerCase()) ||
    c.value.toLowerCase().includes(search.toLowerCase())
  )

  const closeForm = () => { setActiveForm(null); setEditing(null); fetchConfigs() }

  return (
    <div className="management-page">
      <div className="page-header">
        <h2>Game Configs</h2>
        <Button variant="primary" onClick={() => { setActiveForm('create'); setEditing(null) }}>+ New Config</Button>
      </div>
      {error && <div className="error-banner">{error}</div>}
      <div className="flex items-center gap-3 mb-4">
        <input type="text" placeholder="Search keys or values..." value={search}
          onChange={e => setSearch(e.target.value)}
          className="flex-1 px-3 py-2 bg-surface-muted border-2 border-border color-text rounded" />
        <Button variant="secondary" onClick={fetchConfigs}>Refresh</Button>
      </div>
      {loading ? <div className="loading">Loading configs...</div> : (
        <DataTable
          columns={[
            { header: 'Key', accessor: 'key', render: (_, row): ReactNode => (
              <div><span className="text-text text-sm font-medium">{humanizeKey(row.key)}</span><br />
              <code className="text-primary text-xs">{row.key}</code></div>)},
            { header: 'Value', accessor: 'value', render: (val) => <ConfigValueCell value={val as string} /> },
            { header: 'Actions', accessor: '_actions', render: (_, row): ReactNode => (
              <div className="flex gap-2">
                <Button variant="accent" size="sm" onClick={() => { setEditing(row); setActiveForm('edit') }}>Edit</Button>
                <Button variant="danger" size="sm" onClick={() => setDeleteTarget(row)}>Delete</Button>
              </div>)},
          ]}
          data={filtered} getKey={row => row.id}
          emptyMessage={configs.length === 0 ? 'No configs found. Create one below.' : 'No configs match your search.'}
        />
      )}
      {activeForm && <ConfigForm editing={editing} onDone={closeForm} />}
      {deleteTarget && <DeleteConfigModal target={deleteTarget} onConfirm={handleDelete} onCancel={() => setDeleteTarget(null)} />}
    </div>
  )
}