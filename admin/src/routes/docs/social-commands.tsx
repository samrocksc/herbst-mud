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
        <strong>TL;DR:</strong> Social commands are player-to-player and player-to-world
        communication tools. Predefined commands (say, shout, whisper, tell, emote) plus
        configurable custom socials via the admin panel.
      </InfoBox>

      <Section title="Predefined Socials">
        <Table
          headers={["Command", "Range", "Description"]}
          rows={[
            ["say &lt;text&gt;", "same_room", "Speak to everyone in the same room."],
            ["shout &lt;text&gt;", "adjacent", "Shout to adjacent rooms."],
            ["whisper &lt;player&gt; &lt;text&gt;", "same_room", "Private message to a player in the same room."],
            ["tell &lt;player&gt; &lt;text&gt;", "global", "Private message to a player anywhere in the world."],
            ["emote &lt;text&gt;", "same_room", "Perform an action visible to the room (e.g., /me smiles)."],
          ]}
        />
      </Section>

      <Section title="Custom Socials">
        <p className="text-text-muted mb-3">
          Admins can create custom social commands at /socials. Each social has:
        </p>
        <Table
          headers={["Field", "Description"]}
          rows={[
            ["name", "Command name (e.g., laugh, wave)."],
            ["message_format", "Template string. $sender is the player's name."],
            ["channel", "same_room, adjacent, or global."],
            ["cooldown_secs", "Optional cooldown between uses."],
          ]}
        />
      </Section>

      <Section title="Channel Configuration">
        <p className="text-text-muted mb-3">
          Chat channels are configured at /channels. Channels route messages to specific
          audiences and can be restricted by class, level, or faction.
        </p>
        <Table
          headers={["Channel Field", "Description"]}
          rows={[
            ["name", "Channel identifier (e.g., ooc, auction, guild)."],
            ["range", "same_room, adjacent, or global."],
            ["min_level", "Minimum character level to use."],
            ["required_faction", "Optional faction requirement."],
          ]}
        />
      </Section>

      <Section title="SendMessage Effect">
        <p className="text-text-muted mb-3">
          The SendMessage effect type can trigger social-style messages from abilities,
          quests, or events. This lets game systems communicate with players through
          the same channel infrastructure.
        </p>
      </Section>
    </div>
  );
}
