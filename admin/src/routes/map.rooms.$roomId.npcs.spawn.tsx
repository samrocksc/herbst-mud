import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { apiGet, apiPost } from '../utils/apiFetch'
import { PageHeader } from '../components/PageHeader'
import { Button } from '../components/Button'
import { NumberField } from '../components/fields/NumberField'
import { FormField } from '../components/fields/FormField'
import { FormError } from '../components/fields/FormError'
import { ResourceSearchSelect } from '../components/ResourceSearchSelect'
import { SearchableSelect } from '../components/SearchableSelect'
import { RESOURCE_ENDPOINTS } from '../utils/resourceEndpoints'

const API = `${window.location.origin}/api`

type NPCTemplate = Readonly<{
  id: string
  name: string
  race: string
  level: number
  respawn_rooms: string[]
  respawn_cooldown: number
}>

function parseRoomIds(raw: string): string[] {
  return raw.split(',').map((s) => s.trim()).filter((s) => s.length > 0)
}

export const Route = createFileRoute('/map/rooms/$roomId/npcs/spawn')({
  component: NPCSpawnPage,
})

function NPCSpawnPage() {
  const { roomId } = Route.useParams()
  const roomIdNum = Number(roomId)
  const navigate = useNavigate()
  const queryClient = useQueryClient()

  const { data: templates = [], isLoading: templatesLoading } = useQuery({
    queryKey: ['npc-templates'],
    queryFn: () => apiGet<NPCTemplate[]>(`${API}/npc-templates`),
  })

  const createMutation = useMutation({
    mutationFn: (input: Record<string, unknown>) => apiPost(`${API}/npc-instances`, input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['npc-instances'] })
      navigate({ to: '/map', search: { room: roomIdNum } })
    },
  })

  const [templateId, setTemplateId] = useState('')
  const [level, setLevel] = useState(0)
  const [hitpoints, setHitpoints] = useState(0)
  const [spawnRoomId, setSpawnRoomId] = useState<number | string | null>(roomIdNum)
  const [respawnCooldown, setRespawnCooldown] = useState(0)
  const [respawnRooms, setRespawnRooms] = useState('')

  const selectedTemplate = templateId ? templates.find((t) => t.id === templateId) ?? null : null

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!templateId) return
    const payload: Record<string, unknown> = { template_id: templateId, room_id: Number(spawnRoomId) || roomIdNum }
    if (level > 0) payload.level = level
    if (hitpoints > 0) payload.hitpoints = hitpoints
    if (respawnCooldown > 0) payload.respawn_cooldown = respawnCooldown
    const parsedRooms = parseRoomIds(respawnRooms)
    if (parsedRooms.length > 0) payload.respawn_rooms = parsedRooms
    createMutation.mutate(payload)
  }

  const npcTemplateOptions = templates.map((t) => ({ id: t.id, name: `${t.name} (${t.race}, lv.${t.level})` }))

  return (
    <div className="p-6">
      <PageHeader title="Spawn NPC Instance" showBack backTo="/map" />
      <div className="bg-surface p-6 border border-border rounded max-w-[600px]">
        <form onSubmit={handleSubmit} className="space-y-3">
          {createMutation.isError && (
            <FormError message={(createMutation.error as Error)?.message || 'Failed to spawn NPC instance'} />
          )}
          <div>
            <label className="text-text-muted text-xs block mb-1">NPC Template *</label>
            {templatesLoading ? (
              <div className="text-text-muted text-xs">Loading templates...</div>
            ) : (
              <SearchableSelect options={npcTemplateOptions} value={templateId}
                onChange={setTemplateId} placeholder="Search by name or ID..." />
            )}
          </div>
          {selectedTemplate && (
            <div className="p-2 bg-surface-muted border border-border rounded text-xs text-text-muted space-y-0.5">
              <div>Respawn cooldown: {selectedTemplate.respawn_cooldown}s</div>
              <div>Respawn rooms: {(selectedTemplate.respawn_rooms ?? []).length > 0 ? selectedTemplate.respawn_rooms.join(', ') : 'none'}</div>
            </div>
          )}
          <div className="flex gap-2">
            <NumberField label="Level (0 = default)" value={level} onChange={setLevel} min={0} />
            <NumberField label="HP (0 = default)" value={hitpoints} onChange={setHitpoints} min={0} />
          </div>
          <ResourceSearchSelect
            label="Room"
            value={spawnRoomId}
            onChange={(id) => setSpawnRoomId(id)}
            {...RESOURCE_ENDPOINTS.rooms}
            placeholder="Search rooms by name or ID..."
          />
          <div className="flex gap-2">
            <NumberField label="Respawn Cooldown (0 = default)" value={respawnCooldown} onChange={setRespawnCooldown} min={0} />
            <FormField label="Respawn Rooms (empty = default)" value={respawnRooms} onChange={setRespawnRooms} placeholder="e.g. 101, 102, 103" />
          </div>
          <div className="flex gap-2 pt-2">
            <Button type="submit" variant="primary" disabled={!templateId || createMutation.isPending}>
              {createMutation.isPending ? 'Spawning...' : 'Spawn Instance'}
            </Button>
            <Button type="button" variant="secondary" onClick={() => navigate({ to: '/map', search: { room: roomIdNum } })}>
              Cancel
            </Button>
          </div>
        </form>
      </div>
    </div>
  )
}
