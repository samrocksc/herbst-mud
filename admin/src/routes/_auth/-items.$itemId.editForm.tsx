import { useState } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useNavigate } from '@tanstack/react-router'
import { apiPut, apiDelete } from '../../utils/apiFetch'
import { Button } from '../../components/Button'
import { CombatFieldsEditor, type CombatFields } from '../../components/CombatFieldsEditor'
import { NumberField, SelectField, CheckboxField, TextareaField, FormError } from '../../components/FormFields'
import { SLOT_OPTIONS, ITEM_TYPE_OPTIONS } from '../../components/itemConstants'

export type TemplateEditForm = Readonly<{
  name: string; description: string; slot: string; level: number; weight: number
  item_type: string; color: string; is_visible: boolean; is_immovable: boolean
  effect_type: string; effect_value: number; effect_duration: number
  is_container: boolean; container_capacity: number; is_locked: boolean
  key_item_id: string; reveal_condition: string
}> & CombatFields

export function TemplateEditForm({ template, itemId, onDone }: Readonly<{
  template: TemplateEditForm; itemId: string; onDone: () => void
}>) {
  const queryClient = useQueryClient()
  const navigate = useNavigate()
  const [form, setForm] = useState<TemplateEditForm>(() => ({ ...template }))
  const [deleteConfirm, setDeleteConfirm] = useState(false)
  const API = `${window.location.origin}`

  const updateMutation = useMutation({
    mutationFn: (body: Record<string, unknown>) => apiPut(`${API}/api/equipment-templates/${itemId}`, body),
    onSuccess: () => { queryClient.invalidateQueries({ queryKey: ['item-template', itemId] }); onDone() },
  })

  const deleteMutation = useMutation({
    mutationFn: () => apiDelete(`${API}/api/equipment-templates/${itemId}`),
    onSuccess: () => { queryClient.invalidateQueries({ queryKey: ['item-templates'] }); navigate({ to: '/items' }) },
  })

  const set = <K extends keyof TemplateEditForm>(key: K, value: TemplateEditForm[K]) =>
    setForm(prev => ({ ...prev, [key]: value }))

  const handleSave = () => {
    const body: Record<string, unknown> = { ...form }
    try { body.stats = JSON.parse(form.stats) } catch { body.stats = {} }
    updateMutation.mutate(body)
  }

  return (
    <div className="bg-surface-muted rounded-lg p-6 border border-border mb-6">
      <h2 className="mt-0 mb-4 text-text text-lg font-semibold">Edit Template</h2>
      <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
        <div><label className="text-text-muted text-xs block mb-1">Name</label>
          <input type="text" value={form.name} onChange={(e) => set('name', e.target.value)}
            className="w-full p-2 bg-surface border border-border rounded text-text text-sm" /></div>
        <SelectField label="Slot" value={form.slot} onChange={(v) => set('slot', v)} options={[...SLOT_OPTIONS]} />
        <SelectField label="Item Type" value={form.item_type} onChange={(v) => set('item_type', v)} options={[...ITEM_TYPE_OPTIONS]} />
        <NumberField label="Level" value={form.level} onChange={(v) => set('level', v)} />
        <NumberField label="Weight" value={form.weight} onChange={(v) => set('weight', v)} />
        <div><label className="text-text-muted text-xs block mb-1">Color</label>
          <input type="text" value={form.color} onChange={(e) => set('color', e.target.value)}
            className="w-full p-2 bg-surface border border-border rounded text-text text-sm" /></div>
        <CheckboxField label="Visible" checked={form.is_visible} onChange={(v) => set('is_visible', v)} />
        <CheckboxField label="Immovable" checked={form.is_immovable} onChange={(v) => set('is_immovable', v)} />
        <CheckboxField label="Container" checked={form.is_container} onChange={(v) => set('is_container', v)} />
      </div>
      <TextareaField label="Description" value={form.description} onChange={(v) => set('description', v)} rows={3} />
      <div className="mt-4 pt-4 border-t border-border">
        <CombatFieldsEditor form={form} onChange={(u) => setForm(prev => ({ ...prev, ...u }))} slot={form.slot} />
      </div>
      {updateMutation.isError && <FormError message={(updateMutation.error as Error)?.message ?? 'Failed to save'} />}
      <div className="flex gap-2 mt-4">
        <Button variant="primary" onClick={handleSave} disabled={updateMutation.isPending}>{updateMutation.isPending ? 'Saving...' : 'Save'}</Button>
        <Button variant="secondary" onClick={onDone}>Cancel</Button>
        <Button variant="danger" onClick={() => deleteConfirm ? deleteMutation.mutate() : setDeleteConfirm(true)} disabled={deleteMutation.isPending}>
          {deleteConfirm ? 'Confirm Delete?' : 'Delete'}</Button>
        {deleteConfirm && <Button variant="secondary" onClick={() => setDeleteConfirm(false)}>Cancel Delete</Button>}
      </div>
      {deleteMutation.isError && <div className="mt-2 text-danger text-sm">Failed to delete: {(deleteMutation.error as Error)?.message}</div>}
    </div>
  )
}