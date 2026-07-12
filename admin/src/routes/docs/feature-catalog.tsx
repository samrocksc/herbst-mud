import { createFileRoute } from "@tanstack/react-router";
import { PageHeader } from "../../components/PageHeader";

export const Route = createFileRoute("/docs/feature-catalog")({
  component: FeatureCatalog,
});

type Feature = {
  id: string;
  category: string;
  name: string;
  status: "Complete" | "In Progress" | "Planned" | "Not Implemented";
  scope: string;
  description: string;
};

const STATUS_BADGES: Record<string, string> = {
  Complete: "bg-success/20 text-success",
  "In Progress": "bg-primary/20 text-primary",
  Planned: "bg-accent/20 text-accent",
  "Not Implemented": "bg-surface-muted text-text-muted",
};

const STATUS_ICONS: Record<string, string> = {
  Complete: "✅",
  "In Progress": "🛠️",
  Planned: "🔜",
  "Not Implemented": "❌",
};

const FEATURES: Feature[] = [
  // Authentication & Users
  { id: "auth-root", category: "Authentication", name: "Root User Authentication", status: "Complete", scope: "Web, API", description: "Admin authentication with root user" },
  { id: "auth-user-creation", category: "Authentication", name: "User Account Creation", status: "Complete", scope: "SSH, Web, API", description: "Create new user accounts via registration" },
  { id: "auth-user-login", category: "Authentication", name: "User Login", status: "Complete", scope: "SSH, Web, API", description: "Authenticate existing users" },
  { id: "auth-char-auth", category: "Authentication", name: "Character Authentication", status: "Complete", scope: "SSH, Web", description: "Authenticate characters within a world" },
  { id: "auth-password-obfuscation", category: "Authentication", name: "Password Obfuscation", status: "Complete", scope: "SSH", description: "Passwords show as asterisks during input" },
  { id: "auth-email-validation", category: "Authentication", name: "Email Validation", status: "Complete", scope: "SSH, API", description: "Validate email format during registration" },
  { id: "auth-password-confirm", category: "Authentication", name: "Password Confirmation", status: "Complete", scope: "SSH", description: "Confirm password during account creation" },

  // Character Creation & Management
  { id: "char-dynamic-creation", category: "Character Creation", name: "Dynamic Character Creation", status: "Complete", scope: "SSH", description: "Create characters with dynamic race/class selection" },
  { id: "char-class-selection", category: "Character Creation", name: "Character Class Selection", status: "Complete", scope: "Web, SSH", description: "Select character class during creation" },
  { id: "char-race-selection", category: "Character Creation", name: "Character Race Selection", status: "Complete", scope: "SSH", description: "Select character race from API-provided list" },
  { id: "char-name-validation", category: "Character Creation", name: "Character Name Validation", status: "Complete", scope: "SSH", description: "Validate name format and length" },
  { id: "char-list", category: "Character Creation", name: "Character List", status: "Complete", scope: "SSH", description: "View all owned characters" },
  { id: "char-selection", category: "Character Creation", name: "Character Selection", status: "Complete", scope: "SSH", description: "Select which character to play" },
  { id: "char-default-class", category: "Character Creation", name: "Default Class (Survivor)", status: "Complete", scope: "Web", description: "Auto-assign survivor class when none selected" },
  { id: "char-crud", category: "Character Creation", name: "Character CRUD API", status: "Complete", scope: "API", description: "Full CRUD operations for characters" },

  // World & MUD Selection
  { id: "world-selection", category: "World Selection", name: "MUD/World Selection", status: "Complete", scope: "SSH", description: "Select which world to play in" },
  { id: "world-multiserver", category: "World Selection", name: "Multiple World Support", status: "Complete", scope: "API, Admin", description: "Support multiple game worlds" },
  { id: "world-scoping", category: "World Selection", name: "World Scoping for Entities", status: "Complete", scope: "Admin", description: "NPCs, items, factions scoped to world" },
  { id: "world-export", category: "World Selection", name: "World Export", status: "Complete", scope: "Admin", description: "Export world data as JSON" },
  { id: "world-import", category: "World Selection", name: "World Import", status: "Complete", scope: "Admin", description: "Import world data from JSON" },

  // Combat System
  { id: "combat-initiation", category: "Combat", name: "Combat Initiation", status: "Complete", scope: "SSH, Web", description: "Start combat by targeting NPC" },
  { id: "combat-auto-attack", category: "Combat", name: "Auto-Attack", status: "Complete", scope: "SSH, Web", description: "Default attack on each combat tick" },
  { id: "combat-crits", category: "Combat", name: "Critical Hits (d20=20)", status: "Complete", scope: "SSH, Web", description: "Double damage on natural 20" },
  { id: "combat-fumbles", category: "Combat", name: "Fumbles (d20=1)", status: "Complete", scope: "SSH, Web", description: "Fail action on natural 1" },
  { id: "combat-flee", category: "Combat", name: "Flee Combat", status: "Complete", scope: "SSH, Web", description: "Attempt to escape combat" },
  { id: "combat-npc-death", category: "Combat", name: "NPC Death & Respawn", status: "Complete", scope: "SSH, Web", description: "NPCs respawn after death" },
  { id: "combat-screen", category: "Combat", name: "Combat Screen Redesign", status: "Complete", scope: "SSH", description: "Dedicated combat TUI screen" },
  { id: "combat-hotkeys", category: "Combat", name: "Combat Hotkeys (1-5)", status: "Complete", scope: "SSH, Web", description: "Use skills via hotkeys" },
  { id: "combat-potion-hotkey", category: "Combat", name: "Potion Hotkey (R)", status: "Complete", scope: "SSH, Web", description: "Use potions via hotkey" },
  { id: "combat-mana-deduction", category: "Combat", name: "Mana Deduction", status: "Complete", scope: "Web", description: "Deduct mana on spell cast" },

  // Abilities & Skills
  { id: "ability-crud", category: "Abilities & Skills", name: "Ability CRUD API", status: "Complete", scope: "Admin, API", description: "Full CRUD for abilities" },
  { id: "ability-effects", category: "Abilities & Skills", name: "Ability Effect Linking", status: "Complete", scope: "Admin", description: "Link effects to abilities" },
  { id: "ability-types", category: "Abilities & Skills", name: "Ability Types", status: "Complete", scope: "Admin", description: "Support active/passive abilities" },
  { id: "ability-cooldowns", category: "Abilities & Skills", name: "Ability Cooldowns", status: "Complete", scope: "SSH, Web", description: "Track and display cooldowns" },
  { id: "ability-costs", category: "Abilities & Skills", name: "Ability Mana/Stamina Costs", status: "Complete", scope: "SSH, Web", description: "Deduct resources on cast" },

  // Effects System
  { id: "effect-crud", category: "Effects System", name: "Effect CRUD API", status: "Complete", scope: "Admin, API", description: "Full CRUD for effects" },
  { id: "effect-types", category: "Effects System", name: "Effect Types", status: "Complete", scope: "Admin", description: "HP change, XP gain, message, teleport, etc." },
  { id: "effect-xp-gain", category: "Effects System", name: "XP Gain Effect", status: "Complete", scope: "Admin, API", description: "Grant XP to characters" },
  { id: "effect-xp-multiplier", category: "Effects System", name: "XP Multiplier Effect", status: "Complete", scope: "Admin, API", description: "Temporary XP boost" },
  { id: "effect-bind-point", category: "Effects System", name: "Bind Point Set", status: "Complete", scope: "Admin", description: "Set character bind point" },
  { id: "effect-teleport", category: "Effects System", name: "Teleport Effect", status: "Complete", scope: "Admin", description: "Teleport character to room" },
  { id: "effect-message", category: "Effects System", name: "Message Effects", status: "Complete", scope: "Admin", description: "Send custom messages" },
  { id: "effect-tag", category: "Effects System", name: "Tag Add/Remove", status: "Complete", scope: "Admin", description: "Add/remove character tags" },
  { id: "effect-nested", category: "Effects System", name: "Apply Effect (Nested)", status: "Complete", scope: "Admin", description: "Chain effects together" },
  { id: "effect-stack", category: "Effects System", name: "Stack Modes", status: "Complete", scope: "Admin", description: "Replace, refresh, additive stacking" },

  // NPC System
  { id: "npc-crud", category: "NPC System", name: "NPC Template CRUD", status: "Complete", scope: "Admin", description: "Create/manage NPC templates" },
  { id: "npc-instance", category: "NPC System", name: "NPC Instance Management", status: "Complete", scope: "Admin", description: "Spawn instances of templates" },
  { id: "npc-disposition", category: "NPC System", name: "NPC Disposition", status: "Complete", scope: "SSH, Web", description: "Hostile, friendly, shopkeeper" },
  { id: "npc-shopkeeper", category: "NPC System", name: "Shopkeeper NPCs", status: "Complete", scope: "SSH, Web", description: "Trade with shopkeeper NPCs" },
  { id: "npc-equipment", category: "NPC System", name: "NPC Equipment", status: "Complete", scope: "Admin", description: "Equip NPC instances with items" },
  { id: "npc-respawn", category: "NPC System", name: "NPC Respawn", status: "Complete", scope: "SSH", description: "NPCs respawn after death" },
  { id: "npc-greeting", category: "NPC System", name: "NPC Greeting", status: "Complete", scope: "Admin", description: "Set NPC greeting message" },

  // Item System
  { id: "item-crud", category: "Item System", name: "Item Template CRUD", status: "Complete", scope: "Admin", description: "Create/manage equipment templates" },
  { id: "item-instance", category: "Item System", name: "Item Instance Spawning", status: "Complete", scope: "Admin", description: "Spawn items into rooms" },
  { id: "item-types", category: "Item System", name: "Item Types", status: "Complete", scope: "Admin", description: "Weapons, armor, potions, accessories" },
  { id: "item-damage", category: "Item System", name: "Weapon Damage", status: "Complete", scope: "SSH, Web", description: "Dice-based damage calculation" },
  { id: "item-armor", category: "Item System", name: "Armor Rating", status: "Complete", scope: "SSH", description: "Armor reduces incoming damage" },
  { id: "item-slots", category: "Item System", name: "Equipment Slots", status: "Complete", scope: "SSH", description: "Main hand, off hand, chest, etc." },

  // Quest System
  { id: "quest-crud", category: "Quest System", name: "Quest CRUD API", status: "Complete", scope: "Admin, API", description: "Create/manage quests" },
  { id: "quest-progress", category: "Quest System", name: "Quest Progress Tracking", status: "Complete", scope: "SSH", description: "Track objective progress" },
  { id: "quest-objectives", category: "Quest System", name: "Quest Objectives", status: "Complete", scope: "Admin", description: "Kill, collect, explore objectives" },
  { id: "quest-rewards", category: "Quest System", name: "Quest Rewards", status: "Complete", scope: "SSH", description: "XP and item rewards on completion" },
  { id: "quest-status", category: "Quest System", name: "Quest Status", status: "Complete", scope: "SSH", description: "Active, completed, abandoned" },
  { id: "quest-log", category: "Quest System", name: "Quest Log Display", status: "Complete", scope: "SSH", description: "Show quest log TUI" },

  // Crafting System
  { id: "craft-recipes", category: "Crafting System", name: "Recipe Management", status: "Complete", scope: "Admin", description: "Create/edit recipes" },
  { id: "craft-stations", category: "Crafting System", name: "Crafting Stations", status: "Complete", scope: "SSH", description: "Multiple station types" },
  { id: "craft-command", category: "Crafting System", name: "Craft Command", status: "Complete", scope: "SSH", description: "Craft an item" },
  { id: "craft-list", category: "Crafting System", name: "Recipes List", status: "Complete", scope: "SSH", description: "List available recipes" },
  { id: "craft-output", category: "Crafting System", name: "Crafting Output", status: "Complete", scope: "SSH", description: "Create item instances" },

  // Map & Rooms
  { id: "room-crud", category: "Map & Rooms", name: "Room CRUD API", status: "Complete", scope: "Admin, API", description: "Full CRUD for rooms" },
  { id: "room-coords", category: "Map & Rooms", name: "Room Coordinates", status: "Complete", scope: "Admin", description: "X, Y, Z coordinates" },
  { id: "room-exits", category: "Map & Rooms", name: "Room Exits", status: "Complete", scope: "Admin", description: "Directional exits" },
  { id: "room-bidirectional", category: "Map & Rooms", name: "Bidirectional Exits", status: "Complete", scope: "Admin", description: "Create two-way connections" },
  { id: "room-drag-drop", category: "Map & Rooms", name: "Map Drag & Drop", status: "Complete", scope: "Admin", description: "Create rooms visually" },
  { id: "room-exit-drag", category: "Map & Rooms", name: "Map Exit by Dragging", status: "Complete", scope: "Admin", description: "Connect rooms via drag" },
  { id: "room-navigation", category: "Map & Rooms", name: "Room Navigation", status: "Complete", scope: "SSH", description: "Move between rooms" },
  { id: "room-cleanup", category: "Map & Rooms", name: "Room Cleanup", status: "Complete", scope: "Admin", description: "Remove orphan exits" },

  // Admin Panel Features
  { id: "admin-dashboard", category: "Admin Panel", name: "Admin Dashboard", status: "Complete", scope: "Web", description: "Overview with stats" },
  { id: "admin-npc-count", category: "Admin Panel", name: "Active NPC Count", status: "Complete", scope: "Web", description: "Display active NPC count" },
  { id: "admin-char-management", category: "Admin Panel", name: "Character Management", status: "Complete", scope: "Admin", description: "List/view characters" },
  { id: "admin-player-modal", category: "Admin Panel", name: "Player Detail Modal", status: "Complete", scope: "Admin", description: "View player details" },
  { id: "admin-factions", category: "Admin Panel", name: "Faction Management", status: "Complete", scope: "Admin", description: "Create/edit factions" },
  { id: "admin-faction-categories", category: "Admin Panel", name: "Faction Categories", status: "Complete", scope: "Admin", description: "Organize factions by category" },
  { id: "admin-achievements", category: "Admin Panel", name: "Achievement Management", status: "Complete", scope: "Admin", description: "CRUD for achievements" },
  { id: "admin-channels", category: "Admin Panel", name: "Channel Management", status: "Complete", scope: "Admin", description: "Configure chat channels" },
  { id: "admin-triggers", category: "Admin Panel", name: "Trigger Management", status: "Complete", scope: "Admin", description: "Event triggers" },
  { id: "admin-config", category: "Admin Panel", name: "Config Management", status: "Complete", scope: "Admin", description: "Global configuration" },
  { id: "admin-races", category: "Admin Panel", name: "Race Management", status: "Complete", scope: "Admin", description: "Create/edit races" },

  // SSH Client Commands
  { id: "ssh-attack", category: "SSH Commands", name: "Attack Commands", status: "Complete", scope: "SSH", description: "Attack, kill, fight commands" },
  { id: "ssh-examine", category: "SSH Commands", name: "Examine Command", status: "Complete", scope: "SSH", description: "Examine targets" },
  { id: "ssh-look", category: "SSH Commands", name: "Look Command", status: "Complete", scope: "SSH", description: "Look around room" },
  { id: "ssh-move", category: "SSH Commands", name: "Movement Commands", status: "Complete", scope: "SSH", description: "North/south/east/west/up/down" },
  { id: "ssh-inventory", category: "SSH Commands", name: "Inventory Command", status: "Complete", scope: "SSH", description: "View inventory" },
  { id: "ssh-equip", category: "SSH Commands", name: "Equip Command", status: "Complete", scope: "SSH", description: "Show and equip items" },
  { id: "ssh-craft", category: "SSH Commands", name: "Craft Command", status: "Complete", scope: "SSH", description: "Craft items" },
  { id: "ssh-quests", category: "SSH Commands", name: "Quests Command", status: "Complete", scope: "SSH", description: "View quest log" },
  { id: "ssh-chat", category: "SSH Commands", name: "Chat Commands", status: "Complete", scope: "SSH", description: "Chat and social commands" },
  { id: "ssh-reply", category: "SSH Commands", name: "Reply Command", status: "Complete", scope: "SSH", description: "Reply to last whisper" },
  { id: "ssh-ignore", category: "SSH Commands", name: "Ignore Command", status: "Complete", scope: "SSH", description: "Ignore a player" },
  { id: "ssh-logout", category: "SSH Commands", name: "Logout Command", status: "Complete", scope: "SSH", description: "Save and exit" },

  // Web Client Features
  { id: "web-login", category: "Web Client", name: "Login Screen", status: "Complete", scope: "Web", description: "Login with credentials" },
  { id: "web-char-select", category: "Web Client", name: "Character Selection", status: "Complete", scope: "Web", description: "Choose character to play" },
  { id: "web-game-screen", category: "Web Client", name: "Game Screen", status: "Complete", scope: "Web", description: "Main gameplay interface" },
  { id: "web-combat-hud", category: "Web Client", name: "Combat HUD", status: "Complete", scope: "Web", description: "5-slot combat interface" },
  { id: "web-combat-log", category: "Web Client", name: "Combat Log", status: "Complete", scope: "Web", description: "Display combat events" },
  { id: "web-char-stats", category: "Web Client", name: "Character Stats", status: "Complete", scope: "Web", description: "Display HP, Mana, XP" },
  { id: "web-room-text", category: "Web Client", name: "Room Text", status: "Complete", scope: "Web", description: "Display room description" },
  { id: "web-char-list", category: "Web Client", name: "Character List", status: "Complete", scope: "Web", description: "View characters in room" },
  { id: "web-actions", category: "Web Client", name: "Action Buttons", status: "Complete", scope: "Web", description: "Attack, Talk, Trade, Examine" },

  // API Endpoints
  { id: "api-auth", category: "API Endpoints", name: "Authentication API", status: "Complete", scope: "API", description: "User and character authentication" },
  { id: "api-characters", category: "API Endpoints", name: "Characters API", status: "Complete", scope: "API", description: "Character CRUD and craft operations" },
  { id: "api-abilities", category: "API Endpoints", name: "Abilities API", status: "Complete", scope: "API", description: "Ability and effect management" },
  { id: "api-npcs", category: "API Endpoints", name: "NPCs API", status: "Complete", scope: "API", description: "NPC templates and instances" },
  { id: "api-items", category: "API Endpoints", name: "Items API", status: "Complete", scope: "API", description: "Equipment templates and instances" },
  { id: "api-recipes", category: "API Endpoints", name: "Recipes API", status: "Complete", scope: "API", description: "Crafting recipes" },
  { id: "api-quests", category: "API Endpoints", name: "Quests API", status: "Complete", scope: "API", description: "Quest management" },
  { id: "api-rooms", category: "API Endpoints", name: "Rooms API", status: "Complete", scope: "API", description: "Room CRUD and exits" },
  { id: "api-races", category: "API Endpoints", name: "Races API", status: "Complete", scope: "API", description: "Race management" },
  { id: "api-factions", category: "API Endpoints", name: "Factions API", status: "Complete", scope: "API", description: "Faction and category management" },
  { id: "api-achievements", category: "API Endpoints", name: "Achievements API", status: "Complete", scope: "API", description: "Achievement CRUD" },
  { id: "api-channels", category: "API Endpoints", name: "Channels API", status: "Complete", scope: "API", description: "Chat channel management" },
  { id: "api-triggers", category: "API Endpoints", name: "Triggers API", status: "Complete", scope: "API", description: "Trigger management" },
  { id: "api-admin", category: "API Endpoints", name: "Admin API", status: "Complete", scope: "API", description: "Export, import, world operations" },
];

