import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'

export const Route = createFileRoute('/_auth/npcs')({
  component: NPCsManagement,
})

interface NPC {
  id: number
  name: string
  description: string
  type: string
  level: number
  behavior: string
  hp: number
  mp: number
}

function NPCsManagement() {
  const [npcs] = useState<NPC[]>([
    { id: 1, name: 'Guard Captain', description: 'Protects the town gate', type: 'guard', level: 10, behavior: 'passive', hp: 200, mp: 50 },
    { id: 2, name: 'Shopkeeper Bob', description: 'Sells supplies', type: 'merchant', level: 5, behavior: 'passive', hp: 100, mp: 100 },
    { id: 3, name: 'Forest Troll', description: 'An aggressive forest creature', type: 'enemy', level: 15, behavior: 'aggressive', hp: 300, mp: 0 },
    { id: 4, name: 'Elder Sage', description: 'Hands out quests', type: 'quest', level: 20, behavior: 'passive', hp: 150, mp: 200 },
  ])

  const [showForm, setShowForm] = useState(false)
  const [newNPC, setNewNPC] = useState({
    name: '',
    description: '',
    type: 'enemy',
    level: 1,
    behavior: 'passive',
    hp: 100,
    mp: 0
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    console.log('Creating NPC:', newNPC)
    setShowForm(false)
    setNewNPC({ name: '', description: '', type: 'enemy', level: 1, behavior: 'passive', hp: 100, mp: 0 })
  }

  const typeColors: Record<string, string> = {
    enemy: '#e74c3c',
    merchant: '#27ae60',
    quest: '#f39c12',
    guard: '#3498db',
    civilian: '#95a5a6'
  }

  return (
    <div className="management-page">
      <div className="page-header">
        <h2>NPCs Management</h2>
        <button onClick={() => setShowForm(!showForm)}>
          {showForm ? 'Cancel' : '+ Add NPC'}
        </button>
      </div>

      {showForm && (
        <div className="form-card">
          <h3>Add New NPC</h3>
          <form onSubmit={handleSubmit}>
            <div className="form-row">
              <label>Name:</label>
              <input 
                type="text" 
                value={newNPC.name}
                onChange={(e) => setNewNPC({...newNPC, name: e.target.value})}
                required 
              />
            </div>
            <div className="form-row">
              <label>Description:</label>
              <textarea 
                value={newNPC.description}
                onChange={(e) => setNewNPC({...newNPC, description: e.target.value})}
              />
            </div>
            <div className="form-row">
              <label>Type:</label>
              <select 
                value={newNPC.type}
                onChange={(e) => setNewNPC({...newNPC, type: e.target.value})}
              >
                <option value="enemy">Enemy</option>
                <option value="merchant">Merchant</option>
                <option value="quest">Quest Giver</option>
                <option value="guard">Guard</option>
                <option value="civilian">Civilian</option>
              </select>
            </div>
            <div className="form-row">
              <label>Level:</label>
              <input 
                type="number" 
                min="1" 
                max="100"
                value={newNPC.level}
                onChange={(e) => setNewNPC({...newNPC, level: parseInt(e.target.value)})}
              />
            </div>
            <div className="form-row">
              <label>Behavior:</label>
              <select 
                value={newNPC.behavior}
                onChange={(e) => setNewNPC({...newNPC, behavior: e.target.value})}
              >
                <option value="passive">Passive</option>
                <option value="aggressive">Aggressive</option>
                <option value="flee">Flee</option>
              </select>
            </div>
            <div className="form-row">
              <label>HP:</label>
              <input 
                type="number" 
                min="1"
                value={newNPC.hp}
                onChange={(e) => setNewNPC({...newNPC, hp: parseInt(e.target.value)})}
              />
            </div>
            <div className="form-row">
              <label>MP:</label>
              <input 
                type="number" 
                min="0"
                value={newNPC.mp}
                onChange={(e) => setNewNPC({...newNPC, mp: parseInt(e.target.value)})}
              />
            </div>
            <button type="submit">Create NPC</button>
          </form>
        </div>
      )}

      <div className="npcs-grid">
        {npcs.map((npc) => (
          <div key={npc.id} className="npc-card">
            <div className="npc-type" style={{ backgroundColor: typeColors[npc.type] }}>
              {npc.type}
            </div>
            <h4>{npc.name}</h4>
            <p className="npc-desc">{npc.description}</p>
            <div className="npc-stats">
              <span>Level: {npc.level}</span>
              <span>HP: {npc.hp}</span>
              <span>MP: {npc.mp}</span>
            </div>
            <div className="npc-behavior">
              <span className={`behavior-${npc.behavior}`}>{npc.behavior}</span>
            </div>
            <div className="npc-actions">
              <button>Edit</button>
              <button>Set Spawn</button>
              <button className="danger">Delete</button>
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}