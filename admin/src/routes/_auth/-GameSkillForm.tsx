/* eslint-disable functional/no-mixed-types, functional/prefer-immutable-types */
import { useState } from "react";
import {
  useCreateGameSkill,
  useUpdateGameSkill,
  type GameSkill,
  type GameSkillInput,
} from "../../hooks/useGameSkills";
import { Button } from "../../components/Button";
import { FormField, NumberField, TextareaField, SelectField } from "../../components/FormFields";
import { FormError } from "../../components/fields/FormError";
import { showToast } from "../../components/Toast";

const CATEGORY_OPTS = [
  { value: "weapon", label: "Weapon" },
  { value: "armor", label: "Armor" },
  { value: "craft", label: "Craft" },
  { value: "magic", label: "Magic" },
];

const XP_CURVE_MODE_OPTS = [
  { value: "percentage", label: "Percentage" },
  { value: "hand_coded", label: "Hand-Coded" },
];

type Props = Readonly<{
  skill?: GameSkill | null
  onSubmit: () => void
  onCancel: () => void
}>

export function GameSkillForm({ skill, onSubmit, onCancel }: Props) {
  const createMutation = useCreateGameSkill();
  const updateMutation = useUpdateGameSkill();
  const isEditing = !!skill;
  const isLoading = createMutation.isPending || updateMutation.isPending;
  const error = createMutation.error || updateMutation.error;

  const [name, setName] = useState(skill?.name ?? "");
  const [displayName, setDisplayName] = useState(skill?.display_name ?? "");
  const [description, setDescription] = useState(skill?.description ?? "");
  const [category, setCategory] = useState(skill?.category ?? "weapon");
  const [maxLevel, setMaxLevel] = useState(skill?.max_level ?? 100);
  const [xpCurveMode, setXpCurveMode] = useState(skill?.xp_curve_mode ?? "percentage");

  // xp_curve_data handling: percentage → single number, hand_coded → JSON string
  const initialPercentage =
    skill?.xp_curve_data && typeof skill.xp_curve_data.percentage === "number"
      ? skill.xp_curve_data.percentage
      : 50;
  const initialThresholdsJson =
    skill?.xp_curve_data && Array.isArray(skill.xp_curve_data.thresholds)
      ? JSON.stringify(skill.xp_curve_data.thresholds, null, 2)
      : "[0, 100, 300, 600, 1000, 1500, 2200, 3000, 4000, 5500]";

  const [percentage, setPercentage] = useState(initialPercentage);
  const [thresholdsJson, setThresholdsJson] = useState(initialThresholdsJson);
  const [jsonError, setJsonError] = useState("");

  function buildXpCurveData(): Record<string, unknown> {
    if (xpCurveMode === "percentage") {
      return { percentage };
    }
    try {
      const parsed = JSON.parse(thresholdsJson);
      if (!Array.isArray(parsed)) {
        setJsonError("Thresholds must be a JSON array of numbers.");
        return { thresholds: [] };
      }
      setJsonError("");
      return { thresholds: parsed };
    } catch {
      setJsonError("Invalid JSON for thresholds.");
      return { thresholds: [] };
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setJsonError("");
    const xpCurveData = buildXpCurveData();
    if (jsonError) return;

    const input: GameSkillInput = {
      name,
      display_name: displayName,
      description,
      category,
      max_level: maxLevel,
      xp_curve_mode: xpCurveMode,
      xp_curve_data: xpCurveData,
    };

    try {
      if (isEditing) {
        await updateMutation.mutateAsync({ id: skill!.id, input });
        showToast("Skill updated", "success");
      } else {
        await createMutation.mutateAsync(input);
        showToast("Skill created", "success");
      }
      onSubmit();
    } catch (err) {
      console.error("GameSkill save error:", err);
      const message = err instanceof Error ? err.message : "Failed to save skill";
      showToast(message, "error");
    }
  };

  return (
    <div className="form-card space-y-3">
      <h3 className="mt-0 mb-0 text-text text-base font-semibold">
        {isEditing ? "Edit Skill" : "Add New Skill"}
      </h3>
      {error && <FormError message={error.message} />}
      {jsonError && <FormError message={jsonError} />}
      <form onSubmit={handleSubmit} className="space-y-3">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <FormField
            label="Name (machine)"
            value={name}
            onChange={setName}
            required
            tooltip="Machine identifier, e.g. 'blades', 'heavy_armor'"
            placeholder="e.g., blades"
          />
          <FormField
            label="Display Name"
            value={displayName}
            onChange={setDisplayName}
            required
            tooltip="Human-readable name shown to players, e.g. 'Blades', 'Heavy Armor'"
            placeholder="e.g., Blades"
          />
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <SelectField
            label="Category"
            value={category}
            onChange={setCategory}
            options={CATEGORY_OPTS}
            tooltip="Weapon = combat skills, Armor = defense skills, Craft = crafting skills, Magic = spellcasting skills"
          />
          <NumberField
            label="Max Level"
            value={maxLevel}
            onChange={setMaxLevel}
            min={1}
            tooltip="Maximum skill level. Default is 100."
          />
        </div>

        <TextareaField
          label="Description"
          value={description}
          onChange={setDescription}
          rows={3}
          placeholder="What this skill represents and how it's used."
        />

        <div>
          <h4 className="text-sm font-semibold text-text mb-1">XP Curve</h4>
          <p className="text-xs text-text-muted mb-2">
            Configure how XP maps to skill levels. "Percentage" uses a flat percentage of current
            level as the threshold. "Hand-Coded" uses explicit XP thresholds per level.
          </p>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <SelectField
              label="XP Curve Mode"
              value={xpCurveMode}
              onChange={setXpCurveMode}
              options={XP_CURVE_MODE_OPTS}
              tooltip="Percentage = each level requires a percentage of the previous level's XP. Hand-Coded = explicit threshold array."
            />
            {xpCurveMode === "percentage" ? (
              <NumberField
                label="Percentage"
                value={percentage}
                onChange={setPercentage}
                min={1}
                max={1000}
                tooltip="Percentage of current-level XP needed to advance. E.g. 50 means each level needs 50% more XP than the last."
              />
            ) : (
              <div>
                <label className="text-text-muted text-xs block mb-1">
                  Thresholds (JSON array)
                </label>
                <textarea
                  className="w-full p-2 bg-surface border border-border rounded text-text text-sm font-mono resize-y"
                  rows={5}
                  value={thresholdsJson}
                  onChange={(e) => setThresholdsJson(e.target.value)}
                  placeholder="[0, 100, 300, 600, 1000, ...]"
                />
              </div>
            )}
          </div>
        </div>

        <div className="flex gap-2 pt-1">
          <Button type="submit" variant="primary" disabled={isLoading} fullWidth>
            {isLoading ? "Saving..." : isEditing ? "Update Skill" : "Create Skill"}
          </Button>
          <Button variant="secondary" onClick={onCancel} fullWidth>
            Cancel
          </Button>
        </div>
      </form>
    </div>
  );
}