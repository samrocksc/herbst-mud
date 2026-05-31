import { createFileRoute } from "@tanstack/react-router";
import { PageHeader } from "../../components/PageHeader";

export const Route = createFileRoute("/docs/quest-system")({
  component: QuestSystemDoc,
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

function QuestSystemDoc() {
  return (
    <div className="management-page max-w-4xl">
      <PageHeader title="Quest System" backTo="/docs" />

      <InfoBox>
        <strong>TL;DR:</strong> Quests give your players goals to chase. A quest is a list of
        objectives that players complete in order. You can set up kill targets, exploration
        checkpoints, item collection, NPC conversations, and item deliveries. The game tracks
        progress automatically as players play.
      </InfoBox>

      <Section title="Quest Lifecycle">
        <p className="text-text-muted mb-3">
          Every quest goes through a set of states for each character:
        </p>
        <Table
          headers={["State", "What it means"]}
          rows={[
            ["Available", "The player can see the quest and hasn't accepted it yet."],
            ["Active", "The player accepted the quest and is working through objectives."],
            ["Completed", "All objectives are done. Rewards have been handed out."],
            ["Abandoned", "The player gave up. All progress on this quest is lost."],
          ]}
        />
      </Section>

      <Section title="Objective Types">
        <p className="text-text-muted mb-3">
          Each quest has a list of objectives that players tackle in order. A later objective
          won't start counting until the one before it is finished. This lets you build
          multi-step quest lines that feel natural.
        </p>
        <Table
          headers={["Type", "What to target", "How it works"]}
          rows={[
            ["kill", "NPC template ID", "Defeat a specific kind of NPC. The count tracks how many the player has killed."],
            ["explore", "Room ID", "Visit a specific room. Completes automatically when the player walks in."],
            ["collect", "Item template ID", "Gather items. The count tracks how many the player has picked up."],
            ["talk", "NPC template ID", "Have a conversation with a specific NPC. Completes automatically on dialog."],
            ["return", "Item template ID", "Bring an item back to the quest giver to turn it in."],
          ]}
        />
      </Section>

      <Section title="Repeat Modes">
        <p className="text-text-muted mb-3">
          You can let players repeat a quest, or keep it as a one-time deal:
        </p>
        <Table
          headers={["Mode", "How it works"]}
          rows={[
            ["none", "One-time only. Once the player finishes, they can never pick it up again."],
            ["cooldown", "Players can redo the quest after the cooldown period has passed since their last completion."],
            ["always", "Players can re-accept this quest any time, no cooldown required."],
          ]}
        />
        <p className="text-text-muted">
          Prerequisites are always checked before acceptance. A player can't accept a quest
          they already have active.
        </p>
      </Section>

      <Section title="Rewards">
        <p className="text-text-muted mb-3">
          When a player finishes a quest, you can hand out any of these rewards:
        </p>
        <Table
          headers={["Reward", "What it does"]}
          rows={[
            ["XP", "Experience points added straight to the character."],
            ["Items", "Item templates dropped into the character's inventory."],
            ["Effects", "Abilities or effects applied to the character."],
            ["Tags", "Tags added to or removed from the character (great for quest flags)."],
            ["Achievements", "Achievements unlocked for the character."],
          ]}
        />
        <InfoBox>
          The admin panel currently only exposes the XP reward field. If you need to grant
          items, effects, tags, or achievements as quest rewards, you can set those up
          through the API directly.
        </InfoBox>
      </Section>

      <Section title="Player Commands">
        <Table
          headers={["Command", "What it does"]}
          rows={[
            ["quests / quest / q", "Shows the quest tracker with all active, completed, and abandoned quests."],
            ["quest accept &lt;id&gt;", "Accept a quest by its ID. The game checks prerequisites before allowing this."],
            ["quest abandon &lt;id&gt;", "Abandon an active quest. Warning: all progress is lost."],
          ]}
        />
        <p className="text-text-muted">
          Players don't need to manually track quest progress. The game updates objectives
          automatically as they kill NPCs, enter rooms, pick up items, or talk to NPCs.
        </p>
      </Section>

      <Section title="Admin API">
        <Table
          headers={["Method", "Endpoint", "What it does"]}
          rows={[
            ["GET", "/api/quests", "List all quest definitions."],
            ["POST", "/api/quests", "Create a new quest."],
            ["GET", "/api/quests/:id", "Get a single quest definition by ID."],
            ["PUT", "/api/quests/:id", "Update a quest definition."],
            ["DELETE", "/api/quests/:id", "Delete a quest. This fails if any players have progress records for it."],
            ["GET", "/api/characters/:id/quests", "List quest progress for a character."],
            ["POST", "/api/characters/:id/quests", "Accept a quest on behalf of a character (body: {\"quest_id\": N})."],
            ["PUT", "/api/characters/:id/quests/:qid/check", "Manually check or increment progress on a single quest objective."],
            ["PUT", "/api/characters/:id/quests/:qid/abandon", "Abandon a quest on behalf of a character."],
            ["POST", "/api/characters/:id/quests/check-all", "Bulk check all matching quests for a character at once."],
          ]}
        />
      </Section>
    </div>
  );
}