🔴 Feature: Examine Quests - Issue Look-08 / Issue #16
  As a player
  I want to examine quests for details and progress
  So that I can track what I need to do

  Background:
    Given the player has active quests
    And quest system is initialized

  Scenario: Examine quest shows quest details
    Given the player has a quest
    When the player types "examine quest <N>" or "quest <N>"
    Then the quest name and description should be shown
    And the current objective should be highlighted
    And any relevant hints should be provided

  Scenario: Examine quest shows progress
    Given a quest has multiple objectives
    When the player examines the quest
    Then completed objectives should be marked
    And remaining objectives should be listed
    And progress percentage should be shown

  Scenario: Examine quest shows reward preview
    Given a quest has rewards
    When the player examines the quest
    Then the expected rewards should be listed
    And XP, items, and gold should be indicated

  Scenario: Examine current quest shorthand
    Given the player has one active quest
    When the player types "quest" (no args)
    Then the current/primary quest should display
    And full quest details should be shown

  Scenario: Quest journal lists all active quests
    Given the player has multiple active quests
    When the player types "quests" or "journal"
    Then all active quests should be listed
    And each quest should show name and brief status

  Scenario: Completed quests can be reviewed
    Given the player has completed quests
    When the player types "quests completed"
    Then completed quest history should be shown
    And the player can see what they've accomplished

  Scenario: Quest auto-updates when objective completes
    Given a player is working on a quest
    When an objective is completed (e.g., item collected, mob killed)
    Then the quest progress should update automatically
    And a message should indicate quest progress

  Scenario: Quest completion is announced
    Given all objectives of a quest are met
    When the player interacts with the quest giver or completes final step
    Then a completion message should display
    And rewards should be granted
    And the quest should move to completed list
