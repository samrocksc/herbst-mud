# Skill/Talent Swap UI Design

*Designer: Mikey 🐢*
*Date: 2026-03-18*

---

## Overview

Terminal UI mockups for swapping equipped skills and talents in herbst-mud.

---

## Screen 1: View Available Skills

```
╔═══════════════════════════════════════════════════════════╗
║                    ⚔️ YOUR SKILLS ⚔️                        ║
╠═══════════════════════════════════════════════════════════╣
║ WEAPON MASTERY (2/4 equipped)                              ║
║ ─────────────────────────────────────────────────────────  ║
║ 🟢 [1] Sword Combat      ★★★☆☆  (lvl 3)  [EQUIPPED]       ║
║ 🔵 [2] Axe Mastery       ★★☆☆☆  (lvl 2)                    ║
║ 🔵 [3] Staff Magic       ★★★★★  (lvl 5)                    ║
║ 🔵 [4] Bow Proficiency   ★☆☆☆☆  (lvl 1)  [EQUIPPED]       ║
║                                                           ║
║ MAGIC ARCANE (1/2 equipped)                               ║
║ ─────────────────────────────────────────────────────────  ║
║ 🟢 [5] Fire Bolt         ★★★☆☆  (lvl 3)  [EQUIPPED]       ║
║ 🔵 [6] Ice Shard         ★★☆☆☆  (lvl 2)                    ║
║                                                           ║
╠═══════════════════════════════════════════════════════════╣
║ Commands: [1-6] equip | [S]kill details | [Q]uit            ║
╚═══════════════════════════════════════════════════════════╝
```

---

## Screen 2: View Available Talents

```
╔═══════════════════════════════════════════════════════════╗
║                    ✨ YOUR TALENTS ✨                      ║
╠═══════════════════════════════════════════════════════════╣
║ COMBAT ARSENAL (2/4 slots)                                ║
║ ─────────────────────────────────────────────────────────  ║
║ 🟢 [1] Power Strike    [1]  ★★★☆☆  (lvl 3) [EQUIPPED]    ║
║ 🟢 [2] Shield Block    [2]  ★★☆☆☆  (lvl 2) [EQUIPPED]    ║
║ 🔵 [3] Dual Wield      [3]  ★★★★★  (lvl 5)               ║
║ 🔵 [4] Battle Cry      [4]  ★☆☆☆☆  (lvl 1)               ║
║ 🔵 [5] Parry          [--]  ★★★☆☆  (lvl 3)               ║
║                                                           ║
║ BUFF DEPLOYMENT (0/2 slots)                              ║
║ ─────────────────────────────────────────────────────────  ║
║ 🔵 [6] Healing Word    [H]  ★★★★☆  (lvl 4)               ║
║ 🔵 [7] Stone Skin     [K]  ★☆☆☆☆  (lvl 1)               ║
║                                                           ║
╠═══════════════════════════════════════════════════════════╣
║ Hotbar: 1️⃣Power Strike  2️⃣Shield Block  ⬜  ⬜           ║
║ Commands: [1-7] equip | [T]alent details | [Q]uit         ║
╚═══════════════════════════════════════════════════════════╝
```

---

## Screen 3: Equip Success Feedback

```
════════════════════════════════════════════════════════════
  ✅ SWAP SUCCESS!
  
  ⚔️ Removed: Bow Proficiency
  ⚔️ Equipped: Axe Mastery
  
  Your weapon loadout: Sword | Axe | Fire Bolt | Bow
════════════════════════════════════════════════════════════
```

---

## Screen 4: Invalid Swap Error

```
════════════════════════════════════════════════════════════
  ❌ CANNOT EQUIP
  
  Reason: You don't have the required skill!
  
  "Dual Wield requires: Sword Combat lvl 3 (you have lvl 2)"
  "Train sword combat to unlock this talent."
════════════════════════════════════════════════════════════
```

---

## In-Combat Mode (Simplified)

```
╔═══════════════════════════════════════════════════════════╗
║ Combat: Giant Rat [HP: ████████░░ 80%]                    ║
╠═══════════════════════════════════════════════════════════╣
║ Your talents: [1]Power Strike  [2]Shield  [3]🟣 [4]🟣    ║
║                                                               ║
║ > You see a rusty dagger here.                              ║
║ Commands: [1-4] use talent | [S]kills | [Q]uit combat      ║
╚═══════════════════════════════════════════════════════════╝
```

---

## UX Principles

1. **Clear State** - Equipped items always highlighted green 🟢
2. **Context Awareness** - Combat mode hides non-essential UI
3. **Helpful Errors** - Tell user WHY a swap failed
4. **Visual Hotbar** - Always show current talent bindings
5. **Color Coding**: 
   - 🟢 Green = Equipped/Available
   - 🔵 Blue = Available but not equipped
   - 🟣 Purple = In Progress/Active

---

*Design by Mikey - Make it look good, play even better! 🐢🍕*