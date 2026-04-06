Feature: Class System - Implementation
  As a player
  I want to choose a character class with unique abilities
  So that I can play with a specialized combat and skill style

  Background:
    Given the class system is implemented
    And the following classes are available: Warrior, Mage, Rogue, Priest

  Scenario: All four base classes are available
    When I request the list of available classes
    Then the response should include: Warrior, Mage, Rogue, Priest
    And the count should be exactly 4 classes

  Scenario: Each class has a description
    When I query class details
    Then each class should have a description field explaining the class role

  # Warrior Class
  Scenario: Warrior has correct base stats
    When I create or examine a Warrior character
    Then the base health should be: high (around 120-150)
    And the base mana should be: low (around 20-40)
    And the base strength should be: high (around 15-20)
    And the base agility should be: medium (around 8-12)
    And the base intelligence should be: low (around 3-6)
    And the base wisdom should be: low (around 3-6)

  Scenario: Warrior can equip heavy weapons
    When I check Warrior equipment restrictions
    Then Warriors should be able to equip: swords, axes, maces, shields
    And Warriors should NOT be able to equip: staffs, wands

  Scenario: Warrior combat style description
    When I examine the Warrior class
    Then the combat style should be: "melee" or "physical"
    And the role should be: "tank" or "damage dealer"

  # Mage Class
  Scenario: Mage has correct base stats
    When I create or examine a Mage character
    Then the base health should be: low (around 60-80)
    And the base mana should be: high (around 100-130)
    And the base strength should be: low (around 3-6)
    And the base agility should be: medium (around 8-12)
    And the base intelligence should be: high (around 15-20)
    And the base wisdom should be: medium (around 10-15)

  Scenario: Mage can equip magical weapons
    When I check Mage equipment restrictions
    Then Mages should be able to equip: staffs, wands, robes
    And Mages should NOT be able to equip: heavy armor, shields, axes

  Scenario: Mage combat style description
    When I examine the Mage class
    Then the combat style should be: "ranged" or "magical"
    And the role should be: "damage dealer" or "caster"

  # Rogue Class
  Scenario: Rogue has correct base stats
    When I create or examine a Rogue character
    Then the base health should be: medium (around 90-110)
    And the base mana should be: medium (around 50-70)
    And the base strength should be: medium (around 10-14)
    And the base agility should be: high (around 15-20)
    And the base intelligence should be: medium (around 8-12)
    And the base wisdom should be: low (around 5-8)

  Scenario: Rogue can equip light weapons and stealth gear
    When I check Rogue equipment restrictions
    Then Rogues should be able to equip: daggers, bows, leather armor
    And Rogues should NOT be able to equip: heavy armor, two-handed weapons

  Scenario: Rogue combat style description
    When I examine the Rogue class
    Then the combat style should be: "stealth" or "melee"
    And the role should be: "damage dealer" or "scout"

  # Priest Class
  Scenario: Priest has correct base stats
    When I create or examine a Priest character
    Then the base health should be: medium (around 80-100)
    And the base mana should be: high (around 90-120)
    And the base strength should be: low (around 4-7)
    And the base agility should be: low (around 5-8)
    And the base intelligence should be: medium (around 10-14)
    And the base wisdom should be: high (around 15-20)

  Scenario: Priest can equip holy items and light armor
    When I check Priest equipment restrictions
    Then Priests should be able to equip: maces, holy staffs, robes
    And Priests should NOT be able to equip: heavy armor, swords, axes

  Scenario: Priest combat style description
    When I examine the Priest class
    Then the combat style should be: "healer" or "support"
    And the role should be: "healer" or "buffer"

  # Class Selection
  Scenario: Class is selected during character creation
    Given I am creating a new character
    When I select class "Warrior"
    Then the character should have Warrior base stats
    And the character should have the correct equipment restrictions
    And the character should be able to use Warrior-specific skills

  Scenario: Class cannot be changed after character creation
    Given I have a character with class "Warrior"
    When I attempt to change the class to "Mage"
    Then the operation should be denied
    Or the character should need to be recreated

  Scenario: Invalid class name is rejected
    Given I am creating a character
    When I attempt to select class "SuperClass"
    Then the response should be a validation error
    And the error should list valid class options

  # Class Talents/Skills (Chef subclass from PR #213)
  Scenario: Class has associated talents
    When I examine a character's class
    Then I should see a list of available talents for that class

  Scenario: Chef class (special) has pizza-themed talents
    Given the Chef class is implemented from PR #213
    When I examine the Chef class talents
    Then there should be pizza-themed combat abilities
    And the class should follow the same stat structure as other classes
