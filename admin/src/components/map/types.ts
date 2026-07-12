export type Room = Readonly<{
  id: number
  name: string
  description: string
  isStartingRoom?: boolean
  isRootRoom?: boolean
  exits: Record<string, number>
  atmosphere?: string
  posZ?: number
  tags?: string[]
  version?: number
  zoneIds?: string[]
}>

export type NPC = Readonly<{
  id: number
  name: string
  class: string
  race: string
  level: number
  currentRoomId: number
}>

export type Equipment = Readonly<{
  id: number
  name: string
  description?: string
  roomId?: number
}>

export type EquipmentTemplate = Readonly<{
  equipment_template_id: string; name: string; description: string; slot: string
  level: number; weight: number; item_type: string; color: string
  is_visible: boolean; is_immovable: boolean; effect_type: string
  effect_value: number; effect_duration: number; is_container: boolean
  container_capacity: number; is_locked: boolean
}>

export type ItemInstanceView = Readonly<{
  id: number; name: string; description: string; slot: string
  level: number; weight: number; isEquipped: boolean; isImmovable: boolean
  color: string; isVisible: boolean; itemType: string
  equipment_template_id: string; effect_type: string; effect_value: number
  effect_duration: number; healing: number; effect: string
  isContainer: boolean; containerCapacity: number; isLocked: boolean
}>

export type SpawnFormData = {
  template_id: string; name: string; description: string; slot: string
  level: number; weight: number; color: string; room_id: number
}

export type EditFormData = {
  name: string; description: string; slot: string
  level: number; weight: number; color: string
}
