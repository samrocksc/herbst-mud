import { describe, it, expect } from 'vitest';
import { rollD20, rollDamage, calculatePlayerDamage } from './combat';

describe('combat', () => {
  describe('rollD20', () => {
    it('should return a result with roll between 1 and 20', () => {
      const result = rollD20();
      expect(result.roll).toBeGreaterThanOrEqual(1);
      expect(result.roll).toBeLessThanOrEqual(20);
      expect(result.total).toBeGreaterThanOrEqual(1);
      expect(result.total).toBeLessThanOrEqual(20);
    });

    it('should apply modifier to total', () => {
      const result = rollD20(5);
      expect(result.total).toBe(result.roll + 5);
    });
  });

  describe('rollDamage', () => {
    it('should return minimum damage with 0 dice', () => {
      expect(rollDamage(6, 0, 2)).toBe(2);
    });

    it('should return damage within expected range', () => {
      // 1d6 + 2 should be 3-8
      const result = rollDamage(6, 1, 2);
      expect(result).toBeGreaterThanOrEqual(3);
      expect(result).toBeLessThanOrEqual(8);
    });

    it('should return minimum with modifier only', () => {
      expect(rollDamage(6, 0, 5)).toBe(5);
    });
  });

  describe('calculatePlayerDamage', () => {
    it('should calculate damage with weapon', () => {
      const equipment = [
        {
          id: 1,
          name: 'Sword',
          slot: 'main_hand',
          damage_dice_count: 1,
          damage_dice_sides: 6,
          damage_bonus: 2,
          isEquipped: true,
        },
      ] as any;

      // With STR=10, strMod=0, so damage should be 1d6 + 2 = 3-8
      const result = calculatePlayerDamage(10, equipment);
      expect(result).toBeGreaterThanOrEqual(3);
      expect(result).toBeLessThanOrEqual(8);
    });

    it('should use fists damage with no weapon', () => {
      // With STR=10, fists damage = 1 + floor(10/5) = 3
      const result = calculatePlayerDamage(10, []);
      expect(result).toBe(3);
    });

    it('should use fists damage when equipment is undefined', () => {
      const result = calculatePlayerDamage(10, undefined);
      expect(result).toBe(3);
    });
  });
});
