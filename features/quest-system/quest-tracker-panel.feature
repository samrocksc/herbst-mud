Feature: Quest Tracker Panel
  Players can view their quest log with active, available, and completed quests

  Background:
    Given player is logged in
    And player has at least one active quest

  Scenario: View quest tracker with 'quests' command
    When player enters "quests"
    Then quest tracker panel is displayed
    And panel shows "QUEST LOG" title
    And quests are sorted by status

  Scenario: View quest tracker with 'q' alias
    When player enters "q"
    Then quest tracker panel is displayed
    And same content as "quests" command

  Scenario: View quest tracker with 'quest' command
    When player enters "quest"
    Then quest tracker panel is displayed
    And same content as "quests" command

  Scenario: Quest shows correct status - In Progress
    Given player has quest "Prove Yourself" in progress
    When player enters "quests"
    Then quest "Prove Yourself" shows status "In Progress"
    And status displayed in yellow color

  Scenario: Quest shows correct status - Available
    Given player has quest "Ooze Samples" available
    When player enters "quests"
    Then quest "Ooze Samples" shows status "Available"
    And status displayed in purple color

  Scenario: Quest shows correct status - Completed
    Given player has quest "First Blood" completed
    When player enters "quests"
    Then quest "First Blood" shows status "Completed"
    And status displayed in green color with strikethrough

  Scenario: Objective progress tracking
    Given quest "Prove Yourself" has objective "Kill Scrap Rat"
    And current progress is 2 of 3
    When player enters "quests"
    Then objective shows "Kill Scrap Rat (2/3)"
    And completed objectives show checkmark
    And incomplete objectives show circle

  Scenario: Quest details include giver
    Given quest "Prove Yourself" is given by "Guard Marco"
    When player enters "quests"
    Then quest displays "Giver: Guard Marco"

  Scenario: Quest details include rewards
    Given quest "Ooze Samples" rewards "repair_kit"
    When player enters "quests"
    Then quest displays "Reward: repair_kit"

  Scenario: Quest summary footer shows counts
    Given player has 3 active quests
    And player has 2 available quests
    And player has 5 completed quests
    When player enters "quests"
    Then footer shows "Active: 3 | Available: 2 | Completed: 5"

  Scenario: Quest shows description
    Given quest "Prove Yourself" has description "Guard Marco needs you to prove your worth..."
    When player enters "quests"
    Then quest description is displayed

  Scenario: Quest objectives are listed
    Given quest "Prove Yourself" has objectives:
      | objective           | progress |
      | Kill 3 Scrap Rats  | 2/3      |
      | Return to Guard    | 0/1      |
    When player enters "quests"
    Then all objectives are listed
    And progress shown for each

  Scenario: No quests shows empty message
    Given player has no quests
    When player enters "quests"
    Then message "No quests yet. Explore the world to find adventures!"

  Scenario: Quest tracker panel formatting
    Given player has quests
    When quest tracker is rendered
    Then panel uses box-drawing characters
    And quest names are bold
    And status badges have colors

  Scenario: Quest tracker accessible in combat
    Given player is in combat
    And player has active quest with objective to defeat enemy
    When player defeats enemy
    Then quest objective updates
    And quest tracker still accessible
