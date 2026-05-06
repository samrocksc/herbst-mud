import { useState, useEffect, useCallback } from 'react'
import type { ApiError } from '../api'

export type QueryState<T> = Readonly<{
  data: T
  loading: boolean
  error: ApiError | null
  refetch: () => void
}>

export type MutationState = Readonly<{
  loading: boolean
  error: ApiError | null
}>

/**
 * Generic GET query hook.
 * Automatically fetches on mount and provides a refetch trigger.
 * Supply an initialData value so `data` is never undefined.
 */
export function useApiQuery<T>(
  endpoint: string,
  initialData: T,
): QueryState<T> {
  const [data, setData] = useState<T>(initialData)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<ApiError | null>(null)

  const fetchData = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const res = await fetch(endpoint, {
        headers: { Authorization: `Bearer ${localStorage.getItem('token') ?? ''}` },
      })
      if (!res.ok) {
        const err = await res.json().catch(() => ({ error: res.statusText }))
        throw new Error(err.error ?? `HTTP ${res.status}`)
      }
      const payload = await res.json()
      setData(payload)
    } catch (e: unknown) {
      setError({
        message: e instanceof Error ? e.message : 'Unknown error',
        status: e instanceof Response ? e.status : undefined,
      })
    } finally {
      setLoading(false)
    }
  }, [endpoint])

  useEffect(() => {
    void fetchData()
  }, [fetchData])

  return { data, loading, error, refetch: fetchData }
}

/**
 * Generic mutation hook for POST / PUT / DELETE.
 * Returns a mutate function and loading / error state.
 */
export function useApiMutation<T>(
  endpoint: string,
  method: 'POST' | 'PUT' | 'DELETE',
): Readonly<{
  mutate: (body?: unknown) => Promise<T | null>
  loading: boolean
  error: ApiError | null
  reset: () => void
}> {
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<ApiError | null>(null)

  const mutate = useCallback(
    async (body?: unknown) => {
      setLoading(true)
      setError(null)
      try {
        const headers: Record<string, string> = {
          Authorization: `Bearer ${localStorage.getItem('token') ?? ''}`,
        }
        if (method !== 'DELETE') {
          headers['Content-Type'] = 'application/json'
        }
        const res = await fetch(endpoint, {
          method,
          headers,
          body: body ? JSON.stringify(body) : undefined,
        })
        if (!res.ok) {
          const err = await res.json().catch(() => ({ error: res.statusText }))
          throw new Error(err.error ?? `HTTP ${res.status}`)
        }
        if (res.status === 204) return null
        return (await res.json()) as T
      } catch (e: unknown) {
        const err: ApiError = {
          message: e instanceof Error ? e.message : 'Unknown error',
          status: e instanceof Response ? e.status : undefined,
        }
        setError(err)
        throw err
      } finally {
        setLoading(false)
      }
    },
    [endpoint, method],
  )

  const reset = useCallback(() => {
    setError(null)
  }, [])

  return { mutate, loading, error, reset }
}
