const API_BASE = (() => {
  // Production build override (set at docker build time)
  if (import.meta.env.VITE_API_URL) return import.meta.env.VITE_API_URL;
  // Same-origin in production (nginx reverse-proxy), port 8080 for local dev
  const { protocol, hostname, port } = window.location;
  if (port === "8080") return `${protocol}//${hostname}:8080`;
  if (port === "5174" || port === "5173") return `${protocol}//${hostname}:8080`;
  return `${protocol}//${hostname}`;
})();

function getToken(): string | null {
  return localStorage.getItem("herbst_token");
}

function headers(): HeadersInit {
  const h: Record<string, string> = { "Content-Type": "application/json" };
  const token = getToken();
  if (token) {
    return { ...h, Authorization: `Bearer ${token}` };
  }
  return h;
}

export type Ability = {
  readonly id: number;
  readonly name: string;
  readonly description: string;
  readonly ability_type: string;
  readonly cost: number;
  readonly cooldown: number;
  readonly mana_cost: number;
  readonly stamina_cost: number;
  readonly hp_cost: number;
};

export type User = {
  readonly id: number;
  readonly email: string;
  readonly is_admin: boolean;
  readonly allowed_worlds: string;
};

export type World = {
  readonly name: string;
  readonly file: string;
};

export type Character = {
  readonly id: number;
  readonly name: string;
  readonly isNPC: boolean;
  readonly is_admin: boolean;
  readonly currentRoomId: number;
  readonly currentWorld: string;
  readonly hitpoints: number;
  readonly max_hitpoints: number;
  readonly stamina: number;
  readonly max_stamina: number;
  readonly mana: number;
  readonly max_mana: number;
  readonly race: string;
  readonly gender: string;
  readonly level: number;
  readonly class: string;
};

export type Race = {
  readonly name: string;
  readonly display_name: string;
  readonly description: string;
  readonly stat_modifiers: Readonly<Record<string, unknown>>;
  readonly skill_grants: Readonly<Record<string, unknown>>;
  readonly equipment_slots: readonly string[];
};

export type EquipmentItem = {
  readonly id: number;
  readonly name: string;
  readonly description: string;
  readonly slot: string;
  readonly level: number;
  readonly weight: number;
  readonly isEquipped: boolean;
  readonly rarity: string;
  readonly color: string;
  readonly armor_rating: number;
  readonly armor_type: string;
  readonly damage_dice_count: number;
  readonly damage_dice_sides: number;
  readonly damage_bonus: number;
  readonly damage_type: string;
  readonly weapon_type: string;
  readonly is_two_handed: boolean;
  readonly ownerId: number | null;
};

export type Gender = {
  readonly name: string;
  readonly display_name: string;
  readonly subject_pronoun: string;
  readonly object_pronoun: string;
  readonly possessive_pronoun: string;
};

export async function login(email: string, password: string): Promise<{ readonly token: string; readonly user: User }> {
  const res = await fetch(`${API_BASE}/users/auth`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password }),
  });
  if (!res.ok) {
    const err = await res.json().catch((): { readonly error: string } => ({ error: "Login failed" }));
    throw new Error(err.error || "Login failed");
  }
  const data = await res.json();
  return { token: data.token, user: data };
}

export async function me(): Promise<User> {
  const res = await fetch(`${API_BASE}/users/me`, {
    headers: headers(),
  });
  if (!res.ok) {
    const err = await res.json().catch((): { readonly error: string } => ({ error: "Session invalid" }));
    throw new Error(err.error || "Session invalid");
  }
  return res.json();
}

export async function listWorlds(): Promise<{ readonly worlds: readonly World[]; readonly count: number }> {
  const res = await fetch(`${API_BASE}/admin/export/worlds`, {
    headers: headers(),
  });
  if (!res.ok) {
    const err = await res.json().catch((): { readonly error: string } => ({ error: "Failed to load worlds" }));
    throw new Error(err.error || "Failed to load worlds");
  }
  return res.json();
}

export async function listMyCharacters(): Promise<readonly Character[]> {
  const res = await fetch(`${API_BASE}/api/me/characters`, {
    headers: headers(),
  });
  if (!res.ok) {
    const err = await res.json().catch((): { readonly error: string } => ({ error: "Failed to load characters" }));
    throw new Error(err.error || "Failed to load characters");
  }
  return res.json();
}

export async function createCharacter(input: {
  readonly name: string;
  readonly race: string;
  readonly gender: string;
  readonly world: string;
}): Promise<Character> {
  const res = await fetch(`${API_BASE}/api/me/characters`, {
    method: "POST",
    headers: headers(),
    body: JSON.stringify(input),
  });
  if (!res.ok) {
    const err = await res.json().catch((): { readonly error: string } => ({ error: "Failed to create character" }));
    throw new Error(err.error || "Failed to create character");
  }
  return res.json();
}

export async function listRaces(): Promise<readonly Race[]> {
  const res = await fetch(`${API_BASE}/races`, { headers: headers() });
  if (!res.ok) {
    const err = await res.json().catch((): { readonly error: string } => ({ error: "Failed to load races" }));
    throw new Error(err.error || "Failed to load races");
  }
  return res.json();
}

