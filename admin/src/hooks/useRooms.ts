import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { apiGet, apiPut, apiPost, apiDelete } from '../utils/apiFetch'

type Room = Readonly<{
  id: number
  name: string
  description: string
  isStartingRoom: boolean
  exits: Record<string, number>
  posX: number
  posY: number
  version: number
}>

type RoomInput = {
  name: string
  description: string
  isStartingRoom: boolean
  exits: Record<string, number>
  posX: number
  posY: number
}

type RoomUpdate = Partial<RoomInput> & {
  version?: number
}

const API_BASE = `${window.location.origin}`

export function useRooms() {
  const queryClient = useQueryClient()

  const roomsQuery = useQuery<Room[]>({
    queryKey: ['rooms'],
    queryFn: () => apiGet<Room[]>(`${API_BASE}/rooms`),
  })

  const createMutation = useMutation({
    mutationFn: (input: RoomInput) => apiPost<Room>(`${API_BASE}/rooms`, input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['rooms'] })
    },
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, update }: { id: number; update: RoomUpdate }) => 
      apiPut(`${API_BASE}/rooms/${id}`, update),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['rooms'] })
      queryClient.invalidateQueries({ queryKey: ['rooms', variables.id] })
    },
  })

  const deleteMutation = useMutation({
    mutationFn: (id: number) => apiDelete(`${API_BASE}/rooms/${id}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['rooms'] })
    },
  })

  return {
    rooms: roomsQuery.data ?? [],
    isLoading: roomsQuery.isLoading,
    isError: roomsQuery.isError,
    error: roomsQuery.error,
    createRoom: createMutation.mutate,
    updateRoom: updateMutation.mutate,
    deleteRoom: deleteMutation.mutate,
    isUpdating: updateMutation.isPending,
    isCreating: createMutation.isPending,
    isDeleting: deleteMutation.isPending,
  }
}
