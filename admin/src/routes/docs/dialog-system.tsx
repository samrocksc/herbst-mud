import { createFileRoute } from "@tanstack/react-router";
import { PageHeader } from "../../components/PageHeader";

export const Route = createFileRoute("/docs/dialog-system")({
  component: DialogSystemDoc,
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

function DialogSystemDoc() {
  return (
    <div className="management-page max-w-4xl">
      <PageHeader title="Dialog System" backTo="/docs" />

      <InfoBox>
        <strong>TL;DR:</strong> NPCs have dialog trees made of nodes. Each node has a speaker line,
        player options (as branches), and optional on-enter effects. Players walk the tree with
        the talk command.
      </InfoBox>

      <Section title="DialogNode Entity">
        <Table
          headers={["Field", "Description"]}
          rows={[
            ["npc_template_id", "Which NPC template this node belongs to."],
            ["order", "Sort order within the NPC's dialog. Walking is sequential."],
            ["speaker_line", "What the NPC says when this node is reached."],
            ["player_options", "JSON array of options the player can choose. Each has text and target node order."],
            ["on_enter_effects", "JSON array of effects triggered when the node is entered (future use)."],
          ]}
        />
      </Section>

      <Section title="Tree Walking">
        <p className="text-text-muted mb-3">
          When a player uses the <code>talk</code> command on an NPC, the game finds the NPC's
          first dialog node (lowest order) that the player hasn't completed. Player options in
          each node branch to specific target node orders, allowing branching conversations.
        </p>
        <Table
          headers={["Step", "Behavior"]}
          rows={[
            ["Start", "Find lowest-order node not yet shown."],
            ["Display", "Show the speaker_line."],
            ["Options", "Show player_options as numbered choices."],
            ["Branch", "Player picks option, advances to target node order."],
            ["End", "Node with no player_options ends the conversation."],
          ]}
        />
      </Section>

      <Section title="Player Options Structure">
        <p className="text-text-muted mb-3">
          Each option in player_options JSON:
        </p>
        <Table
          headers={["Field", "Description"]}
          rows={[
            ["text", "What the player sees as the choice text."],
            ["target_order", "Which node order to advance to when chosen."],
          ]}
        />
      </Section>

      <Section title="On-Enter Effects">
        <p className="text-text-muted mb-3">
          The on_enter_effects field can trigger abilities or effects when a dialog node is reached.
          This supports quest triggers, reputation changes, or buff application via dialog.
        </p>
      </Section>

      <Section title="Player Commands">
        <Table
          headers={["Command", "Description"]}
          rows={[
            ["talk", "Start or continue a conversation with the nearest NPC."],
          ]}
        />
      </Section>

      <Section title="Admin UI">
        <p className="text-text-muted">
          Dialog nodes are managed on the NPC template detail page in the admin panel.
          Each NPC template shows its dialog nodes sorted by order.
        </p>
      </Section>
    </div>
  );
}
