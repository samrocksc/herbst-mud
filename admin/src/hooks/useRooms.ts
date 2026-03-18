import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'

const API_BASE = 'http://localhost:8080'

export interface Room {
  id: number
  name: string
  description: string
  isStartingRoom: boolean
  exits: Record<string, number>
  atmosphere?: string
  x?: number
  y?: number
  zLevel?: number
}

export interface RoomInput {
  name: string
  description: string
  isStartingRoom?: boolean
  exits?: Record<string, number>
  atmosphere?: string
  x?: number
  y?: number
  zLevel?: number
}

export function useRooms() {
  return useQuery({
    queryKey: ['rooms'],
    queryFn: async (): Promise<Room[]> => {
      const response = await fetch(`${API_BASE}/rooms`)
      if (!response.ok) throw new Error('Failed to fetch rooms')
      return response.json()
    }
  })
}

export function useRoom(id: number | null) {
  return useQuery({
    queryKey: ['room', id],
    queryFn: async (): Promise<Room | null> => {
      if (!id) return null
      const response = await fetch(`${API_BASE}/rooms/${id}`)
      if (!response.ok) throw new Error('Failed to fetch room')
      return response.json()
    },
    enabled: !!id
  })
}

export function useCreateRoom() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async (room: RoomInput): Promise<Room> => {
      const response = await fetch(`${API_BASE}/rooms`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(room)
      })
      if (!response.ok) throw new Error('Failed to create room')
      return response.json()
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['rooms'] })
    }
  })
}

export function useUpdateRoom() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async ({ id, room }: { id: number; room: RoomInput }): Promise<Room> => {
      const response = await fetch(`${API_BASE}/rooms/${id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(room)
      })
      if (!response.ok) throw new Error('Failed to update room')
      return response.json()
    },
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: ['rooms'] })
      queryClient.invalidateQueries({ queryKey: ['room', id] })
    }
  })
}

export function useDeleteRoom() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async (id: number): Promise<void> => {
      const response = await fetch(`${API_BASE}/rooms/${id}`, {
        method: 'DELETE'
      })
      if (!response.ok) throw new Error('Failed to delete room')
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['rooms'] })
    }
  })
}

// Helper to get all available directions
export const DIRECTIONS = ['north', 'south', 'east', 'west', 'up', 'down'] as const
export type Direction = typeof DIRECTIONS[number]

// Get opposite direction
export function getOppositeDirection(dir: Direction): Direction {
  const opposites: Record<Direction, Direction> = {
    north: 'south',
    south: 'north',
    east: 'west',
    west: 'east',
    up: 'down',
    down: 'up'
  }
  return opposites[dir]
}