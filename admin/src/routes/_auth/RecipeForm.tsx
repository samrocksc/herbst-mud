/* eslint-disable react-refresh/only-export-components, functional/no-mixed-types */
import { useState } from "react";
import { Button } from "../../components/Button";
import { FormField, TextareaField, NumberField, FormError } from "../../components/fields";
import { SearchableSelect } from "../../components/SearchableSelect";
import { useEquipmentTemplates } from "../../hooks/useEquipmentTemplates";
import { useWorldStore } from "../../contexts/WorldStoreContext";
import type { Recipe, RecipeInput, CraftingInput, CraftingOutput } from "../../hooks/useRecipes";

type RecipeFormProps = Readonly<{
  recipe: Recipe | null
  onSubmit: (data: RecipeInput) => void
  onCancel: () => void
}>

function recipeToForm(r: Recipe): RecipeInput {
  return {
    name: r.name,
    display_name: r.display_name,
    description: r.description ?? "",
    required_station_tag: r.required_station_tag,
    required_class: r.required_class ?? "",
    required_skill_level: r.required_skill_level ?? 0,
    required_skill: r.required_skill ?? "",
    inputs: r.inputs ? [...r.inputs] : [],
    outputs: r.outputs ? [...r.outputs] : [],
    craft_time_secs: r.craft_time_secs ?? 3,
    world_id: r.world_id ?? "default",
  };
}

export { recipeToForm };

function InputRow({
  input,
  onChange,
  onRemove,
  templates,
}: Readonly<{
  input: CraftingInput
  onChange: (input: CraftingInput) => void
  onRemove: () => void
  templates: ReadonlyArray<{ id: string; name: string }>
}>) {
  return (
    <div className="flex gap-2 items-center mb-2">
      <div className="flex-1">
        <SearchableSelect
          options={templates}
          value={input.equipment_template_id}
          onChange={(v) => onChange({ ...input, equipment_template_id: v })}
          placeholder="Search equipment template..."
        />
      </div>
      <input
        type="number"
        placeholder="Qty"
        min="1"
        value={input.quantity}
        onChange={(e) => onChange({ ...input, quantity: parseInt(e.target.value) || 0 })}
        className="w-20 px-3 py-1.5 bg-surface border border-border rounded text-text text-sm"
      />
      <label className="flex items-center gap-1 text-sm">
        <input
          type="checkbox"
          checked={input.consumed}
          onChange={(e) => onChange({ ...input, consumed: e.target.checked })}
          className="rounded"
        />
        Consumed
      </label>
      <Button variant="danger" size="sm" onClick={onRemove}>×</Button>
    </div>
  );
}

function OutputRow({
  output,
  onChange,
  onRemove,
  templates,
}: Readonly<{
  output: CraftingOutput
  onChange: (output: CraftingOutput) => void
  onRemove: () => void
  templates: ReadonlyArray<{ id: string; name: string }>
}>) {
  return (
    <div className="flex gap-2 items-center mb-2">
      <div className="flex-1">
        <SearchableSelect
          options={templates}
          value={output.equipment_template_id}
          onChange={(v) => onChange({ ...output, equipment_template_id: v })}
          placeholder="Search equipment template..."
        />
      </div>
      <input
        type="number"
        placeholder="Qty"
        min="1"
        value={output.quantity}
        onChange={(e) => onChange({ ...output, quantity: parseInt(e.target.value) || 0 })}
        className="w-20 px-3 py-1.5 bg-surface border border-border rounded text-text text-sm"
      />
      <Button variant="danger" size="sm" onClick={onRemove}>×</Button>
    </div>
  );
}

