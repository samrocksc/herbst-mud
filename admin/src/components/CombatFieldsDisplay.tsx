import { isWeaponSlot, isArmorSlot } from './equipConstants';

/** Read-only display of weapon/armor combat fields. */
export function CombatFieldsDisplay(props: Readonly<{
  slot: string
  damage_dice_count: number
  damage_dice_sides: number
  damage_bonus: number
  weapon_type: string
  damage_type: string
  is_two_handed: boolean
  armor_rating: number
  armor_type: string
  rarity: string
  skill_requirement: string
  skill_requirement_level: number
  stats: Record<string, unknown> | null
}>) {
  const showWeapon = isWeaponSlot(props.slot);
  const showArmor = isArmorSlot(props.slot);

  return (
    <div className="space-y-3">
      <h3 className="text-text text-sm font-semibold border-b border-border pb-1">Combat Stats</h3>
      <div className="grid grid-cols-2 md:grid-cols-3 gap-x-6 gap-y-3">
        <DetailField label="Rarity" value={props.rarity || 'common'} />
        {showWeapon && (
          <>
            <DetailField label="Damage" value={props.damage_dice_count > 0
              ? `${props.damage_dice_count}d${props.damage_dice_sides}${props.damage_bonus > 0 ? `+${props.damage_bonus}` : ''}`
              : '—'} />
            <DetailField label="Weapon Type" value={props.weapon_type || '—'} />
            <DetailField label="Damage Type" value={props.damage_type || '—'} />
            <DetailField label="Two-Handed" value={props.is_two_handed ? 'Yes' : 'No'} />
          </>
        )}
        {showArmor && (
          <>
            <DetailField label="Armor Rating" value={String(props.armor_rating)} />
            <DetailField label="Armor Type" value={props.armor_type || '—'} />
          </>
        )}
        <DetailField label="Skill Requirement" value={props.skill_requirement || '—'} />
        <DetailField label="Skill Req. Level" value={String(props.skill_requirement_level)} />
      </div>
      {props.stats && Object.keys(props.stats).length > 0 && <StatsDisplay stats={props.stats as Record<string, number>} />}
    </div>
  );
}

function DetailField({ label, value }: Readonly<{ label: string; value: string }>) {
  return (<div><span className="text-text-muted text-xs block mb-0.5">{label}</span><span className="text-text text-sm font-medium">{value}</span></div>);
}

function StatsDisplay({ stats }: Readonly<{ stats: Record<string, number> }>) {
  return (
    <div className="pt-2 border-t border-border">
      <span className="text-text-muted text-xs block mb-1">Stats</span>
      <div className="grid grid-cols-3 gap-x-4 gap-y-1">
        {Object.entries(stats).map(([stat, val]) => (
          <div key={stat} className="flex justify-between text-sm">
            <span className="text-text-muted">{stat}</span>
            <span className="text-text font-medium">{val}</span>
          </div>
        ))}
      </div>
    </div>
  );
}