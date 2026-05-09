import { createFileRoute, useNavigate, Link } from '@tanstack/react-router'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useState, useEffect } from 'react'
import { apiGet, apiPost, apiPut, apiDelete } from '../../utils/apiFetch'
import { PageHeader } from '../../components/PageHeader'
import { DataTable, type Column } from '../../components/DataTable'
import { DeleteConfirmation } from '../../components/DeleteConfirmation'
import { Modal } from '../../components/Modal'
import { Button } from '../../components/Button'
import { HooksPanel } from '../../components/HooksPanel'
import {
  FormField,
  NumberField,
  TextareaField,
  FormError,
} from '../../components/FormFields'

export const Route = createFileRoute('/_auth/npcs/$npcId/')({
  component: NpcTemplateDetail,
})

// ─── Types ──────────────────────────────────────────────────────────────────

type NPCTemplate = Readonly<{
  id: string
  name: string
  description: string
  race: string
  disposition: string
  level: number
  xp_value: number
  skills: Record<string, number>
  trades_with: string[]
  greeting: string
  respawn_rooms: string[]
  respawn_cooldown: number
}>

type NPCInstance = Readonly<{
  id: number
  name: string
  npc_template_id: string
  instance_number: number
  room_id: number
  level: number
  race: string
  hitpoints: number
  max_hitpoints: number
  stamina: number
  max_stamina: number
  mana: number
  max_mana: number
  isNPC: boolean
  is_instance: boolean
}>

type SpawnForm = Readonly<{
  room_id: number
  instance_number: number
  instance_name: string
}>

type EditForm = Readonly<{
  name: string
  description: string
  race: string
  disposition: string
  level: number
  xp_value: number
  skills: string
  trades_with: string
  greeting: string
  respawn_rooms: string
  respawn_cooldown: number
}>

const API = `${window.location.origin}`

// ─── Instance table columns factory ─────────────────────────────────────────

function buildInstanceColumns(npcId: string): Column<NPCInstance>[] {
  return [
    { header: 'ID', accessor: 'id', align: 'center' },
    {
      header: 'Name',
      accessor: 'name',
      render: (_: unknown, row: NPCInstance) => (
        <Link
          to="/npcs/$npcId/instances/$instanceId"
          params={{ npcId, instanceId: String(row.id) }}
          className="no-underline text-primary hover:underline font-bold"
        >
          {row.name}
        </Link>
      ),
    },
    {
      header: 'Location',
      accessor: 'room_id',
      render: (_: unknown, row: NPCInstance) => (
        <Link
          to="/map"
          className="text-primary no-underline hover:underline text-xs"
        >
          Room #{row.room_id}
        </Link>
      ),
    },
    { header: 'Instance #', accessor: 'instance_number', align: 'center' },
    {
      header: 'HP',
      accessor: 'hitpoints',
      render: (_: unknown, row: NPCInstance) => (
        <span>{row.hitpoints}/{row.max_hitpoints}</span>
      ),
    },
  ]
}

// ─── Helpers to build edit form from template ───────────────────────────────

function templateToEditForm(t: NPCTemplate): EditForm {
  const skills = t.skills ?? {}
  const trades = t.trades_with ?? []
  const rooms = t.respawn_rooms ?? []
  return {
    name: t.name,
    description: t.description,
    race: t.race,
    disposition: t.disposition,
    level: t.level,
    xp_value: t.xp_value,
    skills: Object.entries(skills)
      .map(([k, v]) => `${k}:${v}`)
      .join('\\n'),
    trades_with: trades.join('\\n'),
    greeting: t.greeting ?? '',
    respawn_rooms: rooms.join('\\n'),
    respawn_cooldown: t.respawn_cooldown,
  }
}

