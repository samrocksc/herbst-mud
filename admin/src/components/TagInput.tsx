import { useState, useRef, useEffect } from 'react';
import type { KeyboardEvent } from 'react';
import { Button } from './Button';
import { TooltipIcon } from './Tooltip';

export type TagInputProps = Readonly<{
  /** Current selected tags */
  value: string[]
  /** Called when tags change */
  onChange: (tags: string[]) => void
  /** All known tags for autocomplete */
  availableTags?: string[]
  /** Placeholder shown when input is empty */
  placeholder?: string
  /** Disable the input */
  disabled?: boolean
  /** Label shown above the input */
  label?: string
  /** Optional tooltip explaining this field */
  tooltip?: string
}>

/**
 * Chip-style tag input with fuzzy autocomplete.
 *
 * - Type to filter known tags
 * - Press Enter to add the typed text as a new tag
 * - Click a suggestion to add that tag
 * - Click x on a chip to remove it
 */
export function TagInput({
  value,
  onChange,
  availableTags = [],
  placeholder = 'Add a tag…',
  disabled = false,
  label,
  tooltip,
}: TagInputProps) {
  const [inputValue, setInputValue] = useState('');
  const [isOpen, setIsOpen] = useState(false);
  const [highlightedIndex, setHighlightedIndex] = useState(0);
  const inputRef = useRef<HTMLInputElement>(null);
  const containerRef = useRef<HTMLDivElement>(null);

  // Close dropdown when clicking outside
  useEffect(() => {
    const handler = (e: MouseEvent) => {
      if (containerRef.current && !containerRef.current.contains(e.target as Node)) {
        setIsOpen(false);
      }
    };
    document.addEventListener('mousedown', handler);
    return () => document.removeEventListener('mousedown', handler);
  }, []);

  const normalised = inputValue.toLowerCase().trim();

  /** Tags that match the current input (excludes already-selected) */
  const suggestions = availableTags.filter(
    (t) =>
      t.toLowerCase().includes(normalised) &&
      !value.includes(t)
  );

  /** Whether the exact input text is a new tag (not in availableTags) */
  const inputIsNewTag =
    normalised.length > 0 &&
    !availableTags.map((t) => t.toLowerCase()).includes(normalised);

  const addTag = (tag: string) => {
    const trimmed = tag.trim();
    if (!trimmed || value.includes(trimmed)) return;
    onChange([...value, trimmed]);
    setInputValue('');
    setIsOpen(false);
    inputRef.current?.focus();
  };

  const removeTag = (tag: string) => {
    onChange(value.filter((t) => t !== tag));
  };

  const handleKeyDown = (e: KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      if (suggestions.length > 0) {
        addTag(suggestions[highlightedIndex]);
      } else if (inputIsNewTag) {
        addTag(inputValue);
      } else if (normalised.length > 0) {
        // Existing but not selected — add it anyway
        addTag(inputValue);
      }
    } else if (e.key === 'ArrowDown') {
      e.preventDefault();
      setIsOpen(true);
      setHighlightedIndex((i) => Math.min(i + 1, suggestions.length - (inputIsNewTag ? 0 : 1)));
    } else if (e.key === 'ArrowUp') {
      e.preventDefault();
      setHighlightedIndex((i) => Math.max(i - 1, 0));
    } else if (e.key === 'Escape') {
      setIsOpen(false);
      setInputValue('');
    } else if (e.key === 'Backspace' && inputValue === '' && value.length > 0) {
      removeTag(value[value.length - 1]);
    }
  };

  const showDropdown = isOpen && !disabled && (
    suggestions.length > 0 || inputIsNewTag
  );

  return (
    <div ref={containerRef} className="tag-input-container">
      {label && (
        <label className="form-label">
          {label}
          {tooltip && <TooltipIcon content={tooltip} />}
        </label>
      )}

      {/* Selected chips */}
      {value.length > 0 && (
        <div className="tag-chips">
          {value.map((tag) => (
            <span key={tag} className="tag-chip">
              {tag}
              <Button
                type="button"
                variant="ghost"
                size="sm"
                className="tag-chip-remove"
                onClick={() => removeTag(tag)}
                disabled={disabled}
                aria-label={`Remove ${tag}`}
              >
                x
              </Button>
            </span>
          ))}
        </div>
      )}

      {/* Input + dropdown */}
      <div className="tag-input-wrapper">
        <input
          ref={inputRef}
          type="text"
          value={inputValue}
          onChange={(e) => {
            setInputValue(e.target.value);
            setIsOpen(true);
            setHighlightedIndex(0);
          }}
          onFocus={() => setIsOpen(true)}
          onKeyDown={handleKeyDown}
          placeholder={placeholder}
          disabled={disabled}
          className="tag-text-input"
          autoComplete="off"
        />

        {showDropdown && (
          <ul className="tag-dropdown" role="listbox">
            {suggestions.map((tag, i) => (
              <li
                key={tag}
                role="option"
                aria-selected={i === highlightedIndex}
                className={`tag-dropdown-item${i === highlightedIndex ? ' highlighted' : ''}`}
                onMouseDown={(e) => { e.preventDefault(); addTag(tag); }}
                onMouseEnter={() => setHighlightedIndex(i)}
              >
                {tag}
              </li>
            ))}
            {inputIsNewTag && (
              <li
                role="option"
                aria-selected={highlightedIndex === suggestions.length}
                className={`tag-dropdown-item tag-dropdown-item-create${
                  highlightedIndex === suggestions.length ? ' highlighted' : ''
                }`}
                onMouseDown={(e) => { e.preventDefault(); addTag(inputValue); }}
                onMouseEnter={() => setHighlightedIndex(suggestions.length)}
              >
                + Create &ldquo;{inputValue}&rdquo;
              </li>
            )}
          </ul>
        )}
      </div>
    </div>
  );
}