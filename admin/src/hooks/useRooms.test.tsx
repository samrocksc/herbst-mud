import React from 'react'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { renderHook, waitFor } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { useRooms, useRoom, useCreateRoom, useUpdateRoom, useDeleteRoom, roomToNode, roomToEdges } from './useRooms'

// Mock fetch globally
global.fetch = vi.fn()

const createWrapper = () => {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: { retry: false },
      mutations: { retry: false },
    },
  })
  
  return ({ children }: { children: React.ReactNode }) => (
    <QueryClientProvider client={queryClient}>
      {children}
    </QueryClientProvider>
  )
}

describe('useRooms', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('fetches rooms successfully', async () => {
    const mockRooms = [
      { id: 1, name: 'Town Square', description: 'Central hub', x: 100, y: 100, zLevel: 0, exits: { north: 2 } },
      { id: 2, name: 'Main Street', description: 'A street', x: 100, y: 200, zLevel: 0, exits: { south: 1 } },
    ]

    ;(global.fetch as ReturnType<typeof vi.fn>).mockResolvedValueOnce({
      ok: true,
      json: async () => mockRooms,
    })

    const { result } = renderHook(() => useRooms(), { wrapper: createWrapper() })

    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    
    expect(result.current.data).toEqual(mockRooms)
    expect(global.fetch).toHaveBeenCalledWith('/api/rooms')
  })

  it('handles fetch error', async () => {
    ;(global.fetch as ReturnType<typeof vi.fn>).mockResolvedValueOnce({
      ok: false,
      status: 500,
    })

    const { result } = renderHook(() => useRooms(), { wrapper: createWrapper() })

    await waitFor(() => expect(result.current.isError).toBe(true))
    
    expect(result.current.error).toBeDefined()
  })

  it('shows loading state initially', () => {
    ;(global.fetch as ReturnType<typeof vi.fn>).mockImplementation(
      () => new Promise(() => {}) // Never resolves
    )

    const { result } = renderHook(() => useRooms(), { wrapper: createWrapper() })

    expect(result.current.isLoading).toBe(true)
    expect(result.current.isFetching).toBe(true)
  })
})

describe('useRoom', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('fetches single room successfully', async () => {
    const mockRoom = { id: 1, name: 'Town Square', description: 'Central hub', x: 100, y: 100, zLevel: 0, exits: {} }

    ;(global.fetch as ReturnType<typeof vi.fn>).mockResolvedValueOnce({
      ok: true,
      json: async () => mockRoom,
    })

    const { result } = renderHook(() => useRoom(1), { wrapper: createWrapper() })

    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    
    expect(result.current.data).toEqual(mockRoom)
    expect(global.fetch).toHaveBeenCalledWith('/api/rooms/1')
  })

  it('does not fetch when id is 0', () => {
    const { result } = renderHook(() => useRoom(0), { wrapper: createWrapper() })

    expect(result.current.isFetching).toBe(false)
    expect(global.fetch).not.toHaveBeenCalled()
  })
})

describe('useCreateRoom', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('creates room successfully', async () => {
    const newRoom = { name: 'New Room', description: 'A new room' }
    const createdRoom = { id: 3, ...newRoom, x: 0, y: 0, zLevel: 0, exits: {} }

    ;(global.fetch as ReturnType<typeof vi.fn>).mockResolvedValueOnce({
      ok: true,
      json: async () => createdRoom,
    })

    const { result } = renderHook(() => useCreateRoom(), { wrapper: createWrapper() })

    result.current.mutate(newRoom)

    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    
    expect(global.fetch).toHaveBeenCalledWith('/api/rooms', expect.objectContaining({
      method: 'POST',
      body: JSON.stringify(newRoom),
    }))
  })
})

describe('useUpdateRoom', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('updates room successfully', async () => {
    const updateInput = { name: 'Updated Room' }
    const updatedRoom = { id: 1, name: 'Updated Room', description: 'Central hub', x: 100, y: 100, zLevel: 0, exits: {} }

    ;(global.fetch as ReturnType<typeof vi.fn>).mockResolvedValueOnce({
      ok: true,
      json: async () => updatedRoom,
    })

    const { result } = renderHook(() => useUpdateRoom(), { wrapper: createWrapper() })

    result.current.mutate({ id: 1, input: updateInput })

    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    
    expect(global.fetch).toHaveBeenCalledWith('/api/rooms/1', expect.objectContaining({
      method: 'PUT',
      body: JSON.stringify(updateInput),
    }))
  })
})

describe('useDeleteRoom', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('deletes room successfully', async () => {
    ;(global.fetch as ReturnType<typeof vi.fn>).mockResolvedValueOnce({
      ok: true,
    })

    const { result } = renderHook(() => useDeleteRoom(), { wrapper: createWrapper() })

    result.current.mutate(1)

    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    
    expect(global.fetch).toHaveBeenCalledWith('/api/rooms/1', expect.objectContaining({
      method: 'DELETE',
    }))
  })
})

describe('roomToNode', () => {
  it('transforms room to ReactFlow node', () => {
    const room = {
      id: 1,
      name: 'Town Square',
      description: 'Central hub',
      x: 100,
      y: 200,
      zLevel: 1,
      exits: { north: 2, east: 3 },
      items: [1, 2],
      npcs: [5],
    }

    const node = roomToNode(room)

    expect(node).toEqual({
      id: '1',
      type: 'room',
      position: { x: 100, y: 200 },
      data: {
        name: 'Town Square',
        description: 'Central hub',
        zLevel: 1,
        exits: { north: 2, east: 3 },
        items: [1, 2],
        npcs: [5],
      },
    })
  })

  it('handles missing optional fields', () => {
    const room = {
      id: 1,
      name: 'Room',
      description: '',
      zLevel: 0,
      exits: {},
    }

    const node = roomToNode(room as Room)

    expect(node.position).toEqual({ x: 0, y: 0 })
    expect(node.data.items).toEqual([])
    expect(node.data.npcs).toEqual([])
  })
})

describe('roomToEdges', () => {
  it('transforms room exits to edges', () => {
    const room = {
      id: 1,
      name: 'Room',
      description: '',
      x: 0,
      y: 0,
      zLevel: 0,
      exits: { north: 2, south: 3, east: 4 },
    } as Room

    const edges = roomToEdges(room)

    expect(edges).toHaveLength(3)
    expect(edges).toContainEqual({
      id: '1-2-north',
      source: '1',
      target: '2',
      label: 'north',
      type: 'smoothstep',
    })
    expect(edges).toContainEqual({
      id: '1-3-south',
      source: '1',
      target: '3',
      label: 'south',
      type: 'smoothstep',
    })
    expect(edges).toContainEqual({
      id: '1-4-east',
      source: '1',
      target: '4',
      label: 'east',
      type: 'smoothstep',
    })
  })

  it('returns empty array for room with no exits', () => {
    const room = {
      id: 1,
      name: 'Room',
      description: '',
      x: 0,
      y: 0,
      zLevel: 0,
    } as Room

    const edges = roomToEdges(room)

    expect(edges).toHaveLength(0)
  })
})