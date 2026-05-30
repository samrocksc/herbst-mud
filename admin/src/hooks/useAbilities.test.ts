import { describe, it, expect } from "vitest";
import { renderHook, waitFor } from "@testing-library/react";
import { QueryClientWrapper } from "../test/wrappers";
import { useAbilities, useAbility, useCreateAbility, useUpdateAbility, useDeleteAbility } from "./useAbilities";
import { makeAbilityInput } from "../test/factories";

describe("useAbilities", () => {
  it("returns a list", async () => {
    const { result } = renderHook(() => useAbilities(), { wrapper: QueryClientWrapper });
    await waitFor(() => expect(result.current.data).toBeDefined());
    expect(result.current.data).toHaveLength(2);
  });

  it("returns abilities with expected shape", async () => {
    const { result } = renderHook(() => useAbilities(), { wrapper: QueryClientWrapper });
    await waitFor(() => expect(result.current.data).toBeDefined());
    const ability = result.current.data![0];
    expect(ability).toHaveProperty("id");
    expect(ability).toHaveProperty("name");
    expect(ability).toHaveProperty("ability_type");
  });

  it("passes filters to query key", async () => {
    const { result } = renderHook(() => useAbilities({ type: "magic", abilityClass: "active" }), { wrapper: QueryClientWrapper });
    await waitFor(() => expect(result.current.data).toBeDefined());
    expect(result.current.isSuccess).toBe(true);
  });
});

describe("useAbility", () => {
  it("returns single ability by id", async () => {
    const { result } = renderHook(() => useAbility(1), { wrapper: QueryClientWrapper });
    await waitFor(() => expect(result.current.data).toBeDefined());
    expect(result.current.data!.id).toBe(1);
  });

  it("returns undefined when id is null", () => {
    const { result } = renderHook(() => useAbility(null), { wrapper: QueryClientWrapper });
    expect(result.current.data).toBeUndefined();
  });

  it("does not fetch when id is null", () => {
    const { result } = renderHook(() => useAbility(null), { wrapper: QueryClientWrapper });
    expect(result.current.fetchStatus).toBe("idle");
  });
});

describe("useCreateAbility", () => {
  it("creates ability and invalidates list", async () => {
    const { result } = renderHook(() => useCreateAbility(), { wrapper: QueryClientWrapper });
    const input = makeAbilityInput({ name: "New Ability" });

    result.current.mutate(input);

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
  });

  it("returns created ability", async () => {
    const { result } = renderHook(() => useCreateAbility(), { wrapper: QueryClientWrapper });
    const input = makeAbilityInput({ name: "Thunder Strike" });

    result.current.mutate(input);

    await waitFor(() => expect(result.current.data).toBeDefined());
    expect(result.current.data!.name).toBe("Thunder Strike");
  });
});

describe("useUpdateAbility", () => {
  it("updates ability and invalidates list", async () => {
    const { result } = renderHook(() => useUpdateAbility(), { wrapper: QueryClientWrapper });
    const input = makeAbilityInput({ name: "Updated Fireball" });

    result.current.mutate({ id: 1, input });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
  });

  it("returns updated ability", async () => {
    const { result } = renderHook(() => useUpdateAbility(), { wrapper: QueryClientWrapper });
    const input = makeAbilityInput({ name: "Mega Fireball" });

    result.current.mutate({ id: 1, input });

    await waitFor(() => expect(result.current.data).toBeDefined());
    expect(result.current.data!.name).toBe("Mega Fireball");
  });
});

describe("useDeleteAbility", () => {
  it("deletes ability", async () => {
    const { result } = renderHook(() => useDeleteAbility(), { wrapper: QueryClientWrapper });

    result.current.mutate(1);

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
  });
});