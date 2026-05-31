import { createFileRoute } from "@tanstack/react-router";
import { PageHeader } from "../../components/PageHeader";

export const Route = createFileRoute("/docs/crafting-system")({
  component: CraftingSystemDoc,
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

function CraftingSystemDoc() {
  return (
    <div className="management-page max-w-4xl">
      <PageHeader title="Crafting System" backTo="/docs" />

      <InfoBox>
        <strong>TL;DR:</strong> Crafting lets your players turn materials into gear. You define
        recipes that say what items go in, what comes out, and whether the player needs to be
        at a specific station (like a forge) to make it happen.
      </InfoBox>

      <Section title="Recipe Entity">
        <p className="text-text-muted mb-3">
          Each crafting recipe is a CraftingRecipe entity. You can create and edit recipes in the admin
          panel at /recipes. Here is what each field does:
        </p>
        <Table
          headers={["Field", "What it does"]}
          rows={[
            ["name", "A unique internal name for this recipe (like iron_sword). Players will not see this."],
            ["display_name", "The name players see (like Iron Sword)."],
            ["description", "Flavor text that tells players what they are making."],
            ["required_station_tag", "The room tag a player must be near to craft this (like forge). Leave blank if they can craft it anywhere."],
            ["required_class", "Restrict this recipe to a specific class. Leave blank if any class can craft it."],
            ["required_skill", "A skill the player must have (like Blacksmithing). Leave blank for no skill requirement."],
            ["required_skill_level", "The minimum level in that skill. Only matters if you set required_skill."],
            ["inputs", "A JSON array of CraftingInput objects. These are the materials the player provides."],
            ["outputs", "A JSON array of CraftingOutput objects. These are the items the player receives."],
            ["craft_time_secs", "How long the craft takes in seconds. This is reserved for future use."],
          ]}
        />
      </Section>

      <Section title="Inputs and Outputs">
        <p className="text-text-muted mb-3">
          Inputs are the materials the player brings to a recipe. When the craft succeeds, the inputs
          are removed from the player's inventory (unless you set consumed to false).
        </p>
        <Table
          headers={["Field", "What it does"]}
          rows={[
            ["equipment_template_id", "The ID of the equipment template the player provides."],
            ["quantity", "How many of that item the player needs for one craft."],
            ["consumed", "Whether the input gets used up. Set true (default) to consume it, false to return it after crafting."],
          ]}
        />
        <p className="text-text-muted mb-3 mt-4">
          Outputs are what the player gets when the craft finishes:
        </p>
        <Table
          headers={["Field", "What it does"]}
          rows={[
            ["equipment_template_id", "The ID of the equipment template the player receives."],
            ["quantity", "How many of that item the player gets per craft."],
          ]}
        />
      </Section>

      <Section title="Stations">
        <p className="text-text-muted mb-3">
          Some recipes require the player to be in the right place. You do this by tagging rooms.
          If a recipe has <code>required_station_tag: forge</code>, the player can only craft it in a
          room tagged <code>forge</code>. You can add tags to rooms in the admin panel.
        </p>
        <Table
          headers={["Example Tag", "What kind of crafting it unlocks"]}
          rows={[
            ["forge", "Metal weapons and armor"],
            ["alchemy_lab", "Potions and alchemical items"],
            ["enchanting_table", "Magical item enchanting"],
            ["workbench", "General crafting"],
          ]}
        />
      </Section>

      <Section title="Player Commands">
        <p className="text-text-muted mb-3">
          Players interact with crafting through these in-game commands:
        </p>
        <Table
          headers={["Command", "What it does"]}
          rows={[
            ["craft <recipe>", "Try to craft something by recipe name. The game checks the station, class, skill, and inputs before allowing it."],
            ["recipes", "See all the recipes you know and whether you meet the requirements."],
            ["stations", "See which crafting stations exist in the current world."],
          ]}
        />
      </Section>

      <Section title="Admin UI">
        <p className="text-text-muted">
          You can create, edit, and delete recipes at /recipes. Inputs and outputs are edited as
          JSON arrays in the form. If you are not sure what to put there, start simple: one input,
          one output, and build from there.
        </p>
      </Section>
    </div>
  );
}