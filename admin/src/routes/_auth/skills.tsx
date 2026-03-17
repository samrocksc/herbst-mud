import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'

export const Route = createFileRoute('/_auth/skills')({
  component: SkillsManagement,
})

interface Skill {
  id: number
  name: string
  description: string
  category: string
  maxLevel: number
}

function SkillsManagement() {
  const [skills] = useState<Skill[]>([
    { id: 1, name: 'Sword Slash', description: 'Basic sword attack', category: 'Combat', maxLevel: 10 },
    { id: 2, name: 'Fireball', description: 'Launch a fireball', category: 'Magic', maxLevel: 5 },
    { id: 3, name: 'Heal', description: 'Restore HP', category: 'Magic', maxLevel: 7 },
    { id: 4, name: 'Sneak', description: 'Move silently', category: 'Stealth', maxLevel: 5 },
    { id: 5, name: 'Pick Lock', description: 'Open locks without keys', category: 'Stealth', maxLevel: 3 },
  ])

  const [showForm, setShowForm] = useState(false)
  const [newSkill, setNewSkill] = useState({ name: '', description: '', category: 'Combat', maxLevel: 5 })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    console.log('Creating skill:', newSkill)
    setShowForm(false)
    setNewSkill({ name: '', description: '', category: 'Combat', maxLevel: 5 })
  }

  const categories = [...new Set(skills.map(s => s.category))]

  return (
    <div className="management-page">
      <div className="page-header">
        <h2>Skills & Talents Management</h2>
        <button onClick={() => setShowForm(!showForm)}>
          {showForm ? 'Cancel' : '+ Add Skill'}
        </button>
      </div>

      {showForm && (
        <div className="form-card">
          <h3>Add New Skill</h3>
          <form onSubmit={handleSubmit}>
            <div className="form-row">
              <label>Name:</label>
              <input 
                type="text" 
                value={newSkill.name}
                onChange={(e) => setNewSkill({...newSkill, name: e.target.value})}
                required 
              />
            </div>
            <div className="form-row">
              <label>Description:</label>
              <textarea 
                value={newSkill.description}
                onChange={(e) => setNewSkill({...newSkill, description: e.target.value})}
              />
            </div>
            <div className="form-row">
              <label>Category:</label>
              <select 
                value={newSkill.category}
                onChange={(e) => setNewSkill({...newSkill, category: e.target.value})}
              >
                <option value="Combat">Combat</option>
                <option value="Magic">Magic</option>
                <option value="Stealth">Stealth</option>
                <option value="Crafting">Crafting</option>
                <option value="Social">Social</option>
              </select>
            </div>
            <div className="form-row">
              <label>Max Level:</label>
              <input 
                type="number" 
                min="1" 
                max="100"
                value={newSkill.maxLevel}
                onChange={(e) => setNewSkill({...newSkill, maxLevel: parseInt(e.target.value)})}
              />
            </div>
            <button type="submit">Create Skill</button>
          </form>
        </div>
      )}

      {categories.map(cat => (
        <div key={cat} className="skill-category">
          <h3>{cat}</h3>
          <div className="skills-grid">
            {skills.filter(s => s.category === cat).map((skill) => (
              <div key={skill.id} className="skill-card">
                <h4>{skill.name}</h4>
                <p>{skill.description}</p>
                <div className="skill-meta">
                  <span>Max Level: {skill.maxLevel}</span>
                </div>
                <div className="skill-actions">
                  <button>Edit</button>
                  <button className="danger">Delete</button>
                </div>
              </div>
            ))}
          </div>
        </div>
      ))}
    </div>
  )
}