Feature: Admin Race CRUD with world scoping
  As a game administrator
  I want to create races that are scoped to the world I picked in the UI
  So that races don't silently land in the wrong world and trip a unique-constraint error

  Background:
    Given I am authenticated as an admin

  Scenario: Creating a duplicate race in the same world returns a friendly conflict
    Given a race "Human" exists in world "1"
    When I POST to "/api/races?world_id=1" with:
      | name        | Human     |
      | display_name| Human     |
      | description | Duplicate |
    Then the response status should be 409
    And the response should contain "already exists"

  Scenario: Creating a race in a different world succeeds even if a same-named race exists in world 1
    Given a race "Human" exists in world "1"
    When I POST to "/api/races?world_id=2" with:
      | name        | Human     |
      | display_name| Human     |
      | description | Same name, different world |
    Then the response status should be 201
    And the response should contain "world_id"
    And a race "Human" should exist in world "2"

  Scenario: Creating a race with world_id "default" is normalized to world 1
    Given a race "Human" exists in world "1"
    When I POST to "/api/races?world_id=default" with:
      | name        | Human     |
      | display_name| Human     |
      | description | Default sentinel |
    Then the response status should be 409
    And the response should contain "already exists"

  Scenario: A new race in world 2 is actually stored in world 2, not silently leaked to world 1
    When I POST to "/api/races?world_id=2" with:
      | name        | Turtle    |
      | display_name| Turtle    |
      | description | Shell-back |
    Then the response status should be 201
    And the stored race "Turtle" should have world_id "2"
    And a race "Turtle" should not exist in world "1"

  Scenario: List races in world 1 does not include races created in world 2
    Given a race "Elf" exists in world "2"
    When I GET "/api/races?world_id=1"
    Then the response status should be 200
    And the response should not contain "Elf"
