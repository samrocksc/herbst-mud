import { useState, useEffect, useCallback } from 'react'

// ─── Types ──────────────────────────────────────────────────────────────────

export type ApiState<T> = Readonly<{
  data: T | null
  loading: boolean
  error: string | null
}>

// ─── Helpers ──────────────────────────────────────────────────────────────

/** Build a standard RequestInit with auth headers. */
function buildRequest(
  method: string,
  body?: unknown,
): RequestInit {
  const token = localStorage.getItem('token') ?? ''
  const init: RequestInit = {
    method,
    headers: {
      Authorization: `Bearer ${token}`,
    } as Record<string, string>,
  }
  if (body !== undefined) {
    ;(init.headers as Record<string, string>)['Content-Type'] = 'application/json'
    init.body = JSON.stringify(body)
  }
  return init
}

/** Wrap fetch with auth and JSON parsing, throwing on non-OK status. */
async function apiFetch<T>(input: string | URL | Request, init?: RequestInit): Promise<T> {
  const res = await fetch(input, init)
  if (!res.ok) {
    const text = await res.text().catch(() => 'HTTP error')
    throw new Error(`${res.status} ${res.statusText}: ${text}`)
  }
  return (await res.json()) as T
}

// ─── Query hook (GET) ─────────────────────────────────────────────────────

export function useApiQuery<T>(endpoint: string): Readonly<ApiState<T> & { refetch: () => void }> {
  const [data, setData] = useState<T | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const refetch = useCallback(() => {
    let cancelled = false
    setLoading(true)
    setError(null)
    void (async () => {
      try {
        const result = await apiFetch<T>(`${window.location.origin}${endpoint}`, buildRequest('GET'))
        if (!cancelled) setData(result)
      } catch (e) {
        if (!cancelled) setError(e instanceof Error ? e.message : 'Unknown error')
      } finally {
        if (!cancelled) setLoading(false)
      }
    })()
    return () => { cancelled = true }
  }, [endpoint])

  useEffect(() => {
    const cleanup = refetch()
    return cleanup
  }, [refetch])

  return { data, loading, error, refetch }
}

// ─── Mutation hook (POST / PUT / DELETE / PATCH) ───────────────────────────

export function useApiMutation<T, Body = unknown>(
  method: string,
  endpoint: string,
): Readonly<{
  mutate: (body?: Body) => Promise<T | null>
  loading: boolean
  error: string | null
}> {
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const mutate = useCallback(
    async (body?: Body): Promise<T | null> => {
      setLoading(true)
      setError(null)
      try {
        const result = await apiFetch<T>(
          `${window.location.origin}${endpoint}`,
          buildRequest(method, body),
        )
        return result
      } catch (e) {
        const msg = e instanceof Error ? e.message : 'Unknown error'
        setError(msg)
        return null
      } finally {
        setLoading(false)
      }
    },
    [method, endpoint],
  )

  return { mutate, loading, error }
}
