import { Button } from "../../components/Button";
import type { FactionCategory } from "./factionTypes";

export function CategorySidebar({
  categories,
  setShowCreateCategory,
}: Readonly<{
  categories: FactionCategory[]
  setShowCreateCategory: (v: boolean) => void
}>) {
  return (
    <>
      <div className="p-3 border-b border-border">
        <h2 className="m-0 text-text text-lg">Categories</h2>
        <p className="text-text-muted text-xs mt-1 mb-0">{categories.length} categories</p>
      </div>
      <div className="flex-1 overflow-y-auto p-3">
        <div className="flex flex-col gap-1">
          {categories.map((c) => (
            <div key={c.id} className="p-2 rounded text-xs text-text">
              <div className="font-bold">{c.display_name || c.name}</div>
              <div className="text-text-muted text-xs">{c.description || "—"}</div>
            </div>
          ))}
          {categories.length === 0 && <div className="text-text-muted text-center py-4">No categories</div>}
        </div>
      </div>
      <div className="p-3 border-t border-border">
        <Button variant="primary" size="md" fullWidth onClick={() => setShowCreateCategory(true)}>
          + Add Category
        </Button>
      </div>
    </>
  );
}