/** Character stats display component. */

export function CharacterStats({ character }: Readonly<{
  character: Readonly<{
    id: number
    name: string
    isNPC: boolean
    currentRoomId: number
    hitpoints: number
    max_hitpoints: number
    stamina: number
    max_stamina: number
    mana: number
    max_mana: number
    race: string
    class: string
    level: number
    xp: number
    strength: number
    dexterity: number
    intelligence: number
    wisdom: number
    constitution: number
    gender: string
    description: string
    is_admin: boolean
    is_immortal: boolean
  }>
}>) {
  return (
    <div className="bg-surface-muted rounded-lg p-6 border border-border mb-6">
      <h2 className="mt-0 mb-4 text-text text-lg font-semibold">Character Stats</h2>
      <div className="grid grid-cols-2 md:grid-cols-3 gap-x-6 gap-y-3">
        <DetailField label="ID" value={String(character.id)} />
        <DetailField label="Name" value={character.name} />
        <DetailField label="Race" value={character.race} />
        <DetailField label="Class" value={character.class} />
        <DetailField label="Level" value={String(character.level)} />
        <DetailField label="XP" value={String(character.xp)} />
        <DetailField label="Room" value={String(character.currentRoomId)} />
        <DetailField label="HP" value={`${character.hitpoints} / ${character.max_hitpoints}`} />
        <DetailField label="Stamina" value={`${character.stamina} / ${character.max_stamina}`} />
        <DetailField label="Mana" value={`${character.mana} / ${character.max_mana}`} />
        <DetailField label="STR" value={String(character.strength)} />
        <DetailField label="DEX" value={String(character.dexterity)} />
        <DetailField label="INT" value={String(character.intelligence)} />
        <DetailField label="WIS" value={String(character.wisdom)} />
        <DetailField label="CON" value={String(character.constitution)} />
        <DetailField label="Gender" value={character.gender || "—"} />
        <BoolBadge value={character.isNPC} label="NPC" />
        <BoolBadge value={character.is_admin} label="Admin" />
        <BoolBadge value={character.is_immortal} label="Immortal" />
      </div>
      {character.description && (
        <div className="mt-4 pt-4 border-t border-border">
          <span className="text-text-muted text-xs block mb-1">Description</span>
          <span className="text-text text-sm">{character.description}</span>
        </div>
      )}
    </div>
  );
}

function DetailField({ label, value }: Readonly<{ label: string; value: string }>) {
  return (<div><span className="text-text-muted text-xs block mb-0.5">{label}</span><span className="text-text text-sm font-medium">{value}</span></div>);
}

function BoolBadge({ value, label }: Readonly<{ value: boolean; label: string }>) {
  const cls = value
    ? "inline-block px-2 py-0.5 rounded text-xs font-medium bg-green-900/30 text-green-400 border border-green-700/40"
    : "inline-block px-2 py-0.5 rounded text-xs font-medium bg-red-900/30 text-red-400 border border-red-700/40";
  return (<div><span className="text-text-muted text-xs block mb-0.5">{label}</span><span className={cls}>{value ? "Yes" : "No"}</span></div>);
}