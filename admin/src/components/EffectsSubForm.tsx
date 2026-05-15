import { useState } from 'react';
import {
  useEffects,
  useCreateEffect,
  useDeleteEffect,
  type AbilityEffect,
  type EffectInput,
} from '../hooks/useEffects';
import { NumberField, SelectField } from './FormFields';
import { Button } from './Button';

const EFFECT_TYPE_OPTS = [
  { value: 'damage', label: 'Damage' },
  { value: 'heal', label: 'Heal' },
  { value: 'buff', label: 'Buff' },
  { value: 'debuff', label: 'Debuff' },
  { value: 'dot', label: 'DoT' },
  { value: 'hot', label: 'HoT' },
  { value: 'stun', label: 'Stun' },
  { value: 'accuracy_boost', label: 'Accuracy Boost' },
  { value: 'dodge_all', label: 'Dodge All' },
];

const DAMAGE_SUBTYPE_OPTS = [
  { value: '', label: '— None —' },
  { value: 'slashing', label: 'Slashing' },
  { value: 'piercing', label: 'Piercing' },
  { value: 'bludgeoning', label: 'Bludgeoning' },
  { value: 'fire', label: 'Fire' },
  { value: 'cold', label: 'Cold' },
  { value: 'lightning', label: 'Lightning' },
  { value: 'poison', label: 'Poison' },
  { value: 'psychic', label: 'Psychic' },
];

const TARGET_OPTS = [
  { value: 'enemy', label: 'Enemy' },
  { value: 'self', label: 'Self' },
  { value: 'ally', label: 'Ally' },
  { value: 'area', label: 'Area' },
  { value: 'random_enemy', label: 'Random Enemy' },
];

const SCALING_STAT_OPTS = [
  { value: '', label: '— None —' },
  { value: 'strength', label: 'Strength' },
  { value: 'dexterity', label: 'Dexterity' },
  { value: 'constitution', label: 'Constitution' },
  { value: 'intelligence', label: 'Intelligence' },
  { value: 'wisdom', label: 'Wisdom' },
];

const EMPTY_EFFECT: EffectInput = {
  effect_type: 'damage',
  damage_subtype: '',
  target: 'enemy',
  value: 0,
  duration: 0,
  scaling_stat: '',
  scaling_ratio: 0,
  sort_order: 0,
};

function EffectRow({
  effect,
  onDelete,
  isDeleting,
}: Readonly<{
  effect: AbilityEffect
  onDelete: () => void
  isDeleting: boolean
}>) {
  const typeLabel = EFFECT_TYPE_OPTS.find((o) => o.value === effect.effect_type)?.label ?? effect.effect_type;
  const targetLabel = TARGET_OPTS.find((o) => o.value === effect.target)?.label ?? effect.target;
  const statLabel = effect.scaling_stat
    ? (SCALING_STAT_OPTS.find((o) => o.value === effect.scaling_stat)?.label ?? effect.scaling_stat)
    : null;

  return (
    <div className="flex items-center gap-2 p-2 bg-card-alt rounded text-sm">
      <span className="talent-effect">{typeLabel}</span>
      {effect.damage_subtype && (
        <span className="talent-effect-value">{effect.damage_subtype}</span>
      )}
      <span className="text-muted">→ {targetLabel}</span>
      {effect.value > 0 && <span className="talent-effect-value">{effect.value}</span>}
      {effect.duration > 0 && <span className="text-muted">{effect.duration}t</span>}
      {statLabel && <span className="text-muted">scales {statLabel}</span>}
      {effect.scaling_ratio > 0 && <span className="text-muted">×{effect.scaling_ratio}</span>}
      <span className="text-muted ml-auto">#{effect.sort_order}</span>
      <Button variant="danger" size="sm" onClick={onDelete} disabled={isDeleting} type="button">
        {isDeleting ? '...' : '×'}
      </Button>
    </div>
  );
}

function NewEffectForm({
  abilityId,
  sortOrder,
}: Readonly<{
  abilityId: number
  sortOrder: number
}>) {
  const createEffect = useCreateEffect();
  const [form, setForm] = useState<EffectInput>({ ...EMPTY_EFFECT, sort_order: sortOrder });
  const set = (patch: Partial<EffectInput>) => setForm((prev) => ({ ...prev, ...patch }));

  const handleAdd = () => {
    createEffect.mutate({ abilityId, input: form }, { onSuccess: () => setForm({ ...EMPTY_EFFECT, sort_order: sortOrder + 1 }) });
  };

  return (
    <div className="space-y-2 p-2 border border-dashed rounded">
      <div className="grid grid-cols-3 gap-2">
        <SelectField label="Type" value={form.effect_type} onChange={(v) => set({ effect_type: v })} options={EFFECT_TYPE_OPTS} />
        <SelectField label="Target" value={form.target ?? 'enemy'} onChange={(v) => set({ target: v })} options={TARGET_OPTS} />
        <SelectField label="Dmg Subtype" value={form.damage_subtype ?? ''} onChange={(v) => set({ damage_subtype: v })} options={DAMAGE_SUBTYPE_OPTS} />
      </div>
      <div className="grid grid-cols-3 gap-2">
        <NumberField label="Value" value={form.value ?? 0} onChange={(v) => set({ value: v })} />
        <NumberField label="Duration" value={form.duration ?? 0} onChange={(v) => set({ duration: v })} />
        <NumberField label="Sort Order" value={form.sort_order ?? 0} onChange={(v) => set({ sort_order: v })} />
      </div>
      <div className="grid grid-cols-2 gap-2">
        <SelectField label="Scaling Stat" value={form.scaling_stat ?? ''} onChange={(v) => set({ scaling_stat: v })} options={SCALING_STAT_OPTS} />
        <NumberField label="Scaling Ratio" value={form.scaling_ratio ?? 0} onChange={(v) => set({ scaling_ratio: v })} step={0.1} />
      </div>
      <Button variant="primary" size="sm" fullWidth disabled={createEffect.isPending} onClick={handleAdd} type="button">
        {createEffect.isPending ? 'Adding...' : '+ Add Effect'}
      </Button>
    </div>
  );
}

export function EffectsSubForm({ abilityId }: Readonly<{ abilityId: number }>) {
  const { data: effects, isLoading } = useEffects(abilityId);
  const deleteEffect = useDeleteEffect();
  const [showAdd, setShowAdd] = useState(false);

  const sorted = [...(effects ?? [])].sort((a, b) => a.sort_order - b.sort_order);
  const nextOrder = sorted.length > 0 ? sorted[sorted.length - 1].sort_order + 1 : 0;

  return (
    <div className="space-y-2">
      <div className="flex items-center justify-between">
        <h4 className="text-sm font-semibold text-text m-0">Effects</h4>
        <Button variant="ghost" size="sm" onClick={() => setShowAdd(!showAdd)} type="button">
          {showAdd ? 'Cancel' : '+ Add'}
        </Button>
      </div>

      {isLoading && <div className="text-muted text-sm">Loading effects...</div>}

      {sorted.length === 0 && !isLoading && (
        <p className="text-muted text-sm m-0">No effects. Add one to define what this ability does.</p>
      )}

      {sorted.map((effect) => (
        <EffectRow
          key={effect.id}
          effect={effect}
          onDelete={() => deleteEffect.mutate({ id: effect.id, abilityId })}
          isDeleting={deleteEffect.isPending}
        />
      ))}

      {showAdd && <NewEffectForm abilityId={abilityId} sortOrder={nextOrder} />}
    </div>
  );
}