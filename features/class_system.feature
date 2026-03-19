Feature: Class System
  As a game developer
  I want a Class system
  So that players can choose different character classes

  Scenario: Available classes
    When I list all available classes
    Then I should see: Warrior, Mage, Rogue, Priest
    And each class should have base stats

  Scenario: Class base stats
    Given the "Warrior" class
    Then it should have: high health, medium mana, high strength
    And "Mage" should have: low health, high mana, high intelligence
    And "Rogue" should have: medium health, medium mana, high agility
    And "Priest" should have: medium health, high mana, high wisdom

  Scenario: Class selection
    Given I am creating a character
    When I select class "Warrior"
    Then my character should have Warrior base stats