import { createFileRoute } from "@tanstack/react-router";
import { PageHeader } from "../../components/PageHeader";

export const Route = createFileRoute("/docs/competencies")({
  component: CompetenciesDoc,
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

function CompetenciesDoc() {
  return (
    <div className="management-page max-w-4xl">
      <PageHeader title="Competencies" backTo="/docs" />

      <InfoBox>
        <strong>TL;DR:</strong> Competencies are trainable skill categories with level thresholds.
        Characters earn XP in a competency by using related abilities. Level thresholds define
        XP gates and unlock bonuses.
      </InfoBox>

      <Section title="CompetencyCategory Entity">
        <Table
          headers={["Field", "Description"]}
          rows={[
            ["name", "Category name (e.g., Blades, Fire Magic, Stealth)."],
            ["description", "Player-facing description of the competency."],
            ["skill_name", "The underlying skill being trained (e.g., blades, fire_magic)."],
            ["category", "Grouping: combat, magic, crafting, social."],
          ]}
        />
      </Section>

      <Section title="CompetencyLevelThreshold Entity">
        <Table
          headers={["Field", "Description"]}
          rows={[
            ["category_id", "Which CompetencyCategory this threshold belongs to."],
            ["level", "The level this threshold represents (1, 2, 3...)."],
            ["xp_required", "Total XP needed to reach this level."],
            ["bonus", "Bonus granted when this level is reached (JSON)."],
          ]}
        />
        <p className="text-text-muted mt-3">
          Example: Blades category has thresholds at level 1 (0 XP), level 2 (100 XP, +5% crit),
          level 3 (300 XP, +10% crit), etc. Bonuses stack as level increases.
        </p>
      </Section>

      <Section title="CharacterCompetency Entity">
        <Table
          headers={["Field", "Description"]}
          rows={[
            ["character_id", "Which character has this competency."],
            ["category_id", "Which competency category."],
            ["level", "Current level in this competency."],
            ["xp", "Current XP in this competency."],
          ]}
        />
      </Section>

      <Section title="How Training Works">
        <p className="text-text-muted mb-3">
          When a character uses an ability that has a skill requirement:
        </p>
        <ul className="list-disc pl-6 text-text-muted mb-4 space-y-1">
          <li>The ability checks if the character meets the skill level.</li>
          <li>On use, the character earns XP in that skill (e.g., hitting with a sword trains Blades).</li>
          <li>When XP crosses a threshold, level increases.</li>
          <li>Level thresholds can grant bonuses (crit chance, damage, etc.).</li>
        </ul>
      </Section>

      <Section title="Skill Page">
        <p className="text-text-muted mb-3">
          The admin Skills page at /skills manages competency categories and their
          level thresholds. Character competency progress is visible on the
          character detail page.
        </p>
      </Section>

      <Section title="Admin API">
        <Table
          headers={["Method", "Endpoint", "Description"]}
          rows={[
            ["GET", "/api/competency-categories", "List all competency categories."],
            ["POST", "/api/competency-categories", "Create a competency category."],
            ["GET", "/api/competency-categories/:id/thresholds", "Get level thresholds for a category."],
            ["POST", "/api/competency-categories/:id/thresholds", "Add a level threshold."],
            ["GET", "/api/characters/:id/competencies", "Get character competency progress."],
          ]}
        />
      </Section>
    </div>
  );
}
