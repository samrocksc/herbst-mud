import { createFileRoute, useNavigate, Link } from '@tanstack/react-router'
import { useEffect, useState, useCallback } from 'react'

export const Route = createFileRoute('/items')({
  component: ItemManager,
})

interface Item {
  id: number
  name: string
  description: string
  slot: string
  level: number
  weight: number
  isEquipped: boolean
  isImmovable: boolean
  color: string
  isVisible: boolean
  itemType: string
}

interface ItemForm {
  name: string
  description: string
  slot: string
  level: number
  weight: number
  isImmovable: boolean
  isVisible: boolean
  itemType: string
}

const SLOTS = ['weapon', 'head', 'chest', 'legs', 'feet', 'hands', 'accessory', 'none']
const ITEM_TYPES = ['weapon', 'armor', 'consumable', 'key', 'misc']

function ItemManager() {
  const navigate = useNavigate()
  const [items, setItems] = useState<Item[]>([])
  const [loading, setLoading] = useState(true)
  const [searchQuery, setSearchQuery] = useState('')
  const [selectedItem, setSelectedItem] = useState<Item | null>(null)
  const [editingItem, setEditingItem] = useState<Item | null>(null)
  const [saving, setSaving] = useState(false)
  const [confirmDelete, setConfirmDelete] = useState<number | null>(null)
  const [showCreateForm, setShowCreateForm] = useState(false)
  const [form, setForm] = useState<ItemForm>({
    name: '',
    description: '',
    slot: 'none',
    level: 1,
    weight: 1,
    isImmovable: false,
    isVisible: true,
    itemType: 'misc'
  })

  useEffect(() => {
    const token = localStorage.getItem('token')
    if (!token) {
      navigate({ to: '/login' })
      return
    }

    fetch('http://localhost:8080/equipment')
      .then(res => res.json())
      .then(data => {
        setItems(Array.isArray(data) ? data : [])
        setLoading(false)
      })
      .catch(err => {
        console.error('Failed to load items:', err)
        setLoading(false)
      })
  }, [navigate])

  const filteredItems = items.filter(item =>
    item.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
    item.itemType?.toLowerCase().includes(searchQuery.toLowerCase()) ||
    item.slot?.toLowerCase().includes(searchQuery.toLowerCase())
  )

  const handleCreateItem = useCallback(async () => {
    if (!form.name) {
      alert('Please enter an item name')
      return
    }

    setSaving(true)
    try {
      const token = localStorage.getItem('token')

      const response = await fetch('http://localhost:8080/equipment', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify(form)
      })

      if (!response.ok) {
        throw new Error('Failed to create item')
      }

      // Refresh item list
      const itemsResponse = await fetch('http://localhost:8080/equipment')
      const itemsData = await itemsResponse.json()
      setItems(Array.isArray(itemsData) ? itemsData : [])

      // Reset form
      setForm({
        name: '',
        description: '',
        slot: 'none',
        level: 1,
        weight: 1,
        isImmovable: false,
        isVisible: true,
        itemType: 'misc'
      })
      setShowCreateForm(false)
    } catch (err) {
      console.error('Create item error:', err)
      alert('Failed to create item')
    } finally {
      setSaving(false)
    }
  }, [form])

  const handleUpdateItem = useCallback(async () => {
    if (!editingItem) return

    setSaving(true)
    try {
      const token = localStorage.getItem('token')

      const response = await fetch(`http://localhost:8080/equipment/${editingItem.id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify(form)
      })

      if (!response.ok) {
        throw new Error('Failed to update item')
      }

      // Refresh item list
      const itemsResponse = await fetch('http://localhost:8080/equipment')
      const itemsData = await itemsResponse.json()
      setItems(Array.isArray(itemsData) ? itemsData : [])

      setEditingItem(null)
      setSelectedItem(null)
    } catch (err) {
      console.error('Update item error:', err)
      alert('Failed to update item')
    } finally {
      setSaving(false)
    }
  }, [editingItem, form])

  const handleDeleteItem = useCallback(async (itemId: number) => {
    try {
      const token = localStorage.getItem('token')

      const response = await fetch(`http://localhost:8080/equipment/${itemId}`, {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${token}`
        }
      })

      if (!response.ok) {
        throw new Error('Failed to delete item')
      }

      // Refresh item list
      const itemsResponse = await fetch('http://localhost:8080/equipment')
      const itemsData = await itemsResponse.json()
      setItems(Array.isArray(itemsData) ? itemsData : [])
      setSelectedItem(null)
      setConfirmDelete(null)
    } catch (err) {
      console.error('Delete item error:', err)
      alert('Failed to delete item')
    }
  }, [])

  const startEditing = (item: Item) => {
    setEditingItem(item)
    setForm({
      name: item.name,
      description: item.description || '',
      slot: item.slot || 'none',
      level: item.level || 1,
      weight: item.weight || 1,
      isImmovable: item.isImmovable || false,
      isVisible: item.isVisible !== false,
      itemType: item.itemType || 'misc'
    })
    setShowCreateForm(false)
  }

  if (loading) {
    return <div className="p-8 text-text">Loading items...</div>
  }

  return (
    <div className="flex h-screen bg-surface">
      {/* Left Sidebar */}
      <div className="w-[280px] bg-surface-muted border-r border-border flex flex-col">
        <div className="p-4 border-b border-border">
          <Link
            to="/dashboard"
            className="block text-primary no-underline p-2 rounded bg-surface-dark text-center mb-2 hover:bg-surface-darker"
          >
            ← Dashboard
          </Link>
          <Link
            to="/map"
            className="block text-text-muted no-underline p-2 rounded bg-surface-dark text-center mb-2 hover:bg-surface-darker"
          >
            Map Builder
          </Link>
          <Link
            to="/npcs"
            className="block text-text-muted no-underline p-2 rounded bg-surface-dark text-center hover:bg-surface-darker"
          >
            NPC Manager
          </Link>
        </div>

        <div className="p-3 border-b border-border">
          <h2 className="m-0 text-text text-lg">Item Manager</h2>
          <p className="text-text-muted text-xs mt-1 mb-0">{items.length} items</p>
        </div>

        {/* Search */}
        <div className="p-3 border-b border-border">
          <input
            type="text"
            placeholder="Search items..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
          />
        </div>

        {/* Item List */}
        <div className="flex-1 overflow-y-auto p-3">
          <div className="flex flex-col gap-1">
            {filteredItems.map(item => (
              <div
                key={item.id}
                onClick={() => {
                  setSelectedItem(item)
                  setEditingItem(null)
                  setShowCreateForm(false)
                }}
                className={`p-2 cursor-pointer rounded text-xs ${
                  selectedItem?.id === item.id ? 'text-primary bg-surface-dark' : 'text-text'
                }`}
              >
                <div className="font-bold">{item.name}</div>
                <div className="text-text-muted">
                  {item.itemType} • {item.slot || 'no slot'}
                </div>
              </div>
            ))}
            {filteredItems.length === 0 && (
              <div className="text-text-muted text-center py-4">No items found</div>
            )}
          </div>
        </div>

        {/* Create Item Button */}
        <div className="p-3 border-t border-border">
          <button
            onClick={() => {
              setShowCreateForm(true)
              setSelectedItem(null)
              setEditingItem(null)
              setForm({
                name: '',
                description: '',
                slot: 'none',
                level: 1,
                weight: 1,
                isImmovable: false,
                isVisible: true,
                itemType: 'misc'
              })
            }}
            className="w-full p-2 bg-primary border-none rounded text-white cursor-pointer hover:bg-primary-hover"
          >
            + Add Item
          </button>
        </div>
      </div>

      {/* Main Content */}
      <div className="flex-1 overflow-y-auto p-6">
        {showCreateForm ? (
          <div className="max-w-[600px] mx-auto">
            <h2 className="mt-0 mb-4 text-text">Create New Item</h2>

            <div className="bg-surface-muted rounded-lg p-4 border border-border">
              <div className="mb-4">
                <label className="text-text-muted text-xs block mb-1">Name *</label>
                <input
                  type="text"
                  value={form.name}
                  onChange={(e) => setForm({ ...form, name: e.target.value })}
                  placeholder="Item name"
                  className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
                />
              </div>

              <div className="mb-4">
                <label className="text-text-muted text-xs block mb-1">Description</label>
                <textarea
                  value={form.description}
                  onChange={(e) => setForm({ ...form, description: e.target.value })}
                  rows={3}
                  className="w-full p-2 bg-surface border border-border rounded text-text text-sm resize-y"
                />
              </div>

              <div className="grid grid-cols-2 gap-4 mb-4">
                <div>
                  <label className="text-text-muted text-xs block mb-1">Type</label>
                  <select
                    value={form.itemType}
                    onChange={(e) => setForm({ ...form, itemType: e.target.value })}
                    className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
                  >
                    {ITEM_TYPES.map(type => (
                      <option key={type} value={type}>{type}</option>
                    ))}
                  </select>
                </div>
                <div>
                  <label className="text-text-muted text-xs block mb-1">Slot</label>
                  <select
                    value={form.slot}
                    onChange={(e) => setForm({ ...form, slot: e.target.value })}
                    className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
                  >
                    {SLOTS.map(slot => (
                      <option key={slot} value={slot}>{slot}</option>
                    ))}
                  </select>
                </div>
              </div>

              <div className="grid grid-cols-2 gap-4 mb-4">
                <div>
                  <label className="text-text-muted text-xs block mb-1">Level</label>
                  <input
                    type="number"
                    value={form.level}
                    onChange={(e) => setForm({ ...form, level: parseInt(e.target.value) || 1 })}
                    min={1}
                    className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
                  />
                </div>
                <div>
                  <label className="text-text-muted text-xs block mb-1">Weight</label>
                  <input
                    type="number"
                    value={form.weight}
                    onChange={(e) => setForm({ ...form, weight: parseInt(e.target.value) || 1 })}
                    min={0}
                    className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
                  />
                </div>
              </div>

              <div className="mb-4">
                <label className="flex items-center gap-2 text-text-muted text-xs cursor-pointer">
                  <input
                    type="checkbox"
                    checked={form.isVisible}
                    onChange={(e) => setForm({ ...form, isVisible: e.target.checked })}
                    className="cursor-pointer"
                  />
                  Visible to players
                </label>
              </div>

              <div className="mb-4">
                <label className="flex items-center gap-2 text-text-muted text-xs cursor-pointer">
                  <input
                    type="checkbox"
                    checked={form.isImmovable}
                    onChange={(e) => setForm({ ...form, isImmovable: e.target.checked })}
                    className="cursor-pointer"
                  />
                  Immovable (cannot be picked up)
                </label>
              </div>

              <div className="flex gap-2">
                <button
                  onClick={handleCreateItem}
                  disabled={saving}
                  className="flex-1 p-2 bg-primary border-none rounded text-white cursor-pointer disabled:opacity-70"
                >
                  {saving ? 'Creating...' : 'Create Item'}
                </button>
                <button
                  onClick={() => setShowCreateForm(false)}
                  className="flex-1 p-2 bg-surface-dark border border-border rounded text-text-muted cursor-pointer"
                >
                  Cancel
                </button>
              </div>
            </div>
          </div>
        ) : editingItem ? (
          <div className="max-w-[600px] mx-auto">
            <h2 className="mt-0 mb-4 text-text">Edit Item</h2>

            <div className="bg-surface-muted rounded-lg p-4 border border-border">
              <div className="mb-4">
                <label className="text-text-muted text-xs block mb-1">Name *</label>
                <input
                  type="text"
                  value={form.name}
                  onChange={(e) => setForm({ ...form, name: e.target.value })}
                  className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
                />
              </div>

              <div className="mb-4">
                <label className="text-text-muted text-xs block mb-1">Description</label>
                <textarea
                  value={form.description}
                  onChange={(e) => setForm({ ...form, description: e.target.value })}
                  rows={3}
                  className="w-full p-2 bg-surface border border-border rounded text-text text-sm resize-y"
                />
              </div>

              <div className="grid grid-cols-2 gap-4 mb-4">
                <div>
                  <label className="text-text-muted text-xs block mb-1">Type</label>
                  <select
                    value={form.itemType}
                    onChange={(e) => setForm({ ...form, itemType: e.target.value })}
                    className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
                  >
                    {ITEM_TYPES.map(type => (
                      <option key={type} value={type}>{type}</option>
                    ))}
                  </select>
                </div>
                <div>
                  <label className="text-text-muted text-xs block mb-1">Slot</label>
                  <select
                    value={form.slot}
                    onChange={(e) => setForm({ ...form, slot: e.target.value })}
                    className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
                  >
                    {SLOTS.map(slot => (
                      <option key={slot} value={slot}>{slot}</option>
                    ))}
                  </select>
                </div>
              </div>

              <div className="grid grid-cols-2 gap-4 mb-4">
                <div>
                  <label className="text-text-muted text-xs block mb-1">Level</label>
                  <input
                    type="number"
                    value={form.level}
                    onChange={(e) => setForm({ ...form, level: parseInt(e.target.value) || 1 })}
                    className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
                  />
                </div>
                <div>
                  <label className="text-text-muted text-xs block mb-1">Weight</label>
                  <input
                    type="number"
                    value={form.weight}
                    onChange={(e) => setForm({ ...form, weight: parseInt(e.target.value) || 1 })}
                    className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
                  />
                </div>
              </div>

              <div className="mb-4">
                <label className="flex items-center gap-2 text-text-muted text-xs cursor-pointer">
                  <input type="checkbox" checked={form.isVisible} onChange={(e) => setForm({ ...form, isVisible: e.target.checked })} className="cursor-pointer" />
                  Visible to players
                </label>
              </div>

              <div className="mb-4">
                <label className="flex items-center gap-2 text-text-muted text-xs cursor-pointer">
                  <input type="checkbox" checked={form.isImmovable} onChange={(e) => setForm({ ...form, isImmovable: e.target.checked })} className="cursor-pointer" />
                  Immovable
                </label>
              </div>

              <div className="flex gap-2">
                <button
                  onClick={handleUpdateItem}
                  disabled={saving}
                  className="flex-1 p-2 bg-primary border-none rounded text-white cursor-pointer disabled:opacity-70"
                >
                  {saving ? 'Saving...' : 'Save Changes'}
                </button>
                <button
                  onClick={() => setEditingItem(null)}
                  className="flex-1 p-2 bg-surface-dark border border-border rounded text-text-muted cursor-pointer"
                >
                  Cancel
                </button>
              </div>
            </div>
          </div>
        ) : selectedItem ? (
          <div className="max-w-[600px] mx-auto">
            <div className="flex justify-between items-center mb-4">
              <h2 className="m-0 text-text">{selectedItem.name}</h2>
              <button
                onClick={() => setSelectedItem(null)}
                className="bg-transparent border-none text-text-muted cursor-pointer text-xl"
              >
                ×
              </button>
            </div>

            <div className="bg-surface-muted rounded-lg p-4 border border-border">
              <div className="mb-4">
                <div className="text-text-muted text-xs mb-1">ID: #{selectedItem.id}</div>
                <div className="text-text">{selectedItem.description || 'No description'}</div>
              </div>

              <div className="grid grid-cols-2 gap-4 mb-4">
                <div>
                  <label className="text-text-muted text-xs block mb-1">Type</label>
                  <div className="text-text">{selectedItem.itemType}</div>
                </div>
                <div>
                  <label className="text-text-muted text-xs block mb-1">Slot</label>
                  <div className="text-text">{selectedItem.slot || 'none'}</div>
                </div>
              </div>

              <div className="grid grid-cols-2 gap-4 mb-4">
                <div>
                  <label className="text-text-muted text-xs block mb-1">Level</label>
                  <div className="text-text">{selectedItem.level}</div>
                </div>
                <div>
                  <label className="text-text-muted text-xs block mb-1">Weight</label>
                  <div className="text-text">{selectedItem.weight}</div>
                </div>
              </div>

              <div className="grid grid-cols-2 gap-4 mb-4">
                <div>
                  <label className="text-text-muted text-xs block mb-1">Visible</label>
                  <div className="text-text">{selectedItem.isVisible ? 'Yes' : 'No'}</div>
                </div>
                <div>
                  <label className="text-text-muted text-xs block mb-1">Immovable</label>
                  <div className="text-text">{selectedItem.isImmovable ? 'Yes' : 'No'}</div>
                </div>
              </div>

              <div className="flex gap-2">
                <button
                  onClick={() => startEditing(selectedItem)}
                  className="flex-1 p-2 bg-primary border-none rounded text-white cursor-pointer hover:bg-primary-hover"
                >
                  Edit Item
                </button>
                <button
                  onClick={() => {
                    if (confirmDelete === selectedItem.id) {
                      handleDeleteItem(selectedItem.id)
                    } else {
                      setConfirmDelete(selectedItem.id)
                    }
                  }}
                  className={`flex-1 p-2 border-none rounded text-white cursor-pointer ${
                    confirmDelete === selectedItem.id
                      ? 'bg-warning hover:bg-warning/80'
                      : 'bg-danger hover:bg-danger-hover'
                  }`}
                >
                  {confirmDelete === selectedItem.id ? 'Confirm Delete?' : 'Delete Item'}
                </button>
              </div>
            </div>
          </div>
        ) : (
          <div className="flex flex-col items-center justify-center h-full text-text-muted">
            <p className="mb-2">Select an item from the list or create a new one</p>
            <p className="text-xs text-center">
              Items can be equipment, consumables, or other objects<br/>
              that players can interact with in the game world
            </p>
          </div>
        )}
      </div>
    </div>
  )
}