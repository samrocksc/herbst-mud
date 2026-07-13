/* eslint-disable functional/no-mixed-types, functional/prefer-immutable-types, functional/immutable-data */
import { createFileRoute } from "@tanstack/react-router";
import { useState, useEffect } from "react";
import { useWorldConfig, useUpdateWorldConfig, DEFAULT_CONFIG, type WorldConfig } from "../../hooks/useWorldConfig";
import { PageHeader } from "../../components/PageHeader";
import { PageContainer } from "../../components/PageContainer";
import { Button } from "../../components/Button";
import { NumberField, SelectField, CheckboxField, FormError } from "../../components/fields";
import { showToast } from "../../components/Toast";

export const Route = createFileRoute("/_auth/world-config")({
  component: WorldConfigPage,
});

const LEVEL_CURVE_MODE_OPTS = [
  { value: "percentage", label: "Percentage" },
  { value: "hand_coded", label: "Hand-Coded" },
];

function WorldConfigPage() {
  const { data: config, isLoading, error } = useWorldConfig();
  const updateMutation = useUpdateWorldConfig();

  type MutableConfig = {
    level_curve: { mode: "percentage" | "hand_coded"; base_xp: number; percentage: number; max_level: number; thresholds?: number[] };
    stat_growth: { hp_per_level: number; mana_per_level: number; stamina_per_level: number };
    skill_xp: { usage_diminishing_returns: boolean; usage_cap_per_hour: number; anti_grind_kill_threshold: number };
    reclass: { allowed: boolean; cost: number; min_level: number; cooldown_seconds: number; skill_retention: number };
    rerace: { allowed: boolean; cost: number };
  };

  const [form, setForm] = useState<MutableConfig>(DEFAULT_CONFIG);
  const [thresholdsJson, setThresholdsJson] = useState<string>("[0, 100, 300, 600, 1000, 1500, 2200, 3000, 4000, 5500]");
  const [jsonError, setJsonError] = useState("");
  const [saveError, setSaveError] = useState("");

  useEffect(() => {
    if (config) {
      // Merge with defaults to ensure all sections exist (API may not return reclass/rerace)
      const merged = JSON.parse(JSON.stringify({
        ...DEFAULT_CONFIG,
        ...config,
      })) as MutableConfig;
      setForm(merged);
      if (config.level_curve?.thresholds) {
        setThresholdsJson(JSON.stringify(config.level_curve.thresholds, null, 2));
      }
    }
  }, [config]);

  const setLevelCurve = <K extends keyof MutableConfig["level_curve"]>(key: K, value: MutableConfig["level_curve"][K]) =>
    setForm(prev => ({ ...prev, level_curve: { ...prev.level_curve, [key]: value } }));

  const setStatGrowth = <K extends keyof MutableConfig["stat_growth"]>(key: K, value: MutableConfig["stat_growth"][K]) =>
    setForm(prev => ({ ...prev, stat_growth: { ...prev.stat_growth, [key]: value } }));

  const setSkillXp = <K extends keyof MutableConfig["skill_xp"]>(key: K, value: MutableConfig["skill_xp"][K]) =>
    setForm(prev => ({ ...prev, skill_xp: { ...prev.skill_xp, [key]: value } }));

  const setReclass = <K extends keyof MutableConfig["reclass"]>(key: K, value: MutableConfig["reclass"][K]) =>
    setForm(prev => ({ ...prev, reclass: { ...prev.reclass, [key]: value } }));

  const setRerace = <K extends keyof MutableConfig["rerace"]>(key: K, value: MutableConfig["rerace"][K]) =>
    setForm(prev => ({ ...prev, rerace: { ...prev.rerace, [key]: value } }));

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSaveError("");
    setJsonError("");

    const payload: MutableConfig = { ...form };

    if (form.level_curve.mode === "hand_coded") {
      try {
        const parsed = JSON.parse(thresholdsJson);
        if (!Array.isArray(parsed)) {
          setJsonError("Thresholds must be a JSON array of numbers.");
          return;
        }
        payload.level_curve = { ...form.level_curve, thresholds: parsed };
      } catch {
        setJsonError("Invalid JSON for thresholds.");
        return;
      }
    } else {
      payload.level_curve = { ...form.level_curve, thresholds: undefined };
    }

    try {
      await updateMutation.mutateAsync(payload);
      showToast("World config saved", "success");
    } catch (err) {
      const msg = err instanceof Error ? err.message : "Failed to save world config";
      setSaveError(msg);
      showToast(msg, "error");
    }
  };

  if (isLoading) {
    return (
      <PageContainer>
        <PageHeader title="World Config" backTo="/dashboard" />
        <div className="loading">Loading world config...</div>
      </PageContainer>
    );
  }

  if (error) {
    return (
      <PageContainer>
        <PageHeader title="World Config" backTo="/dashboard" />
        <div className="error-banner">{error instanceof Error ? error.message : "Failed to load world config"}</div>
      </PageContainer>
    );
  }

  const isSaving = updateMutation.isPending;

  return (
    <PageContainer>
      <PageHeader title="World Config" backTo="/dashboard" />
      <p className="text-sm text-text-muted mb-4">
        Configure the XP system for this world: level curve, stat growth per level, skill XP anti-grind settings,
        and reclass/rerace rules.
      </p>

      {saveError && <FormError message={saveError} />}
      {jsonError && <FormError message={jsonError} />}

      <form onSubmit={handleSubmit} className="space-y-6 max-w-[800px]">
        {/* Level Curve Section */}
        <section className="form-card space-y-3">
          <h3 className="mt-0 mb-1 text-text text-base font-semibold">Level Curve</h3>
          <p className="text-xs text-text-muted mb-2">
            Configure how character XP maps to levels. "Percentage" uses a base XP and a percentage increase per level.
            "Hand-Coded" uses an explicit array of XP thresholds.
          </p>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <SelectField
              label="Mode"
              value={form.level_curve.mode}
              onChange={(v) => setLevelCurve("mode", v as "percentage" | "hand_coded")}
              options={LEVEL_CURVE_MODE_OPTS}
              tooltip="Percentage = each level requires a percentage more XP. Hand-Coded = explicit threshold array."
            />
            <NumberField
              label="Max Level"
              value={form.level_curve.max_level}
              onChange={(v) => setLevelCurve("max_level", v)}
              min={1}
              tooltip="Maximum character level."
            />
          </div>
          {form.level_curve.mode === "percentage" ? (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <NumberField
                label="Base XP"
                value={form.level_curve.base_xp}
                onChange={(v) => setLevelCurve("base_xp", v)}
                min={1}
                tooltip="XP needed for level 1 → 2."
              />
              <NumberField
                label="Percentage (%)"
                value={form.level_curve.percentage}
                onChange={(v) => setLevelCurve("percentage", v)}
                min={1}
                max={1000}
                tooltip="Each level requires this % more XP than the previous."
              />
            </div>
          ) : (
            <div>
              <label className="text-text-muted text-xs block mb-1">
                Thresholds (JSON array of numbers)
              </label>
              <textarea
                className="w-full p-2 bg-surface border border-border rounded text-text text-sm font-mono resize-y"
                rows={6}
                value={thresholdsJson}
                onChange={(e) => setThresholdsJson(e.target.value)}
                placeholder="[0, 100, 300, 600, 1000, ...]"
              />
              <p className="text-xs text-text-muted mt-1">
                Each entry is the total XP needed to reach that level (index 0 = level 1).
              </p>
            </div>
          )}
        </section>

        {/* Stat Growth Section */}
        <section className="form-card space-y-3">
          <h3 className="mt-0 mb-1 text-text text-base font-semibold">Stat Growth</h3>
          <p className="text-xs text-text-muted mb-2">
            Base stat gains per level. Racial multipliers are applied on top of these values.
          </p>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <NumberField
              label="HP per Level"
              value={form.stat_growth.hp_per_level}
              onChange={(v) => setStatGrowth("hp_per_level", v)}
              min={0}
              tooltip="Base HP gained per character level."
            />
            <NumberField
              label="Mana per Level"
              value={form.stat_growth.mana_per_level}
              onChange={(v) => setStatGrowth("mana_per_level", v)}
              min={0}
              tooltip="Base mana gained per character level."
            />
            <NumberField
              label="Stamina per Level"
              value={form.stat_growth.stamina_per_level}
              onChange={(v) => setStatGrowth("stamina_per_level", v)}
              min={0}
              tooltip="Base stamina gained per character level."
            />
          </div>
        </section>

        {/* Skill XP Section */}
        <section className="form-card space-y-3">
          <h3 className="mt-0 mb-1 text-text text-base font-semibold">Skill XP</h3>
          <p className="text-xs text-text-muted mb-2">
            Anti-grind controls for skill usage-based XP. These prevent macro/automation abuse.
          </p>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <NumberField
              label="Usage Cap per Hour"
              value={form.skill_xp.usage_cap_per_hour}
              onChange={(v) => setSkillXp("usage_cap_per_hour", v)}
              min={0}
              tooltip="Maximum skill usage XP events counted per hour per skill."
            />
            <NumberField
              label="Anti-Grind Kill Threshold"
              value={form.skill_xp.anti_grind_kill_threshold}
              onChange={(v) => setSkillXp("anti_grind_kill_threshold", v)}
              min={0}
              tooltip="Max kills of the same NPC type that grant skill XP before diminishing returns kick in."
            />
          </div>
          <CheckboxField
            label="Usage Diminishing Returns"
            checked={form.skill_xp.usage_diminishing_returns}
            onChange={(v) => setSkillXp("usage_diminishing_returns", v)}
          />
        </section>

        {/* Reclass Section */}
        <section className="form-card space-y-3">
          <h3 className="mt-0 mb-1 text-text text-base font-semibold">Reclass</h3>
          <p className="text-xs text-text-muted mb-2">
            Rules for changing character class within the same race.
          </p>
          <CheckboxField
            label="Allow Reclass"
            checked={form.reclass.allowed}
            onChange={(v) => setReclass("allowed", v)}
          />
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <NumberField
              label="Cost (gold)"
              value={form.reclass.cost}
              onChange={(v) => setReclass("cost", v)}
              min={0}
              tooltip="Gold cost to reclass."
            />
            <NumberField
              label="Min Level"
              value={form.reclass.min_level}
              onChange={(v) => setReclass("min_level", v)}
              min={1}
              tooltip="Minimum character level to reclass."
            />
          </div>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <NumberField
              label="Cooldown (seconds)"
              value={form.reclass.cooldown_seconds}
              onChange={(v) => setReclass("cooldown_seconds", v)}
              min={0}
              tooltip="Cooldown between reclass attempts."
            />
            <div>
              <FieldLabelHack label="Skill Retention (0–1)" tooltip="Fraction of skill levels retained when reclassing. 1.0 = keep all, 0.0 = lose all." />
              <input
                type="text"
                inputMode="decimal"
                value={String(form.reclass.skill_retention)}
                onChange={(e) => {
                  const raw = e.target.value;
                  if (raw === "" || raw === "-") { setReclass("skill_retention", 0); return; }
                  const parsed = parseFloat(raw);
                  if (!isNaN(parsed)) setReclass("skill_retention", parsed);
                }}
                className="w-full px-3 py-2 bg-surface border border-border rounded text-sm text-text focus:outline-none focus:border-primary"
              />
            </div>
          </div>
        </section>

        {/* Rerace Section */}
        <section className="form-card space-y-3">
          <h3 className="mt-0 mb-1 text-text text-base font-semibold">Rerace</h3>
          <p className="text-xs text-text-muted mb-2">
            Rules for changing character race.
          </p>
          <CheckboxField
            label="Allow Rerace"
            checked={form.rerace.allowed}
            onChange={(v) => setRerace("allowed", v)}
          />
          <NumberField
            label="Cost (gold)"
            value={form.rerace.cost}
            onChange={(v) => setRerace("cost", v)}
            min={0}
            tooltip="Gold cost to rerace."
          />
        </section>

        {/* Save */}
        <div className="flex gap-2 pt-1">
          <Button type="submit" variant="primary" disabled={isSaving} fullWidth>
            {isSaving ? "Saving..." : "Save World Config"}
          </Button>
        </div>
      </form>
    </PageContainer>
  );
}

/**
 * Small helper to render a field label with tooltip, matching FieldLabel styling.
 * Used for the float skill_retention input that NumberField doesn't support (it parses int).
 */
function FieldLabelHack({ label, tooltip }: Readonly<{ label: string; tooltip?: string }>) {
  return (
    <label className="text-text-muted text-xs block mb-1" title={tooltip}>
      {label}
    </label>
  );
}