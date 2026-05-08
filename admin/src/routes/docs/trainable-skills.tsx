import { createFileRoute } from '@tanstack/react-router'
import { PageHeader } from '../../components/PageHeader'

export const Route = createFileRoute('/docs/trainable-skills')({
  component: TrainableSkillsDoc,
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

function TrainableSkillsDoc() {
  return (
    <div className="management-page max-w-4xl">
      <PageHeader title="Trainable Skills" backTo="/docs" />

      <InfoBox>
        <strong>TL;DR:</strong> Trainable skills are <strong>proficiency multipliers</strong> — "Blades",
        "Fire Magic", "Pizza Making". Characters train them over time to improve effectiveness
        with related abilities. They are <strong>not</strong> combat actions (those are Abilities).
      </InfoBox>

      <Section title="Skill Categories">
        <Table
          headers={['Category', 'Examples', 'Affects']}
          rows={[
            ['Combat', 'Blades, Knives, Staves, Brawling, Martial Arts, Bows, Thrown', 'Damage with related weapons'],
            ['Magic', 'Fire, Water, Wind, Earth, Light, Dark', 'Spell power and mana efficiency'],
            ['Defense', 'Light Armor, Cloth Armor, Heavy Armor, Shields', 'Damage reduction, block chance'],
            ['Utility', 'Tech, Pizza Making, Crafting, Trading', 'Non-combat success rates, item quality'],
          ]}
        />
      </Section>

      <Section title="How Skills Relate to Abilities">
        <p className="text-text-muted mb-3">
          Every <strong>Ability</strong> (combat action) has a <code>required_tag</code> or implicit category.
          If a character's <strong>Trainable Skill</strong> in that category is high, the ability gets bonuses:
        </p>
        <div className="bg-surface-muted rounded-lg p-4 font-mono text-sm mb-3">
          ability_bonus = skill_level × 2%
        </div>
        <p className="text-text-muted mb-2">Example:</p>
        <ul className="text-sm text-text-muted space-y-1 mb-3">
          <li>Character has <strong>Blades 25</strong></li>
          <li>Uses "Slash" ability (tagged <code>sword, blade</code>)</li>
          <li>Base damage: 50</li>
          <li>Bonus: 50 × (25 × 0.02) = <strong>+25 damage</strong></li>
          <li>Final: 75 damage</li>
        </ul>
        <p className="text-text-muted">
          At skill level 50: +100% bonus (double damage). At 100: +200% (triple damage).
          This is why specialization beats generalization.
        </p>
      </Section>

      <Section title="Training Mechanics">
        <p className="text-text-muted mb-3">Skills improve through use. Each successful action adds XP to the relevant skill.</p>
        <Table
          headers={['Action', 'XP Gain', 'Skill Trained']}
          rows={[
            ['Hit with sword', '2–5', 'Blades'],
            ['Cast fire spell', '3–6', 'Fire Magic'],
            ['Block with shield', '1–3', 'Shields'],
            ['Craft item', '5–10', 'Crafting'],
            ['Trade successfully', '1–2', 'Trading'],
            ['Parry attack', '3–5', 'Blades or Martial Arts'],
          ]}
        />
        <div className="bg-surface-muted rounded-lg p-4 font-mono text-sm mt-3 mb-3">
          level_up_threshold = current_level × 10
        </div>
        <p className="text-text-muted">
          To go from level 5 → 6, you need 50 XP. From 50 → 51, you need 500 XP.
          Higher levels require exponentially more effort.
        </p>
      </Section>

      <Section title="Requirements Format">
        <p className="text-text-muted mb-3">
          The <strong>Requirements</strong> field on the Skills page uses a key:value format.
          All conditions must be met.
        </p>
        <Table
          headers={['Key', 'Meaning', 'Example']}
          rows={[
            ['level', 'Minimum character level', 'level:10'],
            ['str', 'Minimum Strength', 'str:15'],
            ['dex', 'Minimum Dexterity', 'dex:12'],
            ['con', 'Minimum Constitution', 'con:8'],
            ['int', 'Minimum Intelligence', 'int:10'],
            ['wis', 'Minimum Wisdom', 'wis:8'],
            ['cha', 'Minimum Charisma', 'cha:5'],
          ]}
        />
        <p className="text-text-muted mt-3">Example full requirement: <code>level:5,str:10</code> means level 5 AND 10 STR.</p>
      </Section>
    </div>
  )
}
