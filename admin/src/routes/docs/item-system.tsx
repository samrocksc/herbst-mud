import { createFileRoute } from '@tanstack/react-router'
import { PageHeader } from '../../components/PageHeader'

export const Route = createFileRoute('/docs/item-system')({
  component: ItemSystemDoc,
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

function ItemSystemDoc() {
  return (
    <div className="management-page max-w-4xl">
      <PageHeader title="Item System" backTo="/docs" />

      <InfoBox>
        <strong>TL;DR:</strong> Items have templates (static blueprints) and instances (specific copies).
        Equipment provides damage/armor. Tags connect items to abilities. Slots limit what can be equipped.
      </InfoBox>

      <Section title="Equipment Slots">
        <Table
          headers={['Slot', 'What Goes Here', 'Notes']}
          rows={[
            ['head', 'Helmets, hats, headbands', 'Usually armor-focused.'],
            ['body', 'Armor, robes, shirts', 'Highest armor value slot.'],
            ['hands', 'Gloves, gauntlets', 'Often adds DEX or damage.'],
            ['legs', 'Pants, greaves', 'Armor + movement.'],
            ['feet', 'Boots, shoes', 'Often adds DEX or dodge.'],
            ['main_hand', 'Sword, staff, pistol', 'Primary weapon. Damage used in attacks.'],
            ['off_hand', 'Shield, dagger, orb', 'Secondary item. Can add armor or utility.'],
            ['both_hands', 'Greatsword, bow, rifle', 'Occupies main + off hand. Cannot equip off_hand item.'],
          ]}
        />
      </Section>

      <Section title="Damage & Armor Calculation">
        <div className="bg-surface-muted rounded-lg p-4 font-mono text-sm mb-3">
          total_damage = weapon_damage + (STR × 0.5) + dice_roll
          total_armor = sum(all_equipped_armor_values) + (CON × 0.3)
          net_damage = total_damage − total_armor
          final_damage = max(net_damage, 1)
        </div>
        <p className="text-text-muted mb-2">Key rules:</p>
        <ul className="text-sm text-text-muted space-y-1">
          <li><strong>Weapon damage:</strong> Item's <code>damage</code> field + dice (e.g. 2d6 = 2–12).</li>
          <li><strong>Armor:</strong> Sum of all equipped items. Negative armor = <strong>cursed</strong> (increases damage taken).</li>
          <li><strong>Minimum damage:</strong> 1 HP per hit. Armor cannot reduce below 1.</li>
          <li><strong>Magic damage:</strong> Uses INT instead of STR. Bypasses physical armor (uses WIS resistance).</li>
        </ul>
      </Section>

      <Section title="Tags & Abilities">
        <p className="text-text-muted mb-3">
          Item <strong>tags</strong> are keywords that connect items to abilities. Abilities have a
          <code>required_tag</code> field — the character must have an item with that tag equipped.
        </p>
        <Table
          headers={['Tag', 'Typical Items', 'Required By']}
          rows={[
            ['sword', 'Longsword, katana, rapier', 'Slash, Parry, Riposte abilities'],
            ['blade', 'Knives, daggers, swords', 'Stab, Backstab, Quick Strike'],
            ['staff', 'Quarterstaff, bo staff', 'Bash, Sweep, Deflect'],
            ['fire', 'Flame torch, fire gem', 'Fireball, Ignite, Flame Ward'],
            ['shield', 'Buckler, tower shield', 'Shield Bash, Block, Raise Shield'],
            ['potion', 'Health potion, mana potion', 'Item use (no tag requirement)'],
          ]}
        />
        <p className="text-text-muted mt-3">
          Tags are also used for filtering and crafting recipes. A swordsmith needs items tagged
          <code>metal</code> and <code>blade</code> to forge new weapons.
        </p>
      </Section>

      <Section title="Item Visibility">
        <Table
          headers={['Visibility', 'Meaning']}
          rows={[
            ['Visible', 'Players see this item when they look at the room.'],
            ['Hidden', 'Only found via search/examine commands.'],
            ['Immovable', 'Cannot be picked up. Part of the room description.'],
          ]}
        />
      </Section>
    </div>
  )
}
