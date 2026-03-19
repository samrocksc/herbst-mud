Feature: Class System Implementation (Issue #19)
  As a player
  I want to choose a character class with unique abilities
  So that I can play with a distinct playstyle

  Background:
    Given the game has a class system
    And I am creating a new character

  Scenario: Available classes
    Given I am at class selection
    Then I should see five classes available:
    | Class      | Focus              |
    | Warrior    | Combat and strength|
    | Magician   | Magic and spells  |
    | Thief      | Stealth and speed |
    | Charlatan  | Trickery and charm|
    | Vigilante  | Justice and skill |

  Scenario: Choose Warrior class
    Given I am creating a character
    When I select "Warrior" class
    Then my character should have Warrior class
    And I should receive Warrior starting stats
    And I should gain Warrior abilities

  Scenario: Choose Magician class
    Given I am creating a character
    When I select "Magician" class
    Then my character should have Magician class
    And I should receive Magician starting stats
    And I should gain access to spell abilities

  Scenario: Choose Thief class
    Given I am creating a character
    When I select "Thief" class
    Then my character should have Thief class
    And I should receive Thief starting stats
    And I should gain stealth-based abilities

  Scenario: Choose Charlatan class
    Given I am creating a character
    When I select "Charlatan" class
    Then my character should have Charlatan class
    And I should receive Charlatan starting stats
    And I should gain trickery abilities

  Scenario: Choose Vigilante class
    Given I am creating a character
    When I select "Vigilante" class
    Then my character should have Vigilante class
    And I should receive Vigilante starting stats
    And I should gain justice-focused abilities

  Scenario: Class affects starting stats
    Given different classes have different stat bonuses
    When I choose a class
    Then the class should modify my base stats
    And the modifications should match the class theme

  Scenario Outline: Class stat bonuses
    Given I choose <class> class
    Then I should receive <primary_stat> bonus
    And my <secondary_stat> should be above average

    Examples:
      | class     | primary_stat | secondary_stat |
      | Warrior   | Strength     | Fortitude      |
      | Magician  | Intellect    | Wisdom         |
      | Thief     | Dexterity    | Wisdom         |
      | Charlatan | Wisdom       | Dexterity      |
      | Vigilante | Strength     | Dexterity      |