function editFormToPayload(form: EditForm) {
  const skills: Record<string, number> = {}
  for (const line of form.skills.split('\\n')) {
    const trimmed = line.trim()
    if (!trimmed) continue
    const colonIdx = trimmed.lastIndexOf(':')
    if (colonIdx > 0) {
      skills[trimmed.slice(0, colonIdx)] = parseInt(trimmed.slice(colonIdx + 1)) || 0
    } else {
      skills[trimmed] = 0
    }
  }
  return {
    name: form.name,
    description: form.description,
    race: form.race,
    disposition: form.disposition,
    level: form.level,
    xp_value: form.xp_value,
    skills,
    trades_with: form.trades_with.split('\\n').map((s) => s.trim()).filter(Boolean),
    greeting: form.greeting,
    respawn_rooms: form.respawn_rooms.split('\\n').map((s) => s.trim()).filter(Boolean),
    respawn_cooldown: form.respawn_cooldown,
  }
}

// ─── Component ──────────────────────────────────────────────────────────────

export function NpcTemplateDetail() {
  const { npcId } = Route.useParams()
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const [editing, setEditing] = useState(false)
  const [showDeleteModal, setShowDeleteModal] = useState(false)
  const [showSpawnModal, setShowSpawnModal] = useState(false)
  const [form, setForm] = useState<EditForm | null>(null)
  const [spawnForm, setSpawnForm] = useState<SpawnForm>({
    room_id: 1,
    instance_number: 1,
    instance_name: '',
  })

  const templateQuery = useQuery<NPCTemplate>({
    queryKey: ['npc-templates', npcId],
    queryFn: () => apiGet<NPCTemplate>(`${API}/api/npc-templates/${npcId}`),
  })

  const instancesQuery = useQuery<NPCInstance[]>({
    queryKey: ['npc-instances'],
    queryFn: () => apiGet<NPCInstance[]>(`${API}/api/npc-instances`),
  })

  const templateInstances = (instancesQuery.data ?? []).filter((inst) => {
    return inst.npc_template_id === npcId
  })

  useEffect(() => {
    if (editing && templateQuery.data && !form) {
      setForm(templateToEditForm(templateQuery.data))
    }
  }, [editing, templateQuery.data, form])

  const deleteMutation = useMutation({
    mutationFn: () => apiDelete(`${API}/api/npc-templates/${npcId}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['npc-templates'] })
      navigate({ to: '/npcs' })
    },
  })

  const updateMutation = useMutation({
    mutationFn: (body: ReturnType<typeof editFormToPayload>) =>
      apiPut(`${API}/api/npc-templates/${npcId}`, body),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['npc-templates', npcId] })
      queryClient.invalidateQueries({ queryKey: ['npc-templates'] })
      setEditing(false)
      setForm(null)
    },
  })

  const spawnMutation = useMutation({
    mutationFn: (input: SpawnForm) =>
      apiPost<NPCInstance>(`${API}/api/npc-instances`, {
        template_id: npcId,
        room_id: input.room_id,
        instance_number: input.instance_number,
        instance_name: input.instance_name || undefined,
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['npc-instances'] })
      setShowSpawnModal(false)
      setSpawnForm({ room_id: 1, instance_number: 1, instance_name: '' })
    },
  })

  const handleDelete = () => {
    deleteMutation.mutate()
  }

  const handleSpawn = () => {
    spawnMutation.mutate(spawnForm)
  }

  const handleSpawnModalClose = () => {
    setShowSpawnModal(false)
    setSpawnForm({ room_id: 1, instance_number: 1, instance_name: '' })
  }

  const startEditing = () => {
    if (templateQuery.data) {
      setForm(templateToEditForm(templateQuery.data))
      setEditing(true)
    }
  }

  const cancelEditing = () => {
    setEditing(false)
    setForm(null)
  }

  const handleSave = () => {
    if (!form) return
    updateMutation.mutate(editFormToPayload(form))
  }

  if (templateQuery.isLoading) {
    return (
      <div className="p-8">
        <PageHeader title="Loading..." backTo="/npcs" />
        <div className="text-text-muted">Loading NPC template...</div>
      </div>
    )
  }

  if (templateQuery.error || !templateQuery.data) {
    return (
      <div className="p-8">
        <PageHeader title="Error" backTo="/npcs" />
        <div className="text-danger">
          Failed to load template: {templateQuery.error?.message ?? 'Unknown error'}
        </div>
      </div>
    )
  }

  const template = templateQuery.data

  return (
    <div className="p-8">
      <PageHeader
        title={template.name}
        backTo="/npcs"
        actions={
          <div className="flex items-center gap-2">
            {editing ? (
              <>
                <Button
                  variant="primary"
                  size="sm"
                  onClick={handleSave}
                  disabled={updateMutation.isPending}
                >
                  {updateMutation.isPending ? 'Saving...' : 'Save'}
                </Button>
                <Button variant="secondary" size="sm" onClick={cancelEditing}>
                  Cancel
                </Button>
              </>
            ) : (
              <>
                <Button variant="primary" size="sm" onClick={startEditing}>
                  Edit
                </Button>
                <Button
                  variant="danger"
                  size="sm"
                  onClick={() => setShowDeleteModal(true)}
                  disabled={deleteMutation.isPending}
                >
                  Delete
                </Button>
              </>
            )}
          </div>
        }
      />

      <div className="max-w-2xl">
        <div className="bg-surface-muted rounded-lg p-6 border border-border mb-6">
          {editing && form ? (
            <TemplateEditForm
              form={form}
              onChange={setForm}
              saveError={updateMutation.isError ? (updateMutation.error as Error)?.message : null}
            />
          ) : (
            <>
              <h2 className="mt-0 mb-4 text-text text-lg font-semibold">Template Stats</h2>
              <div className="grid grid-cols-2 gap-x-6 gap-y-3">
                <DetailField label="ID" value={template.id} />
                <DetailField label="Name" value={template.name} />
                <DetailField label="Race" value={template.race} />
                <DetailField label="Disposition" value={template.disposition} />
                <DetailField label="Level" value={String(template.level)} />
                <DetailField label="XP Value" value={String(template.xp_value)} />
                <DetailField label="Respawn Cooldown" value={`${template.respawn_cooldown}s`} />
                <DetailField label="Respawn Rooms" value={(template.respawn_rooms ?? []).join(', ') || '—'} />
                <DetailField label="Greeting" value={template.greeting || '—'} />
                <DetailField label="Trades With" value={(template.trades_with ?? []).join(', ') || '—'} />
              </div>

              {template.description && (
                <div className="mt-4 pt-4 border-t border-border">
                  <span className="text-text-muted text-xs block mb-1">Description</span>
                  <span className="text-text text-sm">{template.description}</span>
                </div>
              )}

              {Object.keys(template.skills).length > 0 && (
                <div className="mt-4 pt-4 border-t border-border">
                  <span className="text-text-muted text-xs block mb-2">Skills</span>
                  <div className="grid grid-cols-3 gap-x-4 gap-y-1">
                    {Object.entries(template.skills).map(([skill, value]) => (
                      <div key={skill} className="flex justify-between text-sm">
                        <span className="text-text-muted">{skill}</span>
                        <span className="text-text font-medium">{value}</span>
                      </div>
                    ))}
                  </div>
                </div>
              )}
            </>
          )}
        </div>
      </div>

      <div className="mt-6">
        <div className="flex items-center justify-between mb-4">
          <h2 className="m-0 text-text text-lg font-semibold">NPC Instances</h2>
          <Button variant="primary" size="sm" onClick={() => setShowSpawnModal(true)} disabled={editing}>
            + Add Instance
          </Button>
        </div>

        {instancesQuery.isLoading && (
          <div className="p-8 text-text-muted text-center text-xs">Loading instances...</div>
        )}

        {instancesQuery.isError && (
          <div className="p-4 bg-danger/10 border border-danger rounded text-danger text-xs">
            Failed to load instances: {instancesQuery.error?.message ?? 'Unknown error'}
          </div>
        )}

        {!instancesQuery.isLoading && !instancesQuery.isError && (
          <DataTable<NPCInstance>
            columns={buildInstanceColumns(npcId)}
            data={templateInstances}
            getKey={(row) => row.id}
            emptyMessage="No instances of this template."
            variant="dark"
          />
        )}
      </div>

      <HooksPanel npcTemplateId={npcId} />

      <DeleteConfirmation
        open={showDeleteModal}
        title="Delete NPC Template"
        message="Are you sure? This will permanently delete this NPC template."
        onConfirm={handleDelete}
        onCancel={() => setShowDeleteModal(false)}
        isLoading={deleteMutation.isPending}
      />

      <Modal isOpen={showSpawnModal} onClose={handleSpawnModalClose} title="Add NPC Instance">
        <div className="flex flex-col gap-4">
          <div>
            <label className="text-text-muted text-xs block mb-1">Room ID *</label>
            <input
              type="text"
              inputMode="numeric"
              value={spawnForm.room_id || ''}
              onChange={(e) =>
                setSpawnForm({ ...spawnForm, room_id: parseInt(e.target.value) || 0 })
              }
              placeholder="Enter room number"
              className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
            />
          </div>

          <div>
            <label className="text-text-muted text-xs block mb-1">Instance Number *</label>
            <input
              type="text"
              inputMode="numeric"
              value={spawnForm.instance_number || ''}
              onChange={(e) =>
                setSpawnForm({ ...spawnForm, instance_number: parseInt(e.target.value) || 1 })
              }
              placeholder="1"
              className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
            />
          </div>

          <div>
            <label className="text-text-muted text-xs block mb-1">
              Instance Name <span className="text-text-muted">(optional, defaults to template name)</span>
            </label>
            <input
              type="text"
              value={spawnForm.instance_name}
              onChange={(e) =>
                setSpawnForm({ ...spawnForm, instance_name: e.target.value })
              }
              placeholder={template.name}
              className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
            />
          </div>

          {spawnMutation.isError && (
            <div className="p-2 bg-danger/10 border border-danger rounded text-danger text-xs">
              Failed to spawn instance: {spawnMutation.error?.message ?? 'Unknown error'}
            </div>
          )}

          <div className="flex gap-2">
            <Button variant="primary" onClick={handleSpawn} disabled={spawnMutation.isPending}>
              {spawnMutation.isPending ? 'Spawning...' : 'Spawn'}
            </Button>
            <Button variant="secondary" onClick={handleSpawnModalClose}>
              Cancel
            </Button>
          </div>
        </div>
      </Modal>
    </div>
  )
}

function DetailField({ label, value }: Readonly<{ label: string; value: string }>) {
  return (
    <div>
      <span className="text-text-muted text-xs block mb-0.5">{label}</span>
      <span className="text-text text-sm font-medium">{value}</span>
    </div>
  )
}

function TemplateEditForm({ form, onChange, saveError }: Readonly<{
  form: EditForm;
  onChange: (val: EditForm) => void;
  saveError: string | null;
}>) {
  return (
    <div className="space-y-4">
      <div className="grid grid-cols-2 gap-4">
        <FormField
          label="Name"
          value={form.name}
          onChange={(val) => onChange({ ...form, name: val })}
        />
        <FormField
          label="Race"
          value={form.race}
          onChange={(val) => onChange({ ...form, race: val })}
        />
        <FormField
          label="Disposition"
          value={form.disposition}
          onChange={(val) => onChange({ ...form, disposition: val })}
        />
        <NumberField
          label="Level"
          value={form.level}
          onChange={(val) => onChange({ ...form, level: val })}
        />
        <NumberField
          label="XP Value"
          value={form.xp_value}
          onChange={(val) => onChange({ ...form, xp_value: val })}
        />
        <NumberField
          label="Respawn Cooldown (s)"
          value={form.respawn_cooldown}
          onChange={(val) => onChange({ ...form, respawn_cooldown: val })}
        />
        <TextareaField
          label="Greeting"
          value={form.greeting}
          onChange={(val) => onChange({ ...form, greeting: val })}
        />
        <TextareaField
          label="Trades With (one per line)"
          value={form.trades_with}
          onChange={(val) => onChange({ ...form, trades_with: val })}
        />
        <TextareaField
          label="Respawn Rooms (one per line)"
          value={form.respawn_rooms}
          onChange={(val) => onChange({ ...form, respawn_rooms: val })}
        />
        <TextareaField
          label="Skills (skill:value, one per line)"
          value={form.skills}
          onChange={(val) => onChange({ ...form, skills: val })}
        />
      </div>
      <TextareaField
        label="Description"
        value={form.description}
        onChange={(val) => onChange({ ...form, description: val })}
      />
      {saveError && <FormError message={saveError} />}
    </div>
  )
}
