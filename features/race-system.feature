Feature: Race System
  As a player
  I want to select a mutant race with unique bonuses and abilities
  So that my character has distinct advantages and playstyle

  Background:
    Given I am creating a new character
    And I have reached the race selection screen

  Scenario: View available races
    Given I am on the race selection screen
    When I view the race options
    Then I should see 5 available races
    And the races should be: human, turtle, rabbit, rat, rhino

  Scenario: Select human race
    Given I am on the race selection screen
    When I select "human"
    Then my character race should be set to "human"
    And I should receive +1 to all stats (Versatile ability)
    And I should receive +1 to every skill level

  Scenario: Select turtle race
    Given I am on the race selection screen
    When I select "turtle"
    Then my character race should be set to "turtle"
    And I should receive +CON bonus
    And I should receive -INT penalty
    And I should gain "Innate Block" passive ability
    And I should see the block chance description

  Scenario: Select rabbit race
    Given I am on the race selection screen
    When I select "rabbit"
    Then my character race should be set to "rabbit"
    And I should receive +DEX bonus (+1 level)
    And I should receive -STR penalty
    And I should gain "Quick Strike" ability
    And I should see faster attack speed description

  Scenario: Select rat race
    Given I am on the race selection screen
    When I select "rat"
    Then my character race should be set to "rat"
    And I should receive +WIS bonus
    And I should receive -STR penalty
    And I should gain "Stealth" ability
    And I should see improved stealth description

  Scenario: Select rhino race
    Given I am on the race selection screen
    When I select "rhino"
    Then my character race should be set to "rhino"
    And I should receive +STR bonus
    And I should receive -WIS penalty
    And I should gain "Gore" ability
    And I should see gore chance description that scales with level

  Scenario: Turtle innate block triggers in combat
    Given I have created a turtle character
    And I am in combat
    When an enemy attacks me with melee
    Then there should be a low percentage chance to block
    And when block triggers, I should see visual feedback
    And the attack should be negated

  Scenario: Rabbit quick strike in combat
    Given I have created a rabbit character
    And I am in combat
    When I attack an enemy
    Then my dexterity should include the +1 level bonus
    And my attack speed should be faster than baseline

  Scenario: Rat stealth ability
    Given I have created a rat character
    When I attempt to avoid detection
    Then my stealth capability should be enhanced
    And I should have higher chance to avoid NPC detection

  Scenario: Rhino gore ability scales with level
    Given I have created a rhino character
    And I am in combat
    When I attack in melee
    Then there should be a random chance to trigger Gore
    And the base percentage should increase slightly per level

  Scenario: Race affects skill leveling
    Given I have created a human character
    When I level up skills
    Then all skills should receive +1 level bonus (Versatile ability)

  Scenario: Race persists across sessions
    Given I have created a character with turtle race
    When I log out and log back in
    Then my race should still be "turtle"
    And my stat bonuses should be applied
    And my racial ability should still be active

  Scenario: API returns character race
    Given I have created a character
    When I call GET /characters/:id/race
    Then I should receive the character's race
    And the race should be one of: human, turtle, rabbit, rat, rhino

  Scenario: API updates character race
    Given I have created a character with race "human"
    When I call PUT /characters/:id/race with {"race": "turtle"}
    Then my race should be updated to "turtle"
    And my bonuses should be recalculated
    And my abilities should be updated

  Scenario: Invalid race rejected
    Given I have created a character
    When I call PUT /characters/:id/race with {"race": "dragon"}
    Then I should receive a 400 error
    And the error should indicate invalid race value