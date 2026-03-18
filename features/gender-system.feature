Feature: Gender System Implementation (Issue #21)
  As a player
  I want to choose my character's gender/pronouns
  So that my character reflects my identity

  Background:
    Given the game has a gender system
    And I am creating a new character

  Scenario: Available pronouns
    Given I am at gender/pronoun selection
    Then I should see pronoun options:
    | Pronouns |
    | he/him   |
    | she/her  |
    | it/its   |
    | they/them|

  Scenario: Choose he/him pronouns
    Given I am creating a character
    When I select "he/him" pronouns
    Then my character should use he/him pronouns
    And game text should reflect this choice

  Scenario: Choose she/her pronouns
    Given I am creating a character
    When I select "she/her" pronouns
    Then my character should use she/her pronouns
    And game text should reflect this choice

  Scenario: Choose it/its pronouns
    Given I am creating a character
    When I select "it/its" pronouns
    Then my character should use it/its pronouns
    And game text should reflect this choice

  Scenario: Choose they/them pronouns
    Given I am creating a character
    When I select "they/them" pronouns
    Then my character should use they/them pronouns
    And game text should reflect this choice

  Scenario: Extensible pronoun system
    Given the gender system is extensible
    When new pronoun sets are added
    Then they should be available without code changes
    And the system should accept custom pronouns

  Scenario: Pronouns in game messages
    Given my character uses "they/them" pronouns
    When another player looks at my character
    Then the description should use "they/them"
    And actions should use correct verb forms

  Scenario: Gender stored on character
    Given a Character entity
    Then it should have a gender field
    And the field should store the selected pronouns
    And the field should be queryable