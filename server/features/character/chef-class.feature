Feature: Chef Class - Pizza Combat
  As a player
  I want to play as a Chef with pizza combat abilities
  So that I can support my team with buffs while dealing ranged pizza damage

  Background:
    Given the system is running
    And the database is initialized with Chef skills and talents

  Scenario: Chef class is available
    Given a character is being created
    When the player selects class "chef"
    Then the character should be created with class "chef"
    And the character should have Chef-specific bonuses

  Scenario: Chef Pizzaiolo specialty
    Given a character with class "chef"
    When the player selects specialty "pizzaiolo"
    Then the character should have specialty "pizzaiolo"
    And the character should have access to pizza combat talents

  Scenario: Chef starting skills
    Given a character with class "chef" and specialty "pizzaiolo"
    Then the character should have skill "cooking" at level 1
    And the character should have skill "pizza_combat" at level 1
    And the character should have skill "foraging" at level 1

  Scenario: Chef starting pizza combat talents
    Given a character with class "chef" and specialty "pizzaiolo"
    Then the character should have access to talent "dough_ball"
    And the character should have access to talent "sauce_splash"
    And the character should have access to talent "pizza_cutter_dash"
    And the character should have access to talent "recipe_book"

  Scenario: Dough Ball attack
    Given a Chef character
    When the character uses talent "dough_ball"
    Then the attack should be a ranged attack
    And the attack should use flour as the projectile

  Scenario: Sauce Splash attack
    Given a Chef character
    When the character uses talent "sauce_splash"
    Then the attack should have a chance to blind the enemy

  Scenario: Pizza Cutter Dash attack
    Given a Chef character
    When the character uses talent "pizza_cutter_dash"
    Then the attack should hit all adjacent enemies
    And the attack should be a spin attack

  Scenario: Pizza Meteor ultimate attack
    Given a Chef character
    When the character uses talent "pizza_meteor"
    Then the attack should deal massive damage to one target

  Scenario: Cooking skills for buffs
    Given a Chef character
    When the character uses talent "recipe_book"
    Then the character should be able to learn new dishes

  Scenario: Foraging ability
    Given a Chef character
    When the character uses talent "forage"
    Then the character should find mutant ingredients

  Scenario: Pizza stall for passive income
    Given a Chef character
    When the character uses talent "open_pizza_stall"
    Then the character should generate passive income from selling pizza

  Scenario: Food Fight AoE attack
    Given a Chef character
    When the character uses talent "food_fight"
    Then the attack should affect all enemies in the room