import { useCallback, useEffect, useState } from "react";
import { type Character, type EquipmentItem, type Race, getCharacterEquipment, equipItem, unequipItem, listRaces } from "../lib/api";
import { Button } from "../ui";

const DEFAULT_SLOTS = ["head", "body", "hands", "legs", "feet", "main_hand"];

const SLOT_LABELS: Record<string, string> = {
  head: "Head",
  body: "Body",
  hands: "Hands",
  legs: "Legs",
  feet: "Feet",
  main_hand: "Main Hand",
  off_hand: "Off Hand",
  ring: "Ring",
  neck: "Neck",
  back: "Back",
  waist: "Waist",
  wrists: "Wrists",
  shell: "Shell",
};

const RARITY_COLORS: Record<string, string> = {
  common: "var(--mud-muted)",
  uncommon: "var(--mud-success)",
  rare: "var(--mud-info)",
  epic: "var(--mud-accent)",
  legendary: "var(--mud-warning)",
};

function slotLabel(slot: string): string {
  return SLOT_LABELS[slot] || slot.replace(/_/g, " ").replace(/\b\w/g, (c) => c.toUpperCase());
}

function rarityColor(rarity: string): string {
  return RARITY_COLORS[rarity?.toLowerCase()] || "var(--mud-muted)";
}

type EquipmentScreenProps = {
  readonly character: Character;
  readonly onClose: () => void;
};

export default function EquipmentScreen({ character, onClose }: Readonly<EquipmentScreenProps>) {
  const [equipped, setEquipped] = useState<readonly EquipmentItem[]>([]);
  const [inventory, setInventory] = useState<readonly EquipmentItem[]>([]);
  const [slots, setSlots] = useState<readonly string[]>(DEFAULT_SLOTS);
  const [selectedSlot, setSelectedSlot] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const load = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const [items, races] = await Promise.all([
        getCharacterEquipment(character.id),
        listRaces(),
      ]);
      const race = races.find((r: Race) => r.name === character.race);
      const raceSlots = race?.equipment_slots?.length ? race.equipment_slots : DEFAULT_SLOTS;
      setSlots(raceSlots);
      setEquipped(items.filter((i) => i.isEquipped));
      setInventory(items.filter((i) => !i.isEquipped));
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load equipment");
    } finally {
      setLoading(false);
    }
  }, [character.id, character.race]);

  useEffect(() => {
    void load();
  }, [load]);

  useEffect(() => {
    const handler = (e: KeyboardEvent) => {
      if (e.key === "Escape") onClose();
    };
    window.addEventListener("keydown", handler);
    return () => window.removeEventListener("keydown", handler);
  }, [onClose]);

  const handleEquip = async (item: EquipmentItem) => {
    try {
      await equipItem(item.id, character.id);
      await load();
      setSelectedSlot(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Equip failed");
    }
  };

  const handleUnequip = async (item: EquipmentItem) => {
    try {
      await unequipItem(item.id, character.id);
      await load();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Unequip failed");
    }
  };

  const handleSlotClick = (slot: string) => {
    setSelectedSlot((prev) => (prev === slot ? null : slot));
  };

  const matchingInventory = selectedSlot
    ? inventory.filter((i) => i.slot === selectedSlot)
    : [];

  return (
    <div className="fixed inset-0 z-50 bg-background flex flex-col">
      {/* Header */}
      <header className="shrink-0 flex items-center justify-between px-4 py-3 border-b border-border bg-surface">
        <div className="flex items-center gap-2 text-xs font-mono">
          <span className="font-bold text-accent">{character.name}</span>
          <span className="text-muted">Lv.{character.level}</span>
          <span className="text-muted">{character.race}</span>
          <span className="text-muted">{character.class}</span>
        </div>
        <Button variant="ghost" size="sm" onClick={onClose}>
          &#x2715;
        </Button>
      </header>

      {/* Error Banner */}
      {error && (
        <div className="shrink-0 bg-danger/20 text-danger text-xs px-4 py-2 font-mono">
          {error}
        </div>
      )}

      {/* Body */}
      <div className="flex-1 min-h-0 overflow-y-auto p-4">
        {loading ? (
          <p className="text-xs text-muted text-center py-8 font-mono">Loading equipment...</p>
        ) : (
          <div className="space-y-4">
            <p className="text-[10px] text-muted uppercase tracking-wider font-mono">Equipment Slots</p>

            <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-3">
              {slots.map((slot) => {
                const item = equipped.find((i) => i.slot === slot);
                const isSelected = selectedSlot === slot;
                return (
                  <button
                    key={slot}
                    type="button"
                    onClick={() => handleSlotClick(slot)}
                    className={`relative text-left rounded border p-3 transition-colors font-mono ${
                      isSelected
                        ? "border-accent bg-accent/10"
                        : "border-border hover:bg-surface-alt"
                    }`}
                  >
                    <div className="flex items-center justify-between mb-1">
                      <span className="text-[10px] text-muted uppercase tracking-wider">
                        {slotLabel(slot)}
                      </span>
                      {item && (
                        <span
                          className="text-[10px] hover:text-danger"
                          style={{ cursor: "pointer" }}
                          onClick={(e) => {
                            e.stopPropagation();
                            void handleUnequip(item);
                          }}
                          title="Unequip"
                        >
                          &#x2715;
                        </span>
                      )}
                    </div>
                    {item ? (
                      <div className="space-y-0.5">
                        <div className="flex items-center gap-1.5">
                          <span
                            className="inline-block w-2 h-2 rounded-full"
                            style={{ backgroundColor: rarityColor(item.rarity) }}
                          />
                          <span className="text-xs font-bold">{item.name}</span>
                        </div>
                        <span className="text-[10px] text-muted">
                          Lv.{item.level} {item.rarity}
                        </span>
                      </div>
                    ) : (
                      <span className="text-xs text-muted italic">Empty</span>
                    )}
                  </button>
                );
              })}
            </div>

            {/* Expanded item list */}
            {selectedSlot != null && (
              <div className="border border-border rounded bg-surface p-3 space-y-2">
                <p className="text-[10px] text-muted uppercase tracking-wider font-mono">
                  Available for {slotLabel(selectedSlot)}
                </p>
                {matchingInventory.length === 0 ? (
                  <p className="text-xs text-muted italic font-mono">
                    No items available for this slot.
                  </p>
                ) : (
                  <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-2">
                    {matchingInventory.map((item) => (
                      <button
                        key={item.id}
                        type="button"
                        onClick={() => void handleEquip(item)}
                        className="text-left rounded border border-border hover:bg-surface-alt transition-colors p-2 font-mono"
                      >
                        <div className="flex items-center gap-1.5">
                          <span
                            className="inline-block w-2 h-2 rounded-full"
                            style={{ backgroundColor: rarityColor(item.rarity) }}
                          />
                          <span className="text-xs font-bold">{item.name}</span>
                        </div>
                        <span className="text-[10px] text-muted">
                          Lv.{item.level} {item.weight}wt {item.rarity}
                        </span>
                      </button>
                    ))}
                  </div>
                )}
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  );
}
