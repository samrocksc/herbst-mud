🔴 Feature: Readable Items - Issue Look-07 / Issue #15
  As a player
  I want to read books, signs, and notes
  So that I can learn lore and find quest clues

  Background:
    Given the player has or encounters a readable item
    And readable items exist in the game world

  Scenario: Read command shows item text
    Given a readable item exists (scroll, book, sign)
    When the player types "read <item>"
    Then the text content should be displayed
    And the text should be readable in the output pane

  Scenario: Sign displays short message
    Given a sign is placed in a room
    When the player types "read sign"
    Then the sign's message should be displayed
    And the message should be brief (1-3 lines)

  Scenario: Book displays multi-page content
    Given a book exists with multiple pages
    When the player reads the book
    Then multi-page content should be navigable
    And the player can use "next page" / "prev page" commands

  Scenario: Cannot read non-readable items
    Given an item that is not readable
    When the player types "read <item>"
    Then an error message should indicate the item cannot be read
    And the player should not crash

  Scenario: Reading may trigger quest progress
    Given a readable item contains quest-relevant information
    When the player reads the item
    Then the quest log should update
    And a message should indicate quest progress

  Scenario: Skill check gates some readable content
    Given a book contains advanced arcane text
    And reading it requires a skill check
    When the player reads without the required skill
    Then some or all content should be garbled/missing
    And the player should know what skill is needed
