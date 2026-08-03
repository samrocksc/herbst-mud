import type { CombatTarget, CombatLogEntry, CharacterSkill } from "../lib/types";
import CombatTargetList from "./CombatTargetList";
import CombatVitals from "./CombatVitals";
import CombatActionBar from "./CombatActionBar";
import CombatLog from "./CombatLog";

type Props = {
  readonly round: number;
  readonly targets: readonly CombatTarget[];
  readonly combatLog: readonly CombatLogEntry[];
  readonly queuedAction: string | null;
  readonly playerHP: number;
  readonly playerMaxHP: number;
  readonly playerStamina: number;
  readonly playerMaxStamina: number;
  readonly playerMana: number;
  readonly playerMaxMana: number;
  readonly skills: readonly CharacterSkill[];
  readonly potionCount: number;
  readonly onSkill: (slot: number) => void;
  readonly onPotion: () => void;
  readonly onFlee: () => void;
};

export default function CombatScreen({
  round,
  targets,
  combatLog,
  queuedAction,
  playerHP,
  playerMaxHP,
  playerStamina,
  playerMaxStamina,
  playerMana,
  playerMaxMana,
  skills,
  potionCount,
  onSkill,
  onPotion,
  onFlee,
}: Props) {
  // console.log is disabled in production - enable for debugging
  // console.log(`[CombatScreen] Render - playerMana=${playerMana}, playerMaxMana=${playerMaxMana}`);
  return (
    <div className="shrink-0 bg-surface border-t border-border flex flex-col" style={{ maxHeight: "60vh" }}>
      {/* Header */}
      <div className="shrink-0 flex items-center justify-between px-3 py-2 border-b border-border bg-danger/5">
        <span className="text-xs font-bold text-danger uppercase tracking-wider">
          ⚔ Combat — Round {round}
        </span>
        {queuedAction && (
          <span className="text-[10px] text-accent">Queued: {queuedAction}</span>
        )}
      </div>

      {/* Sticky Vitals/Target Section - Fixed at top, above scrollable combat log */}
      <div className="shrink-0 px-3 py-3 space-y-4 border-b border-border/50 bg-surface">
        <CombatTargetList targets={targets} />
        <CombatVitals
          hp={playerHP}
          maxHp={playerMaxHP}
          stamina={playerStamina}
          maxStamina={playerMaxStamina}
          mana={playerMana}
          maxMana={playerMaxMana}
        />
      </div>

      {/* Scrollable Combat Log - Auto-scrolls to bottom */}
      <CombatLog entries={combatLog} />

      {/* Action Bar */}
      <CombatActionBar
        skills={skills}
        potionCount={potionCount}
        onSkill={onSkill}
        onPotion={onPotion}
        onFlee={onFlee}
        queuedAction={queuedAction}
      />
    </div>
  );
}
