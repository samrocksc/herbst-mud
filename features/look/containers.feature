🔴 Feature: Containers - Issue Look-06
  As a player
  I want to look inside containers and take items from them
  So that I can interact with storage and loot

  Background:
    Given the player is in a room
    And a container exists (chest, barrel, bag)

  Scenario: Container can be opened and closed
    Given a closed container exists
    When the player types "open <container>"
    Then the container should be marked as open
    And the player can then look inside
    And the player can type "close <container>" to close it

  Scenario: Look inside container shows items
    Given an open container exists
    When the player types "look in <container>" or "open <container>"
    Then the items inside should be listed
    And item names and possibly quantities should be shown

  Scenario: Take item from container
    Given an open container has items
    When the player types "take <item> from <container>"
    Then the item should be moved to player inventory
    And the container's contents should update
    And a message should confirm the take action

  Scenario: Put item in container
    Given a container is open
    And the player has an item in inventory
    When the player types "put <item> in <container>"
    Then the item should move from inventory to container
    And the container contents should update

  Scenario: Cannot take from closed container
    Given a closed container exists
    When the player types "take <item> from <container>"
    Then an error should indicate the container is closed
    And the player should be prompted to open it first

  Scenario: Container has a capacity limit
    Given a container has limited capacity
    When too many items are added
    Then some items should not fit
    And a message should indicate the container is full
