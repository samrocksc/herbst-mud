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
        <strong>TL;DR:</strong> Crafting lets players create equipment by providing inputs
        and optionally being at the right station. Recipes define what inputs are needed,
        what station is required, and what outputs are produced.
      </InfoBox>

      <Section title="Recipe Entity">
        <p className="text-text-muted mb-3">
          Recipes are defined as CraftingRecipe entities and managed via the admin UI at /recipes.
        </p>
        <Table
          headers={["Field", "Description"]}
          rows={[
            ["name", "Unique machine-readable identifier (e.g., iron_sword)"],
            ["display_name", "Player-facing name (e.g., Iron Sword)"],
            ["description", "What the recipe produces and its flavor text."],
            ["required_station_tag", "Room tag needed to craft (e.g., forge). Leave blank for anywhere."],
            ["required_class", "Class restriction. Leave blank for any class."],
            ["required_skill", "Skill name required (e.g., Blacksmithing). Leave blank for none."],
            ["required_skill_level", "Minimum skill level required."],
            ["inputs", "JSON array of CraftingInput objects."],
            ["outputs", "JSON array of CraftingOutput objects."],
            ["craft_time_secs", "Time in seconds the craft takes (future use)."],
          ]}
        />
      </Section>

      <Section title="Inputs and Outputs">
        <p className="text-text-muted mb-3">
          Each input specifies an equipment template and quantity. Inputs are consumed on successful craft.
        </p>
        <Table
          headers={["Field", "Description"]}
          rows={[
            ["equipment_template_id", "ID of the equipment template to consume."],
            ["quantity", "Number of units consumed per craft."],
            ["consumed", "Whether the input is consumed (true) or returned (false)."],
          ]}
        />
        <p className="text-text-muted mb-3 mt-4">
          Outputs define what is created:
        </p>
        <Table
          headers={["Field", "Description"]}
          rows={[
            ["equipment_template_id", "ID of the equipment template produced."],
            ["quantity", "Number of units produced per craft."],
          ]}
        />
      </Section>

      <Section title="Stations">
        <p className="text-text-muted mb-3">
          Rooms act as crafting stations via tags. A room with tag <code>forge</code> satisfies
          <code>required_station_tag: forge</code>. Players can tag rooms in the admin panel.
        </p>
        <Table
          headers={["Example Tag", "Associated Crafting"]}
          rows={[
            ["forge", "Metal weapons and armor"],
            ["alchemy_lab", "Potions and alchemical items"],
            ["enchanting_table", "Magical item enchanting"],
            ["workbench", "General crafting"],
          ]}
        />
      </Section>

      <Section title="Player Commands">
        <Table
          headers={["Command", "Description"]}
          rows={[
            ["craft &lt;recipe&gt;", "Attempt to craft a recipe by name. Checks station, class, skill, and inputs."],
            ["recipes", "List all known recipes and whether requirements are met."],
            ["stations", "List crafting stations available in the current world."],
          ]}
        />
      </Section>

      <Section title="Admin UI">
        <p className="text-text-muted">
          The admin UI at /recipes provides full CRUD for crafting recipes.
          Inputs and outputs are edited as JSON arrays in the form.
        </p>
      </Section>
    </div>
  );
}
