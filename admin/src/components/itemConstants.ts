/** Shared item constants for the admin UI. */

/** Slot options for item select dropdowns. */
export const SLOT_OPTIONS = [
  { value: '', label: '-- none --' },
  { value: 'head', label: 'Head' },
  { value: 'neck', label: 'Neck' },
  { value: 'chest', label: 'Chest' },
  { value: 'back', label: 'Back' },
  { value: 'hands', label: 'Hands' },
  { value: 'legs', label: 'Legs' },
  { value: 'feet', label: 'Feet' },
  { value: 'finger_left', label: 'Finger (Left)' },
  { value: 'finger_right', label: 'Finger (Right)' },
  { value: 'main_hand', label: 'Main Hand' },
  { value: 'off_hand', label: 'Off Hand' },
  { value: 'tail', label: 'Tail' },
  { value: 'horn', label: 'Horn' },
  { value: 'wings', label: 'Wings' },
] as const;

/** Item type options for select dropdowns. */
export const ITEM_TYPE_OPTIONS = [
  { value: 'misc', label: 'Misc' },
  { value: 'weapon', label: 'Weapon' },
  { value: 'armor', label: 'Armor' },
  { value: 'consumable', label: 'Consumable' },
  { value: 'quest', label: 'Quest' },
  { value: 'container', label: 'Container' },
  { value: 'potion', label: 'Potion' },
] as const;