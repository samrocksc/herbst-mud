Feature: Readable Items
  As a player
  I want to read books, signs, and scrolls
  So that I can learn about the game world

  Scenario: Read a book
    Given I have a "ancient tome" that contains readable text
    When I type "read ancient tome"
    Then I should see the book's content
    And I should see page navigation if multiple pages

  Scenario: Read a sign in the room
    Given I am in a room with a sign saying "Welcome to the Junkyard"
    When I enter the room
    Then the sign text should automatically display

  Scenario: Read an item without text
    Given I have a "rock" that has no readable text
    When I type "read rock"
    Then I should see "There is nothing to read"

  Scenario: Navigate book pages
    Given I have a "magic book" with 3 pages
    When I type "read magic book"
    And I type "turn page"
    Then I should see page 2
    When I type "turn page"
    Then I should see page 3
    And there should be no more pages after

  Scenario: Read a scroll
    Given I have a "scroll of fire" with spell instructions
    When I type "read scroll of fire"
    Then I should see the scroll's magical text