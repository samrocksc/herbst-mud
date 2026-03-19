Feature: Look - Container System (look-06)
  As a player
  I want to use containers to organize my inventory
  So that I can carry more items and keep my inventory tidy

  Background:
    Given the server is running
    And I am logged in as a test character
    And I have items in my inventory

  @donnie
  Scenario: View container contents
    Given I have a bag container in my inventory
    When I type "open bag" or "look in bag"
    Then I should see the contents of the container
    And each item should be listed with name and quantity

  @donnie
  Scenario: Container has weight limit
    Given I have a container with capacity
    When I try to add items beyond capacity
    Then I should see a "container full" message
    And the item should remain in my inventory

  @donnie
  Scenario: Put item in container
    Given I have a container and an item
    When I type "put [item] in [container]"
    Then the item should move to the container
    And I should see a success message

  @donnie
  Scenario: Take item from container
    Given I have items stored in a container
    When I type "take [item] from [container]"
    Then the item should move to my inventory
    And the container should have one less item

  @donnie
  Scenario: Container shows in inventory
    Given I own a container
    When I type "inventory" or "i"
    Then I should see the container listed
    And it should show how many items it contains

  @donnie
  Scenario: Nested containers
    Given I have a container inside another container
    When I open the outer container
    Then I should see the inner container listed

  @donnie
  Scenario: Container item icon display
    Given I have a container in my inventory
    When I view my inventory
    Then I should see a bag/backpack icon indicator
    And the container should show "🎒" or similar icon

  @donnie
  Scenario: Drop container to room
    Given I have a container in my inventory
    When I type "drop [container]"
    Then the container should appear on the floor
    And items inside should remain in the container

  @donnie
  Scenario: Get container from room
    Given a container is on the floor
    When I type "get [container]"
    Then the container should be in my inventory
    And all items inside should come with it