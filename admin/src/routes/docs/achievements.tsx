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
        <strong>TL;DR:</strong> Achievements are rewards you set up for players who accomplish
        specific things in the game. When a character does something that matches the criteria you
        defined, the achievement unlocks. They get XP and a badge on their profile. You manage them
        at /achievements.
      </InfoBox>

      <Section title="Achievement Entity">
        <p className="text-text-muted mb-3">
          Each achievement has a name, description, optional icon, XP reward, and criteria. Here is
          what every field controls:
        </p>
        <Table
          headers={["Field", "What it does"]}
          rows={[
            ["name", "A unique internal name for this achievement (like first_blood or dragon_slayer)."],
            ["description", "Text that tells players what they did to earn this. This is what they see on their profile."],
            ["icon", "An optional emoji or icon. Shows up next to the achievement on the character page."],
            ["xp_reward", "How much XP the character gets when this achievement unlocks."],
            ["criteria", "A JSON object that defines what triggers the achievement. See the Criteria section below."],
          ]}
        />
      </Section>

      <Section title="Criteria">
        <p className="text-text-muted mb-3">
          Criteria tell the game what has to happen for the achievement to unlock. When an in-game
          event matches the criteria, the game grants the achievement automatically.
        </p>
        <Table
          headers={["Criteria Field", "What it means"]
          }
          rows={[
            ["event", "The type of event that counts. Options are: kill_npc, complete_quest, explore_room, or reach_level."],
            ["target_id", "The specific thing to match against. For example, if the event is kill_npc, this would be the NPC template ID."],
            ["count", "How many times the event has to happen before the achievement unlocks."],
            ["world_id", "Optional. Restricts the achievement to a specific world. Leave blank if it should count everywhere."],
          ]}
        />
      </Section>

      <Section title="Rewards">
        <p className="text-text-muted mb-3">
          When a character unlocks an achievement, here is what they get:
        </p>
        <ul className="list-disc pl-6 text-text-muted mb-4 space-y-1">
          <li>The XP reward is added to their total XP.</li>
          <li>The achievement is recorded on their character.</li>
          <li>If you set an icon, it shows up on their profile page.</li>
        </ul>
      </Section>

      <Section title="Admin UI">
        <p className="text-text-muted mb-3">
          You can create, edit, and delete achievements at /achievements. The criteria field is
          edited as raw JSON, so make sure the structure matches what the game expects.
        </p>
        <Table
          headers={["Method", "Endpoint", "What it does"]}
          rows={[
            ["GET", "/api/achievements", "List all achievements."],
            ["POST", "/api/achievements", "Create a new achievement."],
            ["GET", "/api/achievements/:id", "Look up a specific achievement by ID."],
            ["PUT", "/api/achievements/:id", "Update an existing achievement."],
            ["DELETE", "/api/achievements/:id", "Delete an achievement."],
            ["GET", "/api/characters/:id/achievements", "See which achievements a character has unlocked."],
          ]}
        />
      </Section>
    </div>
  );
}