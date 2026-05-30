import { describe, it, expect, vi } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";
import { SearchableSelect } from "./SearchableSelect";
import type { SearchableSelectOption } from "./SearchableSelect";

const baseOptions: SearchableSelectOption[] = [
  { id: "1", name: "Fireball" },
  { id: "2", name: "Icebolt" },
  { id: "3", name: "Lightning Bolt" },
];

describe("SearchableSelect", () => {
  const getInput = () => screen.getByRole("textbox", { name: /search/i });

  describe("rendering", () => {
    it("renders input field", () => {
      render(
        <SearchableSelect
          options={baseOptions}
          value=""
          onChange={vi.fn()}
        />,
      );

      expect(getInput()).toBeInTheDocument();
    });

    it("renders placeholder", () => {
      render(
        <SearchableSelect
          options={baseOptions}
          value=""
          onChange={vi.fn()}
          placeholder="Search abilities..."
        />,
      );

      expect(screen.getByPlaceholderText("Search abilities...")).toBeInTheDocument();
    });

    it("renders default placeholder", () => {
      render(
        <SearchableSelect
          options={baseOptions}
          value=""
          onChange={vi.fn()}
        />,
      );

      expect(screen.getByPlaceholderText("Search...")).toBeInTheDocument();
    });

    it("renders label when provided", () => {
      render(
        <SearchableSelect
          options={baseOptions}
          value=""
          onChange={vi.fn()}
          label="Select Ability"
        />,
      );

      expect(screen.getByText("Select Ability")).toBeInTheDocument();
    });

    it("renders disabled state", () => {
      render(
        <SearchableSelect
          options={baseOptions}
          value=""
          onChange={vi.fn()}
          disabled={true}
        />,
      );

      expect(getInput()).toBeDisabled();
    });
  });

  describe("display value", () => {
    it("shows selected option label", () => {
      render(
        <SearchableSelect
          options={baseOptions}
          value="1"
          onChange={vi.fn()}
        />,
      );

      expect(screen.getByDisplayValue("Fireball (1)")).toBeInTheDocument();
    });

    it("shows (not selected) for empty value", () => {
      render(
        <SearchableSelect
          options={baseOptions}
          value=""
          onChange={vi.fn()}
        />,
      );

      expect(screen.getByDisplayValue("(not selected)")).toBeInTheDocument();
    });
  });

  describe("dropdown opening", () => {
    it("opens dropdown on focus", () => {
      render(
        <SearchableSelect
          options={baseOptions}
          value=""
          onChange={vi.fn()}
        />,
      );

      fireEvent.focus(getInput());

      expect(screen.getByText("Fireball (1)")).toBeInTheDocument();
    });

    it("shows all options when opened", () => {
      render(
        <SearchableSelect
          options={baseOptions}
          value=""
          onChange={vi.fn()}
        />,
      );

      fireEvent.focus(getInput());

      expect(screen.getByText("Fireball (1)")).toBeInTheDocument();
      expect(screen.getByText("Icebolt (2)")).toBeInTheDocument();
      expect(screen.getByText("Lightning Bolt (3)")).toBeInTheDocument();
    });
  });

  describe("search filtering", () => {
    it("filters options when typing", () => {
      render(
        <SearchableSelect
          options={baseOptions}
          value=""
          onChange={vi.fn()}
        />,
      );

      const input = getInput();
      fireEvent.focus(input);

      fireEvent.change(input, { target: { value: "fire" } });

      // Options with highlight markup should appear in dropdown
      expect(screen.getByText(/Fire/i)).toBeInTheDocument();
    });

    it("shows no results message when no match", () => {
      render(
        <SearchableSelect
          options={baseOptions}
          value=""
          onChange={vi.fn()}
        />,
      );

      const input = getInput();
      fireEvent.focus(input);

      fireEvent.change(input, { target: { value: "xyz" } });

      expect(screen.getByText("No results")).toBeInTheDocument();
    });
  });

  describe("option selection", () => {
    it("calls onChange with option id on click", () => {
      const handleChange = vi.fn();
      render(
        <SearchableSelect
          options={baseOptions}
          value=""
          onChange={handleChange}
        />,
      );

      fireEvent.focus(getInput());
      fireEvent.click(screen.getByText("Fireball (1)"));

      expect(handleChange).toHaveBeenCalledWith("1");
    });

    it("closes dropdown after selection", () => {
      render(
        <SearchableSelect
          options={baseOptions}
          value=""
          onChange={vi.fn()}
        />,
      );

      fireEvent.focus(getInput());
      expect(screen.getByText("Fireball (1)")).toBeInTheDocument();

      fireEvent.click(screen.getByText("Fireball (1)"));
      expect(screen.queryByText("Icebolt (2)")).not.toBeInTheDocument();
    });
  });

  describe("keyboard navigation", () => {
    it("opens dropdown with ArrowDown", () => {
      render(
        <SearchableSelect
          options={baseOptions}
          value=""
          onChange={vi.fn()}
        />,
      );

      fireEvent.keyDown(getInput(), { key: "ArrowDown" });

      expect(screen.getByText("Fireball (1)")).toBeInTheDocument();
    });

    it("closes dropdown with Escape", () => {
      render(
        <SearchableSelect
          options={baseOptions}
          value=""
          onChange={vi.fn()}
        />,
      );

      fireEvent.focus(getInput());
      expect(screen.getByText("Fireball (1)")).toBeInTheDocument();

      fireEvent.keyDown(getInput(), { key: "Escape" });

      expect(screen.queryByText("Icebolt (2)")).not.toBeInTheDocument();
    });

    it("navigates options with arrow keys", () => {
      render(
        <SearchableSelect
          options={baseOptions}
          value=""
          onChange={vi.fn()}
        />,
      );

      fireEvent.keyDown(getInput(), { key: "ArrowDown" });
      expect(screen.getByText("Fireball (1)")).toBeInTheDocument();

      fireEvent.keyDown(getInput(), { key: "ArrowDown" });
      fireEvent.keyDown(getInput(), { key: "ArrowDown" });

      // Selection should change, no errors
      expect(getInput()).toBeInTheDocument();
    });

    it("selects highlighted option with Enter", () => {
      const handleChange = vi.fn();
      render(
        <SearchableSelect
          options={baseOptions}
          value=""
          onChange={handleChange}
        />,
      );

      fireEvent.focus(getInput());
      fireEvent.change(getInput(), { target: { value: "ice" } });

      fireEvent.keyDown(getInput(), { key: "Enter" });

      expect(handleChange).toHaveBeenCalled();
    });
  });

  describe("click outside", () => {
    it("closes dropdown when clicking outside", () => {
      render(
        <SearchableSelect
          options={baseOptions}
          value=""
          onChange={vi.fn()}
        />,
      );

      fireEvent.focus(getInput());
      expect(screen.getByText("Fireball (1)")).toBeInTheDocument();

      fireEvent.mouseDown(document.body);

      expect(screen.queryByText("Icebolt (2)")).not.toBeInTheDocument();
    });
  });
});