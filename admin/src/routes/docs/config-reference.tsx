import { createFileRoute } from "@tanstack/react-router";
import { PageHeader } from "../../components/PageHeader";

export const Route = createFileRoute("/docs/config-reference")({
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
        <strong>TL;DR:</strong> These config keys control how your game behaves globally. You edit them on
        the <strong>Config</strong> page in the admin panel. Most changes take effect right away without
        restarting anything.
      </InfoBox>

      <Section title="Combat Settings">
        <p className="text-text-muted mb-3">
          These settings control the pace and feel of combat. Want faster, more lethal fights? Lower the
          tick interval. Want a gritty system where hits are rare but devastating? Raise min_damage and
          lower base_hit_chance.
        </p>
        <Table
          headers={["Key", "Default", "What it does"]}
          rows={[
            ["tick_interval", "1.5s", "How many seconds between combat ticks. Lower means faster combat."],
            ["base_hit_chance", "50%", "The starting hit chance before any stats are considered."],
            ["crit_base_chance", "5%", "The baseline critical hit percentage for everyone."],
            ["min_damage", "1", "The floor for damage. Even the weakest hit deals at least this much."],
            ["max_hit_chance", "95%", "The ceiling for hit chance. It never quite reaches 100%."],
          ]}
        />
      </Section>

      <Section title="XP Settings">
        <p className="text-text-muted mb-3">
          These control how fast players level up and how the game rewards or penalizes them based on the
          level gap between player and enemy.
        </p>
        <Table
          headers={["Key", "Default", "What it does"]}
          rows={[
            ["xp_level_base", "100", "The XP needed to reach level 2. Everything else scales from here."],
            ["xp_level_multiplier", "1.5", "Each level requires 1.5 times the previous level's XP."],
            ["xp_gap_penalty_start", "5", "When you are this many levels below the enemy, XP starts dropping."],
            ["xp_gap_bonus_start", "5", "When you are this many levels above the enemy, XP starts increasing."],
            ["xp_max_level", "100", "The level cap. Players stop earning XP once they hit this."],
          ]}
        />
        <div className="bg-surface-muted rounded-lg p-4 font-mono text-sm mt-3 mb-3">
          xp_for_level(n) = base × multiplier^(n-2)
        </div>
        <p className="text-text-muted">
          Example: Level 3 needs 100 × 1.5^1 = 150 XP. Level 10 needs 100 × 1.5^8 = roughly 2,562 XP.
          The curve gets steep fast, which keeps progression meaningful at higher levels.
        </p>
      </Section>

      <Section title="NPC Settings">
        <p className="text-text-muted mb-3">
          These control how NPCs populate the world and how tough they are at each level.
        </p>
        <Table
          headers={["Key", "Default", "What it does"]}
          rows={[
            ["npc_respawn_default", "60s", "How long after an NPC dies before a new one appears."],
            ["npc_max_instances", "50", "The cap on how many NPCs can be alive in the world at once."],
            ["npc_hp_per_level", "15", "How much HP each NPC level adds on top of the base."],
            ["npc_damage_per_level", "1.5", "How much base damage each NPC level adds."],
          ]}
        />
      </Section>

      <Section title="Player Settings">
        <p className="text-text-muted mb-3">
          These set the starting conditions for new characters and the penalty for dying.
        </p>
        <Table
          headers={["Key", "Default", "What it does"]}
          rows={[
            ["player_start_hp", "100", "How many HP a brand-new character begins with."],
            ["player_start_mana", "50", "How much mana a brand-new character begins with."],
            ["player_start_stamina", "75", "How much stamina a brand-new character begins with."],
            ["player_max_inventory", "30", "How many item slots a character has."],
            ["player_death_penalty_xp", "10%", "What fraction of current-level XP a player loses when they die."],
          ]}
        />
      </Section>

      <Section title="World Settings">
        <p className="text-text-muted mb-3">
          These shape the overall feel of your game world, from its name to whether players can fight each
          other.
        </p>
        <Table
          headers={["Key", "Default", "What it does"]}
          rows={[
            ["world_name", "New Venice", "The display name of your game world. Shown to players on login."],
            ["world_time_scale", "4", "Game minutes per real second. At 4, one real second equals four game minutes."],
            ["weather_enabled", "true", "Turn the weather system on or off."],
            ["pvp_enabled", "false", "Allow player-vs-player combat. Off by default."],
            ["chat_range", "same_room", "How far chat messages travel. Options: same_room, adjacent, global."],
          ]}
        />
      </Section>
    </div>
  );
}