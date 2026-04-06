import { createFileRoute } from '@tanstack/react-router'
import { useState, useMemo } from 'react'
import { useSkills, useCreateSkill, useUpdateSkill, useDeleteSkill, type Skill, type SkillInput } from '../../hooks/useSkills'

export const Route = createFileRoute('/_auth/skills')({
  component: SkillsManagement,
})

const EMPTY_SKILL: SkillInput = {
  name: '',
  description: '',
  type: 'weapon',
  tags: '',
  primary_stat: 'STR',
  level_req: 1,
  cooldown: 0,
  mana_cost: 0,
  stamina_cost: 0,
  classless: false,
  effects: ''
}

function SkillForm({
  skill,
  onSubmit,
  onCancel,
  isLoading
}: {
  skill: Skill | null
  onSubmit: (data: SkillInput) => void
  onCancel: () => void
  isLoading: boolean
}) {
  const [formData, setFormData] = useState<SkillInput>(() => {
    if (skill) {
      return {
        ...skill,
        tags: skill.tags.join(', '),
        effects: skill.effects ? JSON.stringify(skill.effects, null, 2) : ''
      }
    }
    return EMPTY_SKILL
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    onSubmit(formData)
  }

  return (
    <div className="form-card">
      <h3>{skill ? 'Edit Skill' : 'Add New Skill'}</h3>
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
          <label>Type:</label>
          <select
            value={formData.type}
            onChange={(e) => setFormData({ ...formData, type: e.target.value as SkillInput['type'] })}
          >
            <option value="weapon">Weapon</option>
            <option value="magic">Magic</option>
            <option value="utility">Utility</option>
          </select>
        </div>

        <div className="form-row">
          <label>Tags (comma-separated):</label>
          <input
            type="text"
            value={formData.tags}
            onChange={(e) => setFormData({ ...formData, tags: e.target.value })}
            placeholder="e.g., classless, melee, healing"
          />
        </div>

        <div className="form-row">
          <label>Primary Stat:</label>
          <select
            value={formData.primary_stat}
            onChange={(e) => setFormData({ ...formData, primary_stat: e.target.value as SkillInput['primary_stat'] })}
          >
            <option value="STR">Strength (STR)</option>
            <option value="DEX">Dexterity (DEX)</option>
            <option value="INT">Intelligence (INT)</option>
            <option value="WIS">Wisdom (WIS)</option>
          </select>
        </div>

        <div className="form-row-group">
          <div className="form-row">
            <label>Level Req:</label>
            <input
              type="number"
              min="1"
              max="100"
              value={formData.level_req}
              onChange={(e) => setFormData({ ...formData, level_req: parseInt(e.target.value) || 1 })}
            />
          </div>
          <div className="form-row">
            <label>Cooldown (s):</label>
            <input
              type="number"
              min="0"
              value={formData.cooldown}
              onChange={(e) => setFormData({ ...formData, cooldown: parseInt(e.target.value) || 0 })}
            />
          </div>
        </div>

        <div className="form-row-group">
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

        <div className="form-row checkbox">
          <label>
            <input
              type="checkbox"
              checked={formData.classless}
              onChange={(e) => setFormData({ ...formData, classless: e.target.checked })}
            />
            Classless (available to all characters)
          </label>
        </div>

        <div className="form-row">
          <label>Effects (JSON, optional):</label>
          <textarea
            value={formData.effects || ''}
            onChange={(e) => setFormData({ ...formData, effects: e.target.value })}
            placeholder='{"damage": 10, "heal": 5}'
            rows={3}
          />
        </div>

        <div className="form-actions">
          <button type="submit" disabled={isLoading}>
            {isLoading ? 'Saving...' : skill ? 'Update Skill' : 'Create Skill'}
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
  skill,
  onConfirm,
  onCancel,
  isLoading
}: {
  skill: Skill
  onConfirm: () => void
  onCancel: () => void
  isLoading: boolean
}) {
  return (
    <div className="modal-overlay" onClick={onCancel}>
      <div className="modal-content modal-sm" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h3>Delete Skill</h3>
          <button className="modal-close" onClick={onCancel}>×</button>
        </div>
        <div className="modal-body">
          <p>Are you sure you want to delete <strong>{skill.name}</strong>?</p>
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

function SkillsManagement() {
  const [filterType, setFilterType] = useState<string>('')
  const [filterTag, setFilterTag] = useState<string>('')
  const [showForm, setShowForm] = useState(false)
  const [editingSkill, setEditingSkill] = useState<Skill | null>(null)
  const [deletingSkill, setDeletingSkill] = useState<Skill | null>(null)

  const createSkill = useCreateSkill()
  const updateSkill = useUpdateSkill()
  const deleteSkill = useDeleteSkill()

  const { data: skills, isLoading, error } = useSkills({
    type: filterType || undefined,
    tag: filterTag || undefined
  })

  const allTags = useMemo(() => {
    if (!skills) return []
    const skillArray = Array.isArray(skills) ? skills : (skills as { skills: Skill[] }).skills ?? []
    const tags = new Set<string>()
    skillArray.forEach((s: Skill) => s.tags?.forEach((t: string) => tags.add(t)))
    return Array.from(tags).sort()
  }, [skills])

  const handleSubmit = async (formData: SkillInput) => {
    if (editingSkill) {
      await updateSkill.mutateAsync({ id: editingSkill.id, input: formData })
    } else {
      await createSkill.mutateAsync(formData)
    }
    setShowForm(false)
    setEditingSkill(null)
  }

  const handleEdit = (skill: Skill) => {
    setEditingSkill(skill)
    setShowForm(true)
  }

  const handleDelete = async () => {
    if (deletingSkill) {
      await deleteSkill.mutateAsync(deletingSkill.id)
      setDeletingSkill(null)
    }
  }

  const handleCancelForm = () => {
    setShowForm(false)
    setEditingSkill(null)
  }

  if (isLoading) return <div className="loading">Loading skills...</div>
  if (error) return <div className="error">Failed to load skills: {error.message}</div>

  return (
    <div className="management-page">
      <div className="page-header">
        <h2>Skills & Talents Management</h2>
        <button onClick={() => { setEditingSkill(null); setShowForm(true) }}>
          + Add Skill
        </button>
      </div>

      <div className="filters-bar">
        <div className="filter-group">
          <label>Type:</label>
          <select value={filterType} onChange={(e) => setFilterType(e.target.value)}>
            <option value="">All Types</option>
            <option value="weapon">Weapon</option>
            <option value="magic">Magic</option>
            <option value="utility">Utility</option>
          </select>
        </div>
        <div className="filter-group">
          <label>Tag:</label>
          <select value={filterTag} onChange={(e) => setFilterTag(e.target.value)}>
            <option value="">All Tags</option>
            {allTags.map(tag => (
              <option key={tag} value={tag}>{tag}</option>
            ))}
          </select>
        </div>
        {(filterType || filterTag) && (
          <button className="btn-clear-filters" onClick={() => { setFilterType(''); setFilterTag('') }}>
            Clear Filters
          </button>
        )}
      </div>

      {showForm && (
        <SkillForm
          skill={editingSkill}
          onSubmit={handleSubmit}
          onCancel={handleCancelForm}
          isLoading={createSkill.isPending || updateSkill.isPending}
        />
      )}

      <div className="skills-grid">
        {skills?.map((skill) => (
          <div key={skill.id} className="skill-card">
            <div className="skill-header">
              <h4>{skill.name}</h4>
              {skill.classless && (
                <span className="badge badge-classless">Classless</span>
              )}
            </div>
            <p className="skill-desc">{skill.description}</p>
            <div className="skill-meta">
              <span className={`skill-type skill-type-${skill.type}`}>{skill.type}</span>
              <span className="skill-stat">{skill.primary_stat}</span>
            </div>
            <div className="skill-tags">
              {(skill.tags ?? []).map(tag => (
                <span key={tag} className="tag">{tag}</span>
              ))}
            </div>
            <div className="skill-costs">
              <span title="Mana Cost">MP: {skill.mana_cost}</span>
              <span title="Stamina Cost">SP: {skill.stamina_cost}</span>
              <span title="Cooldown">CD: {skill.cooldown}s</span>
            </div>
            <div className="skill-level-req">Level {skill.level_req}+</div>
            <div className="skill-actions">
              <button onClick={() => handleEdit(skill)}>Edit</button>
              <button className="danger" onClick={() => setDeletingSkill(skill)}>Delete</button>
            </div>
          </div>
        ))}
      </div>

      {skills?.length === 0 && (
        <div className="empty-state">
          <p>No skills found. {filterType || filterTag ? 'Try clearing filters.' : 'Create your first skill!'}</p>
        </div>
      )}

      {deletingSkill && (
        <DeleteConfirmation
          skill={deletingSkill}
          onConfirm={handleDelete}
          onCancel={() => setDeletingSkill(null)}
          isLoading={deleteSkill.isPending}
        />
      )}
    </div>
  )
}
