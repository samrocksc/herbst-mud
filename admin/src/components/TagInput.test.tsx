import { describe, it, expect, vi } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";
import { TagInput } from "./TagInput";

const availableTags = ["fire", "ice", "lightning", "earth", "water"];

describe("TagInput", () => {
  describe("rendering", () => {
    it("renders input field", () => {
      render(
        <TagInput
          value={[]}
          onChange={vi.fn()}
        />,
      );

      expect(screen.getByRole("textbox")).toBeInTheDocument();
    });

    it("renders placeholder", () => {
      render(
        <TagInput
          value={[]}
          onChange={vi.fn()}
          placeholder="Add a tag..."
        />,
      );

      expect(screen.getByPlaceholderText("Add a tag...")).toBeInTheDocument();
    });

    it("renders label when provided", () => {
      render(
        <TagInput
          value={[]}
          onChange={vi.fn()}
          label="Tags"
        />,
      );

      expect(screen.getByText("Tags")).toBeInTheDocument();
    });

    it("renders disabled state", () => {
      render(
        <TagInput
          value={[]}
          onChange={vi.fn()}
          disabled={true}
        />,
      );

      expect(screen.getByRole("textbox")).toBeDisabled();
    });
  });

  describe("selected tags", () => {
    it("renders selected tags", () => {
      render(
        <TagInput
          value={["fire", "ice"]}
          onChange={vi.fn()}
        />,
      );

      expect(screen.getByText("fire")).toBeInTheDocument();
      expect(screen.getByText("ice")).toBeInTheDocument();
    });

    it("renders remove button for each tag", () => {
      render(
        <TagInput
          value={["fire", "ice"]}
          onChange={vi.fn()}
        />,
      );

      const removeButtons = screen.getAllByRole("button", { name: /remove/i });
      expect(removeButtons).toHaveLength(2);
    });

    it("shows nothing when no tags selected", () => {
      render(
        <TagInput
          value={[]}
          onChange={vi.fn()}
        />,
      );

      expect(screen.queryByRole("button", { name: /remove/i })).not.toBeInTheDocument();
    });
  });

  describe("adding tags with Enter", () => {
    it("calls onChange when pressing Enter with new tag", () => {
      const handleChange = vi.fn();
      render(
        <TagInput
          value={[]}
          onChange={handleChange}
          availableTags={availableTags}
        />,
      );

      const input = screen.getByRole("textbox");
      fireEvent.change(input, { target: { value: "lightning" } });
      fireEvent.keyDown(input, { key: "Enter" });

      expect(handleChange).toHaveBeenCalledWith(["lightning"]);
    });

    it("clears input after adding tag", () => {
      const handleChange = vi.fn();
      render(
        <TagInput
          value={[]}
          onChange={handleChange}
          availableTags={availableTags}
        />,
      );

      const input = screen.getByRole("textbox");
      fireEvent.change(input, { target: { value: "test-tag" } });
      fireEvent.keyDown(input, { key: "Enter" });

      // Input should be cleared after adding tag
      expect(input).toHaveValue("");
    });
  });

  describe("removing tags", () => {
    it("removes tag when clicking remove button", () => {
      const handleChange = vi.fn();
      render(
        <TagInput
          value={["fire", "ice"]}
          onChange={handleChange}
          availableTags={availableTags}
        />,
      );

      const removeButton = screen.getByRole("button", { name: "Remove fire" });
      fireEvent.click(removeButton);

      expect(handleChange).toHaveBeenCalledWith(["ice"]);
    });

    it("removes last tag with backspace on empty input", () => {
      const handleChange = vi.fn();
      render(
        <TagInput
          value={["fire", "ice"]}
          onChange={handleChange}
          availableTags={availableTags}
        />,
      );

      const input = screen.getByRole("textbox");
      fireEvent.keyDown(input, { key: "Backspace" });

      expect(handleChange).toHaveBeenCalledWith(["fire"]);
    });
  });

  describe("autocomplete dropdown", () => {
    it("shows dropdown when typing", () => {
      render(
        <TagInput
          value={[]}
          onChange={vi.fn()}
          availableTags={availableTags}
        />,
      );

      const input = screen.getByRole("textbox");
      fireEvent.focus(input);
      fireEvent.change(input, { target: { value: "li" } });

      expect(screen.getByRole("listbox")).toBeInTheDocument();
    });

    it("dropdown contains suggestions", () => {
      render(
        <TagInput
          value={[]}
          onChange={vi.fn()}
          availableTags={availableTags}
        />,
      );

      const input = screen.getByRole("textbox");
      fireEvent.focus(input);
      fireEvent.change(input, { target: { value: "fi" } });

      const options = screen.getAllByRole("option");
      expect(options.length).toBeGreaterThan(0);
    });
  });

  describe("keyboard navigation", () => {
    it("opens dropdown with arrow keys", () => {
      render(
        <TagInput
          value={[]}
          onChange={vi.fn()}
          availableTags={availableTags}
        />,
      );

      const input = screen.getByRole("textbox");
      fireEvent.focus(input);
      fireEvent.change(input, { target: { value: "l" } });

      fireEvent.keyDown(input, { key: "ArrowDown" });

      expect(screen.getByRole("listbox")).toBeInTheDocument();
    });

    it("closes dropdown with Escape", () => {
      render(
        <TagInput
          value={[]}
          onChange={vi.fn()}
          availableTags={availableTags}
        />,
      );

      const input = screen.getByRole("textbox");
      fireEvent.focus(input);
      fireEvent.change(input, { target: { value: "l" } });

      expect(screen.getByRole("listbox")).toBeInTheDocument();

      fireEvent.keyDown(input, { key: "Escape" });

      expect(screen.queryByRole("listbox")).not.toBeInTheDocument();
    });
  });

  describe("click outside", () => {
    it("closes dropdown when clicking outside", () => {
      render(
        <TagInput
          value={[]}
          onChange={vi.fn()}
          availableTags={availableTags}
        />,
      );

      const input = screen.getByRole("textbox");
      fireEvent.focus(input);
      fireEvent.change(input, { target: { value: "l" } });

      expect(screen.getByRole("listbox")).toBeInTheDocument();

      fireEvent.mouseDown(document.body);

      expect(screen.queryByRole("listbox")).not.toBeInTheDocument();
    });
  });

  describe("accessibility", () => {
    it("input has correct attributes", () => {
      render(
        <TagInput
          value={[]}
          onChange={vi.fn()}
        />,
      );

      const input = screen.getByRole("textbox");
      expect(input).toHaveAttribute("type", "text");
      expect(input).toHaveAttribute("autocomplete", "off");
    });
  });
});