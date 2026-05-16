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
        <strong>TL;DR:</strong> Factions are in-game groups with <strong>standing</strong> scores
        (−100 to +100). Standing affects NPC aggression, shop prices, quest access, and
        faction-specific abilities. Players gain/lose standing through actions.
      </InfoBox>

      <Section title="Standing Scale">
        <Table
          headers={["Standing", "Label", "Effects"]}
          rows={[
            ["+81 to +100", "Revered", "50% shop discount. Exclusive quests. Faction abilities unlocked."],
            ["+51 to +80", "Honored", "25% discount. Better quests. NPCs help in combat."],
            ["+11 to +50", "Friendly", "10% discount. Standard quests available."],
            ["−10 to +10", "Neutral", "No modifiers. Standard interactions."],
            ["−50 to −11", "Unfriendly", "25% price markup. NPCs refuse some services."],
            ["−80 to −51", "Hostile", "50% markup. NPCs may attack on sight."],
            ["−100 to −81", "Hated", "Kill on sight. Shops closed. Bounty placed on player."],
          ]}
        />
      </Section>

      <Section title="Faction Categories">
        <Table
          headers={["Category", "Examples", "Typical Goals"]}
          rows={[
            ["Political", "City councils, noble houses", "Territory control, taxes, laws"],
            ["Religious", "Churches, cults, monastic orders", "Convert followers, holy quests"],
            ["Criminal", "Thieves guilds, smugglers", "Profit, territory, avoiding guards"],
            ["Guild", "Craft guilds, mercenary companies", "Skill mastery, contracts, reputation"],
            ["Race", "Elf clans, dwarf holds", "Species pride, ancient grudges"],
          ]}
        />
      </Section>

      <Section title="Changing Standing">
        <p className="text-text-muted mb-3">Standing changes through player actions:</p>
        <Table
          headers={["Action", "Standing Change", "Faction"]}
          rows={[
            ["Complete faction quest", "+10 to +25", "Quest-giver's faction"],
            ["Kill faction enemy", "+5", "Enemy of target's faction"],
            ["Kill faction member", "−20 to −50", "Victim's faction"],
            ["Donate items/gold", "+1 to +5", "Receiving faction"],
            [" betray faction quest", "−30", "Betrayed faction"],
            ["Wear faction emblem", "+1 per day (max +10)", "Emblem faction"],
            ["Attack faction ally", "−10", "Ally's faction"],
          ]}
        />
      </Section>

      <Section title="Faction Abilities">
        <p className="text-text-muted mb-3">
          At <strong>Revered</strong> standing, factions grant unique abilities:
        </p>
        <Table
          headers={["Faction", "Ability", "Effect"]}
          rows={[
            ["Surf Wardens", "Wave Rush", "Water-based knockback. DEX scaling."],
            ["Dune Traders", "Sand Veil", "Temporary invisibility in desert zones."],
            ["Tinkerers", "Gadget Barrage", "Tech damage + stun. INT scaling."],
            ["The Vine Climb", "Overgrowth", "Roots enemies in place. WIS scaling."],
          ]}
        />
        <p className="text-text-muted mt-3">
          These abilities appear in the character's ability list once standing reaches 81+.
          They can be equipped in any skill slot like classless skills.
        </p>
      </Section>
    </div>
  );
}
