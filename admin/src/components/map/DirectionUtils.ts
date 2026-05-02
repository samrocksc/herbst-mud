export const OPPOSITE_DIR: Record<string, string> = {
  north: 'south',
  south: 'north',
  east: 'west',
  west: 'east',
  northeast: 'southwest',
  southwest: 'northeast',
  northwest: 'southeast',
  southeast: 'northwest',
  up: 'down',
  down: 'up',
}

export const DIRECTION_OFFSETS: Record<string, { dx: number; dy: number }> = {
  north: { dx: 0, dy: -120 },
  south: { dx: 0, dy: 120 },
  east: { dx: 150, dy: 0 },
  west: { dx: -150, dy: 0 },
  northeast: { dx: 106, dy: -106 },
  northwest: { dx: -106, dy: -106 },
  southeast: { dx: 106, dy: 106 },
  southwest: { dx: -106, dy: 106 },
}

export const ALL_DIRECTIONS = [
  'north',
  'northeast',
  'east',
  'southeast',
  'south',
  'southwest',
  'west',
  'northwest',
  'up',
  'down',
]
