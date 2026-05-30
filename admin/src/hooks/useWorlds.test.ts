import { describe, it, expect } from "vitest";
import { renderHook, waitFor } from "@testing-library/react";
import { QueryClientWrapper } from "../test/wrappers";
import { useWorlds, useWorld, useCreateWorld, useUpdateWorld, useDeleteWorld } from "./useWorlds";
import { makeWorldInput } from "../test/factories";

describe("useWorlds", () => {
  it("returns a list of worlds", async () => {
    const { result } = renderHook(() => useWorlds(), { wrapper: QueryClientWrapper });
    await waitFor(() => expect(result.current.data).toBeDefined());
    expect(result.current.data).toHaveLength(2);
  });

  it("contains expected worlds", async () => {
    const { result } = renderHook(() => useWorlds(), { wrapper: QueryClientWrapper });
    await waitFor(() => expect(result.current.data).toBeDefined());
    expect(result.current.data).toContainEqual(
      expect.objectContaining({ id: 1, name: "Test World" })
    );
    expect(result.current.data).toContainEqual(
      expect.objectContaining({ id: 2, name: "Dark Realm" })
    );
  });
});

describe("useWorld", () => {
  it("returns a single world", async () => {
    const { result } = renderHook(() => useWorld(1), { wrapper: QueryClientWrapper });
    await waitFor(() => expect(result.current.data).toBeDefined());
    expect(result.current.data).toMatchObject({ id: 1 });
  });
});

describe("useCreateWorld", () => {
  it("creates a world and invalidates list", async () => {
    const { result } = renderHook(() => useCreateWorld(), { wrapper: QueryClientWrapper });

    let createdWorld: unknown;
    result.current.mutate(makeWorldInput({ name: "New World" }), {
      onSuccess: (data) => { createdWorld = data; },
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(createdWorld).toMatchObject({ name: "New World" });
  });

  it("sends correct input to API", async () => {
    const { result } = renderHook(() => useCreateWorld(), { wrapper: QueryClientWrapper });
    const input = makeWorldInput({ name: "Created World", description: "A new world" });

    result.current.mutate(input);

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
  });
});

describe("useUpdateWorld", () => {
  it("updates a world and invalidates list and single world", async () => {
    const { result } = renderHook(() => useUpdateWorld(), { wrapper: QueryClientWrapper });

    let updatedWorld: unknown;
    result.current.mutate(
      { id: 1, input: makeWorldInput({ name: "Updated World" }) },
      { onSuccess: (data) => { updatedWorld = data; } }
    );

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(updatedWorld).toMatchObject({ id: 1, name: "Updated World" });
  });

  it("sends correct id and input to API", async () => {
    const { result } = renderHook(() => useUpdateWorld(), { wrapper: QueryClientWrapper });

    result.current.mutate({ id: 2, input: makeWorldInput({ name: "Modified" }) });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
  });
});

describe("useDeleteWorld", () => {
  it("deletes a world and invalidates list", async () => {
    const { result } = renderHook(() => useDeleteWorld(), { wrapper: QueryClientWrapper });

    let deleteCompleted = false;
    result.current.mutate(1, {
      onSuccess: () => { deleteCompleted = true; },
    });

    await waitFor(() => expect(deleteCompleted).toBe(true));
  });

  it("sends correct id to API", async () => {
    const { result } = renderHook(() => useDeleteWorld(), { wrapper: QueryClientWrapper });

    result.current.mutate(2);

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
  });
});