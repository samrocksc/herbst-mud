Feature: Admin Ability CRUD
  As a game administrator
  I want to create, read, update, and delete abilities via the admin API
  So that I can define combat skills for players and NPCs

  Scenario: Create an active combat ability
    Given I am authenticated as an admin
    When I POST to "/api/abilities" with:
      | name            | Slash              |
      | description     | A quick blade strike |
      | ability_type    | combat             |
      | ability_class   | active             |
      | cost            | 0                  |
      | mana_cost       | 0                  |
      | stamina_cost    | 10                 |
      | cooldown_seconds  | 3                |
    Then the response status should be 201

  Scenario: Create a passive ability
    Given I am authenticated as an admin
    When I POST to "/api/abilities" with:
      | name            | Toughness          |
      | ability_type    | defensive          |
      | ability_class   | passive            |
      | requirements    | 1                  |
    Then the response status should be 201

  Scenario: List abilities
    Given I am authenticated as an admin
    When I GET "/api/abilities"
    Then the response status should be 200

  Scenario: Update an ability
    Given I am authenticated as an admin
    And an ability "Slash" exists
    When I PUT to "/api/abilities/{id}" with:
      | name | Power Slash      |
      | stamina_cost | 15 |
    Then the response status should be 200
    And the ability name should be "Power Slash"

  Scenario: Delete an ability
    Given I am authenticated as an admin
    And an ability "Slash" exists
    When I DELETE "/api/abilities/{id}"
    Then the response status should be 200

  Scenario: Reject ability with empty name
    Given I am authenticated as an admin
    When I POST to "/api/abilities" with:
      | name |                |
    Then the response status should be 400
