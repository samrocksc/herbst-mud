Feature: Player Class-Specific Abilities

  Background:
    Given I am authenticated as sma
    And I have a character named "smack" in Ooze Surfers
    And my character class is "trash_mage"

  Scenario: Class-specific abilities are visible in the equip panel
    Given I am playing as "smack"
    When I open the ability equip panel
    Then I should see "Trash Bolt" as an available ability
    And I should see "Junk Shield" as an available ability
    And I should see "Putrid Spray" as an available ability
    And I should see "Salvage Aura" as an available ability

  Scenario: Trash Mage abilities are eligible for a trash_mage character
    Given I am playing as "smack"
    When I fetch my character skills from the API
    Then the faction_abilities for "trash_mage" should all be eligible
    And the faction_abilities for "foot_clank" should NOT be eligible

  Scenario: Equipping a class-specific ability works
    Given I am playing as "smack"
    And the ability equip panel is open
    When I equip "Trash Bolt" to an available slot
    Then the slot should show "Trash Bolt"
    And the server log should show POST /characters/5/abilities with status 201

  Scenario: Foot Clank abilities are not eligible for a trash_mage character
    Given I am playing as "smack"
    When I fetch my character skills from the API
    Then "Mech Blade Slash" should have eligible=false
    And "Cloak Field" should have eligible=false
    And "Servo Stomp" should have eligible=false
    And "System Reboot" should have eligible=false

  Scenario: No duplicated class categories exist
    Given I query the faction_categories table
    Then there should be exactly one category named "class" for world_id "2"
    And there should be no category named "Classes" for world_id "2"
    And there should be no category named "Professions" for world_id "2"