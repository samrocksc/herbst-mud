Feature: Combat Screen Redesign
  As a player
  I want a dedicated combat screen with all 4 equipped abilities and a potion visible
  So that I can fight tactically instead of guessing what my hotkeys do

  Background:
    Given the player is in combat or presses Tab to preview combat mode
    And the player has up to 4 equipped combat skills plus a potion slot

  Scenario: Entering combat screen from adventurer view
    Given the player is on ScreenPlaying (adventurer view)
    When the player presses Tab or initiates combat
    Then the screen switches to ScreenCombat
    And the top 60% shows room text / combat log
    And the bottom 40% shows the combat HUD

  Scenario: Combat HUD displays 4 equipped abilities + potion
    Given the player is on ScreenCombat
    Then the bottom HUD shows exactly 5 slots:
      | Slot | Content         |
      | 1    | Equipped skill 1 |
      | 2    | Equipped skill 2 |
      | 3    | Equipped skill 3 |
      | 4    | Equipped skill 4 |
      | 5    | Potion           |
    And each slot displays:
      - Skill name (or "No skill" if empty)
      - Remaining cooldown as a radial overlay or countdown number
      - Disabled visual when on cooldown
    And the potion slot shows potion count (e.g. "3x")

  Scenario: Cooldown rendering updates in real time
    Given skill 1 has a 3.0s cooldown remaining
    When 1.5s elapses
    Then the HUD updates the cooldown indicator for slot 1
    And the indicator reaches 0 when the cooldown expires

  Scenario: Using a skill via hotkey during combat
    Given the player is on ScreenCombat
    And slot 1 is off cooldown
    When the player presses "1"
    Then the game sends the corresponding combat action
    And slot 1 enters cooldown state
    And the cooldown indicator begins counting down

  Scenario: Potion hotkey triggers potion use
    Given the player has potions in inventory
    When the player presses "5"
    Then a potion is consumed
    And the potion count decrements
    And a heal effect applies to the player

  Scenario: Tab toggles between combat and adventurer view
    Given the player is on ScreenCombat
    When the player presses Tab
    Then the screen returns to ScreenPlaying
    And the combat HUD is hidden
    When the player presses Tab again
    Then ScreenCombat returns with HUD visible
    And the previous cooldown states are preserved

  Scenario: Web-client combat HUD
    Given the player is using the browser client
    When in combat or Tab is pressed
    Then the web-client shows the same 5-slot bottom HUD
    And the HUD is overlaid on the game viewport
    And tapping a slot triggers the same action as the corresponding hotkey

  Scenario: Mobile / touch support for combat HUD
    Given the player is on a touch device
    Then the 5 HUD slots are large tappable buttons (>= 64x64dp)
    And each tap triggers the skill or potion
    And a haptic feedback or toast confirms the action
