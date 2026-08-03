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
  readonly npc_template_id?: string;
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
  readonly exits: readonly RoomExit[] | null;
  readonly characters: readonly RoomCharacter[] | null;
  readonly items: readonly RoomItem[] | null;
};

export type ConversationScreenPayload = {
  readonly view_type: "conversation";
  readonly npc_name: string;
  readonly npc_template_id: string;
  readonly nodes: Record<string, DialogNode>;
  readonly current_node_id: string;
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
  readonly caster_message?: string;
  readonly recipient_message?: string;
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

export type CharacterPanelTab = "inventory" | "skills" | "abilities" | "delete";

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
  readonly payload: RoomScreenPayload | ConversationScreenPayload;
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
  readonly kind: "hit" | "miss" | "crit" | "heal" | "system" | "queue" | "flee" | "defeat" | "error";
};

export type VitalsPayload = {
  readonly hp: number;
  readonly max_hp: number;
  readonly stamina: number;
  readonly max_stamina: number;
  readonly mana: number;
  readonly max_mana: number;
};

export type DialogResponse = {
  readonly label?: string;
  readonly next_node_id?: string;
  readonly command?: string;
  readonly condition?: string;
  readonly quest_offer_id?: string;
  readonly decline_node_id?: string;
  readonly effects?: readonly number[];
};

export type DialogNode = {
  readonly id: string;
  readonly npc_template_id: string;
  readonly npc_text: string;
  readonly responses: readonly DialogResponse[];
  readonly is_entry: boolean;
  readonly entry_condition?: string;
  readonly on_enter_effects?: readonly number[];
};

// Shop types (NPC-006 - Vendor NPCs)

export type ShopItem = {
  readonly id: number;
  readonly shop_id: number;
  readonly equipment_template_id: number;
  readonly name: string;
  readonly description: string;
  readonly category: string;
  readonly price: number;
  readonly quantity: number;
  readonly max_stock: number;
  readonly is_enabled: boolean;
  readonly last_restocked?: number;
  // Equipment-specific fields
  readonly slot?: string;
  readonly rarity?: string;
  readonly damage_dice_count?: number;
  readonly damage_dice_sides?: number;
  readonly damage_bonus?: number;
  readonly damage_type?: string;
};

export type ShopSession = {
  readonly shop_id: number;
  readonly name: string;
  readonly vendor_name: string;
  readonly npc_template_id?: string;
  readonly currency_name: string;
  readonly gold: number;
  readonly items: readonly ShopItem[];
};
