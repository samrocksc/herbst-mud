import { createFileRoute } from "@tanstack/react-router";
import { PageHeader } from "../../components/PageHeader";

export const Route = createFileRoute("/docs/examine-skill")({
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
        <strong>TL;DR:</strong> The <strong>examine</strong> skill lets players look closely at rooms,
        items, and NPCs to find hidden details that everyone else misses. Higher skill level means more
        secrets revealed. It runs on INT and WIS. The max level is 100.
      </InfoBox>

      <Section title="How Examine Works">
        <p className="text-text-muted mb-3">
          When a player types <code>examine</code> (or <code>examine [target]</code>), the game checks
          whether they notice any hidden details in what they are looking at. The roll goes like this:
        </p>
        <div className="bg-surface-muted rounded-lg p-4 font-mono text-sm mb-3">
          examine_roll = skill_level + INT + (WIS / 2) + random(1, 10)
        </div>
        <p className="text-text-muted mb-2">
          There are two ways a hidden detail can be revealed:
        </p>
        <ul className="text-sm text-text-muted space-y-1 mb-3">
          <li><strong>Automatic:</strong> If the player's skill level meets or beats the detail's minimum level, they see it. No roll needed.</li>
          <li><strong>Check:</strong> The game rolls against a difficulty class (DC). If the examine roll is equal to or higher than the DC, the detail is revealed. INT or WIS can be used for the check.</li>
        </ul>
        <p className="text-text-muted text-sm">
          Tip: When you set up hidden details, think about which approach fits. Use "automatic" for
          things any observant character would eventually notice. Use "check" for genuinely hard-to-find
          secrets that reward dedicated examine builds.
        </p>
      </Section>

      <Section title="Skill Levels and Bonuses">
        <p className="text-text-muted mb-3">
          As the examine skill climbs, players unlock progressively better discoveries:
        </p>
        <Table
          headers={["Examine Level", "Bonus", "What players can find"]}
          rows={[
            ["0 to 25", "0%", "Just the basic description. Nothing hidden yet."],
            ["26 to 50", "10%", "Hidden items in rooms start appearing."],
            ["51 to 75", "25%", "NPC weaknesses and item stats become visible."],
            ["76 to 90", "50%", "Secret exits and quest hints reveal themselves."],
            ["91 to 100", "75%", "Everything: lore, backstory, all the secrets."],
          ]}
        />
        <p className="text-text-muted mt-3">
          The bonus percentage also improves shop prices (because you can appraise items better) and loot
          quality (because you spot the valuable stuff others miss).
        </p>
      </Section>

      <Section title="XP Rewards">
        <p className="text-text-muted mb-3">
          Players earn XP for paying attention to the world. Here is how the rewards break down:
        </p>
        <Table
          headers={["Action", "XP Gained"]}
          rows={[
            ["First time examining a room", "+1 XP"],
            ["Discovering a hidden detail", "+2 XP"],
            ["Revealing a hidden exit", "+5 XP"],
            ["Decrypting a secret message", "+10 XP"],
          ]}
        />
        <div className="bg-surface-muted rounded-lg p-4 font-mono text-sm mt-3 mb-3">
          level_up_every = 10 XP
          max_level = 100
        </div>
        <p className="text-text-muted">
          At max level (100), examine automatically reveals every hidden detail with no roll needed. Getting
          there takes roughly 1,000 successful examinations. That sounds like a lot, but curious players who
          examine everything will get there naturally over time.
        </p>
      </Section>

      <Section title="Hidden Detail Format">
        <p className="text-text-muted mb-3">
          When you create rooms, items, or NPCs in the admin panel, you can attach hidden details that only
          high-examine characters will discover. Here are the fields you need to fill in:
        </p>
        <Table
          headers={["Property", "What to put here"]}
          rows={[
            ["Text", "The description the player sees when the detail is revealed."],
            ["Min Examine Level", "The minimum skill level needed to see this. Used in automatic mode."],
            ["Mode", "Pick \"automatic\" for a skill threshold, or \"check\" for a roll against a DC."],
            ["DC", "The Difficulty Class for check mode. Higher numbers mean harder to find."],
            ["Stat", "Which stat to roll with in check mode: \"INT\" or \"WIS\"."],
          ]}
        />
        <p className="text-text-muted mt-3 text-sm">
          Tip: Mix automatic and check-mode details in the same room. That way both casual players and
          dedicated examiner builds get rewarded for exploring.
        </p>
      </Section>
    </div>
  );
}