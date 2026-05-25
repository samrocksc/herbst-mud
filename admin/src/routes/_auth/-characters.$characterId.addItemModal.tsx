import { useState } from "react";
import { useEquipmentTemplates } from "../../hooks/useEquipmentTemplates";
import { Button } from "../../components/Button";

/** Modal for selecting an item template and spawning it into a character's inventory. */
export function AddItemModal({ open, onClose, onSpawn, isLoading, error }: Readonly<{
  open: boolean
  onClose: () => void
  onSpawn: (templateId: string) => void
  isLoading: boolean
  error: string | null
}>) {
  const [search, setSearch] = useState("");
  const { data: templates, isLoading: loadingTemplates } = useEquipmentTemplates();

  if (!open) return null;

  const filtered = (templates ?? []).filter((t) =>
    t.name.toLowerCase().includes(search.toLowerCase())
  );

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h3>Add Item to Inventory</h3>
          <Button variant="ghost" size="sm" onClick={onClose}>x</Button>
        </div>
        <div className="modal-body space-y-3">
          <input
            type="text"
            className="w-full bg-surface border border-border rounded-md px-3 py-2 text-text text-sm"
            placeholder="Search item templates..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            autoFocus
          />
          {loadingTemplates && <div className="text-text-muted text-xs">Loading templates...</div>}
          <div className="max-h-64 overflow-y-auto space-y-1">
            {filtered.map((t) => (
              <button
                key={t.id}
                className="w-full text-left px-3 py-2 rounded text-sm hover:bg-surface-muted text-text"
                onClick={() => onSpawn(String(t.id))}
                disabled={isLoading}
              >
                <span className="font-medium">{t.name}</span>
                <span className="text-text-muted ml-2">{t.slot} · Lvl {t.level}</span>
              </button>
            ))}
            {!loadingTemplates && filtered.length === 0 && (
              <div className="text-text-muted text-xs py-2">No templates found</div>
            )}
          </div>
          {error && <div className="text-danger text-xs">{error}</div>}
        </div>
        <div className="modal-footer">
          <Button variant="secondary" size="sm" onClick={onClose}>Cancel</Button>
        </div>
      </div>
    </div>
  );
}