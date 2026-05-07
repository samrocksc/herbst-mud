import { createFileRoute } from '@tanstack/react-router'
import { PageHeader } from '../../components/PageHeader'

export const Route = createFileRoute('/docs/ability-system')({
  component: AbilitySystemDoc,
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

function AbilitySystemDoc() {
  return (
    <div className="management-page max-w-4xl">
      <PageHeader title="Ability System" backTo="/docs" />

      <InfoBox>
        <strong>TL;DR:</strong> Abilities are combat actions players press <kbd className="bg-surface-dark text-text-inverse px-1 rounded">1–5</kbd> to
        activate. They cost mana/stamina, have cooldowns, and scale with character stats.
        The <strong>5 classless skills</strong> are available to every character regardless of class.
      </InfoBox>

      <Section title="The 5 Classless Skills">
        <p className="text-text-muted mb-3">
          Every character gets these 5 combat abilities in slots 1–5. They provide tactical
          options regardless of class. Skills cost mana/stamina to activate and have cooldowns.
        </p>
        <Table
          headers={['Slot', 'Skill', 'Effect', 'Mana', 'Stamina', 'Cooldown']}
          rows={[
            ['1', 'Concentrate', '+WIS to hit rolls for 4 rounds', '10', '0', '8 rounds'],
            ['2', 'Haymaker', '+STR damage, −DEX to hit for 1 attack', '0', '15', '6 rounds'],
            ['3', 'Back-off', 'Guaranteed dodge vs all attacks this round', '0', '25', '10 rounds'],
            ['4', 'Scream', '+DEX/STR, −WIS/INT for 2 rounds', '5', '10', '12 rounds'],
            ['5', 'Slap', 'DEX vs CON: stun for 1 round', '0', '12', '8 rounds'],
          ]}
        />
      </Section>

      <Section title="Ability Fields Reference">
        <p className="text-text-muted mb-3">
          These are the fields you see on the Abilities management page and what they mean in combat.
        </p>
        <Table
          headers={['Field', 'Meaning']}
          rows={[
            [
              'Name',
              'The command name players type, e.g. "concentrate". Also used as the API slug if slug is empty.',
            ],
            ['Description', 'Flavor text shown to players when they examine or help this ability.'],
            [
              'Skill Type',
              'combat=direct attack, magic=spell, utility=non-combat, healing=restore HP, support=buff allies',
            ],
            [
              'Required Tag',
              'Comma-separated item tags. Character must have an item with this tag equipped to use this ability.',
            ],
            ['Level Req', 'Minimum character level to learn this ability.'],
            ['Cost', 'Legacy flat energy cost. Prefer mana_cost/stamina_cost for new abilities.'],
            ['Cooldown (s)', 'Seconds before ability can be reused. Back-off=10s, Haymaker=6s.'],
            [
              'Effect Type',
              'damage/heal/buff/debuff/dot/hot, or special: concentrate/haymaker/scream/slap/backoff',
            ],
            ['Effect Value', 'Base magnitude. Scales with Scaling Stat if set.'],
            ['Effect Duration', 'Combat ticks the effect lasts. 1=one round, 4=Concentrate\'s full duration.'],
            ['Mana Cost', 'MP drained on activation. If mana < cost, ability fails.'],
            ['Stamina Cost', 'SP drained on activation. Fighters use this resource.'],
            ['HP Cost', 'Self-damage to cast. Berserker and blood magic abilities.'],
            ['Scaling Stat', 'STR=damage, DEX=dodge/accuracy, WIS=healing, INT=magic power.'],
            [
              'Scaling %/point',
              '0.05 = 5%. Formula: final = base + (stat × pct × base). STR=10, base=50, pct=0.05 → 75.',
            ],
            ['Proc Chance', '0.0–1.0 chance the effect triggers on the proc_event. 0.3 = 30%.'],
            [
              'Proc Event',
              'on_hit=attack lands, on_crit=critical hit, on_dodge=dodging, on_kill=enemy death',
            ],
            ['Skill Class', 'active=press button, passive=always on, toggle=on/off switch.'],
          ]}
        />
      </Section>

      <Section title="Effect Types">
        <Table
          headers={['Type', 'What It Does']}
          rows={[
            ['damage', 'Subtracts HP from target. Amount = effect_value + scaling.'],
            ['heal', 'Restores HP to target (self or ally). Amount = effect_value + scaling.'],
            ['buff', 'Raises a stat temporarily. e.g. +WIS from Concentrate.'],
            ['debuff', 'Lowers a stat temporarily. e.g. −INT from Scream.'],
            ['dot', 'Damage Over Time — repeats every combat tick for duration.'],
            ['hot', 'Heal Over Time — repeats every combat tick for duration.'],
            [
              'concentrate',
              'Special: +WIS to hit rolls for 4 rounds. Best for casters who need accuracy.',
            ],
            [
              'haymaker',
              'Special: +STR damage but −DEX to hit for 1 attack. High risk, high reward.',
            ],
            [
              'scream',
              'Special: Trade INT/WIS for DEX/STR for 2 rounds. Desperation burst damage.',
            ],
            [
              'slap',
              'Special: DEX vs CON check. If passed, target stunned for 1 round.',
            ],
            [
              'backoff',
              'Special: Guaranteed dodge against ALL attacks this round. Emergency escape.',
            ],
          ]}
        />
      </Section>

      <Section title="Scaling Formula">
        <div className="bg-surface-muted rounded-lg p-4 font-mono text-sm mb-3">
          final_effect = base_value + (stat_value × scaling_percent × base_value)
        </div>
        <p className="text-text-muted mb-3">Example:</p>
        <ul className="text-sm text-text-muted space-y-1 mb-3">
          <li>Base damage: <code className="bg-surface-dark text-text-inverse px-1 rounded">50</code></li>
          <li>Scaling Stat: <code className="bg-surface-dark text-text-inverse px-1 rounded">STR</code></li>
          <li>STR value: <code className="bg-surface-dark text-text-inverse px-1 rounded">10</code></li>
          <li>Scaling %/point: <code className="bg-surface-dark text-text-inverse px-1 rounded">0.05</code> (5%)</li>
          <li>Final damage: <code className="bg-surface-dark text-text-inverse px-1 rounded">50 + (10 × 0.05 × 50) = 75</code></li>
        </ul>
        <p className="text-text-muted">
          If <strong>Scaling Stat</strong> is empty, the ability uses the base value with no scaling.
        </p>
      </Section>

      <Section title="Skill Combos">
        <Table
          headers={['Combo', 'How It Works']}
          rows={[
            ['Scream → Haymaker', 'Scream buffs STR, Haymaker uses it. Trade INT/WIS for massive damage.'],
            [
              'Concentrate → Back-off',
              'Use Concentrate for accuracy. Back-off when things get dicey. Next attacks hit hard.',
            ],
            ['Slap → Haymaker', 'Stun prevents enemy action. Free Haymaker with no retaliation risk.'],
            ['Back-off → Potion → Attack', 'Dodge round gives time to heal. Re-engage with full resources.'],
          ]}
        />
      </Section>

      <Section title="Stat Synergies">
        <Table
          headers={['Skill', 'Primary Stat', 'Best For']}
          rows={[
            ['Concentrate', 'WIS', 'Clerics, Mages, Rangers'],
            ['Haymaker', 'STR', 'Warriors, Barbarians, Fighters'],
            ['Back-off', 'DEX', 'Rogues, Rangers, Monks'],
            ['Scream', 'CON (magnitude)', 'Tanky builds'],
            ['Slap', 'DEX', 'Rogues, fast attackers'],
          ]}
        />
      </Section>
    </div>
  )
}
