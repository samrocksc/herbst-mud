import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'
import {
  useWeaponSkills,
  useCreateWeaponSkill,
  useUpdateWeaponSkill,
  useDeleteWeaponSkill,
  type TrainableSkill,
  type TrainableSkillInput,
} from '../../hooks/useWeaponSkills'
import { PageHeader } from '../../components/PageHeader'
import { DataTable, type Column } from '../../components/DataTable'
import { Button } from '../../components/Button'

export const Route = createFileRoute('/_auth/skills')({
  component: TrainableSkillsManagement,
})

const EMPTY_SKILL: TrainableSkillInput = {
  name: '',
  description: '',
  skill_category: '',
  requirements: '',
}

function SkillForm({
  skill,
  onSubmit,
  onCancel,
  isLoading,
}: {
  skill: TrainableSkill | null
  onSubmit: (data: TrainableSkillInput) => void
  onCancel: () => void
  isLoading: boolean
}) {
  const [formData, setFormData] = useState<TrainableSkillInput>(() => {
    if (skill) {
      return {
        id: skill.id,
        name: skill.name,
        description: skill.description,
        skill_category: skill.skill_category,
        requirements: skill.requirements,
      }
    }
    return EMPTY_SKILL
  })

  return (
    <div className="form-card">
      <h3>{skill ? 'Edit Skill' : 'Add New Skill'}</h3>
      <form
        onSubmit={(e) => {
          e.preventDefault()
          onSubmit(formData)
        }}
      >
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
          <label>Category:</label>
          <select
            value={formData.skill_category}
            onChange={(e) => setFormData({ ...formData, skill_category: e.target.value })}
          >
            <option value="">— None —</option>
            <optgroup label="Combat">
              <option value="blades">Blades</option>
              <option value="knives">Knives</option>
              <option value="staves">Staves</option>
              <option value="brawling">Brawling</option>
              <option value="martial">Martial Arts</option>
              <option value="bows">Bows</option>
              <option value="thrown">Thrown Weapons</option>
            </optgroup>
            <optgroup label="Magic">
              <option value="fire_magic">Fire Magic</option>
              <option value="water_magic">Water Magic</option>
              <option value="wind_magic">Wind Magic</option>
              <option value="earth_magic">Earth Magic</option>
              <option value="light_magic">Light Magic</option>
              <option value="dark_magic">Dark Magic</option>
            </optgroup>
            <optgroup label="Defense">
              <option value="light_armor">Light Armor</option>
              <option value="cloth_armor">Cloth Armor</option>
              <option value="heavy_armor">Heavy Armor</option>
              <option value="shields">Shields</option>
            </optgroup>
            <optgroup label="Utility">
              <option value="tech">Tech</option>
              <option value="pizza_making">Pizza Making</option>
              <option value="crafting">Crafting</option>
              <option value="trading">Trading</option>
            </optgroup>
          </select>
        </div>

        <div className="form-row">
          <label>Requirements:</label>
          <input
            type="text"
            value={formData.requirements}
            onChange={(e) => setFormData({ ...formData, requirements: e.target.value })}
            placeholder="e.g. level:5, str:10"
          />
        </div>

        <div className="form-actions">
          <Button type="submit" variant="primary" disabled={isLoading}>
            {isLoading ? 'Saving...' : skill ? 'Update Skill' : 'Create Skill'}
          </Button>
          <Button variant="secondary" onClick={onCancel}>
            Cancel
          </Button>
        </div>
      </form>
    </div>
  )
}

function DeleteConfirmation({
  skill,
  onConfirm,
  onCancel,
  isLoading,
}: {
  skill: TrainableSkill
  onConfirm: () => void
  onCancel: () => void
  isLoading: boolean
}) {
  return (
    <div className="modal-overlay" onClick={onCancel}>
      <div className="modal-content modal-sm" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h3>Delete Skill</h3>
          <Button variant="ghost" size="sm" onClick={onCancel} aria-label="Close">
            ×
          </Button>
        </div>
        <div className="modal-body">
          <p>
            Are you sure you want to delete <strong>{skill.name}</strong>?
          </p>
          <p className="text-muted">This action cannot be undone.</p>
        </div>
        <div className="modal-footer">
          <Button variant="danger" onClick={onConfirm} disabled={isLoading}>
            {isLoading ? 'Deleting...' : 'Delete'}
          </Button>
          <Button variant="secondary" onClick={onCancel}>
            Cancel
          </Button>
        </div>
      </div>
    </div>
  )
}

const COLUMNS: Column<TrainableSkill>[] = [
  {
    header: 'Name',
    accessor: 'name',
    render: (_: unknown, row: TrainableSkill) => <strong>{row.name}</strong>,
  },
  {
    header: 'Description',
    accessor: 'description',
  },
  {
    header: 'Category',
    accessor: 'skill_category',
    render: (val: unknown) =>
      val ? <span className="talent-effect">{String(val)}</span> : <span className="text-muted">—</span>,
  },
  {
    header: 'Requirements',
    accessor: 'requirements',
  },
  {
    header: 'Actions',
    accessor: '_actions',
    render: () => null, // value unused — actions rendered below
  },
]

function TrainableSkillsManagement() {
  const [showForm, setShowForm] = useState(false)
  const [editing, setEditing] = useState<TrainableSkill | null>(null)
  const [deleting, setDeleting] = useState<TrainableSkill | null>(null)

  const create = useCreateWeaponSkill()
  const update = useUpdateWeaponSkill()
  const remove = useDeleteWeaponSkill()
  const { data: skills, isLoading, error } = useWeaponSkills()

  const handleSubmit = async (formData: TrainableSkillInput) => {
    if (editing) {
      await update.mutateAsync({ id: editing.id, input: formData })
    } else {
      await create.mutateAsync(formData)
    }
    setShowForm(false)
    setEditing(null)
  }

  if (isLoading) return <div className="loading">Loading skills...</div>
  if (error) return <div className="error">Failed to load skills: {error.message}</div>

  const columns: Column<TrainableSkill>[] = [
    ...COLUMNS.slice(0, 4),
    {
      header: 'Actions',
      accessor: '_actions',
      render: (_: unknown, row: TrainableSkill) => (
        <span className="inline-flex gap-2">
          <Button variant="accent" size="sm" onClick={() => { setEditing(row); setShowForm(true) }}>
            Edit
          </Button>
          <Button variant="danger" size="sm" onClick={() => setDeleting(row)}>
            Delete
          </Button>
        </span>
      ),
    },
  ]

  return (
    <div className="management-page">
      <PageHeader
        title="Skills"
        backTo="/dashboard"
        actions={
          <Button
            variant="primary"
            onClick={() => {
              setEditing(null)
              setShowForm(true)
            }}
          >
            + Add Skill
          </Button>
        }
      />

      {showForm && (
        <SkillForm
          skill={editing}
          onSubmit={handleSubmit}
          onCancel={() => {
            setShowForm(false)
            setEditing(null)
          }}
          isLoading={create.isPending || update.isPending}
        />
      )}

      <DataTable
        columns={columns}
        data={skills ?? []}
        getKey={(row: TrainableSkill) => row.id}
        emptyMessage="No skills found. Add your first trainable skill!"
      />

      {deleting && (
        <DeleteConfirmation
          skill={deleting}
          onConfirm={async () => {
            await remove.mutateAsync(deleting.id)
            setDeleting(null)
          }}
          onCancel={() => setDeleting(null)}
          isLoading={remove.isPending}
        />
      )}
    </div>
  )
}