const CATEGORIES = [
  "Authentication", "Character Creation", "World Selection", "Combat",
  "Abilities & Skills", "Effects System", "NPC System", "Item System",
  "Quest System", "Crafting System", "Map & Rooms", "Admin Panel",
  "SSH Commands", "Web Client", "API Endpoints",
];

function FeatureCatalog() {
  return (
    <div className="management-page max-w-5xl">
      <PageHeader title="Feature Catalog" backTo="/docs" />

      <p className="text-text-muted mb-6">
        This catalog documents all features in herbst-mud. Features are organized by
        domain and labeled with their implementation status. For detailed technical
        documentation, see the individual system pages linked in the sidebar.
      </p>

      <div className="space-y-8">
        {CATEGORIES.map((category) => {
          const categoryFeatures = FEATURES.filter((f) => f.category === category);

          if (categoryFeatures.length === 0) return null;

          return (
            <div key={category}>
              <h2 className="text-lg font-semibold text-text mb-3 pb-2 border-b border-border">
                {category}
              </h2>
              <div className="overflow-x-auto">
                <table className="w-full text-sm border border-border rounded-lg">
                  <thead>
                    <tr className="bg-surface-muted">
                      <th className="text-left px-3 py-2 font-semibold">Status</th>
                      <th className="text-left px-3 py-2 font-semibold">Name</th>
                      <th className="text-left px-3 py-2 font-semibold">Scope</th>
                      <th className="text-left px-3 py-2 font-semibold">Description</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-border">
                    {categoryFeatures.map((feature) => (
                      <tr key={feature.id}>
                        <td className="px-3 py-2">
                          <span
                            className={`inline-flex items-center gap-1 px-2 py-1 rounded text-xs font-medium ${STATUS_BADGES[feature.status]}`}
                          >
                            {STATUS_ICONS[feature.status]}
                            {feature.status}
                          </span>
                        </td>
                        <td className="px-3 py-2 font-medium text-text">{feature.name}</td>
                        <td className="px-3 py-2 text-text-muted">{feature.scope}</td>
                        <td className="px-3 py-2 text-text-muted">{feature.description}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>
          );
        })}
      </div>

      <div className="mt-8 border-t border-border pt-6">
        <h2 className="text-lg font-semibold text-text mb-3">Not Yet Implemented</h2>
        <p className="text-text-muted mb-3">
          Features marked as "Not Implemented" or "Planned" are not yet complete:
        </p>
        <ul className="list-disc pl-6 text-text-muted space-y-1">
          <li>World Creation UI - Admin UI to create new worlds</li>
          <li>Character Deletion - UI to delete characters</li>
          <li>Trading System - Full merchant trading</li>
          <li>Guild/Faction Reputation - Reputation mechanics</li>
          <li>Mail System - In-game mail</li>
          <li>Auction House - Player trading platform</li>
          <li>Player Kills (PK) - Player-vs-player combat</li>
          <li>Rate Limiting - API rate limiting</li>
          <li>Webhook Events - External webhook integration</li>
        </ul>
        <p className="text-text-muted text-sm mt-4">
          For complete status tracking, see the main Feature Catalog document at
          <code className="ml-1 text-primary">docs/FEATURE_CATALOG.md</code>.
        </p>
      </div>
    </div>
  );
}
