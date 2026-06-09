import { Button } from "../../components/Button";
import type { FactionCategory } from "./-factionTypes";

export function CategorySidebar({
  categories,
  setShowCreateCategory,
  onEditCategory,
  onDeleteCategory,
}: Readonly<{
  categories: FactionCategory[]
  setShowCreateCategory: (v: boolean) => void
  onEditCategory: (c: FactionCategory) => void
  onDeleteCategory: (c: FactionCategory) => void
}>) {
  return (
    <>
      <div className="p-3 border-b border-border flex items-center justify-between">
        <div>
          <h2 className="m-0 text-text text-lg">Categories</h2>
          <p className="text-text-muted text-xs mt-1 mb-0">{categories.length} total</p>
        </div>
        <Button variant="primary" size="sm" onClick={() => setShowCreateCategory(true)}>
          + New
        </Button>
      </div>
      <div className="flex-1 overflow-y-auto p-3">
        <div className="flex flex-col gap-2">
          {categories.map((c) => {
            const memberTagsCount = (c as FactionCategory & { faction_count?: number }).faction_count ?? 0;
            return (
              <div key={c.id} className="p-2 rounded border border-border bg-surface-muted">
                <div className="flex items-start justify-between gap-1">
                  <div className="flex-1 min-w-0">
                    <div className="font-bold text-sm text-text">{c.display_name || c.name}</div>
                    <div className="text-text-muted text-xs">
                      {memberTagsCount > 0 && <span>{memberTagsCount} factions</span>}
                      {memberTagsCount > 0 && c.max_memberships != null && <span> · </span>}
                      {c.max_memberships != null && (
                        <span>max {c.max_memberships === 1 ? "single" : c.max_memberships}</span>
                      )}
                    </div>
                    {c.description && (
                      <div className="text-text-muted text-xs mt-1 italic line-clamp-2">{c.description}</div>
                    )}
                    {c.auto_join && (
                      <span className="inline-block mt-1 px-1.5 py-0.5 rounded text-xs bg-success/20 text-success">
                        Auto-join
                      </span>
                    )}
                  </div>
                </div>
                <div className="flex gap-1 mt-2">
                  <button
                    type="button"
                    onClick={() => onEditCategory(c)}
                    className="text-xs text-text-muted hover:text-primary"
                  >
                    Edit
                  </button>
                  <span className="text-text-muted">·</span>
                  <button
                    type="button"
                    onClick={() => onDeleteCategory(c)}
                    className="text-xs text-text-muted hover:text-danger"
                  >
                    Delete
                  </button>
                </div>
              </div>
            );
          })}
          {categories.length === 0 && (
            <div className="text-center py-6 space-y-2">
              <p className="text-text-muted text-sm">No categories yet</p>
              <Button variant="accent" size="sm" onClick={() => setShowCreateCategory(true)}>
                + Create your first category
              </Button>
            </div>
          )}
        </div>
      </div>
    </>
  );
}
