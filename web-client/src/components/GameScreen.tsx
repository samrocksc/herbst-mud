import { useEffect, useRef, useState, useCallback } from "react";
import { useTheme } from "../lib/theme";
import { type Character, type Ability, getCharacterAbilities, listClasslessAbilities, equipAbility, unequipAbility } from "../lib/api";
import { type CharacterSkill, type CharacterPanelTab, type InventoryItem, type RoomScreenPayload } from "../lib/types";
import { useMUDSocket } from "../hooks/useMUDSocket";
import { Button } from "../ui";
import { useCombatEngine } from "../hooks/useCombatEngine";
import CombatScreen from "./CombatScreen";
import Scrollback from "./Scrollback";
import RoomScreen from "./RoomScreen";
import HotkeyBar from "./HotkeyBar";
import InputBar from "./InputBar";
import CharacterPanel from "./CharacterPanel";

export type GameScreenProps = {
  worldName: string;
  character: Character;
  onDisconnect: () => void;
  token: string;
};

const HOTKEY_BINDINGS: Record<string, string> = {
  l: "look",
  e: "examine",
  r: "use potion",
};

const INITIAL_SKILLS: readonly CharacterSkill[] = [
  { slot: 1, name: null },
  { slot: 2, name: null },
  { slot: 3, name: null },
  { slot: 4, name: null },
];

const INITIAL_INVENTORY: readonly InventoryItem[] = [];

