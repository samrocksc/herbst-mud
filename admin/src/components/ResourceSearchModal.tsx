import { useState, useRef, useEffect, useMemo } from 'react';
import { Modal } from './Modal';
import { Button } from './Button';
import { useResourceSearch } from './useResourceSearch';
import { fuzzyMatch, highlightMatch } from './fuzzyMatch';

type Props = Readonly<{
  isOpen: boolean
  onClose: () => void
  title: string
  resourceType: string
  apiBase: string
  value?: number | string | null
  onSelect: (id: number | string, name: string) => void
  multi?: boolean
  selectedIds?: (number | string)[]
}>

export function ResourceSearchModal({
  isOpen, onClose, title, resourceType, apiBase,
  value, onSelect, multi = false, selectedIds = [],
}: Props) {
  const [query, setQuery] = useState('');
  const [highlightIdx, setHighlightIdx] = useState(-1);
  const inputRef = useRef<HTMLInputElement>(null);
  const { data: results = [], isLoading } = useResourceSearch(resourceType, apiBase, query);
  const filtered = useMemo(() => {
    if (!query.trim()) return results.slice(0, 50);
    return results.filter(
      (r: { id: number | string; name: string }) => fuzzyMatch(r.name, query) || fuzzyMatch(String(r.id), query),
    );
  }, [results, query]);

  useEffect(() => { setHighlightIdx(-1); }, [filtered.length]);
  useEffect(() => { if (isOpen) { setQuery(''); inputRef.current?.focus(); } }, [isOpen]);

  const handleSelect = (id: number | string, name: string) => {
    onSelect(id, name);
    if (!multi) onClose();
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'ArrowDown') {
      e.preventDefault();
      setHighlightIdx((i) => Math.min(i + 1, filtered.length - 1));
    } else if (e.key === 'ArrowUp') {
      e.preventDefault();
      setHighlightIdx((i) => Math.max(i - 1, 0));
    } else if (e.key === 'Enter' && highlightIdx >= 0 && highlightIdx < filtered.length) {
      e.preventDefault();
      const r = filtered[highlightIdx];
      handleSelect(r.id, r.name);
    } else if (e.key === 'Escape') {
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
          {query ? 'No results found' : 'Type to search...'}
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
                onClick={() => handleSelect(r.id, r.name)}
                className={`px-3 py-2 text-sm cursor-pointer flex items-center gap-2 ${
                  isHighlighted ? 'bg-surface-hover' : ''
                } hover:bg-surface-hover ${isSelected ? 'opacity-60' : ''}`}
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