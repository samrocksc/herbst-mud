Feature: Admin NPC Template CRUD
  As a game administrator
  I want to create, read, update, and delete NPC templates via the admin API
  So that I can manage the game's NPC population

  Scenario: Create a new NPC template
    Given I am authenticated as an admin
    When I POST to "/api/npc-templates" with:
      | name      | Goblin Scout        |
      | level     | 3                   |
      | xp_value  | 150                 |
      | disposition | hostile           |
    Then the response status should be 201
    And the response should contain the created NPC template

  Scenario: List NPC templates
    Given I am authenticated as an admin
    When I GET "/api/npc-templates"
    Then the response status should be 200
    And the response should be a JSON array

  Scenario: Update an NPC template
    Given I am authenticated as an admin
    And an NPC template "Goblin Scout" exists
    When I PUT to "/api/npc-templates/{id}" with:
      | name | Goblin Scout (Elite) |
      | level | 5                   |
    Then the response status should be 200
    And the NPC template name should be "Goblin Scout (Elite)"

  Scenario: Delete an NPC template
    Given I am authenticated as an admin
    And an NPC template "Goblin Scout" exists
    When I DELETE "/api/npc-templates/{id}"
    Then the response status should be 200
    And the NPC template should no longer exist

  Scenario: Reject NPC template with empty name
    Given I am authenticated as an admin
    When I POST to "/api/npc-templates" with:
      | name |                     |
      | level | 1                  |
      | xp_value | 0               |
    Then the response status should be 400
