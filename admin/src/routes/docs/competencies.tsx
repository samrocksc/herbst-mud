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
        <strong>TL;DR:</strong> Competencies are skill categories your characters can train over
        time. As characters use related abilities, they earn XP in that competency. When they hit
        an XP threshold, they level up and may unlock bonuses like extra crit chance or more damage.
      </InfoBox>

      <Section title="CompetencyCategory Entity">
        <p className="text-text-muted mb-3">
          Each CompetencyCategory defines a skill your characters can train. Here are the fields:
        </p>
        <Table
          headers={["Field", "What it does"]}
          rows={[
            ["name", "The name players see for this skill category (like Blades, Fire Magic, or Stealth)."],
            ["description", "A short description that tells players what this competency covers."],
            ["skill_name", "The internal skill name this category trains (like blades or fire_magic). This links it to abilities."],
            ["category", "A broad grouping: combat, magic, crafting, or social. Helps organize the skill list."],
          ]}
        />
      </Section>

      <Section title="CompetencyLevelThreshold Entity">
        <p className="text-text-muted mb-3">
          Level thresholds define the XP milestones for each competency. When a character's XP
          crosses a threshold, they gain that level and any bonus it grants.
        </p>
        <Table
          headers={["Field", "What it does"]}
          rows={[
            ["category_id", "Which CompetencyCategory this threshold belongs to."],
            ["level", "The level number this threshold represents (1, 2, 3, and so on)."],
            ["xp_required", "The total XP the character needs to reach this level."],
            ["bonus", "A JSON object describing the bonus the character gets when they hit this level."],
          ]}
        />
        <p className="text-text-muted mt-3">
          For example, the Blades category might have thresholds at level 1 (0 XP), level 2 (100 XP,
          grants +5% crit), and level 3 (300 XP, grants +10% crit). Bonuses stack as the character
          levels up.
        </p>
      </Section>

      <Section title="CharacterCompetency Entity">
        <p className="text-text-muted mb-3">
          This tracks where each character stands in each competency they are training.
        </p>
        <Table
          headers={["Field", "What it does"]}
          rows={[
            ["character_id", "Which character this record belongs to."],
            ["category_id", "Which competency category the character is training."],
            ["level", "The character's current level in this competency."],
            ["xp", "The character's current XP in this competency."],
          ]}
        />
      </Section>

      <Section title="How Training Works">
        <p className="text-text-muted mb-3">
          When a character uses an ability that requires a skill, here is what happens:
        </p>
        <ul className="list-disc pl-6 text-text-muted mb-4 space-y-1">
          <li>The game checks whether the character meets the skill level requirement for that ability.</li>
          <li>If they use the ability successfully, they earn XP in the related competency. Hitting something with a sword trains Blades, for example.</li>
          <li>When their XP crosses a level threshold, their competency level goes up.</li>
          <li>Level thresholds can grant bonuses like extra crit chance, more damage, and so on.</li>
        </ul>
      </Section>

      <Section title="Skill Page">
        <p className="text-text-muted mb-3">
          You can manage competency categories and their level thresholds at /skills in the admin
          panel. To see how a specific character is progressing, check their character detail page.
        </p>
      </Section>

      <Section title="Admin API">
        <Table
          headers={["Method", "Endpoint", "What it does"]}
          rows={[
            ["GET", "/api/competency-categories", "List all competency categories."],
            ["POST", "/api/competency-categories", "Create a new competency category."],
            ["GET", "/api/competency-categories/:id/thresholds", "See the level thresholds for a specific category."],
            ["POST", "/api/competency-categories/:id/thresholds", "Add a new level threshold to a category."],
            ["GET", "/api/characters/:id/competencies", "See how far a character has progressed in their competencies."],
          ]}
        />
      </Section>
    </div>
  );
}