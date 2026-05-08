import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'
import { useCompetencyCategories, type CompetencyCategory, type CompetencyThreshold } from '../../hooks/useCompetencies'
import { PageHeader } from '../../components/PageHeader'
import { DataTable, type Column } from '../../components/DataTable'
import { Button } from '../../components/Button'

export const Route = createFileRoute('/_auth/skills')({
  component: SkillsManagement,
})

function SkillsManagement() {
  const { data: categories, isLoading, error } = useCompetencyCategories()
  const [selectedId, setSelectedId] = useState<string | null>(null)

  if (isLoading) return <div className="loading">Loading skills...</div>
  if (error) return <div className="error">Failed to load skills: {error.message}</div>

  const selected = categories?.find((c) => c.id === selectedId)

  const columns: Column<CompetencyCategory>[] = [
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
      render: (val: unknown) => `${(val as CompetencyThreshold[]).length}`,
    },
    {
      header: '',
      accessor: '_action',
      render: (_: unknown, row: CompetencyCategory) => (
        <Button variant="ghost" size="sm" onClick={() => setSelectedId(row.id)}>
          View Thresholds
        </Button>
      ),
    },
  ]

  return (
    <div className="management-page">
      <PageHeader title="Skills (Competencies)" backTo="/dashboard" />

      <p className="text-sm text-muted mb-4">
        Skills are leveled proficiencies (Blades, Staves, etc.) that gate equipment use and provide
        damage/defense bonuses as they level up. Each category has 10 levels with increasing XP requirement
        and better multipliers.
      </p>

      <DataTable
        columns={columns}
        data={categories ?? []}
        getKey={(row) => row.id}
        emptyMessage="No competency categories found. Run seed data first."
      />

      {selected && (
        <ThresholdTable category={selected} onClose={() => setSelectedId(null)} />
      )}
    </div>
  )
}

function ThresholdTable({
  category,
  onClose,
}: Readonly<{ category: CompetencyCategory; onClose: () => void }>) {
  const thresholds = [...(category.thresholds ?? [])].sort((a, b) => a.level - b.level)

  const columns: Column<CompetencyThreshold>[] = [
    { header: 'Level', accessor: 'level' },
    {
      header: 'XP Required',
      accessor: 'xp_required',
      render: (val: unknown) => Number(val).toLocaleString(),
    },
    {
      header: 'Dmg Mult',
      accessor: 'damage_multiplier',
      render: (val: unknown) => `${Number(val).toFixed(2)}x`,
    },
    {
      header: 'Def Mult',
      accessor: 'defense_multiplier',
      render: (val: unknown) => `${Number(val).toFixed(2)}x`,
    },
  ]

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content modal-lg" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h3>{category.name} — Level Thresholds</h3>
          <Button variant="ghost" size="sm" onClick={onClose}>×</Button>
        </div>
        <div className="modal-body">
          <p className="text-sm text-muted mb-3">
            XP multiplier: {category.xp_multiplier * 100}% — Each level provides better
            damage and defense multipliers for weapons/armor in this category.
          </p>
          <DataTable
            columns={columns}
            data={thresholds}
            getKey={(row) => row.level}
            emptyMessage="No thresholds defined"
          />
        </div>
      </div>
    </div>
  )
}