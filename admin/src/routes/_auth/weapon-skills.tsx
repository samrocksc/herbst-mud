import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'
import { useWeaponSkills, useCreateWeaponSkill, useUpdateWeaponSkill, useDeleteWeaponSkill, type WeaponSkill, type WeaponSkillInput } from '../../hooks/useWeaponSkills'
import { PageHeader } from '../../components/PageHeader'

export const Route = createFileRoute('/_auth/weapon-skills')({
  component: WeaponSkillsManagement,
})

const EMPTY_WEAPON_SKILL: WeaponSkillInput = {
  name: '',
  description: '',
  skill_category: '',
  effect_type: '',
  effect_value: 0,
  effect_duration: 0,
  cooldown: 0,
  mana_cost: 0,
  stamina_cost: 0,
  requirements: '',
}

function WeaponSkillForm({
  weaponSkill,
  onSubmit,
  onCancel,
  isLoading
}: {
  weaponSkill: WeaponSkill | null
  onSubmit: (data: WeaponSkillInput) => void
  onCancel: () => void
  isLoading: boolean
}) {
  const [formData, setFormData] = useState<WeaponSkillInput>(() => {
    if (weaponSkill) {
      return {
        id: weaponSkill.id,
        name: weaponSkill.name,
        description: weaponSkill.description,
        skill_category: weaponSkill.skill_category,
        effect_type: weaponSkill.effect_type,
        effect_value: weaponSkill.effect_value,
        effect_duration: weaponSkill.effect_duration,
        cooldown: weaponSkill.cooldown,
        mana_cost: weaponSkill.mana_cost,
        stamina_cost: weaponSkill.stamina_cost,
        requirements: weaponSkill.requirements,
      }
    }
    return EMPTY_WEAPON_SKILL
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    onSubmit(formData)
  }

  return (
    <div className="form-card">
      <h3>{weaponSkill ? 'Edit Weapon Skill' : 'Add New Weapon Skill'}</h3>
      <form onSubmit={handleSubmit}>
        <div className="form-row">
          <label>Name:</label>
          <input
            type="text"
            value={formData.name}
            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
            required
          />
        </div>

        <div className="form-row">
          <label>Description:</label>
          <textarea
            value={formData.description}
            onChange={(e) => setFormData({ ...formData, description: e.target.value })}
            rows={3}
          />
        </div>

        <div className="form-row">
          <label>Skill Category:</label>
          <select
            value={formData.skill_category}
            onChange={(e) => setFormData({ ...formData, skill_category: e.target.value })}
          >
            <option value="">— None —</option>
            <option value="blades">Blades</option>
            <option value="staves">Staves</option>
            <option value="knives">Knives</option>
            <option value="martial">Martial Arts</option>
            <option value="brawling">Brawling</option>
            <option value="tech">Tech</option>
            <option value="fire_magic">Fire Magic</option>
            <option value="water_magic">Water Magic</option>
            <option value="wind_magic">Wind Magic</option>
            <option value="earth_magic">Earth Magic</option>
            <option value="light_magic">Light Magic</option>
            <option value="dark_magic">Dark Magic</option>
            <option value="light_armor">Light Armor</option>
            <option value="cloth_armor">Cloth Armor</option>
            <option value="heavy_armor">Heavy Armor</option>
            <option value="shields">Shields</option>
            <option value="bows">Bows</option>
            <option value="thrown">Thrown Weapons</option>
          </select>
        </div>

        <div className="form-row">
          <label>Requirements (e.g. level:5, str:10):</label>
          <input
            type="text"
            value={formData.requirements}
            onChange={(e) => setFormData({ ...formData, requirements: e.target.value })}
            placeholder="e.g., level:5, str:10"
          />
        </div>

        <div className="form-row">
          <label>Effect Type:</label>
          <select
            value={formData.effect_type}
            onChange={(e) => setFormData({ ...formData, effect_type: e.target.value })}
          >
            <option value="">— None —</option>
            <option value="heal">Heal</option>
            <option value="damage">Damage</option>
            <option value="dot">Damage over Time</option>
            <option value="buff_armor">Buff: Armor</option>
            <option value="buff_dodge">Buff: Dodge</option>
            <option value="buff_crit">Buff: Critical</option>
            <option value="debuff">Debuff</option>
          </select>
        </div>

        <div className="form-row-group">
          <div className="form-row">
            <label>Effect Value:</label>
            <input
              type="number"
              min="0"
              value={formData.effect_value}
              onChange={(e) => setFormData({ ...formData, effect_value: parseInt(e.target.value) || 0 })}
            />
          </div>
          <div className="form-row">
            <label>Effect Duration (ticks):</label>
            <input
              type="number"
              min="0"
              value={formData.effect_duration}
              onChange={(e) => setFormData({ ...formData, effect_duration: parseInt(e.target.value) || 0 })}
            />
          </div>
        </div>

        <div className="form-row-group">
          <div className="form-row">
            <label>Cooldown (ticks):</label>
            <input
              type="number"
              min="0"
              value={formData.cooldown}
              onChange={(e) => setFormData({ ...formData, cooldown: parseInt(e.target.value) || 0 })}
            />
          </div>
          <div className="form-row">
            <label>Mana Cost:</label>
            <input
              type="number"
              min="0"
              value={formData.mana_cost}
              onChange={(e) => setFormData({ ...formData, mana_cost: parseInt(e.target.value) || 0 })}
            />
          </div>
          <div className="form-row">
            <label>Stamina Cost:</label>
            <input
              type="number"
              min="0"
              value={formData.stamina_cost}
              onChange={(e) => setFormData({ ...formData, stamina_cost: parseInt(e.target.value) || 0 })}
            />
          </div>
        </div>

        <div className="form-actions">
          <button type="submit" disabled={isLoading}>
            {isLoading ? 'Saving...' : weaponSkill ? 'Update Weapon Skill' : 'Create Weapon Skill'}
          </button>
          <button type="button" className="btn-cancel" onClick={onCancel}>
            Cancel
          </button>
        </div>
      </form>
    </div>
  )
}

function DeleteConfirmation({
  weaponSkill,
  onConfirm,
  onCancel,
  isLoading
}: {
  weaponSkill: WeaponSkill
  onConfirm: () => void
  onCancel: () => void
  isLoading: boolean
}) {
  return (
    <div className="modal-overlay" onClick={onCancel}>
      <div className="modal-content modal-sm" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h3>Delete Weapon Skill</h3>
          <button className="modal-close" onClick={onCancel}>×</button>
        </div>
        <div className="modal-body">
          <p>Are you sure you want to delete <strong>{weaponSkill.name}</strong>?</p>
          <p className="text-muted">This action cannot be undone.</p>
        </div>
        <div className="modal-footer">
          <button className="btn-danger" onClick={onConfirm} disabled={isLoading}>
            {isLoading ? 'Deleting...' : 'Delete'}
          </button>
          <button className="btn-cancel" onClick={onCancel}>Cancel</button>
        </div>
      </div>
    </div>
  )
}

function WeaponSkillsManagement() {
  const [showForm, setShowForm] = useState(false)
  const [editingWeaponSkill, setEditingWeaponSkill] = useState<WeaponSkill | null>(null)
  const [deletingWeaponSkill, setDeletingWeaponSkill] = useState<WeaponSkill | null>(null)

  const createWeaponSkill = useCreateWeaponSkill()
  const updateWeaponSkill = useUpdateWeaponSkill()
  const deleteWeaponSkill = useDeleteWeaponSkill()

  const { data: weaponSkills, isLoading, error } = useWeaponSkills()

  const handleSubmit = async (formData: WeaponSkillInput) => {
    if (editingWeaponSkill) {
      await updateWeaponSkill.mutateAsync({ id: editingWeaponSkill.id, input: formData })
    } else {
      await createWeaponSkill.mutateAsync(formData)
    }
    setShowForm(false)
    setEditingWeaponSkill(null)
  }

  const handleEdit = (weaponSkill: WeaponSkill) => {
    setEditingWeaponSkill(weaponSkill)
    setShowForm(true)
  }

  const handleDelete = async () => {
    if (deletingWeaponSkill) {
      await deleteWeaponSkill.mutateAsync(deletingWeaponSkill.id)
      setDeletingWeaponSkill(null)
    }
  }

  const handleCancelForm = () => {
    setShowForm(false)
    setEditingWeaponSkill(null)
  }

  if (isLoading) return <div className="loading">Loading weapon skills...</div>
  if (error) return <div className="error">Failed to load weapon skills: {error.message}</div>

  return (
    <div className="management-page">
      <PageHeader
        title="Weapon Skills"
        backTo="/dashboard"
        actions={<button className="btn-primary" onClick={() => { setEditingWeaponSkill(null); setShowForm(true) }}>+ Add Weapon Skill</button>}
      />

      {showForm && (
        <WeaponSkillForm
          weaponSkill={editingWeaponSkill}
          onSubmit={handleSubmit}
          onCancel={handleCancelForm}
          isLoading={createWeaponSkill.isPending || updateWeaponSkill.isPending}
        />
      )}

      <div className="table-container">
        <table className="table">
          <thead>
            <tr>
              <th>Name</th>
              <th>Description</th>
              <th>Category</th>
              <th>Effect</th>
              <th>Requirements</th>
              <th>Costs</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {weaponSkills?.map((ws) => (
              <tr key={ws.id}>
                <td><strong>{ws.name}</strong></td>
                <td className="text-muted">{ws.description || '-'}</td>
                <td>
                  <span className={`talent-effect talent-effect-${ws.skill_category}`}>
                    {ws.skill_category}
                  </span>
                </td>
                <td>
                  {ws.effect_type && (
                    <span className="talent-effect">{ws.effect_type}</span>
                  )}
                  {ws.effect_value > 0 && (
                    <span className="talent-effect-value"> {ws.effect_value}</span>
                  )}
                  {ws.effect_duration > 0 && (
                    <span className="text-muted"> ({ws.effect_duration}t)</span>
                  )}
                </td>
                <td className="text-muted">{ws.requirements || '-'}</td>
                <td>
                  {ws.mana_cost > 0 && <span className="cost-badge" title="Mana Cost">MP: {ws.mana_cost}</span>}
                  {ws.stamina_cost > 0 && <span className="cost-badge" title="Stamina Cost">SP: {ws.stamina_cost}</span>}
                  {ws.cooldown > 0 && <span className="cost-badge" title="Cooldown">CD: {ws.cooldown}</span>}
                </td>
                <td>
                  <button onClick={() => handleEdit(ws)}>Edit</button>
                  <button className="danger" onClick={() => setDeletingWeaponSkill(ws)}>Delete</button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>

        {(!weaponSkills || weaponSkills.length === 0) && (
          <div className="empty-state">
            <p>No weapon skills found. Create your first weapon skill!</p>
          </div>
        )}
      </div>

      {deletingWeaponSkill && (
        <DeleteConfirmation
          weaponSkill={deletingWeaponSkill}
          onConfirm={handleDelete}
          onCancel={() => setDeletingWeaponSkill(null)}
          isLoading={deleteWeaponSkill.isPending}
        />
      )}
    </div>
  )
}