export function RecipeForm({ recipe, onSubmit, onCancel }: RecipeFormProps) {
  const { data: templates } = useEquipmentTemplates();
  const { currentWorld } = useWorldStore();
  const templateOptions = (templates ?? []).map((t) => ({ id: t.id, name: t.name }));
  const [form, setForm] = useState<RecipeInput>(() =>
    recipe ? recipeToForm(recipe) : {
      name: "",
      display_name: "",
      description: "",
      required_station_tag: "",
      required_class: "",
      required_skill_level: 0,
      required_skill: "",
      inputs: [],
      outputs: [],
      craft_time_secs: 3,
      world_id: currentWorld || "default",
    }
  );
  const [submitError, setSubmitError] = useState("");

  const set = <K extends keyof RecipeInput>(key: K, value: RecipeInput[K]) =>
    setForm((prev) => ({ ...prev, [key]: value }));

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!form.name.trim()) {
      setSubmitError("Name is required");
      return;
    }
    if (!form.required_station_tag.trim()) {
      setSubmitError("Station tag is required");
      return;
    }
    onSubmit(form);
  };

  const addInput = () => {
    set("inputs", [...form.inputs, { equipment_template_id: "", quantity: 1, consumed: true }]);
  };

  const updateInput = (index: number, input: CraftingInput) => {
    const updated = [...form.inputs];
    updated[index] = input;
    set("inputs", updated);
  };

  const removeInput = (index: number) => {
    set("inputs", form.inputs.filter((_, i) => i !== index));
  };

  const addOutput = () => {
    set("outputs", [...form.outputs, { equipment_template_id: "", quantity: 1 }]);
  };

  const updateOutput = (index: number, output: CraftingOutput) => {
    const updated = [...form.outputs];
    updated[index] = output;
    set("outputs", updated);
  };

  const removeOutput = (index: number) => {
    set("outputs", form.outputs.filter((_, i) => i !== index));
  };

  return (
    <div className="form-card space-y-3">
      <h3 className="mt-0 mb-0 text-text text-base font-semibold">
        {recipe ? "Edit Recipe" : "Add New Recipe"}
      </h3>
      {submitError && <FormError message={submitError} />}
      <form onSubmit={handleSubmit} className="space-y-3">
        <FormField
          label="Name (unique ID)"
          value={form.name}
          onChange={(v) => set("name", v)}
          placeholder="e.g. pepperoni_pizza"
          required
          disabled={!!recipe}
        />
        <FormField
          label="Display Name"
          value={form.display_name}
          onChange={(v) => set("display_name", v)}
          placeholder="e.g. Pepperoni Pizza"
        />
        <TextareaField
          label="Description"
          value={form.description}
          onChange={(v) => set("description", v)}
          rows={2}
          placeholder="What this recipe creates..."
        />
        <div className="grid grid-cols-2 gap-3">
          <FormField
            label="Required Station Tag"
            value={form.required_station_tag}
            onChange={(v) => set("required_station_tag", v)}
            placeholder="e.g. kitchen, forge, alchemy"
            required
          />
          <FormField
            label="Required Class"
            value={form.required_class}
            onChange={(v) => set("required_class", v)}
            placeholder="Leave empty for any class"
          />
        </div>
        <div className="grid grid-cols-3 gap-3">
          <FormField
            label="Required Skill"
            value={form.required_skill}
            onChange={(v) => set("required_skill", v)}
            placeholder="e.g. cooking, smithing"
          />
          <NumberField
            label="Required Skill Level"
            value={form.required_skill_level}
            onChange={(v) => set("required_skill_level", v)}
            min={0}
          />
          <NumberField
            label="Craft Time (seconds)"
            value={form.craft_time_secs}
            onChange={(v) => set("craft_time_secs", v)}
            min={1}
          />
        </div>
        <FormField
          label="World ID"
          value={form.world_id}
          onChange={(v) => set("world_id", v)}
          placeholder="default"
        />

        <div>
          <label className="block text-sm font-medium text-text mb-1">Inputs</label>
          {form.inputs.map((input, i) => (
            <InputRow
              key={i}
              input={input}
              templates={templateOptions}
              onChange={(updated) => updateInput(i, updated)}
              onRemove={() => removeInput(i)}
            />
          ))}
          <Button variant="secondary" size="sm" onClick={addInput} type="button">
            + Add Input
          </Button>
        </div>

        <div>
          <label className="block text-sm font-medium text-text mb-1">Outputs</label>
          {form.outputs.map((output, i) => (
            <OutputRow
              key={i}
              output={output}
              templates={templateOptions}
              onChange={(updated) => updateOutput(i, updated)}
              onRemove={() => removeOutput(i)}
            />
          ))}
          <Button variant="secondary" size="sm" onClick={addOutput} type="button">
            + Add Output
          </Button>
        </div>

        <div className="flex gap-2 pt-1">
          <Button type="submit" variant="primary" fullWidth>
            {recipe ? "Update Recipe" : "Create Recipe"}
          </Button>
          <Button type="button" variant="secondary" onClick={onCancel} fullWidth>
            Cancel
          </Button>
        </div>
      </form>
    </div>
  );
}