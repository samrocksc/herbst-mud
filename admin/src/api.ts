/** Central API client utilities for the admin UI */

export function getToken(): string {
  return localStorage.getItem('token') ?? ''
}

export type ApiError = Readonly<{
  message: string
  status?: number
}>

/**
 * Core fetch wrapper that injects auth header, JSON content-type, and
 * normalises errors.
 */
export async function apiFetch(path: string, opts?: RequestInit): Promise<unknown> {
  const res = await fetch(path, {
    ...opts,
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${getToken()}`,
      ...opts?.headers,
    },
  })
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }))
    throw new Error(err.error ?? `HTTP ${res.status}`)
  }
  if (res.status === 204) return null
  return res.json()
}

/** GET helper with typed return */
export async function apiGet<T>(path: string): Promise<T> {
  return apiFetch(path) as Promise<T>
}

/** POST helper */
export async function apiPost<T>(path: string, body: unknown): Promise<T> {
  return apiFetch(path, { method: 'POST', body: JSON.stringify(body) }) as Promise<T>
}

/** PUT helper */
export async function apiPut<T>(path: string, body: unknown): Promise<T> {
  return apiFetch(path, { method: 'PUT', body: JSON.stringify(body) }) as Promise<T>
}

/** DELETE helper */
export async function apiDelete(path: string): Promise<void> {
  await apiFetch(path, { method: 'DELETE' })
}
