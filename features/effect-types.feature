Feature: All Effect Types
  As a game developer
  I want to test all effect types available in the system
  So that I can verify they work correctly

  # XP-related effects
  Scenario: XP Gain effect
    Given I am authenticated as an admin
    When I POST to "/api/effects" with:
      | name        | Small XP Gain    |
      | effect_type | xp_gain          |
      | parameters  | {"amount": 100}  |
    Then the response status should be 201

  Scenario: XP Set effect
    Given I am authenticated as an admin
    When I POST to "/api/effects" with:
      | name        | Reset XP         |
      | effect_type | xp_set           |
      | parameters  | {"amount": 0}    |
    Then the response status should be 201

  Scenario: XP Multiplier effect
    Given I am authenticated as an admin
    When I POST to "/api/effects" with:
      | name        | Double XP        |
      | effect_type | xp_multiplier    |
      | parameters  | {"multiplier": 2.0, "duration_secs": 300} |
    Then the response status should be 201

  # Character stat effects
  Scenario: HP Change effect
    Given I am authenticated as an admin
    When I POST to "/api/effects" with:
      | name        | Minor Heal       |
      | effect_type | hp_change        |
      | parameters  | {"amount": 25}   |
    Then the response status should be 201

  Scenario: Stamina Change effect
    Given I am authenticated as an admin
    When I POST to "/api/effects" with:
      | name        | Stamina Boost    |
      | effect_type | stamina_change   |
      | parameters  | {"amount": 50}   |
    Then the response status should be 201

  Scenario: Mana Change effect
    Given I am authenticated as an admin
    When I POST to "/api/effects" with:
      | name        | Mana Restoration |
      | effect_type | mana_change      |
      | parameters  | {"amount": 30}   |
    Then the response status should be 201

  # Utility effects
  Scenario: Bind Point Set effect
    Given I am authenticated as an admin
    When I POST to "/api/effects" with:
      | name        | Set Bind Point   |
      | effect_type | bind_point_set   |
      | parameters  | {"room_id": 1}   |
    Then the response status should be 201

  Scenario: Teleport effect
    Given I am authenticated as an admin
    When I POST to "/api/effects" with:
      | name        | Quick Teleport   |
      | effect_type | teleport         |
      | parameters  | {"room_id": 10}  |
    Then the response status should be 201

  # Message effects
  Scenario: Message effect
    Given I am authenticated as an admin
    When I POST to "/api/effects" with:
      | name        | Victory Message  |
      | effect_type | message          |
      | parameters  | {"text": "Victory!", "message_type": "info"} |
    Then the response status should be 201

  # Tag effects
  Scenario: Tag Add effect
    Given I am authenticated as an admin
    When I POST to "/api/effects" with:
      | name        | Give Title       |
      | effect_type | tag_add          |
      | parameters  | {"tag_name": "hero"} |
    Then the response status should be 201

  Scenario: Tag Remove effect
    Given I am authenticated as an admin
    When I POST to "/api/effects" with:
      | name        | Remove Tag       |
      | effect_type | tag_remove       |
      | parameters  | {"tag_name": "visitor"} |
    Then the response status should be 201

  # Nested effects
  Scenario: Apply Effect (nested) effect
    Given an effect "Heal" exists with type hp_change
    When I POST to "/api/effects" with:
      | name        | Chain Heal       |
      | effect_type | apply_effect     |
      | parameters  | {"effect_id": {heal_id}} |
    Then the response status should be 201
