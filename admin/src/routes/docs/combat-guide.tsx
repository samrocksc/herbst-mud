import { createFileRoute } from "@tanstack/react-router";
import { PageHeader } from "../../components/PageHeader";

export const Route = createFileRoute("/docs/combat-guide")({
  component: CombatGuideDoc,
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

function CombatGuideDoc() {
  return (
    <div className="management-page max-w-4xl">
      <PageHeader title="Combat Guide" backTo="/docs" />

      <InfoBox>
        <strong>TL;DR:</strong> Combat runs on a tick system (1.5 seconds per tick). Characters take
        actions by spending ticks. Your six core stats (STR, DEX, CON, INT, WIS, CHA) determine how hard
        you hit, how well you dodge, and how potent your magic is. Press{" "}
        <kbd className="bg-surface-dark text-text-inverse px-1 rounded">1</kbd> through{" "}
        <kbd className="bg-surface-dark text-text-inverse px-1 rounded">5</kbd> during a fight to use your
        classless abilities.
      </InfoBox>

      <Section title="Tick System">
        <p className="text-text-muted mb-3">
          Every 1.5 seconds, the combat clock ticks. Each action costs a certain number of ticks. When your
          accumulated tick debt is paid off, you get to act again. This means faster actions let you act
          more often, while heavy abilities make you wait longer.
        </p>
        <Table
          headers={["Action", "Tick Cost", "What you need to know"]}
          rows={[
            ["Attack", "1", "Your basic swing or shot. Damage comes from your weapon plus your stats."],
            ["Defend", "0", "Raises your armor briefly without costing you a turn. Always worth using."],
            ["Flee", "1", "Try to escape. Your DEX is compared to the enemy's level to see if you succeed."],
            ["Use Item", "1", "Drink a potion, throw a grenade, and so on."],
            ["Activate Skill", "1 to 3", "Classless skills cost 1 tick. Passive abilities can cost up to 3."],
          ]}
        />
      </Section>

      <Section title="Character Stats">
        <p className="text-text-muted mb-3">
          Every character has six stats. These are the foundation for almost everything in combat, from
          damage rolls to hit chance to how many spells you can cast:
        </p>
        <Table
          headers={["Stat", "Abbrev", "What it does for you"]}
          rows={[
            ["Strength", "STR", "More melee damage and higher carrying capacity."],
            ["Dexterity", "DEX", "Better hit chance, better dodge chance, faster attacks."],
            ["Constitution", "CON", "Bigger HP pool, bigger stamina pool, better resistance to stun and poison."],
            ["Intelligence", "INT", "Stronger magic, better tech skills, larger mana pool."],
            ["Wisdom", "WIS", "Keener perception, stronger wind magic and healing, better magical resistance."],
            ["Charisma", "CHA", "Better shop prices, friendlier NPCs, larger party size."],
          ]}
        />
      </Section>

      <Section title="Damage Formula">
        <p className="text-text-muted mb-3">
          When you hit someone, the game calculates damage like this:
        </p>
        <div className="bg-surface-muted rounded-lg p-4 font-mono text-sm mb-3">
          damage = weapon_base + (STR × 0.5) + random(1, weapon_dice_sides) - target_armor
        </div>
        <p className="text-text-muted mb-2">
          If the ability has a <strong>scaling_stat</strong> set, the formula changes to:
        </p>
        <div className="bg-surface-muted rounded-lg p-4 font-mono text-sm mb-3">
          damage = base + (stat_value × scaling_percent × base) - target_armor
        </div>
        <ul className="text-sm text-text-muted space-y-1">
          <li><strong>Weapon base:</strong> The damage number on your weapon item.</li>
          <li><strong>Random roll:</strong> A number between 1 and the weapon's dice sides. A d8 gives you 1 through 8.</li>
          <li><strong>Armor:</strong> Add up the armor values on every piece of gear the target has equipped.</li>
          <li><strong>Minimum damage:</strong> You always deal at least 1 HP, no matter how tough the armor.</li>
        </ul>
      </Section>

      <Section title="Hit Chance and Dodge">
        <p className="text-text-muted mb-3">
          Every attack starts with a 50% chance to hit. DEX shifts that number:
        </p>
        <div className="bg-surface-muted rounded-lg p-4 font-mono text-sm mb-3">
          hit_chance = 50% + (attacker_DEX - target_DEX) × 2.5%
        </div>
        <p className="text-text-muted mb-2">
          Each point of DEX difference between you and your target moves hit chance by 2.5%.
        </p>
        <ul className="text-sm text-text-muted space-y-1 mb-3">
          <li>Minimum hit chance: 5%. There is always a small chance you connect.</li>
          <li>Maximum hit chance: 95%. Nothing is guaranteed.</li>
          <li><strong>Back-off</strong> skill: Gives you 100% dodge for 1 entire round. You avoid everything.</li>
          <li><strong>Concentrate</strong> skill: Adds your WIS to hit rolls for 4 rounds.</li>
        </ul>
      </Section>

      <Section title="Critical Hits">
        <p className="text-text-muted mb-3">
          A critical hit deals <strong>150% damage</strong> before armor is subtracted. Criticals can also
          trigger <code>on_crit</code> proc events on your abilities.
        </p>
        <div className="bg-surface-muted rounded-lg p-4 font-mono text-sm mb-3">
          crit_chance = 5% + (DEX - 10) × 0.5%
        </div>
        <p className="text-text-muted">
          Everyone starts with a 5% crit chance. Each point of DEX above 10 adds 0.5%. Some gear and buffs
          can push this even higher.
        </p>
      </Section>

      <Section title="Combat Flow Example">
        <p className="text-text-muted mb-3">
          Here is what a typical 5-round fight looks like:
        </p>
        <ol className="text-sm text-text-muted space-y-2 mb-3">
          <li>
            <strong>Round 1:</strong> You press <kbd>1</kbd> (Concentrate). Cost: 10 MP, 0 ticks. Now
            you have +WIS to hit for 4 rounds.
          </li>
          <li>
            <strong>Round 2:</strong> You attack. Concentrate is boosting your accuracy. DEX 14 vs enemy
            DEX 10 gives you a 60% base hit. With WIS bonus you are around 70%.
          </li>
          <li>
            <strong>Round 3:</strong> The enemy attacks. You press <kbd>3</kbd> (Back-off). Cost: 25 SP.
            You dodge 100% of incoming attacks this round.
          </li>
          <li>
            <strong>Round 4:</strong> You press <kbd>2</kbd> (Haymaker). Cost: 15 SP, 1 tick. STR 16 means
            damage = 50 + (16 × 0.05 × 50) = 90. But Haymaker lowers your DEX, so your next hit chance drops.
          </li>
          <li>
            <strong>Round 5:</strong> You press <kbd>4</kbd> (Scream). Cost: 5 MP + 10 SP. You gain +DEX/STR
            and the enemy loses WIS/INT for 2 rounds. Your next Haymaker is going to hurt.
          </li>
        </ol>
      </Section>

      <Section title="Status Effects">
        <Table
          headers={["Effect", "What happens", "How long it lasts"]}
          rows={[
            ["Stunned", "You cannot act. You skip your next turn entirely.", "1 round (from Slap)"],
            ["Poisoned", "You lose HP every tick.", "Until cured or the duration runs out"],
            ["Buffed", "One or more of your stats are raised.", "Duration set by the ability (in ticks)"],
            ["Debuffed", "One or more of your stats are lowered.", "Duration set by the ability (in ticks)"],
            ["Concentrating", "Your WIS is added to hit rolls.", "4 rounds"],
            ["Haymaker stance", "Your STR boosts damage, but your DEX lowers hit chance.", "1 attack"],
            ["Screaming", "You gain DEX and STR but lose WIS and INT.", "2 rounds"],
          ]}
        />
      </Section>
    </div>
  );
}