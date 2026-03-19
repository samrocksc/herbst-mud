Feature: Item Description Schema
  As a developer
  I want the item schema to support all look/examine fields
  So that the system can display rich, contextual item information

  Background:
    Given the item schema is properly defined

  Scenario: Items have multiple description fields
    Given an item "rusty_pipe" exists
    Then it has a name field: "rusty_pipe"
    And it has a short_desc field (nullable)
    And it has a description field
    And it has an examine_desc field (nullable, for detailed examine)

  Scenario: short_desc used by look command
    Given an item has short_desc: "A rusty metal pipe"
    When the room displays with look command
    Then the short_desc is shown (not full description)

  Scenario: description used by look at command
    Given an item has description: "A rusted metal pipe, about two feet long"
    When I type "look at rusty_pipe"
    Then the description is shown

  Scenario: examine_desc used by examine command
    Given an item has examine_desc: "Detailed rust patterns indicate this was a quality weapon once..."
    When I type "examine rusty_pipe"
    Then the examine_desc is shown (not short_desc)

  Scenario: examine_desc falls back to description
    Given an item has no examine_desc
    But it has a description
    When I examine the item
    Then the description is shown as fallback

  Scenario: Item type enum is recognized
    Given an item has type: "weapon"
    When the item is displayed
    Then it is recognized as a weapon type
    And other types (armor, consumable, quest, misc) also work

  Scenario: is_immovable flag is respected
    Given an item has is_immovable: true
    When I try to take it
    Then the take fails with appropriate message
    And the item remains in place

  Scenario: is_visible controls room display
    Given an item has is_visible: false
    When I look at the room
    Then the item is not shown in the room listing
    When I have the right conditions
    Then the item can be revealed

  Scenario: is_container flag identifies containers
    Given an item has is_container: true
    When I interact with it
    Then container commands work on it

  Scenario: is_readable flag identifies readable items
    Given an item has is_readable: true
    When I try to read it
    Then the read command works

  Scenario: content field stores readable text
    Given an item has is_readable: true
    And it has content: "Some readable text..."
    When I read the item
    Then the content is displayed

  Scenario: hidden_details JSON field exists
    Given an item has hidden_details JSON array
    When I examine the item
    Then hidden details are revealed based on skill level

  Scenario: on_examine JSON field for event triggers
    Given an item has on_examine event triggers
    When I examine the item at correct level
    Then events fire as defined

  Scenario: color field controls display color
    Given an item has color field set
    When the item is displayed
    Then the item name is shown in that color

  Scenario: Item schema supports all required fields
    When validating the Equipment schema
    Then all required fields are present:
      | field           | type    |
      | name            | string  |
      | shortDesc       | string  |
      | description     | text    |
      | examineDesc     | text    |
      | type            | enum    |
      | isImmvable      | boolean |
      | isVisible       | boolean |
      | isContainer     | boolean |
      | isReadable      | boolean |
      | content         | text    |
      | hiddenDetails   | JSON    |
      | onExamine       | JSON    |
      | color           | string  |
