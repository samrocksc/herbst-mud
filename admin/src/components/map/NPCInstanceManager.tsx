import { Button } from '../Button'
import { NPCSpawnModal } from './NPCSpawnModal'
import { NPCInstanceRow } from './NPCEditRow'
import { useNPCInstances } from './useNPCInstances'

export type NPCInstanceView = Readonly<{
  id: number
  name: string
  npc_template_id: string
  instance_number: number
  room_id: number
  starting_room_id: number
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

export type NPCTemplate = Readonly<{
  id: string
  name: string
  race: string
  level: number
  respawn_rooms: string[]
  respawn_cooldown: number
}>

export type SpawnFormData = {
  template_id: string
  level: number
  hitpoints: number
  room_id: number
  respawn_cooldown: number
  respawn_rooms: string
}

export type EditFormData = {
  level: number
  hitpoints: number
  room_id: number
  starting_room_id: number
}

type Props = Readonly<{ roomId: number }>

export function NPCInstanceManager({ roomId }: Props) {
  const {
    instancesQuery, templatesQuery, createMutation, updateMutation,
    showSpawn, setShowSpawn, editingId, setEditingId, confirmDeleteId, setConfirmDeleteId,
    spawnForm, setSpawnForm, editForm, setEditForm,
    handleSpawn, handleUpdate, handleDelete, startEdit,
  } = useNPCInstances(roomId)

  if (instancesQuery.isLoading) {
    return <div className="mb-3"><strong className="text-warning text-xs">NPCs:</strong><div className="text-text-muted text-[10px] mt-1">Loading...</div></div>
  }
  if (instancesQuery.error) {
    return <div className="mb-3"><strong className="text-warning text-xs">NPCs:</strong><div className="text-danger text-[10px] mt-1">Error loading NPCs</div></div>
  }

  const instances = instancesQuery.data ?? []
  const templates = templatesQuery.data ?? []

  return (
    <div className="mb-3">
      <div className="flex items-center justify-between mb-1">
        <strong className="text-warning text-xs">NPCs:</strong>
        <Button variant="primary" size="sm" className="!px-1.5 !py-0 !text-[10px]"
          onClick={() => { setSpawnForm({ template_id: '', level: 0, hitpoints: 0, room_id: roomId, respawn_cooldown: 0, respawn_rooms: '' }); setShowSpawn(true) }}>
          + Spawn
        </Button>
      </div>
      {instances.length === 0 ? (
        <div className="text-text-muted text-[10px]">No NPCs in this room.</div>
      ) : (
        <div className="mt-1 flex flex-col gap-1">
          {instances.map((inst) => (
            <NPCInstanceRow key={inst.id} inst={inst} editingId={editingId} confirmDeleteId={confirmDeleteId}
              editForm={editForm} setEditForm={setEditForm} startEdit={startEdit}
              handleUpdate={handleUpdate} handleDelete={handleDelete}
              setEditingId={setEditingId} setConfirmDeleteId={setConfirmDeleteId}
              isUpdatePending={updateMutation.isPending} updateError={updateMutation.error} />
          ))}
        </div>
      )}
      <NPCSpawnModal showSpawn={showSpawn} setShowSpawn={setShowSpawn} spawnForm={spawnForm}
        setSpawnForm={setSpawnForm} templates={templates} templatesLoading={templatesQuery.isLoading}
        onSpawn={handleSpawn} isPending={createMutation.isPending} error={createMutation.error}
        selectedTemplate={spawnForm.template_id ? templates.find((t) => t.id === spawnForm.template_id) ?? null : null} />
    </div>
  )
}

