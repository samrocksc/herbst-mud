import { createFileRoute } from "@tanstack/react-router";
import { PageHeader } from "../../components/PageHeader";

export const Route = createFileRoute("/docs/trainable-skills")({
  component: TrainableSkillsDoc,
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

function TrainableSkillsDoc() {
  return (
    <div className="management-page max-w-4xl">
      <PageHeader title="Trainable Skills" backTo="/docs" />

      <InfoBox>
        <strong>TL;DR:</strong> Trainable skills are <strong>proficiency multipliers</strong>. Think
        of them as things like "Blades", "Fire Magic", or even "Pizza Making". Characters get better
        at them over time by using related abilities, and that improvement makes those abilities more
        effective. They are <strong>not</strong> combat actions (those are called Abilities).
      </InfoBox>

      <Section title="Skill Categories">
        <p className="text-text-muted mb-3">
          Skills fall into a few broad categories. Each one affects different aspects of gameplay:
        </p>
        <Table
          headers={["Category", "Examples", "What it affects"]}
          rows={[
            ["Combat", "Blades, Knives, Staves, Brawling, Martial Arts, Bows, Thrown", "Damage with related weapons"],
            ["Magic", "Fire, Water, Wind, Earth, Light, Dark", "Spell power and mana efficiency"],
            ["Defense", "Light Armor, Cloth Armor, Heavy Armor, Shields", "Damage reduction and block chance"],
            ["Utility", "Tech, Pizza Making, Crafting, Trading", "Non-combat success rates and item quality"],
          ]}
        />
      </Section>

      <Section title="How Skills Relate to Abilities">
        <p className="text-text-muted mb-3">
          Every <strong>Ability</strong> (a combat action) is tagged with a skill category. When a
          character uses that ability, their skill level in the matching category boosts its effect.
          The formula is simple:
        </p>
        <div className="bg-surface-muted rounded-lg p-4 font-mono text-sm mb-3">
          ability_bonus = skill_level * 2%
        </div>
        <p className="text-text-muted mb-2">Here is a concrete example:</p>
        <ul className="text-sm text-text-muted space-y-1 mb-3">
          <li>Character has <strong>Blades 25</strong></li>
          <li>They use the "Slash" ability (tagged <code>sword, blade</code>)</li>
          <li>Base damage: 50</li>
          <li>Bonus: 50 * (25 * 0.02) = <strong>+25 damage</strong></li>
          <li>Final damage: 75</li>
        </ul>
        <p className="text-text-muted">
          At skill level 50 you get a +100% bonus (double damage). At 100 you get +200% (triple
          damage). This is why specializing in one skill pays off more than spreading points around.
        </p>
      </Section>

      <Section title="Training Mechanics">
        <p className="text-text-muted mb-3">
          Characters improve skills by using them. Every time they successfully perform an action,
          they earn XP in the related skill. The harder the action, the more XP they get.
        </p>
        <Table
          headers={["Action", "XP Gained", "Skill Trained"]}
          rows={[
            ["Hit with a sword", "2 to 5", "Blades"],
            ["Cast a fire spell", "3 to 6", "Fire Magic"],
            ["Block with a shield", "1 to 3", "Shields"],
            ["Craft an item", "5 to 10", "Crafting"],
            ["Complete a trade", "1 to 2", "Trading"],
            ["Parry an attack", "3 to 5", "Blades or Martial Arts"],
          ]}
        />
        <div className="bg-surface-muted rounded-lg p-4 font-mono text-sm mt-3 mb-3">
          level_up_threshold = current_level * 10
        </div>
        <p className="text-text-muted">
          Going from level 5 to 6 costs 50 XP. Going from 50 to 51 costs 500 XP. Higher levels take
          exponentially more effort, so reaching the top is a real commitment.
        </p>
      </Section>

      <Section title="Requirements Format">
        <p className="text-text-muted mb-3">
          The <strong>Requirements</strong> field on the Skills page lets you set prerequisites.
          Use a key:value format. All conditions must be met for the skill to unlock.
        </p>
        <Table
          headers={["Key", "What it checks", "Example"]}
          rows={[
            ["level", "Minimum character level", "level:10"],
            ["str", "Minimum Strength", "str:15"],
            ["dex", "Minimum Dexterity", "dex:12"],
            ["con", "Minimum Constitution", "con:8"],
            ["int", "Minimum Intelligence", "int:10"],
            ["wis", "Minimum Wisdom", "wis:8"],
            ["cha", "Minimum Charisma", "cha:5"],
          ]}
        />
        <p className="text-text-muted mt-3">
          Example: <code>level:5,str:10</code> means the character must be at least level 5 AND have
          at least 10 Strength. Separate conditions with commas, no spaces.
        </p>
      </Section>
    </div>
  );
}