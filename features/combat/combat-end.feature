🔴 Feature: Combat End (Victory/Defeat Screens) - Issue #12
  As a combat system
  I need victory and defeat screens with loot, XP, and consequences
  So that combat has satisfying resolution

  Background:
    Given combat has been initiated
    And combatants are tracked throughout combat

  # VICTORY
  Scenario: Victory screen shows on enemy defeat
    Given all enemies have 0 HP
    When the last enemy is defeated
    Then a victory screen should display
    And the screen should show combat statistics

  Scenario: Victory screen shows earned XP
    Given the player wins combat
    When the victory screen appears
    Then XP earned should be displayed
    And the value should be based on enemy level and count

  Scenario: Loot is displayed from defeated enemies
    Given enemies are defeated
    When the victory screen appears
    Then dropped items should be listed
    And gold rewards should be displayed
    And the player should be able to claim loot

  Scenario: Loot is added to inventory
    Given the player claims loot after victory
    When the player confirms loot collection
    Then items should be added to player inventory
    And gold should be added to player purse

  Scenario: Loot can be declined
    Given the player sees loot options after victory
    When the player declines loot
    Then the loot should not be added
    And the player proceeds to post-combat state

  # DEFEAT
  Scenario: Defeat screen shows on player death
    Given the player has 0 HP
    When the player is defeated
    Then a defeat screen should display
    And consequences should be shown

  Scenario: Death causes experience loss
    Given the player is defeated in combat
    When the defeat screen appears
    Then a penalty message should be shown
    And the player should lose some XP
    And XP loss should not cause level down

  Scenario: Player respawns after death
    Given the player is defeated
    When the player acknowledges defeat
    Then the player should respawn
    And the player should be at a safe location
    And HP should be restored to full or partial

  # POST-COMBAT STATE
  Scenario: Player returns to exploration after combat
    Given combat ends (victory or defeat)
    When the player dismisses the end screen
    Then normal exploration mode resumes
    And the player can continue playing

  Scenario: Combat log is cleared or archived after combat ends
    Given combat ends
    When the player returns to exploration
    Then the combat log may be archived
    Or a new exploration log begins
