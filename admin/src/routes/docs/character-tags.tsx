import { createFileRoute } from "@tanstack/react-router";
import { PageHeader } from "../../components/PageHeader";

export const Route = createFileRoute("/docs/character-tags")({
  component: CharacterTagsDoc,
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

function CharacterTagsDoc() {
  return (
    <div className="management-page max-w-4xl">
      <PageHeader title="Character Tags" backTo="/docs" />

      <InfoBox>
        <strong>TL;DR:</strong> Tags are key-value metadata attached to characters.
        Used for quest markers, faction tracking, feature flags, and game state.
        Not to be confused with item tags — these are per-character.
      </InfoBox>

      <Section title="CharacterTag Entity">
        <Table
          headers={["Field", "Description"]}
          rows={[
            ["character_id", "Which character owns this tag."],
            ["tag_key", "The tag name (e.g., faction_rank, quest_stage)."],
            ["tag_value", "The tag value (e.g., captain, 3)."],
          ]}
        />
        <p className="text-text-muted mt-3">
          Multiple tags can share the same key (e.g., multiple quest markers).
          The unique constraint is (character_id, tag_key, tag_value).
        </p>
      </Section>

      <Section title="Common Use Cases">
        <Table
          headers={["Use Case", "Key Example", "Value Example"]}
          rows={[
            ["Faction membership", "faction_joined", "iron_guild"],
            ["Faction rank", "faction_rank", "captain"],
            ["Quest progress", "quest_<id>_stage", "3"],
            ["Quest completion", "completed_quest", "42"],
            ["Feature flags", "feature_enabled", "crafting_v2"],
            ["Custom flags", "title", "dragon_slayer"],
          ]}
        />
      </Section>

      <Section title="Quest Integration">
        <p className="text-text-muted mb-3">
          Quest progress often stores intermediate stages as tags. When a quest
          objective advances, the game sets or updates a tag. This allows:
        </p>
        <ul className="list-disc pl-6 text-text-muted mb-4 space-y-1">
          <li>Tracking which stage of a multi-stage quest a character is on.</li>
          <li>Storing NPC kills or items collected as tag counts.</li>
          <li>Marking quest completion for repeatable quests.</li>
        </ul>
      </Section>

      <Section title="Admin UI">
        <p className="text-text-muted mb-3">
          Tags can be managed from the character detail page at /characters/:id/tags.
        </p>
        <Table
          headers={["Method", "Endpoint", "Description"]}
          rows={[
            ["GET", "/api/characters/:id/tags", "List all tags for a character."],
            ["POST", "/api/characters/:id/tags", "Add a tag (body: {tag_key, tag_value})."],
            ["DELETE", "/api/characters/:id/tags/:tagKey", "Remove a tag by key."],
            ["GET", "/api/characters/:id/tags?tag_key=X", "Filter tags by key."],
          ]}
        />
      </Section>
    </div>
  );
}
