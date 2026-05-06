import { createFileRoute, Link } from '@tanstack/react-router'
import { logError } from '../../utils/log'
import { useEffect, useState, useCallback } from 'react'
import { Button } from '../../components/Button'
import { DataTable, type Column } from '../../components/DataTable'
import {
  FormField,
  NumberField,
  TextareaField,
  SelectField,
  CheckboxField,
} from '../../components/FormFields'

export const Route = createFileRoute('/_auth/factions')({
  component: FactionsManagement,
})

type FactionCategory = Readonly<{
  id: number
  name: string
  description?: string
}>

type Faction = Readonly<{
  id: number
  name: string
  description?: string
  category_id?: number
  standing: number
  members?: number[]
  is_universal?: boolean
  created_at?: string
}>

type FactionForm = Readonly<{
  name: string
  description: string
  category_id: number | ''
  standing: number
  is_universal: boolean
}>

const EMPTY_FORM: FactionForm = {
  name: '',
  description: '',
  category_id: '',
  standing: 0,
  is_universal: false,
}

function FactionFormFields({
  form,
  setForm,
  categories,
}: Readonly<{
  form: FactionForm
  setForm: (f: FactionForm) => void
  categories: FactionCategory[]
}>) {
  const set = (patch: Partial<FactionForm>) => setForm({ ...form, ...patch })
  const catOptions = [
    { value: '', label: 'None' },
    ...categories.map((c) => ({ value: String(c.id), label: c.name })),
  ]

  return (
    <div className="bg-surface-muted rounded-lg p-4 border border-border space-y-3">
      <FormField
        label="Name"
        value={form.name}
        onChange={(v) => set({ name: v })}
        placeholder="Faction name"
      />

      <TextareaField
        label="Description"
        value={form.description}
        onChange={(v) => set({ description: v })}
        rows={3}
      />

      <div className="grid grid-cols-2 gap-3">
        <SelectField
          label="Category"
          value={String(form.category_id)}
          onChange={(v) => set({ category_id: v ? Number(v) : '' })}
          options={catOptions}
        />
        <NumberField
          label="Standing"
          value={form.standing}
          onChange={(v) => set({ standing: v })}
        />
      </div>

      <CheckboxField
        label="Universal (applies to all characters)"
        checked={form.is_universal}
        onChange={(v) => set({ is_universal: v })}
      />
    </div>
  )
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
  const [form, setForm] = useState<FactionForm>(EMPTY_FORM)

  const fetchFactions = useCallback(async () => {
    const token = localStorage.getItem('token')
    const res = await fetch(`${window.location.origin}/api/factions`, {
      headers: { Authorization: `Bearer ${token}` },
    })
    const data = await res.json()
    setFactions(Array.isArray(data) ? data : [])
  }, [])

  const fetchCategories = useCallback(async () => {
    const token = localStorage.getItem('token')
    const res = await fetch(`${window.location.origin}/api/faction-categories`, {
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

  const filteredFactions = factions.filter(
    (f) =>
      f.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      (f.description?.toLowerCase().includes(searchQuery.toLowerCase()) ?? false),
  )

  const handleCreateFaction = useCallback(async () => {
    if (!form.name) {
      alert('Please enter a faction name')
      return
    }
    setSaving(true)
    try {
      const token = localStorage.getItem('token')
      const res = await fetch(`${window.location.origin}/api/factions`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ ...form, category_id: form.category_id || null }),
      })
      if (!res.ok) throw new Error('Failed to create faction')
      await fetchFactions()
      setForm(EMPTY_FORM)
      setShowCreateForm(false)
    } catch (err) {
      logError('Error', err)
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
      const res = await fetch(
        `${window.location.origin}/api/factions/${editingFaction.id}`,
        {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({ ...form, category_id: form.category_id || null }),
        },
      )
      if (!res.ok) throw new Error('Failed to update faction')
      await fetchFactions()
      setEditingFaction(null)
      setSelectedFaction(null)
    } catch (err) {
      logError('Error', err)
      alert('Failed to update faction')
    } finally {
      setSaving(false)
    }
  }, [editingFaction, form, fetchFactions])

  const handleDeleteFaction = useCallback(
    async (id: number) => {
      try {
        const token = localStorage.getItem('token')
        const res = await fetch(`${window.location.origin}/api/factions/${id}`, {
          method: 'DELETE',
          headers: { Authorization: `Bearer ${token}` },
        })
        if (!res.ok) throw new Error('Failed to delete')
        await fetchFactions()
        setSelectedFaction(null)
        setConfirmDelete(null)
      } catch (err) {
        logError('Error', err)
        alert('Failed to delete faction')
      }
    },
    [fetchFactions],
  )

  const handleCreateCategory = useCallback(async () => {
    if (!categoryName) {
      alert('Please enter a category name')
      return
    }
    try {
      const token = localStorage.getItem('token')
      const res = await fetch(`${window.location.origin}/api/faction-categories`, {
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
      logError('Error', err)
      alert('Failed to create category')
    }
  }, [categoryName, categoryDesc, fetchCategories])

  const handleDeleteCategory = useCallback(
    async (id: number) => {
      if (!confirm(`Delete this category? Factions in it will lose their category.`)) return
      try {
        const token = localStorage.getItem('token')
        const res = await fetch(`${window.location.origin}/api/faction-categories/${id}`, {
          method: 'DELETE',
          headers: { Authorization: `Bearer ${token}` },
        })
        if (!res.ok) throw new Error('Failed to delete')
        await fetchCategories()
      } catch (err) {
        logError('Error', err)
        alert('Failed to delete category')
      }
    },
    [fetchCategories],
  )

  const categoryColumns: Column<FactionCategory>[] = [
    { header: 'ID', accessor: 'id' },
    { header: 'Name', accessor: 'name', className: 'font-bold' },
    { header: 'Description', accessor: 'description' },
    {
      header: '',
      accessor: 'id',
      align: 'right',
      render: (_, row) => (
        <Button variant="danger" size="sm" onClick={() => handleDeleteCategory(row.id)}>
          Delete
        </Button>
      ),
    },
  ]

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

        <div className="flex p-3 gap-2 border-b border-border">
          <Button
            variant={tab === 'factions' ? 'primary' : 'secondary'}
            size="sm"
            fullWidth
            onClick={() => setTab('factions')}
          >
            Factions
          </Button>
          <Button
            variant={tab === 'categories' ? 'primary' : 'secondary'}
            size="sm"
            fullWidth
            onClick={() => setTab('categories')}
          >
            Categories
          </Button>
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
                {filteredFactions.map((f) => (
                  <div
                    key={f.id}
                    onClick={() => {
                      setSelectedFaction(f)
                      setEditingFaction(null)
                      setShowCreateForm(false)
                    }}
                    className={`p-2 cursor-pointer rounded text-xs ${
                      selectedFaction?.id === f.id
                        ? 'text-primary bg-primary/20 font-medium'
                        : 'text-text'
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
              <Button
                variant="primary"
                size="md"
                fullWidth
                onClick={() => {
                  setShowCreateForm(true)
                  setSelectedFaction(null)
                  setEditingFaction(null)
                  setForm(EMPTY_FORM)
                }}
              >
                + Add Faction
              </Button>
            </div>
          </>
        )}

        {tab === 'categories' && (
          <>
            <div className="p-3 border-b border-border">
              <h2 className="m-0 text-text text-lg">Categories</h2>
              <p className="text-text-muted text-xs mt-1 mb-0">
                {categories.length} categories
              </p>
            </div>
            <div className="flex-1 overflow-y-auto p-3">
              <div className="flex flex-col gap-1">
                {categories.map((c) => (
                  <div key={c.id} className="p-2 rounded text-xs text-text">
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
              <Button
                variant="primary"
                size="md"
                fullWidth
                onClick={() => setShowCreateCategory(true)}
              >
                + Add Category
              </Button>
            </div>
          </>
        )}
      </div>

      {/* Main Content */}
      <div className="flex-1 overflow-y-auto p-6">
        {tab === 'factions' && showCreateForm && (
          <div className="max-w-[600px] mx-auto">
            <h2 className="mt-0 mb-4 text-text">Create Faction</h2>
            <FactionFormFields form={form} setForm={setForm} categories={categories} />
            <div className="flex gap-2 mt-3">
              <Button variant="primary" size="md" fullWidth onClick={handleCreateFaction} disabled={saving}>
                {saving ? 'Creating...' : 'Create Faction'}
              </Button>
              <Button variant="secondary" size="md" fullWidth onClick={() => setShowCreateForm(false)}>
                Cancel
              </Button>
            </div>
          </div>
        )}

        {tab === 'factions' && editingFaction && !showCreateForm && (
          <div className="max-w-[600px] mx-auto">
            <h2 className="mt-0 mb-4 text-text">Edit Faction</h2>
            <FactionFormFields form={form} setForm={setForm} categories={categories} />
            <div className="flex gap-2 mt-3">
              <Button variant="primary" size="md" fullWidth onClick={handleUpdateFaction} disabled={saving}>
                {saving ? 'Saving...' : 'Save Changes'}
              </Button>
              <Button variant="secondary" size="md" fullWidth onClick={() => setEditingFaction(null)}>
                Cancel
              </Button>
            </div>
          </div>
        )}

        {tab === 'factions' && selectedFaction && !showCreateForm && !editingFaction && (
          <div className="max-w-[600px] mx-auto">
            <h2 className="mt-0 mb-4 text-text">{selectedFaction.name}</h2>
            <div className="bg-surface-muted rounded-lg p-4 border border-border space-y-2">
              <DetailRow label="ID" value={String(selectedFaction.id)} />
              <DetailRow label="Name" value={selectedFaction.name} />
              <DetailRow label="Description" value={selectedFaction.description || '—'} />
              <DetailRow label="Standing" value={String(selectedFaction.standing)} />
              <DetailRow label="Universal" value={selectedFaction.is_universal ? 'Yes' : 'No'} />
              <DetailRow label="Members" value={selectedFaction.members?.join(', ') || 'none'} />
              <div className="flex gap-2 mt-3">
                <Button variant="primary" size="md" fullWidth onClick={() => startEditing(selectedFaction)}>
                  Edit
                </Button>
                <Button
                  variant={confirmDelete === selectedFaction.id ? 'secondary' : 'danger'}
                  size="md"
                  fullWidth
                  onClick={() => {
                    if (confirmDelete === selectedFaction.id) {
                      handleDeleteFaction(selectedFaction.id)
                    } else {
                      setConfirmDelete(selectedFaction.id)
                    }
                  }}
                >
                  {confirmDelete === selectedFaction.id ? 'Confirm Delete?' : 'Delete'}
                </Button>
                <Button variant="secondary" size="md" fullWidth onClick={() => setSelectedFaction(null)}>
                  Close
                </Button>
              </div>
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
            <div className="bg-surface-muted rounded-lg p-4 border border-border space-y-3">
              <FormField
                label="Name"
                value={categoryName}
                onChange={setCategoryName}
                placeholder="Category name"
              />
              <TextareaField
                label="Description"
                value={categoryDesc}
                onChange={setCategoryDesc}
                rows={2}
              />
              <div className="flex gap-2">
                <Button variant="primary" size="md" fullWidth onClick={handleCreateCategory}>
                  Create Category
                </Button>
                <Button variant="secondary" size="md" fullWidth onClick={() => setShowCreateCategory(false)}>
                  Cancel
                </Button>
              </div>
            </div>
          </div>
        )}

        {tab === 'categories' && !showCreateCategory && (
          <div className="max-w-[600px] mx-auto">
            <h2 className="mt-0 mb-4 text-text">Faction Categories</h2>
            <DataTable
              columns={categoryColumns}
              data={categories}
              getKey={(c) => c.id}
              emptyMessage="No categories yet"
            />
          </div>
        )}
      </div>
    </div>
  )
}

function DetailRow({ label, value }: Readonly<{ label: string; value: string }>) {
  return (
    <div className="detail-row">
      <label>{label}</label>
      <span>{value}</span>
    </div>
  )
}
