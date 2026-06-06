import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import { useMutation } from "@tanstack/react-query";
import { apiPost } from "../../utils/apiFetch";
import { useCreateRecipe } from "../../hooks/useRecipes";
import { PageHeader } from "../../components/PageHeader";
import { Button } from "../../components/Button";
import { PageContainer } from "../../components/PageContainer";
import { FormField, NumberField, TextareaField, SelectField } from "../../components/FormFields";

export const Route = createFileRoute("/_auth/recipes/new")({
  component: CreateRecipePage,
});

const STATION_OPTIONS = [
  { value: "forge", label: "Forge" },
  { value: "anvil", label: "Anvil" },
  { value: "cooking_utensil", label: "Cooking Utensil" },
  { value: "crafting_table", label: "Crafting Table" },
  { value: "furnace", label: "Furnace" },
  { value: "loom", label: "Loom" },
  { value: "workbench", label: "Workbench" },
  { value: "alchemy_table", label: "Alchemy Table" },
];

const SKILL_OPTIONS = [
  { value: "blades", label: "Blades" },
  { value: "staves", label: "Staves" },
  { value: "blunt", label: "Blunt" },
  { value: "poisons", label: "Poisons" },
  { value: "crafting", label: "Crafting" },
  { value: "cooking", label: "Cooking" },
  { value: "alchemy", label: "Alchemy" },
  { value: "lockpicking", label: "Lockpicking" },
  { value: "search", label: "Search" },
  { value: "stealth", label: "Stealth" },
];

const CLASS_OPTIONS = [
  { value: "any", label: "Any Class" },
  { value: "warrior", label: "Warrior" },
  { value: "rogue", label: "Rogue" },
  { value: "wizard", label: "Wizard" },
  { value: "cleric", label: "Cleric" },
  { value: "ranger", label: "Ranger" },
  { value: "bard", label: "Bard" },
];

