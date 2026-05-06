import { useState, useCallback } from 'react'
import { Button } from './Button'
import { Modal } from './Modal'
import {
  useNPCInstances,
  useCreateNPCInstance,
  useUpdateNPCInstance,
  useDeleteNPCInstance,
} from '../hooks/useNPCInstances'
import type { NPCInstance, NPCInstanceInput, NPCInstanceUpdate } from '../hooks/useNPCInstances'
import { SearchableSelect } from './SearchableSelect'
import { logError } from '../utils/log'

// ─── Types ──────────────────────────────────────────────────────────────────

type NPCTemplate = Readonly<{
  id: string
  name: string
}>

type Room = Readonly<{
  id: number
  name: string
}>

type NPCInstanceManagerProps = Readonly<{
  roomId?: number
  rooms: Room[]
  templates: NPCTemplate[]
  onSelectRoom?: (roomId: number) => void
}>

// ─── Component ──────────────────────────────────────────────────────────────

export function NPCInstanceManager({
  roomId,
  rooms,
  templates,
  onSelectRoom,
}: NPCInstanceManagerProps) {
  const { data: instances = [], isLoading, error } = useNPCInstances(roomId)
  const createMutation = useCreateNPCInstance()
  const updateMutation = useUpdateNPCInstance()
  const deleteMutation = useDeleteNPCInstance()

  const [showCreate, setShowCreate] = useState(false)
  const [editingId, setEditingId] = useState<number | null>(null)
  const [confirmDeleteId, setConfirmDeleteId] = useState<number | null>(null)
  const [createForm, setCreateForm] = useState<NPCInstanceInput>({
    template_id: '',
    room_id: roomId ?? 0,
  })
  const [editForm, setEditForm] = useState<NPCInstanceUpdate>({})

  const getRoomName = useCallback(
    (id: number) => {
      const room = rooms.find((r) => r.id === id)
      return room ? room.name : `Room ${id}`
    },
    [rooms],
  )

  const getTemplateName = useCallback(
    (templateId: string) => {
      const tmpl = templates.find((t) => t.id === templateId)
      return tmpl ? tmpl.name : templateId
    },
    [templates],
  )

  const handleCreate = useCallback(async () => {
    if (!createForm.template_id || !createForm.room_id) return
    try {
      await createMutation.mutateAsync(createForm)
      setShowCreate(false)
      setCreateForm({ template_id: '', room_id: roomId ?? 0 })
    } catch (err) {
      logError('Create NPC instance:', err)
    }
  }, [createForm, createMutation, roomId])

  const startEdit = useCallback((inst: NPCInstance) => {
    setEditingId(inst.id)
    setConfirmDeleteId(null)
    setEditForm({
      room_id: inst.room_id,
      starting_room_id: inst.starting_room_id,
      hitpoints: inst.hitpoints,
    })
  }, [])

  const handleUpdate = useCallback(async () => {
    if (editingId === null) return
    try {
      await updateMutation.mutateAsync({ id: editingId, update: editForm })
      setEditingId(null)
      setEditForm({})
    } catch (err) {
      logError('Update NPC instance:', err)
    }
  }, [editingId, editForm, updateMutation])

  const handleDelete = useCallback(
    async (id: number) => {
      try {
        await deleteMutation.mutateAsync(id)
        setConfirmDeleteId(null)
        if (editingId === id) {
          setEditingId(null)
          setEditForm({})
        }
      } catch (err) {
        logError('Delete NPC instance:', err)
      }
    },
    [deleteMutation, editingId],
  )

  if (isLoading) {
    return (
      <div className="flex-1 flex items-center justify-center text-text-muted text-sm p-4">
        Loading NPC instances...
      </div>
    )
  }

  if (error) {
    return (
      <div className="flex-1 flex items-center justify-center text-danger text-sm p-4">
        Error loading NPC instances
      </div>
    )
  }

  return (
    <div className="flex flex-col h-full">
      {/* Header */}
      <div className="p-3 border-b border-border flex justify-between items-center">
        <h3 className="m-0 text-text text-sm font-semibold">
          NPC Instances{roomId ? ` (Room ${roomId})` : ''}
        </h3>
        <Button variant="primary" size="sm" onClick={() => setShowCreate(true)}>
          + Spawn
        </Button>
      </div>

      {/* Instance list */}
      <div className="flex-1 overflow-y-auto p-2">
        {instances.length === 0 ? (
          <div className="text-text-muted text-xs text-center py-4">
            No NPC instances{roomId ? ' in this room' : ''}.
          </div>
        ) : (
          <div className="flex flex-col gap-1">
            {instances.map((inst) => (
              <div
                key={inst.id}
                className="p-2 bg-surface-muted border border-border rounded text-xs"
              >
                {editingId === inst.id ? (
                  /* Edit view */
                  <div className="space-y-2">
                    <div className="font-medium text-text">{inst.name}</div>
                    <div>
                      <label className="text-text-muted block mb-1">Room</label>
                      <select
                        value={editForm.room_id ?? inst.room_id}
                        onChange={(e) =>
                          setEditForm((f) => ({ ...f, room_id: parseInt(e.target.value, 10) || 0 }))
                        }
                        className="w-full p-1 bg-surface border border-border rounded text-text text-xs"
                      >
                        {rooms.map((r) => (
                          <option key={r.id} value={r.id}>
                            {r.name}
                          </option>
                        ))}
                      </select>
                    </div>
                    <div>
                      <label className="text-text-muted block mb-1">Starting Room</label>
                      <select
                        value={editForm.starting_room_id ?? inst.starting_room_id}
                        onChange={(e) =>
                          setEditForm((f) => ({ ...f, starting_room_id: parseInt(e.target.value, 10) || 0 }))
                        }
                        className="w-full p-1 bg-surface border border-border rounded text-text text-xs"
                      >
                        {rooms.map((r) => (
                          <option key={r.id} value={r.id}>
                            {r.name}
                          </option>
                        ))}
                      </select>
                    </div>
                    <div>
                      <label className="text-text-muted block mb-1">Hitpoints</label>
                      <input
                        type="number"
                        value={editForm.hitpoints ?? inst.hitpoints}
                        onChange={(e) =>
                          setEditForm((f) => ({
                            ...f,
                            hitpoints: parseInt(e.target.value, 10) || 0,
                          }))
                        }
                        className="w-full p-1 bg-surface border border-border rounded text-text text-xs"
                      />
                    </div>
                    <div className="flex gap-1">
                      <Button
                        variant="primary"
                        size="sm"
                        onClick={handleUpdate}
                        disabled={updateMutation.isPending}
                      >
                        {updateMutation.isPending ? 'Saving...' : 'Save'}
                      </Button>
                      <Button variant="ghost" size="sm" onClick={() => setEditingId(null)}>
                        Cancel
                      </Button>
                    </div>
                  </div>
                ) : (
                  /* Detail view */
                  <>
                    <div className="flex justify-between items-start">
                      <div>
                        <div className="font-medium text-text">{inst.name}</div>
                        <div className="text-text-muted text-[10px]">
                          #{inst.id} &middot; Instance {inst.instance_number} &middot;{' '}
                          {getTemplateName(inst.npc_template_id)}
                        </div>
                      </div>
                      <div className="flex gap-1">
                        <Button variant="ghost" size="sm" className="!px-1 !py-0" onClick={() => startEdit(inst)}>
                          ✏️
                        </Button>
                        <Button
                          variant={confirmDeleteId === inst.id ? 'secondary' : 'ghost'}
                          size="sm"
                          className="!px-1 !py-0"
                          onClick={() => {
                            if (confirmDeleteId === inst.id) {
                              handleDelete(inst.id)
                            } else {
                              setConfirmDeleteId(inst.id)
                            }
                          }}
                        >
                          {confirmDeleteId === inst.id ? '❓' : '🗑'}
                        </Button>
                      </div>
                    </div>

                    <div className="mt-1 flex gap-2 text-text-muted">
                      <span>
                        HP: <span className={inst.hitpoints < inst.max_hitpoints * 0.3 ? 'text-danger' : 'text-text'}>{inst.hitpoints}/{inst.max_hitpoints}</span>
                      </span>
                      <span>Lv {inst.level}</span>
                      <span>{inst.race}</span>
                    </div>

                    <div
                      className="mt-1 text-text-muted cursor-pointer hover:text-text transition-colors"
                      onClick={() => onSelectRoom?.(inst.room_id)}
                    >
                      📍 {getRoomName(inst.room_id)}
                    </div>

                    {confirmDeleteId === inst.id && (
                      <div className="mt-1 p-1 bg-danger/10 border border-danger rounded text-[10px] text-text">
                        Confirm delete?{' '}
                        <Button variant="danger" size="sm" className="!px-1 !py-0 !text-[10px]" onClick={() => handleDelete(inst.id)}>
                          Yes
                        </Button>{' '}
                        <Button variant="ghost" size="sm" className="!px-1 !py-0 !text-[10px]" onClick={() => setConfirmDeleteId(null)}>
                          No
                        </Button>
                      </div>
                    )}
                  </>
                )}
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Create modal */}
      <Modal isOpen={showCreate} onClose={() => setShowCreate(false)} title="Spawn NPC Instance">
        <div className="space-y-3">
          {createMutation.isError && (
            <div className="p-2 bg-danger/10 text-danger rounded text-xs">
              {(createMutation.error as Error)?.message || 'Failed to spawn NPC instance'}
            </div>
          )}

          <div>
            <label className="text-text-muted text-xs block mb-1">NPC Template *</label>
            <SearchableSelect
              options={templates}
              value={createForm.template_id}
              onChange={(id) => setCreateForm((f) => ({ ...f, template_id: id }))}
              placeholder="Search by name or ID..."
              disabled={createMutation.isPending}
            />
          </div>

          <div>
            <label className="text-text-muted text-xs block mb-1">Room *</label>
            <select
              value={createForm.room_id}
              onChange={(e) => setCreateForm((f) => ({ ...f, room_id: parseInt(e.target.value, 10) || 0 }))}
              className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
            >
              <option value={0}>-- Select room --</option>
              {rooms.map((r) => (
                <option key={r.id} value={r.id}>
                  {r.name}
                </option>
              ))}
            </select>
          </div>

          <div>
            <label className="text-text-muted text-xs block mb-1">
              Instance Number <span className="text-text-muted">(0 = auto)</span>
            </label>
            <input
              type="number"
              value={createForm.instance_number ?? 0}
              onChange={(e) =>
                setCreateForm((f) => ({
                  ...f,
                  instance_number: parseInt(e.target.value, 10) || 0,
                }))
              }
              min={0}
              className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
            />
          </div>

          <div className="flex gap-2 pt-2">
            <Button
              variant="primary"
              onClick={handleCreate}
              disabled={createMutation.isPending || !createForm.template_id || !createForm.room_id}
            >
              {createMutation.isPending ? 'Spawning...' : 'Spawn Instance'}
            </Button>
            <Button variant="secondary" onClick={() => setShowCreate(false)}>
              Cancel
            </Button>
          </div>
        </div>
      </Modal>
    </div>
  )
}