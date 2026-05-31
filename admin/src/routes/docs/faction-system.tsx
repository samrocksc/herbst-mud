import { createFileRoute } from "@tanstack/react-router";
import { PageHeader } from "../../components/PageHeader";

export const Route = createFileRoute("/docs/faction-system")({
  component: FactionSystemDoc,
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

function FactionSystemDoc() {
  return (
    <div className="management-page max-w-4xl">
      <PageHeader title="Faction System" backTo="/docs" />

      <InfoBox>
        <strong>TL;DR:</strong> Factions are in-game groups that track <strong>standing</strong> with each
        player on a scale from -100 to +100. Standing changes how NPCs treat you: whether they attack you,
        what prices you get at shops, which quests you can access, and whether you unlock special faction-only
        abilities. You earn or lose standing based on what you do.
      </InfoBox>

      <Section title="Standing Scale">
        <p className="text-text-muted mb-3">
          Standing runs from -100 (hated) to +100 (revered). Here is what each range means for the player:
        </p>
        <Table
          headers={["Standing", "Label", "What happens"]}
          rows={[
            ["+81 to +100", "Revered", "50% shop discount. Exclusive quests. Faction abilities unlocked."],
            ["+51 to +80", "Honored", "25% discount. Better quests. Friendly NPCs will help you in combat."],
            ["+11 to +50", "Friendly", "10% discount. Standard quests become available."],
            ["-10 to +10", "Neutral", "No modifiers. Just normal interactions."],
            ["-50 to -11", "Unfriendly", "25% markup on prices. Some NPCs refuse to deal with you."],
            ["-80 to -51", "Hostile", "50% markup. Some NPCs will attack you on sight."],
            ["-100 to -81", "Hated", "Kill on sight. Shops are closed to you. A bounty is placed on your head."],
          ]}
        />
      </Section>

      <Section title="Faction Categories">
        <p className="text-text-muted mb-3">
          When you create a faction, pick a category that fits its theme. The category helps the game know
          what kind of behavior gains or loses standing:
        </p>
        <Table
          headers={["Category", "Examples", "What they care about"]}
          rows={[
            ["Political", "City councils, noble houses", "Territory control, taxes, laws"],
            ["Religious", "Churches, cults, monastic orders", "Converting followers, completing holy quests"],
            ["Criminal", "Thieves guilds, smuggler rings", "Profit, territory, avoiding the guards"],
            ["Guild", "Craft guilds, mercenary companies", "Skill mastery, contracts, building reputation"],
            ["Race", "Elf clans, dwarf holds", "Species pride, ancient grudges"],
          ]}
        />
      </Section>

      <Section title="Changing Standing">
        <p className="text-text-muted mb-3">
          Players shift their standing through the choices they make in the world. Here is how much each
          action moves the needle:
        </p>
        <Table
          headers={["Action", "Standing Change", "Which faction cares"]}
          rows={[
            ["Complete a faction quest", "+10 to +25", "The faction that gave the quest"],
            ["Kill a faction's enemy", "+5", "Whichever faction hates the target"],
            ["Kill a faction member", "-20 to -50", "The victim's faction"],
            ["Donate items or gold", "+1 to +5", "The faction receiving the donation"],
            ["Betray a faction quest", "-30", "The faction you betrayed"],
            ["Wear a faction emblem", "+1 per day (max +10)", "The faction on the emblem"],
            ["Attack a faction's ally", "-10", "The ally's faction"],
          ]}
        />
        <p className="text-text-muted mt-3 text-sm">
          Tip: Betraying a quest hurts a lot. Make sure players feel the weight of that choice in your
          story design. A -30 swing is hard to recover from.
        </p>
      </Section>

      <Section title="Faction Abilities">
        <p className="text-text-muted mb-3">
          When a player reaches <strong>Revered</strong> standing (81+) with a faction, they earn a unique
          ability that nobody else gets:
        </p>
        <Table
          headers={["Faction", "Ability", "What it does"]}
          rows={[
            ["Surf Wardens", "Wave Rush", "Water-based knockback. Scales with DEX."],
            ["Dune Traders", "Sand Veil", "Go temporarily invisible in desert zones."],
            ["Tinkerers", "Gadget Barrage", "Deals tech damage and stuns. Scales with INT."],
            ["The Vine Climb", "Overgrowth", "Roots enemies in place. Scales with WIS."],
          ]}
        />
        <p className="text-text-muted mt-3">
          These abilities show up in the character's ability list once standing hits 81. They can be
          equipped in any skill slot, just like the classless abilities.
        </p>
      </Section>
    </div>
  );
}