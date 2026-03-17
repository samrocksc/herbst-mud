import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'

export const Route = createFileRoute('/_auth/rooms')({
  component: RoomsManagement,
})

interface Room {
  id: number
  name: string
  description: string
  exits: string[]
  items: string[]
}

function RoomsManagement() {
  const [rooms] = useState<Room[]>([
    { id: 1, name: 'Town Square', description: 'The central plaza', exits: ['north', 'south', 'east', 'west'], items: ['Fountain'] },
    { id: 2, name: 'Main Street', description: 'A busy street', exits: ['north', 'south'], items: ['Lamp Post'] },
    { id: 3, name: 'Forest Path', description: 'A winding path through the woods', exits: ['south', 'east'], items: [] },
  ])

  const [showForm, setShowForm] = useState(false)
  const [newRoom, setNewRoom] = useState({ name: '', description: '', exits: '', items: '' })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    const room = {
      ...newRoom,
      exits: newRoom.exits.split(',').map(e => e.trim()).filter(Boolean),
      items: newRoom.items.split(',').map(i => i.trim()).filter(Boolean)
    }
    console.log('Creating room:', room)
    setShowForm(false)
    setNewRoom({ name: '', description: '', exits: '', items: '' })
  }

  return (
    <div className="management-page">
      <div className="page-header">
        <h2>Rooms Management</h2>
        <button onClick={() => setShowForm(!showForm)}>
          {showForm ? 'Cancel' : '+ Add Room'}
        </button>
      </div>

      {showForm && (
        <div className="form-card">
          <h3>Add New Room</h3>
          <form onSubmit={handleSubmit}>
            <div className="form-row">
              <label>Name:</label>
              <input 
                type="text" 
                value={newRoom.name}
                onChange={(e) => setNewRoom({...newRoom, name: e.target.value})}
                required 
              />
            </div>
            <div className="form-row">
              <label>Description:</label>
              <textarea 
                value={newRoom.description}
                onChange={(e) => setNewRoom({...newRoom, description: e.target.value})}
              />
            </div>
            <div className="form-row">
              <label>Exits (comma-separated):</label>
              <input 
                type="text" 
                placeholder="north, south, east"
                value={newRoom.exits}
                onChange={(e) => setNewRoom({...newRoom, exits: e.target.value})}
              />
            </div>
            <div className="form-row">
              <label>Items (comma-separated):</label>
              <input 
                type="text" 
                placeholder="Fountain, Lamp"
                value={newRoom.items}
                onChange={(e) => setNewRoom({...newRoom, items: e.target.value})}
              />
            </div>
            <button type="submit">Create Room</button>
          </form>
        </div>
      )}

      <div className="rooms-grid">
        {rooms.map((room) => (
          <div key={room.id} className="room-card">
            <h4>{room.name}</h4>
            <p className="room-desc">{room.description}</p>
            <div className="room-exits">
              <strong>Exits:</strong> {room.exits.join(', ') || 'None'}
            </div>
            <div className="room-items">
              <strong>Items:</strong> {room.items.join(', ') || 'None'}
            </div>
            <div className="room-actions">
              <button>Edit</button>
              <button>View on Map</button>
              <button className="danger">Delete</button>
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}