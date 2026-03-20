Feature: Examine Skill System
  Examine skill (0-100) determines what hidden details players can discover.

  Background:
    Given player exists
    And examine skill system is initialized

  Scenario: New character has examine skill
    Given new character is created
    Then examine skill starts at level 1
    Or examine skill based on INT

  Scenario: Examine skill gains XP from examining
    Given player examines item for first time
    When examine action completes
    Then examine skill gains 1 XP

  Scenario: Discover hidden detail at skill threshold
    Given item has hidden detail requiring level 25
    And player has examine level 30
    When player examines item
    Then hidden detail is revealed

  Scenario: Hidden detail not revealed at low skill
    Given item has hidden detail requiring level 50
    And player has examine level 30
    When player examines item
    Then hidden detail is not shown
    And message suggests raising examine skill

  Scenario: Examine skill bonus at level 26-50
    Given player has examine level 40
    When checking reveal chance
    Then bonus of 10% applies

  Scenario: Examine skill bonus at level 76-90
    Given player has examine level 80
    When checking reveal chance
    Then bonus of 50% applies

  Scenario: INT stat affects examine
    Given player has INT 15
    When examine check is made
    Then INT bonus of +15 applies

  Scenario: Master examiner at level 100
    Given player has examine level 100
    When examining any item
    Then all hidden details are revealed

  Scenario: Discovering secret gives bonus XP
    Given player reveals hidden compartment
    When check succeeds
    Then examine gains 5 XP