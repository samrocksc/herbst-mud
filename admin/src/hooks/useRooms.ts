import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { apiGet, apiPut, apiPost, apiDelete } from '../utils/apiFetch';
import { useWorldStore } from './useWorldStore';

type Room = Readonly<{
  id: number
  name: string
  description: string
  isStartingRoom: boolean
  isRootRoom: boolean
  exits: Record<string, number>
  posX: number
  posY: number
  posZ: number
  version: number
}>

type RoomInput = {
  name: string
  description: string
  isStartingRoom: boolean
  isRootRoom: boolean
  exits: Record<string, number>
  posX: number
  posY: number
  posZ: number
}

type RoomUpdate = Partial<RoomInput> & {
  version?: number
}

const API_BASE = `${window.location.origin}`;

export function useRooms() {
  const { currentWorld } = useWorldStore();
  const queryClient = useQueryClient();

  const roomsQuery = useQuery<Room[]>({
    queryKey: ['rooms', currentWorld],
    queryFn: () => {
      const params = new URLSearchParams();
      if (currentWorld) params.append('world_id', currentWorld);
      const qs = params.toString() ? `?${params.toString()}` : '';
      return apiGet<Room[]>(`${API_BASE}/api/rooms${qs}`);
    },
  });

  const createMutation = useMutation({
    mutationFn: (input: RoomInput) => apiPost<Room>(`${API_BASE}/api/rooms`, input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['rooms'] });
    },
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, update }: { id: number; update: RoomUpdate }) =>
      apiPut(`${API_BASE}/api/rooms/${id}`, update),
    onMutate: async ({ id, update }) => {
      await queryClient.cancelQueries({ queryKey: ['rooms'] });
      const previousRooms = queryClient.getQueryData<Room[]>(['rooms']);
      queryClient.setQueryData<Room[]>(['rooms'], (old) =>
        old?.map(r => r.id === id ? { ...r, ...update } as Room : r)
      );
      return { previousRooms };
    },
    onError: (_err, _vars, context) => {
      if (context?.previousRooms) {
        queryClient.setQueryData(['rooms'], context.previousRooms);
      }
    },
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: ['rooms'] });
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (id: number) => apiDelete(`${API_BASE}/api/rooms/${id}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['rooms'] });
    },
  });

  const cleanupMutation = useMutation({
    mutationFn: () => apiPost<{ cleaned: number }>(`${API_BASE}/api/rooms/cleanup-orphan-exits`, {}),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['rooms'] });
    },
  });

  const bidirectionalExitMutation = useMutation({
    mutationFn: ({ roomId, direction, targetRoomId }: { roomId: number; direction: string; targetRoomId: number }) =>
      apiPost<{ source: Room; target: Room }>(`${API_BASE}/api/rooms/${roomId}/exits/bidirectional`, { direction, targetRoomId }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['rooms'] });
    },
  });

  const removeBidirectionalExitMutation = useMutation({
    mutationFn: ({ roomId, direction }: { roomId: number; direction: string }) =>
      apiDelete(`${API_BASE}/api/rooms/${roomId}/exits/bidirectional?direction=${encodeURIComponent(direction)}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['rooms'] });
    },
  });

  return {
    rooms: roomsQuery.data ?? [],
    isLoading: roomsQuery.isLoading,
    isError: roomsQuery.isError,
    error: roomsQuery.error,
    createRoom: createMutation.mutate,
    createRoomAsync: createMutation.mutateAsync,
    updateRoom: updateMutation.mutate,
    deleteRoom: deleteMutation.mutate,
    isUpdating: updateMutation.isPending,
    isCreating: createMutation.isPending,
    isDeleting: deleteMutation.isPending,
    cleanupOrphanExits: cleanupMutation.mutate,
    createBidirectionalExit: bidirectionalExitMutation.mutateAsync,
    removeBidirectionalExit: removeBidirectionalExitMutation.mutateAsync,
  };
}
