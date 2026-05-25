Feature: Level 1 Training Dummy NPC
  As a player
  I want an immortal training dummy in a dedicated training room
  So that I can test combat mechanics, rotations, and the new combat screen safely

  Background:
    Given the game world has a Training Room connected to the starting area
    And the Training Room contains an immortal NPC named "Training Dummy"

  Scenario: Locating the Training Room
    Given the player is in The Hole (starting room)
    When the player looks around or checks the map
    Then an exit "training" is visible leading to the Training Room
    And the Training Room is east of The Hole

  Scenario: Training Dummy is present and targetable
    Given the player enters the Training Room
    Then the room description mentions a wooden training dummy
    And "Training Dummy" appears in the NPC list
    And the player can target the dummy with "target dummy" or "t dummy"
    And the dummy shows as Level 1

  Scenario: Dummy has unlimited health and auto-heals
    Given the player is fighting the Training Dummy
    When the player deals damage to the dummy
    Then the dummy takes the damage and shows reduced HP
    And after 3 seconds the dummy heals back to full HP
    And the dummy never dies or drops loot
    And the dummy never flees

  Scenario: Dummy deals zero damage
    Given the player is in combat with the Training Dummy
    When the dummy's combat turn arrives
    Then the dummy attempts an attack
    But the damage dealt is always 0
    And the player receives a message: "The dummy swings harmlessly at you."

  Scenario: Dummy allows testing all 4 combat skills
    Given the player has 4 combat skills equipped
    And the player is targeting the Training Dummy
    When the player uses each skill in sequence
    Then each skill applies its effect
    And cooldowns are triggered normally
    And the dummy remains alive after all skills

  Scenario: Web-client combat screen with dummy
    Given the player is in the browser client
    And the player is in the Training Room with the dummy targeted
    When the combat screen (HUD) is visible
    Then pressing 1-4 uses skills against the dummy
    And pressing 5 uses a potion
    And the combat log shows "You hit Training Dummy for X damage"
    And the dummy auto-heal message appears periodically

  Scenario: Training Room is reset-safe
    Given the game server restarts
    When the player enters the Training Room
    Then the Training Dummy is still present
    And its stats remain Level 1, immortal, 0 damage

  Scenario Outline: Dummy supports different player levels
    Given a player of level <level>
    When they fight the dummy
    Then the dummy still auto-heals and takes 0-damage swings
    And the dummy remains targetable regardless of player level

    Examples:
      | level |
      | 1     |
      | 5     |
      | 10    |
      | 50    |
