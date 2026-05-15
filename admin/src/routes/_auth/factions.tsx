import { createFileRoute } from '@tanstack/react-router';
import { useEffect, useState, useCallback } from 'react';
import { showToast } from '../../components/Toast';
import { apiGet, apiPost } from '../../utils/apiFetch';
import { FactionDetail } from './FactionDetail';
import { CreateFactionView } from './CreateFactionView';
import { CategoryManager, CreateCategoryForm } from './CategoryManager';
import { FactionSidebar } from './FactionSidebar';
import { CategorySidebar } from './CategorySidebar';
import { TabBar } from './FactionsTabBar';
import { EMPTY_FORM, type Faction, type FactionCategory, type FactionForm } from './factionTypes';

export const Route = createFileRoute('/_auth/factions')({
  component: FactionsManagement,
});

function FactionsManagement() {
  const [factions, setFactions] = useState<Faction[]>([]);
  const [categories, setCategories] = useState<FactionCategory[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedFaction, setSelectedFaction] = useState<Faction | null>(null);
  const [showCreateForm, setShowCreateForm] = useState(false);
  const [showCreateCategory, setShowCreateCategory] = useState(false);
  const [tab, setTab] = useState<'factions' | 'categories'>('factions');
  const [form, setForm] = useState<FactionForm>(EMPTY_FORM);
  const [createError, setCreateError] = useState('');
  const [saving, setSaving] = useState(false);

  const refresh = useCallback(async () => {
    try {
      const [f, c] = await Promise.all([
        apiGet<Faction[]>('/api/factions'),
        apiGet<FactionCategory[]>('/api/faction-categories'),
      ]);
      setFactions(Array.isArray(f) ? f : []);
      setCategories(Array.isArray(c) ? c : []);
    } catch { showToast('Failed to load data', 'error'); }
  }, []);

  useEffect(() => { refresh().then(() => setLoading(false)); }, [refresh]);

  const filteredFactions = factions.filter((f) =>
    f.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
    (f.description?.toLowerCase().includes(searchQuery.toLowerCase()) ?? false),
  );

  const handleCreate = async () => {
    if (!form.name) { setCreateError('Faction name is required'); return; }
    setSaving(true); setCreateError('');
    try {
      await apiPost('/api/factions', { ...form, display_name: form.display_name || form.name, category_id: form.category_id || null });
      showToast('Faction created', 'success');
      setForm(EMPTY_FORM); setShowCreateForm(false); refresh();
    } catch (err) {
      setCreateError(err instanceof Error ? err.message : 'Failed to create faction');
      showToast('Failed to create faction', 'error');
    } finally { setSaving(false); }
  };

  if (loading) return <div className="p-8 text-text">Loading...</div>;

  return (
    <div className="flex h-screen bg-surface">
      <div className="w-[280px] bg-surface-muted border-r border-border flex flex-col">
        <TabBar tab={tab} setTab={setTab} />
        {tab === 'factions' && (
          <FactionSidebar factions={factions} searchQuery={searchQuery} setSearchQuery={setSearchQuery}
            filteredFactions={filteredFactions} selectedFaction={selectedFaction}
            setSelectedFaction={setSelectedFaction} setShowCreateForm={setShowCreateForm} setForm={setForm} />
        )}
        {tab === 'categories' && <CategorySidebar categories={categories} setShowCreateCategory={setShowCreateCategory} />}
      </div>
      <div className="flex-1 overflow-y-auto p-6">
        {tab === 'factions' && showCreateForm && (
          <CreateFactionView form={form} setForm={setForm} categories={categories}
            createError={createError} saving={saving} onCreate={handleCreate} onCancel={() => setShowCreateForm(false)} />
        )}
        {tab === 'factions' && !showCreateForm && selectedFaction && (
          <FactionDetail faction={selectedFaction} categories={categories} onRefresh={refresh} />
        )}
        {tab === 'factions' && !showCreateForm && !selectedFaction && (
          <div className="flex items-center justify-center h-full text-text-muted">Select a faction or create a new one</div>
        )}
        {tab === 'categories' && showCreateCategory && (
          <CreateCategoryForm onDone={() => { setShowCreateCategory(false); refresh(); }} />
        )}
        {tab === 'categories' && !showCreateCategory && <CategoryManager categories={categories} />}
      </div>
    </div>
  );
}