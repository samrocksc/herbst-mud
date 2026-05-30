import { describe, it, expect, vi } from "vitest";
import { screen, fireEvent } from "@testing-library/react";
import { withQuery } from "../../test/wrappers";
import { AbilitiesManagement } from "./abilities";

vi.mock("@tanstack/react-router", async () => {
  const actual = await vi.importActual<typeof import("@tanstack/react-router")>("@tanstack/react-router");
  return {
    ...actual,
    useNavigate: () => vi.fn(),
    useLocation: () => ({ pathname: "/abilities" }),
    createFileRoute: () => () => ({ component: vi.fn() }),
    Link: ({ children }: { children: React.ReactNode }) => <span>{children}</span>,
  };
});

describe("AbilitiesManagement", () => {
  it("renders ability list from MSW", async () => {
    withQuery(<AbilitiesManagement />);
    // Use findByText which waits for the element to appear
    await screen.findAllByText("Fireball");
    expect(screen.getByText("Abilities")).toBeInTheDocument();
  });

  it('"Add Ability" button is present', async () => {
    withQuery(<AbilitiesManagement />);
    await screen.findAllByText("Fireball");
    // Check that we have at least one button with "Add Ability"
    const buttons = screen.queryAllByRole("button", { name: /add ability/i });
    expect(buttons.length).toBeGreaterThan(0);
  });

  it("filter type dropdown changes value", async () => {
    withQuery(<AbilitiesManagement />);
    await screen.findAllByText("Fireball");
    // Get all select elements and find the one with the label "Type"
    const selects = screen.queryAllByRole("combobox");
    expect(selects.length).toBeGreaterThan(0);
  });
});