Feature: Race System - Implementation
  As a player
  I want to choose a character race with unique bonuses
  So that I can customize my character's strengths and play style

  Background:
    Given the race system is implemented
    And the following races are available: Human, Elf, Dwarf, Orc

  Scenario: All four base races are available
    When I request the list of available races
    Then the response should include: Human, Elf, Dwarf, Orc
    And the count should be exactly 4 races

  Scenario: Each race has a description
    When I query race details
    Then each race should have a description field explaining the race traits

  # Human Race
  Scenario: Human has correct racial bonuses
    When I create or examine a Human character
    Then Humans should have a +10% experience gain bonus
    And Humans should have no stat penalties
    And Humans should have well-rounded base stats (no extremes)

  Scenario: Human race bonus applies to experience gain
    Given a Human character earns 100 experience points
    When the XP is applied to the character
    Then the effective XP gained should be 110 (100 + 10% bonus)
    And the bonus should stack correctly with other XP modifiers

  Scenario: Human racial description
    When I examine the Human race
    Then the description should mention: "adaptable" or "versatile" or "quick learners"

  # Elf Race
  Scenario: Elf has correct racial bonuses
    When I create or examine an Elf character
    Then Elves should have a +20% magic resistance
    And Elves should have no stat penalties
    And Elves should have slightly higher intelligence and agility

  Scenario: Elf magic resistance reduces magical damage
    Given an Elf character receives 100 magical damage
    When the damage is calculated
    Then the actual damage taken should be reduced by 20%
    And the character should take 80 damage instead of 100

  Scenario: Elf magic resistance does not affect physical damage
    Given an Elf character receives 100 physical damage
    When the damage is calculated
    Then the actual damage taken should be 100 (no reduction)

  Scenario: Elf racial description
    When I examine the Elf race
    Then the description should mention: "magical" or "arcane" or "resistant"

  # Dwarf Race
  Scenario: Dwarf has correct racial bonuses
    When I create or examine a Dwarf character
    Then Dwarves should have a +10% health bonus
    And Dwarves should have slightly higher constitution
    And Dwarves should have slightly lower agility

  Scenario: Dwarf health bonus applies to max health
    Given a Dwarf character has a base max health of 100
    When the character is created
    Then the max health should be 110 (100 + 10% bonus)

  Scenario: Dwarf health bonus also affects current health proportionally
    Given a Dwarf is created with base health matching max health
    Then the starting health should also include the +10% bonus

  Scenario: Dwarf racial description
    When I examine the Dwarf race
    Then the description should mention: "hardy" or "sturdy" or "health"

  # Orc Race
  Scenario: Orc has correct racial bonuses
    When I create or examine an Orc character
    Then Orcs should have a +10% physical damage bonus
    And Orcs should have slightly higher strength
    And Orcs should have slightly lower intelligence

  Scenario: Orc physical damage bonus increases attack power
    Given an Orc character has a base physical attack of 100
    When the character performs a physical attack
    Then the damage dealt should be 110 (100 + 10% bonus)

  Scenario: Orc physical damage bonus does not affect magical attacks
    Given an Orc Mage character performs a magical attack
    When the damage is calculated
    Then the magical damage should not receive the +10% bonus

  Scenario: Orc racial description
    When I examine the Orc race
    Then the description should mention: "powerful" or "fierce" or "strength"

  # Race Selection
  Scenario: Race is selected during character creation
    Given I am creating a new character
    When I select race "Elf"
    Then the character should have Elf racial bonuses applied
    And the character should have Elf stat modifiers

  Scenario: Race cannot be changed after character creation
    Given I have a character with race "Human"
    When I attempt to change the race to "Orc"
    Then the operation should be denied
    Or the character should need to be recreated

  Scenario: Invalid race name is rejected
    Given I am creating a character
    When I attempt to select race "Giant"
    Then the response should be a validation error
    And the error should list valid race options

  # Race and Class Combination
  Scenario: Race and class bonuses stack independently
    Given I create a character with race "Elf" and class "Mage"
    Then the character should have Elf magic resistance
    And the character should have Mage base stats
    And both bonuses should apply without conflict

  Scenario: Dwarf Warrior has high health and strength
    Given I create a character with race "Dwarf" and class "Warrior"
    Then the character should have +10% max health from Dwarf
    And the character should have high strength from Warrior
    And the character should have both racial and class bonuses

  # Race Bonuses with Leveling
  Scenario: Race bonuses persist through leveling
    Given a Human character levels up
    When the new level stats are calculated
    Then the +10% XP bonus should still apply
    And no other race bonuses should be affected

  # Edge Cases
  Scenario: Race name is case-sensitive
    Given I attempt to create a character with race "elf"
    Then the response should be a validation error
    And valid races should be listed as: Human, Elf, Dwarf, Orc
