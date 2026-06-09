import { createFileRoute } from "@tanstack/react-router";
import { useEffect, useState, useCallback } from "react";
import { showToast } from "../../components/Toast";
import { apiGet, apiPost, apiDelete } from "../../utils/apiFetch";
import { DataTable, type Column } from "../../components/DataTable";
import { FactionDetail } from "./FactionDetail";
import { CreateFactionView } from "./CreateFactionView";
import { CategoryManager, CreateCategoryForm, CategoryEditForm } from "./CategoryManager";
import { FactionSidebar } from "./FactionSidebar";
import { CategorySidebar } from "./CategorySidebar";
import { TabBar } from "./FactionsTabBar";
import { EMPTY_FORM, type Faction, type FactionCategory, type FactionForm } from "./-factionTypes";
import { Button } from "../../components/Button";

export const Route = createFileRoute("/_auth/factions")({
  component: FactionsManagement,
});

function FactionsManagement() {
  const [factions, setFactions] = useState<Faction[]>([]);
  const [categories, setCategories] = useState<FactionCategory[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchQuery, setSearchQuery] = useState("");
  const [selectedFaction, setSelectedFaction] = useState<Faction | null>(null);
  const [showCreateForm, setShowCreateForm] = useState(false);
  const [showCreateCategory, setShowCreateCategory] = useState(false);
  const [editingCategory, setEditingCategory] = useState<FactionCategory | null>(null);
  const [tab, setTab] = useState<"factions" | "categories">("factions");
  const [factionView, setFactionView] = useState<"list" | "table">("list");
  const [form, setForm] = useState<FactionForm>(EMPTY_FORM);
  const [createError, setCreateError] = useState("");
  const [saving, setSaving] = useState(false);

  const refresh = useCallback(async () => {
    try {
      const [f, c] = await Promise.all([
        apiGet<Faction[]>("/api/factions"),
        apiGet<FactionCategory[]>("/api/faction-categories"),
      ]);
      setFactions(Array.isArray(f) ? f : []);
      setCategories(Array.isArray(c) ? c : []);
    } catch { showToast("Failed to load data", "error"); }
  }, []);

  useEffect(() => { refresh().then(() => setLoading(false)); }, [refresh]);

  const filteredFactions = factions.filter((f) =>
    f.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
    (f.description?.toLowerCase().includes(searchQuery.toLowerCase()) ?? false),
  );

  const handleCreate = async () => {
    if (!form.name) { setCreateError("Faction name is required"); return; }
    setSaving(true); setCreateError("");
    try {
      await apiPost("/api/factions", { ...form, display_name: form.display_name || form.name, category_id: form.category_id || null });
      showToast("Faction created", "success");
      setForm(EMPTY_FORM); setShowCreateForm(false); refresh();
    } catch (err) {
      setCreateError(err instanceof Error ? err.message : "Failed to create faction");
      showToast("Failed to create faction", "error");
    } finally { setSaving(false); }
  };

  const handleEditCategory = (cat: FactionCategory) => {
    setEditingCategory(cat);
  };

  const handleDeleteCategory = async (cat: FactionCategory) => {
    if (!confirm(`Delete category "${cat.display_name || cat.name}"? This will unlink any factions.`)) return;
    try {
      await apiDelete(`/api/faction-categories/${cat.id}`);
      showToast("Category deleted", "success");
      refresh();
    } catch (err) {
      showToast(err instanceof Error ? err.message : "Failed to delete", "error");
    }
  };

  const factionColumns: Column<Faction>[] = [
    {
      header: "Icon", accessor: "emoji", align: "center",
      render: (_, row) => {
        const emoji = (row as Faction & { emoji?: string }).emoji;
        return <span className="text-base">{emoji ?? ""}</span>;
      },
    },
    {
      header: "Name", accessor: "display_name",
      render: (_, row) => (
        <button type="button" onClick={() => { setSelectedFaction(row); setShowCreateForm(false); }} className="text-primary hover:underline font-bold text-left">
          {row.display_name || row.name}
        </button>
      ),
    },
    { header: "Slug", accessor: "name", render: (v) => <code className="text-xs text-text-muted">{String(v)}</code> },
    {
      header: "Category", accessor: "category_id",
      render: (_, row) => {
        const cat = categories.find((c) => c.id === row.category_id);
        return cat ? <span className="text-xs text-text-muted">{cat.display_name || cat.name}</span> : <span className="text-xs text-text-muted">—</span>;
      },
    },
    { header: "Members", accessor: "member_count", align: "center", render: (v) => String(v ?? 0) },
    {
      header: "Tags", accessor: "member_tags",
      render: (_, row) => (
        <div className="flex flex-wrap gap-1">
          {(row.member_tags ?? []).slice(0, 2).map((t) => (
            <span key={t} className="px-1.5 py-0.5 bg-primary/20 text-primary text-xs rounded">{t}</span>
          ))}
          {(row.member_tags?.length ?? 0) > 2 && <span className="text-xs text-text-muted">+{row.member_tags!.length - 2}</span>}
        </div>
      ),
    },
  ];

  if (loading) return <div className="p-8 text-text">Loading...</div>;

  return (
    <div className="flex h-full min-h-[100dvh] bg-surface">
      <div className="w-[280px] bg-surface-muted border-r border-border flex flex-col shrink-0 overflow-y-auto">
        <TabBar tab={tab} setTab={setTab} />
        {tab === "factions" && factionView === "list" && (
          <FactionSidebar factions={factions} searchQuery={searchQuery} setSearchQuery={setSearchQuery}
            filteredFactions={filteredFactions} selectedFaction={selectedFaction}
            setSelectedFaction={setSelectedFaction} setShowCreateForm={setShowCreateForm} setForm={setForm} />
        )}
        {tab === "factions" && factionView === "table" && (
          <div className="p-3 border-b border-border">
            <input
              type="text"
              placeholder="Search factions..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
            />
          </div>
        )}
        {tab === "categories" && (
          <CategorySidebar
            categories={categories}
            setShowCreateCategory={setShowCreateCategory}
            onEditCategory={handleEditCategory}
            onDeleteCategory={handleDeleteCategory}
          />
        )}
      </div>
      <div className="flex-1 overflow-y-auto p-6">
        {tab === "factions" && (
          <>
            <div className="flex items-center justify-between mb-4">
              <h1 className="text-2xl font-bold text-text m-0">Factions</h1>
              <div className="flex gap-2">
                <div className="flex border border-border rounded overflow-hidden">
                  <button
                    type="button"
                    onClick={() => setFactionView("list")}
                    className={`px-3 py-1 text-xs ${factionView === "list" ? "bg-primary text-white" : "bg-surface text-text hover:bg-surface-hover"}`}
                  >List</button>
                  <button
                    type="button"
                    onClick={() => setFactionView("table")}
                    className={`px-3 py-1 text-xs ${factionView === "table" ? "bg-primary text-white" : "bg-surface text-text hover:bg-surface-hover"}`}
                  >Table</button>
                </div>
                <Button variant="primary" onClick={() => { setShowCreateForm(true); setSelectedFaction(null); setForm(EMPTY_FORM); }}>+ Add Faction</Button>
              </div>
            </div>
            {showCreateForm && (
              <CreateFactionView form={form} setForm={setForm} categories={categories}
                createError={createError} saving={saving} onCreate={handleCreate} onCancel={() => setShowCreateForm(false)} />
            )}
            {!showCreateForm && factionView === "table" && (
              <DataTable<Faction> columns={factionColumns} data={filteredFactions} getKey={(f) => f.id}
                emptyMessage={searchQuery ? "No factions match your search" : "No factions yet. Click + Add Faction to start."} />
            )}
            {!showCreateForm && factionView === "list" && selectedFaction && (
              <FactionDetail faction={selectedFaction} categories={categories} onRefresh={refresh} />
            )}
            {!showCreateForm && factionView === "list" && !selectedFaction && (
              <div className="flex flex-col items-center justify-center h-64 text-text-muted gap-3">
                <p>Select a faction from the sidebar to view details.</p>
                <Button variant="primary" onClick={() => { setShowCreateForm(true); setForm(EMPTY_FORM); }}>+ Create your first faction</Button>
              </div>
            )}
          </>
        )}
        {tab === "categories" && showCreateCategory && (
          <CreateCategoryForm onDone={() => { setShowCreateCategory(false); refresh(); }} />
        )}
        {tab === "categories" && !showCreateCategory && !editingCategory && (
          <CategoryManager categories={categories} onRefresh={refresh} onEdit={handleEditCategory} />
        )}
        {tab === "categories" && editingCategory && (
          <div>
            <button
              type="button"
              onClick={() => setEditingCategory(null)}
              className="text-sm text-text-muted hover:text-primary mb-3"
            >
              ← Back to categories
            </button>
            <CategoryEditForm
              category={editingCategory}
              onDone={() => { setEditingCategory(null); refresh(); }}
            />
          </div>
        )}
      </div>
    </div>
  );
}
