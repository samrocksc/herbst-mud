Feature: Admin Effects CRUD
  As a game administrator
  I want to create, read, update, and delete effects via the admin API
  So that I can define combat and game mechanics

  Scenario: Create a damage effect
    Given I am authenticated as an admin
    When I POST to "/api/effects" with:
      | name        | Fireball Blast     |
      | effect_type | hp_change          |
      | parameters  | {"amount": -50}    |
      | stack_mode  | replace            |
    Then the response status should be 201

  Scenario: Create a heal effect
    Given I am authenticated as an admin
    When I POST to "/api/effects" with:
      | name        | Minor Heal      |
      | effect_type | hp_change       |
      | parameters  | {"amount": 25}  |
      | stack_mode  | refresh         |
    Then the response status should be 201

  Scenario: Create a message effect
    Given I am authenticated as an admin
    When I POST to "/api/effects" with:
      | name        | Victory Fanfare     |
      | effect_type | message             |
      | parameters  | {"text": "You won!", "message_type": "combat"} |
    Then the response status should be 201

  Scenario: List effects
    Given I am authenticated as an admin
    When I GET "/api/effects"
    Then the response status should be 200
    And the response should contain the created effects

  Scenario: Update an effect
    Given I am authenticated as an admin
    And an effect "Fireball Blast" exists
    When I PUT to "/api/effects/{id}" with:
      | parameters | {"amount": -75} |
    Then the response status should be 200

  Scenario: Delete an effect
    Given I am authenticated as an admin
    And an effect "Fireball Blast" exists
    When I DELETE "/api/effects/{id}"
    Then the response status should be 200
    And the effect should no longer exist

  Scenario: Link effect to ability
    Given I am authenticated as an admin
    And an effect "Minor Heal" exists
    And an ability "Heal" exists
    When I POST to "/api/abilities/{ability_id}/effects" with:
      | effect_id   | {effect_id} |
      | effect_type | hp_change   |
      | target      | ally        |
      | value       | 25          |
    Then the response status should be 201

  Scenario: Reject effect with missing name
    Given I am authenticated as an admin
    When I POST to "/api/effects" with:
      | name        |                |
      | effect_type | hp_change      |
    Then the response status should be 400
