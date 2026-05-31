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
        <strong>TL;DR:</strong> Tags are key-value pairs you can stick on any character. Use them to
        track quest progress, faction membership, feature flags, or anything else you need. These are
        character-specific, not the same as item tags.
      </InfoBox>

      <Section title="CharacterTag Entity">
        <p className="text-text-muted mb-3">
          Each tag has three fields: who it belongs to, what it is called, and what its value is.
        </p>
        <Table
          headers={["Field", "What it does"]}
          rows={[
            ["character_id", "Which character this tag belongs to."],
            ["tag_key", "The name of the tag (like faction_rank or quest_stage)."],
            ["tag_value", "The value of the tag (like captain or 3)."],
          ]}
        />
        <p className="text-text-muted mt-3">
          A character can have multiple tags with the same key. For example, a character could have
          several quest markers under the same key. The unique constraint is on the combination of
          (character_id, tag_key, tag_value), so you will not get accidental duplicates.
        </p>
      </Section>

      <Section title="Common Use Cases">
        <p className="text-text-muted mb-3">
          Here are some ways you might use tags in your game:
        </p>
        <Table
          headers={["Use Case", "Key Example", "Value Example"]}
          rows={[
            ["Track which faction a character joined", "faction_joined", "iron_guild"],
            ["Track their rank within that faction", "faction_rank", "captain"],
            ["Remember where they are in a quest", "quest_<id>_stage", "3"],
            ["Mark a quest as finished", "completed_quest", "42"],
            ["Toggle a feature on for a specific player", "feature_enabled", "crafting_v2"],
            ["Give someone a special title", "title", "dragon_slayer"],
          ]}
        />
      </Section>

      <Section title="Quest Integration">
        <p className="text-text-muted mb-3">
          Tags are a natural fit for quest progress. When a quest objective advances, the game sets or
          updates a tag on the character. This is handy for things like:
        </p>
        <ul className="list-disc pl-6 text-text-muted mb-4 space-y-1">
          <li>Remembering which stage of a multi-stage quest the character is on.</li>
          <li>Counting how many NPCs a character has defeated or items they have collected.</li>
          <li>Marking a quest as complete so it does not trigger again for repeatable quests.</li>
        </ul>
      </Section>

      <Section title="Admin UI">
        <p className="text-text-muted mb-3">
          You can add, view, and remove tags from the character detail page at /characters/:id/tags.
        </p>
        <Table
          headers={["Method", "Endpoint", "What it does"]}
          rows={[
            ["GET", "/api/characters/:id/tags", "See all tags on a character."],
            ["POST", "/api/characters/:id/tags", "Add a tag. Send {tag_key, tag_value} in the body."],
            ["DELETE", "/api/characters/:id/tags/:tagKey", "Remove a tag by its key."],
            ["GET", "/api/characters/:id/tags?tag_key=X", "Filter tags by key. Handy if a character has a lot of them."],
          ]}
        />
      </Section>
    </div>
  );
}