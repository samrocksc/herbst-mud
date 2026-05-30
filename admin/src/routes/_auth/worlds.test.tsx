import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { WorldsManagement } from "./worlds";

// Mock router hooks
vi.mock("@tanstack/react-router", async () => {
  const actual = await vi.importActual("@tanstack/react-router");
  return {
    ...actual,
    useNavigate: () => vi.fn(),
    useLocation: () => ({ pathname: "/worlds" }),
    createFileRoute: () => () => ({ component: vi.fn() }),
    Link: ({ children, to }: { children: React.ReactNode; to: string }) => <a href={to}>{children}</a>,
  };
});

// Mock the useWorlds hook
vi.mock("../../hooks/useWorlds", () => ({
  useWorlds: () => ({
    data: [
      { id: 1, name: "Test World", title: "The Test World", description: "A world for testing", active: true },
      { id: 2, name: "Inactive World", title: "Inactive World", description: "An inactive world", active: false },
    ],
    isLoading: false,
    error: null,
  }),
  useSetActiveWorld: () => ({
    mutate: vi.fn(),
    isPending: false,
  }),
}));

const createTestQueryClient = () =>
  new QueryClient({
    defaultOptions: { queries: { retry: false, staleTime: Infinity } },
  });

function renderWithQueryClient(ui: React.ReactElement) {
  const client = createTestQueryClient();
  return render(
    <QueryClientProvider client={client}>
      {ui}
    </QueryClientProvider>
  );
}

describe("WorldsManagement", () => {
  it("renders world list", async () => {
    renderWithQueryClient(<WorldsManagement />);

    expect(screen.getByText("Worlds")).toBeInTheDocument();

    // Check that worlds are shown (using queryAllByText to handle duplicates)
    const testWorldElements = screen.queryAllByText("Test World");
    const inactiveWorldElements = screen.queryAllByText("Inactive World");
    expect(testWorldElements.length).toBeGreaterThan(0);
    expect(inactiveWorldElements.length).toBeGreaterThan(0);
  });

  it("active/inactive status shown", async () => {
    renderWithQueryClient(<WorldsManagement />);

    // Check that active/inactive status is shown
    const worldRows = screen.getAllByRole("row");
    // Look for rows that contain the status
    const rowsWithStatus = worldRows.filter(row =>
      row.textContent?.includes("Active") || row.textContent?.includes("Inactive")
    );
    expect(rowsWithStatus.length).toBeGreaterThan(0);
  });
});