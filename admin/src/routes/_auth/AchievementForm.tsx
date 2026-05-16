import { useState } from "react";
import type { Achievement, AchievementInput } from "../../hooks/useAchievements";
import { Button } from "../../components/Button";
import { FormField } from "../../components/FormFields";
import { NumberField } from "../../components/FormFields";
import { TextareaField } from "../../components/FormFields";
import { FormError } from "../../components/FormFields";

const EMPTY_INPUT: AchievementInput = {
  name: "",
  description: "",
  icon: "",
  xp_reward: 0,
  criteria: "",
};

export function AchievementForm({
  achievement,
  onSubmit,
  onCancel,
  isLoading,
  error,
}: Readonly<{
  achievement: Achievement | null
  onSubmit: (data: AchievementInput) => void
  onCancel: () => void
  isLoading: boolean
  error?: string
}>) {
  const [formData, setFormData] = useState<AchievementInput>(() =>
    achievement
      ? { name: achievement.name, description: achievement.description, icon: achievement.icon, xp_reward: achievement.xp_reward, criteria: achievement.criteria }
      : EMPTY_INPUT,
  );
  const set = <K extends keyof AchievementInput>(key: K, value: AchievementInput[K]) =>
    setFormData(prev => ({ ...prev, [key]: value }));

  return (
    <div className="form-card space-y-3">
      <h3 className="mt-0 mb-0 text-text text-base font-semibold">
        {achievement ? "Edit Achievement" : "Add New Achievement"}
      </h3>
      {error && <FormError message={error} />}
      <form onSubmit={(e) => { e.preventDefault(); onSubmit(formData); }} className="space-y-3">
        <FormField label="Name" value={formData.name} onChange={(v) => set("name", v)} required />
        <TextareaField label="Description" value={formData.description} onChange={(v) => set("description", v)} rows={3} />
        <FormField label="Icon (emoji)" value={formData.icon} onChange={(v) => set("icon", v)} placeholder="e.g. 🏆" />
        <NumberField label="XP Reward" value={formData.xp_reward} onChange={(v) => set("xp_reward", v)} min={0} />
        <TextareaField label="Criteria (JSON)" value={formData.criteria} onChange={(v) => set("criteria", v)} rows={3} placeholder='e.g. {"type":"kill_count","target":10}' />
        <div className="flex gap-2 pt-1">
          <Button type="submit" variant="primary" disabled={isLoading} fullWidth>
            {isLoading ? "Saving..." : achievement ? "Update Achievement" : "Create Achievement"}
          </Button>
          <Button variant="secondary" onClick={onCancel} fullWidth>Cancel</Button>
        </div>
      </form>
    </div>
  );
}

export function DeleteConfirmation({
  achievement,
  onConfirm,
  onCancel,
  isLoading,
}: Readonly<{
  achievement: Achievement
  onConfirm: () => void
  onCancel: () => void
  isLoading: boolean
}>) {
  return (
    <div className="modal-overlay" onClick={onCancel}>
      <div className="modal-content modal-sm" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h3>Delete Achievement</h3>
          <Button variant="ghost" size="sm" onClick={onCancel} aria-label="Close">×</Button>
        </div>
        <div className="modal-body">
          <p>Are you sure you want to delete <strong>{achievement.name}</strong>?</p>
          <p className="text-muted">This action cannot be undone.</p>
        </div>
        <div className="modal-footer">
          <Button variant="danger" onClick={onConfirm} disabled={isLoading}>{isLoading ? "Deleting..." : "Delete"}</Button>
          <Button variant="secondary" onClick={onCancel}>Cancel</Button>
        </div>
      </div>
    </div>
  );
}