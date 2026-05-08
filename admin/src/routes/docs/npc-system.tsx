import { createFileRoute } from '@tanstack/react-router'
import { PageHeader } from '../../components/PageHeader'

export const Route = createFileRoute('/docs/npc-system')({
  component: NPCSystemDoc,
})

function Section({ title, children }: Readonly<{ title: string; children: React.ReactNode }>) {
  return (
    <section className="mb-8">
      <h2 className="text-lg font-semibold text-text mb-3 pb-2 border-b border-border">{title}</h2>
      {children}
    </section>
  )
}

function InfoBox({ children }: Readonly<{ children: React.ReactNode }>) {
  return (
    <div className="bg-primary/10 border border-primary/30 rounded-lg p-4 mb-4 text-sm">
      {children}
    </div>
  )
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
  )
}

function NPCSystemDoc() {
  return (
    <div className="management-page max-w-4xl">
      <PageHeader title="NPC System" backTo="/docs" />

      <InfoBox>
        <strong>TL;DR:</strong> NPCs are defined by <strong>templates</strong> (static blueprints) and
        spawned as <strong>instances</strong> (live entities). When an instance dies, it respawns after a
        cooldown in one of the template's respawn rooms.
      </InfoBox>

      <Section title="Lifecycle">
        <ol className="text-sm text-text-muted space-y-2 mb-3">
          <li>
            <strong>Template:</strong> Static definition (name, level, race, abilities, loot).
            Created in the admin panel. Never changes at runtime.
          </li>
          <li>
            <strong>Instance:</strong> Live NPC in the world. Has current HP, position, state.
            Spawned from template at server start or on respawn.
          </li>
          <li>
            <strong>Combat:</strong> Instance engages players. Uses abilities from template.
            AI behavior: passive, aggressive, or flee.
          </li>
          <li>
            <strong>Death:</strong> Instance removed. Loot dropped. XP awarded to damage contributors.
          </li>
          <li>
            <strong>Respawn:</strong> After <code>respawn_cooldown</code> seconds, a new instance spawns
            in a random room from <code>respawn_rooms</code>.
          </li>
        </ol>
      </Section>

      <Section title="Template Fields">
        <Table
          headers={['Field', 'Meaning']}
          rows={[
            ['Name', 'Display name. e.g. "Junkyard Scrapper"'],
            ['Level', 'Power rating. Determines HP pool, damage output, and XP value.'],
            ['XP Value', 'Experience points awarded when defeated. Split among damage contributors.'],
            ['Race', 'Affects base stats, resistances, and available abilities. References races table.'],
            ['Respawn Cooldown', 'Seconds after death before a new instance appears.'],
            [
              'Respawn Rooms',
              'Comma-separated room IDs. New instance randomly picks one after cooldown.',
            ],
            ['Description', 'Flavor text shown when player examines the NPC.'],
          ]}
        />
      </Section>

      <Section title="Level Scaling">
        <div className="bg-surface-muted rounded-lg p-4 font-mono text-sm mb-3">
          npc_hp = 50 + (level × 15) + (CON × 5)
          npc_damage = 3 + (level × 1.5) + (STR × 0.8)
        </div>
        <p className="text-text-muted mb-2">
          NPCs scale linearly with level. Higher level = more HP, more damage, better loot.
          Race provides the base CON/STR values.
        </p>
        <Table
          headers={['Level Range', 'Tier', 'Typical HP', 'Typical Damage']}
          rows={[
            ['1–5', 'Weakling', '65–125', '4–10'],
            ['6–15', 'Standard', '140–275', '12–25'],
            ['16–30', 'Veteran', '290–500', '27–48'],
            ['31–50', 'Elite', '515–800', '50–78'],
            ['51+', 'Boss', '815+', '80+'],
          ]}
        />
      </Section>

      <Section title="AI Behaviors">
        <Table
          headers={['Behavior', 'Trigger', 'Action']}
          rows={[
            [
              'Passive',
              'Player enters room',
              'Watches but does not attack unless attacked first.',
            ],
            [
              'Aggressive',
              'Player enters room',
              'Immediately attacks if player level ≤ NPC level + 5.',
            ],
            [
              'Flee',
              'HP drops below 25%',
              'Attempts to move to adjacent room. Success = DEX check.',
            ],
            [
              'Healer',
              'Ally HP drops below 50%',
              'Uses healing ability if available. Priority over attack.',
            ],
          ]}
        />
      </Section>

      <Section title="XP Award Formula">
        <div className="bg-surface-muted rounded-lg p-4 font-mono text-sm mb-3">
          xp_per_player = (xp_value × damage_contribution%) × level_gap_multiplier
        </div>
        <Table
          headers={['Level Gap', 'Multiplier', 'Rule']}
          rows={[
            ['Player ≥ NPC level', '1.0', 'Full XP.'],
            ['Player = NPC level − 1 to −5', '0.8–0.5', 'Slight penalty for easy kills.'],
            ['Player = NPC level − 6 to −10', '0.4–0.1', 'Large penalty. Grind elsewhere.'],
            ['Player < NPC level − 10', '0.0', 'No XP. Too easy.'],
            ['Player > NPC level + 5', '1.2–1.5', 'Bonus for challenging yourself.'],
          ]}
        />
        <p className="text-text-muted mt-3">
          Damage contribution % is calculated from total damage dealt to the NPC during the fight.
          If solo: 100%. If party: proportional split.
        </p>
      </Section>
    </div>
  )
}
