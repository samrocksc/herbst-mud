import { describe, it, expect } from "vitest";
import { renderHook, waitFor } from "@testing-library/react";
import { QueryClientWrapper } from "../test/wrappers";
import {
  useEquipmentTemplates,
  useEquipmentTemplate,
  useCreateTemplate,
  useUpdateTemplate,
  useDeleteTemplate,
} from "./useEquipmentTemplates";
import { makeEquipmentTemplateInput } from "../test/factories";

describe("useEquipmentTemplates", () => {
  it("returns a list", async () => {
    const { result } = renderHook(() => useEquipmentTemplates(), { wrapper: QueryClientWrapper });
    await waitFor(() => expect(result.current.data).toBeDefined());
    expect(result.current.data).toHaveLength(2);
  });

  it("returns templates with expected shape", async () => {
    const { result } = renderHook(() => useEquipmentTemplates(), { wrapper: QueryClientWrapper });
    await waitFor(() => expect(result.current.data).toBeDefined());
    const template = result.current.data![0];
    expect(template).toHaveProperty("id");
    expect(template).toHaveProperty("name");
    expect(template).toHaveProperty("slot");
  });
});

describe("useEquipmentTemplate", () => {
  it("returns single template by id", async () => {
    const { result } = renderHook(() => useEquipmentTemplate(1), { wrapper: QueryClientWrapper });
    await waitFor(() => expect(result.current.data).toBeDefined());
    expect(result.current.data!.id).toBe(1);
  });

  it("returns null when id is null", async () => {
    const { result } = renderHook(() => useEquipmentTemplate(null), { wrapper: QueryClientWrapper });
    expect(result.current.data).toBeUndefined();
  });
});

describe("useCreateTemplate", () => {
  it("creates template and invalidates list", async () => {
    const { result } = renderHook(() => useCreateTemplate(), { wrapper: QueryClientWrapper });
    const input = makeEquipmentTemplateInput({ name: "Steel Shield" });

    result.current.mutate(input);

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
  });

  it("returns created template", async () => {
    const { result } = renderHook(() => useCreateTemplate(), { wrapper: QueryClientWrapper });
    const input = makeEquipmentTemplateInput({ name: "Magic Staff" });

    result.current.mutate(input);

    await waitFor(() => expect(result.current.data).toBeDefined());
    expect(result.current.data!.name).toBe("Magic Staff");
  });
});

describe("useUpdateTemplate", () => {
  it("updates template and invalidates list", async () => {
    const { result } = renderHook(() => useUpdateTemplate(), { wrapper: QueryClientWrapper });
    const input = makeEquipmentTemplateInput({ name: "Upgraded Sword" });

    result.current.mutate({ id: 1, input });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
  });

  it("returns updated template", async () => {
    const { result } = renderHook(() => useUpdateTemplate(), { wrapper: QueryClientWrapper });
    const input = makeEquipmentTemplateInput({ name: "Enchanted Dagger" });

    result.current.mutate({ id: 1, input });

    await waitFor(() => expect(result.current.data).toBeDefined());
    expect(result.current.data!.name).toBe("Enchanted Dagger");
  });
});

describe("useDeleteTemplate", () => {
  it("deletes template", async () => {
    const { result } = renderHook(() => useDeleteTemplate(), { wrapper: QueryClientWrapper });

    result.current.mutate(1);

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
  });
});