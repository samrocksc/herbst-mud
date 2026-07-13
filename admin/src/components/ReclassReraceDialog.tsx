/* eslint-disable functional/prefer-immutable-types, react-hooks/purity */
import { useState, useMemo } from "react";
import { Modal } from "./Modal";
import { Button } from "./Button";
import { useClasses, type ClassOption } from "../hooks/useClasses";
import { useRaces, type Race } from "../hooks/useRaces";
import { useReclass, useRerace } from "../hooks/useCharacterHistory";

type DialogMode = "reclass" | "rerace";

type ReclassReraceDialogProps = Readonly<{
  open: boolean
  mode: DialogMode
  characterId: number
  currentName: string
  currentClassName: string
  currentRaceName: string
  onClose: () => void
}>;

/**
 * Shared dialog for reclass and rerace operations.
 *
 * Two-step flow:
 * 1. Select a new class (faction) or race from a searchable list.
 * 2. Confirm the destructive operation with a reason field.
 *
 * Calls POST /characters/:id/reclass or POST /characters/:id/rerace.
 */
export function ReclassReraceDialog({
  open,
  mode,
  characterId,
  currentName,
  currentClassName,
  currentRaceName,
  onClose,
}: ReclassReraceDialogProps) {
  const isReclass = mode === "reclass";
  const [search, setSearch] = useState("");
  const [selectedId, setSelectedId] = useState<number | null>(null);
  const [reason, setReason] = useState("");
  const [step, setStep] = useState<"select" | "confirm">("select");

  const { data: classes, isLoading: classesLoading } = useClasses();
  const { data: races, isLoading: racesLoading } = useRaces();
  const reclassMutation = useReclass();
  const reraceMutation = useRerace();

  const mutation = isReclass ? reclassMutation : reraceMutation;
  const isLoading = isReclass ? classesLoading : racesLoading;
  const error = mutation.error as Error | null;

  const items = useMemo(() => {
    if (isReclass) {
      return (classes ?? []).map((c: ClassOption) => ({
        id: c.id,
        name: c.name,
        displayName: c.display_name ?? c.name,
      }));
    }
    return (races ?? []).map((r: Race) => ({
      id: r.id,
      name: r.name,
      displayName: r.display_name || r.name,
    }));
  }, [isReclass, classes, races]);

  const filtered = items.filter((item) =>
    item.name.toLowerCase().includes(search.toLowerCase()) ||
    item.displayName.toLowerCase().includes(search.toLowerCase()),
  );

  const selectedItem = items.find((i) => i.id === selectedId) ?? null;
  const currentValue = isReclass ? currentClassName : currentRaceName;
  const label = isReclass ? "Class" : "Race";

  const handleClose = () => {
    setSearch("");
    setSelectedId(null);
    setReason("");
    setStep("select");
    onClose();
  };

  const handleExecute = () => {
    if (selectedId === null) return;
    const payload = {
      characterId,
      reason: reason.trim() || undefined,
    };
    const mutate = isReclass
      ? reclassMutation.mutate
      : reraceMutation.mutate;

    if (isReclass) {
      reclassMutation.mutate(
        { ...payload, factionId: selectedId },
        { onSuccess: handleClose },
      );
    } else {
      reraceMutation.mutate(
        { ...payload, raceId: selectedId },
        { onSuccess: handleClose },
      );
    }
  };

  return (
    <Modal
      isOpen={open}
      onClose={handleClose}
      title={isReclass ? `Reclass ${currentName}` : `Rerace ${currentName}`}
    >
      {step === "select" && (
        <div className="space-y-3">
          {/* Current value */}
          <div className="text-xs text-text-muted">
            Current {label.toLowerCase()}: <span className="text-text font-medium">{currentValue || "—"}</span>
          </div>

          {/* Search */}
          <input
            type="text"
            placeholder={`Search ${label.toLowerCase()}…`}
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
            autoFocus
          />

          {isLoading && <div className="text-text-muted text-xs">Loading…</div>}

          {/* List */}
          <div className="max-h-56 overflow-y-auto space-y-1">
            {filtered.map((item) => (
              <button
                key={item.id}
                className={`w-full text-left px-3 py-2 rounded text-sm text-text hover:bg-surface-muted ${
                  selectedId === item.id ? "bg-primary/20 border border-primary" : "border border-transparent"
                }`}
                onClick={() => setSelectedId(item.id)}
                disabled={mutation.isPending}
              >
                <span className="font-medium">{item.displayName}</span>
                {item.name !== item.displayName && (
                  <span className="text-text-muted ml-2">({item.name})</span>
                )}
              </button>
            ))}
            {!isLoading && filtered.length === 0 && (
              <div className="text-text-muted text-xs py-2">No {label.toLowerCase()} found</div>
            )}
          </div>

          {/* Actions */}
          <div className="flex justify-end gap-2 pt-2 border-t border-border">
            <Button variant="secondary" size="sm" onClick={handleClose}>Cancel</Button>
            <Button
              variant="primary"
              size="sm"
              disabled={selectedId === null || mutation.isPending}
              onClick={() => setStep("confirm")}
            >
              Next
            </Button>
          </div>
        </div>
      )}

      {step === "confirm" && selectedItem && (
        <div className="space-y-4">
          {/* Warning */}
          <div className="bg-danger/10 border border-danger rounded-md p-3">
            <div className="text-danger font-semibold text-sm mb-1">⚠ Warning</div>
            <div className="text-text text-sm">
              You are about to change <span className="font-bold">{currentName}</span>'s {label.toLowerCase()} from{" "}
              <span className="font-bold">{currentValue || "—"}</span> to{" "}
              <span className="font-bold">{selectedItem.displayName}</span>.
            </div>
            <div className="text-text-muted text-xs mt-1">
              This is a destructive operation and will be recorded in the character's history.
            </div>
          </div>

          {/* Reason */}
          <div>
            <label className="text-text-muted text-sm block mb-1">Reason (optional)</label>
            <input
              type="text"
              placeholder="Admin reclass/rerace reason…"
              value={reason}
              onChange={(e) => setReason(e.target.value)}
              className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
              autoFocus
            />
          </div>

          {/* Error */}
          {error && (
            <div className="text-danger text-xs bg-danger/10 border border-danger rounded p-2">
              {error.message}
            </div>
          )}

          {/* Actions */}
          <div className="flex justify-end gap-2 pt-2 border-t border-border">
            <Button
              variant="secondary"
              size="sm"
              onClick={() => setStep("select")}
              disabled={mutation.isPending}
            >
              Back
            </Button>
            <Button
              variant="danger"
              size="sm"
              disabled={mutation.isPending}
              onClick={handleExecute}
            >
              {mutation.isPending ? "Executing…" : `Confirm ${isReclass ? "Reclass" : "Rerace"}`}
            </Button>
          </div>
        </div>
      )}
    </Modal>
  );
}