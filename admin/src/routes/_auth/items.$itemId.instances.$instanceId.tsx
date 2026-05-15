import { createFileRoute, useNavigate } from '@tanstack/react-router';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useState } from 'react';
import { apiGet, apiDelete } from '../../utils/apiFetch';
import { PageHeader } from '../../components/PageHeader';
import { Button } from '../../components/Button';
import { DeleteConfirmation } from '../../components/DeleteConfirmation';
import { InstanceDetailView } from './-items.$itemId.instances.$instanceId.detailView';
import { InstanceEditForm } from './-items.$itemId.instances.$instanceId.editForm';
import type { ItemInstance } from '../../hooks/useItemInstances';

export const Route = createFileRoute('/_auth/items/$itemId/instances/$instanceId')({
  component: ItemInstanceDetail,
});

const API = `${window.location.origin}`;

function ItemInstanceDetail() {
  const { itemId, instanceId } = Route.useParams();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [editing, setEditing] = useState(false);
  const [deleteId, setDeleteId] = useState<string | null>(null);

  const { data: instance, isLoading, error } = useQuery<ItemInstance>({
    queryKey: ['item-instances', instanceId],
    queryFn: () => apiGet<ItemInstance>(`${API}/api/item-instances/${instanceId}`),
  });

  const deleteMutation = useMutation({
    mutationFn: () => apiDelete(`${API}/api/item-instances/${instanceId}`),
    onSuccess: () => { queryClient.invalidateQueries({ queryKey: ['item-instances'] }); navigate({ to: '/items/$itemId', params: { itemId } }); },
  });

  if (isLoading) return <div className="p-8"><PageHeader title="Loading..." backTo={`/items/${itemId}`} /></div>;
  if (error || !instance) return <div className="p-8"><PageHeader title="Error" backTo={`/items/${itemId}`} /><div className="text-danger">Failed to load instance</div></div>;

  return (
    <div className="p-8">
      <PageHeader title={instance.name} backTo={`/items/${itemId}`} actions={
        !editing ? (
          <div className="flex items-center gap-2">
            <Button variant="primary" size="sm" onClick={() => setEditing(true)}>Edit</Button>
            <Button variant="danger" size="sm" onClick={() => setDeleteId(instanceId)}>Delete</Button>
          </div>
        ) : undefined
      } />
      {editing ? <InstanceEditForm instance={instance} instanceId={instanceId} onDone={() => setEditing(false)} /> : <InstanceDetailView instance={instance} />}
      <DeleteConfirmation
        open={!!deleteId}
        title="Delete Item Instance"
        message="Are you sure? This will permanently delete this item instance."
        onConfirm={() => deleteMutation.mutate()}
        onCancel={() => setDeleteId(null)}
        isLoading={deleteMutation.isPending}
      />
    </div>
  );
}