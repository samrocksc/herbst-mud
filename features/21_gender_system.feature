Feature: Gender System - Implementation
  As a player
  I want to select my character's gender
  So that my character reflects my identity and has appropriate in-game representation

  Background:
    Given the gender system is implemented
    And the following genders are available: Male, Female, Non-binary, Other

  Scenario: All four gender options are available
    When I request the list of available gender options
    Then the response should include: Male, Female, Non-binary, Other
    And the count should be exactly 4 options

  Scenario: Gender is an optional field during character creation
    When I create a character without specifying gender
    Then the character should still be created successfully
    And the gender should default to "Other" or null

  # Male Gender
  Scenario: Male gender is stored correctly
    When I create a character with gender "Male"
    Then the character should have gender "Male" stored in the database
    And the gender should be returned correctly in API responses

  Scenario: Male gender is displayed correctly
    Given a character with gender "Male" exists
    When I examine the character
    Then the gender should be displayed as "Male"

  # Female Gender
  Scenario: Female gender is stored correctly
    When I create a character with gender "Female"
    Then the character should have gender "Female" stored in the database
    And the gender should be returned correctly in API responses

  Scenario: Female gender is displayed correctly
    Given a character with gender "Female" exists
    When I examine the character
    Then the gender should be displayed as "Female"

  # Non-binary Gender
  Scenario: Non-binary gender is stored correctly
    When I create a character with gender "Non-binary"
    Then the character should have gender "Non-binary" stored
    And the exact string "Non-binary" should be preserved

  Scenario: Non-binary gender is displayed correctly
    Given a character with gender "Non-binary" exists
    When I examine the character
    Then the gender should be displayed as "Non-binary"
    And the hyphen should be preserved

  # Other Gender
  Scenario: Other gender is stored correctly
    When I create a character with gender "Other"
    Then the character should have gender "Other" stored
    And the gender should be returned correctly in API responses

  # Gender Validation
  Scenario: Invalid gender value is rejected
    Given I am creating a character
    When I attempt to set gender to "unknown" or "none"
    Then the validation should fail
    And the error should list valid gender options

  Scenario: Gender is case-sensitive
    Given I attempt to create a character with gender "male"
    Then the validation should fail
    And valid options should be shown as: Male, Female, Non-binary, Other

  Scenario: Empty gender string is handled
    Given I attempt to create a character with gender ""
    Then the character should either default to "Other"
    Or the validation should fail with an error

  # Gender in Character Display
  Scenario: Gender appears in character status
    Given a character with gender "Female" exists
    When I check the character status or sheet
    Then the gender should be listed as part of the character profile

  Scenario: Gender appears in who/player list
    Given characters with different genders are logged in
    When other players view the online player list
    Then each player's gender should be visible in the list

  # Gender in Communication
  Scenario: Pronouns are correctly associated with gender
    Given a character with gender "Female" is created
    When the game refers to the character in narration
    Then female pronouns (she/her) should be used

    Given a character with gender "Male" is created
    When the game refers to the character in narration
    Then male pronouns (he/him) should be used

    Given a character with gender "Non-binary" is created
    When the game refers to the character in narration
    Then neutral pronouns (they/them) should be used

  # Gender Change (if implemented)
  Scenario: Gender can be changed after character creation
    Given a character with gender "Male" exists
    When I attempt to change the gender to "Female"
    Then the gender should be updated successfully
    And the new gender should be reflected in all displays

  # Gender Display in Look Command
  Scenario: Look command shows character gender
    Given a character "JaneDoe" with gender "Female" exists
    When I use the look command on "JaneDoe"
    Then the output should include "Female" or a female indicator

  # Gender Storage
  Scenario: Gender is stored as an enumerated string
    When I examine the database schema for the gender field
    Then the field type should be: string or enum
    And the valid values should be constrained to the four options
