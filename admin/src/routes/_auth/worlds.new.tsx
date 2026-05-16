import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import { useCreateWorld, type WorldInput } from "../../hooks/useWorlds";
import { PageHeader } from "../../components/PageHeader";
import { Button } from "../../components/Button";
import { FormField, TextareaField } from "../../components/FormFields";
import { showToast } from "../../components/Toast";

export const Route = createFileRoute("/_auth/worlds/new")({
  component: CreateWorldPage,
});

const EMPTY_WORLD: WorldInput = {
  name: "",
  title: "",
  description: "",
  active: false,
};

function CreateWorldPage() {
  const navigate = useNavigate();
  const createWorld = useCreateWorld();
  const [formData, setFormData] = useState<WorldInput>(EMPTY_WORLD);

  const set = (patch: Partial<WorldInput>) => setFormData((prev) => ({ ...prev, ...patch }));

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await createWorld.mutateAsync(formData);
      showToast("World created", "success");
      navigate({ to: "/worlds" });
    } catch {
      // Error is toasted by global onError handler
    }
  };

  return (
    <div className="p-6 max-w-[800px] mx-auto">
      <PageHeader title="Create World" showBack backTo="/worlds" />
      <div className="card bg-surface p-6 border border-border rounded">
        <form onSubmit={handleSubmit} className="space-y-4">
          <h3 className="text-text font-semibold mb-4">Basic Information</h3>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <FormField
              label="Name *"
              value={formData.name}
              onChange={(v) => set({ name: v })}
              required
              placeholder="e.g., My Fantasy World"
            />
            <FormField
              label="Title"
              value={formData.title}
              onChange={(v) => set({ title: v })}
              placeholder="Human-readable title"
            />
          </div>
          <TextareaField
            label="Description"
            value={formData.description}
            onChange={(v) => set({ description: v })}
            rows={3}
            placeholder="Describe this world..."
          />

          <div className="flex items-center gap-2 mt-4">
            <label className="flex items-center gap-2 cursor-pointer">
              <input
                type="checkbox"
                checked={formData.active}
                onChange={(e) => set({ active: e.target.checked })}
                className="accent-primary"
              />
              <span className="text-text">Active</span>
            </label>
            <span className="text-text-muted text-sm">(Enable this world for players)</span>
          </div>

          <div className="flex gap-2 justify-end mt-6">
            <Button variant="secondary" onClick={() => navigate({ to: "/worlds" })}>
              Cancel
            </Button>
            <Button variant="primary" type="submit" disabled={createWorld.isPending}>
              {createWorld.isPending ? "Creating..." : "Create World"}
            </Button>
          </div>
        </form>
      </div>
    </div>
  );
}