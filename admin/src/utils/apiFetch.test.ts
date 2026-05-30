import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { apiGet, apiPost, apiPut, apiDelete } from "./apiFetch";

describe("apiFetch", () => {
  beforeEach(() => {
    vi.stubGlobal("fetch", vi.fn());
    vi.stubGlobal("localStorage", { getItem: vi.fn(() => null) });
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  describe("apiGet", () => {
    it("calls fetch with URL", async () => {
      const mockResponse = { abilities: [{ id: 1, name: "Fireball" }] };
      (fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
        ok: true,
        text: () => Promise.resolve(JSON.stringify(mockResponse)),
      });

      await apiGet("/api/abilities");

      expect(fetch).toHaveBeenCalledWith("/api/abilities", expect.any(Object));
    });

    it("unwraps array from response envelope", async () => {
      const mockResponse = { abilities: [{ id: 1 }, { id: 2 }] };
      (fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
        ok: true,
        text: () => Promise.resolve(JSON.stringify(mockResponse)),
      });

      const result = await apiGet<Array<{ id: number }>>("/api/abilities");

      expect(result).toHaveLength(2);
      expect(result).toEqual([{ id: 1 }, { id: 2 }]);
    });
  });

  describe("apiPost", () => {
    it("sends JSON body with POST method", async () => {
      const mockResponse = { id: 1, name: "Fireball" };
      (fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
        ok: true,
        text: () => Promise.resolve(JSON.stringify(mockResponse)),
      });

      await apiPost("/api/abilities", { name: "Fireball" });

      expect(fetch).toHaveBeenCalledWith(
        "/api/abilities",
        expect.objectContaining({
          method: "POST",
          body: JSON.stringify({ name: "Fireball" }),
        }),
      );
    });

    it("sets Content-Type header for POST", async () => {
      const mockResponse = { id: 1 };
      (fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
        ok: true,
        text: () => Promise.resolve(JSON.stringify(mockResponse)),
      });

      await apiPost("/api/abilities", { name: "Test" });

      const call = (fetch as ReturnType<typeof vi.fn>).mock.calls[0];
      const headers = call[1]?.headers as Record<string, string> | undefined;
      expect(headers?.["Content-Type"]).toBe("application/json");
    });
  });

  describe("apiPut", () => {
    it("sends JSON body with PUT method", async () => {
      const mockResponse = { id: 1, name: "Updated" };
      (fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
        ok: true,
        text: () => Promise.resolve(JSON.stringify(mockResponse)),
      });

      await apiPut("/api/abilities/1", { name: "Updated" });

      expect(fetch).toHaveBeenCalledWith(
        "/api/abilities/1",
        expect.objectContaining({
          method: "PUT",
          body: JSON.stringify({ name: "Updated" }),
        }),
      );
    });
  });

  describe("apiDelete", () => {
    it("calls fetch with DELETE method", async () => {
      (fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
        ok: true,
        text: () => Promise.resolve(""),
      });

      await apiDelete("/api/abilities/1");

      expect(fetch).toHaveBeenCalledWith(
        "/api/abilities/1",
        expect.objectContaining({ method: "DELETE" }),
      );
    });
  });

  describe("error handling", () => {
    it("rejects with error message on non-OK response", async () => {
      (fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
        ok: false,
        status: 404,
        statusText: "Not Found",
        json: () => Promise.resolve({ error: "Resource not found" }),
      });

      await expect(apiGet("/api/abilities/999")).rejects.toThrow("Resource not found");
    });

    it("rejects with HTTP status text when no error body", async () => {
      (fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
        ok: false,
        status: 500,
        statusText: "Internal Server Error",
        json: () => Promise.reject(new Error("parse error")),
      });

      await expect(apiGet("/api/abilities")).rejects.toThrow("HTTP 500 Internal Server Error");
    });
  });

  describe("authentication", () => {
    it("adds Bearer token from localStorage", async () => {
      vi.stubGlobal("localStorage", { getItem: vi.fn(() => "test-token-123") });

      const mockResponse = { abilities: [] };
      (fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
        ok: true,
        text: () => Promise.resolve(JSON.stringify(mockResponse)),
      });

      await apiGet("/api/abilities");

      const call = (fetch as ReturnType<typeof vi.fn>).mock.calls[0];
      const headers = call[1]?.headers as Record<string, string> | undefined;
      expect(headers?.["Authorization"]).toBe("Bearer test-token-123");
    });

    it("does not add Authorization header when no token", async () => {
      vi.stubGlobal("localStorage", { getItem: vi.fn(() => null) });

      const mockResponse = { abilities: [] };
      (fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
        ok: true,
        text: () => Promise.resolve(JSON.stringify(mockResponse)),
      });

      await apiGet("/api/abilities");

      const call = (fetch as ReturnType<typeof vi.fn>).mock.calls[0];
      const headers = call[1]?.headers as Record<string, string> | undefined;
      expect(headers?.["Authorization"]).toBeUndefined();
    });
  });

  describe("response parsing", () => {
    it("returns null for empty responses", async () => {
      (fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
        ok: true,
        text: () => Promise.resolve(""),
      });

      const result = await apiGet("/api/abilities");
      expect(result).toBeNull();
    });

    it("returns plain text when JSON parse fails", async () => {
      (fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
        ok: true,
        text: () => Promise.resolve("plain text response"),
      });

      const result = await apiGet("/api/abilities");
      expect(result).toBe("plain text response");
    });

    it("does not unwrap multi-key responses", async () => {
      const mockResponse = { data: [{ id: 1 }], meta: { count: 1 } };
      (fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
        ok: true,
        text: () => Promise.resolve(JSON.stringify(mockResponse)),
      });

      const result = await apiGet<{ data: unknown[]; meta: unknown }>("/api/abilities");

      expect(result).toHaveProperty("data");
      expect(result).toHaveProperty("meta");
    });
  });
});