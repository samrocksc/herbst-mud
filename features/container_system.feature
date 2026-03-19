Feature: Container System
  As a player
  I want to open containers to access stored items
  So that I can manage my inventory through chests and bags

  Scenario: Open an empty container
    Given I am in a room with an empty chest
    When I type "open chest"
    Then I should see "The chest is empty"

  Scenario: Open a container with items
    Given I am in a room with a chest containing "rusty sword"
    When I type "open chest"
    Then I should see "rusty sword" in the container contents

  Scenario: Take item from container
    Given I am in a room with a chest containing "rusty sword"
    And I have opened the chest
    When I type "take rusty sword from chest"
    Then I should have "rusty sword" in my inventory
    And the chest should no longer contain "rusty sword"

  Scenario: Put item in container
    Given I am in a room with a chest
    And I have "rusty sword" in my inventory
    When I type "put rusty sword in chest"
    Then I should not have "rusty sword" in my inventory
    And the chest should contain "rusty sword"

  Scenario: Cannot take item that doesn't exist in container
    Given I am in a room with an empty chest
    When I type "take nonexistent item from chest"
    Then I should see "There is no such item"

  Scenario: Locked container requires key
    Given I am in a room with a locked chest
    And I do not have the key
    When I type "open chest"
    Then I should see "The chest is locked"
    And I should need "chest key"

  Scenario: Unlocked container with key
    Given I am in a room with a locked chest
    And I have the "chest key" in my inventory
    When I type "open chest"
    Then the chest should open