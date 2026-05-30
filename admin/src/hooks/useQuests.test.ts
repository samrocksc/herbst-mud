import { describe, it, expect } from "vitest";
import { renderHook, waitFor } from "@testing-library/react";
import { QueryClientWrapper } from "../test/wrappers";
import {
  useQuests,
  useQuest,
  useQuestLookups,
  useCreateQuest,
  useUpdateQuest,
  useDeleteQuest,
} from "./useQuests";
import { makeQuestInput } from "../test/factories";

describe("useQuests", () => {
  it("returns a list", async () => {
    const { result } = renderHook(() => useQuests(), { wrapper: QueryClientWrapper });
    await waitFor(() => expect(result.current.data).toBeDefined());
    expect(result.current.data).toHaveLength(2);
  });

  it("returns quests with expected shape", async () => {
    const { result } = renderHook(() => useQuests(), { wrapper: QueryClientWrapper });
    await waitFor(() => expect(result.current.data).toBeDefined());
    const quest = result.current.data![0];
    expect(quest).toHaveProperty("id");
    expect(quest).toHaveProperty("name");
    expect(quest).toHaveProperty("objectives");
    expect(quest).toHaveProperty("rewards");
  });
});

describe("useQuest", () => {
  it("returns single quest by id", async () => {
    const { result } = renderHook(() => useQuest(1), { wrapper: QueryClientWrapper });
    await waitFor(() => expect(result.current.data).toBeDefined());
    expect(result.current.data!.id).toBe(1);
  });

  it("returns undefined when id is null", () => {
    const { result } = renderHook(() => useQuest(null), { wrapper: QueryClientWrapper });
    expect(result.current.data).toBeUndefined();
  });
});

describe("useQuestLookups", () => {
  it("returns lookup data with quest types", async () => {
    const { result } = renderHook(() => useQuestLookups(), { wrapper: QueryClientWrapper });
    await waitFor(() => {
      expect(result.current.data).toBeDefined();
    });
    expect(result.current.data).toHaveProperty("quest_types");
  });

  it("includes expected lookup categories", async () => {
    const { result } = renderHook(() => useQuestLookups(), { wrapper: QueryClientWrapper });
    await waitFor(() => {
      expect(result.current.data).toBeDefined();
    });
    expect(result.current.data).toHaveProperty("npcs");
    expect(result.current.data).toHaveProperty("rooms");
    expect(result.current.data).toHaveProperty("items");
    expect(result.current.data).toHaveProperty("effects");
  });
});

describe("useCreateQuest", () => {
  it("creates quest and invalidates list", async () => {
    const { result } = renderHook(() => useCreateQuest(), { wrapper: QueryClientWrapper });
    const input = makeQuestInput({ name: "New Quest" });

    result.current.mutate(input);

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
  });

  it("returns created quest", async () => {
    const { result } = renderHook(() => useCreateQuest(), { wrapper: QueryClientWrapper });
    const input = makeQuestInput({ name: "Collect 10 Herbs" });

    result.current.mutate(input);

    await waitFor(() => expect(result.current.data).toBeDefined());
    expect(result.current.data!.name).toBe("Collect 10 Herbs");
  });
});

describe("useUpdateQuest", () => {
  it("updates quest and invalidates list", async () => {
    const { result } = renderHook(() => useUpdateQuest(), { wrapper: QueryClientWrapper });
    const input = makeQuestInput({ name: "Updated Quest" });

    result.current.mutate({ id: 1, input });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
  });

  it("returns updated quest", async () => {
    const { result } = renderHook(() => useUpdateQuest(), { wrapper: QueryClientWrapper });
    const input = makeQuestInput({ name: "Escort the Merchant" });

    result.current.mutate({ id: 1, input });

    await waitFor(() => expect(result.current.data).toBeDefined());
    expect(result.current.data!.name).toBe("Escort the Merchant");
  });
});

describe("useDeleteQuest", () => {
  it("deletes quest", async () => {
    const { result } = renderHook(() => useDeleteQuest(), { wrapper: QueryClientWrapper });

    result.current.mutate(1);

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
  });
});