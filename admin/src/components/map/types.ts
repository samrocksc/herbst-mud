export type Room = Readonly<{
  id: number
  name: string
  description: string
  isStartingRoom?: boolean
  exits: Record<string, number>
  atmosphere?: string
  posX?: number
  posY?: number
  version?: number
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
