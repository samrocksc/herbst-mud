import { createFileRoute } from "@tanstack/react-router";
import { PageHeader } from "../../components/PageHeader";

export const Route = createFileRoute("/docs/effect-system")({
  component: EffectSystemDoc,
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

function EffectSystemDoc() {
  return (
    <div className="management-page max-w-4xl">
      <PageHeader title="Effect System" backTo="/docs" />

      <InfoBox>
        <strong>TL;DR:</strong> Effects are reusable building blocks for game logic. You define an
        Effect with a name, type, value, and duration. Then you attach it to abilities, quests, or
        events using EffectHooks. When something fires in-game, ActiveEffect tracks the live instance
        on a character.
      </InfoBox>

      <Section title="Effect Entity">
        <p className="text-text-muted mb-3">
          An Effect is the template. You fill in these fields to define what it does:
        </p>
        <Table
          headers={["Field", "What it means"]}
          rows={[
            ["name", "A unique identifier for the effect (e.g., burning, regenerate)."],
            ["description", "What players see when this effect is on them."],
            ["effect_type", "The category: damage, heal, buff, debuff, dot, hot, stun, accuracy_boost, dodge_all, or set_bind_point."],
            ["value", "The base number this effect uses (damage amount, heal amount, etc.)."],
            ["duration", "How long the effect lasts in seconds. Use 0 for instant or permanent effects."],
            ["scaling_stat", "Which stat scales this effect's power (e.g., willpower for a spell that gets stronger with higher willpower)."],
            ["scaling_ratio", "How much the scaling stat contributes. Acts as a multiplier."],
            ["target", "Who receives the effect: self, enemy, ally, or room."],
          ]}
        />
      </Section>

      <Section title="Effect Types">
        <p className="text-text-muted mb-3">
          Each effect falls into one of these categories. Some categories have subtypes for more
          specific behavior:
        </p>
        <Table
          headers={["Type", "What it does", "Subtypes"]}
          rows={[
            ["damage", "Deals damage to the target.", "physical, fire, cold, lightning, acid, mental"],
            ["heal", "Restores health to the target.", "--"],
            ["buff", "Gives a positive stat modifier.", "accuracy_boost, dodge_all, damage_reduction"],
            ["debuff", "Applies a negative stat modifier.", "--"],
            ["dot", "Damage over time. Ticks each combat round.", "--"],
            ["hot", "Healing over time. Ticks each combat round.", "--"],
            ["stun", "Prevents the target from taking any action for the duration.", "--"],
            ["set_bind_point", "Sets the character's respawn point to their current room.", "--"],
          ]}
        />
      </Section>

      <Section title="EffectHook Entity">
        <p className="text-text-muted mb-3">
          EffectHooks are how you wire an Effect into the game. They connect a reusable Effect
          template to specific events:
        </p>
        <Table
          headers={["Field", "What it means"]}
          rows={[
            ["effect_id", "The ID of the Effect you want to apply."],
            ["trigger_event", "When the effect fires. Options: ability_cast, hit, kill, quest_complete, dialog_node_enter."],
            ["target_type", "Who receives the effect: self, enemy, ally, or room."],
            ["target_id", "An optional specific target (like a particular ability ID)."],
            ["conditions", "JSON conditions that must be met before the hook fires."],
            ["priority", "Controls the order when multiple hooks fire on the same event. Lower numbers go first."],
          ]}
        />
      </Section>

      <Section title="ActiveEffect">
        <p className="text-text-muted mb-3">
          When an effect actually lands on a character, it becomes an ActiveEffect. This is the
          live, running instance:
        </p>
        <Table
          headers={["Field", "What it means"]}
          rows={[
            ["character_id", "Which character has the effect right now."],
            ["effect_id", "The Effect template this instance came from."],
            ["source", "Where the effect came from (an ability ID, quest ID, etc.)."],
            ["expires_at", "When the effect wears off. Only set for timed effects."],
            ["stacks", "How many times the effect has stacked. Stacking refreshes duration rather than multiplying the value."],
          ]}
        />
        <p className="text-text-muted mt-3">
          Timed effects are automatically cleaned up when they expire. If an effect stacks,
          the duration refreshes but the value stays the same. This prevents a buff from
          becoming absurdly powerful just because someone applied it multiple times.
        </p>
      </Section>

      <Section title="Legacy Fields">
        <p className="text-text-muted mb-3">
          The Ability entity used to have flat fields like damage, damage_type, healing, and
          apply_effect directly on it. Those are now deprecated. The new way is to create
          Effect entities and link them via the effects edge. The old fields still work for
          backward compatibility, but you should use Effect + EffectHook for any new content
          you create.
        </p>
      </Section>

      <Section title="Scaling Formula">
        <p className="text-text-muted">
          The final value of a scaled effect is: base_value + (scaling_stat_value * scaling_ratio).
          For example, if an effect has a base of 50, the character's willpower is 20, and the
          scaling ratio is 0.5, you get: 50 + (20 * 0.5) = 60 total.
        </p>
      </Section>
    </div>
  );
}