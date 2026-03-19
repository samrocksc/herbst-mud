Feature: Quest Tracker Panel
  As a player
  I want to see my active quests and their progress
  So that I can track my objectives while playing

  Background:
    Given I am logged into the game
    And I have an active character

  Scenario: Display quest log with quests command
    Given I have active quests in my log
    When I type "quests"
    Then I should see the quest tracker panel
    And I should see my active quests listed
    And I should see quest status indicators

  Scenario: Display quest log with q alias
    Given I have active quests in my log
    When I type "q"
    Then I should see the quest tracker panel
    And I should see my active quests listed

  Scenario: Show quest progress
    Given I have a quest "Prove Yourself" with progress 2/3
    When I view my quest log
    Then I should see "Prove Yourself" listed
    And I should see progress "2/3"
    And I should see status "In Progress"

  Scenario: Display available quests
    Given I have a quest "Ooze Samples" that is available
    When I view my quest log
    Then I should see "Ooze Samples" listed
    And I should see status "Available"

  Scenario: Display completed quests
    Given I have a completed quest "First Steps"
    When I view my quest log
    Then I should see "First Steps" listed
    And I should see status "Completed"

  Scenario: Empty quest log message
    Given I have no quests in my log
    When I type "quests"
    Then I should see a message indicating no active quests

  Scenario: Quest panel styling
    Given I have active quests
    When I view my quest log
    Then the quest panel should use Lip Gloss styling
    And the panel should be properly formatted

  Scenario: API integration for quest data
    Given the quest API is available
    When I view my quest log
    Then quest data should be fetched from the API
    And the quest log should reflect current quest status

  Scenario: Handle API unavailability gracefully
    Given the quest API is unavailable
    When I view my quest log
    Then I should see fallback placeholder data
    And I should not see an error message

  Scenario: Handle malformed API response
    Given the quest API returns malformed data
    When I view my quest log
    Then I should see an error indication
    Or I should see fallback data

  Scenario Outline: Quest status indicators
    Given I have a quest with status "<status>"
    When I view my quest log
    Then I should see the status indicator for "<status>"

    Examples:
      | status       |
      | In Progress  |
      | Available    |
      | Completed    |

  Scenario: Quest objectives display
    Given I have a quest "Prove Yourself" with objective "Kill 3 Scrap Rats"
    And I have killed 2 Scrap Rats
    When I view my quest log
    Then I should see "Kill 3 Scrap Rats"
    And I should see progress towards the objective