import { useState } from "react";
import {
  useUpdateAbility,
  useDeleteAbility,
  type Ability,
  type AbilityInput,
} from "../../hooks/useAbilities";
import { useTags } from "../../hooks/useTags";
import { useGameSkills } from "../../hooks/useGameSkills";
import { useNavigate } from "@tanstack/react-router";
import { Button } from "../../components/Button";
import { SearchableSelect } from "../../components/SearchableSelect";
import { EffectsSubForm } from "../../components/EffectsSubForm";
import { DeleteConfirmation } from "../../components/DeleteConfirmation";
import {
  FormField,
  NumberField,
  TextareaField,
  SelectField,
} from "../../components/FormFields";
import { showToast } from "../../components/Toast";

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

export function AbilityEditForm({
  ability,
  abilityId,
  onDone,
}: Readonly<{
  ability: Ability;
  abilityId: number;
  onDone: () => void;
}>) {
  const navigate = useNavigate();
  const updateAbility = useUpdateAbility();
  const deleteAbility = useDeleteAbility();
  const { data: availableTags } = useTags();
  const { data: gameSkills } = useGameSkills();
  const [showDeleteModal, setShowDeleteModal] = useState(false);

  const [formData, setFormData] = useState<AbilityInput>({
    name: ability.name,
    description: ability.description,
    ability_type: ability.ability_type,
    requirements: ability.requirements ?? "",
    cooldown_seconds: ability.cooldown_seconds,
    mana_cost: ability.mana_cost,
    stamina_cost: ability.stamina_cost,
    hp_cost: ability.hp_cost,
    proc_chance: ability.proc_chance,
    proc_event: ability.proc_event ?? "",
    ability_class: ability.ability_class,
    required_tag: ability.required_tag ?? "",
    slug: ability.slug ?? "",
    caster_message: ability.caster_message ?? "",
    recipient_message: ability.recipient_message ?? "",
    required_skill_id: ability.required_skill_id ?? null,
    required_skill_level: ability.required_skill_level ?? 0,
  });

  const set = (patch: Partial<AbilityInput>) =>
    setFormData((prev) => ({ ...prev, ...patch }));

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await updateAbility.mutateAsync({ id: abilityId, input: formData });
      showToast("Ability updated", "success");
      onDone();
    } catch {
      // Error is toasted by global onError handler
    }
  };

  const handleDelete = async () => {
    try {
      await deleteAbility.mutateAsync(abilityId);
      showToast("Ability deleted", "success");
      navigate({ to: "/abilities" });
    } catch {
      // Error is toasted by global onError handler
    }
  };

  // Refinement #6: only show Proc Settings for passive abilities
  const isPassive = formData.ability_class === "passive";

  return (
    <div className="bg-surface-muted rounded-lg p-6 border border-border mb-6">
      <h2 className="mt-0 mb-4 text-text text-lg font-semibold">Edit Ability</h2>
      <form onSubmit={handleSubmit} className="space-y-6">
        {/* Refinement #12: sectioned with left-border accent + heading + help text */}
        <section>
          <h3 className="text-text font-semibold mb-1 border-l-4 border-primary pl-3">Basic Information</h3>
          <p className="text-xs text-muted mb-3 ml-1">Name, type, class. Required: Name.</p>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <FormField
              label="Name"
              value={formData.name}
              onChange={(v) => set({ name: v })}
              required
              placeholder="e.g., Fireball"
            />
            <FormField
              label="Slug (optional)"
              value={formData.slug ?? ""}
              onChange={(v) => set({ slug: v })}
              placeholder="Auto-generated from name if empty"
            />
            <SelectField
              label="Ability Type"
              value={formData.ability_type}
              onChange={(v) => set({ ability_type: v })}
              options={ABILITY_TYPE_OPTS}
              tooltip="Combat = melee, Magic = spells, Utility = misc actions, Healing = restore HP, Support = buff allies, Defensive = ward or shield"
            />
            <SearchableSelect
              label="Required Tag (optional)"
              options={(availableTags ?? []).map((t) => ({ id: t.name, name: t.name }))}
              value={formData.required_tag || ""}
              onChange={(v) => set({ required_tag: v })}
              placeholder="No tag required"
            />
            <SelectField
              label="Ability Class"
              value={formData.ability_class}
              onChange={(v) => set({ ability_class: v })}
              options={ABILITY_CLASS_OPTS}
              tooltip="Active = use on your turn. Passive = auto-triggers from Proc Settings. Toggle = turn on/off."
            />
          </div>
          <div className="mt-4">
            <TextareaField
              label="Description"
              value={formData.description}
              onChange={(v) => set({ description: v })}
              rows={3}
              placeholder="What the ability does — shown in the help panel and combat log."
            />
          </div>
        </section>

        {/* Refinement #1 + #2 + #12: removed legacy Cost/Cooldown (ticks) */}
        <section>
          <h3 className="text-text font-semibold mb-1 border-l-4 border-primary pl-3">Costs & Cooldown</h3>
          <p className="text-xs text-muted mb-3 ml-1">
            Active abilities cost mana, stamina, or HP to use. Set to 0 for abilities that don't consume resources.
          </p>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <NumberField
              label="Cooldown (s)"
              tooltip="Cooldown in seconds before the ability can be used again. 0 = no cooldown."
              value={formData.cooldown_seconds}
              onChange={(v) => set({ cooldown_seconds: v })}
            />
            <FormField
              label="Unlock Tags (JSON)"
              tooltip='JSON prerequisites for unlocking this ability. Example: {"tags":["level:5"]} means the character must be level 5 or higher.'
              value={formData.requirements}
              onChange={(v) => set({ requirements: v })}
              placeholder='{"tags":["level:5"]}'
            />
          </div>
          <div className="grid grid-cols-3 gap-4 mt-4">
            <NumberField
              label="Mana Cost"
              tooltip="Mana points consumed when using this ability. Use for magical abilities."
              value={formData.mana_cost}
              onChange={(v) => set({ mana_cost: v })}
            />
            <NumberField
              label="Stamina Cost"
              tooltip="Stamina points consumed when using this ability. Use for physical abilities."
              value={formData.stamina_cost}
              onChange={(v) => set({ stamina_cost: v })}
            />
            <NumberField
              label="HP Cost"
              tooltip="HP sacrificed to use the ability. Use for blood magic or self-damaging abilities."
              value={formData.hp_cost}
              onChange={(v) => set({ hp_cost: v })}
            />
          </div>
        </section>

        {/* Refinement #6: only show for passive abilities */}
        {isPassive && (
          <section>
            <h3 className="text-text font-semibold mb-1 border-l-4 border-warning pl-3">Proc Settings</h3>
            <p className="text-xs text-muted mb-3 ml-1">
              Passives trigger automatically from combat events. Set the chance and the event that fires them.
            </p>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <NumberField
                label="Proc Chance (0–1)"
                tooltip="Probability of triggering on each relevant event. 0.15 = 15% chance."
                value={formData.proc_chance}
                onChange={(v) => set({ proc_chance: v })}
                step={0.01}
              />
              <FormField
                label="Proc Event"
                tooltip="Combat event that triggers this passive: on_hit, on_hit_received, on_crit, on_kill."
                value={formData.proc_event}
                onChange={(v) => set({ proc_event: v })}
                placeholder="e.g., on_hit, on_crit"
              />
            </div>
          </section>
        )}

        {/* Refinement #7: collapsed, no double-labels */}
        <section>
          <h3 className="text-text font-semibold mb-1 border-l-4 border-primary pl-3">Combat Messages</h3>
          <p className="text-xs text-muted mb-3 ml-1">
            Templates shown in the combat log. Use <code className="bg-muted px-1 rounded">{'{actor}'}</code>{" "}
            for the user and <code className="bg-muted px-1 rounded">{'{target}'}</code> for the recipient.
          </p>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <FormField
              label="Caster Message"
              tooltip='What the actor sees. Example: "{actor} casts fireball"'
              value={formData.caster_message}
              onChange={(v) => set({ caster_message: v })}
              placeholder="e.g., {actor} casts fireball"
            />
            <FormField
              label="Recipient Message"
              tooltip='What the target sees. Example: "{actor} casts fireball at you"'
              value={formData.recipient_message}
              onChange={(v) => set({ recipient_message: v })}
              placeholder="e.g., {actor} casts fireball at you"
            />
          </div>
        </section>

        {/* Phase 2: Skill Requirement */}
        <section>
          <h3 className="text-text font-semibold mb-1 border-l-4 border-primary pl-3">Skill Requirement (optional)</h3>
          <p className="text-xs text-muted mb-3 ml-1">
            Gate this ability behind a skill level. Characters must have the selected skill at or above the required level to use this ability. Leave blank for no requirement.
          </p>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <SelectField
              label="Required Skill"
              value={String(formData.required_skill_id ?? "")}
              onChange={(v) => set({ required_skill_id: v ? Number(v) : null, required_skill_level: v ? formData.required_skill_level : 0 })}
              options={(gameSkills ?? []).map((s) => ({ value: String(s.id), label: s.display_name }))}
              placeholder="— None —"
              tooltip="Select a skill that the character must possess to use this ability."
            />
            <NumberField
              label="Required Skill Level"
              value={formData.required_skill_level}
              onChange={(v) => set({ required_skill_level: v })}
              tooltip="Minimum skill level required. 0 means no level requirement."
            />
          </div>
        </section>

        <EffectsSubForm abilityId={abilityId} />

        <div className="flex gap-2 pt-4 border-t border-border">
          <Button type="submit" variant="primary" disabled={updateAbility.isPending}>
            {updateAbility.isPending ? "Saving..." : "Save Changes"}
          </Button>
          <Button variant="secondary" onClick={onDone} type="button">
            Cancel
          </Button>
          <Button variant="danger" onClick={() => setShowDeleteModal(true)} type="button" className="ml-auto">
            Delete Ability
          </Button>
        </div>
      </form>

      <DeleteConfirmation
        open={showDeleteModal}
        title="Delete Ability"
        message={`Are you sure you want to delete "${ability.name}"? This action cannot be undone.`}
        onConfirm={handleDelete}
        onCancel={() => setShowDeleteModal(false)}
        isLoading={deleteAbility.isPending}
      />
    </div>
  );
}
