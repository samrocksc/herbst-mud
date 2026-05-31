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
        (rooms, NPCs, items, quests). All worlds share the same database but are isolated by a
        <code className="bg-surface-dark text-text-inverse px-1 rounded">world_id</code> field on every entity.
      </InfoBox>

      <Section title="What is a World?">
        <p className="text-text-muted mb-3">
          In Herbst MUD, a world is a complete game environment with its own:
        </p>
        <ul className="text-sm text-text-muted space-y-1 mb-3 list-disc pl-5">
          <li><strong>Content</strong> — rooms, NPC templates, items, quests, abilities, crafting recipes</li>
          <li><strong>Characters</strong> — players and NPCs created within that world</li>
          <li><strong>World settings</strong> — name, title, description, active status</li>
        </ul>
        <p className="text-text-muted">
          Create and manage worlds from the <strong>Config → Worlds</strong> page in the admin panel.
          Switch between worlds using the world selector in the top bar.
        </p>
      </Section>

      <Section title="World Lifecycle">
        <Table
          headers={["Status", "Meaning", "Accessibility"]}
          rows={[
            ["active", "Ready for players", "Players can connect and create characters"],
            ["inactive", "Offline", "Players cannot join; admins can still edit content"],
          ]}
        />
        <p className="text-text-muted mt-3">
          Only <strong>active</strong> worlds appear in the character creation screen and world selector.
          Set a world to inactive to close it for maintenance or events.
        </p>
      </Section>

      <Section title="Multi-World Data Isolation">
        <p className="text-text-muted mb-3">
          Every content entity has a <code className="bg-surface-dark text-text-inverse px-1 rounded">world_id</code> field
          that ties it to a specific world. When you query the API, results are automatically filtered
          to the currently selected world.
        </p>
        <p className="text-text-muted mb-3">
          Entities that are world-scoped:
        </p>
        <ul className="text-sm text-text-muted space-y-1 mb-3 list-disc pl-5">
          <li>Rooms, NPC templates, Items, Abilities</li>
          <li>Quests, Crafting recipes, Factions</li>
          <li>Races, Genders, Tags</li>
        </ul>
        <p className="text-text-muted">
          Shared across all worlds: user accounts, achievements, and the admin panel itself.
          A player can have characters in multiple worlds with the same login.
        </p>
      </Section>

      <Section title="World Content">
        <p className="text-text-muted mb-3">
          Content for each world is created and edited entirely through the admin panel or the REST API.
          There is no external file-based configuration. All rooms, NPCs, items, and quests live in the
          database and are scoped by <code className="bg-surface-dark text-text-inverse px-1 rounded">world_id</code>.
        </p>
        <p className="text-text-muted mb-3">
          To build a new world:
        </p>
        <ol className="text-sm text-text-muted space-y-1 mb-3 list-decimal pl-5">
          <li>Create the world in <strong>Config → Worlds</strong></li>
          <li>Select it from the world dropdown in the top bar</li>
          <li>Add rooms, NPC templates, items, and quests through the admin panel</li>
          <li>Set the world to <strong>active</strong> when ready</li>
        </ol>
        <p className="text-text-muted">
          Use the <strong>Export</strong> and <strong>Import</strong> tools in the admin panel to copy
          content between worlds or back up your work.
        </p>
      </Section>

      <Section title="Cross-World Features">
        <p className="text-text-muted mb-3">
          Some systems are shared across all worlds:
        </p>
        <ul className="text-sm text-text-muted space-y-1 mb-3 list-disc pl-5">
          <li><strong>User accounts</strong> — single login works across all worlds</li>
          <li><strong>Achievements</strong> — unlocked in one world, visible everywhere</li>
          <li><strong>Admin panel</strong> — manage all worlds from one interface</li>
          <li><strong>Game config</strong> — global settings like XP thresholds apply to all worlds</li>
        </ul>
      </Section>

      <Section title="Tips">
        <ul className="text-sm text-text-muted space-y-1 mb-3 list-disc pl-5">
          <li>Start with one world. Add more only when you have a complete, tested game loop.</li>
          <li>Use the <strong>default</strong> world for development and testing before creating a production world.</li>
          <li>Export your world content regularly as a backup — the export JSON contains everything needed to recreate the world.</li>
          <li>World names must be unique. Use descriptive names like <code>herbst-fantasy</code> or <code>neon-nights</code>.</li>
        </ul>
      </Section>
    </div>
  );
}
