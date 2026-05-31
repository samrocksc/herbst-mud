import { createFileRoute } from "@tanstack/react-router";
import { PageHeader } from "../../components/PageHeader";

export const Route = createFileRoute("/docs/item-system")({
  component: ItemSystemDoc,
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

function ItemSystemDoc() {
  return (
    <div className="management-page max-w-4xl">
      <PageHeader title="Item System" backTo="/docs" />

      <InfoBox>
        <strong>TL;DR:</strong> Items are built from <strong>templates</strong> (the static blueprint) and
        spawned as <strong>instances</strong> (the actual copy a player picks up). Equipment gives damage
        and armor. Tags connect items to abilities so you need the right gear to use certain skills. Slots
        control what can be worn at the same time.
      </InfoBox>

      <Section title="Equipment Slots">
        <p className="text-text-muted mb-3">
          Every character has these equipment slots. An item can only go into the slot it was designed for:
        </p>
        <Table
          headers={["Slot", "What goes here", "What to know"]}
          rows={[
            ["head", "Helmets, hats, headbands", "Usually focused on armor."],
            ["body", "Armor, robes, shirts", "The slot with the highest armor values."],
            ["hands", "Gloves, gauntlets", "Often adds DEX or bonus damage."],
            ["legs", "Pants, greaves", "Armor plus movement effects."],
            ["feet", "Boots, shoes", "Often adds DEX or dodge."],
            ["main_hand", "Sword, staff, pistol", "Your primary weapon. Its damage stat is used in attacks."],
            ["off_hand", "Shield, dagger, orb", "A secondary item. Can add armor or utility effects."],
            ["both_hands", "Greatsword, bow, rifle", "Takes up both main and off hand. You cannot equip anything in the off hand while this is equipped."],
          ]}
        />
      </Section>

      <Section title="Damage and Armor Calculation">
        <p className="text-text-muted mb-3">
          When a hit lands, here is how the game figures out what happens:
        </p>
        <div className="bg-surface-muted rounded-lg p-4 font-mono text-sm mb-3">
          total_damage = weapon_damage + (STR × 0.5) + dice_roll
          total_armor = sum(all_equipped_armor_values) + (CON × 0.3)
          net_damage = total_damage - total_armor
          final_damage = max(net_damage, 1)
        </div>
        <p className="text-text-muted mb-2">The key rules are:</p>
        <ul className="text-sm text-text-muted space-y-1">
          <li><strong>Weapon damage:</strong> Comes from the item's <code>damage</code> field, plus a dice roll (like 2d6 giving 2 to 12).</li>
          <li><strong>Armor:</strong> Adds up every equipped piece. If an item has negative armor, it is <strong>cursed</strong> and increases damage taken.</li>
          <li><strong>Minimum damage:</strong> Every hit deals at least 1 HP. Armor cannot reduce it to zero.</li>
          <li><strong>Magic damage:</strong> Uses INT instead of STR. It bypasses physical armor and is resisted by WIS instead.</li>
        </ul>
      </Section>

      <Section title="Tags and Abilities">
        <p className="text-text-muted mb-3">
          Item <strong>tags</strong> are keywords that link items to abilities. An ability's{" "}
          <code>required_tag</code> field means the character has to be wearing an item with that tag to
          use the ability. This is how you make a fireball require a fire gem, or a shield bash require a
          shield.
        </p>
        <Table
          headers={["Tag", "Typical items", "Which abilities need it"]}
          rows={[
            ["sword", "Longsword, katana, rapier", "Slash, Parry, Riposte"],
            ["blade", "Knives, daggers, short swords", "Stab, Backstab, Quick Strike"],
            ["staff", "Quarterstaff, bo staff", "Bash, Sweep, Deflect"],
            ["fire", "Flame torch, fire gem", "Fireball, Ignite, Flame Ward"],
            ["shield", "Buckler, tower shield", "Shield Bash, Block, Raise Shield"],
            ["potion", "Health potion, mana potion", "No tag required. Any character can use these."],
          ]}
        />
        <p className="text-text-muted mt-3">
          Tags are also useful for filtering and crafting. For example, a swordsmith needs items tagged{" "}
          <code>metal</code> and <code>blade</code> to forge new weapons.
        </p>
      </Section>

      <Section title="Item Visibility">
        <p className="text-text-muted mb-3">
          You can control whether players can see an item right away or need to work for it:
        </p>
        <Table
          headers={["Visibility", "What it means"]}
          rows={[
            ["Visible", "Players see this item when they look at the room. Anyone can find it."],
            ["Hidden", "Players only find this by using search or examine commands."],
            ["Immovable", "Cannot be picked up. It stays in the room as part of the description."],
          ]}
        />
        <p className="text-text-muted mt-3 text-sm">
          Tip: Use "Hidden" for quest items and rewards. Pair them with the examine skill's hidden details
          so players feel rewarded for being observant.
        </p>
      </Section>
    </div>
  );
}