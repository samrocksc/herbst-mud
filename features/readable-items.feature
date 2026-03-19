Feature: Look - Readable Items (look-07)
  As a player
  I want to read books, scrolls, and other readable items
  So that I can learn about the game world and find hidden lore

  Background:
    Given the server is running
    And I am logged in as a test character

  @donnie
  Scenario: Read a book
    Given I am holding a readable book
    When I type "read [book]" or "read [book title]"
    Then the book's content should be displayed
    And the text should wrap to terminal width

  @donnie
  Scenario: Read a scroll
    Given I am holding a scroll
    When I type "read [scroll]"
    Then I should see the scroll text
    And I should be able to "recite" or "use" it if it's magic

  @donnie
  Scenario: Readable item shows in look
    Given I am holding a readable item
    When I type "inventory" or "look at [item]"
    Then I should see it has "(readable)" indicator

  @donnie
  Scenario: Read unknown item
    Given I try to read an item that is not readable
    Then I should see "You can't read that"
    And no text should be displayed

  @donnie
  Scenario: Read item with multiple pages
    Given I have a book with multiple pages
    When I type "read [book]" then "page [n]" or "turn page"
    Then I should see the requested page
    And I should see page numbers (e.g., "Page 2 of 5")

  @donnie
  Scenario: Read ancient/damaged text with skill check
    Given I have a damaged scroll requiring decipher
    And my "Literature" or "Ancient Languages" skill
    When I attempt to read it
    Then a skill check may be required
    And success reveals the full text

  @donnie
  Scenario: Quest book auto-read triggers
    Given I have a quest book
    When I read the first page
    Then a quest may be automatically offered
    Or the book content may be quest-gated

  @donnie
  Scenario: Examine shows read hint
    Given I examine a readable item
    Then I should see "Type 'read [item]' to read"
    Or similar hint about reading capability