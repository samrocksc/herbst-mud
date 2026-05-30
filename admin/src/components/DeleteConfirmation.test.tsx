import { describe, it, expect, vi } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";
import { DeleteConfirmation } from "./DeleteConfirmation";
import { Modal } from "./Modal";

vi.mock("./Modal", () => ({
  Modal: ({ isOpen, children, title, onClose }: { isOpen: boolean; children: React.ReactNode; title: string; onClose: () => void }) =>
    isOpen ? (
      <div role="dialog" aria-modal="true" data-testid="modal">
        <div data-testid="modal-title">{title}</div>
        <div data-testid="modal-content">{children}</div>
        <button type="button" onClick={onClose} data-testid="modal-close">Close</button>
      </div>
    ) : null,
}));

describe("DeleteConfirmation", () => {
  describe("rendering", () => {
    it("renders modal when open", () => {
      render(
        <DeleteConfirmation
          open={true}
          message="Are you sure you want to delete this ability?"
          onConfirm={vi.fn()}
          onCancel={vi.fn()}
        />,
      );

      expect(screen.getByTestId("modal")).toBeInTheDocument();
    });

    it("does not render modal when closed", () => {
      render(
        <DeleteConfirmation
          open={false}
          message="Are you sure?"
          onConfirm={vi.fn()}
          onCancel={vi.fn()}
        />,
      );

      expect(screen.queryByTestId("modal")).not.toBeInTheDocument();
    });

    it("renders message text", () => {
      const message = "Are you sure you want to delete Fireball?";
      render(
        <DeleteConfirmation
          open={true}
          message={message}
          onConfirm={vi.fn()}
          onCancel={vi.fn()}
        />,
      );

      expect(screen.getByText(message)).toBeInTheDocument();
    });
  });

  describe("title", () => {
    it("renders default title", () => {
      render(
        <DeleteConfirmation
          open={true}
          message="Delete?"
          onConfirm={vi.fn()}
          onCancel={vi.fn()}
        />,
      );

      expect(screen.getByTestId("modal-title")).toHaveTextContent("Confirm Delete");
    });

    it("renders custom title", () => {
      render(
        <DeleteConfirmation
          open={true}
          title="Delete Ability"
          message="Delete?"
          onConfirm={vi.fn()}
          onCancel={vi.fn()}
        />,
      );

      expect(screen.getByTestId("modal-title")).toHaveTextContent("Delete Ability");
    });
  });

  describe("buttons", () => {
    it("renders Cancel button", () => {
      render(
        <DeleteConfirmation
          open={true}
          message="Delete?"
          onConfirm={vi.fn()}
          onCancel={vi.fn()}
        />,
      );

      expect(screen.getByRole("button", { name: "Cancel" })).toBeInTheDocument();
    });

    it("renders Delete button", () => {
      render(
        <DeleteConfirmation
          open={true}
          message="Delete?"
          onConfirm={vi.fn()}
          onCancel={vi.fn()}
        />,
      );

      expect(screen.getByRole("button", { name: "Delete" })).toBeInTheDocument();
    });

    it("renders loading text when isLoading is true", () => {
      render(
        <DeleteConfirmation
          open={true}
          message="Delete?"
          onConfirm={vi.fn()}
          onCancel={vi.fn()}
          isLoading={true}
        />,
      );

      expect(screen.getByText("Deleting…")).toBeInTheDocument();
      expect(screen.queryByRole("button", { name: "Delete" })).not.toBeInTheDocument();
    });
  });

  describe("interactions", () => {
    it("calls onConfirm when Delete button is clicked", () => {
      const handleConfirm = vi.fn();
      render(
        <DeleteConfirmation
          open={true}
          message="Delete?"
          onConfirm={handleConfirm}
          onCancel={vi.fn()}
        />,
      );

      fireEvent.click(screen.getByRole("button", { name: "Delete" }));

      expect(handleConfirm).toHaveBeenCalledTimes(1);
    });

    it("calls onCancel when Cancel button is clicked", () => {
      const handleCancel = vi.fn();
      render(
        <DeleteConfirmation
          open={true}
          message="Delete?"
          onConfirm={vi.fn()}
          onCancel={handleCancel}
        />,
      );

      fireEvent.click(screen.getByRole("button", { name: "Cancel" }));

      expect(handleCancel).toHaveBeenCalledTimes(1);
    });

    it("disables buttons when isLoading", () => {
      render(
        <DeleteConfirmation
          open={true}
          message="Delete?"
          onConfirm={vi.fn()}
          onCancel={vi.fn()}
          isLoading={true}
        />,
      );

      const cancelBtn = screen.getByRole("button", { name: "Cancel" });
      expect(cancelBtn).toBeDisabled();
    });
  });

  describe("accessibility", () => {
    it("modal has role dialog", () => {
      render(
        <DeleteConfirmation
          open={true}
          message="Delete?"
          onConfirm={vi.fn()}
          onCancel={vi.fn()}
        />,
      );

      expect(screen.getByRole("dialog")).toBeInTheDocument();
    });

    it("modal is marked as modal", () => {
      render(
        <DeleteConfirmation
          open={true}
          message="Delete?"
          onConfirm={vi.fn()}
          onCancel={vi.fn()}
        />,
      );

      expect(screen.getByRole("dialog")).toHaveAttribute("aria-modal", "true");
    });
  });
});