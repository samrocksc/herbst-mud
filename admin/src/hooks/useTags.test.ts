import { describe, it, expect } from "vitest";
import { renderHook, waitFor } from "@testing-library/react";
import { QueryClientWrapper } from "../test/wrappers";
import { useTags, useCreateTag, useUpdateTag } from "./useTags";
import { makeTagInput } from "../test/factories";

describe("useTags", () => {
  it("returns a list", async () => {
    const { result } = renderHook(() => useTags(), { wrapper: QueryClientWrapper });
    await waitFor(() => expect(result.current.data).toBeDefined());
    expect(result.current.data).toHaveLength(2);
  });

  it("contains expected tags", async () => {
    const { result } = renderHook(() => useTags(), { wrapper: QueryClientWrapper });
    await waitFor(() => expect(result.current.data).toBeDefined());
    expect(result.current.data).toContainEqual(
      expect.objectContaining({ id: 1, name: "common" })
    );
    expect(result.current.data).toContainEqual(
      expect.objectContaining({ id: 2, name: "rare" })
    );
  });
});

describe("useCreateTag", () => {
  it("creates a tag and invalidates list", async () => {
    const { result } = renderHook(() => useCreateTag(), { wrapper: QueryClientWrapper });

    let createdTag: unknown;
    result.current.mutate(makeTagInput({ name: "epic", color: "#ff0000" }), {
      onSuccess: (data) => { createdTag = data; },
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(createdTag).toMatchObject({ name: "epic", color: "#ff0000" });
  });

  it("sends correct input to API", async () => {
    const { result } = renderHook(() => useCreateTag(), { wrapper: QueryClientWrapper });
    const input = makeTagInput({ name: "legendary", color: "#ffd700" });

    result.current.mutate(input);

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
  });
});

describe("useUpdateTag", () => {
  it("updates a tag and invalidates list", async () => {
    const { result } = renderHook(() => useUpdateTag(), { wrapper: QueryClientWrapper });

    let updatedTag: unknown;
    result.current.mutate(
      { id: 1, input: { name: "updated-name", color: "#00ff00" } },
      { onSuccess: (data) => { updatedTag = data; } }
    );

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(updatedTag).toMatchObject({ id: 1, name: "updated-name", color: "#00ff00" });
  });

  it("sends correct id and input to API", async () => {
    const { result } = renderHook(() => useUpdateTag(), { wrapper: QueryClientWrapper });

    result.current.mutate({ id: 2, input: { name: "modified", color: "#0000ff" } });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
  });
});