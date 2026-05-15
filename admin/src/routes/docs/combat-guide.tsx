import { createFileRoute } from '@tanstack/react-router';
import { PageHeader } from '../../components/PageHeader';

export const Route = createFileRoute('/docs/combat-guide')({
  component: CombatGuideDoc,
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

function CombatGuideDoc() {
  return (
    <div className="management-page max-w-4xl">
      <PageHeader title="Combat Guide" backTo="/docs" />

      <InfoBox>
        <strong>TL;DR:</strong> Combat is tick-based (1.5s ticks). Characters act by paying tick costs.
        Stats (STR/DEX/CON/INT/WIS/CHA) drive damage, accuracy, dodge, and magic. Press
        <kbd className="bg-surface-dark text-text-inverse px-1 rounded">1–5</kbd> in combat to use classless skills.
      </InfoBox>

      <Section title="Tick System">
        <p className="text-text-muted mb-3">
          The combat clock fires globally every 1.5 seconds. Every action has a tick cost.
          When a character's accumulated tick debt is paid, they can act again.
        </p>
        <Table
          headers={['Action', 'Tick Cost', 'Notes']}
          rows={[
            ['Attack', '1', 'Basic melee or ranged attack. Damage = weapon + stat.'],
            ['Defend', '0', 'Buff action — raises armor briefly without costing a turn.'],
            ['Flee', '1', 'Attempt to escape combat. Success chance = DEX vs enemy level.'],
            ['Use Item', '1', 'Drink potion, throw grenade, etc.'],
            ['Activate Skill', '1–3', 'Classless skills (1 tick) or passive abilities (1–3 ticks).'],
          ]}
        />
      </Section>

      <Section title="Character Stats">
        <Table
          headers={['Stat', 'Abbreviation', 'What It Does']}
          rows={[
            ['Strength', 'STR', 'Melee damage bonus, carrying capacity.'],
            ['Dexterity', 'DEX', 'Accuracy (hit chance), dodge chance, attack speed.'],
            ['Constitution', 'CON', 'HP pool, stamina pool, resistance to stun/poison.'],
            ['Intelligence', 'INT', 'Magic power, tech skill effectiveness, mana pool.'],
            ['Wisdom', 'WIS', 'Perception, wind magic, healing power, resistance.'],
            ['Charisma', 'CHA', 'Shop prices, NPC disposition, party size limit.'],
          ]}
        />
      </Section>

      <Section title="Damage Formula">
        <div className="bg-surface-muted rounded-lg p-4 font-mono text-sm mb-3">
          damage = weapon_base + (STR × 0.5) + random(1, weapon_dice_sides) − target_armor
        </div>
        <p className="text-text-muted mb-2">
          If the attacker has a <strong>scaling_stat</strong> set on their ability, the formula becomes:
        </p>
        <div className="bg-surface-muted rounded-lg p-4 font-mono text-sm mb-3">
          damage = base + (stat_value × scaling_percent × base) − target_armor
        </div>
        <ul className="text-sm text-text-muted space-y-1">
          <li><strong>Weapon base:</strong> Item's <code>damage</code> field.</li>
          <li><strong>Random roll:</strong> 1 to weapon's dice sides (e.g. d8 = 1–8).</li>
          <li><strong>Armor:</strong> Sum of all equipped items' <code>armor</code> values.</li>
          <li><strong>Minimum damage:</strong> 1 (you always deal at least 1 HP).</li>
        </ul>
      </Section>

      <Section title="Hit Chance & Dodge">
        <div className="bg-surface-muted rounded-lg p-4 font-mono text-sm mb-3">
          hit_chance = 50% + (attacker_DEX − target_DEX) × 2.5%
        </div>
        <p className="text-text-muted mb-2">Base 50%. Each point of DEX difference shifts it by 2.5%.</p>
        <ul className="text-sm text-text-muted space-y-1 mb-3">
          <li>Minimum hit chance: 5% (always a glimmer of hope).</li>
          <li>Maximum hit chance: 95% (never guaranteed).</li>
          <li><strong>Back-off</strong> skill: sets dodge to 100% for 1 round vs ALL attacks.</li>
          <li><strong>Concentrate</strong> skill: +WIS to hit rolls for 4 rounds.</li>
        </ul>
      </Section>

      <Section title="Critical Hits">
        <p className="text-text-muted mb-3">
          Critical hits deal <strong>150% damage</strong> (before armor subtraction) and can trigger
          <code>on_crit</code> proc events on abilities.
        </p>
        <div className="bg-surface-muted rounded-lg p-4 font-mono text-sm mb-3">
          crit_chance = 5% + (DEX − 10) × 0.5%
        </div>
        <p className="text-text-muted">
          Base 5% crit. Each DEX above 10 adds 0.5%. Some items and buffs can raise this further.
        </p>
      </Section>

      <Section title="Combat Flow Example">
        <ol className="text-sm text-text-muted space-y-2 mb-3">
          <li>
            <strong>Round 1:</strong> Player presses <kbd>1</kbd> (Concentrate). Cost: 10 MP, 0 ticks.
            Buff active: +WIS to hit for 4 rounds.
          </li>
          <li>
            <strong>Round 2:</strong> Player attacks. Hit chance boosted by Concentrate.
            DEX 14 vs enemy DEX 10 → 60% base + WIS bonus ≈ 70% hit.
          </li>
          <li>
            <strong>Round 3:</strong> Enemy attacks. Player uses <kbd>3</kbd> (Back-off) as reaction.
            Cost: 25 SP. Dodge: 100% this round.
          </li>
          <li>
            <strong>Round 4:</strong> Player presses <kbd>2</kbd> (Haymaker). Cost: 15 SP, 1 tick.
            STR 16 → damage = 50 + (16 × 0.05 × 50) = 90. But −DEX to hit this attack.
          </li>
          <li>
            <strong>Round 5:</strong> Player presses <kbd>4</kbd> (Scream). Cost: 5 MP + 10 SP.
            +DEX/STR, −WIS/INT for 2 rounds. Next Haymaker will be devastating.
          </li>
        </ol>
      </Section>

      <Section title="Status Effects">
        <Table
          headers={['Effect', 'What It Does', 'Duration']}
          rows={[
            ['Stunned', 'Cannot act. Skips next turn.', '1 round (Slap)'],
            ['Poisoned (dot)', 'Loses HP every tick.', 'Until cured or duration ends'],
            ['Buffed', 'Raised stat(s).', 'Ability duration (ticks)'],
            ['Debuffed', 'Lowered stat(s).', 'Ability duration (ticks)'],
            ['Concentrating', '+WIS to hit rolls.', '4 rounds'],
            ['Haymaker stance', '+STR damage, −DEX hit.', '1 attack'],
            ['Screaming', '+DEX/STR, −WIS/INT.', '2 rounds'],
          ]}
        />
      </Section>
    </div>
  );
}
