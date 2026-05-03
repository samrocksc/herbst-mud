/**
 * Centralized fetch wrapper. Handles:
 * - Auto-injecting Authorization Bearer token from localStorage
 * - Auto-setting Content-Type: application/json for mutating methods
 * - Auto-parsing JSON responses and unwrapping known response shapes
 * - Throwing descriptive errors on non-OK responses
 */
async function apiFetch<T = unknown>(input: RequestInfo, init?: RequestInit): Promise<T> {
  const url = typeof input === 'string' ? input : input.url
  const method = (init?.method ?? (typeof input === 'object' ? (input as Request).method : 'GET')).toUpperCase()
  const isMutating = ['POST', 'PUT', 'PATCH', 'DELETE'].includes(method)

  const headers: Record<string, string> = {
    ...(isMutating ? { 'Content-Type': 'application/json' } : {}),
    ...(init?.headers as Record<string, string>),
  }
  const token = localStorage.getItem('token')
  if (token) headers['Authorization'] = `Bearer ${token}`

  const response = await fetch(url, { ...init, headers })

  if (!response.ok) {
    let message = `HTTP ${response.status} ${response.statusText}`
    try {
      const body = await response.json()
      message = body.error || body.message || message
    } catch { /* use HTTP status message */ }
    throw new Error(message)
  }

  const text = await response.text()
  if (!text) return undefined as T
  try {
    const parsed = JSON.parse(text)
    // Unwrap known { key: [...] } response shapes from the backend
    for (const key of ['talents', 'skills', 'npcs', 'characters', 'abilities', 'users', 'items', 'rooms', 'races']) {
      if (Object.prototype.hasOwnProperty.call(parsed, key) && Array.isArray(parsed[key])) {
        return parsed[key] as T
      }
    }
    return parsed as T
  } catch {
    return text as unknown as T
  }
}

export const apiGet = <T>(url: string) => apiFetch<T>(url)
export const apiPost = <T>(url: string, body: unknown) =>
  apiFetch<T>(url, { method: 'POST', body: JSON.stringify(body) })
export const apiPut = <T>(url: string, body: unknown) =>
  apiFetch<T>(url, { method: 'PUT', body: JSON.stringify(body) })
export const apiDelete = <T>(url: string) => apiFetch<T>(url, { method: 'DELETE' })
