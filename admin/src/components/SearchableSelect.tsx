/* eslint-disable functional/no-mixed-types, functional/immutable-data, react-hooks/set-state-in-effect */
 
 
 
 
 
 
import { useState, useEffect, useCallback, useRef, useMemo } from "react";
import { fuzzyMatch, highlightMatch } from "./fuzzyMatch";

// ─── Types ──────────────────────────────────────────────────────────────────

export type SearchableSelectOption = Readonly<{
  id: string
  name: string
}>

type SearchableSelectProps = Readonly<{
  options: SearchableSelectOption[]
  value: string
  onChange: (id: string) => void
  placeholder?: string
  disabled?: boolean
  label?: string
}>

// ─── Helpers ────────────────────────────────────────────────────────────────

function getDisplayLabel(
  value: string,
  options: ReadonlyArray<SearchableSelectOption>,
): string {
  if (!value) return "(not selected)";
  const found = options.find((o) => o.id === value);
  return found ? `${found.name} (${found.id})` : value;
}

// ─── Component ──────────────────────────────────────────────────────────────

export function SearchableSelect({
  options,
  value,
  onChange,
  placeholder = "Search...",
  disabled = false,
  label,
}: SearchableSelectProps) {
  const [query, setQuery] = useState("");
  const [isOpen, setIsOpen] = useState(false);
  const [highlightIdx, setHighlightIdx] = useState(-1);
  const containerRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);

  // ── Filtered options ────────────────────────────────────────────────────

  const filtered = useMemo(() => {
    if (!query.trim()) return options;
    const q = query.trim();
    return options.filter(
      (o) => fuzzyMatch(o.name, q) || fuzzyMatch(o.id, q),
    );
  }, [options, query]);

  // ── Click outside closes dropdown ───────────────────────────────────────

  useEffect(() => {
    function handleDocumentClick(e: MouseEvent) {
      if (
        containerRef.current &&
        !containerRef.current.contains(e.target as Node)
      ) {
        setIsOpen(false);
        setQuery("");
      }
    }
    document.addEventListener("mousedown", handleDocumentClick);
    return () => document.removeEventListener("mousedown", handleDocumentClick);
  }, []);

  // ── Reset highlight when filtered list changes ─────────────────────────

  // Using a ref to track previous length avoids cascading renders from useEffect
  const prevFilteredLength = useRef(filtered.length);
  useEffect(() => {
    if (filtered.length !== prevFilteredLength.current) {
      setHighlightIdx(-1);
    }
    prevFilteredLength.current = filtered.length;
  }, [filtered.length]);

  // ── Handlers ────────────────────────────────────────────────────────────

  const handleFocus = useCallback(() => {
    setIsOpen(true);
    setQuery("");
  }, []);

  const handleInputChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      setQuery(e.target.value);
      setIsOpen(true);
    },
    [],
  );

  const selectOption = useCallback(
    (id: string) => {
      onChange(id);
      setIsOpen(false);
      setQuery("");
    },
    [onChange],
  );

  const handleKeyDown = useCallback(
    (e: React.KeyboardEvent) => {
      if (!isOpen) {
        if (e.key === "ArrowDown" || e.key === "Enter") {
          setIsOpen(true);
          return;
        }
        return;
      }

      switch (e.key) {
        case "ArrowDown": {
          e.preventDefault();
          setHighlightIdx((prev) =>
            prev < filtered.length - 1 ? prev + 1 : 0,
          );
          break;
        }
        case "ArrowUp": {
          e.preventDefault();
          setHighlightIdx((prev) =>
            prev > 0 ? prev - 1 : filtered.length - 1,
          );
          break;
        }
        case "Enter": {
          e.preventDefault();
          if (highlightIdx >= 0 && highlightIdx < filtered.length) {
            selectOption(filtered[highlightIdx].id);
          } else if (filtered.length === 1) {
            selectOption(filtered[0].id);
          }
          break;
        }
        case "Escape": {
          e.preventDefault();
          setIsOpen(false);
          setQuery("");
          inputRef.current?.blur();
          break;
        }
      }
    },
    [isOpen, filtered, highlightIdx, selectOption],
  );

  // ── Render ──────────────────────────────────────────────────────────────

  const displayValue = isOpen ? query : getDisplayLabel(value, options);

  return (
    <div ref={containerRef} className="relative">
      {label && (
        <label className="text-text-muted text-xs block mb-1">{label}</label>
      )}
      <input
        ref={inputRef}
        type="text"
        value={displayValue}
        onChange={handleInputChange}
        onFocus={handleFocus}
        onKeyDown={handleKeyDown}
        placeholder={placeholder}
        disabled={disabled}
        aria-label={label ?? "Search"}
        className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
      />
      {isOpen && (
        <div className="absolute z-50 mt-1 w-full bg-surface border border-border rounded shadow-lg max-h-48 overflow-y-auto">
          {filtered.length === 0 ? (
            <div className="px-3 py-1.5 text-sm text-text-muted">
              No results
            </div>
          ) : (
            filtered.map((opt, idx) => {
              const isHighlighted = idx === highlightIdx;
              const isMatch = fuzzyMatch(opt.name, query.trim()) || fuzzyMatch(opt.id, query.trim());
              const displayName = isMatch
                ? `${opt.name} (${opt.id})`
                : `${opt.name} (${opt.id})`;

              function handleClick() {
                selectOption(opt.id);
              }

              return (
                <div
                  key={opt.id}
                  role="option"
                  aria-selected={isHighlighted}
                  onClick={handleClick}
                  className={`px-3 py-1.5 text-sm text-text cursor-pointer ${
                    isHighlighted ? "bg-surface-hover" : ""
                  } hover:bg-surface-hover`}
                  dangerouslySetInnerHTML={{
                    __html: query.trim()
                      ? highlightMatch(displayName, query.trim())
                      : displayName,
                  }}
                />
              );
            })
          )}
        </div>
      )}
    </div>
  );
}
