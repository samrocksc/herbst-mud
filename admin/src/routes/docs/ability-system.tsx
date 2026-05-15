import { createFileRoute } from '@tanstack/react-router';
import { PageHeader } from '../../components/PageHeader';

export const Route = createFileRoute('/docs/ability-system')({
  component: AbilitySystemDoc,
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

function AbilitySystemDoc() {
  return (
    <div className="management-page max-w-4xl">
      <PageHeader title="Ability System" backTo="/docs" />

      <InfoBox>
        <strong>TL;DR:</strong> Abilities are combat actions players press <kbd className="bg-surface-dark text-text-inverse px-1 rounded">1–5</kbd> to
        activate. They cost mana/stamina, have cooldowns, and scale with character stats.
        Each ability has one or more <strong>effects</strong> that define what happens (damage, heal, buff, etc.).
        The <strong>5 classless abilities</strong> are available to every character regardless of class.
      </InfoBox>

      <Section title="Four-Domain Model">
        <p className="text-text-muted mb-3">
          The system separates concerns across four domains:
        </p>
        <Table
          headers={['Domain', 'What it is', 'Examples', 'Entity']}
          rows={[
            ['Abilities', 'Actions characters perform', 'Concentrate, Haymaker, Fireball', 'Ability'],
            ['Skills', 'Leveled proficiencies (training)', 'Blades, Staves, Light Armor', 'Character columns'],
            ['Stats', 'Numeric attributes', 'Strength, Dexterity, Wisdom', 'Character fields'],
            ['Effects', 'What actually happens', 'Damage, Heal, Buff, Stun', 'AbilityEffect'],
          ]}
        />
      </Section>

      <Section title="The 5 Classless Abilities">
        <p className="text-text-muted mb-3">
          Every character gets these 5 combat abilities in slots 1–5. Each ability has one or more
          effects defined in the AbilityEffect entity.
        </p>
        <Table
          headers={['Slot', 'Ability', 'Effects', 'Mana', 'Stamina', 'Cooldown']}
          rows={[
            ['1', 'Concentrate', 'accuracy_boost (self, 4t, WIS scaled)', '10', '0', '8s'],
            ['2', 'Haymaker', 'damage (enemy, STR scaled) + debuff (self, 1t)', '0', '15', '6s'],
            ['3', 'Back-off', 'dodge_all (self, 1t)', '0', '25', '10s'],
            ['4', 'Scream', 'buff (self, CON scaled, 2t) + debuff (enemy, CON scaled, 2t)', '5', '10', '12s'],
            ['5', 'Slap', 'stun (enemy, 1t, DEX contest)', '0', '12', '8s'],
          ]}
        />
      </Section>

      <Section title="Effect Types">
        <p className="text-text-muted mb-3">
          Effects are generic building blocks — each ability can have multiple effects with different
          targets, values, and scaling. The effect types are:
        </p>
        <Table
          headers={['Type', 'What It Does', 'Target Options']}
          rows={[
            ['damage', 'Subtracts HP from target. Can specify a damage subtype.', 'enemy, area, random_enemy'],
            ['heal', 'Restores HP to target (self or ally).', 'self, ally'],
            ['buff', 'Applies a positive status effect for a duration.', 'self, ally'],
            ['debuff', 'Applies a negative status effect for a duration.', 'enemy'],
            ['dot', 'Damage Over Time — repeats every combat tick for the duration.', 'enemy'],
            ['hot', 'Heal Over Time — repeats every combat tick for the duration.', 'self, ally'],
            ['stun', 'Target skips their next turn. Resisted by CON contest.', 'enemy'],
            ['accuracy_boost', 'Increases hit chance for a duration.', 'self'],
            ['dodge_all', 'Avoids all attacks for the duration.', 'self'],
            ['set_bind_point', 'Sets the character\'s respawn point to the current room.', 'self'],
          ]}
        />
      </Section>

      <Section title="Damage Subtypes">
        <p className="text-text-muted mb-3">
          Damage effects can specify a subtype for future resistance/weakness calculations:
        </p>
        <Table
          headers={['Subtype', 'Description']}
          rows={[
            ['slashing', 'Swords, axes — countered by heavy armor'],
            ['piercing', 'Daggers, arrows — countered by cloth armor'],
            ['bludgeoning', 'Clubs, fists — countered by light armor'],
            ['fire', 'Fireballs, dragon breath'],
            ['cold', 'Ice spells, frost weapons'],
            ['lightning', 'Shock spells, storm abilities'],
            ['poison', 'Toxins, venoms'],
            ['psychic', 'Mind-based attacks'],
          ]}
        />
      </Section>

      <Section title="Ability Fields Reference">
        <Table
          headers={['Field', 'Meaning']}
          rows={[
            ['Name', 'The command name players type, e.g. "concentrate".'],
            ['Description', 'Flavor text shown when examining the ability.'],
            ['Ability Type', 'combat, magic, utility, healing, support, defensive'],
            ['Ability Class', 'active (press to use), passive (auto-trigger), toggle (on/off)'],
            ['Required Tag', 'Item tag required to use this ability.'],
            ['Level Req', 'Minimum character level to learn this ability.'],
            ['Cooldown (s)', 'Seconds before ability can be reused.'],
            ['Mana Cost', 'MP drained on activation.'],
            ['Stamina Cost', 'SP drained on activation.'],
            ['HP Cost', 'Self-damage to cast.'],
            ['Proc Chance', '0.0–1.0 trigger probability (passives only).'],
            ['Proc Event', 'When to roll: on_hit, on_hit_received, on_crit, on_kill'],
          ]}
        />
      </Section>

      <Section title="Scaling Formula">
        <div className="bg-surface-muted rounded-lg p-4 font-mono text-sm mb-3">
          final_value = base_value + (stat_modifier × scaling_ratio × base_value)
        </div>
        <p className="text-text-muted mb-3">
          Where <code className="bg-surface-dark text-text-inverse px-1 rounded">stat_modifier = (stat - 10) / 2</code> (standard D&D modifier).
        </p>
        <p className="text-text-muted mb-3">Example: Haymaker (base 15, STR scaling, ratio 0.5):</p>
        <ul className="text-sm text-text-muted space-y-1 mb-3">
          <li>Character with STR 18 (+4 modifier)</li>
          <li>Final damage: <code className="bg-surface-dark text-text-inverse px-1 rounded">15 + (4 × 0.5 × 15) = 15 + 30 = 45</code></li>
        </ul>
      </Section>

      <Section title="Ability Combos">
        <Table
          headers={['Combo', 'How It Works']}
          rows={[
            ['Scream → Haymaker', 'Scream buffs STR, Haymaker uses it. Trade INT/WIS for massive damage.'],
            ['Concentrate → Back-off', 'Use Concentrate for accuracy. Back-off when things get dicey.'],
            ['Slap → Haymaker', 'Stun prevents enemy action. Free Haymaker with no retaliation risk.'],
            ['Back-off → Potion → Attack', 'Dodge round gives time to heal. Re-engage with full resources.'],
          ]}
        />
      </Section>

      <Section title="Passive Abilities (Formerly Talents)">
        <p className="text-text-muted mb-3">
          Passive abilities use <code className="bg-surface-dark text-text-inverse px-1 rounded">ability_class=passive</code> and
          trigger automatically based on <strong>proc_event</strong> (on_hit, on_crit, etc.)
          with a <strong>proc_chance</strong> (0.0–1.0). When triggered, their effects are applied automatically.
        </p>
      </Section>
    </div>
  );
}