export async function listGenders(): Promise<readonly Gender[]> {
  const res = await fetch(`${API_BASE}/genders`, { headers: headers() });
  if (!res.ok) {
    const err = await res.json().catch((): { readonly error: string } => ({ error: "Failed to load genders" }));
    throw new Error(err.error || "Failed to load genders");
  }
  return res.json();
}

export async function listClasslessAbilities(): Promise<readonly Ability[]> {
  const res = await fetch(`${API_BASE}/abilities/classless`, { headers: headers() });
  if (!res.ok) {
    const err = await res.json().catch((): { readonly error: string } => ({ error: "Failed to load abilities" }));
    throw new Error(err.error || "Failed to load abilities");
  }
  const data = await res.json() as { abilities: readonly Ability[] };
  return data.abilities;
}

export async function getCharacterAbilities(charID: number): Promise<{ slots: (Ability & { slot: number })[] }> {
  const res = await fetch(`${API_BASE}/characters/${charID}/abilities`, { headers: headers() });
  if (!res.ok) {
    const err = await res.json().catch((): { readonly error: string } => ({ error: "Failed to load character abilities" }));
    throw new Error(err.error || "Failed to load character abilities");
  }
  return res.json();
}

export async function equipAbility(charID: number, abilityID: number, slot: number): Promise<void> {
  const res = await fetch(`${API_BASE}/characters/${charID}/abilities`, {
    method: "POST",
    headers: headers(),
    body: JSON.stringify({ ability_id: abilityID, slot }),
  });
  if (!res.ok) {
    const err = await res.json().catch((): { readonly error: string } => ({ error: "Failed to equip ability" }));
    throw new Error(err.error || "Failed to equip ability");
  }
}

export async function unequipAbility(charID: number, slot: number): Promise<void> {
  const res = await fetch(`${API_BASE}/characters/${charID}/abilities/${slot}`, {
    method: "DELETE",
    headers: headers(),
  });
  if (!res.ok) {
    const err = await res.json().catch((): { readonly error: string } => ({ error: "Failed to unequip ability" }));
    throw new Error(err.error || "Failed to unequip ability");
  }
}

export async function getCombatStatus(charID: number): Promise<{ hp: number; maxHp: number; isNPC: boolean }> {
  const res = await fetch(`${API_BASE}/characters/${charID}/combat-status`, { headers: headers() });
  if (!res.ok) {
    const err = await res.json().catch((): { readonly error: string } => ({ error: "Failed to load combat status" }));
    throw new Error(err.error || "Failed to load combat status");
  }
  return res.json();
}

export async function applyDamage(charID: number, damage: number, attackerID?: number): Promise<void> {
  const body: Record<string, unknown> = { damage };
  if (attackerID != null && attackerID > 0) body.attacker_id = attackerID;
  const res = await fetch(`${API_BASE}/characters/${charID}/damage`, {
    method: "POST",
    headers: headers(),
    body: JSON.stringify(body),
  });
  if (!res.ok) {
    const err = await res.json().catch((): { readonly error: string } => ({ error: "Failed to apply damage" }));
    throw new Error(err.error || "Failed to apply damage");
  }
}

export async function healCharacter(charID: number, amount: number): Promise<void> {
  const res = await fetch(`${API_BASE}/characters/${charID}/heal`, {
    method: "POST",
    headers: headers(),
    body: JSON.stringify({ amount }),
  });
  if (!res.ok) {
    const err = await res.json().catch((): { readonly error: string } => ({ error: "Failed to heal" }));
    throw new Error(err.error || "Failed to heal");
  }
}

export async function getCharacterEquipment(charID: number): Promise<readonly EquipmentItem[]> {
  const res = await fetch(`${API_BASE}/equipment?ownerId=${charID}`, { headers: headers() });
  if (!res.ok) {
    const err = await res.json().catch((): { readonly error: string } => ({ error: "Failed to load equipment" }));
    console.error({ endpoint: `/equipment?ownerId=${charID}`, status: res.status, body: err, characterID: charID });
    throw new Error(err.error || "Failed to load equipment");
  }
  return res.json();
}

export async function equipItem(itemID: number, charID: number): Promise<void> {
  const res = await fetch(`${API_BASE}/equipment/${itemID}/equip`, {
    method: "PUT",
    headers: headers(),
    body: JSON.stringify({ character_id: charID }),
  });
  if (!res.ok) {
    const err = await res.json().catch((): { readonly error: string } => ({ error: "Failed to equip item" }));
    console.error({ endpoint: `/equipment/${itemID}/equip`, status: res.status, body: err, characterID: charID });
    throw new Error(err.error || "Failed to equip item");
  }
}

export async function unequipItem(itemID: number, charID: number): Promise<void> {
  const res = await fetch(`${API_BASE}/equipment/${itemID}/unequip`, {
    method: "PUT",
    headers: headers(),
    body: JSON.stringify({ character_id: charID }),
  });
  if (!res.ok) {
    const err = await res.json().catch((): { readonly error: string } => ({ error: "Failed to unequip item" }));
    console.error({ endpoint: `/equipment/${itemID}/unequip`, status: res.status, body: err, characterID: charID });
    throw new Error(err.error || "Failed to unequip item");
  }
}