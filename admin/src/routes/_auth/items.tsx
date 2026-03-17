import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'

export const Route = createFileRoute('/_auth/items')({
  component: ItemsManagement,
})

interface Item {
  id: number
  name: string
  description: string
  type: string
  immovable: boolean
  color: string
}

function ItemsManagement() {
  const [items] = useState<Item[]>([
    { id: 1, name: 'Rusty Sword', description: 'A worn sword', type: 'weapon', immovable: false, color: '#8B4513' },
    { id: 2, name: 'Healing Potion', description: 'Restores health', type: 'consumable', immovable: false, color: '#FF0000' },
    { id: 3, name: 'Ancient Stone', description: 'A heavy stone', type: 'decoration', immovable: true, color: '#808080' },
  ])

  const [showForm, setShowForm] = useState(false)
  const [newItem, setNewItem] = useState({ name: '', description: '', type: 'weapon', immovable: false, color: '#ffffff' })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    console.log('Creating item:', newItem)
    setShowForm(false)
    setNewItem({ name: '', description: '', type: 'weapon', immovable: false, color: '#ffffff' })
  }

  return (
    <div className="management-page">
      <div className="page-header">
        <h2>Items Management</h2>
        <button onClick={() => setShowForm(!showForm)}>
          {showForm ? 'Cancel' : '+ Add Item'}
        </button>
      </div>

      {showForm && (
        <div className="form-card">
          <h3>Add New Item</h3>
          <form onSubmit={handleSubmit}>
            <div className="form-row">
              <label>Name:</label>
              <input 
                type="text" 
                value={newItem.name}
                onChange={(e) => setNewItem({...newItem, name: e.target.value})}
                required 
              />
            </div>
            <div className="form-row">
              <label>Description:</label>
              <textarea 
                value={newItem.description}
                onChange={(e) => setNewItem({...newItem, description: e.target.value})}
              />
            </div>
            <div className="form-row">
              <label>Type:</label>
              <select 
                value={newItem.type}
                onChange={(e) => setNewItem({...newItem, type: e.target.value})}
              >
                <option value="weapon">Weapon</option>
                <option value="armor">Armor</option>
                <option value="consumable">Consumable</option>
                <option value="decoration">Decoration</option>
                <option value="quest">Quest Item</option>
              </select>
            </div>
            <div className="form-row">
              <label>Color (hex):</label>
              <input 
                type="text" 
                value={newItem.color}
                onChange={(e) => setNewItem({...newItem, color: e.target.value})}
              />
              <input 
                type="color" 
                value={newItem.color}
                onChange={(e) => setNewItem({...newItem, color: e.target.value})}
              />
            </div>
            <div className="form-row checkbox">
              <label>
                <input 
                  type="checkbox"
                  checked={newItem.immovable}
                  onChange={(e) => setNewItem({...newItem, immovable: e.target.checked})}
                />
                Immovable (cannot be picked up)
              </label>
            </div>
            <button type="submit">Create Item</button>
          </form>
        </div>
      )}

      <div className="items-grid">
        {items.map((item) => (
          <div key={item.id} className="item-card">
            <div className="item-color" style={{ backgroundColor: item.color }}></div>
            <h4>{item.name}</h4>
            <p className="item-type">{item.type}</p>
            <p className="item-desc">{item.description}</p>
            {item.immovable && <span className="badge">Immovable</span>}
            <div className="item-actions">
              <button>Edit</button>
              <button className="danger">Delete</button>
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}