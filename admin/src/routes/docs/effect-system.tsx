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
        <strong>TL;DR:</strong> Effects are reusable, composable units of game logic. Effect defines
        the template (name, type, value, duration). EffectHook attaches effects to abilities, quests,
        or events. ActiveEffect tracks runtime instances on characters.
      </InfoBox>

      <Section title="Effect Entity">
        <Table
          headers={["Field", "Description"]}
          rows={[
            ["name", "Unique identifier for the effect (e.g., burning, regenerate)."],
            ["description", "Player-facing description of the effect."],
            ["effect_type", "Category: damage, heal, buff, debuff, dot, hot, stun, accuracy_boost, dodge_all, set_bind_point."],
            ["value", "Base numeric value for the effect."],
            ["duration", "Duration in seconds. 0 = instant/permanent."],
            ["scaling_stat", "Stat that scales the effect value (e.g., willpower)."],
            ["scaling_ratio", "Ratio at which the scaling stat affects value."],
            ["target", "self, enemy, ally, or room."],
          ]}
        />
      </Section>

      <Section title="Effect Types">
        <Table
          headers={["Type", "Description", "Subtypes"]}
          rows={[
            ["damage", "Deals damage to target.", "physical, fire, cold, lightning, acid, mental"],
            ["heal", "Restores health.", "—"],
            ["buff", "Positive stat modifier.", "accuracy_boost, dodge_all, damage_reduction"],
            ["debuff", "Negative stat modifier.", "—"],
            ["dot", "Damage over time. Applied each tick.", "—"],
            ["hot", "Healing over time. Applied each tick.", "—"],
            ["stun", "Prevents action for duration.", "—"],
            ["set_bind_point", "Sets character's respawn bind.", "—"],
          ]}
        />
      </Section>

      <Section title="EffectHook Entity">
        <p className="text-text-muted mb-3">
          EffectHooks connect reusable Effect templates to triggers:
        </p>
        <Table
          headers={["Field", "Description"]}
          rows={[
            ["effect_id", "ID of the Effect entity to apply."],
            ["trigger_event", "When the effect fires: ability_cast, hit, kill, quest_complete, dialog_node_enter."],
            ["target_type", "Who receives: self, enemy, ally, room."],
            ["target_id", "Optional specific target (e.g., ability ID)."],
            ["conditions", "JSON conditions that must be met for the hook to fire."],
            ["priority", "Order when multiple hooks fire on the same event."],
          ]}
        />
      </Section>

      <Section title="ActiveEffect">
        <p className="text-text-muted mb-3">
          ActiveEffect tracks runtime effect instances on characters:
        </p>
        <Table
          headers={["Field", "Description"]}
          rows={[
            ["character_id", "Which character has the effect."],
            ["effect_id", "The Effect entity applied."],
            ["source", "Where the effect came from (ability ID, quest ID, etc.)."],
            ["expires_at", "When the effect ends (for timed effects)."],
            ["stacks", "Number of stacks (for stackable effects)."],
          ]}
        />
        <p className="text-text-muted mt-3">
          Expiry scheduling removes effects after their duration. Stacking effects
          refresh duration rather than multiplying value.
        </p>
      </Section>

      <Section title="Legacy Fields">
        <p className="text-text-muted mb-3">
          The Ability entity previously had flat fields for damage, damage_type, healing,
          and apply_effect. These are deprecated in favor of the effects edge.
          Old fields are still read for backward compatibility but new content should
          use Effect + EffectHook.
        </p>
      </Section>

      <Section title="Scaling Formula">
        <p className="text-text-muted">
          Final value = base_value + (scaling_stat_value * scaling_ratio). For a effect
          with base 50, scaling_stat willpower 20, and ratio 0.5: 50 + (20 * 0.5) = 60.
        </p>
      </Section>
    </div>
  );
}
