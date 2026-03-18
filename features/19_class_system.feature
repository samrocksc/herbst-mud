Feature: Class System
  As a player
  I want to choose a class for my character
  So that I can have unique abilities and playstyles

  Scenario: List available classes
    When I request available classes
    Then I should see: Warrior, Mage, Rogue, Priest

  Scenario: Warrior class attributes
    Given I choose the "Warrior" class
    Then I should have high strength
    And I should have high constitution
    And I should have special ability "Battle Cry"

  Scenario: Mage class attributes
    Given I choose the "Mage" class
    Then I should have high intelligence
    And I should have high wisdom
    And I should have special ability "Fireball"

  Scenario: Rogue class attributes
    Given I choose the "Rogue" class
    Then I should have high dexterity
    And I should have special ability "Sneak"

  Scenario: Priest class attributes
    Given I choose the "Priest" class
    Then I should have high wisdom
    And I should have high charisma
    And I should have special ability "Heal"