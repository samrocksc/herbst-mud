import { createFileRoute } from "@tanstack/react-router";
import { PageHeader } from "../../components/PageHeader";

export const Route = createFileRoute("/docs/achievements")({
  component: AchievementsDoc,
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

function AchievementsDoc() {
  return (
    <div className="management-page max-w-4xl">
      <PageHeader title="Achievements" backTo="/docs" />

      <InfoBox>
        <strong>TL;DR:</strong> Achievements are unlocked when character actions match
        defined criteria. They grant XP rewards and display an icon. Managed via /achievements.
      </InfoBox>

      <Section title="Achievement Entity">
        <Table
          headers={["Field", "Description"]}
          rows={[
            ["name", "Unique identifier (e.g., first_blood, dragon_slayer)."],
            ["description", "Player-facing description of how to earn it."],
            ["icon", "Optional emoji or icon for display."],
            ["xp_reward", "XP granted when achievement is unlocked."],
            ["criteria", "Optional JSON defining earn conditions (event-based)."],
          ]}
        />
      </Section>

      <Section title="Criteria">
        <p className="text-text-muted mb-3">
          Criteria is a JSON object that defines what must happen to earn the achievement.
          When events in the game match the criteria, the achievement is granted.
        </p>
        <Table
          headers={["Criteria Field", "Description"]}
          rows={[
            ["event", "Event type: kill_npc, complete_quest, explore_room, reach_level."],
            ["target_id", "Specific ID to match (e.g., NPC template ID)."],
            ["count", "Number of times the event must occur."],
            ["world_id", "Optional: restrict to specific world."],
          ]}
        />
      </Section>

      <Section title="Rewards">
        <p className="text-text-muted mb-3">
          When an achievement is unlocked, the character receives:
        </p>
        <ul className="list-disc pl-6 text-text-muted mb-4 space-y-1">
          <li>XP reward added to their XP total.</li>
          <li>Achievement recorded in their character record.</li>
          <li>Icon displayed in their character profile (if set).</li>
        </ul>
      </Section>

      <Section title="Admin UI">
        <p className="text-text-muted mb-3">
          Achievements are managed at /achievements with full CRUD. Criteria is edited
          as JSON in the form.
        </p>
        <Table
          headers={["Method", "Endpoint", "Description"]}
          rows={[
            ["GET", "/api/achievements", "List all achievements."],
            ["POST", "/api/achievements", "Create an achievement."],
            ["GET", "/api/achievements/:id", "Get achievement by ID."],
            ["PUT", "/api/achievements/:id", "Update an achievement."],
            ["DELETE", "/api/achievements/:id", "Delete an achievement."],
            ["GET", "/api/characters/:id/achievements", "Get achievements for a character."],
          ]}
        />
      </Section>
    </div>
  );
}
