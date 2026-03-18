Feature: Race System
  As a player
  I want to choose a race for my character
  So that I can have unique racial traits

  Scenario: List available races
    When I request available races
    Then I should see: Human, Elf, Dwarf, Orc

  Scenario: Human racial traits
    Given I choose the "Human" race
    Then I should have +10 to all stats

  Scenario: Elf racial traits
    Given I choose the "Elf" race
    Then I should have +20 dexterity
    And I should have +10 intelligence
    And I should have low constitution

  Scenario: Dwarf racial traits
    Given I choose the "Dwarf" race
    Then I should have +20 constitution
    And I should have +10 strength
    And I should have low charisma

  Scenario: Orc racial traits
    Given I choose the "Orc" race
    Then I should have +25 strength
    And I should have +10 constitution
    And I should have -10 intelligence