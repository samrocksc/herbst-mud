import { DashboardIcon } from "./icons/DashboardIcon";
import { ContentIcon } from "./icons/ContentIcon";
import { NPCsIcon } from "./icons/NPCsIcon";
import { ItemsIcon } from "./icons/ItemsIcon";
import { AbilitiesIcon } from "./icons/AbilitiesIcon";
import { EffectsIcon } from "./icons/EffectsIcon";
import { SkillsIcon } from "./icons/SkillsIcon";
import { QuestsIcon } from "./icons/QuestsIcon";
import { CraftingIcon } from "./icons/CraftingIcon";
import { RacesIcon } from "./icons/RacesIcon";
import { FactionsIcon } from "./icons/FactionsIcon";
import { MapIcon } from "./icons/MapIcon";
import { PlayersIcon } from "./icons/PlayersIcon";
import { XpIcon } from "./icons/XpIcon";
import { SocialsIcon } from "./icons/SocialsIcon";
import { ChannelsIcon } from "./icons/ChannelsIcon";
import { ConfigIcon } from "./icons/ConfigIcon";
import { WorldIcon } from "./icons/WorldIcon";
import { TagsIcon } from "./icons/TagsIcon";
import { LogsIcon } from "./icons/LogsIcon";
import { DocsIcon } from "./icons/DocsIcon";
import type { ComponentType } from "react";

export interface NavItem {
  label: string;
  path: string;
  Icon: ComponentType<{ className?: string; stroke?: string }>;
}

export interface NavGroup {
  label: string;
  Icon: ComponentType<{ className?: string; stroke?: string }>;
  items: NavItem[];
}

export const navGroups: NavGroup[] = [
  {
    label: "Dashboard",
    Icon: DashboardIcon,
    items: [{ label: "Dashboard", path: "/dashboard", Icon: DashboardIcon }],
  },
  {
    label: "Content",
    Icon: ContentIcon,
    items: [
      { label: "NPCs", path: "/npcs", Icon: NPCsIcon },
      { label: "Items", path: "/items", Icon: ItemsIcon },
      { label: "Abilities", path: "/abilities", Icon: AbilitiesIcon },
      { label: "Effects", path: "/effects", Icon: EffectsIcon },
      { label: "Skills", path: "/skills", Icon: SkillsIcon },
      { label: "Quests", path: "/quests", Icon: QuestsIcon },
      { label: "Recipes", path: "/recipes", Icon: CraftingIcon },
      { label: "Races", path: "/races", Icon: RacesIcon },
      { label: "Factions", path: "/factions", Icon: FactionsIcon },
      { label: "Map", path: "/map", Icon: MapIcon },
    ],
  },
  {
    label: "Players",
    Icon: PlayersIcon,
    items: [
      { label: "Players", path: "/players", Icon: PlayersIcon },
      { label: "Characters", path: "/characters", Icon: PlayersIcon },
      { label: "XP", path: "/xp", Icon: XpIcon },
    ],
  },
  {
    label: "Social",
    Icon: SocialsIcon,
    items: [
      { label: "Socials", path: "/socials", Icon: SocialsIcon },
      { label: "Channels", path: "/channels", Icon: ChannelsIcon },
    ],
  },
  {
    label: "System",
    Icon: ConfigIcon,
    items: [
      { label: "Config", path: "/config", Icon: ConfigIcon },
      { label: "Worlds", path: "/worlds", Icon: WorldIcon },
      { label: "Tags", path: "/tags", Icon: TagsIcon },
      { label: "Logs", path: "/logs", Icon: LogsIcon },
      { label: "Docs", path: "/docs", Icon: DocsIcon },
    ],
  },
];

export function findActiveGroup(pathname: string): NavGroup {
  for (const group of navGroups) {
    for (const item of group.items) {
      if (pathname === item.path || pathname.startsWith(item.path + "/")) {
        return group;
      }
    }
  }
  return navGroups[0];
}
