import { createFileRoute } from "@tanstack/react-router";
import { PageHeader } from "../../components/PageHeader";

export const Route = createFileRoute("/docs/npc-system")({
  component: NPCSystemDoc,
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

function NPCSystemDoc() {
  return (
    <div className="management-page max-w-4xl">
      <PageHeader title="NPC System" backTo="/docs" />

      <InfoBox>
        <strong>TL;DR:</strong> NPCs are defined as <strong>templates</strong> (the static blueprint you
        create in the admin panel) and then spawned as <strong>instances</strong> (the living, breathing
        version that walks around the game). When an instance dies, it respawns after a cooldown in one of
        the rooms you specified.
      </InfoBox>

      <Section title="Lifecycle">
        <p className="text-text-muted mb-3">
          Every NPC goes through this cycle:
        </p>
        <ol className="text-sm text-text-muted space-y-2 mb-3">
          <li>
            <strong>Template:</strong> You create this in the admin panel. It holds the NPC's name, level,
            race, abilities, and loot. It never changes while the server is running.
          </li>
          <li>
            <strong>Instance:</strong> A live copy of the template that exists in the game world. It has
            current HP, a position, and a state. The server spawns it when the game starts or when a
            respawn timer fires.
          </li>
          <li>
            <strong>Combat:</strong> The instance fights players using abilities from its template. Its
            AI behavior determines whether it attacks on sight, waits to be provoked, or runs away.
          </li>
          <li>
            <strong>Death:</strong> The instance is removed from the world. It drops loot and awards XP to
            everyone who contributed damage.
          </li>
          <li>
            <strong>Respawn:</strong> After <code>respawn_cooldown</code> seconds, a brand new instance
            spawns in a random room chosen from the <code>respawn_rooms</code> list.
          </li>
        </ol>
      </Section>

      <Section title="Template Fields">
        <p className="text-text-muted mb-3">
          When you create an NPC template, here is what each field controls:
        </p>
        <Table
          headers={["Field", "What it does"]}
          rows={[
            ["Name", "The display name players see. Example: \"Junkyard Scrapper\"."],
            ["Level", "The NPC's power rating. This drives its HP pool, damage output, and XP value."],
            ["XP Value", "How many experience points the NPC is worth when defeated. Split among everyone who dealt damage."],
            ["Race", "Determines base stats, resistances, and which abilities the NPC can use. References the races table."],
            ["Respawn Cooldown", "How many seconds after death before a new instance appears."],
            [
              "Respawn Rooms",
              "A comma-separated list of room IDs. When the NPC respawns, it randomly picks one of these rooms to appear in.",
            ],
            ["Description", "Flavor text players see when they examine the NPC."],
          ]}
        />
      </Section>

      <Section title="Level Scaling">
        <p className="text-text-muted mb-3">
          NPC stats scale up with level. The game uses these formulas:
        </p>
        <div className="bg-surface-muted rounded-lg p-4 font-mono text-sm mb-3">
          npc_hp = 50 + (level × 15) + (CON × 5)
          npc_damage = 3 + (level × 1.5) + (STR × 0.8)
        </div>
        <p className="text-text-muted mb-2">
          NPCs scale linearly. Higher level means more HP, harder hits, and better loot. The race
          determines the base CON and STR values that feed into these formulas.
        </p>
        <Table
          headers={["Level Range", "Tier Name", "Typical HP", "Typical Damage"]}
          rows={[
            ["1 to 5", "Weakling", "65 to 125", "4 to 10"],
            ["6 to 15", "Standard", "140 to 275", "12 to 25"],
            ["16 to 30", "Veteran", "290 to 500", "27 to 48"],
            ["31 to 50", "Elite", "515 to 800", "50 to 78"],
            ["51+", "Boss", "815+", "80+"],
          ]}
        />
        <p className="text-text-muted text-sm">
          Tip: When designing encounters, a Standard NPC (level 6-15) is a fair fight for a solo player
          around that level. Boss-tier NPCs should require a party.
        </p>
      </Section>

      <Section title="AI Behaviors">
        <p className="text-text-muted mb-3">
          Each NPC template has an AI behavior that tells it how to react to players:
        </p>
        <Table
          headers={["Behavior", "When it triggers", "What the NPC does"]}
          rows={[
            [
              "Passive",
              "A player enters the room",
              "Watches but does not attack unless attacked first.",
            ],
            [
              "Aggressive",
              "A player enters the room",
              "Immediately attacks if the player's level is low enough (NPC level + 5 or below).",
            ],
            [
              "Flee",
              "The NPC's HP drops below 25%",
              "Tries to run to an adjacent room. Success depends on a DEX check.",
            ],
            [
              "Healer",
              "An ally's HP drops below 50%",
              "Uses a healing ability if it has one. This takes priority over attacking.",
            ],
          ]}
        />
      </Section>

      <Section title="XP Award Formula">
        <p className="text-text-muted mb-3">
          When an NPC dies, the game figures out how much XP each player gets:
        </p>
        <div className="bg-surface-muted rounded-lg p-4 font-mono text-sm mb-3">
          xp_per_player = (xp_value × damage_contribution%) × level_gap_multiplier
        </div>
        <p className="text-text-muted mb-3">
          The level gap multiplier rewards players for taking on challenges and discourages grinding easy
          enemies:
        </p>
        <Table
          headers={["Level Gap", "Multiplier", "What this means"]}
          rows={[
            ["Player is at or above NPC level", "1.0", "Full XP. You earned it."],
            ["Player is 1 to 5 levels above NPC", "0.8 to 0.5", "Slight penalty. Find harder foes."],
            ["Player is 6 to 10 levels above NPC", "0.4 to 0.1", "Heavy penalty. Time to move on."],
            ["Player is more than 10 levels above NPC", "0.0", "No XP at all. Pick on someone your own size."],
            ["Player is 5+ levels below NPC", "1.2 to 1.5", "Bonus XP for challenging yourself."],
          ]}
        />
        <p className="text-text-muted mt-3">
          Damage contribution is based on how much of the total damage you dealt during the fight. Solo
          players get 100%. Parties split it proportionally.
        </p>
      </Section>
    </div>
  );
}