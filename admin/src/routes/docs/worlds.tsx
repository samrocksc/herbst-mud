import { createFileRoute } from "@tanstack/react-router";
import { PageHeader } from "../../components/PageHeader";

export const Route = createFileRoute("/docs/worlds")({
  component: WorldsDoc,
});

function Section({ title, children }: Readonly<{ title: string; children: React.ReactNode }>) {
  return (
    <section className="mb-8">
      <h2 className="text-lg font-semibold text-text mb-3 pb-2 border-b border-border">{title}</h2>
      {children}
    </section>
  );
}

function InfoBox({ children }: Readonly<{ children: React.ReactNode }>) {
  return (
    <div className="bg-primary/10 border border-primary/30 rounded-lg p-4 mb-4 text-sm">
      {children}
    </div>
  );
}

function Table({
  headers,
  rows,
}: Readonly<{ headers: string[]; rows: (string | React.ReactNode)[][] }>) {
  return (
    <div className="overflow-x-auto mb-4">
      <table className="w-full text-sm border border-border rounded-lg">
        <thead>
          <tr className="bg-surface-muted">
            {headers.map((h) => (
              <th key={h} className="text-left px-3 py-2 font-semibold border-b border-border">{h}</th>
            ))}
          </tr>
        </thead>
        <tbody>
          {rows.map((row, i) => (
            <tr key={i} className="border-b border-border last:border-0">
              {row.map((cell, j) => (
                <td key={j} className="px-3 py-2">{cell}</td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

function WorldsDoc() {
  return (
    <div className="management-page max-w-4xl">
      <PageHeader title="Worlds" backTo="/docs" />

      <InfoBox>
        <strong>TL;DR:</strong> A <strong>world</strong> is a self-contained MUD instance with its own content
        (rooms, NPCs, items, quests). Multiple worlds can run on the same server with separate databases.
        Each world has its own setting overrides (PvP, XP rate, permadeath).
      </InfoBox>

      <Section title="What is a World?">
        <p className="text-text-muted mb-3">
          In Herbst MUD, a world is a complete game environment with its own:
        </p>
        <ul className="text-sm text-text-muted space-y-1 mb-3 list-disc pl-5">
          <li><strong>Content files</strong> — rooms, NPCs, items, quests defined in YAML</li>
          <li><strong>Database tables</strong> — separate prefix for characters, game state</li>
          <li><strong>Game settings</strong> — PvP, XP rates, permadeath rules</li>
          <li><strong>Feature flags</strong> — enable/disable systems (magic, hacking, airships)</li>
        </ul>
        <p className="text-text-muted">
          Worlds are configured in <code className="bg-surface-dark text-text-inverse px-1 rounded">content/worlds.yaml</code>.
          The admin panel shows all worlds, but you typically manage content for one world at a time.
        </p>
      </Section>

      <Section title="World Lifecycle">
        <Table
          headers={["Status", "Meaning", "Accessibility"]}
          rows={[
            ["active", "Ready for players", "Players can connect, create characters"],
            ["development", "Work in progress", "Admins can test, players blocked"],
            ["maintenance", "Temporarily offline", "No access, updates in progress"],
          ]}
        />
        <p className="text-text-muted mt-3">
          Only <strong>active</strong> worlds accept player connections. Development worlds
          are for internal testing only.
        </p>
      </Section>

      <Section title="Multi-World Architecture">
        <p className="text-text-muted mb-3">
          Each entity (Room, NPCTemplate, Item, etc.) has a <code className="bg-surface-dark text-text-inverse px-1 rounded">world_id</code> field.
          This allows the same server to serve multiple worlds with isolated content:
        </p>
        <div className="bg-surface-muted rounded-lg p-4 font-mono text-sm mb-3">
          {`default/          → Herbst MUD (fantasy)
  rooms/         → 50+ fantasy rooms
  npcs/          → goblins, elves, dragons
  cyberpunk/     → Neon Nights (development)
  steampunk/     → Steam & Gear (development)`}
        </div>
        <p className="text-text-muted">
          Content is loaded from <code className="bg-surface-dark text-text-inverse px-1 rounded">content/[world_id]/</code> at startup.
          Database tables use the <code className="bg-surface-dark text-text-inverse px-1 rounded">database_prefix</code> to avoid conflicts.
        </p>
      </Section>

      <Section title="World Settings">
        <p className="text-text-muted mb-3">
          Each world has its own settings overrides in <code className="bg-surface-dark text-text-inverse px-1 rounded">content/worlds.yaml</code>:
        </p>
        <Table
          headers={["Setting", "Description", "Example"]}
          rows={[
            ["pvp_enabled", "Allow player-vs-player combat in this world", "true / false"],
            ["permadeath", "Characters cannot respawn after death", "false (classic)"],
            ["xp_multiplier", "Scale XP gains (1.5 = 50% more XP)", "1.0 / 1.5 / 2.0"],
            ["gold_multiplier", "Scale gold drops and rewards", "0.8 / 1.0 / 1.2"],
            ["features", "List of enabled game systems", "[\"magic_system\", \"loot\"]"],
          ]}
        />
      </Section>

      <Section title="Game Time">
        <p className="text-text-muted mb-3">
          World time runs independently of real time, controlled by <code className="bg-surface-dark text-text-inverse px-1 rounded">world_time_scale</code>:
        </p>
        <Table
          headers={["Scale", "1 Real Second =", "Use Case"]}
          rows={[
            ["1", "1 game minute", "Slow, atmospheric exploration"],
            ["4", "4 game minutes", "Default. Moderate pacing"],
            ["10", "10 game minutes", "Fast-paced combat, events"],
            ["60", "1 game hour", "Server stress testing"],
          ]}
        />
        <p className="text-text-muted mt-3">
          The <code className="bg-surface-dark text-text-inverse px-1 rounded">world_name</code> config key sets the display name shown to players
          in the SSH banner and in-game.
        </p>
      </Section>

      <Section title="Cross-World Features">
        <p className="text-text-muted mb-3">
          Some features are shared across all worlds (defined in <code className="bg-surface-dark text-text-inverse px-1 rounded">shared:</code>):
        </p>
        <ul className="text-sm text-text-muted space-y-1 mb-3 list-disc pl-5">
          <li><strong>User accounts</strong> — single login works across all worlds</li>
          <li><strong>Achievements</strong> — unlocked in one world, visible everywhere</li>
          <li><strong>Admin panel</strong> — manage all worlds from one interface</li>
        </ul>
        <p className="text-text-muted">
          This means players can level an alt in the fast-XP world while keeping their achievements.
        </p>
      </Section>

      <Section title="World Content Structure">
        <p className="text-text-muted mb-3">
          Each world lives in <code className="bg-surface-dark text-text-inverse px-1 rounded">content/[world_id]/</code>:
        </p>
        <div className="bg-surface-muted rounded-lg p-4 font-mono text-sm mb-3">
          {`[world_id]/
  rooms/          # Room definitions (room-*.yaml)
  npcs/           # NPC templates (npc-*.yaml)
  items/          # Item definitions (item-*.yaml)
  quests/         # Quest scripts (quest-*.yaml)
  dialogs/        # NPC conversation trees
  loot_tables/    # Drop probability tables
  spawns/         # Spawn rules and schedules`}
        </div>
        <p className="text-text-muted">
          Content is loaded by <code className="bg-surface-dark text-text-inverse px-1 rounded">content/manager.go</code> at server startup.
          Changes to YAML files require a server restart or content reload.
        </p>
      </Section>
    </div>
  );
}