import { describe, it, expect, vi } from "vitest";
import { screen, fireEvent } from "@testing-library/react";
import { withQuery } from "../../test/wrappers";
import { ItemsIndex } from "./items";

vi.mock("@tanstack/react-router", async () => {
  const actual = await vi.importActual<typeof import("@tanstack/react-router")>("@tanstack/react-router");
  return {
    ...actual,
    useNavigate: () => vi.fn(),
    useLocation: () => ({ pathname: "/items" }),
    createFileRoute: () => () => ({ component: vi.fn() }),
    Link: ({ children }: { children: React.ReactNode }) => <span>{children}</span>,
  };
});

vi.mock("../../contexts/WorldStoreContext", () => ({
  useWorldStore: () => ({ currentWorld: null, setWorld: vi.fn() }),
}));

describe("ItemsIndex", () => {
  it("renders item template list from MSW", async () => {
    withQuery(<ItemsIndex />);
    await screen.findAllByText("Iron Sword");
  });

  it("search input accepts text", async () => {
    withQuery(<ItemsIndex />);
    await screen.findAllByText("Iron Sword");
    const searchInput = screen.getByPlaceholderText(/search items/i);
    fireEvent.change(searchInput, { target: { value: "Iron" } });
    expect(searchInput).toHaveValue("Iron");
  });
});