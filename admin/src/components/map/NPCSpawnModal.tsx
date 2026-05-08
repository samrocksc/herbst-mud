import { Button } from '../Button'
import { Modal } from '../Modal'
import { SearchableSelect } from '../SearchableSelect'
import { NumberField } from '../fields/NumberField'
import { FormField } from '../fields/FormField'
import { FormError } from '../fields/FormError'
import type { NPCTemplate, SpawnFormData } from './NPCInstanceManager'

type NPCSpawnModalProps = Readonly<{
  showSpawn: boolean
  setShowSpawn: (v: boolean) => void
  spawnForm: SpawnFormData
  setSpawnForm: React.Dispatch<React.SetStateAction<SpawnFormData>>
  templates: NPCTemplate[]
  templatesLoading: boolean
  onSpawn: () => void
  isPending: boolean
  error: Error | null
  selectedTemplate: NPCTemplate | null
}>

export function NPCSpawnModal({
  showSpawn, setShowSpawn, spawnForm, setSpawnForm,
  templates, templatesLoading, onSpawn, isPending, error, selectedTemplate,
}: NPCSpawnModalProps) {
  if (!showSpawn) return null
  const upd = (patch: Partial<SpawnFormData>) => setSpawnForm((f) => ({ ...f, ...patch }))

  return (
    <Modal isOpen={showSpawn} onClose={() => setShowSpawn(false)} title="Spawn NPC Instance">
      <div className="space-y-3">
        {error && <FormError message={(error as Error)?.message || 'Failed to spawn NPC instance'} />}
        <div>
          <label className="text-text-muted text-xs block mb-1">NPC Template *</label>
          {templatesLoading ? (
            <div className="text-text-muted text-xs">Loading templates...</div>
          ) : (
            <SearchableSelect options={templates} value={spawnForm.template_id}
              onChange={(id) => upd({ template_id: id })} placeholder="Search by name or ID..." disabled={isPending} />
          )}
        </div>
        {selectedTemplate && (
          <div className="p-2 bg-surface-muted border border-border rounded text-xs text-text-muted space-y-0.5">
            <div>Respawn cooldown: {selectedTemplate.respawn_cooldown}s</div>
            <div>Respawn rooms: {selectedTemplate.respawn_rooms.length > 0 ? selectedTemplate.respawn_rooms.join(', ') : 'none'}</div>
          </div>
        )}
        <div className="flex gap-2">
          <NumberField label="Level (0 = default)" value={spawnForm.level} onChange={(v) => upd({ level: v })} min={0} disabled={isPending} />
          <NumberField label="HP (0 = default)" value={spawnForm.hitpoints} onChange={(v) => upd({ hitpoints: v })} min={0} disabled={isPending} />
        </div>
        <NumberField label="Room ID" value={spawnForm.room_id} onChange={(v) => upd({ room_id: v })} disabled={isPending} />
        <div className="flex gap-2">
          <NumberField label="Respawn Cooldown (0 = default)" value={spawnForm.respawn_cooldown} onChange={(v) => upd({ respawn_cooldown: v })} min={0} disabled={isPending} />
          <FormField label="Respawn Rooms (empty = default)" value={spawnForm.respawn_rooms} onChange={(v) => upd({ respawn_rooms: v })} placeholder="e.g. 101, 102, 103" disabled={isPending} />
        </div>
        <div className="flex gap-2 pt-2">
          <Button variant="primary" onClick={onSpawn} disabled={isPending || !spawnForm.template_id}>
            {isPending ? 'Spawning...' : 'Spawn Instance'}
          </Button>
          <Button variant="secondary" onClick={() => setShowSpawn(false)}>Cancel</Button>
        </div>
      </div>
    </Modal>
  )
}