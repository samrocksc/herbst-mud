import { CombatFieldsDisplay } from '../../components/CombatFieldsDisplay'
import type { ItemTemplate } from './items.$itemId'

/** Read-only detail view for an item template. */
export function ItemDetailView({ template }: Readonly<{ template: ItemTemplate }>) {
  return (
    <div className="bg-surface-muted rounded-lg p-6 border border-border">
      <h2 className="mt-0 mb-4 text-text text-lg font-semibold">Item Stats</h2>
      <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
        <DetailField label="ID" value={template.id} />
        <DetailField label="Name" value={template.name} />
        <DetailField label="Slot" value={template.slot} />
        <DetailField label="Level" value={String(template.level)} />
        <DetailField label="Weight" value={String(template.weight)} />
        <DetailField label="Type" value={template.item_type} />
        <DetailField label="Color" value={template.color || '—'} />
        <DetailField label="Visible" value={template.is_visible ? 'Yes' : 'No'} />
        <DetailField label="Immovable" value={template.is_immovable ? 'Yes' : 'No'} />
        <DetailField label="Container" value={template.is_container ? `Yes (${template.container_capacity})` : 'No'} />
        <DetailField label="Locked" value={template.is_locked ? 'Yes' : 'No'} />
        <DetailField label="Key Item" value={template.key_item_id || '—'} />
        {template.effect_type && (
          <>
            <DetailField label="Effect" value={template.effect_type} />
            <DetailField label="Effect Value" value={String(template.effect_value)} />
            <DetailField label="Duration" value={String(template.effect_duration)} />
          </>
        )}
        {template.expires_at && <DetailField label="Expires" value={new Date(template.expires_at).toLocaleDateString()} />}
      </div>
      <p className="text-text text-sm mt-4">{template.description || 'No description.'}</p>
      <div className="mt-4 pt-4 border-t border-border">
        <CombatFieldsDisplay
          slot={template.slot} damage_dice_count={template.damage_dice_count}
          damage_dice_sides={template.damage_dice_sides} damage_bonus={template.damage_bonus}
          weapon_type={template.weapon_type} damage_type={template.damage_type}
          is_two_handed={template.is_two_handed} armor_rating={template.armor_rating}
          armor_type={template.armor_type} rarity={template.rarity}
          skill_requirement={template.skill_requirement}
          skill_requirement_level={template.skill_requirement_level} stats={template.stats ? JSON.parse(template.stats) : null}
        />
      </div>
    </div>
  )
}

function DetailField({ label, value }: Readonly<{ label: string; value: string }>) {
  return (<div><span className="text-text-muted text-xs block mb-0.5">{label}</span><span className="text-text text-sm font-medium">{value}</span></div>)
}