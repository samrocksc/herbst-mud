/**
 * Centralized fetch wrapper. Handles:
 * - Auto-injecting Authorization Bearer token from localStorage
 * - Auto-setting Content-Type: application/json for mutating methods
 * - Auto-parsing JSON responses and unwrapping known response shapes
 * - Returning descriptive errors on non-OK responses
 */
async function apiFetch<T = unknown>(input: RequestInfo, init?: RequestInit): Promise<T> {
  const url = typeof input === "string" ? input : input.url;
  const method = (init?.method ?? (typeof input === "object" ? (input as Request).method : "GET")).toUpperCase();
  const isMutating = ["POST", "PUT", "PATCH", "DELETE"].includes(method);

  const token = localStorage.getItem("token");
  const headers: Record<string, string> = {
    ...(isMutating ? { "Content-Type": "application/json" } : {}),
    ...(init?.headers as Record<string, string>),
    ...(token ? { "Authorization": `Bearer ${token}` } : {}),
  };

  const response = await fetch(url, { ...init, headers });

  if (!response.ok) {
    const getErrorMessage = async (): Promise<string> => {
      try {
        const body = await response.json();
        return body.error || body.message || `HTTP ${response.status} ${response.statusText}`;
      } catch {
        return `HTTP ${response.status} ${response.statusText}`;
      }
    };
    return Promise.reject(new Error(await getErrorMessage()));
  }

  const text = await response.text();
  if (!text) return null as unknown as T;
  try {
    const parsed = JSON.parse(text);
    // Unwrap known { key: [...] } response shapes from the backend,
    // but only when the response is a simple wrapper (1–2 top-level keys)
    const keys = Object.keys(parsed);
    if (keys.length <= 2) {
      const unwrapKeys = ["skills", "npcs", "characters", "abilities", "users", "items", "rooms", "races", "achievements", "factions", "faction_categories", "effects", "hooks", "active_effects", "tags"];
      const matchKey = unwrapKeys.find(k => Object.prototype.hasOwnProperty.call(parsed, k) && Array.isArray(parsed[k]));
      if (matchKey) {
        return parsed[matchKey] as T;
      }
    }
    return parsed as T;
  } catch {
    return text as unknown as T;
  }
}

export const apiGet = <T>(url: string): Promise<T> => apiFetch<T>(url);
export const apiPost = <T>(url: string, body: unknown): Promise<T> =>
  apiFetch<T>(url, { method: "POST", body: JSON.stringify(body) });
export const apiPut = <T>(url: string, body: unknown): Promise<T> =>
  apiFetch<T>(url, { method: "PUT", body: JSON.stringify(body) });
export const apiDelete = <T>(url: string): Promise<T> => apiFetch<T>(url, { method: "DELETE" });
