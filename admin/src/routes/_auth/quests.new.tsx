import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import {
  useCreateQuest,
  useQuestLookups,
  type QuestInput,
  EMPTY_REWARDS,
  type QuestObjective,
} from "../../hooks/useQuests";
import { PageHeader } from "../../components/PageHeader";
import { Button } from "../../components/Button";
import { FormField, TextareaField, NumberField, SelectField } from "../../components/FormFields";
import { showToast } from "../../components/Toast";
import { PageContainer } from "../../components/PageContainer";
import { SearchableSelect } from "../../components/SearchableSelect";
import { ResourceMultiSelect } from "../../components/ResourceMultiSelect";
import { RESOURCE_ENDPOINTS } from "../../utils/resourceEndpoints";
import { useTags } from "../../hooks/useTags";

export const Route = createFileRoute("/_auth/quests/new")({
  component: CreateQuestPage,
});

const REPEAT_MODE_OPTS = [
  { value: "none", label: "None (one-time)" },
  { value: "cooldown", label: "Cooldown" },
  { value: "always", label: "Always repeatable" },
];

const EMPTY_OBJECTIVE: QuestObjective = {
  type: "", target_id: "", tag_filter: "", count: 1, labels: [], hint: "",
};

const EMPTY_QUEST: QuestInput = {
  name: "",
  description: "",
  prerequisite_quest_ids: [],
  objectives: [{ ...EMPTY_OBJECTIVE }],
  rewards: { ...EMPTY_REWARDS },
  repeat_mode: "none",
  cooldown_hours: 0,
  is_active: true,
};

