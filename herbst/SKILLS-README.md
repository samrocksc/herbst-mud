# 🗡️ Classless Skills System

> **All characters have access to 5 classless combat skills.** These are swappable skills that provide tactical options in combat, regardless of class.

---

## Overview

The skill system provides 5 reusable combat abilities that any character can equip in their 5 skill slots. Skills cost mana/stamina to activate and have cooldowns. During combat, press `1-5` to activate the skill in that slot.

---

## The 5 Classless Skills

| Slot | Skill | Description | Effect | Mana | Stamina |
|------|-------|-------------|--------|------|---------|
| 1 | **Concentrate** | Focus your mind to increase accuracy | `+WIS to hit rolls for 4 rounds` | 10 | 0 |
| 2 | **Haymaker** | A powerful but reckless strike | `+STR damage, -DEX to hit for 1 attack` | 0 | 15 |
| 3 | **Back-off** | Use agility to avoid damage | `Guaranteed dodge vs all attacks this round` | 0 | 25 |
| 4 | **Scream** | Release a berserker cry | `+DEX/STR, -WIS/INT for 2 rounds` | 5 | 10 |
| 5 | **Slap** | A quick stunning strike | `DEX vs CON check: stun for 1 round` | 0 | 12 |

### Skill Details

#### 🎯 Concentrate
- **Duration:** 4 rounds
- **Cooldown:** 8 rounds
- **Best for:** Characters with high WIS who need accuracy
- **Stacking:** Can be reapplied before wearing off

#### 💪 Haymaker
- **Duration:** 1 attack (next attack only)
- **Cooldown:** 6 rounds
- **Best for:** STR-based fighters who can afford the accuracy penalty
- **Note:** The damage bonus is significant but you may miss!

#### 💨 Back-off
- **Duration:** 1 round (immediate)
- **Cooldown:** 10 rounds
- **Best for:** Emergency situations, buying time for heals
- **Note:** Works against ALL attacks from ALL enemies this round

#### 😤 Scream
- **Duration:** 2 rounds
- **Cooldown:** 12 rounds
- **Best for:** Desperation moves or burst damage situations
- **Trade-off:** INT/WIS penalties affect spellcasting and perception

#### 👋 Slap
- **Duration:** 1 round stun on target
- **Cooldown:** 8 rounds
- **Best for:** Interrupting enemy combos or escaping
- **Note:** Target resists with CON. Higher DEX = better success rate

---

## Commands

### View Skills
```
skills              - Show your 5 equipped skills in slots 1-5
skill show          - Same as 'skills'
skill all           - Display all 5 classless skills available
```

### Manage Skills
```
skill equip <name> <slot>  - Equip a skill to slot 1-5
skill swap <s1> <s2>        - Swap skills between two slots
```

### Examples
```
skill equip concentrate 1     # Put Concentrate in slot 1
skill equip haymaker 2        # Put Haymaker in slot 2
skill swap 1 3                # Swap slots 1 and 3
```

---

## In Combat

### During Combat Mode
Once combat starts, your 5 skill slots are mapped to number keys:

| Key | Action |
|-----|--------|
| `1` | Activate skill in slot 1 |
| `2` | Activate skill in slot 2 |
| `3` | Activate skill in slot 3 |
| `4` | Activate skill in slot 4 |
| `5` | Activate skill in slot 5 |

### Combat Log
Skill activations and effects appear in the combat log:
- `🎯 Concentrate! +3 accuracy for 4 rounds`
- `💪 Haymaker! +5 damage, -2 accuracy this attack`
- `💨 Back-off! Dodging all attacks this round!`
- `😤 SCREAM! +2 DEX/STR, -2 WIS/INT for 2 rounds`
- `👋 SLAP! Target is stunned for 1 round!`

---

## Strategy Tips

### Skill Combos

1. **Open with Scream → Spam Haymaker**
   - Scream buffs your STR, Haymaker uses it
   - Trade INT/WIS for massive damage output

2. **Concentrate → Back-off (Defensive)**
   - Use Concentrate to improve accuracy
   - Back-off when things get dicey
   - Your next attacks will hit hard

3. **Slap → Haymaker (Control)**
   - Stun prevents enemy action
   - Free Haymaker with no retaliation risk

4. **Back-off → Potion → Attack**
   - Dodge round gives time to heal
   - Re-engage with full resources

### Stat Synergies

| Skill | Primary Stat | Works Best With |
|-------|--------------|-----------------|
| Concentrate | WIS | Clerics, Mages, Rangers |
| Haymaker | STR | Warriors, Barbarians, Fighters |
| Back-off | DEX | Rogues, Rangers, Monks |
| Scream | Any (uses CON for magnitude) | Tanky builds |
| Slap | DEX | Rogues, fast attackers |

---

## Implementation Notes

### Server Endpoints
```
GET    /characters/:id/classless-skills        - Get all skills for character
POST   /characters/:id/classless-skills        - Equip a skill to a slot
PUT    /characters/:id/classless-skills/swap     - Swap skills between slots
```

### Data Flow
1. Skills are stored in `ClasslessSkills` array (client-side const)
2. Character's equipped slots fetched from server on login
3. `combatSkills` struct tracks:
   - `EquippedSkill[5]` - Current loadout
   - `ActiveEffects[]` - Currently active skill effects
   - `Cooldowns map[int]int` - Tick-based cooldown tracking

---

## Future Enhancements

- [ ] Skill progression (use skills to level them up)
- [ ] Skill-specific talents that enhance each ability
- [ ] New skills added through quest rewards
- [ ] Skill cooldown reduction via equipment
- [ ] Combo system for chaining skills

---

**Version:** 1.0  
**Last Updated:** 2026-03-31