function CreateRecipePage() {
  const navigate = useNavigate();
  const { mutate: createRecipe, isPending } = useCreateRecipe();

  const [form, setForm] = useState({
    name: "",
    display_name: "",
    description: "",
    required_station_tag: "",
    required_class: "any",
    required_skill_level: 1,
    required_skill: "",
    inputs: [] as { item_name: string; quantity: number; consumed: boolean }[],
    outputs: [] as { item_name: string; quantity: number }[],
    craft_time_secs: 10,
    world_id: "",
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!form.name.trim()) return;

    // Basic validation
    if (form.inputs.length === 0) {
      alert("Recipe must have at least one input.");
      return;
    }
    if (form.outputs.length === 0) {
      alert("Recipe must have at least one output.");
      return;
    }

    // Transform inputs/outputs to match expected API format
    const recipeData = {
      name: form.name.trim(),
      display_name: form.display_name.trim() || form.name.trim(),
      description: form.description.trim(),
      required_station_tag: form.required_station_tag || "workbench",
      required_class: form.required_class === "any" ? "" : form.required_class,
      required_skill_level: form.required_skill_level,
      required_skill: form.required_skill || "",
      inputs: form.inputs.map((i) => ({
        equipment_template_slug: i.item_name,
        quantity: i.quantity,
        consumed: i.consumed,
      })),
      outputs: form.outputs.map((o) => ({
        equipment_template_slug: o.item_name,
        quantity: o.quantity,
      })),
      craft_time_secs: form.craft_time_secs,
      world_id: form.world_id || "default",
    };

    createRecipe(recipeData, {
      onSuccess: () => {
        navigate({ to: "/recipes" });
      },
      onError: (error: unknown) => {
        console.error("Failed to create recipe:", error);
        const message = error instanceof Error ? error.message : "Failed to create recipe";
        alert(message);
      },
    });
  };

  const set = (patch: Partial<typeof form>) =>
    setForm((prev) => ({ ...prev, ...patch }));

  const addInput = () => {
    setForm((prev) => ({
      ...prev,
      inputs: [
        ...prev.inputs,
        { item_name: "", quantity: 1, consumed: true },
      ],
    }));
  };

  const updateInput = (index: number, patch: Partial<(typeof form)["inputs"][number]>) => {
    setForm((prev) => {
      const newInputs = [...prev.inputs];
      newInputs[index] = { ...newInputs[index], ...patch };
      return { ...prev, inputs: newInputs };
    });
  };

  const removeInput = (index: number) => {
    setForm((prev) => ({
      ...prev,
      inputs: prev.inputs.filter((_, i) => i !== index),
    }));
  };

  const addOutput = () => {
    setForm((prev) => ({
      ...prev,
      outputs: [
        ...prev.outputs,
        { item_name: "", quantity: 1 },
      ],
    }));
  };

  const updateOutput = (index: number, patch: Partial<(typeof form)["outputs"][number]>) => {
    setForm((prev) => {
      const newOutputs = [...prev.outputs];
      newOutputs[index] = { ...newOutputs[index], ...patch };
      return { ...prev, outputs: newOutputs };
    });
  };

  const removeOutput = (index: number) => {
    setForm((prev) => ({
      ...prev,
      outputs: prev.outputs.filter((_, i) => i !== index),
    }));
  };

  return (
    <PageContainer>
      <PageHeader title="Create Recipe" showBack backTo="/recipes" />

      <div className="card bg-surface p-6 border border-border rounded">
        <form onSubmit={handleSubmit} className="space-y-6">
          {/* Basic Information */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label className="text-text-muted text-xs block mb-1">Name *</label>
              <input
                type="text"
                value={form.name}
                onChange={(e) => set({ name: e.target.value })}
                placeholder="recipe_name"
                className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
                required
              />
            </div>
            <div>
              <label className="text-text-muted text-xs block mb-1">Display Name</label>
              <input
                type="text"
                value={form.display_name}
                onChange={(e) => set({ display_name: e.target.value })}
                placeholder="Heavy Plate Boots"
                className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
              />
            </div>
          </div>

          <div>
            <label className="text-text-muted text-xs block mb-1">Description</label>
            <textarea
              value={form.description}
              onChange={(e) => set({ description: e.target.value })}
              placeholder="Crafted with care..."
              rows={2}
              className="w-full p-2 bg-surface border border-border rounded text-text text-sm resize-y"
            />
          </div>

          {/* Requirements */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label className="text-text-muted text-xs block mb-1">Station Tag</label>
              <SelectField
                label="Station Tag"
                value={form.required_station_tag}
                onChange={(v) => set({ required_station_tag: v })}
                options={STATION_OPTIONS}
                placeholder="Select station..."
              />
            </div>
            <div>
              <label className="text-text-muted text-xs block mb-1">Required Class</label>
              <SelectField
                label="Required Class"
                value={form.required_class}
                onChange={(v) => set({ required_class: v })}
                options={CLASS_OPTIONS}
                placeholder="Select class..."
              />
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label className="text-text-muted text-xs block mb-1">Required Skill</label>
              <SelectField
                label="Required Skill"
                value={form.required_skill}
                onChange={(v) => set({ required_skill: v })}
                options={SKILL_OPTIONS}
                placeholder="Select skill..."
              />
            </div>
            <div>
              <label className="text-text-muted text-xs block mb-1">Required Skill Level</label>
              <NumberField
                label="Required Skill Level"
                value={form.required_skill_level}
                onChange={(v) => set({ required_skill_level: v })}
                min={1}
              />
            </div>
          </div>

          <div>
            <label className="text-text-muted text-xs block mb-1">Craft Time (seconds)</label>
            <NumberField
              label="Craft Time Seconds"
              value={form.craft_time_secs}
              onChange={(v) => set({ craft_time_secs: v })}
              min={1}
            />
          </div>

          <div>
            <label className="text-text-muted text-xs block mb-1">World ID</label>
            <input
              type="text"
              value={form.world_id}
              onChange={(e) => set({ world_id: e.target.value })}
              placeholder="default"
              className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
            />
          </div>

          {/* Inputs */}
          <div>
            <h4 className="text-text font-semibold text-sm mt-4 mb-3">Inputs</h4>
            {form.inputs.map((input, i) => (
              <div key={i} className="mb-3 p-3 bg-surface-dark border border-border rounded">
                <div className="grid grid-cols-1 md:grid-cols-5 gap-2 items-center">
                  <div className="md:col-span-2">
                    <label className="text-text-muted text-xs block mb-1">Item Name</label>
                    <input
                      type="text"
                      value={input.item_name}
                      onChange={(e) => updateInput(i, { item_name: e.target.value })}
                      placeholder="item_slug"
                      className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
                    />
                  </div>
                  <div>
                    <label className="text-text-muted text-xs block mb-1">Quantity</label>
                    <NumberField
                      label="Quantity"
                      value={input.quantity}
                      onChange={(v) => updateInput(i, { quantity: v })}
                      min={1}
                    />
                  </div>
                  <div className="flex items-center pt-6">
                    <input
                      type="checkbox"
                      checked={input.consumed}
                      onChange={(e) => updateInput(i, { consumed: e.target.checked })}
                      className="w-4 h-4"
                    />
                    <span className="text-text-muted text-xs ml-2">Consumed</span>
                  </div>
                  <div className="flex items-center justify-end">
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => removeInput(i)}
                      disabled={form.inputs.length <= 1}
                    >
                      Remove
                    </Button>
                  </div>
                </div>
              </div>
            ))}
            <Button variant="secondary" size="sm" onClick={addInput}>
              + Add Input
            </Button>
          </div>

          {/* Outputs */}
          <div>
            <h4 className="text-text font-semibold text-sm mt-6 mb-3">Outputs</h4>
            {form.outputs.map((output, i) => (
              <div key={i} className="mb-3 p-3 bg-surface-dark border border-border rounded">
                <div className="grid grid-cols-1 md:grid-cols-5 gap-2 items-center">
                  <div className="md:col-span-2">
                    <label className="text-text-muted text-xs block mb-1">Item Name</label>
                    <input
                      type="text"
                      value={output.item_name}
                      onChange={(e) => updateOutput(i, { item_name: e.target.value })}
                      placeholder="item_slug"
                      className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
                    />
                  </div>
                  <div>
                    <label className="text-text-muted text-xs block mb-1">Quantity</label>
                    <NumberField
                      label="Quantity"
                      value={output.quantity}
                      onChange={(v) => updateOutput(i, { quantity: v })}
                      min={1}
                    />
                  </div>
                  <div className="flex items-center justify-end md:col-span-2">
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => removeOutput(i)}
                      disabled={form.outputs.length <= 1}
                    >
                      Remove
                    </Button>
                  </div>
                </div>
              </div>
            ))}
            <Button variant="secondary" size="sm" onClick={addOutput}>
              + Add Output
            </Button>
          </div>

          {/* Actions */}
          <div className="flex gap-2 justify-end mt-6">
            <Button variant="secondary" onClick={() => navigate({ to: "/recipes" })}>
              Cancel
            </Button>
            <Button variant="primary" type="submit" disabled={isPending}>
              {isPending ? "Creating..." : "Create Recipe"}
            </Button>
          </div>
        </form>
      </div>
    </PageContainer>
  );
}
