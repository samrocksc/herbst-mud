import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'

export interface Room {
  id: number
  name: string
  description: string
  x: number
  y: number
  zLevel: number
  exits: Record<string, number>
  items?: number[]
  npcs?: number[]
  createdAt?: string
  updatedAt?: string
}

export interface CreateRoomInput {
  name: string
  description?: string
  x?: number
  y?: number
  zLevel?: number
}

export interface UpdateRoomInput {
  name?: string
  description?: string
  x?: number
  y?: number
  zLevel?: number
  exits?: Record<string, number>
}

const API_BASE = '/api'

async function fetchRooms(): Promise<Room[]> {
  const response = await fetch(`${API_BASE}/rooms`)
  if (!response.ok) {
    throw new Error('Failed to fetch rooms')
  }
  return response.json()
}

async function fetchRoom(id: number): Promise<Room> {
  const response = await fetch(`${API_BASE}/rooms/${id}`)
  if (!response.ok) {
    throw new Error(`Failed to fetch room ${id}`)
  }
  return response.json()
}

async function createRoom(input: CreateRoomInput): Promise<Room> {
  const response = await fetch(`${API_BASE}/rooms`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(input),
  })
  if (!response.ok) {
    throw new Error('Failed to create room')
  }
  return response.json()
}

async function updateRoom(id: number, input: UpdateRoomInput): Promise<Room> {
  const response = await fetch(`${API_BASE}/rooms/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(input),
  })
  if (!response.ok) {
    throw new Error(`Failed to update room ${id}`)
  }
  return response.json()
}

async function deleteRoom(id: number): Promise<void> {
  const response = await fetch(`${API_BASE}/rooms/${id}`, {
    method: 'DELETE',
  })
  if (!response.ok) {
    throw new Error(`Failed to delete room ${id}`)
  }
}

// React Query hooks
export function useRooms() {
  return useQuery({
    queryKey: ['rooms'],
    queryFn: fetchRooms,
    staleTime: 30000, // 30 seconds
  })
}

export function useRoom(id: number) {
  return useQuery({
    queryKey: ['room', id],
    queryFn: () => fetchRoom(id),
    enabled: !!id,
  })
}

export function useCreateRoom() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: createRoom,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['rooms'] })
    },
  })
}

export function useUpdateRoom() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: ({ id, input }: { id: number; input: UpdateRoomInput }) => 
      updateRoom(id, input),
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: ['rooms'] })
      queryClient.setQueryData(['room', data.id], data)
    },
  })
}

export function useDeleteRoom() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: deleteRoom,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['rooms'] })
    },
  })
}

// Transform room to ReactFlow node
export function roomToNode(room: Room) {
  return {
    id: room.id.toString(),
    type: 'room',
    position: { x: room.x || 0, y: room.y || 0 },
    data: {
      name: room.name,
      description: room.description,
      zLevel: room.zLevel || 0,
      exits: room.exits || {},
      items: room.items || [],
      npcs: room.npcs || [],
    },
  }
}

// Transform room exits to ReactFlow edges
export function roomToEdges(room: Room): Array<{
  id: string
  source: string
  target: string
  label: string
  type?: string
}> {
  if (!room.exits) return []
  
  return Object.entries(room.exits).map(([direction, targetId]) => ({
    id: `${room.id}-${targetId}-${direction}`,
    source: room.id.toString(),
    target: targetId.toString(),
    label: direction,
    type: 'smoothstep',
  }))
}