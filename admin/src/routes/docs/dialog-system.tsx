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
        <strong>TL;DR:</strong> NPCs talk through dialog trees. Each tree is made of nodes, and each
        node has a line the NPC says plus choices for the player. You can also attach effects that
        fire when a node is reached. Players navigate the tree with the talk command.
      </InfoBox>

      <Section title="DialogNode Entity">
        <p className="text-text-muted mb-3">
          A DialogNode is one step in an NPC's conversation. Fill in these fields:
        </p>
        <Table
          headers={["Field", "What it means"]}
          rows={[
            ["npc_template_id", "Which NPC template this dialog node belongs to."],
            ["order", "Where this node sits in the conversation. Lower numbers come first. Players walk through nodes in order."],
            ["speaker_line", "What the NPC says when the player reaches this node."],
            ["player_options", "A JSON array of choices the player can pick. Each choice has display text and a target node to jump to."],
            ["on_enter_effects", "A JSON array of effects that trigger automatically when this node is reached."],
          ]}
        />
      </Section>

      <Section title="How Tree Walking Works">
        <p className="text-text-muted mb-3">
          When a player uses the <code>talk</code> command on an NPC, the game looks for that NPC's
          first dialog node (the one with the lowest order number) that the player hasn't seen yet.
          Each node can offer choices that branch to different nodes, so you can build real
          back-and-forth conversations.
        </p>
        <Table
          headers={["Step", "What happens"]}
          rows={[
            ["Start", "Find the lowest-order node the player hasn't completed yet."],
            ["Display", "Show the NPC's speaker_line to the player."],
            ["Options", "Show the player_options as numbered choices they can pick."],
            ["Branch", "The player picks an option. The conversation jumps to the target node order."],
            ["End", "If a node has no player_options, the conversation ends there."],
          ]}
        />
      </Section>

      <Section title="Player Options Structure">
        <p className="text-text-muted mb-3">
          Each choice in the player_options JSON array has two fields:
        </p>
        <Table
          headers={["Field", "What it means"]}
          rows={[
            ["text", "The text the player sees for this choice."],
            ["target_order", "The order number of the node to jump to when the player picks this choice."],
          ]}
        />
        <p className="text-text-muted">
          By pointing different choices at different target_order values, you can create branching
          conversations where the player's choices matter.
        </p>
      </Section>

      <Section title="On-Enter Effects">
        <p className="text-text-muted mb-3">
          The on_enter_effects field lets you trigger abilities or effects the moment a player
          reaches a dialog node. This is useful for all kinds of things: starting a quest when an
          NPC offers one, changing a character's reputation, or applying a buff through a
          conversation.
        </p>
      </Section>

      <Section title="Player Commands">
        <Table
          headers={["Command", "What it does"]}
          rows={[
            ["talk", "Start or continue a conversation with the nearest NPC."],
          ]}
        />
      </Section>

      <Section title="Admin UI">
        <p className="text-text-muted">
          You manage dialog nodes on the NPC template detail page in the admin panel. Each NPC
          template shows its dialog nodes sorted by order, so you can see the conversation flow
          at a glance.
        </p>
      </Section>
    </div>
  );
}