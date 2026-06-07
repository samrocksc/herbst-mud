import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import { useCreateSocial, type SocialInput } from "../../hooks/useSocials";
import { PageHeader } from "../../components/PageHeader";
import { Button } from "../../components/Button";
import { FormField, TextareaField } from "../../components/FormFields";
import { showToast } from "../../components/Toast";

export const Route = createFileRoute("/_auth/socials/new")({
  component: SocialCreatePage,
});

const EMPTY_FORM: SocialInput = {
  name: "",
  self_text: "",
  room_text: "",
  target_self_text: "",
  target_text: "",
  target_room_text: "",
};

function SocialCreatePage() {
  const navigate = useNavigate();
  const createMutation = useCreateSocial();
  const [formData, setFormData] = useState<SocialInput>(EMPTY_FORM);

  const set = (patch: Partial<SocialInput>) => setFormData((prev) => ({ ...prev, ...patch }));

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await createMutation.mutateAsync(formData);
      showToast("Social created", "success");
      navigate({ to: "/socials" });
    } catch (err) {
      console.error("Social creation error:", err);
      const message = err instanceof Error ? err.message : "Failed to create social";
      showToast(message, "error");
    }
  };

  return (
    <div className="p-6 max-w-[900px] mx-auto">
      <PageHeader title="New Social Command" backTo="/socials" />
      <form onSubmit={handleSubmit} className="space-y-4">
        <div className="bg-surface-muted rounded-lg p-6 border border-border">
          <h3 className="mt-0 mb-4 text-text text-lg font-semibold">Basic Texts</h3>
          <div className="space-y-4">
            <FormField label="Name" value={formData.name} onChange={(v) => set({ name: v })} placeholder="smile" />
            <TextareaField label="Self Text" value={formData.self_text} onChange={(v) => set({ self_text: v })} rows={2} placeholder="You smile happily." />
            <TextareaField label="Room Text" value={formData.room_text} onChange={(v) => set({ room_text: v })} rows={2} placeholder="{actor} smiles happily." />
          </div>
        </div>

        <div className="bg-surface-muted rounded-lg p-6 border border-border">
          <h3 className="mt-0 mb-4 text-text text-lg font-semibold">Targeted Texts</h3>
          <div className="space-y-4">
            <TextareaField label="Target Self Text" value={formData.target_self_text} onChange={(v) => set({ target_self_text: v })} rows={2} placeholder="You smile at {target}." />
            <TextareaField label="Target Text" value={formData.target_text} onChange={(v) => set({ target_text: v })} rows={2} placeholder="{actor} smiles at you." />
            <TextareaField label="Target Room Text" value={formData.target_room_text} onChange={(v) => set({ target_room_text: v })} rows={2} placeholder="{actor} smiles at {target}." />
          </div>
          <p className="text-text-muted text-xs mt-2">
            Use {"{actor}"} and {"{target}"} as placeholders. Pronoun substitution available: {"{he}"}, {"{him}"}, {"{his}"}.
          </p>
        </div>

        <Button type="submit" variant="primary" disabled={createMutation.isPending} fullWidth>
          {createMutation.isPending ? "Creating..." : "Create Social"}
        </Button>
      </form>
    </div>
  );
}