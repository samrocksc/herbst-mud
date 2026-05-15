export function TabBar({ tab, setTab }: Readonly<{ tab: string; setTab: (t: 'factions' | 'categories') => void }>) {
  return (
    <>
      <div className="p-4 border-b border-border">
        <a href="/dashboard"
          className="block no-underline p-2 rounded border-2 border-black text-center text-sm font-medium bg-surface-muted text-text hover:border-primary transition-colors">
          ← Dashboard
        </a>
      </div>
      <div className="flex p-3 gap-2 border-b border-border">
        {(['factions', 'categories'] as const).map((t) => (
          <button key={t} onClick={() => setTab(t)}
            className={`flex-1 py-1 px-2 rounded text-sm font-medium border ${tab === t ? 'bg-primary text-white border-primary' : 'bg-surface text-text border-border hover:border-primary'}`}>
            {t.charAt(0).toUpperCase() + t.slice(1)}
          </button>
        ))}
      </div>
    </>
  );
}