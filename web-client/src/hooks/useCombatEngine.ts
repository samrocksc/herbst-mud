import { useCallback, useEffect, useRef, useState } from "react";
import type { CombatTarget, CombatLogEntry, CharacterSkill } from "../lib/types";
import type { EquipmentItem } from "../lib/api";
import { getCombatStatus, applyDamage, healCharacter, getCharacterEquipment, adjustMana } from "../lib/api";
import {
  rollD20,
  calculatePlayerDamage,
  calculateEnemyDamage,
  calculateEnemyAC,
  getDexModifier,
  attemptFlee,
} from "../lib/combat";

const TICK_MS = 1500;

// AbilityData represents an ability for combat display
type AbilityData = {
  readonly id: number;
  readonly name: string;
  readonly description?: string;
  readonly casterMessage?: string;
  readonly recipientMessage?: string;
  readonly manaCost?: number;
  readonly staminaCost?: number;
  readonly cooldown?: number;
};

// EffectData represents an ability effect
type EffectData = {
  readonly effectType: string;
  readonly effectMessage?: string;
  readonly target: string;
  readonly value: number;
};

type Props = {
  characterID: number;
  characterLevel: number;
  characterStrength: number;
  initialHP: number;
  initialMaxHP: number;
  initialMana: number;
  initialMaxMana: number;
  skills: readonly CharacterSkill[];
  onLog: (text: string, kind: CombatLogEntry["kind"]) => void;
  onCombatEnd: () => void;
  onPlayerHPChange: (hp: number) => void;
  // Optional: pass ability details to display messages
  abilities?: readonly AbilityData[];
  effects?: readonly EffectData[];
};

