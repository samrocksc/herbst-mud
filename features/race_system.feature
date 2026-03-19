Feature: Race System
  As a game developer
  I want a Race system
  So that players can choose different character races

  Scenario: Available races
    When I list all available races
    Then I should see: Human, Elf, Dwarf, Orc
    And each race should have race bonuses

  Scenario: Race bonuses
    Given the "Human" race
    Then it should have: +10% experience gain
    And "Elf" should have: +20% magic resistance
    And "Dwarf" should have: +10% health
    And "Orc" should have: +10% physical damage

  Scenario: Race selection
    Given I am creating a character
    When I select race "Elf"
    Then my character should have Elf racial bonuses