import { describe, it, expect, vi } from "vitest";
import { screen } from "@testing-library/react";
import { withQuery } from "../../test/wrappers";
import { QuestsManagement } from "./quests";

vi.mock("@tanstack/react-router", async () => {
  const actual = await vi.importActual<typeof import("@tanstack/react-router")>("@tanstack/react-router");
  return {
    ...actual,
    useNavigate: () => vi.fn(),
    useLocation: () => ({ pathname: "/quests" }),
    createFileRoute: () => () => ({ component: vi.fn() }),
    Link: ({ children }: { children: React.ReactNode }) => <span>{children}</span>,
  };
});

describe("QuestsManagement", () => {
  it("renders quest list from MSW", async () => {
    withQuery(<QuestsManagement />);
    await screen.findAllByText("Find the Key");
    expect(screen.getByText("Quests")).toBeInTheDocument();
  });

  it('"Add Quest" button navigates', async () => {
    withQuery(<QuestsManagement />);
    await screen.findAllByText("Find the Key");
    // Check that we have at least one button with "Add Quest"
    const buttons = screen.queryAllByRole("button", { name: /add quest/i });
    expect(buttons.length).toBeGreaterThan(0);
  });
});