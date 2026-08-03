/** DialogNode represents a single node in an NPC's dialog tree. */
export interface DialogNode {
  readonly id: string;
  readonly npc_template_id: string;
  readonly npc_text: string;
  readonly responses: readonly DialogResponse[];
  readonly is_entry: boolean;
  readonly entry_condition?: string;
  readonly on_enter_effects?: readonly number[];
}

/** DialogResponse represents a player response option in a dialog node. */
export interface DialogResponse {
  readonly label: string;
  readonly next_node_id: string;
  readonly condition?: string;
  readonly quest_offer_id?: string;
  readonly decline_node_id?: string;
  readonly effects?: readonly number[];
}