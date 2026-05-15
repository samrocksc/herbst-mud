import { useState, useEffect, useRef, useMemo } from 'react';
import { useResourceSearch } from './useResourceSearch';
import { fuzzyMatch, highlightMatch } from './fuzzyMatch';
import { FieldLabel } from './fields/FieldLabel';

type Props = Readonly<{
  label: string
  value: number | string | null | undefined
  onChange: (id: number | string | null) => void
  resourceType: string
  apiBase: string
  tooltip?: string
  disabled?: boolean
  placeholder?: string
}>

export function ResourceSearchSelect({
  label, value, onChange, resourceType, apiBase, tooltip, disabled, placeholder,
}: Props) {
  const [query, setQuery] = useState('');
  const [isOpen, setIsOpen] = useState(false);
  const [highlightIdx, setHighlightIdx] = useState(-1);
  const containerRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);
  const { data: results = [], isLoading } = useResourceSearch(resourceType, apiBase, query);

  // Ensure results is always an array for safety
  const safeResults = Array.isArray(results) ? results : [];

  // Find the selected resource name for display
  const selectedName = useMemo(() => {
    if (value == null || value === '') return '';
    const found = safeResults.find((r: { id: number | string; name: string }) => String(r.id) === String(value));
    return found ? found.name : '';
  }, [safeResults, value]);

  // Filter results with fuzzy match
  const filtered = useMemo(() => {
    if (!query.trim()) return safeResults.slice(0, 50);
    return safeResults.filter(
      (r: { id: number | string; name: string }) => fuzzyMatch(r.name, query) || fuzzyMatch(String(r.id), query),
    );
  }, [safeResults, query]);

  useEffect(() => { setHighlightIdx(-1); }, [filtered.length]);

  // Click outside closes dropdown
  useEffect(() => {
    function handleClick(e: MouseEvent) {
      if (containerRef.current && !containerRef.current.contains(e.target as Node)) {
        setIsOpen(false);
      }
    }
    document.addEventListener('mousedown', handleClick);
    return () => document.removeEventListener('mousedown', handleClick);
  }, []);

  const displayValue = isOpen
    ? query
    : value != null && value !== ''
      ? selectedName
        ? `${selectedName} (#${value})`
        : `#${value}`
      : '';

  const handleSelect = (id: number | string, name: string) => {
    onChange(id);
    setIsOpen(false);
    setQuery('');
    inputRef.current?.blur();
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (!isOpen && (e.key === 'ArrowDown' || e.key === 'Enter')) {
      setIsOpen(true);
      return;
    }
    switch (e.key) {
      case 'ArrowDown': {
        e.preventDefault();
        setHighlightIdx((i) => Math.min(i + 1, filtered.length - 1));
        break;
      }
      case 'ArrowUp': {
        e.preventDefault();
        setHighlightIdx((i) => Math.max(i - 1, 0));
        break;
      }
      case 'Enter': {
        e.preventDefault();
        if (highlightIdx >= 0 && highlightIdx < filtered.length) {
          const r = filtered[highlightIdx];
          handleSelect(r.id, r.name);
        } else if (filtered.length === 1) {
          handleSelect(filtered[0].id, filtered[0].name);
        }
        break;
      }
      case 'Escape': {
        e.preventDefault();
        setIsOpen(false);
        setQuery('');
        inputRef.current?.blur();
        break;
      }
    }
  };

  const fieldId = label.toLowerCase().replace(/\s+/g, '-');

  return (
    <div ref={containerRef} className="relative">
      <FieldLabel htmlFor={fieldId} tooltip={tooltip}>{label}</FieldLabel>
      <input
        ref={inputRef}
        id={fieldId}
        type="text"
        value={displayValue}
        onChange={(e) => { setQuery(e.target.value); setIsOpen(true); }}
        onFocus={() => setIsOpen(true)}
        onKeyDown={handleKeyDown}
        placeholder={placeholder ?? `Search ${resourceType}...`}
        disabled={disabled}
        className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
      />
      {isOpen && (
        <div className="absolute z-50 mt-1 w-full bg-surface border border-border rounded shadow-lg max-h-48 overflow-y-auto">
          {isLoading ? (
            <div className="px-3 py-1.5 text-sm text-text-muted">Searching...</div>
          ) : filtered.length === 0 ? (
            <div className="px-3 py-1.5 text-sm text-text-muted">
              {query.trim() ? 'No results found' : 'Type to search...'}
            </div>
          ) : (
            filtered.map((r: { id: number | string; name: string }, idx: number) => {
              const isHighlighted = idx === highlightIdx;
              const displayName = `${r.name} (#${r.id})`;
              return (
                <div
                  key={r.id}
                  role="option"
                  aria-selected={isHighlighted}
                  onClick={() => handleSelect(r.id, r.name)}
                  className={`px-3 py-1.5 text-sm cursor-pointer flex items-center gap-2 ${
                    isHighlighted ? 'bg-surface-hover' : ''
                  } hover:bg-surface-hover ${
                    String(r.id) === String(value) ? 'opacity-60' : ''
                  }`}
                >
                  <span
                    className="text-text"
                    dangerouslySetInnerHTML={{
                      __html: query.trim()
                        ? highlightMatch(displayName, query.trim())
                        : displayName,
                    }}
                  />
                </div>
              );
            })
          )}
        </div>
      )}
    </div>
  );
}
