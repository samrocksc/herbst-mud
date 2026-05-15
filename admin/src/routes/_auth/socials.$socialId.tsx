import { createFileRoute, useNavigate } from '@tanstack/react-router';
import { useState } from 'react';
import { useSocial, useUpdateSocial, useDeleteSocial, type SocialInput } from '../../hooks/useSocials';
import { PageHeader } from '../../components/PageHeader';
import { Button } from '../../components/Button';
import { FormField, TextareaField } from '../../components/FormFields';
import { DeleteConfirmation } from '../../components/DeleteConfirmation';
import { showToast } from '../../components/Toast';

export const Route = createFileRoute('/_auth/socials/$socialId')({
  component: SocialDetailPage,
});

function SocialDetailPage() {
  const socialId = Route.useParams().socialId;
  const navigate = useNavigate();
  const { data: social, isLoading, error } = useSocial(Number(socialId));
  const updateMutation = useUpdateSocial();
  const deleteMutation = useDeleteSocial();
  const [showDelete, setShowDelete] = useState(false);
  const [formData, setFormData] = useState<SocialInput | null>(null);

  if (isLoading) return <div className="loading">Loading social...</div>;
  if (error) return <div className="error">Failed to load social: {error.message}</div>;
  if (!social) return <div className="error">Social not found</div>;

  const current = formData ?? {
    name: social.name,
    self_text: social.self_text,
    room_text: social.room_text,
    target_self_text: social.target_self_text,
    target_text: social.target_text,
    target_room_text: social.target_room_text,
  };

  const set = (patch: Partial<SocialInput>) => setFormData({ ...current, ...patch });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await updateMutation.mutateAsync({ id: Number(socialId), input: current });
      showToast('Social updated', 'success');
      setFormData(null);
    } catch { /* toasted globally */ }
  };

  const handleDelete = async () => {
    try {
      await deleteMutation.mutateAsync(Number(socialId));
      showToast('Social deleted', 'success');
      navigate({ to: '/socials' });
    } catch { /* toasted globally */ }
  };

  return (
    <div className="p-6 max-w-[900px] mx-auto">
      <PageHeader title={social.name} backTo="/socials" />
      <form onSubmit={handleSubmit} className="space-y-4">
        <div className="bg-surface-muted rounded-lg p-6 border border-border">
          <h3 className="mt-0 mb-4 text-text text-lg font-semibold">Basic Texts</h3>
          <div className="space-y-4">
            <FormField label="Name" value={current.name} onChange={(v) => set({ name: v })} placeholder="smile" />
            <TextareaField label="Self Text" value={current.self_text} onChange={(v) => set({ self_text: v })} rows={2} placeholder="You smile happily." />
            <TextareaField label="Room Text" value={current.room_text} onChange={(v) => set({ room_text: v })} rows={2} placeholder="{actor} smiles happily." />
          </div>
        </div>

        <div className="bg-surface-muted rounded-lg p-6 border border-border">
          <h3 className="mt-0 mb-4 text-text text-lg font-semibold">Targeted Texts</h3>
          <div className="space-y-4">
            <TextareaField label="Target Self Text" value={current.target_self_text} onChange={(v) => set({ target_self_text: v })} rows={2} placeholder="You smile at {target}." />
            <TextareaField label="Target Text" value={current.target_text} onChange={(v) => set({ target_text: v })} rows={2} placeholder="{actor} smiles at you." />
            <TextareaField label="Target Room Text" value={current.target_room_text} onChange={(v) => set({ target_room_text: v })} rows={2} placeholder="{actor} smiles at {target}." />
          </div>
          <p className="text-text-muted text-xs mt-2">
            Use {'{actor}'} and {'{target}'} as placeholders. Pronoun substitution available: {'{he}'}, {'{him}'}, {'{his}'}.
          </p>
        </div>

        <div className="flex gap-2">
          <Button type="submit" variant="primary" disabled={updateMutation.isPending}>
            {updateMutation.isPending ? 'Saving...' : 'Save Changes'}
          </Button>
          <Button variant="danger" onClick={() => setShowDelete(true)} type="button">Delete</Button>
        </div>
      </form>

      <DeleteConfirmation
        open={showDelete}
        title="Delete Social Command"
        message="Are you sure you want to delete this social command? This action cannot be undone."
        onConfirm={handleDelete}
        onCancel={() => setShowDelete(false)}
        isLoading={deleteMutation.isPending}
      />
    </div>
  );
}