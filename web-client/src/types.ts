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

// --- WebSocket protocol types ---

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