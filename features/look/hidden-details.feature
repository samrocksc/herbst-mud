Feature: Hidden Details
  As a player
  I want to discover hidden details about objects
  So that I can find secrets and unlock more information

  Background:
    Given player is exploring

  Scenario: Hidden details require skill check
    Given object has hidden detail with difficulty 20
    When player examines object
    Then skill check should be performed

  Scenario: High perception reveals hidden details
    Given player has perception skill at level 30
    And object has hidden detail with difficulty 20
    When player examines object
    Then hidden detail should be revealed

  Scenario: Low perception doesn't reveal hidden details
    Given player has perception skill at level 5
    And object has hidden detail with difficulty 20
    When player examines object
    Then hidden detail should NOT be revealed
    And message should show "You notice nothing unusual"

  Scenario: Hidden detail shows after multiple attempts
    Given player has perception 15
    And hidden detail difficulty is 20
    When player examines object 5 times
    Then eventually hidden detail should be revealed

  Scenario: Hidden detail shows skill threshold
    Given hidden secret requires perception 25
    When player has perception 25 or higher
    Then secret is revealed

  Scenario: Hidden detail remains after first reveal
    Given player has revealed hidden detail
    When player examines object again
    Then detail should be visible without skill check

  Scenario: Different skills reveal different hidden things
    Given object has hidden details
    And hidden detail A requires perception
    And hidden detail B requires investigation
    When player has perception but not investigation
    Then only perception detail should be revealed

  Scenario: Hidden items can be found
    Given room has hidden item "old coin"
    And item has perception difficulty 30
    When player has perception 30+
    Then item becomes visible

  Scenario: GM/debug mode reveals all hidden
    Given player is in debug mode
    When player examines any object
    Then all hidden details should be shown

  Scenario Outline: Perception check results
    Given player has perception <perception>
    And hidden detail has difficulty <difficulty>
    When player attempts to reveal
    Then result should be <result>

    Examples:
      | perception | difficulty | result    |
      | 10         | 20          | fail      |
      | 20         | 20          | pass      |
      | 30         | 20          | pass      |
      | 15         | 25          | fail      |
      | 26         | 25          | pass      |