export default function GameScreen({
  worldName,
  character,
  onDisconnect,
  token,
}: Readonly<GameScreenProps>) {
  const { state, lines, roomScreen, debugLog, connect, send, disconnect, pushLocal } =
    useMUDSocket();
  const { theme, toggle } = useTheme();
  const scrollRef = useRef<HTMLDivElement>(null);

  const [commandHistory, setCommandHistory] = useState<string[]>([]);
  const [historyIndex, setHistoryIndex] = useState(-1);
  const [showDebug, setShowDebug] = useState(false);
  const [expandedRoomId, setExpandedRoomId] = useState<string | null>(null);
  const [panelOpen, setPanelOpen] = useState(false);
  const [panelTab, setPanelTab] = useState<CharacterPanelTab>("inventory");
  const [skills, setSkills] = useState<readonly CharacterSkill[]>(INITIAL_SKILLS);
  const [availableAbilities, setAvailableAbilities] = useState<readonly Ability[]>([]);
  const [combatMode, setCombatMode] = useState(false);
  const combatModeRef = useRef(combatMode);
  combatModeRef.current = combatMode;
  const [pendingTargets, setPendingTargets] = useState<Set<number>>(new Set());
  const [potionCount] = useState(0);

  const {
    inCombat,
    targets: combatTargets,
    combatLog,
    round: combatRound,
    queuedAction,
    playerHP: combatPlayerHP,
    startCombat,
    queueAction,
  } = useCombatEngine({
    characterID: character.id,
    characterLevel: character.level,
    characterStrength: 10, // TODO: fetch from server
    initialHP: character.hitpoints,
    initialMaxHP: character.max_hitpoints,
    skills,
    onLog: (text, kind) => {
      const styleMap: Record<string, "system" | "output" | "error" | "input" | "ping"> = {
        hit: "output", crit: "output", miss: "output", heal: "output",
        system: "system", queue: "system", flee: "output", defeat: "error",
      };
      pushLocal(text, styleMap[kind] ?? "output");
    },
    onCombatEnd: () => {
      pushLocal("Combat ended.", "system");
      handleSubmit("look");
    },
    onPlayerHPChange: (hp) => {
      void hp;
    },
  });

  const loadAbilities = useCallback(async () => {
    try {
      const [charAbilities, classless] = await Promise.all([
        getCharacterAbilities(character.id),
        listClasslessAbilities(),
      ]);
      setAvailableAbilities(classless);
      const newSlots = new Map<number, Ability>();
      const newSkills: CharacterSkill[] = [];
      for (let i = 0; i < 6; i++) {
        const entry = charAbilities.slots[i] as (Ability & { slot?: number }) | null;
        if (entry?.slot != null && entry.name) {
          newSlots.set(entry.slot, entry);
          newSkills.push({ slot: entry.slot, name: entry.name });
        }
      }
      if (newSkills.length > 0) {
        setSkills(newSkills);
      }
    } catch (err) {
      pushLocal(`Failed to load abilities: ${err instanceof Error ? err.message : String(err)}`, "error");
    }
  }, [character.id, pushLocal]);

  const handleEquip = useCallback(async (abilityID: number, slot: number) => {
    try {
      await equipAbility(character.id, abilityID, slot);
      pushLocal(`Equipped ability in slot ${slot}`, "system");
      await loadAbilities();
    } catch (err) {
      pushLocal(`Equip failed: ${err instanceof Error ? err.message : String(err)}`, "error");
    }
  }, [character.id, loadAbilities, pushLocal]);

  const handleUnequip = useCallback(async (slot: number) => {
    try {
      await unequipAbility(character.id, slot);
      pushLocal(`Unequipped slot ${slot}`, "system");
      await loadAbilities();
    } catch (err) {
      pushLocal(`Unequip failed: ${err instanceof Error ? err.message : String(err)}`, "error");
    }
  }, [character.id, loadAbilities, pushLocal]);

  const handleToggleExpand = useCallback((id: string) => {
    setExpandedRoomId((prev) => (prev === id ? null : id));
  }, []);

  const handleTogglePending = useCallback((id: number) => {
    setPendingTargets((prev) => {
      const next = new Set(prev);
      if (next.has(id)) {
        next.delete(id);
      } else {
        next.add(id);
      }
      return next;
    });
  }, []);

  const handleConfirmAttack = useCallback(
    async (char: RoomScreenPayload["characters"][number]) => {
      setPendingTargets((prev) => {
        const next = new Set(prev);
        next.delete(char.id);
        return next;
      });
      await startCombat([{ id: char.id, name: char.name, hp: 0, maxHp: 0 }]);
    },
    [startCombat]
  );

  void worldName;

  useEffect(() => {
    const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
    const { hostname, port } = window.location;
    // Same-origin in production (nginx proxy), port 8080 for local dev
    const wsHost = port === "8080" || port === "5174" || port === "5173" ? `${hostname}:8080` : hostname;
    const url = `${protocol}//${wsHost}/ws?token=${encodeURIComponent(token)}&character_id=${character.id}`;
    connect(url);
    return () => disconnect();
  }, [connect, disconnect, token, character.id]);

  // Load abilities when connected and once after mount
  useEffect(() => {
    if (state === "connected") {
      loadAbilities();
    }
  }, [state, loadAbilities]);

  useEffect(() => {
    if (scrollRef.current) {
      scrollRef.current.scrollTop = scrollRef.current.scrollHeight;
    }
  }, [lines]);

  const handleSubmit = useCallback(
    (text: string) => {
      if (!text.trim()) return;
      pushLocal(`> ${text}`, "input");
      setCommandHistory((prev) => [text, ...prev].slice(0, 50));
      setHistoryIndex(-1);
      send("command", text.trim());
    },
    [send, pushLocal],
  );

  const openPanel = useCallback(
    (tab: CharacterPanelTab) => {
      if (panelOpen && panelTab === tab) {
        setPanelOpen(false);
      } else {
        setPanelTab(tab);
        setPanelOpen(true);
      }
    },
    [panelOpen, panelTab],
  );

  const closePanel = useCallback(() => {
    setPanelOpen(false);
  }, []);

  const handleSkillSwap = useCallback((from: number, to: number) => {
    setSkills((prev) => {
      const arr = prev.map((s) => ({ ...s }));
      const a = arr.find((s) => s.slot === from);
      const b = arr.find((s) => s.slot === to);
      if (a && b) {
        const tmp = a.name;
        a.name = b.name;
        b.name = tmp;
      }
      return arr;
    });
  }, []);

  useEffect(() => {
    const handler = (e: KeyboardEvent) => {
      if (e.target instanceof HTMLInputElement) return;
      const key = e.key.toLowerCase();

      if (inCombat) {
        if (key >= "1" && key <= "4") {
          e.preventDefault();
          const sk = skills.find((s) => s.slot === Number(key));
          queueAction(sk?.name ?? "attack");
          return;
        }
        if (key === "5" || key === "r") {
          e.preventDefault();
          queueAction("use potion");
          return;
        }
        if (key === "f") {
          e.preventDefault();
          queueAction("flee");
          return;
        }
        return; // block all other keys in combat
      }

      if (key === "i") { e.preventDefault(); openPanel("inventory"); return; }
      if (key === "s") { e.preventDefault(); openPanel("skills"); return; }
      if (key === "a") { e.preventDefault(); openPanel("abilities"); return; }
      if (key === "tab") { e.preventDefault(); setCombatMode((v) => !v); return; }
      if (key === "l") { e.preventDefault(); handleSubmit("look"); return; }
      if (key === "e") { e.preventDefault(); handleSubmit("examine"); return; }
      if (HOTKEY_BINDINGS[key]) { e.preventDefault(); handleSubmit(HOTKEY_BINDINGS[key]); }
    };
    window.addEventListener("keydown", handler);
    return () => window.removeEventListener("keydown", handler);
  }, [handleSubmit, openPanel, inCombat, queueAction, skills]);

  const handleHotkey = useCallback(
    (slot: string) => {
      const num = Number(slot);
      if (!Number.isNaN(num)) {
        const ability = skills.find((s) => s.slot === num);
        if (ability?.name) {
          handleSubmit(ability.name);
          return;
        }
      }
      if (HOTKEY_BINDINGS[slot]) handleSubmit(HOTKEY_BINDINGS[slot]);
    },
    [handleSubmit, skills],
  );

  const handleTapExit = useCallback(
    (exit: { direction: string; label: string }) => {
      pushLocal(`Moving ${exit.direction} toward ${exit.label}...`, "system");
      setExpandedRoomId(null);
      handleSubmit(exit.direction);
    },
    [handleSubmit, pushLocal],
  );

  return (
    <div className="flex flex-col h-screen w-full bg-background text-foreground font-mono overflow-hidden">
      <header className="shrink-0 flex items-center justify-between px-3 py-2 border-b border-border bg-surface">
        <div className="flex items-center gap-2 text-xs">
          <span className="font-bold text-accent">{character.name}</span>
          <span className="text-muted">&bull;</span>
          <span className="text-muted">Lv.{character.level}</span>
          <span className="text-muted">&bull;</span>
          <span className="text-danger">
            HP {character.hitpoints}/{character.max_hitpoints}
          </span>
          <span className="text-muted">&bull;</span>
          <span className="text-warning">
            STA {character.stamina}/{character.max_stamina}
          </span>
          <span className="text-muted">&bull;</span>
          <span className="text-info">
            MANA {character.mana}/{character.max_mana}
          </span>
        </div>
        <div className="flex items-center gap-2">
          <Button variant="ghost" size="sm" onClick={() => openPanel("inventory")}>
            I
          </Button>
          <Button variant="ghost" size="sm" onClick={() => openPanel("skills")}>
            S
          </Button>
          <Button variant="ghost" size="sm" onClick={() => openPanel("abilities")}>
            A
          </Button>
          <span
            className="px-2 py-0.5 rounded border border-border text-[10px]"
            style={{
              color:
                state === "connected"
                  ? "var(--mud-success)"
                  : state === "connecting"
                    ? "var(--mud-warning)"
                    : "var(--mud-danger)",
            }}
          >
            {state === "connected"
              ? "\u25cf online"
              : state === "connecting"
                ? "\u25d0 connecting..."
                : "\u25cb offline"}
          </span>
          <Button variant="ghost" size="sm" onClick={() => setCombatMode((v) => !v)} disabled={inCombat} title="Toggle combat mode (Tab)">
            {combatMode ? "⚔️" : "🛡️"}
          </Button>
          <Button variant="ghost" size="sm" onClick={toggle}>
            {theme === "dark" ? "\uD83C\uDF19" : "\u2600\uFE0F"}
          </Button>
          <Button
            variant={showDebug ? "primary" : "ghost"}
            size="sm"
            onClick={() => setShowDebug((v) => !v)}
            title="Toggle debug log"
          >
            {"\uD83E\uDEB2"}
          </Button>
          <Button variant="danger" size="sm" onClick={onDisconnect}>
            Disconnect
          </Button>
        </div>
      </header>

      {/* Main area: mobile gets panel overlay, desktop gets side-by-side */}
      <div className="flex-1 min-h-0 flex relative">
        {/* Scrollback + Room panel column */}
        <div className="flex-1 min-w-0 flex flex-col">
          <div ref={scrollRef} className="flex-1 min-h-0 overflow-y-auto px-3 pt-2">
            <Scrollback lines={lines} />
          </div>

          {inCombat ? (
            <CombatScreen
              round={combatRound}
              targets={combatTargets}
              combatLog={combatLog}
              queuedAction={queuedAction}
              playerHP={combatPlayerHP}
              playerMaxHP={character.max_hitpoints}
              playerStamina={character.stamina}
              playerMaxStamina={character.max_stamina}
              playerMana={character.mana}
              playerMaxMana={character.max_mana}
              skills={skills}
              potionCount={potionCount}
              onSkill={(slot) => {
                const sk = skills.find((s) => s.slot === slot);
                queueAction(sk?.name ?? "attack");
              }}
              onPotion={() => queueAction("use potion")}
              onFlee={() => queueAction("flee")}
            />
          ) : roomScreen ? (
            <RoomScreen
              room={roomScreen}
              onTapExit={handleTapExit}
              onCommand={handleSubmit}
              expandedId={expandedRoomId}
              onToggleExpand={handleToggleExpand}
              pendingTargets={pendingTargets}
              onTogglePending={handleTogglePending}
              onConfirmAttack={handleConfirmAttack}
            />
          ) : (
            <div className="shrink-0 bg-surface border-t border-border px-3 py-4 text-center text-xs text-muted">
              {state === "connecting"
                ? "Connecting to world..."
                : state === "connected"
                  ? "Waiting for room data..."
                  : "Disconnected"}
            </div>
          )}
        </div>

        {/* Desktop sidebar panel */}
        {panelOpen && (
          <div className="hidden md:flex w-64 shrink-0 border-l border-border">
            <CharacterPanel
              activeTab={panelTab}
              onTabChange={setPanelTab}
              onClose={closePanel}
              skills={skills}
              onSkillSwap={handleSkillSwap}
              inventory={INITIAL_INVENTORY}
              availableAbilities={availableAbilities}
              onEquip={handleEquip}
              onUnequip={handleUnequip}
            />
          </div>
        )}
      </div>

      {/* Mobile overlay panel */}
      {panelOpen && (
        <div className="flex md:hidden fixed inset-0 z-50">
          <div className="absolute inset-0 bg-black/50" onClick={closePanel} />
          <div className="relative ml-auto w-full max-w-xs bg-surface border-l border-border flex flex-col">
            <CharacterPanel
              activeTab={panelTab}
              onTabChange={setPanelTab}
              onClose={closePanel}
              skills={skills}
              onSkillSwap={handleSkillSwap}
              inventory={INITIAL_INVENTORY}
              availableAbilities={availableAbilities}
              onEquip={handleEquip}
              onUnequip={handleUnequip}
            />
          </div>
        </div>
      )}

      {!inCombat && (
        <HotkeyBar onActivate={handleHotkey} skills={skills} />
      )}

      <InputBar
        onSubmit={handleSubmit}
        history={commandHistory}
        historyIndex={historyIndex}
        setHistoryIndex={setHistoryIndex}
      />

      {showDebug && (
        <div
          className="shrink-0 border-t border-border bg-black/40"
          style={{ maxHeight: "160px" }}
        >
          <div className="flex items-center justify-between px-3 py-1 border-b border-border/50">
            <span className="text-[10px] text-muted uppercase tracking-wider">
              Debug Log ({debugLog.length})
            </span>
            <Button variant="ghost" size="sm" onClick={() => setShowDebug(false)}>
              Close
            </Button>
          </div>
          <div
            className="overflow-y-auto px-3 py-2 space-y-1"
            style={{ maxHeight: "120px" }}
          >
            {debugLog.length === 0 ? (
              <span className="text-[10px] text-muted italic">
                No debug entries yet...
              </span>
            ) : (
              debugLog.map((entry) => (
                <div key={entry.id} className="flex gap-2 text-[10px] leading-tight">
                  <span className="shrink-0 opacity-60">
                    {new Date(entry.timestamp).toLocaleTimeString("en-US", {
                      hour12: false,
                      hour: "2-digit",
                      minute: "2-digit",
                      second: "2-digit",
                    })}
                  </span>
                  <span
                    className={`shrink-0 px-1 rounded font-bold ${
                      entry.direction === "send"
                        ? "bg-success/20 text-success"
                        : entry.direction === "recv"
                          ? "bg-info/20 text-info"
                          : entry.direction === "state"
                            ? "bg-warning/20 text-warning"
                            : entry.direction === "error"
                              ? "bg-danger/20 text-danger"
                              : "bg-muted/20 text-muted"
                    }`}
                  >
                    {entry.direction.toUpperCase()}
                  </span>
                  <span className="text-foreground truncate" title={entry.payload}>
                    {entry.label}
                    {entry.payload ? ` \u2192 ${entry.payload}` : ""}
                  </span>
                </div>
              ))
            )}
          </div>
        </div>
      )}
    </div>
  );
}