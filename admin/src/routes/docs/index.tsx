import { createFileRoute, Link } from "@tanstack/react-router";
import { PageHeader } from "../../components/PageHeader";
import { DocsIcon } from "../../components/icons/DocsIcon";

export const Route = createFileRoute("/docs/")({
  component: DocsIndex,
});

const DOC_PAGES = [
  {
    title: "Ability System",
    path: "/docs/ability-system",
    desc: "How combat abilities work, effect types, scaling, mana costs, and the classless skill system.",
  },
  {
    title: "Combat Guide",
    path: "/docs/combat-guide",
    desc: "How tick-based combat works, the damage formula, dodge and parry, and skill combos.",
  },
  {
    title: "Trainable Skills",
    path: "/docs/trainable-skills",
    desc: "How weapon and magic proficiencies connect to abilities and training mechanics.",
  },
  {
    title: "NPC System",
    path: "/docs/npc-system",
    desc: "How NPCs spawn, respawn, scale with level, and what race effects do.",
  },
  {
    title: "Item System",
    path: "/docs/item-system",
    desc: "Equipment slots, how damage and armor are calculated, item tags, and categories.",
  },
  {
    title: "Faction System",
    path: "/docs/faction-system",
    desc: "How standing mechanics work, faction categories, and how factions affect gameplay.",
  },
  {
    title: "Quest System",
    path: "/docs/quest-system",
    desc: "How to set up quest objectives, track player progress, configure repeat modes, and hand out rewards.",
  },
  {
    title: "Examine Skill",
    path: "/docs/examine-skill",
    desc: "The examine command, hidden details, skill levels, and difficulty class checks.",
  },
  {
    title: "Config Reference",
    path: "/docs/config-reference",
    desc: "What each config key does and how it shapes the game world.",
  },
  {
    title: "Bind Points & Root Room",
    path: "/docs/bind-points",
    desc: "How the root room works, bind points, respawn mechanics, and reconnect positioning.",
  },
  {
    title: "Worlds",
    path: "/docs/worlds",
    desc: "How multi-world architecture works, world settings, content structure, and game time.",
  },
  {
    title: "Achievements",
    path: "/docs/achievements",
    desc: "Achievement criteria, XP rewards, and how to manage achievements as an admin.",
  },
  {
    title: "Character Tags",
    path: "/docs/character-tags",
    desc: "Key-value metadata on characters for quest markers, feature flags, and more.",
  },
  {
    title: "Competencies",
    path: "/docs/competencies",
    desc: "Trainable skill categories, how level thresholds work, and XP progression.",
  },
  {
    title: "Crafting System",
    path: "/docs/crafting-system",
    desc: "How recipes work, crafting stations, inputs and outputs, and crafting commands.",
  },
  {
    title: "Dialog System",
    path: "/docs/dialog-system",
    desc: "How NPC dialog trees work, branching conversations, and on-enter effects.",
  },
  {
    title: "Social Commands",
    path: "/docs/social-commands",
    desc: "How say, shout, whisper, tell, and emote work, plus how to set up custom socials.",
  },
  {
    title: "Effect System",
    path: "/docs/effect-system",
    desc: "How effects work as reusable game logic, effect hooks, active effects, types, and scaling.",
  },
];

function DocsIndex() {
  return (
    <div className="management-page">
      <PageHeader title="Documentation" backTo="/dashboard" />

      <p className="text-text-muted mb-6 max-w-2xl">
        These pages explain how the systems behind the admin panel actually work in-game.
        You can hover over form fields in the admin pages for quick tooltips, or dig into
        these docs when you need the full picture.
      </p>

      <div className="grid grid-cols-[repeat(auto-fit,minmax(280px,1fr))] gap-4">
        {DOC_PAGES.map((page) => (
          <Link
            key={page.path}
            to={page.path}
            className="block bg-surface-muted rounded-lg p-5 no-underline text-text border border-border transition-colors hover:border-primary hover:bg-surface-muted/70"
          >
            <div className="flex items-center gap-2 mb-2">
              <DocsIcon stroke="var(--color-primary)" />
              <span className="font-semibold">{page.title}</span>
            </div>
            <div className="text-text-muted text-sm">{page.desc}</div>
          </Link>
        ))}
      </div>
    </div>
  );
}