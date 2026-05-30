import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { CreateAbilityPage } from "./abilities.new";

// Mock router hooks
vi.mock("@tanstack/react-router", async () => {
  const actual = await vi.importActual("@tanstack/react-router");
  return {
    ...actual,
    useNavigate: () => vi.fn(),
    createFileRoute: () => () => ({ component: vi.fn() }),
    useRouter: () => ({ navigate: vi.fn() }),
    Link: ({ children }: { children: React.ReactNode }) => <span>{children}</span>,
  };
});

// Mock the hooks
vi.mock("../../hooks/useAbilities", () => ({
  useCreateAbility: () => ({
    mutateAsync: vi.fn().mockResolvedValue({ id: 1, name: "Test Ability" }),
    isPending: false,
  }),
  useTags: () => ({
    data: [],
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

describe("CreateAbilityPage", () => {
  it("renders form", async () => {
    renderWithQueryClient(<CreateAbilityPage />);

    // Check for the page title specifically
    const title = screen.getByRole('heading', { name: /Create Ability/i });
    expect(title).toBeInTheDocument();

    expect(screen.getByLabelText("Name *")).toBeInTheDocument();

    // Check for the submit button specifically
    const submitButton = screen.getByRole('button', { name: /Create Ability/i });
    expect(submitButton).toBeInTheDocument();
  });
});