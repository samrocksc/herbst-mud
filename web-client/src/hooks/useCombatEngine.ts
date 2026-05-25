import { useCallback, useEffect, useRef, useState } from "react";
import type { CombatTarget, CombatLogEntry, CharacterSkill } from "../lib/types";
import { getCombatStatus, applyDamage, healCharacter } from "../lib/api";
import {
  rollD20,
  calculatePlayerDamage,
  calculateEnemyDamage,
  calculateEnemyAC,
  getDexModifier,
  attemptFlee,
} from "../lib/combat";

const TICK_MS = 1500;

type Props = {
  characterID: number;
  characterLevel: number;
  characterStrength: number;
  initialHP: number;
  initialMaxHP: number;
  skills: readonly CharacterSkill[];
  onLog: (text: string, kind: CombatLogEntry["kind"]) => void;
  onCombatEnd: () => void;
  onPlayerHPChange: (hp: number) => void;
};

export function useCombatEngine({
  characterID,
  characterLevel,
  characterStrength,
  initialHP,
  initialMaxHP,
  onLog,
  onCombatEnd,
  onPlayerHPChange,
}: Props) {
  const [inCombat, setInCombat] = useState(false);
  const [targets, setTargets] = useState<CombatTarget[]>([]);
  const [combatLog, setCombatLog] = useState<CombatLogEntry[]>([]);
  const [round, setRound] = useState(0);
  const [queuedAction, setQueuedAction] = useState<string | null>(null);
  const [playerHP, setPlayerHP] = useState(initialHP);

  const targetsRef = useRef(targets);
  const playerHPRef = useRef(playerHP);
  const inCombatRef = useRef(inCombat);
  const roundRef = useRef(round);
  const queuedRef = useRef(queuedAction);
  const tickRef = useRef<ReturnType<typeof setInterval> | null>(null);
  const isProcessingRef = useRef(false);

  targetsRef.current = targets;
  playerHPRef.current = playerHP;
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

  const fetchTargetHP = useCallback(async (targetID: number) => {
    try {
      const status = await getCombatStatus(targetID);
      return status.hp;
    } catch {
      return null;
    }
  }, []);

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

      const damage = calculatePlayerDamage(characterStrength);
      const finalDamage = isCrit ? damage * 2 : damage;

      await applyDamage(target.id, finalDamage, characterID);

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
    [characterID, characterLevel, characterStrength, addLog, fetchTargetHP]
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
  }, [characterLevel, addLog, onCombatEnd, performPlayerAttack, performEnemyTurn]);

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

      if (tickRef.current) clearInterval(tickRef.current);
      tickRef.current = setInterval(() => {
        if (isProcessingRef.current) return;
        isProcessingRef.current = true;
        processTick().catch((err) => console.error("Combat tick failed:", err)).finally(() => {
          isProcessingRef.current = false;
        });
      }, TICK_MS);
    },
    [addLog, processTick]
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
    startCombat,
    queueAction,
    exitCombat,
  };
}