import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'
import {
  useCompetencyCategories,
  useDeleteCompetencyCategory,
  type CompetencyCategory,
} from '../../hooks/useCompetencies'
import { PageHeader } from '../../components/PageHeader'
import { DataTable, type Column } from '../../components/DataTable'
import { Button } from '../../components/Button'
import { SkillForm } from './SkillForm'

export const Route = createFileRoute('/_auth/skills')({
  component: SkillsManagement,
})

const COLUMNS: Column<CompetencyCategory>[] = [
  { header: 'ID', accessor: 'id' },
  { header: 'Name', accessor: 'name' },
  {
    header: 'XP Mult',
    accessor: 'xp_multiplier',
    render: (val: unknown) => `${Number(val) * 100}%`,
  },
  {
    header: 'Levels',
    accessor: 'thresholds',
    render: (val: unknown) => `${(val as unknown[]).length}`,
  },
]

function DeleteConfirmation({
  category,
  onConfirm,
  onCancel,
  isLoading,
}: Readonly<{
  category: CompetencyCategory
  onConfirm: () => void
  onCancel: () => void
  isLoading: boolean
}>) {
  return (
    <div className="modal-overlay" onClick={onCancel}>
      <div className="modal-content modal-sm" onClick={e => e.stopPropagation()}>
        <div className="modal-header">
          <h3>Delete Skill</h3>
          <Button variant="ghost" size="sm" onClick={onCancel} aria-label="Close">×</Button>
        </div>
        <div className="modal-body">
          <p>Are you sure you want to delete <strong>{category.name}</strong>?</p>
          <p className="text-muted">This action cannot be undone.</p>
        </div>
        <div className="modal-footer">
          <Button variant="danger" onClick={onConfirm} disabled={isLoading}>
            {isLoading ? 'Deleting...' : 'Delete'}
          </Button>
          <Button variant="secondary" onClick={onCancel}>Cancel</Button>
        </div>
      </div>
    </div>
  )
}

function SkillsManagement() {
  const [showForm, setShowForm] = useState(false)
  const [editingCategory, setEditingCategory] = useState<CompetencyCategory | null>(null)
  const [deletingCategory, setDeletingCategory] = useState<CompetencyCategory | null>(null)
  const deleteMutation = useDeleteCompetencyCategory()

  const { data: categories, isLoading, error } = useCompetencyCategories()

  const handleEdit = (cat: CompetencyCategory) => {
    setEditingCategory(cat)
    setShowForm(true)
  }

  const handleDelete = async () => {
    if (!deletingCategory) return
    try {
      await deleteMutation.mutateAsync(deletingCategory.id)
      setDeletingCategory(null)
    } catch { /* error is in mutation state */ }
  }

  const columns: Column<CompetencyCategory>[] = [
    ...COLUMNS,
    {
      header: 'Actions',
      accessor: '_actions',
      render: (_: unknown, row: CompetencyCategory) => (
        <>
          <Button variant="accent" size="sm" onClick={() => handleEdit(row)}>Edit</Button>
          <Button variant="danger" size="sm" className="ml-2" onClick={() => setDeletingCategory(row)}>Delete</Button>
        </>
      ),
    },
  ]

  if (isLoading) return <div className="loading">Loading skills...</div>
  if (error) return <div className="error">Failed to load skills: {error.message}</div>

  return (
    <div className="management-page">
      <PageHeader
        title="Skills (Competencies)"
        backTo="/dashboard"
        actions={
          <Button variant="primary" onClick={() => { setEditingCategory(null); setShowForm(true) }}>
            + Add Skill
          </Button>
        }
      />
      <p className="text-sm text-muted mb-4">
        Skills are leveled proficiencies (Blades, Staves, etc.) that gate equipment use and provide
        damage/defense bonuses as they level up. Each category has 10 levels with increasing XP requirement
        and better multipliers.
      </p>

      {showForm && (
        <SkillForm
          category={editingCategory}
          onSubmit={() => { setShowForm(false); setEditingCategory(null) }}
          onCancel={() => { setShowForm(false); setEditingCategory(null) }}
        />
      )}

      <DataTable
        columns={columns}
        data={categories ?? []}
        getKey={(row) => row.id}
        emptyMessage="No competency categories found. Create your first skill!"
      />

      {deletingCategory && (
        <DeleteConfirmation
          category={deletingCategory}
          onConfirm={handleDelete}
          onCancel={() => setDeletingCategory(null)}
          isLoading={deleteMutation.isPending}
        />
      )}
    </div>
  )
}