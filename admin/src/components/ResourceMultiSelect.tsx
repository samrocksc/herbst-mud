import { useState } from 'react';
import { FieldLabel } from './fields/FieldLabel';
import { ResourceSearchModal } from './ResourceSearchModal';
import { useResourceExists } from './useResourceExists';
import { Button } from './Button';

type ChipProps = Readonly<{
  id: number | string
  resourceType: string
  apiBase: string
  onRemove: () => void
  disabled?: boolean
}>

function ResourceChip({ id, resourceType, apiBase, onRemove, disabled }: ChipProps) {
  const { name } = useResourceExists(resourceType, apiBase, id);
  return (
    <span className="inline-flex items-center gap-1 px-2 py-0.5 bg-primary/10 border border-primary/30 rounded text-sm text-text">
      {name ?? `#${id}`}
      <Button
        variant="ghost"
        size="sm"
        onClick={onRemove}
        disabled={disabled}
        className="text-text-muted hover:text-danger px-0.5 py-0 text-xs"
        aria-label={`Remove ${name ?? id}`}
      >
        ✕
      </Button>
    </span>
  );
}

type Props = Readonly<{
  label: string
  value: (number | string)[]
  onChange: (ids: (number | string)[]) => void
  resourceType: string
  apiBase: string
  tooltip?: string
  disabled?: boolean
}>

export function ResourceMultiSelect({
  label, value, onChange, resourceType, apiBase, tooltip, disabled,
}: Props) {
  const [modalOpen, setModalOpen] = useState(false);

  const handleSelect = (id: number | string, _name: string) => {
    if (!value.includes(id)) {
      onChange([...value, id]);
    }
  };

  const handleRemove = (id: number | string) => {
    onChange(value.filter((v) => v !== id));
  };

  return (
    <div>
      <FieldLabel tooltip={tooltip}>{label}</FieldLabel>
      {value.length > 0 && (
        <div className="flex flex-wrap gap-1.5 mb-2">
          {value.map((id) => (
            <ResourceChip
              key={id}
              id={id}
              resourceType={resourceType}
              apiBase={apiBase}
              onRemove={() => handleRemove(id)}
              disabled={disabled}
            />
          ))}
        </div>
      )}
      <Button
        variant="secondary"
        size="sm"
        onClick={() => setModalOpen(true)}
        disabled={disabled}
      >
        + Add
      </Button>
      <ResourceSearchModal
        isOpen={modalOpen}
        onClose={() => setModalOpen(false)}
        title={`Search ${label}`}
        resourceType={resourceType}
        apiBase={apiBase}
        multi
        selectedIds={value}
        onSelect={handleSelect}
      />
    </div>
  );
}