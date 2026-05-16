import { CombatFieldsDisplay } from "../../components/CombatFieldsDisplay";
import type { ItemInstance } from "../../hooks/useItemInstances";

/** Read-only detail view for an item instance. */
export function InstanceDetailView({ instance }: Readonly<{ instance: ItemInstance }>) {
  return (
    <div className="max-w-2xl">
      <div className="bg-surface-muted rounded-lg p-6 border border-border mb-6">
        <h2 className="mt-0 mb-4 text-text text-lg font-semibold">Instance Stats</h2>
        <div className="grid grid-cols-2 gap-x-6 gap-y-3">
          <DetailField label="ID" value={String(instance.id)} />
          <DetailField label="Name" value={instance.name} />
          <DetailField label="Type" value={instance.itemType} />
          <DetailField label="Slot" value={instance.slot} />
          <DetailField label="Level" value={String(instance.level)} />
          <DetailField label="Weight" value={String(instance.weight)} />
          {instance.ownerId ? <DetailField label="Owner ID" value={<span className="text-primary">{String(instance.ownerId)}</span>} /> : <DetailField label="Owner ID" value="None" />}
          {instance.roomId && !instance.ownerId ? <DetailField label="Room ID" value={String(instance.roomId)} /> : <DetailField label="Room ID" value={instance.roomId ? String(instance.roomId) : "None"} />}
          <DetailField label="Color" value={instance.color || "none"} />
          <BoolBadge value={instance.isVisible} label="Visible" />
          <BoolBadge value={instance.isImmovable} label="Immovable" />
          <BoolBadge value={instance.isEquipped} label="Equipped" />
        </div>
        <div className="mt-4 pt-4 border-t border-border">
          <CombatFieldsDisplay slot={instance.slot} damage_dice_count={instance.damage_dice_count}
            damage_dice_sides={instance.damage_dice_sides} damage_bonus={instance.damage_bonus}
            weapon_type={instance.weapon_type} damage_type={instance.damage_type}
            is_two_handed={instance.is_two_handed} armor_rating={instance.armor_rating}
            armor_type={instance.armor_type} rarity={instance.rarity}
            skill_requirement={instance.skill_requirement} skill_requirement_level={instance.skill_requirement_level}
            stats={instance.stats} />
        </div>
      </div>
      {instance.ownerId && <LocationNote text="Held by Character" value={`#${instance.ownerId}`} primary />}
      {!instance.ownerId && instance.roomId && <LocationNote text="In Room" value={`#${instance.roomId}`} />}
    </div>
  );
}

function LocationNote({ text, value, primary }: Readonly<{ text: string; value: string; primary?: boolean }>) {
  return (
    <div className="bg-surface-muted rounded-lg p-4 border border-border mb-6">
      <span className="text-text-muted text-sm">{text} <span className={primary ? "text-primary font-medium" : "text-text font-medium"}>{value}</span></span>
    </div>
  );
}

function DetailField({ label, value }: Readonly<{ label: string; value: React.ReactNode }>) {
  return (<div><span className="text-text-muted text-xs block mb-0.5">{label}</span><span className="text-text text-sm font-medium">{value}</span></div>);
}

function BoolBadge({ value, label }: Readonly<{ value: boolean; label: string }>) {
  const cls = value ? "inline-block px-2 py-0.5 rounded text-xs font-medium bg-green-900/30 text-green-400 border border-green-700/40" : "inline-block px-2 py-0.5 rounded text-xs font-medium bg-red-900/30 text-red-400 border border-red-700/40";
  return (<div><span className="text-text-muted text-xs block mb-0.5">{label}</span><span className={cls}>{value ? "Yes" : "No"}</span></div>);
}