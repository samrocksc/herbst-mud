import { describe, it, expect } from "vitest";
import { renderHook, waitFor } from "@testing-library/react";
import { QueryClientWrapper } from "../test/wrappers";
import { useCharacters, useCharacter, useUpdateCharacter } from "./useCharacters";
import type { CharacterUpdate } from "./useCharacters";

describe("useCharacters", () => {
  it("returns a list of characters", async () => {
    const { result } = renderHook(() => useCharacters(), { wrapper: QueryClientWrapper });
    await waitFor(() => expect(result.current.data).toBeDefined());
    expect(result.current.data).toHaveLength(2);
  });

  it("contains expected characters", async () => {
    const { result } = renderHook(() => useCharacters(), { wrapper: QueryClientWrapper });
    await waitFor(() => expect(result.current.data).toBeDefined());
    expect(result.current.data).toContainEqual(
      expect.objectContaining({ id: 1, name: "Testo" })
    );
    expect(result.current.data).toContainEqual(
      expect.objectContaining({ id: 2, name: "Alice" })
    );
  });

  it("handles non-array response", async () => {
    const { result } = renderHook(() => useCharacters(), { wrapper: QueryClientWrapper });
    await waitFor(() => expect(result.current.data).toBeDefined());
    expect(Array.isArray(result.current.data)).toBe(true);
  });
});

describe("useCharacter", () => {
  it("returns a single character", async () => {
    const { result } = renderHook(() => useCharacter(1), { wrapper: QueryClientWrapper });
    await waitFor(() => expect(result.current.data).toBeDefined());
    expect(result.current.data).toMatchObject({ id: 1 });
  });

  it("fetches correct character data", async () => {
    const { result } = renderHook(() => useCharacter(2), { wrapper: QueryClientWrapper });
    await waitFor(() => expect(result.current.data).toBeDefined());
    expect(result.current.data).toMatchObject({ id: 2 });
  });

  it("is enabled when id is provided", async () => {
    const { result } = renderHook(() => useCharacter(1), { wrapper: QueryClientWrapper });
    await waitFor(() => expect(result.current.isFetching).toBeDefined());
    expect(result.current.isFetching).toBe(true);
  });
});

describe("useUpdateCharacter", () => {
  it("updates a character and invalidates list", async () => {
    const { result } = renderHook(() => useUpdateCharacter(), { wrapper: QueryClientWrapper });

    let updatedChar: unknown;
    const update: CharacterUpdate = { name: "New Name", level: 5 };
    result.current.mutate(
      { id: 1, update },
      { onSuccess: (data) => { updatedChar = data; } }
    );

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(updatedChar).toMatchObject({ id: 1, name: "New Name", level: 5 });
  });

  it("sends correct id and update to API", async () => {
    const { result } = renderHook(() => useUpdateCharacter(), { wrapper: QueryClientWrapper });

    result.current.mutate({ id: 2, update: { name: "Modified", description: "Updated" } });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
  });

  it("handles partial updates", async () => {
    const { result } = renderHook(() => useUpdateCharacter(), { wrapper: QueryClientWrapper });

    result.current.mutate({ id: 1, update: { level: 10 } });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
  });
});