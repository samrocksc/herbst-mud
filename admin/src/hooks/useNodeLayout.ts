import { useMemo } from 'react'
import { resolveOverlaps } from '../components/map/ExitLines'
import { DIRECTION_OFFSETS } from '../components/map/DirectionUtils'
import { ORPHAN_COLS } from '../components/map/constants'
import type { Room } from '../components/map/types'

export function useNodeLayout(rooms: Room[], currentZLevel: number) {
  const zLevels = useMemo(() => {
    const zMap = new Map<number, number>()
    const visited = new Set<number>()
    const assign = (id: number, z: number) => {
      if (visited.has(id)) return
      visited.add(id)
      zMap.set(id, z)
      const room = rooms.find(r => r.id === id)
      if (!room) return
      for (const [dir, targetId] of Object.entries(room.exits || {})) {
        if (!targetId) continue
        const tz = dir === 'up' ? z + 1 : dir === 'down' ? z - 1 : z
        assign(targetId, tz)
      }
    }
    const start = rooms.find(r => r.isStartingRoom) || rooms[0]
    if (start) assign(start.id, 0)
    for (const r of rooms) {
      if (!zMap.has(r.id)) zMap.set(r.id, 0)
    }
    return zMap
  }, [rooms])

  const nodePositions = useMemo(() => {
    const positions = new Map<number, { x: number; y: number }>()
    const positioned = new Set<number>()
    const posRoom = (id: number, x: number, y: number) => {
      const rz = zLevels.get(id) || 0
      if (rz !== currentZLevel) return
      if (positioned.has(id)) return
      positioned.add(id)
      positions.set(id, { x, y })
      const room = rooms.find(r => r.id === id)
      if (!room) return
      for (const [dir, tid] of Object.entries(room.exits || {})) {
        if (!tid || dir === 'up' || dir === 'down') continue
        const off = DIRECTION_OFFSETS[dir] || { dx: 150, dy: 0 }
        const tr = rooms.find(r => r.id === tid)
        const tx = tr?.posX != null ? tr.posX : x + off.dx
        const ty = tr?.posY != null ? tr.posY : y + off.dy
        posRoom(tid, tx, ty)
      }
    }
    const start = rooms.find(r => r.isStartingRoom) || rooms[0]
    if (start) {
      posRoom(start.id, start.posX ?? 400, start.posY ?? 300)
    }
    let oi = 0
    for (const r of rooms) {
      const rz = zLevels.get(r.id) || 0
      if (rz === currentZLevel && !positioned.has(r.id)) {
        const col = oi % ORPHAN_COLS
        const row = Math.floor(oi / ORPHAN_COLS)
        positions.set(r.id, { x: r.posX ?? 800 + col * 180, y: r.posY ?? 300 + row * 120 })
        oi++
      }
    }
    return resolveOverlaps(positions, 50)
  }, [rooms, zLevels, currentZLevel])

  return { zLevels, nodePositions }
}