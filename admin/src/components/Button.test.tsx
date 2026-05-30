import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import { Button } from "./Button";

describe("Button", () => {
  describe("rendering", () => {
    it("renders button with text", () => {
      render(<Button>Click me</Button>);

      expect(screen.getByText("Click me")).toBeInTheDocument();
    });

    it("renders as button element", () => {
      render(<Button>Test</Button>);

      const button = screen.getByRole("button", { name: "Test" });
      expect(button.tagName).toBe("BUTTON");
    });
  });

  describe("variants", () => {
    it("renders primary variant by default", () => {
      render(<Button>Primary</Button>);

      const button = screen.getByRole("button", { name: "Primary" });
      expect(button).toHaveClass("bg-primary");
    });

    it("renders secondary variant", () => {
      render(<Button variant="secondary">Secondary</Button>);

      const button = screen.getByRole("button", { name: "Secondary" });
      expect(button).toHaveClass("bg-transparent");
    });

    it("renders danger variant", () => {
      render(<Button variant="danger">Delete</Button>);

      const button = screen.getByRole("button", { name: "Delete" });
      expect(button).toHaveClass("bg-danger");
    });

    it("renders ghost variant", () => {
      render(<Button variant="ghost">Ghost</Button>);

      const button = screen.getByRole("button", { name: "Ghost" });
      expect(button).toHaveClass("bg-transparent");
    });

    it("renders outline variant", () => {
      render(<Button variant="outline">Outline</Button>);

      const button = screen.getByRole("button", { name: "Outline" });
      expect(button).toHaveClass("border-border");
    });

    it("renders success variant", () => {
      render(<Button variant="success">Save</Button>);

      const button = screen.getByRole("button", { name: "Save" });
      expect(button).toHaveClass("bg-success");
    });

    it("renders accent variant", () => {
      render(<Button variant="accent">Accent</Button>);

      const button = screen.getByRole("button", { name: "Accent" });
      expect(button).toHaveClass("bg-accent");
    });
  });

  describe("sizes", () => {
    it("renders small size", () => {
      render(<Button size="sm">Small</Button>);

      const button = screen.getByRole("button", { name: "Small" });
      expect(button).toHaveClass("px-2.5", "py-1", "text-xs");
    });

    it("renders medium size by default", () => {
      render(<Button size="md">Medium</Button>);

      const button = screen.getByRole("button", { name: "Medium" });
      expect(button).toHaveClass("px-4", "py-2", "text-sm");
    });

    it("renders large size", () => {
      render(<Button size="lg">Large</Button>);

      const button = screen.getByRole("button", { name: "Large" });
      expect(button).toHaveClass("px-6", "py-3", "text-base");
    });
  });

  describe("disabled state", () => {
    it("renders disabled button", () => {
      render(<Button disabled>Disabled</Button>);

      const button = screen.getByRole("button", { name: "Disabled" });
      expect(button).toBeDisabled();
    });

    it("applies disabled opacity styles", () => {
      render(<Button disabled>Disabled</Button>);

      const button = screen.getByRole("button", { name: "Disabled" });
      expect(button).toHaveClass("disabled:opacity-50", "disabled:cursor-not-allowed");
    });
  });

  describe("full width", () => {
    it("applies full width class", () => {
      render(<Button fullWidth>Full Width</Button>);

      const button = screen.getByRole("button", { name: "Full Width" });
      expect(button).toHaveClass("w-full");
    });

    it("does not apply full width by default", () => {
      render(<Button>Normal</Button>);

      const button = screen.getByRole("button", { name: "Normal" });
      expect(button).not.toHaveClass("w-full");
    });
  });

  describe("custom className", () => {
    it("merges custom className", () => {
      render(<Button className="custom-class">Custom</Button>);

      const button = screen.getByRole("button", { name: "Custom" });
      expect(button).toHaveClass("custom-class");
    });
  });

  describe("HTML attributes", () => {
    it("forwards type attribute when specified", () => {
      render(<Button type="submit">Submit</Button>);

      const button = screen.getByRole("button", { name: "Submit" });
      expect(button).toHaveAttribute("type", "submit");
    });

    it("forwards onClick handler", () => {
      const handleClick = vi.fn();
      render(<Button onClick={handleClick}>Click</Button>);

      const button = screen.getByRole("button", { name: "Click" });
      button.click();

      expect(handleClick).toHaveBeenCalledTimes(1);
    });

    it("forwards custom data attributes", () => {
      render(<Button data-testid="custom-button">Custom</Button>);

      expect(screen.getByTestId("custom-button")).toBeInTheDocument();
    });
  });

  describe("accessibility", () => {
    it("has correct focus attributes", () => {
      render(<Button>Focusable</Button>);

      const button = screen.getByRole("button", { name: "Focusable" });
      expect(button).toHaveClass("focus-visible:outline-none", "focus-visible:ring-2");
    });
  });
});