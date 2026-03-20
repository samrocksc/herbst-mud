Feature: Hidden Details System
  As a player
  I want hidden details to be revealed through examination
  So that exploration is rewarding

  Background:
    Given the game is running
    And items can have hidden details

  Scenario: Automatic detail revealed at skill threshold
    Given item "fountain" has hidden detail "Coins glint at the bottom"
    And the detail requires examine level 0
    When the player with examine level 10 examines the fountain
    Then "Coins glint at the bottom" should be revealed

  Scenario: Automatic detail not revealed below threshold
    Given item "fountain" has hidden detail requiring level 75
    When the player with examine level 50 examines the fountain
    Then the hidden detail should not be shown

  Scenario: Check mode requires rolling against DC
    Given item has hidden detail with mode "check" and DC 30
    When the player examines the item
    Then a roll should be performed against DC 30
    And if roll >= 30, the detail is revealed
    And if roll < 30, the detail is not revealed

  Scenario: Perception check uses WIS
    Given item has hidden detail with stat "WIS" and DC 25
    And the player has WIS of 15
    When the player examines the item
    Then the check should be (examine level + 15 + random) vs DC 25

  Scenario: Multiple hidden details on one item
    Given item "fountain" has 3 hidden details
    And detail 1 requires level 0 (automatic)
    And detail 2 requires level 30 (automatic)
    And detail 3 requires level 75 (automatic)
    When the player with examine level 40 examines the fountain
    Then detail 1 should be revealed
    And detail 2 should be revealed
    And detail 3 should not be revealed

  Scenario: Check mode can fail even at high skill
    Given item has hidden detail with mode "check" and DC 50
    When the player with examine level 80 examines the item
    Then the detail may or may not be revealed based on roll

  Scenario: Hidden details show source in output
    Given hidden detail is revealed through skill threshold
    When the examine output is displayed
    Then the detail should indicate it was revealed by "skill_threshold"

  Scenario: Hidden details show as unrevealed with requirement
    Given item has hidden detail requiring level 75
    When the player with examine level 50 examines the item
    Then the output should show "Requires examine level 75"

  Scenario: XP capped at 3 per examine
    Given an item has 5 hidden details that can be revealed
    When the player examines the item
    Then maximum 3 XP should be granted

  Scenario: Examine skill gains XP from revealing details
    Given the player reveals 2 hidden details
    When examine skill XP is calculated
    Then the player should gain 4 XP (2 details * 2 XP)

  Scenario: Hidden detail with check mode shows failure message
    Given item has hidden detail with mode "check" and DC 40
    And the player's check fails
    When the player examines the item
    Then the hidden detail is not revealed
    And feedback "You notice nothing unusual" is shown

  Scenario: Revealed detail persists across sessions
    Given the player revealed hidden detail "secret_note"
    When the player reconnects to the game
    Then the detail remains revealed
    And no XP is granted for re-examining

  Scenario: Multiple check failures do not stack penalty
    Given item has multiple hidden details with check mode
    When the player fails multiple checks
    Then only one failure message is shown
    And the details remain hidden