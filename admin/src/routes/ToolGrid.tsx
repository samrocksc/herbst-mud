import { Link } from '@tanstack/react-router';

const TOOLS = [
  { to: '/map', emoji: '🗺️', title: 'Map Builder', desc: 'View and edit room layout, connections, and z-levels' },
  { to: '/npcs', emoji: '👤', title: 'NPC Manager', desc: 'Create, edit, and manage NPCs and their locations' },
  { to: '/items', emoji: '📦', title: 'Item Manager', desc: 'Create, edit, and manage items and equipment' },
  { to: '/export', emoji: '💾', title: 'Export / Import', desc: 'Backup and restore game world data' },
  { to: '/players', emoji: '🎮', title: 'Player Manager', desc: 'Manage players and reset passwords' },
  { to: '/characters', emoji: '🧙', title: 'Character Manager', desc: 'View and edit characters, rooms, stats, and bind points' },
  { to: '/abilities', emoji: '⚡', title: 'Abilities Manager', desc: 'Create, edit, and manage abilities' },
  { to: '/quests', emoji: '📜', title: 'Quest Manager', desc: 'Create, edit, and manage quests and objectives' },
  { to: '/logs', emoji: '📋', title: 'Log Viewer', desc: 'View and filter application logs with live tail' },
  { to: '/skills', emoji: '🎯', title: 'Skills Manager', desc: 'Manage trainable skill specializations' },
  { to: '/factions', emoji: '⚔️', title: 'Factions Manager', desc: 'Manage factions, categories, and member standing' },
] as const;

export function ToolGrid() {
  return (
    <div className="grid grid-cols-[repeat(auto-fit,minmax(250px,1fr))] gap-4">
      {TOOLS.map((tool) => (
        <Link
          key={tool.to}
          to={tool.to as any}
          className="block bg-surface-muted rounded-lg p-6 no-underline text-text border border-border transition-colors hover:border-primary hover:bg-surface-muted/70"
        >
          <div className="text-2xl mb-2">{tool.emoji}</div>
          <div className="font-bold mb-1">{tool.title}</div>
          <div className="text-text-muted text-sm">{tool.desc}</div>
        </Link>
      ))}
    </div>
  );
}