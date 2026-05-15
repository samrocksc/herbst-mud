import { createFileRoute } from '@tanstack/react-router';
import { PageHeader } from '../../components/PageHeader';

export const Route = createFileRoute('/docs/examine-skill')({
  component: ExamineSkillDoc,
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

function ExamineSkillDoc() {
  return (
    <div className="management-page max-w-4xl">
      <PageHeader title="Examine Skill" backTo="/docs" />

      <InfoBox>
        <strong>TL;DR:</strong> The <strong>examine</strong> skill lets players inspect rooms, items, and
        NPCs to discover hidden details. Higher skill = more details revealed. Based on INT and WIS.
        Max level: 100.
      </InfoBox>

      <Section title="How Examine Works">
        <p className="text-text-muted mb-3">
          Players type <code>examine</code> (or <code>examine [target]</code>) to inspect their surroundings.
          The game rolls against hidden details to see what gets revealed.
        </p>
        <div className="bg-surface-muted rounded-lg p-4 font-mono text-sm mb-3">
          examine_roll = skill_level + INT + (WIS / 2) + random(1, 10)
        </div>
        <p className="text-text-muted mb-2">
          Two reveal modes:
        </p>
        <ul className="text-sm text-text-muted space-y-1 mb-3">
          <li><strong>Automatic:</strong> Revealed if skill_level ≥ min_examine_level. No roll needed.</li>
          <li><strong>Check:</strong> Revealed if examine_roll ≥ DC (difficulty class). Uses INT or WIS for the check.</li>
        </ul>
      </Section>

      <Section title="Skill Levels & Bonuses">
        <Table
          headers={['Examine Level', 'Bonus', 'What It Means']}
          rows={[
            ['0–25', '0%', 'Basic description only.'],
            ['26–50', '10%', 'Reveals hidden items in rooms.'],
            ['51–75', '25%', 'Reveals NPC weaknesses and item stats.'],
            ['76–90', '50%', 'Reveals secret exits and quest hints.'],
            ['91–100', '75%', 'Reveals everything including lore and backstory.'],
          ]}
        />
        <p className="text-text-muted mt-3">
          The bonus % applies to shop prices (better appraisal) and loot quality (spotting valuables).
        </p>
      </Section>

      <Section title="XP Rewards">
        <Table
          headers={['Action', 'XP Gain']}
          rows={[
            ['First-time examine of a room', '+1 XP'],
            ['Discover a hidden detail', '+2 XP'],
            ['Reveal a hidden exit', '+5 XP'],
            ['Decrypt a secret message', '+10 XP'],
          ]}
        />
        <div className="bg-surface-muted rounded-lg p-4 font-mono text-sm mt-3 mb-3">
          level_up_every = 10 XP
          max_level = 100
        </div>
        <p className="text-text-muted">
          At level 100, examine reveals all hidden details automatically with no roll needed.
          This takes approximately 1,000 successful examinations to achieve.
        </p>
      </Section>

      <Section title="Hidden Detail Format">
        <p className="text-text-muted mb-3">
          When designing rooms, items, or NPCs in the admin panel, you can add hidden details
          that only high-examine characters discover:
        </p>
        <Table
          headers={['Property', 'Meaning']}
          rows={[
            ['Text', 'What the player sees when revealed.'],
            ['Min Examine Level', 'Minimum skill level to see this detail (automatic mode).'],
            ['Mode', '"automatic" = skill threshold, "check" = roll against DC.'],
            ['DC', 'Difficulty Class for check mode. Higher = harder to reveal.'],
            ['Stat', '"INT" or "WIS" — which stat adds to the roll in check mode.'],
          ]}
        />
      </Section>
    </div>
  );
}
