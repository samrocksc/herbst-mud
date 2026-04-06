🔴 Feature: Examine Skill - Issue Look-03
  As a player
  I want to use skills when examining things
  So that hidden details are revealed based on my abilities

  Background:
    Given the player has skills
    And examinable objects may have hidden details

  Scenario: High perception reveals hidden details
    Given a player has high perception skill (60+)
    And an object has a hidden detail requiring perception 50
    When the player examines the object
    Then the hidden detail should be revealed
    And a message should indicate the discovery

  Scenario: Low perception does not reveal hidden details
    Given a player has low perception skill (20)
    And an object has a hidden detail requiring perception 50
    When the player examines the object
    Then the hidden detail should NOT be revealed
    And the examine output should be normal

  Scenario: Skill check is shown in examine output
    Given an object has hidden details at various skill thresholds
    When the player examines the object
    Then the output should indicate what skill level would be needed
    And the player should know what to improve

  Scenario: Examine skill improves with use
    Given the player examines many objects
    When examine skill reaches milestone thresholds
    Then new hidden details should become visible
    And the skill should feel progressively more useful

  Scenario: Different skills reveal different things
    Given an object has multiple hidden aspects
    And perception skill reveals one aspect
    And investigation skill reveals another
    When the player examines with appropriate skills
    Then different skills reveal different hidden details
