import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import { PageHeader } from "./PageHeader";
import { Route as RouterContext } from "@tanstack/react-router";
import type { ReactNode } from "react";

vi.mock("@tanstack/react-router", async () => {
  const actual = await vi.importActual("@tanstack/react-router");
  return {
    ...(actual as object),
    Link: ({ to, children, ...props }: { to: string; children: ReactNode }) => (
      <a href={`#${to}`} {...props}>{children}</a>
    ),
  };
});

describe("PageHeader", () => {
  describe("title", () => {
    it("renders the title", () => {
      render(<PageHeader title="Abilities" />);

      expect(screen.getByText("Abilities")).toBeInTheDocument();
    });

    it("renders title with h1 tag", () => {
      render(<PageHeader title="My Title" />);

      const h1 = screen.getByText("My Title");
      expect(h1.tagName).toBe("H1");
    });
  });

  describe("back button", () => {
    it("renders back button when showBack and backTo are provided", () => {
      render(<PageHeader title="Edit Ability" showBack backTo="/abilities" />);

      expect(screen.getByText("← Dashboard")).toBeInTheDocument();
    });

    it("renders back button with custom label", () => {
      render(<PageHeader title="Edit" showBack backTo="/abilities" backLabel="← Abilities" />);

      expect(screen.getByText("← Abilities")).toBeInTheDocument();
    });

    it("does not render back button when showBack is false", () => {
      render(<PageHeader title="Abilities" showBack={false} backTo="/abilities" />);

      expect(screen.queryByText("← Dashboard")).not.toBeInTheDocument();
    });

    it("does not render back button when backTo is not provided", () => {
      render(<PageHeader title="Abilities" showBack />);

      expect(screen.queryByRole("link", { name: /dashboard/i })).not.toBeInTheDocument();
    });
  });

  describe("actions", () => {
    it("renders actions when provided", () => {
      render(
        <PageHeader
          title="Abilities"
          actions={<button type="button">New Ability</button>}
        />,
      );

      expect(screen.getByRole("button", { name: "New Ability" })).toBeInTheDocument();
    });

    it("renders multiple actions", () => {
      render(
        <PageHeader
          title="Abilities"
          actions={
            <>
              <button type="button">Edit</button>
              <button type="button">Delete</button>
            </>
          }
        />,
      );

      expect(screen.getByRole("button", { name: "Edit" })).toBeInTheDocument();
      expect(screen.getByRole("button", { name: "Delete" })).toBeInTheDocument();
    });

    it("does not render actions div when no actions", () => {
      render(<PageHeader title="Abilities" />);

      const container = screen.getByText("Abilities").parentElement;
      expect(container?.children.length).toBeLessThanOrEqual(1);
    });
  });

  describe("styling", () => {
    it("applies heading styles", () => {
      render(<PageHeader title="Test" />);

      const h1 = screen.getByText("Test");
      expect(h1).toHaveClass("text-lg", "font-bold");
    });

    it("truncates long titles", () => {
      render(<PageHeader title="A Very Long Title That Should Be Truncated" />);

      const h1 = screen.getByText("A Very Long Title That Should Be Truncated");
      expect(h1).toHaveClass("truncate");
    });
  });
});