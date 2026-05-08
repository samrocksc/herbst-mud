import { Modal } from '../Modal'
import { Button } from '../Button'
import { FormField } from '../fields/FormField'
import { NumberField } from '../fields/NumberField'
import { SelectField } from '../fields/SelectField'
import { FormError } from '../fields/FormError'
import type { SpawnFormData, EquipmentTemplate } from './types'

type ItemSpawnModalProps = Readonly<{
  isOpen: boolean
  onClose: () => void
  spawnForm: SpawnFormData
  setSpawnForm: React.Dispatch<React.SetStateAction<SpawnFormData>>
  onSpawn: () => void
  isPending: boolean
  error: Error | null
  templates: EquipmentTemplate[]
  templatesLoading: boolean
  selectedTemplate: EquipmentTemplate | undefined
  applyTemplateDefaults: (templateId: string) => void
}>

export function ItemSpawnModal({
  isOpen, onClose, spawnForm, setSpawnForm, onSpawn,
  isPending, error, templates, templatesLoading,
  selectedTemplate, applyTemplateDefaults,
}: ItemSpawnModalProps) {
  const templateOptions = templates.map((t) => ({
    value: t.equipment_template_id,
    label: `${t.name} (${t.slot}, lv.${t.level})`,
  }))

  return (
    <Modal isOpen={isOpen} onClose={onClose} title="Spawn Item Instance">
      <div className="space-y-3">
        {error && <FormError message={error.message || 'Failed to spawn item instance'} />}
        <SelectField
          label="Equipment Template *"
          value={spawnForm.template_id}
          onChange={applyTemplateDefaults}
          options={templateOptions}
          placeholder="-- Select template --"
          disabled={templatesLoading}
        />
        {selectedTemplate && (
          <div className="p-2 bg-surface-muted border border-border rounded text-xs text-text-muted space-y-0.5">
            <div>Slot: {selectedTemplate.slot}</div>
            <div>Level: {selectedTemplate.level} Weight: {selectedTemplate.weight}</div>
            <div>Type: {selectedTemplate.item_type}</div>
          </div>
        )}
        <FormField label="Name" value={spawnForm.name} onChange={(v) => setSpawnForm((f) => ({ ...f, name: v }))} />
        <FormField label="Description" value={spawnForm.description} onChange={(v) => setSpawnForm((f) => ({ ...f, description: v }))} />
        <div className="flex gap-2">
          <div className="flex-1">
            <FormField label="Slot" value={spawnForm.slot} onChange={(v) => setSpawnForm((f) => ({ ...f, slot: v }))} />
          </div>
          <div className="flex-1">
            <NumberField label="Level" value={spawnForm.level} onChange={(v) => setSpawnForm((f) => ({ ...f, level: v }))} min={0} />
          </div>
        </div>
        <div className="flex gap-2">
          <div className="flex-1">
            <NumberField label="Weight" value={spawnForm.weight} onChange={(v) => setSpawnForm((f) => ({ ...f, weight: v }))} min={0} />
          </div>
          <div className="flex-1">
            <FormField label="Color" value={spawnForm.color} onChange={(v) => setSpawnForm((f) => ({ ...f, color: v }))} placeholder="#8b5cf6" />
          </div>
        </div>
        <NumberField label="Room ID" value={spawnForm.room_id} onChange={(v) => setSpawnForm((f) => ({ ...f, room_id: v }))} />
        <div className="flex gap-2 pt-2">
          <Button variant="primary" onClick={onSpawn} disabled={isPending || !spawnForm.template_id}>
            {isPending ? 'Spawning...' : 'Spawn Instance'}
          </Button>
          <Button variant="secondary" onClick={onClose}>Cancel</Button>
        </div>
      </div>
    </Modal>
  )
}