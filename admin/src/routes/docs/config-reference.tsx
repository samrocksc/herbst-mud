import { createFileRoute } from '@tanstack/react-router';
import { PageHeader } from '../../components/PageHeader';

export const Route = createFileRoute('/docs/config-reference')({
  component: ConfigReferenceDoc,
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

function ConfigReferenceDoc() {
  return (
    <div className="management-page max-w-4xl">
      <PageHeader title="Config Reference" backTo="/docs" />

      <InfoBox>
        <strong>TL;DR:</strong> Config keys control global game behavior. Edit these in the
        <strong>Config</strong> admin page. Changes take effect immediately (no restart needed for most).
      </InfoBox>

      <Section title="Combat Settings">
        <Table
          headers={['Key', 'Default', 'Description']}
          rows={[
            ['tick_interval', '1.5s', 'Seconds between combat ticks. Lower = faster combat.'],
            ['base_hit_chance', '50%', 'Starting hit chance before stat modifiers.'],
            ['crit_base_chance', '5%', 'Base critical hit percentage.'],
            ['min_damage', '1', 'Minimum damage per hit (cannot be reduced to 0).'],
            ['max_hit_chance', '95%', 'Hit chance cap. Never 100% (always a miss chance).'],
          ]}
        />
      </Section>

      <Section title="XP Settings">
        <Table
          headers={['Key', 'Default', 'Description']}
          rows={[
            ['xp_level_base', '100', 'XP needed for level 2.'],
            ['xp_level_multiplier', '1.5', 'Each level needs 1.5× the previous.'],
            ['xp_gap_penalty_start', '5', 'Levels below enemy before XP penalty kicks in.'],
            ['xp_gap_bonus_start', '5', 'Levels above enemy before XP bonus kicks in.'],
            ['xp_max_level', '100', 'Level cap. No XP gained at max.'],
          ]}
        />
        <div className="bg-surface-muted rounded-lg p-4 font-mono text-sm mt-3 mb-3">
          xp_for_level(n) = base × multiplier^(n−2)
        </div>
        <p className="text-text-muted">Example: Level 3 = 100 × 1.5¹ = 150 XP. Level 10 = 100 × 1.5⁸ ≈ 2,562 XP.</p>
      </Section>

      <Section title="NPC Settings">
        <Table
          headers={['Key', 'Default', 'Description']}
          rows={[
            ['npc_respawn_default', '60s', 'Default cooldown before NPC respawns.'],
            ['npc_max_instances', '50', 'Maximum live NPCs in the world.'],
            ['npc_hp_per_level', '15', 'HP added per NPC level.'],
            ['npc_damage_per_level', '1.5', 'Base damage added per NPC level.'],
          ]}
        />
      </Section>

      <Section title="Player Settings">
        <Table
          headers={['Key', 'Default', 'Description']}
          rows={[
            ['player_start_hp', '100', 'HP for new characters.'],
            ['player_start_mana', '50', 'Mana for new characters.'],
            ['player_start_stamina', '75', 'Stamina for new characters.'],
            ['player_max_inventory', '30', 'Item slots per character.'],
            ['player_death_penalty_xp', '10%', 'XP lost on death (percentage of current level).'],
          ]}
        />
      </Section>

      <Section title="World Settings">
        <Table
          headers={['Key', 'Default', 'Description']}
          rows={[
            ['world_name', 'New Venice', 'Display name of the game world.'],
            ['world_time_scale', '4', 'Game minutes per real second. 4 = 1 real second = 4 game minutes.'],
            ['weather_enabled', 'true', 'Toggle weather system.'],
            ['pvp_enabled', 'false', 'Allow player-vs-player combat.'],
            ['chat_range', 'same_room', 'How far chat messages travel: same_room, adjacent, global.'],
          ]}
        />
      </Section>
    </div>
  );
}
