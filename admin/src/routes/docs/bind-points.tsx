import { createFileRoute } from "@tanstack/react-router";
import { PageHeader } from "../../components/PageHeader";

export const Route = createFileRoute("/docs/bind-points")({
  component: BindPointsDoc,
});

function Section({ title, children }: Readonly<{ title: string; children: React.ReactNode }>) {
  return (
    <section className="mb-8">
      <h2 className="text-lg font-semibold text-text mb-3 pb-2 border-b border-border">{title}</h2>
      {children}
    </section>
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
            <tr key={i} className="border-b border-border last:border-b-0">
              {row.map((cell, j) => (
                <td key={j} className="px-3 py-2 text-text-muted">{cell}</td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

function BindPointsDoc() {
  return (
    <div className="p-6 max-w-3xl mx-auto">
      <PageHeader title="Bind Points & Root Room" backTo="/docs" />

      <Section title="Root Room">
        <p className="text-text-muted mb-3">
          The <strong className="text-text">root room</strong> is the room where all new characters
          enter the world. You can only have one root room at a time. If you set a new room as root,
          the old one loses the flag automatically.
        </p>
        <p className="text-text-muted mb-3">
          On the map, the root room shows up with a <strong className="text-accent">home icon</strong> (🏠).
          You might also see an older "starting room" flag kicking around, but that's being phased
          out in favor of the root room system.
        </p>
        <Table
          headers={["Field", "Type", "What it does"]}
          rows={[
            ["isRootRoom", "boolean", "Only one room can be root. New characters spawn here on their first login."],
          ]}
        />
      </Section>

      <Section title="Bind Points">
        <p className="text-text-muted mb-3">
          A <strong className="text-text">bind point</strong> is the room where a character respawns
          after dying or after being away for a long time. Each character has a{" "}
          <code className="bg-surface-dark text-text-inverse px-1 rounded">respawnRoomId</code> field
          that stores their current bind point.
        </p>
        <Table
          headers={["Field", "Type", "What it does"]}
          rows={[
            ["respawnRoomId", "integer", "The room where this character respawns. New characters default to the root room."],
            ["lastSeenAt", "timestamp", "When the character was last online. Used to decide where they appear on reconnect."],
          ]}
        />
      </Section>

      <Section title="Setting a Bind Point">
        <p className="text-text-muted mb-3">
          You set bind points using the <code className="bg-surface-dark text-text-inverse px-1 rounded">set_bind_point</code>{" "}
          effect type. When an ability with this effect fires (with target set to self), the character's{" "}
          <code className="bg-surface-dark text-text-inverse px-1 rounded">respawnRoomId</code> updates
          to whatever room they're currently standing in.
        </p>
        <p className="text-text-muted mb-3">
          A common pattern is to have an NPC offer to "bind your spirit here" as a service. You
          create an ability with the set_bind_point effect, give it to an NPC, and when the NPC
          casts it on the player (or the player uses it on themselves), their respawn point is
          updated to that location.
        </p>
        <Table
          headers={["Property", "Value"]}
          rows={[
            ["Effect Type", "set_bind_point"],
            ["Target", "self"],
            ["Mana Cost", "0"],
            ["Cooldown", "0"],
            ["Result", "Character's respawnRoomId is set to their current room"],
          ]}
        />
      </Section>

      <Section title="Reconnect Positioning">
        <p className="text-text-muted mb-3">
          When a player logs back in after being offline, the game picks where they show up:
        </p>
        <Table
          headers={["Condition", "Where they appear"]}
          rows={[
            ["Offline less than 1 hour", "Right back where they were (currentRoomId)"],
            ["Offline 1 hour or more", "At their bind point (respawnRoomId)"],
            ["Brand-new character (first login)", "The root room (isRootRoom)"],
            ["After dying in combat", "Their bind point (respawnRoomId)"],
          ]}
        />
        <p className="text-text-muted mt-3">
          The idea is simple. If someone just disconnected by accident, they come back where they
          left off. If they've been gone a while, they respawn at their bind point. And if the
          character has no bind point set, they fall back to the root room.
        </p>
      </Section>

      <Section title="Death & Respawn">
        <p className="text-text-muted mb-3">
          When a character dies in combat, here's what happens:
        </p>
        <ol className="list-decimal list-inside text-text-muted space-y-1 ml-2">
          <li>HP is restored to maximum</li>
          <li>The character is warped to their bind point (respawnRoomId)</li>
          <li>A 15-second "ethereal form" period kicks in, preventing them from jumping right back into a fight</li>
          <li>If no bind point is set, the root room is used instead</li>
        </ol>
      </Section>

      <Section title="API Reference">
        <Table
          headers={["Endpoint", "Method", "What it does"]}
          rows={[
            ["PUT /characters/:id", "PATCH", "Update respawnRoomId and lastSeenAt on a character"],
            ["POST /rooms", "POST", "Create a new room with the isRootRoom flag"],
            ["PUT /rooms/:id", "PUT", "Update isRootRoom. Setting a new root room clears the old one automatically."],
          ]}
        />
      </Section>
    </div>
  );
}