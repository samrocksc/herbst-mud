🔴 Feature: Quest System & Quest Tracker - Issues #19, #20
  As a player
  I want a quest tracker that shows my active quests
  So that I always know what to do next

  Background:
    Given the player has quests
    And the quest system is active

  # QUEST TRACKER PANEL
  Scenario: Quest tracker shows active quest name
    Given the player has an active quest
    When the quest tracker panel is displayed
    Then the quest name should be visible
    And the quest should be highlighted as current

  Scenario: Quest tracker shows current objective
    Given the player has an active quest with objectives
    When the quest tracker panel is displayed
    Then the current objective should be shown
    And the objective text should describe what to do

  Scenario: Quest tracker shows objective progress
    Given a quest objective has a count (e.g., "Kill 5 rats")
    When the player has killed some rats
    Then the tracker should show progress: "3/5 rats killed"
    And progress should update in real-time

  Scenario: Quest tracker shows multiple active quests
    Given the player has multiple active quests
    When the quest tracker panel is displayed
    Then all active quests should be listed
    And the primary quest should be highlighted

  Scenario: Quest tracker can be toggled open/closed
    Given the player is in game
    When the player types "quests" or uses a UI toggle
    Then the quest tracker panel should appear
    And it can be dismissed to return to normal view

  # QUEST TRACKER COMMANDS
  Scenario: Quest tracker shows help text
    Given the quest tracker is open
    When the player views the panel
    Then available commands should be listed
    And navigation hints should be provided

  Scenario: Select quest to view details
    Given the quest tracker is open
    When the player selects a quest number
    Then detailed quest info should display
    And all objectives should be visible

  Scenario: Abandon quest option
    Given the player wants to abandon a quest
    When the player selects abandon
    Then a confirmation prompt should appear
    And if confirmed, the quest should be removed

  Scenario: Quest tracker shows completed quests
    Given the player has completed quests
    When the player views completed quests section
    Then completed quests should be listed
    And rewards received should be shown

  # QUEST TRACKER PANEL UI
  Scenario: Quest tracker uses panel layout
    Given the quest tracker UI is implemented
    When the panel is rendered
    Then it should use proper panel styling
    And borders and padding should follow UI conventions

  Scenario: Quest tracker scrolls if too many quests
    Given the player has many quests
    When the quest list exceeds panel height
    Then the list should be scrollable
    And the player can navigate with arrow keys

  # QUEST PROGRESSION
  Scenario: Quest objectives update automatically
    Given a quest objective is to collect an item
    When the player picks up the item
    Then the quest tracker should update immediately

  Scenario: Quest completion triggers notification
    Given all objectives of a quest are met
    When the quest is ready to turn in
    Then a notification should appear
    And the quest should be highlighted in the tracker

  Scenario: Turn in quest at quest giver
    Given a quest is complete
    When the player interacts with the quest giver
    Then the quest should be turned in
    And rewards should be granted
    And XP should be applied
