import { describe, it, expect, vi } from "vitest";
import { screen, fireEvent } from "@testing-library/react";
import { withQuery } from "../../test/wrappers";
import { NPCTemplatesIndex } from "./npcs";

vi.mock("@tanstack/react-router", async () => {
  const actual = await vi.importActual<typeof import("@tanstack/react-router")>("@tanstack/react-router");
  return {
    ...actual,
    useLocation: () => ({ pathname: "/npcs" }),
    createFileRoute: () => () => ({ component: vi.fn() }),
    Link: ({ children }: { children: React.ReactNode }) => <span>{children}</span>,
  };
});

vi.mock("../../contexts/WorldStoreContext", () => ({
  useWorldStore: () => ({ currentWorld: null }),
}));

describe("NPCTemplatesIndex", () => {
  it("renders NPC template list from MSW", async () => {
    withQuery(<NPCTemplatesIndex />);
    await screen.findAllByText("Village Guard");
  });

  it("search input accepts text", async () => {
    withQuery(<NPCTemplatesIndex />);
    await screen.findAllByText("Village Guard");
    const searchInput = screen.getByPlaceholderText(/search templates/i);
    fireEvent.change(searchInput, { target: { value: "Village" } });
    expect(searchInput).toHaveValue("Village");
  });
});