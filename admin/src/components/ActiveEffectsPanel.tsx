import { useState } from 'react'
import { useActiveEffects, useRemoveActiveEffect, useApplyEffect, type ActiveEffect } from '../hooks/useActiveEffects'
import { useEffectDefs } from '../hooks/useEffectDefs'
import { Button } from './Button'
import { DeleteConfirmation } from './DeleteConfirmation'
import { showToast } from './Toast'
import { SelectField } from './fields'

function formatDuration(startedAt: string, expiresAt: string | null): string {
  if (!expiresAt) return 'Permanent'
  const now = Date.now()
  const end = new Date(expiresAt).getTime()
  const remaining = Math.max(0, end - now)
  if (remaining === 0) return 'Expired'
  const mins = Math.floor(remaining / 60000)
  const secs = Math.floor((remaining % 60000) / 1000)
  return mins > 0 ? `${mins}m ${secs}s` : `${secs}s`
}

function formatParams(params: Record<string, unknown>): string {
  const entries = Object.entries(params)
  if (entries.length === 0) return '—'
  return entries.map(([k, v]) => `${k}=${v}`).join(', ')
}

export function ActiveEffectsPanel({ characterId }: { characterId: number }) {
  const { data: effects = [], isLoading, error } = useActiveEffects(characterId)
  const { data: effectDefs = [] } = useEffectDefs()
  const remove = useRemoveActiveEffect()
  const apply = useApplyEffect()
  const [confirmDelete, setConfirmDelete] = useState<number | null>(null)
  const [selectedEffectId, setSelectedEffectId] = useState<number>(effectDefs[0]?.id ?? 0)
  const [showApply, setShowApply] = useState(false)

  const handleApply = () => {
    if (!selectedEffectId) return
    apply.mutate({ characterId, effectId: selectedEffectId }, {
      onSuccess: () => { setShowApply(false); showToast('Effect applied', 'success') },
    })
  }

  const handleRemove = (effectId: number) => {
    remove.mutate({ characterId, effectId }, {
      onSuccess: () => { setConfirmDelete(null); showToast('Effect removed', 'success') },
    })
  }

  return (
    <div className="mt-6">
      <div className="flex items-center justify-between mb-4">
        <h2 className="m-0 text-text text-lg font-semibold">
          Active Effects{effects.length > 0 ? ` (${effects.length})` : ''}
        </h2>
        <Button variant="primary" size="sm" onClick={() => setShowApply(!showApply)}>+ Apply Effect</Button>
      </div>

      {showApply && (
        <div className="bg-surface rounded border border-border p-4 mb-4 flex items-end gap-3">
          <div className="flex-1">
            <SelectField
              label="Effect"
              value={String(selectedEffectId)}
              onChange={(v) => setSelectedEffectId(Number(v))}
              options={effectDefs.map((e) => ({ value: String(e.id), label: `${e.name} (${e.effect_type})` }))}
            />
          </div>
          <Button variant="primary" onClick={handleApply} disabled={apply.isPending}>
            {apply.isPending ? 'Applying…' : 'Apply'}
          </Button>
        </div>
      )}

      {error && <div className="error-banner mb-3">{error instanceof Error ? error.message : 'Failed to load effects'}</div>}

      {isLoading ? (
        <div className="text-text-muted text-sm py-4 text-center">Loading effects…</div>
      ) : effects.length === 0 ? (
        <div className="text-text-muted text-sm py-4 text-center">No active effects.</div>
      ) : (
        <div className="space-y-2">
          {effects.map((ae: ActiveEffect) => (
            <div key={ae.id} className="bg-surface-muted rounded border border-border p-3 flex items-center justify-between gap-4">
              <div className="flex-1 min-w-0">
                <div className="flex items-center gap-2">
                  <span className="font-medium text-text">{ae.effect_name}</span>
                  <code className="text-xs text-accent bg-surface px-1.5 py-0.5 rounded">{ae.effect_type}</code>
                  {ae.stack_count > 1 && <span className="text-xs text-text-muted">×{ae.stack_count}</span>}
                  <span className="text-xs text-text-muted">{formatDuration(ae.started_at, ae.expires_at)}</span>
                </div>
                <div className="text-xs text-text-muted mt-0.5">{formatParams(ae.parameters)}</div>
              </div>
              <Button variant="danger" size="sm" onClick={() => setConfirmDelete(ae.id)}>Remove</Button>
            </div>
          ))}
        </div>
      )}

      {confirmDelete !== null && (
        <DeleteConfirmation
          open={confirmDelete !== null}
          title="Remove Effect"
          message="Remove this active effect from the character?"
          onConfirm={() => handleRemove(confirmDelete)}
          onCancel={() => setConfirmDelete(null)}
          isLoading={remove.isPending}
        />
      )}
    </div>
  )
}