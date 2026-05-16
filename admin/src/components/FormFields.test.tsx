 
import { describe, it, expect, vi } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";
import {
  FormField,
  NumberField,
  TextareaField,
  SelectField,
  CheckboxField,
  FormError,
} from "./FormFields";

// ─── FormField (text input) ──────────────────────────────────────────────

describe("FormField", () => {
  it("renders label and input", () => {
    render(<FormField label="Name" value="" onChange={() => {}} />);
    expect(screen.getByLabelText("Name")).toBeInTheDocument();
  });

  it("calls onChange with typed value", () => {
    const onChange = vi.fn();
    render(<FormField label="Name" value="" onChange={onChange} />);
    fireEvent.change(screen.getByLabelText("Name"), { target: { value: "Gizmo" } });
    expect(onChange).toHaveBeenCalledWith("Gizmo");
  });

  it("renders placeholder", () => {
    render(<FormField label="Name" value="" onChange={() => {}} placeholder="Enter name" />);
    expect(screen.getByPlaceholderText("Enter name")).toBeInTheDocument();
  });

  it("disables input when disabled=true", () => {
    render(<FormField label="Name" value="" onChange={() => {}} disabled />);
    expect(screen.getByLabelText("Name")).toBeDisabled();
  });

  it("uses custom id when provided", () => {
    render(<FormField label="Full Name" id="custom-id" value="" onChange={() => {}} />);
    expect(screen.getByLabelText("Full Name")).toHaveAttribute("id", "custom-id");
  });

  it("spreads extra HTML attributes onto input", () => {
    render(<FormField label="Name" value="" onChange={() => {}} data-testid="name-field" />);
    expect(screen.getByTestId("name-field")).toBeInTheDocument();
  });
});

// ─── NumberField ─────────────────────────────────────────────────────────

describe("NumberField", () => {
  it("renders label and input with numeric keyboard hint", () => {
    render(<NumberField label="Level" value={1} onChange={() => {}} />);
    const input = screen.getByLabelText("Level");
    expect(input).toBeInTheDocument();
    expect(input).toHaveAttribute("inputMode", "numeric");
  });

  it("calls onChange with parsed integer", () => {
    const onChange = vi.fn();
    render(<NumberField label="Level" value={0} onChange={onChange} />);
    fireEvent.change(screen.getByLabelText("Level"), { target: { value: "42" } });
    expect(onChange).toHaveBeenCalledWith(42);
  });

  it("calls onChange with 0 for empty input", () => {
    const onChange = vi.fn();
    render(<NumberField label="Level" value={5} onChange={onChange} />);
    fireEvent.change(screen.getByLabelText("Level"), { target: { value: "" } });
    expect(onChange).toHaveBeenCalledWith(0);
  });

  it("ignores non-numeric input", () => {
    const onChange = vi.fn();
    render(<NumberField label="Level" value={5} onChange={onChange} />);
    fireEvent.change(screen.getByLabelText("Level"), { target: { value: "abc" } });
    expect(onChange).not.toHaveBeenCalled();
  });

  it("allows negative numbers", () => {
    const onChange = vi.fn();
    render(<NumberField label="Level" value={5} onChange={onChange} />);
    fireEvent.change(screen.getByLabelText("Level"), { target: { value: "-3" } });
    expect(onChange).toHaveBeenCalledWith(-3);
  });

  it("displays empty string when value is NaN", () => {
    render(<NumberField label="Level" value={NaN} onChange={() => {}} />);
    expect(screen.getByLabelText("Level")).toHaveDisplayValue("");
  });

  it("displays empty string when value is Infinity", () => {
    render(<NumberField label="Level" value={Infinity} onChange={() => {}} />);
    expect(screen.getByLabelText("Level")).toHaveDisplayValue("");
  });
});

// ─── TextareaField ───────────────────────────────────────────────────────

describe("TextareaField", () => {
  it("renders label and textarea", () => {
    render(<TextareaField label="Description" value="" onChange={() => {}} />);
    expect(screen.getByLabelText("Description")).toBeInTheDocument();
  });

  it("calls onChange with typed value", () => {
    const onChange = vi.fn();
    render(<TextareaField label="Description" value="" onChange={onChange} />);
    fireEvent.change(screen.getByLabelText("Description"), { target: { value: "A dark cave." } });
    expect(onChange).toHaveBeenCalledWith("A dark cave.");
  });

  it("uses specified rows", () => {
    render(<TextareaField label="Description" value="" onChange={() => {}} rows={5} />);
    expect(screen.getByLabelText("Description")).toHaveAttribute("rows", "5");
  });

  it("defaults to 3 rows", () => {
    render(<TextareaField label="Description" value="" onChange={() => {}} />);
    expect(screen.getByLabelText("Description")).toHaveAttribute("rows", "3");
  });

  it("disables textarea when disabled=true", () => {
    render(<TextareaField label="Description" value="" onChange={() => {}} disabled />);
    expect(screen.getByLabelText("Description")).toBeDisabled();
  });
});

// ─── SelectField ─────────────────────────────────────────────────────────

describe("SelectField", () => {
  const options = [
    { value: "hostile", label: "Hostile" },
    { value: "neutral", label: "Neutral" },
    { value: "friendly", label: "Friendly" },
  ];

  it("renders label and select with options", () => {
    render(<SelectField label="Disposition" value="neutral" onChange={() => {}} options={options} />);
    expect(screen.getByLabelText("Disposition")).toBeInTheDocument();
    expect(screen.getByText("Hostile")).toBeInTheDocument();
    expect(screen.getByText("Neutral")).toBeInTheDocument();
  });

  it("calls onChange with selected value", () => {
    const onChange = vi.fn();
    render(<SelectField label="Disposition" value="neutral" onChange={onChange} options={options} />);
    fireEvent.change(screen.getByLabelText("Disposition"), { target: { value: "hostile" } });
    expect(onChange).toHaveBeenCalledWith("hostile");
  });

  it("renders placeholder option when provided", () => {
    render(
      <SelectField
        label="Disposition"
        value=""
        onChange={() => {}}
        options={options}
        placeholder="— Select —"
      />,
    );
    expect(screen.getByText("— Select —")).toBeInTheDocument();
  });
});

// ─── CheckboxField ───────────────────────────────────────────────────────

describe("CheckboxField", () => {
  it("renders checked checkbox", () => {
    render(<CheckboxField label="Visible" checked={true} onChange={() => {}} />);
    expect(screen.getByLabelText("Visible")).toBeChecked();
  });

  it("renders unchecked checkbox", () => {
    render(<CheckboxField label="Visible" checked={false} onChange={() => {}} />);
    expect(screen.getByLabelText("Visible")).not.toBeChecked();
  });

  it("calls onChange with toggled boolean", () => {
    const onChange = vi.fn();
    render(<CheckboxField label="Visible" checked={false} onChange={onChange} />);
    fireEvent.click(screen.getByLabelText("Visible"));
    expect(onChange).toHaveBeenCalledWith(true);
  });
});

// ─── FormError ───────────────────────────────────────────────────────────

describe("FormError", () => {
  it("renders error message", () => {
    render(<FormError message="Something went wrong" />);
    expect(screen.getByText("Something went wrong")).toBeInTheDocument();
  });
});
