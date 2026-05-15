import { Link } from '@tanstack/react-router';
import { isWeaponSlot, isArmorSlot } from './equipConstants';
import { formatSlotName, RarityBadge } from './equipViewHelpers';

type EquippedItem = Readonly<{
  id: number; name: string; slot: string; level: number; itemType: string
  isEquipped: boolean; damage_dice_count: number; damage_dice_sides: number
  damage_bonus: number; damage_type: string; weapon_type: string
  armor_rating: number; armor_type: string; rarity: string; equipment_template_id: string
}>

/** Grid of slot cards showing equipped/empty status. */
export function EquippedSlotGrid({ slots, items }: Readonly<{ slots: string[]; items: EquippedItem[] }>) {
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
      {slots.map((slot) => (
        <SlotCard key={slot} slot={slot} item={items.find((i) => i.slot === slot)} />
      ))}
    </div>
  );
}

/** List of unequipped items in inventory. */
export function UnequippedSection({ items }: Readonly<{ items: EquippedItem[] }>) {
  return (
    <div className="mt-4 pt-4 border-t border-border">
      <h3 className="text-text text-sm font-semibold mb-2">Inventory (Unequipped)</h3>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-2">
        {items.map((item) => (
          <div key={item.id} className="rounded border border-border bg-surface p-2">
            <Link to="/items/$itemId/instances/$instanceId"
              params={{ itemId: item.equipment_template_id || String(item.id), instanceId: String(item.id) }}
              className="text-primary hover:underline text-sm font-medium no-underline">{item.name}</Link>
            <span className="text-text-muted text-xs ml-2">Lv {item.level} {item.itemType}</span>
          </div>
        ))}
      </div>
    </div>
  );
}

function SlotCard({ slot, item }: Readonly<{ slot: string; item: EquippedItem | undefined }>) {
  return (
    <div className={`rounded border p-3 ${item ? 'border-primary/40 bg-primary/5' : 'border-border bg-surface'}`}>
      <div className="flex items-center justify-between mb-1">
        <span className="text-text-muted text-xs font-medium uppercase tracking-wide">{formatSlotName(slot)}</span>
        {item && <RarityBadge rarity={item.rarity} />}
      </div>
      {item ? <ItemInfo item={item} slot={slot} /> : <span className="text-text-muted text-xs italic">Empty</span>}
    </div>
  );
}

function ItemInfo({ item, slot }: Readonly<{ item: EquippedItem; slot: string }>) {
  return (
    <div>
      <Link to="/items/$itemId/instances/$instanceId"
        params={{ itemId: item.equipment_template_id || String(item.id), instanceId: String(item.id) }}
        className="text-primary hover:underline text-sm font-medium no-underline">{item.name}</Link>
      <div className="text-text-muted text-xs mt-1">
        Lv {item.level} {item.itemType}
        {isWeaponSlot(slot) && item.damage_dice_count > 0 && (
          <> &mdash; {item.damage_dice_count}d{item.damage_dice_sides}{item.damage_bonus > 0 ? `+${item.damage_bonus}` : ''} {item.damage_type}</>
        )}
        {isArmorSlot(slot) && item.armor_rating > 0 && (
          <> &mdash; AC {item.armor_rating}</>
        )}
      </div>
    </div>
  );
}