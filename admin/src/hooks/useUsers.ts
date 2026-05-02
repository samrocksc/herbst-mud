import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'

const API_BASE = `${window.location.origin}`

export interface User {
  id: number
  email: string
  is_admin: boolean
  created_at: string
  character_name?: string
}

export function useUsers() {
  return useQuery({
    queryKey: ['users'],
    queryFn: async (): Promise<User[]> => {
      const response = await fetch(`${API_BASE}/users`)
      if (!response.ok) throw new Error('Failed to fetch users')
      return response.json()
    }
  })
}

export function useUser(id: number | null) {
  return useQuery({
    queryKey: ['user', id],
    queryFn: async (): Promise<User | null> => {
      if (!id) return null
      const response = await fetch(`${API_BASE}/users/${id}`)
      if (!response.ok) throw new Error('Failed to fetch user')
      return response.json()
    },
    enabled: !!id
  })
}

export function useResetPassword() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async (id: number): Promise<{ message: string }> => {
      const response = await fetch(`${API_BASE}/users/${id}/reset-password`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' }
      })
      if (!response.ok) throw new Error('Failed to reset password')
      return response.json()
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] })
    }
  })
}
