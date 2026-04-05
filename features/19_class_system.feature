🔴 Feature: Class System - Implementation
  As a player
  I want to choose a class for my character
  So that I can have specialized abilities and playstyles

  Background:
    Given the database is connected
    And the class system is implemented

  Scenario: Character can have Warrior class
    Given I create a character
    When I set the class to "Warrior"
    Then the character should be a Warrior
    And the character should have Warrior-specific abilities
    And the character should gain Strength bonuses

  Scenario: Character can have Magician class
    Given I create a character
    When I set the class to "Magician"
    Then the character should be a Magician
    And the character should have Magician-specific spells
    And the character should gain Intellect bonuses

  Scenario: Character can have Thief class
    Given I create a character
    When I set the class to "Thief"
    Then the character should be a Thief
    And the character should have Thief-specific skills
    And the character should gain Dexterity bonuses

  Scenario: Character can have Charlatan class
    Given I create a character
    When I set the class to "Charlatan"
    Then the character should be a Charlatan
    And the character should have Charlatan-specific abilities
    And the character should gain social and deception bonuses

  Scenario: Character can have Vigilante class
    Given I create a character
    When I set the class to "Vigilante"
    Then the character should be a Vigilante
    And the character should have Vigilante-specific abilities
    And the character should gain balanced combat bonuses

  Scenario: Class selection during character creation
    Given I am creating a new character
    When I reach the class selection screen
    Then I should see all available classes
    And I should be able to select one class
    And the class should be saved to my character

  Scenario: Class affects starting stats
    Given I create a Warrior character
    Then the character should have higher Strength
    When I create a Magician character
    Then the character should have higher Intellect
    When I create a Thief character
    Then the character should have higher Dexterity

  Scenario: Class affects available skills
    Given I create a character with class "Warrior"
    Then the character should have access to Warrior skills
    And the character should NOT have access to Magician spells

  Scenario: Class cannot be changed after character creation
    Given a character "Hero" has class "Warrior"
    When I attempt to change the class to "Magician"
    Then the class change should be denied
    Or the character should be able to change class at special locations

  Scenario: Invalid class is rejected
    Given I am creating a new character
    When I select an invalid class "Dragon"
    Then I should see an error "Invalid class selection"
    And the character should not be created

  Scenario: Warrior class has highest health
    Given I create a Warrior character
    And I create a Magician character
    And I create a Thief character
    Then the Warrior should have the highest base HP
    And the Magician should have the lowest base HP

  Scenario: Magician class has highest mana
    Given I create a Magician character
    And I create a Warrior character
    Then the Magician should have the highest base Mana
    And the Warrior should have no or lowest base Mana

  Scenario: Thief class has highest initiative
    Given I create a Thief character
    Then the character should have highest initiative
    Or should have a stealth/avoidance bonus

  Scenario: Class is required for character
    Given I am creating a new character
    When I proceed without selecting a class
    Then I should see an error "Class selection is required"
    And the character should not be created

  Scenario: Each class has unique starting equipment
    Given I create a character with class "Warrior"
    Then the character should have Warrior-specific starting equipment
    When I create a character with class "Magician"
    Then the character should have Magician-specific starting equipment

  Scenario: Class affects damage calculation
    Given a Warrior attacks with a sword
    Then damage should be calculated with Strength modifier
    Given a Magician casts a spell
    Then damage should be calculated with Intellect modifier

  Scenario: Class determines ability score priority
    Given I view class "Warrior" details
    Then the priority stats should be Strength > Fortitude > Dexterity
    Given I view class "Magician" details
    Then the priority stats should be Intellect > Wisdom > Dexterity