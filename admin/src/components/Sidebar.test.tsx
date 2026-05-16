/* eslint-disable functional/prefer-immutable-types */
import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";

// Mock the TanStack router hooks
vi.mock("@tanstack/react-router", () => ({
  Link: ({ children, to, className, title }: { children: React.ReactNode; to: string; className?: string; title?: string }) => (
    <a href={to} className={className} title={title}>
      {children}
    </a>
  ),
  useLocation: () => ({ pathname: "/dashboard" }),
}));

describe("Sidebar", () => {
  beforeEach(() => {
    localStorage.clear();
  });

  it("renders all navigation items", async () => {
    const { Sidebar } = await import("./Sidebar");
    render(<Sidebar />);

    expect(screen.getByText("Dashboard")).toBeInTheDocument();
    expect(screen.getByText("Items")).toBeInTheDocument();
    expect(screen.getByText("Skills")).toBeInTheDocument();
    expect(screen.getByText("Map")).toBeInTheDocument();
    expect(screen.getByText("NPCs")).toBeInTheDocument();
    expect(screen.getByText("Players")).toBeInTheDocument();
  });

  it("renders the header with title", async () => {
    const { Sidebar } = await import("./Sidebar");
    render(<Sidebar />);

    expect(screen.getByText("Herbst MUD")).toBeInTheDocument();
  });

  it("renders a collapse toggle button with correct aria-label", async () => {
    const { Sidebar } = await import("./Sidebar");
    render(<Sidebar />);

    const toggleButton = screen.getByRole("button", { name: "Collapse sidebar" });
    expect(toggleButton).toBeInTheDocument();
  });

  it("toggles to collapsed state when collapse button is clicked", async () => {
    const { Sidebar } = await import("./Sidebar");
    render(<Sidebar />);

    const toggleButton = screen.getByRole("button", { name: "Collapse sidebar" });
    fireEvent.click(toggleButton);

    // After collapsing, the button label should change to "Expand sidebar"
    expect(screen.getByRole("button", { name: "Expand sidebar" })).toBeInTheDocument();
  });

  it("persists collapsed state to localStorage", async () => {
    const { Sidebar } = await import("./Sidebar");
    render(<Sidebar />);

    const toggleButton = screen.getByRole("button", { name: "Collapse sidebar" });
    fireEvent.click(toggleButton);

    expect(localStorage.getItem("sidebar-collapsed")).toBe("true");
  });

  it("reads initial collapsed state from localStorage", async () => {
    localStorage.setItem("sidebar-collapsed", "true");

    const { Sidebar } = await import("./Sidebar");
    render(<Sidebar />);

    // If initially collapsed, the button should say "Expand sidebar"
    expect(screen.getByRole("button", { name: "Expand sidebar" })).toBeInTheDocument();
  });

  it("defaults to expanded when localStorage is empty", async () => {
    const { Sidebar } = await import("./Sidebar");
    render(<Sidebar />);

    expect(screen.getByRole("button", { name: "Collapse sidebar" })).toBeInTheDocument();
  });

  it("shows tooltips on nav items when collapsed", async () => {
    localStorage.setItem("sidebar-collapsed", "true");

    const { Sidebar } = await import("./Sidebar");
    render(<Sidebar />);

    // When collapsed, nav items should have title attributes
    const dashboardLink = screen.getByText("Dashboard").closest("a");
    expect(dashboardLink).toHaveAttribute("title", "Dashboard");
  });

  it("hides tooltips on nav items when expanded", async () => {
    const { Sidebar } = await import("./Sidebar");
    render(<Sidebar />);

    // When expanded, nav items should NOT have title attributes
    const dashboardLink = screen.getByText("Dashboard").closest("a");
    expect(dashboardLink).not.toHaveAttribute("title");
  });
});