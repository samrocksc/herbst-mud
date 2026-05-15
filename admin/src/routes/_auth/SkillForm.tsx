import { useState } from 'react';
import {
  useCreateCompetencyCategory,
  useUpdateCompetencyCategory,
  type CompetencyCategory,
  type CompetencyThresholdInput,
} from '../../hooks/useCompetencies';
import { Button } from '../../components/Button';
import { FormField, NumberField } from '../../components/FormFields';
import { FormError } from '../../components/fields/FormError';

const DEFAULT_THRESHOLDS: CompetencyThresholdInput[] = [
  { level: 1, xp_required: 0, damage_multiplier: 1.0, defense_multiplier: 1.0 },
  { level: 2, xp_required: 100, damage_multiplier: 1.05, defense_multiplier: 1.05 },
  { level: 3, xp_required: 300, damage_multiplier: 1.1, defense_multiplier: 1.1 },
  { level: 4, xp_required: 600, damage_multiplier: 1.15, defense_multiplier: 1.15 },
  { level: 5, xp_required: 1000, damage_multiplier: 1.2, defense_multiplier: 1.2 },
  { level: 6, xp_required: 1500, damage_multiplier: 1.3, defense_multiplier: 1.25 },
  { level: 7, xp_required: 2200, damage_multiplier: 1.4, defense_multiplier: 1.3 },
  { level: 8, xp_required: 3000, damage_multiplier: 1.5, defense_multiplier: 1.35 },
  { level: 9, xp_required: 4000, damage_multiplier: 1.65, defense_multiplier: 1.4 },
  { level: 10, xp_required: 5500, damage_multiplier: 1.8, defense_multiplier: 1.5 },
];

type Props = Readonly<{
  category?: CompetencyCategory | null
  onSubmit: () => void
  onCancel: () => void
}>

export function SkillForm({ category, onSubmit, onCancel }: Props) {
  const createMutation = useCreateCompetencyCategory();
  const updateMutation = useUpdateCompetencyCategory();
  const isEditing = !!category;
  const isLoading = createMutation.isPending || updateMutation.isPending;
  const error = createMutation.error || updateMutation.error;

  const [id, setId] = useState(category?.id ?? '');
  const [name, setName] = useState(category?.name ?? '');
  const [xpMult, setXpMult] = useState(category?.xp_multiplier ?? 0.2);
  const [thresholds, setThresholds] = useState<CompetencyThresholdInput[]>(
    category?.thresholds?.length ? category.thresholds : DEFAULT_THRESHOLDS,
  );

  const updateThreshold = (level: number, field: keyof CompetencyThresholdInput, value: number) => {
    setThresholds(prev => prev.map(t => t.level === level ? { ...t, [field]: value } : t));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      if (isEditing) {
        await updateMutation.mutateAsync({ id: category!.id, input: { name, xp_multiplier: xpMult, thresholds } });
      } else {
        await createMutation.mutateAsync({ id, name, xp_multiplier: xpMult, thresholds });
      }
      onSubmit();
    } catch { /* error is in mutation state */ }
  };

  return (
    <div className="form-card space-y-3">
      <h3 className="mt-0 mb-0 text-text text-base font-semibold">
        {isEditing ? 'Edit Skill' : 'Add New Skill'}
      </h3>
      {error && <FormError message={error.message} />}
      <form onSubmit={handleSubmit} className="space-y-3">
        {!isEditing && (
          <FormField label="ID (slug)" value={id} onChange={setId} tooltip="Unique identifier, e.g. 'blades', 'staves'" />
        )}
        <FormField label="Name" value={name} onChange={setName} tooltip="Display name, e.g. 'Blades', 'Staves'" />
        <NumberField label="XP Multiplier" value={xpMult} onChange={setXpMult} tooltip="Multiplier applied to raw XP before storing" />

        <div>
          <h4 className="text-sm font-semibold text-text mb-1">Level Thresholds</h4>
          <div className="overflow-x-auto">
            <table className="w-full text-xs">
              <thead>
                <tr className="text-text-muted">
                  <th className="text-left py-1 px-1">Lvl</th>
                  <th className="text-right py-1 px-1">XP Req</th>
                  <th className="text-right py-1 px-1">Dmg Mult</th>
                  <th className="text-right py-1 px-1">Def Mult</th>
                </tr>
              </thead>
              <tbody>
                {thresholds.map(t => (
                  <tr key={t.level} className="border-t border-border">
                    <td className="py-1 px-1 text-text-muted">{t.level}</td>
                    <td className="py-1 px-1">
                      <input type="number" className="w-20 input text-right" value={t.xp_required}
                        onChange={e => updateThreshold(t.level, 'xp_required', Number(e.target.value))} />
                    </td>
                    <td className="py-1 px-1">
                      <input type="number" step="0.01" className="w-20 input text-right" value={t.damage_multiplier}
                        onChange={e => updateThreshold(t.level, 'damage_multiplier', Number(e.target.value))} />
                    </td>
                    <td className="py-1 px-1">
                      <input type="number" step="0.01" className="w-20 input text-right" value={t.defense_multiplier}
                        onChange={e => updateThreshold(t.level, 'defense_multiplier', Number(e.target.value))} />
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>

        <div className="flex gap-2 pt-1">
          <Button type="submit" variant="primary" disabled={isLoading} fullWidth>
            {isLoading ? 'Saving...' : isEditing ? 'Update Skill' : 'Create Skill'}
          </Button>
          <Button variant="secondary" onClick={onCancel} fullWidth>Cancel</Button>
        </div>
      </form>
    </div>
  );
}