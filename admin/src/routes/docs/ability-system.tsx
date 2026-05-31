import { createFileRoute } from "@tanstack/react-router";
import { PageHeader } from "../../components/PageHeader";

export const Route = createFileRoute("/docs/ability-system")({
  component: AbilitySystemDoc,
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

function AbilitySystemDoc() {
  return (
    <div className="management-page max-w-4xl">
      <PageHeader title="Ability System" backTo="/docs" />

      <InfoBox>
        <strong>TL;DR:</strong> Abilities are the combat actions players press{" "}
        <kbd className="bg-surface-dark text-text-inverse px-1 rounded">1</kbd> through{" "}
        <kbd className="bg-surface-dark text-text-inverse px-1 rounded">5</kbd> to use. Each one costs
        mana or stamina, has a cooldown, and gets stronger as the character's stats go up. Every ability
        is built from one or more <strong>effects</strong> that tell the game what to actually do (deal
        damage, heal, buff, stun, and so on). The <strong>5 classless abilities</strong> are available
        to every character, no class required.
      </InfoBox>

      <Section title="Four-Domain Model">
        <p className="text-text-muted mb-3">
          The system is built around four domains that each handle a different part of how characters work:
        </p>
        <Table
          headers={["Domain", "What it covers", "Examples", "Where it lives"]}
          rows={[
            ["Abilities", "Actions your character can perform in combat", "Concentrate, Haymaker, Fireball", "Ability entity"],
            ["Skills", "Proficiencies that level up with training", "Blades, Staves, Light Armor", "Character columns"],
            ["Stats", "Core numeric attributes that drive everything else", "Strength, Dexterity, Wisdom", "Character fields"],
            ["Effects", "The actual thing that happens when an ability fires", "Damage, Heal, Buff, Stun", "AbilityEffect entity"],
          ]}
        />
      </Section>

      <Section title="The 5 Classless Abilities">
        <p className="text-text-muted mb-3">
          Every character starts with these 5 combat abilities bound to slots 1 through 5. Each ability
          is made up of one or more effects defined in the AbilityEffect entity.
        </p>
        <Table
          headers={["Slot", "Ability", "What it does", "Mana", "Stamina", "Cooldown"]}
          rows={[
            ["1", "Concentrate", "Boosts your accuracy for 4 ticks, scaled by WIS", "10", "0", "8s"],
            ["2", "Haymaker", "Deals STR-scaled damage to the enemy, but debuffs you for 1 tick", "0", "15", "6s"],
            ["3", "Back-off", "You dodge everything for 1 tick. Great for emergencies.", "0", "25", "10s"],
            ["4", "Scream", "Buffs you (CON scaled) and debuffs the enemy (CON scaled) for 2 ticks", "5", "10", "12s"],
            ["5", "Slap", "Stuns the enemy for 1 tick, contested by DEX", "0", "12", "8s"],
          ]}
        />
      </Section>

      <Section title="Effect Types">
        <p className="text-text-muted mb-3">
          Effects are the building blocks that make abilities work. You can stack multiple effects on a
          single ability, each with its own target, value, and scaling stat. Here are all the effect
          types you can use:
        </p>
        <Table
          headers={["Type", "What it does", "Who you can target"]}
          rows={[
            ["damage", "Takes HP away from the target. Pick a damage subtype to pair it with resistances.", "enemy, area, random_enemy"],
            ["heal", "Restores HP to yourself or an ally.", "self, ally"],
            ["buff", "Gives a positive status effect that lasts for a set number of ticks.", "self, ally"],
            ["debuff", "Slaps a negative status effect on the target for a duration.", "enemy"],
            ["dot", "Damage Over Time. Deals damage each combat tick for the whole duration.", "enemy"],
            ["hot", "Heal Over Time. Restores HP each combat tick for the whole duration.", "self, ally"],
            ["stun", "Target skips their next turn. The target can resist with a CON contest.", "enemy"],
            ["accuracy_boost", "Raises your hit chance for a set duration.", "self"],
            ["dodge_all", "You avoid every attack for the duration.", "self"],
            ["set_bind_point", "Sets your respawn point to the current room.", "self"],
          ]}
        />
      </Section>

      <Section title="Damage Subtypes">
        <p className="text-text-muted mb-3">
          When you create a damage effect, you can assign a subtype. This lets you set up resistances
          and weaknesses later (for example, fire-resistant armor or a cold-vulnerable monster):
        </p>
        <Table
          headers={["Subtype", "When to use it"]}
          rows={[
            ["slashing", "Swords, axes, anything that cuts. Countered by heavy armor."],
            ["piercing", "Daggers, arrows, anything that punctures. Countered by cloth armor."],
            ["bludgeoning", "Clubs, fists, anything that bludgeons. Countered by light armor."],
            ["fire", "Fireballs, dragon breath, anything that burns."],
            ["cold", "Ice spells, frost-enchanted weapons."],
            ["lightning", "Shock spells, storm-based abilities."],
            ["poison", "Toxins, venoms, acids."],
            ["psychic", "Mind attacks, illusions, mental assaults."],
          ]}
        />
      </Section>

      <Section title="Ability Fields Reference">
        <Table
          headers={["Field", "What it means"]}
          rows={[
            ["Name", "The command players type to use it, like \"concentrate\"."],
            ["Description", "Flavor text players see when they examine the ability."],
            ["Ability Type", "What kind of ability it is: combat, magic, utility, healing, support, or defensive."],
            ["Ability Class", "How it activates: active (press to use), passive (triggers on its own), or toggle (flip on/off)."],
            ["Required Tag", "An item tag the character must have equipped to use this ability."],
            ["Level Req", "The minimum character level needed to learn this ability."],
            ["Cooldown (s)", "How many seconds before the ability can be used again."],
            ["Mana Cost", "How much MP the ability drains when activated."],
            ["Stamina Cost", "How much SP the ability drains when activated."],
            ["HP Cost", "How much self-damage the ability deals to cast it."],
            ["Proc Chance", "A 0.0 to 1.0 probability that the ability triggers. Passives only."],
            ["Proc Event", "When the game rolls for a proc: on_hit, on_hit_received, on_crit, or on_kill."],
          ]}
        />
      </Section>

      <Section title="Scaling Formula">
        <p className="text-text-muted mb-3">
          When an ability scales with a stat, the game uses this formula to figure out the final value:
        </p>
        <div className="bg-surface-muted rounded-lg p-4 font-mono text-sm mb-3">
          final_value = base_value + (stat_modifier × scaling_ratio × base_value)
        </div>
        <p className="text-text-muted mb-3">
          The stat modifier uses the standard D&D formula:{" "}
          <code className="bg-surface-dark text-text-inverse px-1 rounded">stat_modifier = (stat - 10) / 2</code>.
        </p>
        <p className="text-text-muted mb-3">Here is a worked example with Haymaker (base 15, STR scaling, ratio 0.5):</p>
        <ul className="text-sm text-text-muted space-y-1 mb-3">
          <li>A character with STR 18 has a +4 modifier.</li>
          <li>Final damage: <code className="bg-surface-dark text-text-inverse px-1 rounded">15 + (4 × 0.5 × 15) = 15 + 30 = 45</code></li>
        </ul>
        <p className="text-text-muted text-sm">
          Tip: When designing new abilities, a higher scaling_ratio makes the ability reward stat investment
          more. A ratio of 0.5 means a +4 modifier doubles the base value. Lower ratios keep things flatter.
        </p>
      </Section>

      <Section title="Ability Combos">
        <p className="text-text-muted mb-3">
          Players will figure out synergies between abilities. Here are the intended combos for the classless
          set:
        </p>
        <Table
          headers={["Combo", "Why it works"]}
          rows={[
            ["Scream then Haymaker", "Scream buffs your STR. Haymaker scales off STR. You trade INT/WIS for big damage."],
            ["Concentrate then Back-off", "Concentrate boosts your accuracy. When things get rough, Back-off gives you a dodge turn."],
            ["Slap then Haymaker", "Stunned enemies cannot act. That means a free Haymaker with zero risk of retaliation."],
            ["Back-off, then potion, then attack", "The dodge round buys you time to heal up. Come back swinging with full resources."],
          ]}
        />
      </Section>

      <Section title="Passive Abilities (Formerly Talents)">
        <p className="text-text-muted mb-3">
          Passive abilities use <code className="bg-surface-dark text-text-inverse px-1 rounded">ability_class=passive</code> and
          fire automatically based on their <strong>proc_event</strong> (like on_hit or on_crit) with a{" "}
          <strong>proc_chance</strong> between 0.0 and 1.0. When the event happens and the chance roll
          succeeds, the ability's effects are applied automatically.
        </p>
        <p className="text-text-muted text-sm">
          Tip: You used to create these as "Talents." Now they are just abilities with a different ability
          class. The old Talent system has been folded into this one.
        </p>
      </Section>
    </div>
  );
}