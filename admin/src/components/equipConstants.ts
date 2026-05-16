/** Shared equipment constants for the admin UI. */

/** Master catalog of all valid equipment slot names. */
export const SLOT_CATALOG = [
  "head", "neck", "chest", "back", "hands", "legs", "feet",
  "finger_left", "finger_right",
  "main_hand", "off_hand",
  "tail", "horn", "wings", "shell",
] as const;

/** Default 11 humanoid slots. */
export const DEFAULT_HUMANOID_SLOTS: string[] = [
  "head", "neck", "chest", "back", "hands", "legs", "feet",
  "finger_left", "finger_right", "main_hand", "off_hand",
];

/** Weapon type options. */
export const WEAPON_TYPES = [
  { value: "", label: "— none —" },
  { value: "sword", label: "Sword" },
  { value: "axe", label: "Axe" },
  { value: "spear", label: "Spear" },
  { value: "knife", label: "Knife" },
  { value: "martial", label: "Martial" },
  { value: "staff", label: "Staff" },
  { value: "pipe", label: "Pipe" },
] as const;

/** Damage type options. */
export const DAMAGE_TYPES = [
  { value: "", label: "— none —" },
  { value: "slashing", label: "Slashing" },
  { value: "piercing", label: "Piercing" },
  { value: "bludgeoning", label: "Bludgeoning" },
  { value: "fire", label: "Fire" },
] as const;

/** Armor type options. */
export const ARMOR_TYPES = [
  { value: "", label: "— none —" },
  { value: "light", label: "Light" },
  { value: "cloth", label: "Cloth" },
  { value: "heavy", label: "Heavy" },
] as const;

/** Rarity options. */
export const RARITY_OPTIONS = [
  { value: "common", label: "Common" },
  { value: "uncommon", label: "Uncommon" },
  { value: "rare", label: "Rare" },
  { value: "epic", label: "Epic" },
  { value: "legendary", label: "Legendary" },
] as const;

/** Slot names that are weapon-related. */
export const WEAPON_SLOTS = ["main_hand", "off_hand", "hands"];

/** Slot names that are armor-related. */
export const ARMOR_SLOTS = ["head", "neck", "chest", "back", "legs", "feet"];

/** Returns true if the slot is weapon-related. */
export const isWeaponSlot = (slot: string): boolean =>
  WEAPON_SLOTS.some((s) => slot.includes(s));

/** Returns true if the slot is armor-related. */
export const isArmorSlot = (slot: string): boolean =>
  ARMOR_SLOTS.some((s) => slot === s);