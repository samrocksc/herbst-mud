export type OutputLine = {
  readonly text: string;
  readonly style: "default" | "error" | "success" | "combat" | "damage" | "heal" | "chat" | "system" | "prompt" | "room_description" | "event";
  readonly timestamp: number;
};

export type RoomExit = {
  readonly direction: string;
  readonly target: number;
  readonly label: string;
};

export type RoomCharacter = {
  readonly name: string;
  readonly type: "npc" | "player";
  readonly id: number;
  readonly hostile: boolean;
};

export type RoomItem = {
  readonly id: number;
  readonly name: string;
  readonly takeable: boolean;
  readonly examinable: boolean;
  readonly readable?: boolean;
  readonly lootable?: boolean;
};

export type RoomScreenPayload = {
  readonly view_type: "room";
  readonly id: number;
  readonly title: string;
  readonly description: string;
  readonly ascii_art_id?: string;
  readonly exits: readonly RoomExit[];
  readonly characters: readonly RoomCharacter[];
  readonly items: readonly RoomItem[];
};

export type Ability = {
  readonly id: number;
  readonly name: string;
  readonly description: string;
  readonly ability_type: string;
  readonly cooldown: number;
  readonly mana_cost: number;
  readonly stamina_cost: number;
  readonly hp_cost: number;
};

export type CharacterSkill = {
  readonly slot: number;
  readonly name: string | null;
};

export type InventoryItem = {
  readonly id: number;
  readonly name: string;
  readonly quantity?: number;
  readonly type?: string;
  readonly description?: string;
};

export type CharacterPanelTab = "inventory" | "skills" | "abilities";

export type ClientMessage = {
  readonly type: "command";
  readonly payload: { readonly input: string };
};

export type OutputMessage = {
  readonly type: "output";
  readonly payload: {
    readonly lines: readonly OutputLine[];
  };
};

export type ScreenMessage = {
  readonly type: "screen";
  readonly payload: RoomScreenPayload;
};

export type EventMessage = {
  readonly type: "event";
  readonly payload: Readonly<Record<string, unknown>>;
};

export type ErrorMessage = {
  readonly type: "error";
  readonly payload: { readonly message: string };
};

export type ServerMessage = OutputMessage | ScreenMessage | EventMessage | ErrorMessage;

export type CombatTarget = {
  readonly id: number;
  readonly name: string;
  readonly hp: number;
  readonly maxHp: number;
  readonly level?: number;
};

export type CombatLogEntry = {
  readonly timestamp: number;
  readonly text: string;
  readonly kind: "hit" | "miss" | "crit" | "heal" | "system" | "queue" | "flee" | "defeat";
};