export function useCombatEngine({
  characterID,
  characterLevel,
  characterStrength,
  initialHP,
  initialMaxHP,
  initialMana,
  initialMaxMana,
  onLog,
  onCombatEnd,
  onPlayerHPChange,
  abilities = [],
  effects = [],
}: Props) {
  const [inCombat, setInCombat] = useState(false);
  const [targets, setTargets] = useState<CombatTarget[]>([]);
  const [combatLog, setCombatLog] = useState<CombatLogEntry[]>([]);
  const [round, setRound] = useState(0);
  const [queuedAction, setQueuedAction] = useState<string | null>(null);
  const [playerHP, setPlayerHP] = useState(initialHP);
  const [playerMana, setPlayerMana] = useState(initialMana);
  const [playerMaxMana, setPlayerMaxMana] = useState(initialMaxMana);

  const targetsRef = useRef(targets);
  const playerHPRef = useRef(playerHP);
  const playerManaRef = useRef(playerMana);
  const inCombatRef = useRef(inCombat);
  const roundRef = useRef(round);
  const queuedRef = useRef(queuedAction);
  const tickRef = useRef<ReturnType<typeof setInterval> | null>(null);
  const isProcessingRef = useRef(false);

  targetsRef.current = targets;
  playerHPRef.current = playerHP;
  playerManaRef.current = playerMana;
  inCombatRef.current = inCombat;
  roundRef.current = round;
  queuedRef.current = queuedAction;

  const addLog = useCallback(
    (text: string, kind: CombatLogEntry["kind"]) => {
      const entry: CombatLogEntry = { timestamp: Date.now(), text, kind };
      setCombatLog((prev) => [...prev.slice(-99), entry]);
      onLog(text, kind);
    },
    [onLog]
  );

  // formatCombatMessage generates the combat message based on ability templates
  const formatCombatMessage = useCallback((ability: AbilityData, effect: EffectData, targetName: string): string => {
    // Priority: CasterMessage -> RecipientMessage -> EffectMessage -> Default

    // Helper to replace all placeholders
    const replacePlaceholders = (msg: string): string => {
      return msg
        .replace("{actor}", "You")
        .replace("{recipient}", targetName)
        .replace("{target}", targetName)
        .replace("{targetName}", targetName)
        .replace("{abilityName}", ability.name)
        .replace("{skillName}", ability.name)
        .replace("{effectType}", effect.effectType)
        .replace("{value}", String(effect.value));
    };

    if (ability.casterMessage) {
      return replacePlaceholders(ability.casterMessage);
    }

    if (targetName && ability.recipientMessage) {
      return replacePlaceholders(ability.recipientMessage);
    }

    if (effect.effectMessage) {
      return replacePlaceholders(effect.effectMessage);
    }

    // Default messages based on effect type
    if (effect.effectType === "damage") {
      return `Your ${ability.name} hits ${targetName} for ${effect.value} damage!`;
    }
    if (effect.effectType === "heal") {
      return `You cast ${ability.name} and heal ${effect.value} HP!`;
    }
    if (effect.effectType === "stun") {
      return `${ability.name} stuns ${targetName}!`;
    }

    return `${ability.name} used on ${targetName}!`;
  }, []);

  const fetchTargetHP = useCallback(async (targetID: number) => {
    try {
      const status = await getCombatStatus(targetID);
      return status.hp;
    } catch {
      return null;
    }
  }, []);

  // Fetch equipment at combat start
  // Use a ref for immediate access to equipment (state updates are async)
  const equipmentRef = useRef<readonly EquipmentItem[]>([]);
  const [equipmentState, setEquipmentState] = useState<readonly EquipmentItem[]>([]);

  const fetchEquipment = useCallback(async () => {
    if (characterID > 0) {
      try {
        const equp = await getCharacterEquipment(characterID);
        equipmentRef.current = equp;
        setEquipmentState(equp);
        console.log("[Combat] Equipment loaded:", equp.length, "items");
        if (equp.length > 0) {
          const mainHand = equp.find(e => e.slot === "main_hand" && e.damage_dice_count > 0);
          if (mainHand) {
            console.log("[Combat] Main hand weapon:", mainHand.name, "- damage:", mainHand.damage_dice_count + "d" + mainHand.damage_dice_sides + "+" + mainHand.damage_bonus);
          } else {
            console.log("[Combat] No main hand weapon found");
          }
        }
      } catch (err) {
        console.error("Failed to fetch equipment:", err);
      }
    }
  }, [characterID]);

  // Log equipment when it changes (for debugging)
  useEffect(() => {
    if (equipmentState.length > 0) {
      const mainHand = equipmentState.find(e => e.slot === "main_hand" && e.damage_dice_count > 0);
      if (mainHand) {
        console.log("[Combat] Weapon equipped:", mainHand.name, "dmg:", mainHand.damage_dice_count + "d" + mainHand.damage_dice_sides + "+" + mainHand.damage_bonus);
      }
    }
  }, [equipmentState]);

  const performPlayerAttack = useCallback(
    async (target: CombatTarget) => {
      const dexMod = getDexModifier(characterLevel);
      const { roll, total, isCrit, isFumble } = rollD20(dexMod);
      const targetAC = calculateEnemyAC(target.level ?? 1);

      if (isFumble) {
        addLog("🎲 FUMBLE! Natural 1 — You stumble badly!", "miss");
        return;
      }

      if (total < targetAC && !isCrit) {
        addLog(`🎲 Miss! (d20=${roll} + ${dexMod} DEX = ${total} vs AC ${targetAC})`, "miss");
        return;
      }

      const damage = calculatePlayerDamage(characterStrength, equipmentRef.current);
      const finalDamage = isCrit ? damage * 2 : damage;
      console.log("[Combat] Attack:", target.name, "damage:", damage, "final:", finalDamage, "isCrit:", isCrit, "equipment count:", equipmentRef.current.length);

      const dmgResult = await applyDamage(target.id, finalDamage, characterID);
      if (dmgResult.messages && dmgResult.messages.length > 0) {
        for (const msg of dmgResult.messages) {
          addLog(msg, "system");
        }
      }
      if (dmgResult.defeated && dmgResult.gold_dropped && dmgResult.gold_dropped > 0) {
        addLog(`💰 ${target.name} drops ${dmgResult.gold_dropped} gold!`, "system");
      }

      const newHP = await fetchTargetHP(target.id);
      if (newHP != null) {
        const updated = targetsRef.current.map((t) =>
          t.id === target.id ? { ...t, hp: newHP } : t
        );
        targetsRef.current = updated;
        setTargets(updated);
      }

      if (isCrit) {
        addLog(`⚔ CRITICAL HIT! ${finalDamage} damage!`, "crit");
      } else {
        addLog(`⚔ You hit ${target.name} for ${finalDamage} damage!`, "hit");
      }
    },
    [characterID, characterLevel, characterStrength, addLog, fetchTargetHP, equipmentRef]
  );

  const performPlayerAbility = useCallback(
    async (target: CombatTarget, ability: AbilityData, effect: EffectData) => {
      // Check mana/stamina cost (use manaRef for current mana value)
      if (ability.manaCost && playerManaRef.current < ability.manaCost) {
        addLog(`Not enough mana for ${ability.name}`, "error");
        return;
      }
      if (ability.staminaCost && playerManaRef.current < ability.staminaCost) {
        addLog(`Not enough stamina for ${ability.name}`, "error");
        return;
      }

      // Apply costs - deduct from server and update local state
      if (ability.manaCost) {
        try {
          const result = await adjustMana(characterID, -ability.manaCost);
          setPlayerMana(result.mana);
          setPlayerMaxMana(result.maxMana);
          console.log(`[COMBAT] Mana after deduction: ${result.mana}/${result.maxMana}`);
        } catch (err) {
          console.error(`Failed to deduct mana: ${err}`);
        }
      }

      // Show the cast message first (if available)
      const castMsg = formatCombatMessage(ability, effect, target.name);
      if (castMsg) {
        addLog(castMsg, "system");
      }

      if (effect.effectType === "damage") {
        const dmgResult = await applyDamage(target.id, effect.value, characterID);
        if (dmgResult.messages && dmgResult.messages.length > 0) {
          for (const msg of dmgResult.messages) {
            addLog(msg, "system");
          }
        }
        if (dmgResult.defeated && dmgResult.gold_dropped && dmgResult.gold_dropped > 0) {
          addLog(`💰 ${target.name} drops ${dmgResult.gold_dropped} gold!`, "system");
        }
        // Log the damage message showing ability name and damage dealt
        addLog(`Your ${ability.name} hits ${target.name} for ${effect.value} damage!`, "hit");
      }

      const newHP = await fetchTargetHP(target.id);
      if (newHP != null) {
        const updated = targetsRef.current.map((t) =>
          t.id === target.id ? { ...t, hp: newHP } : t
        );
        targetsRef.current = updated;
        setTargets(updated);
      }
    },
    [characterID, addLog, formatCombatMessage, fetchTargetHP]
  );

  const performEnemyTurn = useCallback(async () => {
    const currentTargets = targetsRef.current;
    const aliveTargets = currentTargets.filter((t) => t.hp > 0);
    if (aliveTargets.length === 0) return;

    const target = aliveTargets[0];
    const enemyDexMod = Math.floor((target.level ?? 1) / 3);
    const { roll, total, isCrit, isFumble } = rollD20(enemyDexMod);
    const playerAC = calculateEnemyAC(characterLevel);

    if (isFumble) {
      addLog(`🎲 ${target.name} FUMBLES! (rolled 1)`, "miss");
      return;
    }

    if (total < playerAC && !isCrit) {
      addLog(`🎲 ${target.name} misses! (d20=${roll} + ${enemyDexMod} = ${total} vs AC ${playerAC})`, "miss");
      return;
    }

    const damage = calculateEnemyDamage(target.level ?? 1, isCrit);
    await applyDamage(characterID, damage);

    const newHP = Math.max(0, playerHPRef.current - damage);
    setPlayerHP(newHP);
    onPlayerHPChange(newHP);

    if (isCrit) {
      addLog(`⚔ ${target.name} critical hit! ${damage} damage!`, "crit");
    } else {
      addLog(`⚔ ${target.name} hits you for ${damage} damage!`, "hit");
    }

    if (newHP <= 0) {
      addLog("☠ You have been defeated!", "defeat");
      setInCombat(false);
      inCombatRef.current = false;
      if (tickRef.current) {
        clearInterval(tickRef.current);
        tickRef.current = null;
      }
      await healCharacter(characterID, initialMaxHP);
      onCombatEnd();
    }
  }, [characterID, characterLevel, initialMaxHP, addLog, onCombatEnd, onPlayerHPChange]);

  const processTick = useCallback(async () => {
    if (!inCombatRef.current) return;

    setRound((r) => r + 1);

    const action = queuedRef.current;
    setQueuedAction(null);

    if (action === "flee") {
      const { success, roll, total } = attemptFlee(characterLevel);
      if (success) {
        addLog(`🏃 Escape successful! (d20=${roll} + ${Math.floor(characterLevel / 2)} = ${total} vs DC 12)`, "flee");
        setInCombat(false);
        inCombatRef.current = false;
        if (tickRef.current) {
          clearInterval(tickRef.current);
          tickRef.current = null;
        }
        onCombatEnd();
        return;
      }
      addLog(`🏃 Escape failed! (d20=${roll} + ${Math.floor(characterLevel / 2)} = ${total} vs DC 12)`, "flee");
    } else if (action === "attack" || action == null) {
      const aliveTargets = targetsRef.current.filter((t) => t.hp > 0);
      if (aliveTargets.length > 0) {
        await performPlayerAttack(aliveTargets[0]);
      }
    } else {
      // Try to find and cast the ability
      const queuedAbility = abilities.find((a) => a.name === action);
      if (queuedAbility && queuedAbility.id > 0) {
        const aliveTargets = targetsRef.current.filter((t) => t.hp > 0);
        if (aliveTargets.length > 0) {
          const target = aliveTargets[0];
          // Find the first damage effect for this ability
          const damageEffect = effects.find(
            (e) => e.effectType === "damage" && e.target === "enemy"
          ) || effects[0];
          if (damageEffect) {
            await performPlayerAbility(target, queuedAbility, damageEffect);
          }
        }
      } else {
        // Unknown action - treat as attack
        const aliveTargets = targetsRef.current.filter((t) => t.hp > 0);
        if (aliveTargets.length > 0) {
          await performPlayerAttack(aliveTargets[0]);
        }
      }
    }

    if (!inCombatRef.current) return;

    await performEnemyTurn();

    const remainingAlive = targetsRef.current.filter((t) => t.hp > 0);
    if (remainingAlive.length === 0 && inCombatRef.current) {
      addLog("✦ All targets defeated!", "system");
      setInCombat(false);
      inCombatRef.current = false;
      if (tickRef.current) {
        clearInterval(tickRef.current);
        tickRef.current = null;
      }
      onCombatEnd();
    }
  }, [characterLevel, addLog, onCombatEnd, performPlayerAttack, performEnemyTurn, abilities, effects, performPlayerAbility]);

  const startCombat = useCallback(
    async (newTargets: CombatTarget[]) => {
      const refreshed = await Promise.all(
        newTargets.map(async (t) => {
          try {
            const status = await getCombatStatus(t.id);
            return { ...t, hp: status.hp, maxHp: status.maxHp };
          } catch {
            return t;
          }
        })
      );

      setTargets(refreshed);
      setCombatLog([]);
      setRound(1);
      setQueuedAction(null);
      setInCombat(true);
      inCombatRef.current = true;
      roundRef.current = 1;

      addLog(`⚔ Combat started with ${refreshed.map((t) => t.name).join(", ")}!`, "system");

      // Fetch equipment for damage calculation - await it before combat starts
      await fetchEquipment();

      if (tickRef.current) clearInterval(tickRef.current);
      tickRef.current = setInterval(() => {
        if (isProcessingRef.current) return;
        isProcessingRef.current = true;
        processTick().catch((err) => console.error("Combat tick failed:", err)).finally(() => {
          isProcessingRef.current = false;
        });
      }, TICK_MS);
    },
    [addLog, processTick, fetchEquipment]
  );

  const queueAction = useCallback(
    (action: string) => {
      if (!inCombatRef.current) return;
      setQueuedAction(action);
      addLog(`⏱ Queued: ${action}`, "queue");
    },
    [addLog]
  );

  const exitCombat = useCallback(() => {
    setInCombat(false);
    inCombatRef.current = false;
    setQueuedAction(null);
    if (tickRef.current) {
      clearInterval(tickRef.current);
      tickRef.current = null;
    }
    onCombatEnd();
  }, [onCombatEnd]);

  useEffect(() => {
    return () => {
      if (tickRef.current) clearInterval(tickRef.current);
      tickRef.current = null;
    };
  }, []);

  return {
    inCombat,
    targets,
    combatLog,
    round,
    queuedAction,
    playerHP,
    playerMana,
    playerMaxMana,
    startCombat,
    queueAction,
    exitCombat,
  };
}