function CreateQuestPage() {
  const navigate = useNavigate();
  const createQuest = useCreateQuest();
  const { data: lookups, isLoading: lookupsLoading } = useQuestLookups();
  const { data: tags } = useTags();
  const [formData, setFormData] = useState<QuestInput>(EMPTY_QUEST);
  const set = (patch: Partial<QuestInput>) => setFormData((prev) => ({ ...prev, ...patch }));

  if (lookupsLoading) return <div className="loading">Loading options...</div>;

  const addObjective = () => {
    const objs = [...(formData.objectives ?? []), { ...EMPTY_OBJECTIVE }];
    set({ objectives: objs });
  };
  const updateObjective = (i: number, patch: Partial<QuestObjective>) => {
    const objs = formData.objectives?.map((o, idx) => idx === i ? { ...o, ...patch } : o) ?? [];
    set({ objectives: objs });
  };
  const removeObjective = (i: number) => {
    const objs = formData.objectives?.filter((_, idx) => idx !== i) ?? [];
    set({ objectives: objs });
  };

  const getTargetsForType = (type: string) => {
    if (!lookups) return [];
    switch (type) {
      case "kill": return lookups.npcs;
      case "explore": return lookups.rooms;
      case "collect": return lookups.items;
      default: return [];
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await createQuest.mutateAsync(formData);
      showToast("Quest created", "success");
      navigate({ to: "/quests" });
    } catch { /* toasted globally */ }
  };

  return (
    <PageContainer>
      <PageHeader title="Create Quest" showBack backTo="/quests" />
      <div className="bg-surface p-6 border border-border rounded">
        <form onSubmit={handleSubmit} className="space-y-3">
          <FormField label="Name" value={formData.name ?? ""} onChange={(v) => set({ name: v })} />
          <TextareaField label="Description" value={formData.description ?? ""} onChange={(v) => set({ description: v })} rows={3} />
          <SelectField label="Repeat Mode" value={formData.repeat_mode ?? "none"} onChange={(v) => set({ repeat_mode: v })} options={REPEAT_MODE_OPTS} />
          {(formData.repeat_mode === "cooldown") && (
            <NumberField label="Cooldown (hours)" value={formData.cooldown_hours ?? 0} onChange={(v) => set({ cooldown_hours: v })} />
          )}
          <div className="flex items-center gap-2">
            <input type="checkbox" checked={formData.is_active ?? true} onChange={(e) => set({ is_active: e.target.checked })} id="quest-active" />
            <label htmlFor="quest-active" className="text-sm text-text">Active</label>
          </div>

          {/* Prerequisite Quests */}
          <ResourceMultiSelect
            label="Prerequisite Quests"
            value={formData.prerequisite_quest_ids ?? []}
            onChange={(ids) => set({ prerequisite_quest_ids: ids as string[] })}
            {...RESOURCE_ENDPOINTS.quests}
          />

          <div className="border-t border-border pt-3 mt-3">
            <div className="flex items-center justify-between mb-2">
              <h4 className="text-sm font-semibold text-text">Objectives</h4>
              <Button variant="ghost" size="sm" onClick={addObjective}>+ Objective</Button>
            </div>
            {(formData.objectives ?? []).map((obj, i) => {
              const targetOptions = getTargetsForType(obj.type);
              const tagOptions = (tags ?? []).map(t => ({ id: t.name, name: t.name }));
              return (
                <div key={i} className="grid grid-cols-7 gap-2 mb-2 items-end">
                  <SelectField
                    label="Type"
                    value={obj.type}
                    onChange={(v) => updateObjective(i, { type: v, target_id: "" })}
                    options={[
                      { value: "kill", label: "Kill NPC" },
                      { value: "explore", label: "Explore Room" },
                      { value: "collect", label: "Collect Item" },
                    ]}
                  />
                  <SearchableSelect
                    label="Target"
                    value={obj.target_id || ""}
                    onChange={(v) => updateObjective(i, { target_id: v })}
                    options={targetOptions.map(t => ({ id: t.id, name: t.name }))}
                    placeholder="Select target..."
                  />
                  <SearchableSelect
                    label="Tag Filter"
                    value={obj.tag_filter || ""}
                    onChange={(v) => updateObjective(i, { tag_filter: v })}
                    options={tagOptions}
                    placeholder="Filter by tag..."
                  />
                  <NumberField label="Count" value={obj.count} onChange={(v) => updateObjective(i, { count: v })} />
                  <FormField label="Label" value={obj.labels?.[0] ?? ""} onChange={(v) => updateObjective(i, { labels: [v] })} placeholder="Kill Rats" />
                  <FormField label="Hint" value={obj.hint} onChange={(v) => updateObjective(i, { hint: v })} placeholder="Optional hint" />
                  <Button variant="danger" size="sm" onClick={() => removeObjective(i)}>×</Button>
                </div>
              );
            })}
          </div>

          <div className="border-t border-border pt-3 mt-3">
            <h4 className="text-sm font-semibold text-text mb-2">Rewards</h4>
            <NumberField label="XP" value={formData.rewards?.xp ?? 0} onChange={(v) => set({ rewards: { ...formData.rewards ?? EMPTY_REWARDS, xp: v } })} />

            <ResourceMultiSelect
              label="Item Rewards"
              value={(formData.rewards?.item_ids ?? []) as string[]}
              onChange={(ids) => set({ rewards: { ...formData.rewards ?? EMPTY_REWARDS, item_ids: ids as string[] } })}
              {...RESOURCE_ENDPOINTS.equipmentTemplates}
            />

            <div className="mt-3">
              <label className="text-sm text-muted mb-1 block">Tags to Add</label>
              <div className="flex flex-wrap gap-2">
                {(lookups?.tags ?? []).map(tag => (
                  <button
                    key={tag.id}
                    type="button"
                    onClick={() => {
                      const current = formData.rewards?.tag_adds ?? [];
                      const newTags = current.includes(tag.id)
                        ? current.filter(t => t !== tag.id)
                        : [...current, tag.id];
                      set({ rewards: { ...formData.rewards ?? EMPTY_REWARDS, tag_adds: newTags } });
                    }}
                    className={`px-2 py-1 text-xs rounded border ${
                      (formData.rewards?.tag_adds ?? []).includes(tag.id)
                        ? "bg-primary/20 border-primary text-text"
                        : "bg-surface border-border text-muted hover:border-primary"
                    }`}
                  >
                    {tag.name}
                  </button>
                ))}
              </div>
            </div>

            <ResourceMultiSelect
              label="Effect Rewards"
              value={formData.rewards?.effect_ids ?? []}
              onChange={(ids) => set({ rewards: { ...formData.rewards ?? EMPTY_REWARDS, effect_ids: ids as number[] } })}
              {...RESOURCE_ENDPOINTS.effectDefs}
            />

            <ResourceMultiSelect
              label="Achievement Rewards"
              value={formData.rewards?.achievement_ids ?? []}
              onChange={(ids) => set({ rewards: { ...formData.rewards ?? EMPTY_REWARDS, achievement_ids: ids as number[] } })}
              {...RESOURCE_ENDPOINTS.achievements}
            />
          </div>

          <div className="flex gap-2 justify-end pt-4">
            <Button variant="secondary" onClick={() => navigate({ to: "/quests" })}>Cancel</Button>
            <Button type="submit" variant="primary" disabled={createQuest.isPending}>
              {createQuest.isPending ? "Creating..." : "Create Quest"}
            </Button>
          </div>
        </form>
      </div>
    </PageContainer>
  );
}