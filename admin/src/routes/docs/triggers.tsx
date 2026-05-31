import { createFileRoute } from "@tanstack/react-router";
import { PageContainer } from "../../components/PageContainer";

export const Route = createFileRoute("/docs/triggers")({
  component: TriggersDocs,
});

function TriggersDocs() {
  return (
    <PageContainer>
      <div className="flex items-center gap-2 mb-4">
        <Link
          to="/docs/"
          className="no-underline px-2.5 py-1.5 rounded border border-border hover:border-primary transition-colors text-sm font-medium"
        >
          &larr; Documentation
        </Link>
        <h1 className="text-lg sm:text-xl font-bold text-text">Triggers</h1>
      </div>

      <div className="prose prose-invert max-w-none">
        <div className="prose prose-invert max-w-none">
          <p className="text-text text-lg mb-6">
            Triggers are the glue between your game world and the things that happen in it.
            When a player interacts with a room, item, or object, a trigger fires and makes
            something happen. Think of them as "when X happens, do Y" rules.
          </p>

          <h2 className="text-text font-semibold text-xl mt-8 mb-4">Overview</h2>
          <p className="text-text mb-4">
            Every trigger has five pieces:
          </p>
          <ul className="text-text list-disc pl-5 space-y-2 mb-6">
            <li><strong>Trigger Type:</strong> What player action fires the trigger (use, touch, press, enter, examine)</li>
            <li><strong>Target Type:</strong> What happens when it fires (recipe, effect, dialog node)</li>
            <li><strong>Target ID:</strong> Which specific thing to execute</li>
            <li><strong>Target Object:</strong> The room or equipment the trigger is attached to</li>
            <li><strong>Condition:</strong> An optional SPICE expression that must be true before the trigger fires</li>
          </ul>

          <h2 className="text-text font-semibold text-xl mt-8 mb-4">Trigger Types</h2>

          <h3 className="text-text font-bold text-lg mt-4 mb-2">use</h3>
          <p className="text-text mb-4">Fires when a player uses an item. Good for potions, spell scrolls, or anything the player activates.</p>

          <h3 className="text-text font-bold text-lg mt-4 mb-2">touch</h3>
          <p className="text-text mb-4">Fires when a player touches something. Think statues, levers, mystical orbs, or anything tactile.</p>

          <h3 className="text-text font-bold text-lg mt-4 mb-2">press</h3>
          <p className="text-text mb-4">Fires when a player presses a button or switch.</p>

          <h3 className="text-text font-bold text-lg mt-4 mb-2">enter</h3>
          <p className="text-text mb-4">Fires when a player walks into a room. Room triggers only.</p>

          <h3 className="text-text font-bold text-lg mt-4 mb-2">examine</h3>
          <p className="text-text mb-4">Fires when a player examines an object closely.</p>

          <h2 className="text-text font-semibold text-xl mt-8 mb-4">Target Types</h2>

          <h3 className="text-text font-bold text-lg mt-4 mb-2">recipe</h3>
          <p className="text-text mb-4">
            Runs a crafting recipe when triggered. The player receives whatever the recipe produces.
            This is handy for setups where using a tool or workstation automatically crafts something.
          </p>

          <h3 className="text-text font-bold text-lg mt-4 mb-2">effect</h3>
          <p className="text-text mb-4">
            Applies an active or passive effect to the player or target. Effects can:
          </p>
          <ul className="text-text list-disc pl-5 mb-4">
            <li>Heal or damage the target</li>
            <li>Buff or debuff stats</li>
            <li>Apply status effects like stun or invisibility</li>
            <li>Grant temporary abilities</li>
          </ul>

          <h3 className="text-text font-bold text-lg mt-4 mb-2">dialog_node</h3>
          <p className="text-text mb-4">
            Opens a specific NPC dialog node, starting or advancing a conversation. Great for:
          </p>
          <ul className="text-text list-disc pl-5 mb-4">
            <li>Kicking off quests when players interact with objects</li>
            <li>Triggering story moments when players find special items</li>
            <li>Revealing secrets through dialogue</li>
          </ul>

          <h2 className="text-text font-semibold text-xl mt-8 mb-4">Conditions (SPICE)</h2>

          <p className="text-text mb-4">
            SPICE expressions let you gate a trigger behind a condition. Before the trigger fires,
            the expression is evaluated. If it returns <code>false</code>, nothing happens. Leave
            the condition blank if you want the trigger to always fire.
          </p>

          <h3 className="text-text font-bold text-lg mt-4 mb-2">Variables You Can Use</h3>
          <ul className="text-text list-disc pl-5 mb-4">
            <li><code>player_level</code> - The character's current level</li>
            <li><code>player_class</code> - The character's class name</li>
            <li><code>player_race</code> - The character's race</li>
            <li><code>room_id</code> - The ID of the current room</li>
            <li><code>has_tag(&quot;tag_name&quot;)</code> - Check whether a character has a specific tag</li>
            <li><code>has_trait(&quot;trait_name&quot;)</code> - Check whether a character has a specific trait</li>
            <li><code>stat(&quot;stat_name&quot;)</code> - Read one of the character's stats (str, dex, etc.)</li>
          </ul>

          <h3 className="text-text font-bold text-lg mt-4 mb-2">Examples</h3>

          <div className="bg-surface border border-border rounded p-4 mb-4">
            <h4 className="text-text font-semibold mb-2">Only allow characters level 10 and above</h4>
            <code className="block text-sm text-primary">player_level &gt;= 10</code>
          </div>

          <div className="bg-surface border border-border rounded p-4 mb-4">
            <h4 className="text-text font-semibold mb-2">Only mages can use this</h4>
            <code className="block text-sm text-primary">player_class === &quot;mage&quot;</code>
          </div>

          <div className="bg-surface border border-border rounded p-4 mb-4">
            <h4 className="text-text font-semibold mb-2">Requires a specific tag</h4>
            <code className="block text-sm text-primary">has_tag(&quot;guard_key&quot;)</code>
          </div>

          <div className="bg-surface border border-border rounded p-4 mb-4">
            <h4 className="text-text font-semibold mb-2">Combining conditions</h4>
            <code className="block text-sm text-primary">player_level &gt;= 10 &amp;&amp; has_tag(&quot;quest_complete&quot;)</code>
          </div>

          <h2 className="text-text font-semibold text-xl mt-8 mb-4">Common Use Cases</h2>

          <h3 className="text-text font-bold text-lg mt-4 mb-2">Door Traps</h3>
          <p className="text-text mb-4">
            <ol className="text-text list-decimal pl-5">
              <li>Create a room trigger with type <code>enter</code></li>
              <li>Set target type to <code>effect</code></li>
              <li>Pick a damage or stun effect</li>
              <li>Optionally set a condition like <code>player_level &lt; 15</code> so only low-level players get hit</li>
            </ol>
          </p>

          <h3 className="text-text font-bold text-lg mt-4 mb-2">Auto-Crafting Station</h3>
          <p className="text-text mb-4">
            <ol className="text-text list-decimal pl-5">
              <li>Create an equipment trigger with type <code>use</code></li>
              <li>Set target type to <code>recipe</code></li>
              <li>Link it to a crafting recipe</li>
              <li>When a player uses the station, the recipe is crafted automatically</li>
            </ol>
          </p>

          <h3 className="text-text font-bold text-lg mt-4 mb-2">NPC Quest Trigger</h3>
          <p className="text-text mb-4">
            <ol className="text-text list-decimal pl-5">
              <li>Create an equipment trigger with type <code>touch</code></li>
              <li>Set target type to <code>dialog_node</code></li>
              <li>Link it to a dialog node that offers a quest</li>
              <li>When the player touches the quest item, the conversation starts</li>
            </ol>
          </p>

          <h2 className="text-text font-semibold text-xl mt-8 mb-4">Debugging Triggers</h2>

          <p className="text-text mb-4">
            If a trigger isn't working, check these things:
          </p>
          <ul className="text-text list-disc pl-5 mb-4">
            <li>Make sure the trigger is <strong>Enabled</strong></li>
            <li>Verify the Target ID actually exists. The search dropdown will validate this for you.</li>
            <li>If you set a condition, make sure it evaluates to true for your test character</li>
            <li>Confirm you're in the correct World. Triggers are world-specific.</li>
          </ul>

          <div className="bg-primary/10 border border-primary/20 rounded p-4 mt-6">
            <h3 className="text-primary font-semibold mb-2">Tip: Trigger Ordering</h3>
            <p className="text-text text-sm">
              When you have multiple triggers on the same object, they fire in order of their ID.
              Lower IDs go first. Use this to create primary and secondary effects that fire in sequence.
            </p>
          </div>

          <Link
            to="/triggers"
            className="inline-block mt-8 px-4 py-2 bg-primary hover:bg-primary/80 text-white rounded font-medium transition-colors"
          >
            Back to Triggers
          </Link>
        </div>
      </div>
    </PageContainer>
  );
}