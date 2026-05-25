/**
 * Combat utility functions.
 * Mirrors the SSH client's dice logic in herbst/dice/ and game_combat.go.
 */

export type DiceResult = {
  roll: number;
  total: number;
  isCrit: boolean;
  isFumble: boolean;
};

export function rollD20(modifier = 0): DiceResult {
  const roll = Math.floor(Math.random() * 20) + 1;
  return {
    roll,
    total: roll + modifier,
    isCrit: roll === 20,
    isFumble: roll === 1,
  };
}

export function rollDamage(sides: number, count: number, modifier = 0): number {
  let total = modifier;
  for (let i = 0; i < count; i++) {
    total += Math.floor(Math.random() * sides) + 1;
  }
  return Math.max(1, total);
}

/** Base player damage: 1 + strength/5 (matches server tryAttack logic) */
export function calculatePlayerDamage(strength: number): number {
  return Math.max(1, 1 + Math.floor(strength / 5));
}

/** Base enemy damage: level + 2 */
export function calculateEnemyDamage(level: number, isCrit = false): number {
  const dmg = Math.max(1, level + 2);
  return isCrit ? dmg * 2 : dmg;
}

/** Player AC: base 10 + level/2 */
export function calculatePlayerAC(level: number): number {
  return 10 + Math.floor(level / 2);
}

/** Enemy AC: base 10 + level/2 */
export function calculateEnemyAC(level: number): number {
  return 10 + Math.floor(level / 2);
}

/** DEX modifier: level / 3 */
export function getDexModifier(level: number): number {
  return Math.floor(level / 3);
}

/** STR modifier: (strength - 10) / 2 */
export function getStrModifier(strength: number): number {
  return Math.floor((strength - 10) / 2);
}

/** Flee check: d20 + level/2 vs DC 12 */
export function attemptFlee(level: number): { success: boolean; roll: number; total: number } {
  const { roll, total } = rollD20(Math.floor(level / 2));
  return { success: total >= 12, roll, total };
}