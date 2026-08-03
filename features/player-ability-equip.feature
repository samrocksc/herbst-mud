Feature: Player Web Client — Classless abilities are seeded and equippable

  Background:
    Given I am authenticated as sma
    And I have a character named "smack" in Ooze Surfers
    And smack has at least one empty skill slot

  Scenario: Default classless abilities are seeded on first boot
    Given the server is freshly started
    And the abilities table is empty
    When the server starts up
    Then the abilities table should contain at least 5 rows
    And the abilities table should include an ability named "Haymaker"
    And the abilities table should include an ability named "Concentrate"
    And the abilities table should include an ability named "Back-off"
    And the abilities table should include an ability named "Scream"
    And the abilities table should include an ability named "Slap"
    And each seeded ability should have at least one ability_effect

  Scenario: GET /abilities/classless returns the seeded abilities
    Given the server is running with default seeds applied
    When I GET "/abilities/classless"
    Then the response should include an "abilities" array
    And the array should have at least 5 entries
    And each entry should have a non-empty "name" field

  Scenario: Reproduce the regression (before fix)
    Given the server is freshly started
    And the abilities table is empty
    And InitSkillsAndAbilities is not called from main.go
    When the server starts up
    Then the abilities table should still be empty
    And GET /abilities/classless should return an empty array
    And the player's skill slots should all render as empty in the UI

  Scenario: Admin wipe re-seeds abilities
    Given the server is running with default seeds applied
    When an admin POSTs to "/admin/wipe" with wipe_skills=true
    Then the abilities table should still contain at least 5 rows
    And the server log should record a re-seed event

  Scenario: Player can equip a classless ability after seeding
    Given the server is running with default seeds applied
    And smack is in Ooze Surfers
    When I open the equipment screen as smack
    Then the available abilities list should include "Haymaker"
    When I equip "Haymaker" in slot 1
    Then GET /characters/5/abilities slot 1 should have name "Haymaker"
    And smack's HUD should show "1: Haymaker" in the skill bar

  Scenario: Dangling character_abilities rows are cleaned up
    Given smack has character_abilities rows with NULL ability_characters
    When the cleanup migration runs
    Then no character_abilities rows for smack should have NULL ability_characters
    And smack should have 4 empty skill slots available
