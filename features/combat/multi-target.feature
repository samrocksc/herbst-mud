Feature: Multi-Target Combat & Targeting
  As a player
  I want to fight multiple enemies and target them individually
  So that group combat is manageable and strategic

  Background:
    Given a combat is active with 3 enemies: "Scrap Rat A", "Scrap Rat B", "Raccoon Scout"
    And the player is engaged in combat

  Scenario: Target specific enemy by name
    Given multiple enemies are present
    When the player types "attack rat a"
    Then the player targets "Scrap Rat A"
    And the target is highlighted in turn order

  Scenario: Target enemy by index number
    Given multiple enemies are present
    When the player types "attack 2"
    Then the player targets the second enemy

  Scenario: Invalid target shows error
    Given multiple enemies are present
    When the player types "attack nonexistent"
    Then the game shows error: "no target found: nonexistent"
    And no action is taken

  Scenario: Invalid index shows error
    Given 3 enemies are present
    When the player types "attack 99"
    Then the game shows error: "invalid target index"
    And no action is taken

  Scenario: Turn order shows all combatants with targeted indicator
    Given multiple enemies are present
    And the player has selected a target
    When the turn order is displayed
    Then each combatant shows their initiative score
    And the player's current target is marked with [TARGETED]

  Scenario: AoE attack hits all enemies
    Given the player has an AoE ability available
    When the player types "attack all"
    Then all enemies take damage
    And the combat log shows "[Player] hits all enemies!"

  Scenario: Single target command only hits chosen enemy
    Given the player has targeted "Scrap Rat A"
    When the player attacks
    Then only "Scrap Rat A" takes damage
    And "Scrap Rat B" is unaffected

  Scenario: Target persists between rounds
    Given the player has targeted "Scrap Rat A"
    When a combat tick passes
    Then the target remains "Scrap Rat A"
    Until the player changes it

  Scenario: Can change target mid-combat
    Given the player has targeted "Scrap Rat A"
    When the player types "attack rat b"
    Then the target changes to "Scrap Rat B"
    And the new target is marked

  Scenario: Targeting a dead enemy shows error
    Given "Scrap Rat A" is dead
    When the player types "attack rat a"
    Then the game shows error: "target is no longer in combat"

  Scenario: Target command without argument targets current
    Given the player has an active target
    When the player types "target"
    Then the current target is confirmed
    And targeting information is displayed

  Scenario: Partial name matching works
    Given enemies named "Scrap Rat A" and "Scrap Rat B" are present
    When the player types "attack rat"
    Then the game targets the matching enemy
    Or shows a disambiguation list if multiple matches

  Scenario: Turn order sorted by initiative
    Given combatants with different initiative values
    When the turn order is displayed
    Then they are sorted highest to lowest initiative
