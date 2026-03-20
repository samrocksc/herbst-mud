Feature: Combat End (Victory/Defeat Screens)
  As a player
  I want clear combat end screens with rewards or consequences
  So that combat feels meaningful and rewarding

  Background:
    Given a combat is active
    And the player is engaged with an "Old Scrap" enemy

  Scenario: Victory screen shows on enemy defeat
    Given the enemy's HP reaches 0
    When combat ends
    Then the victory screen is displayed
    And the title shows "VICTORY!"
    And the enemy's name is shown: "Old Scrap has been defeated!"

  Scenario: Victory screen shows loot
    Given the enemy drops "scrap_machete" and 23 coins
    When combat ends in victory
    Then the loot section shows:
      | item          | type    |
      | scrap_machete | weapon  |
      | 23 coins      | misc    |

  Scenario: Victory screen shows XP gained
    Given defeating "Old Scrap" grants 150 XP
    When combat ends in victory
    Then the XP section shows 150 XP gained

  Scenario: Victory screen shows skill ups
    Given the player used blades and brawling during combat
    And blades gained 2 XP, brawling gained 1 XP
    When combat ends in victory
    Then the skill up section shows:
      | skill     | gained | new level |
      | blades    | +2     | 47        |
      | brawling  | +1     | 15        |

  Scenario: Defeat screen shows on player death
    Given the player's HP reaches 0
    When combat ends
    Then the defeat screen is displayed
    And the title shows "DEFEAT"

  Scenario: Defeat shows XP loss
    Given the player loses 10% of current level progress
    When combat ends in defeat
    Then the consequences show XP loss amount

  Scenario: Defeat drops random item
    Given the player has items in inventory
    When combat ends in defeat
    Then 1 random item is dropped
    And the dropped item is shown in consequences

  Scenario: Defeat respawns player at last save point
    Given the player's last save point is "Foggy Gate"
    When combat ends in defeat
    Then the player respawns at Foggy Gate
    And the respawn location is shown

  Scenario: Equipment left in zone after defeat
    Given the player has equipment equipped
    When combat ends in defeat
    Then the equipment is left at the combat location
    And the player can reclaim within 5 minutes

  Scenario: Multiple enemy victory shows all loot
    Given 3 enemies are defeated
    And each drops different loot
    When combat ends in victory
    Then all dropped items are shown
    And XP from all enemies is summed

  Scenario: Loot rarity shows correct color coding
    Given enemy drops items of different rarities
    When loot is displayed
    Then common items show in default color
    And rare items show in blue/cyan
    And legendary items show in gold/orange

  Scenario: Victory XP calculation
    Given enemy XP values and bonuses
    When XP is calculated for victory
    Then the total matches expected sum

  Scenario: Defeat respawn timer for equipment
    Given the player dies and leaves equipment
    When 5 minutes pass without reclaim
    Then the equipment is permanently lost

  Scenario: Victory shows combat duration
    Given the combat lasted 8 ticks
    When combat ends in victory
    Then the summary shows "Combat time: 8 ticks"

  Scenario: Defeat penalty reduces level progress
    Given player has 450/1000 XP in current level
    When player is defeated
    Then 10% of current level progress is lost
    And player loses 45 XP

  Scenario: Quest progress saves on victory
    Given the enemy is a quest target
    When combat ends in victory
    Then the quest progress is updated
    And any quest items are awarded

  Scenario: Loot auto-pickup option
    Given the player has auto-loot enabled
    When combat ends in victory
    Then all loot is automatically added to inventory
    And the loot summary is still shown
