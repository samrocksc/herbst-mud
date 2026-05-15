import { createFileRoute } from '@tanstack/react-router';
import { PageHeader } from '../../components/PageHeader';

export const Route = createFileRoute('/docs/bind-points')({
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
          The <strong className="text-text">root room</strong> is the single room designated as the world's
          spawn point. Only one room can be the root room at a time — setting a new room as root
          automatically clears the flag from the previous one.
        </p>
        <p className="text-text-muted mb-3">
          On the map, the root room is shown with a <strong className="text-accent">home icon</strong> (🏠).
          The legacy "starting room" flag still exists but is being replaced by the root room concept.
        </p>
        <Table
          headers={['Field', 'Type', 'Description']}
          rows={[
            ['isRootRoom', 'boolean', 'Only one room can be root. New characters spawn here.'],
          ]}
        />
      </Section>

      <Section title="Bind Points">
        <p className="text-text-muted mb-3">
          A <strong className="text-text">bind point</strong> is a room where a character will respawn
          after death or after being offline for an extended period. Each character has a{' '}
          <code className="bg-surface-dark text-text-inverse px-1 rounded">respawnRoomId</code> field
          that tracks their current bind point.
        </p>
        <Table
          headers={['Field', 'Type', 'Description']}
          rows={[
            ['respawnRoomId', 'integer', 'Room ID where the character respawns. Defaults to the root room on creation.'],
            ['lastSeenAt', 'timestamp', 'When the character was last online. Used for reconnect positioning.'],
          ]}
        />
      </Section>

      <Section title="Setting a Bind Point">
        <p className="text-text-muted mb-3">
          Bind points are set via the <code className="bg-surface-dark text-text-inverse px-1 rounded">set_bind_point</code>{' '}
          effect type. When an ability with this effect fires (target: self), the character's{' '}
          <code className="bg-surface-dark text-text-inverse px-1 rounded">respawnRoomId</code> is updated
          to their current room.
        </p>
        <p className="text-text-muted mb-3">
          In practice, NPCs can cast bind-point abilities on players as a reward or service. For example,
          an NPC in a distant town could offer to "bind your spirit here" — setting the player's respawn
          point to that location.
        </p>
        <Table
          headers={['Property', 'Value']}
          rows={[
            ['Effect Type', 'set_bind_point'],
            ['Target', 'self'],
            ['Mana Cost', '0'],
            ['Cooldown', '0'],
            ['Result', 'Character\'s respawnRoomId set to current room'],
          ]}
        />
      </Section>

      <Section title="Reconnect Positioning">
        <p className="text-text-muted mb-3">
          When a player reconnects after being offline, the game decides where they appear:
        </p>
        <Table
          headers={['Condition', 'Spawn Location']}
          rows={[
            ['Offline < 1 hour', 'Last known room (currentRoomId)'],
            ['Offline >= 1 hour', 'Bind point (respawnRoomId)'],
            ['New character (first login)', 'Root room (isRootRoom)'],
            ['Death', 'Bind point (respawnRoomId)'],
          ]}
        />
        <p className="text-text-muted mt-3">
          This means a quick disconnect and reconnect puts you back where you were, while a longer absence
          returns you to your bind point. If no bind point is set, the root room is used as the fallback.
        </p>
      </Section>

      <Section title="Death & Respawn">
        <p className="text-text-muted mb-3">
          When a character dies in combat:
        </p>
        <ol className="list-decimal list-inside text-text-muted space-y-1 ml-2">
          <li>HP is restored to maximum</li>
          <li>Character is moved to their bind point (respawnRoomId)</li>
          <li>A 15-second "ethereal form" period prevents immediate re-engagement</li>
          <li>If no bind point is set, the root room is used</li>
        </ol>
      </Section>

      <Section title="API Reference">
        <Table
          headers={['Endpoint', 'Method', 'Description']}
          rows={[
            ['PUT /characters/:id', 'PATCH', 'Update respawnRoomId and lastSeenAt on a character'],
            ['POST /rooms', 'POST', 'Create room with isRootRoom flag'],
            ['PUT /rooms/:id', 'PUT', 'Update isRootRoom (enforces single root room)'],
          ]}
        />
      </Section>
    </div>
  );
}