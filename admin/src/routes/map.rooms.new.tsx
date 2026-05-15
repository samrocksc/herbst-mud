import { createFileRoute, useNavigate } from '@tanstack/react-router';
import { useState } from 'react';
import { PageHeader } from '../components/PageHeader';
import { Button } from '../components/Button';
import { FormField } from '../components/fields/FormField';
import { TextareaField } from '../components/fields/TextareaField';
import { FormError } from '../components/fields/FormError';
import { useRooms } from '../hooks/useRooms';

export const Route = createFileRoute('/map/rooms/new')({
  component: CreateRoomPage,
});

function CreateRoomPage() {
  const navigate = useNavigate();
  const { createRoomAsync, rooms } = useRooms();
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [error, setError] = useState<string | null>(null);

  // Check if any room currently has isRootRoom=true
  const hasRootRoom = rooms.some(r => r.isRootRoom);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) return;
    setError(null);
    try {
      const room = await createRoomAsync({
        name: name.trim(),
        description: description.trim(),
        isStartingRoom: false,
        // Auto-set root room if none exists
        isRootRoom: !hasRootRoom,
        exits: {},
        posX: 0,
        posY: 0,
        posZ: 0,
      });
      navigate({ to: '/map', search: { room: room.id } });
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create room');
    }
  };

  return (
    <div className="p-6 max-w-[600px] mx-auto">
      <PageHeader title="Create New Room" showBack backTo="/map" />
      <div className="bg-surface p-6 border border-border rounded">
        <form onSubmit={handleSubmit} className="space-y-4">
          {error && <FormError message={error} />}
          <FormField label="Room Name" value={name} onChange={setName} placeholder="Enter room name" required />
          <TextareaField label="Description" value={description} onChange={setDescription} rows={4} placeholder="Enter room description" />
          {!hasRootRoom && (
            <div className="text-text-muted text-sm bg-surface-muted p-2 rounded border border-border">
              No root room exists yet. This room will be set as the root room (where new characters spawn).
            </div>
          )}
          <div className="flex gap-2 pt-2">
            <Button type="submit" variant="primary" disabled={!name.trim()}>
              Create Room
            </Button>
            <Button type="button" variant="secondary" onClick={() => navigate({ to: '/map' })}>
              Cancel
            </Button>
          </div>
        </form>
      </div>
    </div>
  );
}
