import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import {
  useCreateAbility,
  type AbilityInput,
} from "../../hooks/useAbilities";
import { useTags } from "../../hooks/useTags";
import { PageHeader } from "../../components/PageHeader";
import { Button } from "../../components/Button";
import { TagInput } from "../../components/TagInput";
import {
  FormField,
  NumberField,
  TextareaField,
  SelectField,
} from "../../components/FormFields";
import { showToast } from "../../components/Toast";
import { PageContainer } from "../../components/PageContainer";

export const Route = createFileRoute("/_auth/abilities/new")({
  component: CreateAbilityPage,
});

const ABILITY_TYPE_OPTS = [
  { value: "combat", label: "Combat" },
  { value: "magic", label: "Magic" },
  { value: "utility", label: "Utility" },
  { value: "healing", label: "Healing" },
  { value: "support", label: "Support" },
  { value: "defensive", label: "Defensive" },
];

const ABILITY_CLASS_OPTS = [
  { value: "active", label: "Active" },
  { value: "passive", label: "Passive" },
  { value: "toggle", label: "Toggle" },
];

const EMPTY_ABILITY: AbilityInput = {
  name: "",
  description: "",
  ability_type: "combat",
  requirements: "1",
  cost: 0,
  cooldown: 0,
  cooldown_seconds: 0,
  mana_cost: 0,
  stamina_cost: 0,
  hp_cost: 0,
  proc_chance: 0,
  proc_event: "",
  ability_class: "active",
  required_tag: "",
  slug: "",
};

function CreateAbilityPage() {
  const navigate = useNavigate();
  const createAbility = useCreateAbility();
  const { data: availableTags } = useTags();
  const [formData, setFormData] = useState<AbilityInput>(EMPTY_ABILITY);

  const selectedTags = formData.required_tag
    ? formData.required_tag.split(",").map((t) => t.trim()).filter(Boolean)
    : [];

  const set = (patch: Partial<AbilityInput>) => setFormData((prev) => ({ ...prev, ...patch }));

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await createAbility.mutateAsync(formData);
      showToast("Ability created", "success");
      navigate({ to: "/abilities" });
    } catch {
      // Error is toasted by global onError handler
    }
  };

  return (
    <PageContainer>
      <PageHeader title="Create Ability" showBack backTo="/abilities" />
      <div className="card bg-surface p-6 border border-border rounded">
        <form onSubmit={handleSubmit} className="space-y-4">
          <h3 className="text-text font-semibold mb-4">Basic Information</h3>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <FormField label="Name *" value={formData.name} onChange={(v) => set({ name: v })} required />
            <FormField label="Slug (optional)" value={formData.slug ?? ""} onChange={(v) => set({ slug: v })} placeholder="Auto-generated from name if empty" />
            <SelectField label="Ability Type" value={formData.ability_type} onChange={(v) => set({ ability_type: v })} options={ABILITY_TYPE_OPTS} />
            <TagInput
              label="Required Tag (optional)"
              value={selectedTags}
              onChange={(tags) => set({ required_tag: tags.join(", ") })}
              availableTags={(availableTags ?? []).map((t) => t.name)}
              placeholder="e.g., sword, fire, healing"
            />
            <SelectField label="Ability Class" value={formData.ability_class} onChange={(v) => set({ ability_class: v })} options={ABILITY_CLASS_OPTS} />
          </div>
          <TextareaField label="Description" value={formData.description} onChange={(v) => set({ description: v })} rows={3} />

          <h3 className="text-text font-semibold mt-6 mb-4">Costs & Cooldown</h3>
          <div className="grid grid-cols-3 gap-4">
            <NumberField label="Cost" value={formData.cost} onChange={(v) => set({ cost: v })} />
            <NumberField label="Cooldown (s)" value={formData.cooldown_seconds} onChange={(v) => set({ cooldown_seconds: v })} />
            <FormField label="Unlock Tags (JSON)" value={formData.requirements} onChange={(v) => set({ requirements: v })} placeholder='{"tags":["level:5"]}' />
          </div>
          <div className="grid grid-cols-3 gap-4">
            <NumberField label="Mana Cost" value={formData.mana_cost} onChange={(v) => set({ mana_cost: v })} />
            <NumberField label="Stamina Cost" value={formData.stamina_cost} onChange={(v) => set({ stamina_cost: v })} />
            <NumberField label="HP Cost" value={formData.hp_cost} onChange={(v) => set({ hp_cost: v })} />
          </div>

          <h3 className="text-text font-semibold mt-6 mb-4">Proc Settings</h3>
          <div className="grid grid-cols-2 gap-4">
            <NumberField label="Proc Chance (0–1)" value={formData.proc_chance} onChange={(v) => set({ proc_chance: v })} step={0.01} />
            <FormField label="Proc Event" value={formData.proc_event} onChange={(v) => set({ proc_event: v })} placeholder="e.g., on_hit, on_crit" />
          </div>

          <div className="flex gap-2 justify-end mt-6">
            <Button variant="secondary" onClick={() => navigate({ to: "/abilities" })}>Cancel</Button>
            <Button variant="primary" type="submit" disabled={createAbility.isPending}>
              {createAbility.isPending ? "Creating..." : "Create Ability"}
            </Button>
          </div>
        </form>
      </div>
    </PageContainer>
  );
}
