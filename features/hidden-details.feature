Feature: Look - Hidden Details System (look-04)
  As a player with the Examine skill
  I want to discover hidden details on items and in rooms
  So that I can uncover secrets and gain experience

  Background:
    Given the server is running
    And I am logged in as a test character
    And my character has an Examine skill level

  @donnie
  Scenario: Automatic reveal of hidden details based on skill threshold
    Given an item has a hidden detail with "automatic" mode
    And the detail requires an Examine skill of 15
    When my Examine skill is level 20
    Then the hidden detail should be automatically revealed
    And I should see the "skill_threshold" source message

  @donnie
  Scenario: Hidden detail remains hidden when skill is too low
    Given an item has a hidden detail with "automatic" mode
    And the detail requires an Examine skill of 25
    When my Examine skill is level 15
    Then the hidden detail should remain hidden

  @donnie
  Scenario: Check-based reveal with successful roll
    Given an item has a hidden detail with "check" mode
    And the difficulty class is 30
    And the detail requires an Examine skill of 20
    When my Examine skill is level 25
    And my INT stat is 12
    And my WIS stat is 10
    And the skill check passes the DC
    Then the hidden detail should be revealed
    And I should see the "check_passed" source message
    And my check roll should be displayed

  @donnie
  Scenario: Check-based reveal with failed roll
    Given an item has a hidden detail with "check" mode
    And the difficulty class is 50
    And the detail requires an Examine skill of 20
    When my Examine skill is level 25
    And my INT stat is 10
    And my WIS stat is 10
    And the skill check fails the DC
    Then the hidden detail should remain hidden

  @donnie
  Scenario: Gaining XP from revealing hidden details
    Given I successfully reveal a hidden detail
    When I discover the hidden detail for the first time
    Then I should gain 2 XP (ExamineXPDiscover)
    When I reveal another hidden detail
    Then I should gain the discover XP again

  @donnie
  Scenario: Level up from Examine skill usage
    Given my Examine skill is level 5
    And my Examine XP is 8
    When I successfully reveal a hidden detail
    Then my skill should level up to 6
    And my Examine XP should reset to 0

  @donnie
  Scenario: Examine bonus percentage based on skill level
    Given my Examine skill is level 76
    Then I should receive a 50% bonus on examine actions

  @donnie
  Scenario: Max skill level cap
    Given my Examine skill is level 100
    When I gain XP
    Then my skill level should remain at 100
    And excess XP should be handled appropriately