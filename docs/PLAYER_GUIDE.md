# Player Guide

> 🔵 Last Updated: 2026-04-05

Welcome to **Herbst MUD**! This guide will help you get started playing the game.

## Table of Contents

- [Getting Started](#getting-started)
- [Connecting to the Game](#connecting-to-the-game)
- [Basic Commands](#basic-commands)
- [Movement](#movement)
- [Combat](#combat)
- [Skills System](#skills-system)
- [Talents](#talents)
- [Equipment](#equipment)
- [Quests](#quests)
- [Tips & Tricks](#tips--tricks)

---

## Getting Started

### What is a MUD?

**MUD** stands for **Multi-User Dungeon** (or Dimension). It's a text-based multiplayer game where you explore virtual worlds, solve puzzles, defeat enemies, and interact with other players - all through text commands!

### Creating Your Account

1. The admin panel manages user accounts
2. Create a user account with an email and password
3. Create a character with stats and appearance
4. Log into the game via SSH

---

## Connecting to the Game

### Via SSH

```bash
ssh -p 4444 game.herbst-mud.com
```

Or using your game's specific URL:

```bash
ssh -p 4444 your-droplet-ip
```

**Default Port:** 4444

### Login Process

1. Connect via SSH
2. Enter your username (character name)
3. Enter your password
4. You'll appear in the starting room!

---

## Basic Commands

### Navigation Commands

| Command | Short | Description |
|---------|-------|-------------|
| `look` | `l` | Look around the current room |
| `exits` | - | Show available exits |
| `peer <direction>` | - | Look into adjacent room |

### Information Commands

| Command | Description |
|---------|-------------|
| `score` | Show character stats, HP, level |
| `inventory` | `i` | Show items you're carrying |
| `equipment` | `eq` | Show worn equipment |
| `who` | List players online |

### Interaction Commands

| Command | Description |
|---------|-------------|
| `get <item>` | Pick up an item |
| `drop <item>` | Drop an item |
| `examine <thing>` | Look closely at something |
| `use <item>` | Use a consumable item |
### Communication Commands

| Command | Description |
|---------|-------------|
| `say <message>` | Speak to everyone in the room |
| `tell <player> <message>` | Send a private message |
| `shout <message>` | Yell to adjacent rooms |

---

## Movement

### Direction Commands

Move in any cardinal direction:

| Command | Short | Direction |
|---------|-------|-----------|
| `north` | `n` | Move north |
| `south` | `s` | Move south |
| `east` | `e` | Move east |
| `west` | `w` | Move west |
| `northeast` | `ne` | Move northeast |
| `southeast` | `se` | Move southeast |
| `southwest` | `sw` | Move southwest |
| `northwest` | `nw` | Move northwest |
| `up` | `u` | Move up |
| `down` | `d` | Move down |

### Peeking

Before entering a room, you can peek into it:

```
peer north
```

This shows you the room name and description without entering.

---

## Combat

### Combat Basics

Combat in Herbst MUD is **tick-based**:
- Each "tick" is approximately 1.5 seconds
- Characters and NPCs take turns based on speed
- Combat continues until someone flees, dies, or is defeated

### Engaging in Combat

```
attack <target>
```

Or shorthand:

```
a <target>
```

Example:
```
attack goblin
a goblin
```

### Combat Status

After entering combat:
- You'll see attack notifications each tick
- Your HP will decrease when hit
- NPCs also take damage
- Combat ends when HP reaches 0

### Combat Commands

During combat, you have several options:

| Key | Command | Effect |
|-----|---------|--------|
| `1` | Concentrate | Focus your attack for bonus damage |
| `2` | Haymaker | Powerful strike with momentum |
| `3` | Back-off | Create distance, defensive maneuver |
| `4` | Scream | Intimidate your opponent |
| `r` | Potion | Use a healing potion from inventory |
| `q` | Flee | Attempt to escape combat |

**Note:** These classless skills are available to all characters regardless of class.

### Combat Display

During combat you'll see:
```
⚔️ In Combat vs. Aragorn (85/100 HP)

[1] Concentrate  [2] Haymaker
[3] Back-off     [4] Scream
[R] Use Potion   [Q] Flee

Aragorn uses Nature's Blessing! Heals for 5 HP
You hit Aragorn for 12 damage!
```

### NPC Special Abilities

Some NPCs have special skills they can use in combat:

| NPC | Ability | Effect |
|-----|---------|--------|
| Aragorn | Nature's Blessing | Heals 5% of max HP (druid magic) |

NPCs choose to use skills based on their health and timing. Watch for visual cues like 🌿 to identify special abilities!

### Fleeing Combat

To escape from combat:

```
flee <direction>
```

Example:
```
flee north
flee n
```

**Note:** Fleeing has a chance to fail based on your stats!

### Corpse System

When creatures die, they leave corpses:

| Command | Description |
|---------|-------------|
| `search` | Look for hidden or missed items in room/corpse |
| `search corpse` | Look for items on a corpse |
| `get <item> corpse` | Take item from corpse |

---

## Skills System

### Overview

**Skills** are permanent abilities your character learns. They improve over time and remain forever.

### Skill Types

| Skill | Description |
|-------|-------------|
| Weapon Skills | Sword, axe, dagger mastery |
| Defense | Dodge, parry, block |
| Magic | Spell casting, mana pool |
| Utility | Crafting, gathering |

### Viewing Skills

```
skills
```

### Improving Skills

Skills improve through use:
- Hit enemies with swords → Sword skill increases
- Dodge attacks → Dodge skill increases
- Cast spells → Magic skills increase

---

## Talents

### Overview

**Talents** are special abilities that can be equipped and swapped. You can have up to 6 talents equipped at once.

### Types of Talents

| Talent | Effect |
|--------|--------|
| `Power Attack` | Next attack deals double damage |
| `Quick Heal` | Regenerate HP over time |
| `Mighty Blow` | High damage attack with cooldown |
| `Evasion` | Temporarily increase dodge chance |

### Viewing Talents

```
talents
```

Shows equipped talents and available slots.

### Using Talents

Talents are used by slot number or name during combat:

```
talent 1
# or
talent power
```

### Swapping Talents

Out of combat, you can swap talents:

```
talent swap <slot> <talent_id>
```

### Talent Cooldowns

- Each talent has a cooldown period
- Check remaining cooldown with `talents` command
- Cooldowns reset when leaving combat

---

## Equipment

### Equipment Slots

Your character has multiple equipment slots:

| Slot | Description |
|------|-------------|
| `head` | Helmets, hats, crowns |
| `body` | Armor, robes, clothing |
| `hands` | Weapons, gloves |
| `feet` | Boots, shoes |

### Equipping Items

```
wear <item>
equip <item>
```

Example:
```
wear leather armor
wield iron sword
```

### Unequipping Items

```
remove <item>
unequip <item>
```

### Equipment Statistics

Equipment affects your character:
- **Armor Class (AC):** Reduces incoming damage
- **Damage:** Increases attack damage
- **Stat Bonuses:** +STR, +DEX, etc.

View equipment stats:

```
score
equipment
```

---

## Quests

### Starting Quests

Talk to NPCs to start quests:

```
talk gizmo
```

### Quest Objectives

Quests have specific objectives:
- Visit locations
- Defeat enemies
- Collect items
- Talk to NPCs

### Checking Quest Status

```
quests
```

Shows active quests and progress.

### Completing Quests

Return to the quest giver when complete:

```
talk <npc>
```

### Quest Rewards

Quests reward:
- **Experience Points (XP)** - Level up!
- **Gold** - Buy items
- **Items** - Special equipment
- **Reputation** - NPC standing

---

## Character Progression

### Leveling Up

Gain XP through:
- Defeating enemies
- Completing quests
- Exploring new areas

When you have enough XP:
- You'll level up automatically
- Max HP increases
- Stats may improve
- New abilities unlock

### Statistics

| Stat | Affects |
|------|---------|
| **STR** | Physical damage |
| **DEX** | Dodge, hit chance, order |
| **CON** | HP, regeneration |
| **INT** | Magic power |
| **WIS** | Mana, magic resistance |
| **CHA** | NPC interactions |

View stats: `score`

---

## Interactive NPCs

### Gizmo - The Helpful Robot

Located in the Junkyard, Gizmo teaches new players:

```
talk gizmo
```

**Services:**
- Tutorial quest ("The Beginning")
- Information about the game
- Heals injured players

### Other NPCs

Look for NPCs throughout the world:

```
look
```

Shows NPCs in the current room.

---

## Tips & Tricks

### New Player Tips

1. **Talk to Gizmo first** - Start in the Junkyard and talk to Gizmo for the tutorial
2. **Check your score often** - Monitor HP, especially in combat
3. **Use peer before entering** rooms - Check for dangerous enemies
4. **Equip starting gear** - Don't forget to equip your starting items
5. **Search corpses** - Enemies often drop useful items

### Combat Tips

1. **Monitor HP** - Flee before dying (at 20% HP recommended)
2. **Use talents wisely** - Save powerful talents for tough enemies
3. **Equip appropriate weapons** - Different enemies weak to different damage types
4. **Passive healing** - Characters regenerate HP slowly over time

### Exploration Tips

1. **Map mentally** - Pay attention to room connections
2. **Visit new rooms** - New areas give XP
3. **Look at everything** - `examine` reveals hidden details
4. **Search rooms** - Hidden items exist!

### Social Tips

1. **Be polite** - Other players are people too
2. **Help newbies** - Remember when you were new
3. **Form parties** - Fight stronger enemies together
4. **Trade items** - Communication economy

---

## Quick Reference Card

```
┌─────────────────────────────────────────────────┐
│  HERBST MUD - Quick Reference                   │
├─────────────────────────────────────────────────┤
│  MOVEMENT    |  COMBAT       |  INVENTORY       │
│  n s e w     |  attack       |  inventory (i)   │
│  ne se sw nw |  a            |  get <item>      │
│  u d         |  flee <dir>   |  drop <item>     │
│  peer <dir>  |               |  wear <item>     │
│              |  TALENTS      |  remove <item>   │
│  INFO        |  talents      |                  │
│  look (l)    |  talent <n>   |  COMMUNICATION   │
│  score       |               |  say             │
│  exits       |  NPC          |  tell <who>      │
│              |  talk <npc>   |  shout           │
│  QUESTS      |               |                  │
│  quests      |  OTHER        |                  │
│              |  search       |                  │
└─────────────────────────────────────────────────┘
```

---

## Getting Help

### In-Game Help

```
help
help <topic>
```

### Admin Panel

Access admin features via web browser at your game's admin URL.

### Troubleshooting

| Problem | Solution |
|---------|----------|
| Can't connect | Check SSH port 4444, verify server address |
| Wrong password | Reset via admin panel |
| Stuck in combat | Flee or call for help |
| Lost in map | Use `look` and `exits` frequently |

---

🔵 Document version: 2026-04-04

**Happy Adventuring!** 🎮
