/* eslint-disable functional/no-mixed-types */
import { useState } from "react";
import { FieldLabel, INPUT_CLASS } from "./fields/FieldLabel";
import { ResourceSearchModal } from "./ResourceSearchModal";
import { useResourceExists } from "./useResourceExists";

type Props = Readonly<{
  label: string
  value: number | string | null | undefined
  onChange: (id: number | string | null) => void
  resourceType: string
  apiBase: string
  tooltip?: string
  disabled?: boolean
}>

export function ResourceIdField({
  label, value, onChange, resourceType, apiBase, tooltip, disabled,
}: Props) {
  const [modalOpen, setModalOpen] = useState(false);
  const { exists, name, isValidating } = useResourceExists(resourceType, apiBase, value);
  const fieldId = label.toLowerCase().replace(/\s+/g, "-");
  const displayValue = value != null && value !== "" ? String(value) : "";

  const statusIndicator = isValidating ? (
    <span className="text-text-muted text-xs">...</span>
  ) : exists === true && name ? (
    <span className="text-success text-xs truncate max-w-40" title={name}>&quot;{name}&quot; ✓</span>
  ) : exists === false ? (
    <span className="text-danger text-xs">Not found</span>
  ) : null;

  const borderColor = exists === false ? "border-danger" : "border-border";

  return (
    <div>
      <FieldLabel htmlFor={fieldId} tooltip={tooltip}>{label}</FieldLabel>
      <div className="flex items-center gap-2">
        <input
          id={fieldId}
          type="text"
          inputMode="numeric"
          value={displayValue}
          onChange={(e) => {
            const raw = e.target.value;
            if (raw === "") { onChange(null); return; }
            const parsed = parseInt(raw, 10);
            if (!isNaN(parsed)) onChange(parsed);
          }}
          placeholder={`Enter ${resourceType} ID`}
          disabled={disabled}
          className={`${INPUT_CLASS.replace("border-border", borderColor)} w-28`}
        />
        <button
          type="button"
          onClick={() => setModalOpen(true)}
          disabled={disabled}
          className="px-2 py-1.5 text-sm bg-surface border border-border rounded hover:bg-surface-hover text-text-muted"
          title="Search"
        >
          🔍
        </button>
        {statusIndicator}
      </div>
      <ResourceSearchModal
        isOpen={modalOpen}
        onClose={() => setModalOpen(false)}
        title={`Search ${label}`}
        resourceType={resourceType}
        apiBase={apiBase}
        value={value}
        onSelect={(id) => onChange(id)}
      />
    </div>
  );
}