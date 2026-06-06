Feature: Admin Item Template CRUD
  As a game administrator
  I want to create, read, update, and delete item templates via the admin API
  So that I can stock the game world with equipment

  Scenario: Create a weapon item template
    Given I am authenticated as an admin
    When I POST to "/api/equipment-templates" with:
      | name               | Iron Sword       |
      | slot               | main_hand        |
      | item_type          | weapon           |
      | level              | 1                |
      | damage_dice_count  | 1                |
      | damage_dice_sides  | 8                |
      | damage_bonus       | 1                |
      | damage_type        | slashing         |
      | weapon_type        | sword            |
    Then the response status should be 201

  Scenario: Create an armor item template
    Given I am authenticated as an admin
    When I POST to "/api/equipment-templates" with:
      | name          | Leather Vest      |
      | slot          | chest             |
      | item_type     | armor             |
      | armor_rating  | 3                 |
      | armor_type    | light             |
    Then the response status should be 201

  Scenario: List item templates
    Given I am authenticated as an admin
    When I GET "/api/equipment-templates"
    Then the response status should be 200

  Scenario: Update an item template
    Given I am authenticated as an admin
    And an item template "Iron Sword" exists
    When I PUT to "/api/equipment-templates/{id}" with:
      | name | Iron Sword +1    |
      | level | 3               |
    Then the response status should be 200

  Scenario: Delete an item template
    Given I am authenticated as an admin
    And an item template "Iron Sword" exists
    When I DELETE "/api/equipment-templates/{id}"
    Then the response status should be 200

  Scenario: Reject item template with empty name
    Given I am authenticated as an admin
    When I POST to "/api/equipment-templates" with:
      | name |                   |
      | item_type | weapon      |
    Then the response status should be 400
