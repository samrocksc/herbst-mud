🔴 Feature: Examine Events - Issue Look-09 / Issue #17
  As a player
  I want to examine events and understand world events
  So that I can participate in timed or special activities

  Background:
    Given the event system is active
    And events may occur in the game world

  Scenario: Examine shows current world events
    Given world events are active
    When the player types "events" or "examine event"
    Then active events should be listed
    And each event should show name, time remaining, and description

  Scenario: Examine event shows participation details
    Given an event is active
    When the player examines the event
    Then how to participate should be explained
    And rewards for participation should be hinted

  Scenario: Event countdown is shown
    Given an event has a time limit
    When the player checks the event
    Then time remaining should be displayed
    And urgency should be communicated

  Scenario: Event completion is tracked
    Given the player participates in an event
    When the player completes event objectives
    Then participation progress should be tracked
    And rewards should be granted when event ends

  Scenario: Events have categories
    Given multiple events exist
    When the player views events
    Then events should be categorized (PvP, PvE, crafting, social)
    And the player can filter by category

  Scenario: Seasonal/holiday events appear at times
    Given the game has seasonal event system
    When a seasonal event time arrives
    Then the event should appear automatically
    And the event should be accessible to all players

  Scenario: Event completion message
    Given an event ends with player participation
    When the event rewards are distributed
    Then a summary message should display
    And the player should see what they earned

  Scenario: No active events shows helpful message
    Given no events are currently active
    When the player types "events"
    Then a message should indicate no current events
    And when the next event is expected may be hinted

  Scenario: Event notifications appear
    Given an event is about to start or end
    When the notification threshold is reached
    Then players should receive a game-world notification
    And the notification should be in the output pane

  Scenario: Failed event participation shows why
    Given a player tries to join an event they don't qualify for
    When the player attempts to participate
    Then a message should explain why (level, quest, location)
