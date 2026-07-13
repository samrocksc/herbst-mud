/* eslint-disable functional/immutable-data */
import { useState } from "react";

/**
 * Reusable editor for a map of string → number (e.g. resistances, vulnerabilities).
 * Renders dynamic add/remove rows with a text input for the key and a number input for the value.
 */
type KeyValueEditorProps = Readonly<{
  label: string
  value: Record<string, number>
  onChange: (value: Record<string, number>) => void
  keyPlaceholder?: string
  valueLabel?: string
  tooltip?: string
  suggestions?: readonly string[]
}>

export function KeyValueEditor({
  label,
  value,
  onChange,
  keyPlaceholder = "damage type",
  valueLabel = "%",
  tooltip,
  suggestions,
}: KeyValueEditorProps) {
  const entries = Object.entries(value);
  const [newKey, setNewKey] = useState("");
  const [newValue, setNewValue] = useState("");

  const updateEntry = (oldKey: string, field: "key" | "value", raw: string) => {
    const next: Record<string, number> = {};
    for (const [k, v] of entries) {
      if (k === oldKey) {
        if (field === "key") {
          const trimmed = raw.trim();
          if (trimmed) next[trimmed] = v;
        } else {
          const parsed = parseFloat(raw);
          next[k] = isNaN(parsed) ? 0 : parsed;
        }
      } else {
        next[k] = v;
      }
    }
    onChange(next);
  };

  const removeEntry = (key: string) => {
    const next: Record<string, number> = {};
    for (const [k, v] of entries) {
      if (k !== key) next[k] = v;
    }
    onChange(next);
  };

  const addEntry = () => {
    const trimmedKey = newKey.trim();
    if (!trimmedKey) return;
    const parsed = parseFloat(newValue);
    const numVal = isNaN(parsed) ? 0 : parsed;
    if (value[trimmedKey] !== undefined) return; // don't overwrite
    onChange({ ...value, [trimmedKey]: numVal });
    setNewKey("");
    setNewValue("");
  };

  const listId = label.toLowerCase().replace(/\s+/g, "-") + "-suggestions";

  return (
    <div>
      <label className="text-text-muted text-xs block mb-1" title={tooltip}>
        {label}
      </label>
      {entries.length > 0 && (
        <div className="space-y-1 mb-2">
          {entries.map(([k, v]) => (
            <div key={k} className="flex gap-2 items-center">
              <input
                type="text"
                list={suggestions ? listId : undefined}
                value={k}
                onChange={(e) => updateEntry(k, "key", e.target.value)}
                className="flex-1 px-2 py-1 bg-surface border border-border rounded text-sm text-text focus:outline-none focus:border-primary"
                placeholder={keyPlaceholder}
              />
              <input
                type="text"
                inputMode="decimal"
                value={Number.isFinite(v) ? String(v) : "0"}
                onChange={(e) => updateEntry(k, "value", e.target.value)}
                className="w-20 px-2 py-1 bg-surface border border-border rounded text-sm text-text focus:outline-none focus:border-primary"
                placeholder="0"
              />
              <span className="text-xs text-text-muted w-4">{valueLabel}</span>
              <button
                type="button"
                onClick={() => removeEntry(k)}
                className="text-danger hover:text-danger-hover text-sm px-1"
                aria-label={`Remove ${k}`}
              >
                ×
              </button>
            </div>
          ))}
          {suggestions && (
            <datalist id={listId}>
              {suggestions.map((s) => (
                <option key={s} value={s} />
              ))}
            </datalist>
          )}
        </div>
      )}
      <div className="flex gap-2 items-center">
        <input
          type="text"
          list={suggestions ? listId : undefined}
          value={newKey}
          onChange={(e) => setNewKey(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === "Enter") {
              e.preventDefault();
              addEntry();
            }
          }}
          className="flex-1 px-2 py-1 bg-surface border border-border rounded text-sm text-text focus:outline-none focus:border-primary"
          placeholder={keyPlaceholder}
        />
        <input
          type="text"
          inputMode="decimal"
          value={newValue}
          onChange={(e) => setNewValue(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === "Enter") {
              e.preventDefault();
              addEntry();
            }
          }}
          className="w-20 px-2 py-1 bg-surface border border-border rounded text-sm text-text focus:outline-none focus:border-primary"
          placeholder="0"
        />
        <span className="text-xs text-text-muted w-4">{valueLabel}</span>
        <button
          type="button"
          onClick={addEntry}
          className="text-primary hover:text-primary-hover text-sm px-2 py-1 border border-border rounded"
        >
          + Add
        </button>
      </div>
    </div>
  );
}