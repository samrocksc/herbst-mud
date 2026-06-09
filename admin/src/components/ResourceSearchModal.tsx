/* eslint-disable functional/no-mixed-types, functional/prefer-immutable-types, functional/immutable-data, react-hooks/set-state-in-effect */
 
 
 
 
 
import { useState, useRef, useEffect, useMemo } from "react";
import { Modal } from "./Modal";
import { Button } from "./Button";
import { useResourceSearch } from "./useResourceSearch";
import { fuzzyMatch, highlightMatch } from "./fuzzyMatch";

type Props = Readonly<{
  isOpen: boolean
  onClose: () => void
  title: string
  resourceType: string
  apiBase: string
  worldId?: string
  value?: number | string | null
  onSelect: (id: number | string) => void
  multi?: boolean
  selectedIds?: (number | string)[]
}>

export function ResourceSearchModal({
  isOpen, onClose, title, resourceType, apiBase, worldId,
  value, onSelect, multi = false, selectedIds = [],
}: Props) {
  const [query, setQuery] = useState("");
  const [highlightIdx, setHighlightIdx] = useState(-1);
  const inputRef = useRef<HTMLInputElement>(null);
  const { data: resultsRaw = [], isLoading } = useResourceSearch(resourceType, apiBase, query, worldId);
  // Guard: ensure results is an array (apiGet may return error object on failure)
  const results = Array.isArray(resultsRaw) ? resultsRaw : [];
  const filtered = useMemo(() => {
    if (!query.trim()) {
      // No query: show 5 most recent (highest IDs first) as suggestions
      return [...results].sort((a, b) => Number(b.id) - Number(a.id)).slice(0, 5);
    }
    return results.filter(
      (r: { id: number | string; name: string }) => fuzzyMatch(r.name, query) || fuzzyMatch(String(r.id), query),
    );
  }, [results, query]);

  // Use a ref to track previous filtered length to avoid cascading renders
  const prevFilteredLength = useRef(filtered.length);
  useEffect(() => {
    if (filtered.length !== prevFilteredLength.current) {
      setHighlightIdx(-1);
    }
    prevFilteredLength.current = filtered.length;
  }, [filtered.length]);

  useEffect(() => {
    if (isOpen) {
      setQuery("");
      // Trigger initial load of suggestions
      setTimeout(() => {
        // Force re-render with empty query to load initial results
      }, 0);
      inputRef.current?.focus();
    }
  }, [isOpen]);

  const handleSelect = (id: number | string) => {
    onSelect(id);
    if (!multi) onClose();
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "ArrowDown") {
      e.preventDefault();
      setHighlightIdx((i) => Math.min(i + 1, filtered.length - 1));
    } else if (e.key === "ArrowUp") {
      e.preventDefault();
      setHighlightIdx((i) => Math.max(i - 1, 0));
    } else if (e.key === "Enter" && highlightIdx >= 0 && highlightIdx < filtered.length) {
      e.preventDefault();
      const r = filtered[highlightIdx];
      handleSelect(r.id);
    } else if (e.key === "Escape") {
      e.preventDefault();
      onClose();
    }
  };

  if (!isOpen) return null;

  return (
    <Modal isOpen={isOpen} onClose={onClose} title={title}>
      <input
        ref={inputRef}
        type="text"
        value={query}
        onChange={(e) => setQuery(e.target.value)}
        onKeyDown={handleKeyDown}
        placeholder="Search by name..."
        className="w-full p-2 bg-surface border border-border rounded text-text text-sm mb-3"
        autoFocus
      />
      {isLoading ? (
        <div className="text-text-muted text-sm p-4 text-center">Loading...</div>
      ) : filtered.length === 0 ? (
        <div className="text-text-muted text-sm p-4 text-center">
          {query ? "No results found" : "Type to search..."}
        </div>
      ) : (
        <div className="max-h-64 overflow-y-auto">
          {filtered.map((r: { id: number | string; name: string }, idx: number) => {
            const isSelected = r.id === value || selectedIds.includes(r.id);
            const isHighlighted = idx === highlightIdx;
            return (
              <div
                key={r.id}
                role="option"
                aria-selected={isHighlighted}
                onClick={() => handleSelect(r.id)}
                className={`px-3 py-2 text-sm cursor-pointer flex items-center gap-2 ${
                  isHighlighted ? "bg-surface-hover" : ""
                } hover:bg-surface-hover ${isSelected ? "opacity-60" : ""}`}
              >
                {isSelected && <span className="text-primary">✓</span>}
                <span className="text-text-muted">#{r.id}</span>
                <span
                  className="text-text"
                  dangerouslySetInnerHTML={{
                    __html: query.trim() ? highlightMatch(r.name, query.trim()) : r.name,
                  }}
                />
              </div>
            );
          })}
        </div>
      )}
      {multi && (
        <div className="mt-3 flex justify-end">
          <Button variant="secondary" size="sm" onClick={onClose}>Done</Button>
        </div>
      )}
    </Modal>
  );
}
