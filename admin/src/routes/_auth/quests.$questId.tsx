import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import {
  useQuest,
  useUpdateQuest,
  useDeleteQuest,
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

export const Route = createFileRoute("/_auth/quests/$questId")({
  component: QuestDetailPage,
});

const REPEAT_MODE_OPTS = [
  { value: "none", label: "None (one-time)" },
  { value: "cooldown", label: "Cooldown" },
  { value: "always", label: "Always repeatable" },
];

const EMPTY_OBJECTIVE: QuestObjective = {
  type: "kill", target_id: "", tag_filter: "", count: 1, labels: [], hint: "",
};

function QuestDetailPage() {
  const questId = Route.useParams().questId;
  const navigate = useNavigate();
  const { data: quest, isLoading, error } = useQuest(Number(questId));
  const { data: lookups, isLoading: lookupsLoading } = useQuestLookups();
  const { data: tags } = useTags();
  const updateQuest = useUpdateQuest();
  const deleteQuest = useDeleteQuest();

  const [formData, setFormData] = useState<QuestInput | null>(null);
  const [confirmDelete, setConfirmDelete] = useState(false);

  if (isLoading) return <div className="loading">Loading quest...</div>;
  if (lookupsLoading) return <div className="loading">Loading options...</div>;
  if (error) return <div className="error">Failed to load quest: {error.message}</div>;
  if (!quest) return <div className="error">Quest not found</div>;

  const current = formData ?? {
    name: quest.name,
    description: quest.description,
    prerequisite_quest_ids: quest.prerequisite_quest_ids ?? [],
    objectives: quest.objectives ?? [],
    rewards: quest.rewards ?? EMPTY_REWARDS,
    repeat_mode: quest.repeat_mode,
    cooldown_hours: quest.cooldown_hours,
    is_active: quest.is_active,
    main_type: quest.main_type ?? "general",
  };

  const set = (patch: Partial<QuestInput>) => setFormData({ ...current, ...patch });

  const addObjective = () => {
    const objs = [...(current.objectives ?? []), { ...EMPTY_OBJECTIVE }];
    set({ objectives: objs });
  };
  const updateObjective = (i: number, patch: Partial<QuestObjective>) => {
    const objs = current.objectives?.map((o, idx) => idx === i ? { ...o, ...patch } : o) ?? [];
    set({ objectives: objs });
  };
  const removeObjective = (i: number) => {
    const objs = current.objectives?.filter((_, idx) => idx !== i) ?? [];
    set({ objectives: objs });
  };

  // Get targets filtered by objective type
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
      await updateQuest.mutateAsync({ id: Number(questId), input: current });
      showToast("Quest updated", "success");
      setFormData(null);
    } catch { /* toasted globally */ }
  };

  const handleDelete = async () => {
    try {
      await deleteQuest.mutateAsync(Number(questId));
      showToast("Quest deleted", "success");
      navigate({ to: "/quests" });
    } catch { /* toasted globally */ }
  };

  return (
    <PageContainer>
      <PageHeader title={quest.name} backTo="/quests" />
      <form onSubmit={handleSubmit} className="form-card space-y-3">
        <FormField label="Name" value={current.name ?? ""} onChange={(v) => set({ name: v })} />
        <TextareaField label="Description" value={current.description ?? ""} onChange={(v) => set({ description: v })} rows={3} />
        <SelectField label="Quest Type" value={current.main_type ?? "general"} onChange={(v) => set({ main_type: v })} options={[
          { value: "general", label: "General" },
          { value: "hunter", label: "Hunter (Kill NPCs)" },
          { value: "collector", label: "Collector (Gather Items)" },
          { value: "explorer", label: "Explorer (Visit Rooms)" },
        ]} />
        <SelectField label="Repeat Mode" value={current.repeat_mode ?? "none"} onChange={(v) => set({ repeat_mode: v })} options={REPEAT_MODE_OPTS} />
        {(current.repeat_mode === "cooldown") && (
          <NumberField label="Cooldown (hours)" value={current.cooldown_hours ?? 0} onChange={(v) => set({ cooldown_hours: v })} />
        )}
        <div className="flex items-center gap-2">
          <input type="checkbox" checked={current.is_active ?? true} onChange={(e) => set({ is_active: e.target.checked })} id="quest-active-edit" />
          <label htmlFor="quest-active-edit" className="text-sm text-text">Active</label>
        </div>

        {/* Prerequisite Quests */}
        <ResourceMultiSelect
          label="Prerequisite Quests"
          value={current.prerequisite_quest_ids ?? []}
          onChange={(ids) => set({ prerequisite_quest_ids: ids as string[] })}
          {...RESOURCE_ENDPOINTS.quests}
        />

        <div className="border-t border-border pt-3 mt-3">
          <div className="flex items-center justify-between mb-2">
            <div>
              <h4 className="text-sm font-semibold text-text">Objectives</h4>
              <p className="text-xs text-text-muted mt-0.5">What the player must do to complete this quest. Add multiple objectives for multi-step quests.</p>
            </div>
            <Button variant="ghost" size="sm" onClick={addObjective}>+ Objective</Button>
          </div>
          {(current.objectives ?? []).map((obj, i) => {
            const targetOptions = getTargetsForType(obj.type);
            const tagOptions = (tags ?? []).map(t => ({ id: t.name, name: t.name }));
            return (
              <div key={i} className="bg-surface-muted border border-border rounded p-3 mb-2 relative">
                <div className="flex items-center justify-between mb-2">
                  <span className="text-xs font-semibold text-primary uppercase tracking-wider">Objective {i + 1}</span>
                  <Button variant="danger" size="sm" onClick={() => removeObjective(i)}>×</Button>
                </div>
                <div className="grid grid-cols-2 sm:grid-cols-3 gap-2 mb-2">
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
                  <NumberField label="Count" value={obj.count} onChange={(v) => updateObjective(i, { count: v })} />
                </div>
                {obj.type && (
                  <div className="mb-2">
                    <SearchableSelect
                      label="Target"
                      value={obj.target_id || ""}
                      onChange={(v) => updateObjective(i, { target_id: v })}
                      options={targetOptions.map(t => ({ id: t.id, name: t.name }))}
                      placeholder={obj.type === "kill" ? "Select NPC..." : obj.type === "explore" ? "Select room..." : "Select item..."}
                    />
                  </div>
                )}
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-2 mb-2">
                  <FormField label="Label" value={obj.labels?.[0] ?? ""} onChange={(v) => updateObjective(i, { labels: [v] })} placeholder={obj.type === "kill" ? "Kill the Rats" : obj.type === "explore" ? "Find the Sewer" : "Gather Herbs"} />
                  <FormField label="Hint" value={obj.hint} onChange={(v) => updateObjective(i, { hint: v })} placeholder="Optional hint for players" />
                </div>
                <div className="border-t border-border/50 pt-2 mt-1">
                  <SearchableSelect
                    label="Tag Filter"
                    value={obj.tag_filter || ""}
                    onChange={(v) => updateObjective(i, { tag_filter: v })}
                    options={tagOptions}
                    placeholder="Filter by tag (optional)..."
                  />
                </div>
              </div>
            );
          })}
        </div>

        <div className="border-t border-border pt-3 mt-3">
          <h4 className="text-sm font-semibold text-text mb-2">Rewards</h4>
          <NumberField label="XP" value={current.rewards?.xp ?? 0} onChange={(v) => set({ rewards: { ...current.rewards ?? EMPTY_REWARDS, xp: v } })} />

          {/* Item Rewards */}
          <ResourceMultiSelect
            label="Item Rewards"
            value={(current.rewards?.item_ids ?? []) as string[]}
            onChange={(ids) => set({ rewards: { ...current.rewards ?? EMPTY_REWARDS, item_ids: ids as string[] } })}
            {...RESOURCE_ENDPOINTS.equipmentTemplates}
          />

          {/* Tag Add Rewards */}
          <div className="mt-3">
            <label className="text-sm text-muted mb-1 block">Tags to Add</label>
            <div className="flex flex-wrap gap-2">
              {(lookups?.tags ?? []).map(tag => (
                <button
                  key={tag.id}
                  type="button"
                  onClick={() => {
                    const currentTags = current.rewards?.tag_adds ?? [];
                    const newTags = currentTags.includes(tag.id)
                      ? currentTags.filter(t => t !== tag.id)
                      : [...currentTags, tag.id];
                    set({ rewards: { ...current.rewards ?? EMPTY_REWARDS, tag_adds: newTags } });
                  }}
                  className={`px-2 py-1 text-xs rounded border ${
                    (current.rewards?.tag_adds ?? []).includes(tag.id)
                      ? "bg-primary/20 border-primary text-text"
                      : "bg-surface border-border text-muted hover:border-primary"
                  }`}
                >
                  {tag.name}
                </button>
              ))}
            </div>
          </div>

          {/* Effect Rewards */}
          <ResourceMultiSelect
            label="Effect Rewards"
            value={current.rewards?.effect_ids ?? []}
            onChange={(ids) => set({ rewards: { ...current.rewards ?? EMPTY_REWARDS, effect_ids: ids as number[] } })}
            {...RESOURCE_ENDPOINTS.effectDefs}
          />

          {/* Achievement Rewards */}
          <ResourceMultiSelect
            label="Achievement Rewards"
            value={current.rewards?.achievement_ids ?? []}
            onChange={(ids) => set({ rewards: { ...current.rewards ?? EMPTY_REWARDS, achievement_ids: ids as number[] } })}
            {...RESOURCE_ENDPOINTS.achievements}
          />
        </div>

        <div className="flex gap-2 pt-1">
          <Button type="submit" variant="primary" disabled={updateQuest.isPending}>
            {updateQuest.isPending ? "Saving..." : "Save Changes"}
          </Button>
          <Button variant="danger" onClick={() => setConfirmDelete(true)}>Delete Quest</Button>
        </div>
      </form>

      {confirmDelete && (
        <div className="modal-overlay" onClick={() => setConfirmDelete(false)}>
          <div className="modal-content modal-sm" onClick={(e) => e.stopPropagation()}>
            <div className="modal-header">
              <h3>Delete Quest</h3>
              <Button variant="ghost" size="sm" onClick={() => setConfirmDelete(false)} aria-label="Close">×</Button>
            </div>
            <div className="modal-body">
              <p>Are you sure you want to delete <strong>{quest.name}</strong>?</p>
              <p className="text-muted">This action cannot be undone.</p>
            </div>
            <div className="modal-footer">
              <Button variant="danger" onClick={handleDelete} disabled={deleteQuest.isPending}>
                {deleteQuest.isPending ? "Deleting..." : "Delete"}
              </Button>
              <Button variant="secondary" onClick={() => setConfirmDelete(false)}>Cancel</Button>
            </div>
          </div>
        </div>
      )}
    </PageContainer>
  );
}