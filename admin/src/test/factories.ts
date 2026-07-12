import type { Ability, AbilityInput } from "../hooks/useAbilities";
import type { EquipmentTemplate, EquipmentTemplateInput } from "../hooks/useEquipmentTemplates";
import type { Quest, QuestInput } from "../hooks/useQuests";
import type { World, WorldInput } from "../hooks/useWorlds";
import type { Character, CharacterUpdate } from "../hooks/useCharacters";
import type { Tag, TagInput } from "../hooks/useTags";
import type { Race, RaceInput } from "../hooks/useRaces";

export function makeAbility(overrides?: Partial<Ability>): Ability {
  return {
    id: 1,
    name: "Fireball",
    description: "A ball of fire",
    ability_type: "magic",
    cost: 10,
    cooldown: 0,
    cooldown_seconds: 5,
    requirements: "{}",
    mana_cost: 10,
    stamina_cost: 0,
    hp_cost: 0,
    slug: "fireball",
    required_tag: "",
    ability_class: "active",
    proc_chance: 0,
    proc_event: "",
    faction_skills: null,
    ...overrides,
  };
}

export function makeAbilityInput(overrides?: Partial<AbilityInput>): AbilityInput {
  return {
    name: "Fireball",
    description: "A ball of fire",
    ability_type: "magic",
    requirements: "{}",
    cost: 10,
    cooldown: 0,
    cooldown_seconds: 5,
    mana_cost: 10,
    stamina_cost: 0,
    hp_cost: 0,
    proc_chance: 0,
    proc_event: "",
    ability_class: "active",
    required_tag: "",
    ...overrides,
  };
}

export function makeEquipmentTemplate(overrides?: Partial<EquipmentTemplate>): EquipmentTemplate {
  return {
    id: 1,
    slug: "iron-sword",
    name: "Iron Sword",
    description: "A basic iron sword",
    slot: "hand",
    level: 1,
    weight: 2,
    item_type: "weapon",
    stats: {},
    color: "#888888",
    is_visible: true,
    is_immovable: false,
    effect_type: "",
    effect_value: 0,
    effect_duration: 0,
    is_container: false,
    container_capacity: 0,
    is_locked: false,
    key_item_id: "",
    reveal_condition: "",
    armor_rating: 0,
    armor_type: "",
    rarity: "common",
    skill_requirement: "",
    skill_requirement_level: 0,
    damage_dice_count: 1,
    damage_dice_sides: 6,
    damage_bonus: 0,
    damage_type: "slashing",
    weapon_type: "sword",
    is_two_handed: false,
    ...overrides,
  };
}

export function makeEquipmentTemplateInput(overrides?: Partial<EquipmentTemplateInput>): EquipmentTemplateInput {
  return {
    name: "Iron Sword",
    slug: "iron-sword",
    description: "A basic iron sword",
    slot: "hand",
    level: 1,
    weight: 2,
    item_type: "weapon",
    ...overrides,
  };
}

export function makeNPCTemplate(overrides?: Partial<Record<string, unknown>>): Record<string, unknown> {
  return {
    id: 1,
    name: "Village Guard",
    description: "A generic guard",
    class: "warrior",
    race: "human",
    level: 1,
    ...overrides,
  };
}

export function makeNPCInstance(overrides?: Partial<Record<string, unknown>>): Record<string, unknown> {
  return {
    id: 1,
    name: "Village Guard",
    template_id: 1,
    room_id: 5,
    currentRoomId: 5,
    ...overrides,
  };
}

export function makeQuest(overrides?: Partial<Quest>): Quest {
  return {
    id: 1,
    name: "Find the Key",
    description: "Retrieve the ancient key",
    prerequisite_quest_ids: [],
    objectives: [{ type: "collect", target_id: "1", tag_filter: "", count: 1, labels: [], hint: "" }],
    rewards: { xp: 100, item_ids: [], effect_ids: [], tag_adds: [], tag_removes: [], achievement_ids: [] },
    repeat_mode: "once",
    cooldown_hours: 0,
    is_active: true,
    main_type: "main",
    ...overrides,
  };
}

export function makeQuestInput(overrides?: Partial<QuestInput>): QuestInput {
  return {
    name: "Find the Key",
    description: "Retrieve the ancient key",
    main_type: "main",
    prerequisite_quest_ids: [],
    objectives: [],
    rewards: { xp: 0, item_ids: [], effect_ids: [], tag_adds: [], tag_removes: [], achievement_ids: [] },
    repeat_mode: "once",
    cooldown_hours: 0,
    is_active: true,
    ...overrides,
  };
}

export function makeWorld(overrides?: Partial<World>): World {
  return {
    id: 1,
    name: "Test World",
    title: "The Test World",
    description: "A world for testing",
    active: true,
    ...overrides,
  };
}

export function makeWorldInput(overrides?: Partial<WorldInput>): WorldInput {
  return {
    name: "Test World",
    title: "The Test World",
    description: "A world for testing",
    active: true,
    ...overrides,
  };
}

export function makeCharacter(overrides?: Partial<Character>): Character {
  return {
    id: 1,
    name: "Testo",
    isNPC: false,
    currentRoomId: 5,
    startingRoomId: 5,
    respawnRoomId: 5,
    is_admin: false,
    is_immortal: false,
    is_test: false,
    currentWorld: "1",
    hitpoints: 100,
    max_hitpoints: 100,
    stamina: 50,
    max_stamina: 50,
    mana: 30,
    max_mana: 30,
    race: "human",
    class: "warrior",
    gender: "male",
    description: "A test character",
    level: 1,
    xp: 0,
    strength: 10,
    dexterity: 10,
    constitution: 10,
    intelligence: 10,
    wisdom: 10,
    lastSeenAt: null,
    ...overrides,
  };
}

export function makeTag(overrides?: Partial<Tag>): Tag {
  return {
    id: 1,
    name: "common",
    color: "#888888",
    ...overrides,
  };
}

export function makeTagInput(overrides?: Partial<TagInput>): TagInput {
  return {
    name: "common",
    color: "#888888",
    ...overrides,
  };
}

export function makeRace(overrides?: Partial<Race>): Race {
  return {
    id: 1,
    name: "human",
    display_name: "Human",
    description: "A standard human",
    stat_modifiers: null,
    skill_grants: [],
    ability_modifiers: [],
    equipment_slots: [],
    requirement_tags: [],
    color: "#ffccaa",
    tags: [],
    ...overrides,
  };
}

export function makeRaceInput(overrides?: Partial<RaceInput>): RaceInput {
  return {
    name: "human",
    display_name: "Human",
    description: "A standard human",
    stat_modifiers: "",
    skill_grants: [],
    equipment_slots: [],
    requirement_tags: [],
    color: "#ffccaa",
    tags: [],
    ...overrides,
  };
}

export function makeRoom(overrides?: Partial<Record<string, unknown>>): Record<string, unknown> {
  return {
    id: 1,
    name: "Town Square",
    description: "The center of town",
    isStartingRoom: true,
    isRootRoom: true,
    exits: {},
    posZ: 0,
    version: 1,
    ...overrides,
  };
}
