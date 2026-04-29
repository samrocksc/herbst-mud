import { createFileRoute, Link } from '@tanstack/react-router'
import { useEffect, useState, useCallback } from 'react'

export const Route = createFileRoute('/_auth/factions')({
  component: FactionsManagement,
})

interface FactionCategory {
  id: number
  name: string
  description?: string
}

interface Faction {
  id: number
  name: string
  description?: string
  category_id?: number
  standing: number
  members?: number[]
  is_universal?: boolean
  created_at?: string
}

interface FactionForm {
  name: string
  description: string
  category_id: number | ''
  standing: number
  is_universal: boolean
}

function FactionsManagement() {
  const [factions, setFactions] = useState<Faction[]>([])
  const [categories, setCategories] = useState<FactionCategory[]>([])
  const [loading, setLoading] = useState(true)
  const [searchQuery, setSearchQuery] = useState('')
  const [selectedFaction, setSelectedFaction] = useState<Faction | null>(null)
  const [editingFaction, setEditingFaction] = useState<Faction | null>(null)
  const [saving, setSaving] = useState(false)
  const [confirmDelete, setConfirmDelete] = useState<number | null>(null)
  const [showCreateForm, setShowCreateForm] = useState(false)
  const [showCreateCategory, setShowCreateCategory] = useState(false)
  const [categoryName, setCategoryName] = useState('')
  const [categoryDesc, setCategoryDesc] = useState('')
  const [tab, setTab] = useState<'factions' | 'categories'>('factions')
  const [form, setForm] = useState<FactionForm>({
    name: '',
    description: '',
    category_id: '',
    standing: 0,
    is_universal: false,
  })

  const fetchFactions = useCallback(async () => {
    const token = localStorage.getItem('token')
    const res = await fetch('http://localhost:8080/api/factions', {
      headers: { Authorization: `Bearer ${token}` },
    })
    const data = await res.json()
    setFactions(Array.isArray(data) ? data : [])
  }, [])

  const fetchCategories = useCallback(async () => {
    const token = localStorage.getItem('token')
    const res = await fetch('http://localhost:8080/api/faction-categories', {
      headers: { Authorization: `Bearer ${token}` },
    })
    const data = await res.json()
    setCategories(Array.isArray(data) ? data : [])
  }, [])

  useEffect(() => {
    fetchFactions()
    fetchCategories()
    setLoading(false)
  }, [fetchFactions, fetchCategories])

  const filteredFactions = factions.filter(f =>
    f.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
    (f.description?.toLowerCase().includes(searchQuery.toLowerCase()) ?? false)
  )

  const handleCreateFaction = useCallback(async () => {
    if (!form.name) {
      alert('Please enter a faction name')
      return
    }
    setSaving(true)
    try {
      const token = localStorage.getItem('token')
      const res = await fetch('http://localhost:8080/api/factions', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ ...form, category_id: form.category_id || null }),
      })
      if (!res.ok) throw new Error('Failed to create faction')
      await fetchFactions()
      setForm({ name: '', description: '', category_id: '', standing: 0, is_universal: false })
      setShowCreateForm(false)
    } catch (err) {
      console.error(err)
      alert('Failed to create faction')
    } finally {
      setSaving(false)
    }
  }, [form, fetchFactions])

  const handleUpdateFaction = useCallback(async () => {
    if (!editingFaction) return
    setSaving(true)
    try {
      const token = localStorage.getItem('token')
      const res = await fetch(`http://localhost:8080/api/factions/${editingFaction.id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ ...form, category_id: form.category_id || null }),
      })
      if (!res.ok) throw new Error('Failed to update faction')
      await fetchFactions()
      setEditingFaction(null)
      setSelectedFaction(null)
    } catch (err) {
      console.error(err)
      alert('Failed to update faction')
    } finally {
      setSaving(false)
    }
  }, [editingFaction, form, fetchFactions])

  const handleDeleteFaction = useCallback(async (id: number) => {
    try {
      const token = localStorage.getItem('token')
      const res = await fetch(`http://localhost:8080/api/factions/${id}`, {
        method: 'DELETE',
        headers: { Authorization: `Bearer ${token}` },
      })
      if (!res.ok) throw new Error('Failed to delete')
      await fetchFactions()
      setSelectedFaction(null)
      setConfirmDelete(null)
    } catch (err) {
      console.error(err)
      alert('Failed to delete faction')
    }
  }, [fetchFactions])

  const handleCreateCategory = useCallback(async () => {
    if (!categoryName) {
      alert('Please enter a category name')
      return
    }
    try {
      const token = localStorage.getItem('token')
      const res = await fetch('http://localhost:8080/api/faction-categories', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ name: categoryName, description: categoryDesc }),
      })
      if (!res.ok) throw new Error('Failed to create category')
      await fetchCategories()
      setCategoryName('')
      setCategoryDesc('')
      setShowCreateCategory(false)
    } catch (err) {
      console.error(err)
      alert('Failed to create category')
    }
  }, [categoryName, categoryDesc, fetchCategories])

  const handleDeleteCategory = useCallback(async (id: number) => {
    if (!confirm(`Delete this category? Factions in it will lose their category.`)) return
    try {
      const token = localStorage.getItem('token')
      const res = await fetch(`http://localhost:8080/api/faction-categories/${id}`, {
        method: 'DELETE',
        headers: { Authorization: `Bearer ${token}` },
      })
      if (!res.ok) throw new Error('Failed to delete')
      await fetchCategories()
    } catch (err) {
      console.error(err)
      alert('Failed to delete category')
    }
  }, [fetchCategories])

  const startEditing = (faction: Faction) => {
    setEditingFaction(faction)
    setForm({
      name: faction.name,
      description: faction.description || '',
      category_id: faction.category_id ?? '',
      standing: faction.standing,
      is_universal: faction.is_universal ?? false,
    })
    setShowCreateForm(false)
  }

  if (loading) {
    return <div className="p-8 text-text">Loading...</div>
  }

  return (
    <div className="flex h-screen bg-surface">
      {/* Left Sidebar */}
      <div className="w-[280px] bg-surface-muted border-r border-border flex flex-col">
        <div className="p-4 border-b border-border">
          <Link
            to="/dashboard"
            className="block no-underline p-2 rounded border-2 border-black text-center text-sm font-medium bg-surface-muted text-text hover:border-primary transition-colors"
          >
            ← Dashboard
          </Link>
        </div>

        {/* Tabs */}
        <div className="flex p-3 gap-2 border-b border-border">
          <button
            onClick={() => setTab('factions')}
            className={`flex-1 p-2 rounded text-sm cursor-pointer border-2 transition-colors ${
              tab === 'factions'
                ? 'bg-primary border-primary text-white font-medium'
                : 'bg-surface-muted border-border text-text-muted hover:border-primary'
            }`}
          >
            Factions
          </button>
          <button
            onClick={() => setTab('categories')}
            className={`flex-1 p-2 rounded text-sm cursor-pointer border-2 transition-colors ${
              tab === 'categories'
                ? 'bg-primary border-primary text-white font-medium'
                : 'bg-surface-muted border-border text-text-muted hover:border-primary'
            }`}
          >
            Categories
          </button>
        </div>

        {tab === 'factions' && (
          <>
            <div className="p-3 border-b border-border">
              <h2 className="m-0 text-text text-lg">Factions</h2>
              <p className="text-text-muted text-xs mt-1 mb-0">{factions.length} factions</p>
            </div>

            <div className="p-3 border-b border-border">
              <input
                type="text"
                placeholder="Search factions..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
              />
            </div>

            <div className="flex-1 overflow-y-auto p-3">
              <div className="flex flex-col gap-1">
                {filteredFactions.map(f => (
                  <div
                    key={f.id}
                    onClick={() => { setSelectedFaction(f); setEditingFaction(null); setShowCreateForm(false); }}
                    className={`p-2 cursor-pointer rounded text-xs ${
                      selectedFaction?.id === f.id ? 'text-primary bg-surface-dark' : 'text-text'
                    }`}
                  >
                    <div className="font-bold">{f.name}</div>
                    <div className="text-text-muted">
                      {f.is_universal ? 'universal' : `standing: ${f.standing}`}
                    </div>
                  </div>
                ))}
                {filteredFactions.length === 0 && (
                  <div className="text-text-muted text-center py-4">No factions found</div>
                )}
              </div>
            </div>

            <div className="p-3 border-t border-border">
              <button
                onClick={() => {
                  setShowCreateForm(true)
                  setSelectedFaction(null)
                  setEditingFaction(null)
                  setForm({ name: '', description: '', category_id: '', standing: 0, is_universal: false })
                }}
                className="w-full p-2 bg-primary border-2 border-black rounded text-white cursor-pointer hover:bg-primary-hover"
              >
                + Add Faction
              </button>
            </div>
          </>
        )}

        {tab === 'categories' && (
          <>
            <div className="p-3 border-b border-border">
              <h2 className="m-0 text-text text-lg">Categories</h2>
              <p className="text-text-muted text-xs mt-1 mb-0">{categories.length} categories</p>
            </div>

            <div className="flex-1 overflow-y-auto p-3">
              <div className="flex flex-col gap-1">
                {categories.map(c => (
                  <div
                    key={c.id}
                    className="p-2 rounded text-xs text-text"
                  >
                    <div className="font-bold">{c.name}</div>
                    <div className="text-text-muted text-xs">{c.description || '—'}</div>
                  </div>
                ))}
                {categories.length === 0 && (
                  <div className="text-text-muted text-center py-4">No categories</div>
                )}
              </div>
            </div>

            <div className="p-3 border-t border-border">
              <button
                onClick={() => setShowCreateCategory(true)}
                className="w-full p-2 bg-primary border-2 border-black rounded text-white cursor-pointer hover:bg-primary-hover"
              >
                + Add Category
              </button>
            </div>
          </>
        )}
      </div>

      {/* Main Content */}
      <div className="flex-1 overflow-y-auto p-6">
        {tab === 'factions' && showCreateForm && (
          <div className="max-w-[600px] mx-auto">
            <h2 className="mt-0 mb-4 text-text">Create Faction</h2>
            <div className="bg-surface-muted rounded-lg p-4 border border-border">
              <div className="mb-4">
                <label className="text-text-muted text-xs block mb-1">Name *</label>
                <input type="text" value={form.name}
                  onChange={e => setForm({ ...form, name: e.target.value })}
                  placeholder="Faction name"
                  className="w-full p-2 bg-surface border border-border rounded text-text text-sm" />
              </div>
              <div className="mb-4">
                <label className="text-text-muted text-xs block mb-1">Description</label>
                <textarea value={form.description} rows={3}
                  onChange={e => setForm({ ...form, description: e.target.value })}
                  className="w-full p-2 bg-surface border border-border rounded text-text text-sm resize-y" />
              </div>
              <div className="grid grid-cols-2 gap-4 mb-4">
                <div>
                  <label className="text-text-muted text-xs block mb-1">Category</label>
                  <select value={form.category_id}
                    onChange={e => setForm({ ...form, category_id: e.target.value ? Number(e.target.value) : '' })}
                    className="w-full p-2 bg-surface border border-border rounded text-text text-sm">
                    <option value="">None</option>
                    {categories.map(c => <option key={c.id} value={c.id}>{c.name}</option>)}
                  </select>
                </div>
                <div>
                  <label className="text-text-muted text-xs block mb-1">Standing</label>
                  <input type="number" value={form.standing}
                    onChange={e => setForm({ ...form, standing: parseInt(e.target.value) || 0 })}
                    className="w-full p-2 bg-surface border border-border rounded text-text text-sm" />
                </div>
              </div>
              <div className="mb-4">
                <label className="flex items-center gap-2 text-text-muted text-xs cursor-pointer">
                  <input type="checkbox" checked={form.is_universal}
                    onChange={e => setForm({ ...form, is_universal: e.target.checked })}
                    className="cursor-pointer" />
                  Universal (applies to all characters)
                </label>
              </div>
              <div className="flex gap-2">
                <button onClick={handleCreateFaction} disabled={saving}
                  className="flex-1 p-2 bg-primary border-2 border-black rounded text-white disabled:opacity-70">
                  {saving ? 'Creating...' : 'Create Faction'}
                </button>
                <button onClick={() => setShowCreateForm(false)}
                  className="flex-1 p-2 bg-surface-dark border border-border rounded text-text-muted">
                  Cancel
                </button>
              </div>
            </div>
          </div>
        )}

        {tab === 'factions' && editingFaction && !showCreateForm && (
          <div className="max-w-[600px] mx-auto">
            <h2 className="mt-0 mb-4 text-text">Edit Faction</h2>
            <div className="bg-surface-muted rounded-lg p-4 border border-border">
              <div className="mb-4">
                <label className="text-text-muted text-xs block mb-1">Name *</label>
                <input type="text" value={form.name}
                  onChange={e => setForm({ ...form, name: e.target.value })}
                  className="w-full p-2 bg-surface border border-border rounded text-text text-sm" />
              </div>
              <div className="mb-4">
                <label className="text-text-muted text-xs block mb-1">Description</label>
                <textarea value={form.description} rows={3}
                  onChange={e => setForm({ ...form, description: e.target.value })}
                  className="w-full p-2 bg-surface border border-border rounded text-text text-sm resize-y" />
              </div>
              <div className="grid grid-cols-2 gap-4 mb-4">
                <div>
                  <label className="text-text-muted text-xs block mb-1">Category</label>
                  <select value={form.category_id}
                    onChange={e => setForm({ ...form, category_id: e.target.value ? Number(e.target.value) : '' })}
                    className="w-full p-2 bg-surface border border-border rounded text-text text-sm">
                    <option value="">None</option>
                    {categories.map(c => <option key={c.id} value={c.id}>{c.name}</option>)}
                  </select>
                </div>
                <div>
                  <label className="text-text-muted text-xs block mb-1">Standing</label>
                  <input type="number" value={form.standing}
                    onChange={e => setForm({ ...form, standing: parseInt(e.target.value) || 0 })}
                    className="w-full p-2 bg-surface border border-border rounded text-text text-sm" />
                </div>
              </div>
              <div className="mb-4">
                <label className="flex items-center gap-2 text-text-muted text-xs cursor-pointer">
                  <input type="checkbox" checked={form.is_universal}
                    onChange={e => setForm({ ...form, is_universal: e.target.checked })}
                    className="cursor-pointer" />
                  Universal
                </label>
              </div>
              <div className="flex gap-2">
                <button onClick={handleUpdateFaction} disabled={saving}
                  className="flex-1 p-2 bg-primary border-2 border-black rounded text-white disabled:opacity-70">
                  {saving ? 'Saving...' : 'Save Changes'}
                </button>
                <button onClick={() => setEditingFaction(null)}
                  className="flex-1 p-2 bg-surface-dark border border-border rounded text-text-muted">
                  Cancel
                </button>
              </div>
            </div>
          </div>
        )}

        {tab === 'factions' && selectedFaction && !showCreateForm && !editingFaction && (
          <div className="max-w-[600px] mx-auto">
            <h2 className="mt-0 mb-4 text-text">{selectedFaction.name}</h2>
            <div className="bg-surface-muted rounded-lg p-4 border border-border">
              <div className="detail-row">
                <label>ID</label><span>{selectedFaction.id}</span>
              </div>
              <div className="detail-row">
                <label>Name</label><span>{selectedFaction.name}</span>
              </div>
              <div className="detail-row">
                <label>Description</label>
                <span>{selectedFaction.description || '—'}</span>
              </div>
              <div className="detail-row">
                <label>Standing</label><span>{selectedFaction.standing}</span>
              </div>
              <div className="detail-row">
                <label>Universal</label>
                <span>{selectedFaction.is_universal ? 'Yes' : 'No'}</span>
              </div>
              <div className="detail-row">
                <label>Members</label>
                <span>{selectedFaction.members?.join(', ') || 'none'}</span>
              </div>
              <div className="flex gap-2 mt-4">
                <button onClick={() => startEditing(selectedFaction)}
                  className="flex-1 p-2 bg-primary border-2 border-black rounded text-white cursor-pointer">
                  Edit
                </button>
                <button onClick={() => setConfirmDelete(selectedFaction.id)}
                  className="flex-1 p-2 bg-danger border border-border rounded text-text-muted cursor-pointer">
                  Delete
                </button>
                <button onClick={() => setSelectedFaction(null)}
                  className="flex-1 p-2 bg-surface-dark border border-border rounded text-text-muted cursor-pointer">
                  Close
                </button>
              </div>
              {confirmDelete === selectedFaction.id && (
                <div className="mt-3 p-3 bg-danger/20 border border-danger rounded">
                  <p className="text-text text-sm mb-2">Confirm deletion of "{selectedFaction.name}"?</p>
                  <div className="flex gap-2">
                    <button onClick={() => handleDeleteFaction(selectedFaction.id)}
                      className="flex-1 p-2 bg-danger text-white rounded cursor-pointer">
                      Confirm Delete
                    </button>
                    <button onClick={() => setConfirmDelete(null)}
                      className="flex-1 p-2 bg-surface-dark border border-border rounded text-text-muted">
                      Cancel
                    </button>
                  </div>
                </div>
              )}
            </div>
          </div>
        )}

        {tab === 'factions' && !showCreateForm && !editingFaction && !selectedFaction && (
          <div className="flex items-center justify-center h-full text-text-muted">
            Select a faction or create a new one
          </div>
        )}

        {tab === 'categories' && showCreateCategory && (
          <div className="max-w-[500px] mx-auto">
            <h2 className="mt-0 mb-4 text-text">Create Category</h2>
            <div className="bg-surface-muted rounded-lg p-4 border border-border">
              <div className="mb-4">
                <label className="text-text-muted text-xs block mb-1">Name *</label>
                <input type="text" value={categoryName}
                  onChange={e => setCategoryName(e.target.value)}
                  placeholder="Category name"
                  className="w-full p-2 bg-surface border border-border rounded text-text text-sm" />
              </div>
              <div className="mb-4">
                <label className="text-text-muted text-xs block mb-1">Description</label>
                <textarea value={categoryDesc} rows={2}
                  onChange={e => setCategoryDesc(e.target.value)}
                  className="w-full p-2 bg-surface border border-border rounded text-text text-sm resize-y" />
              </div>
              <div className="flex gap-2">
                <button onClick={handleCreateCategory}
                  className="flex-1 p-2 bg-primary border-2 border-black rounded text-white cursor-pointer">
                  Create Category
                </button>
                <button onClick={() => setShowCreateCategory(false)}
                  className="flex-1 p-2 bg-surface-dark border border-border rounded text-text-muted">
                  Cancel
                </button>
              </div>
            </div>
          </div>
        )}

        {tab === 'categories' && !showCreateCategory && (
          <div className="max-w-[600px] mx-auto">
            <h2 className="mt-0 mb-4 text-text">Faction Categories</h2>
            <div className="bg-surface-muted rounded-lg border border-border overflow-hidden">
              <table className="w-full text-text text-sm">
                <thead>
                  <tr className="border-b border-border bg-surface-dark">
                    <th className="text-left p-3">ID</th>
                    <th className="text-left p-3">Name</th>
                    <th className="text-left p-3">Description</th>
                    <th className="text-right p-3">Actions</th>
                  </tr>
                </thead>
                <tbody>
                  {categories.map(c => (
                    <tr key={c.id} className="border-b border-border last:border-0">
                      <td className="p-3">{c.id}</td>
                      <td className="p-3 font-bold">{c.name}</td>
                      <td className="p-3 text-text-muted">{c.description || '—'}</td>
                      <td className="p-3 text-right">
                        <button onClick={() => handleDeleteCategory(c.id)}
                          className="text-danger text-xs hover:underline cursor-pointer bg-transparent border-0 p-0">
                          Delete
                        </button>
                      </td>
                    </tr>
                  ))}
                  {categories.length === 0 && (
                    <tr>
                      <td colSpan={4} className="p-8 text-center text-text-muted">
                        No categories yet
                      </td>
                    </tr>
                  )}
                </tbody>
              </table>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
