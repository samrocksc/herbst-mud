import { useState } from "react";
import type { Ability, CharacterSkill, CharacterPanelTab, InventoryItem } from "../lib/types";
import { Button } from "../ui";

type CharacterPanelProps = {
  readonly activeTab: CharacterPanelTab;
  readonly onTabChange: (tab: CharacterPanelTab) => void;
  readonly onClose: () => void;
  readonly skills: readonly CharacterSkill[];
  readonly onSkillSwap: (from: number, to: number) => void;
  readonly inventory: readonly InventoryItem[];
  readonly availableAbilities: readonly Ability[];
  readonly onEquip: (abilityID: number, slot: number) => Promise<void>;
  readonly onUnequip: (slot: number) => Promise<void>;
};

function AbilitiesView({
  skills,
  available,
  onEquip,
  onUnequip,
}: {
  readonly skills: readonly CharacterSkill[];
  readonly available: readonly Ability[];
  readonly onEquip: (abilityID: number, slot: number) => Promise<void>;
  readonly onUnequip: (slot: number) => Promise<void>;
}) {
  const [selectedSlot, setSelectedSlot] = useState<number | null>(null);

  const handleSlotClick = (slot: number) => {
    setSelectedSlot(slot);
  };

  const handleAbilityClick = async (ability: Ability) => {
    if (selectedSlot == null) return;
    await onEquip(ability.id, selectedSlot);
    setSelectedSlot(null);
  };

  return (
    <div className="space-y-2">
      <p className="text-[10px] text-muted uppercase tracking-wider">Equipped</p>
      {skills.map((sk) => (
        <button
          key={sk.slot}
          type="button"
          onClick={() => handleSlotClick(sk.slot)}
          className={`w-full flex items-center justify-between px-2 py-1.5 rounded border transition-colors font-mono ${
            selectedSlot === sk.slot
              ? "border-accent bg-accent/10"
              : "border-border hover:bg-surface-alt"
          }`}
        >
          <span className="text-[10px] text-muted font-bold">Slot {sk.slot}</span>
          <span className="text-xs">{sk.name ?? "Unassigned"}</span>
          {sk.name ? (
            <span
              className="text-[10px] text-muted hover:text-danger ml-2"
              style={{ cursor: "pointer" }}
              onClick={(e) => {
                e.stopPropagation();
                onUnequip(sk.slot);
              }}
            >
              ✕
            </span>
          ) : null}
        </button>
      ))}

      {selectedSlot != null && (
        <>
          <p className="text-[10px] text-muted uppercase tracking-wider mt-3">
            Select ability for Slot {selectedSlot}
          </p>
          {available.length === 0 ? (
            <p className="text-xs text-muted text-center py-2">No abilities available.</p>
          ) : (
            <div className="space-y-1 max-h-40 overflow-y-auto">
              {available.map((ability) => (
                <button
                  key={ability.id}
                  type="button"
                  onClick={() => handleAbilityClick(ability)}
                  className="w-full text-left px-2 py-1.5 rounded border border-border hover:bg-surface-alt transition-colors"
                >
                  <span className="text-xs font-mono block">{ability.name}</span>
                  <span className="text-[10px] text-muted">{ability.description}</span>
                </button>
              ))}
            </div>
          )}
        </>
      )}
    </div>
  );
}

function SkillsView({
  skills,
  onSwap,
}: {
  readonly skills: readonly CharacterSkill[];
  readonly onSwap: (from: number, to: number) => void;
}) {
  const [selectedSlot, setSelectedSlot] = useState<number | null>(null);

  return (
    <div className="space-y-2">
      {skills.map((sk) => {
        const isSelected = selectedSlot === sk.slot;
        return (
          <button
            key={sk.slot}
            type="button"
            onClick={() => {
              if (selectedSlot === null) {
                setSelectedSlot(sk.slot);
              } else if (selectedSlot === sk.slot) {
                setSelectedSlot(null);
              } else {
                onSwap(selectedSlot, sk.slot);
                setSelectedSlot(null);
              }
            }}
            className={`w-full flex items-center justify-between px-2 py-1.5 rounded border transition-colors font-mono ${
              isSelected
                ? "border-accent bg-accent/10"
                : "border-border hover:bg-surface-alt"
            }`}
          >
            <span className="text-[10px] text-muted font-bold">Slot {sk.slot}</span>
            <span className="text-xs">{sk.name ?? "Unassigned"}</span>
          </button>
        );
      })}
    </div>
  );
}

function InventoryView({ inventory }: { readonly inventory: readonly InventoryItem[] }) {
  if (inventory.length === 0) {
    return (
      <p className="text-xs text-muted text-center py-4">
        Your inventory is empty.
      </p>
    );
  }

  return (
    <div className="space-y-1">
      {inventory.map((item) => (
        <div
          key={item.id}
          className="flex items-center justify-between px-2 py-1.5 rounded border border-border"
        >
          <span className="text-xs font-mono">{item.name}</span>
          {item.quantity != null ? (
            <span className="text-[10px] text-muted">x{item.quantity}</span>
          ) : null}
        </div>
      ))}
    </div>
  );
}

export default function CharacterPanel({
  activeTab,
  onTabChange,
  onClose,
  skills,
  onSkillSwap,
  inventory,
  availableAbilities,
  onEquip,
  onUnequip,
}: Readonly<CharacterPanelProps>) {
  return (
    <div className="flex flex-col h-full bg-surface border-l border-border">
      <div className="shrink-0 flex items-center justify-between px-3 py-2 border-b border-border md:border-b-0">
        <div className="flex gap-1">
          <button
            type="button"
            onClick={() => onTabChange("skills")}
            className={`px-2 py-1 text-xs font-mono rounded transition-colors ${
              activeTab === "skills"
                ? "bg-accent text-background"
                : "text-muted hover:text-foreground"
            }`}
          >
            Skills
          </button>
          <button
            type="button"
            onClick={() => onTabChange("abilities")}
            className={`px-2 py-1 text-xs font-mono rounded transition-colors ${
              activeTab === "abilities"
                ? "bg-accent text-background"
                : "text-muted hover:text-foreground"
            }`}
          >
            Abilities
          </button>
          <button
            type="button"
            onClick={() => onTabChange("inventory")}
            className={`px-2 py-1 text-xs font-mono rounded transition-colors ${
              activeTab === "inventory"
                ? "bg-accent text-background"
                : "text-muted hover:text-foreground"
            }`}
          >
            Inventory
          </button>
        </div>
        <Button variant="ghost" size="sm" onClick={onClose}>
          &#x2715;
        </Button>
      </div>

      <div className="flex-1 min-h-0 overflow-y-auto px-3 py-2">
        {activeTab === "skills" ? (
          <SkillsView skills={skills} onSwap={onSkillSwap} />
        ) : activeTab === "abilities" ? (
          <AbilitiesView
            skills={skills}
            available={availableAbilities}
            onEquip={onEquip}
            onUnequip={onUnequip}
          />
        ) : (
          <InventoryView inventory={inventory} />
        )}
      </div>
    </div>
  );
}