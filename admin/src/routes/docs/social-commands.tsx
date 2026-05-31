import { createFileRoute } from "@tanstack/react-router";
import { PageHeader } from "../../components/PageHeader";

export const Route = createFileRoute("/docs/social-commands")({
  component: SocialCommandsDoc,
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

function SocialCommandsDoc() {
  return (
    <div className="management-page max-w-4xl">
      <PageHeader title="Social Commands" backTo="/docs" />

      <InfoBox>
        <strong>TL;DR:</strong> Social commands let players talk to each other and interact with
        the world around them. You get five built-in commands (say, shout, whisper, tell, emote)
        and you can create as many custom socials as you want from the admin panel.
      </InfoBox>

      <Section title="Built-In Social Commands">
        <p className="text-text-muted mb-3">
          These commands come with the game and are always available:
        </p>
        <Table
          headers={["Command", "Reach", "What it does"]}
          rows={[
            ["say &lt;text&gt;", "Same room", "Talk to everyone in the same room."],
            ["shout &lt;text&gt;", "Adjacent rooms", "Shout loud enough for nearby rooms to hear."],
            ["whisper &lt;player&gt; &lt;text&gt;", "Same room", "Send a private message to someone in the same room."],
            ["tell &lt;player&gt; &lt;text&gt;", "Global", "Send a private message to anyone in the world."],
            ["emote &lt;text&gt;", "Same room", "Perform an action that everyone in the room can see (e.g., /me smiles)."],
          ]}
        />
      </Section>

      <Section title="Custom Socials">
        <p className="text-text-muted mb-3">
          You can create custom social commands from the admin panel at /socials. Each social
          has these fields:
        </p>
        <Table
          headers={["Field", "What it means"]}
          rows={[
            ["name", "The command players type (e.g., laugh, wave, dance)."],
            ["message_format", "The template for what other people see. Use $sender for the player's name."],
            ["channel", "How far the message reaches: same_room, adjacent, or global."],
            ["cooldown_secs", "Optional cooldown between uses. Set to 0 or leave blank for no cooldown."],
          ]}
        />
        <p className="text-text-muted">
          Custom socials are great for roleplaying. A "bow" command that shows "Aleron bows deeply"
          adds flavor that makes the world feel alive.
        </p>
      </Section>

      <Section title="Channel Configuration">
        <p className="text-text-muted mb-3">
          Chat channels are set up at /channels. Channels control who can hear a message and
          how far it travels. You can restrict channels by class, level, or faction.
        </p>
        <Table
          headers={["Channel Field", "What it does"]}
          rows={[
            ["name", "The channel's identifier (e.g., ooc, auction, guild)."],
            ["range", "How far the message travels: same_room, adjacent, or global."],
            ["min_level", "The minimum character level required to use this channel."],
            ["required_faction", "Optional. Restrict this channel to members of a specific faction."],
          ]}
        />
      </Section>

      <Section title="SendMessage Effect">
        <p className="text-text-muted mb-3">
          The SendMessage effect type lets you trigger social-style messages from abilities,
          quests, or events. This means game systems can talk to players through the same
          channel infrastructure that players use to talk to each other.
        </p>
      </Section>
    </div>
  );
}