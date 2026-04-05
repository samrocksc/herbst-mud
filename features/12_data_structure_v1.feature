🔴 Feature: Data Structure V1
  As a game developer
  I want proper data structures for the game entities
  So that the game has a solid foundation for all game mechanics

  Background:
    Given the database is connected
    And the ent schema is generated

  Scenario: User entity structure
    Given I create a new user
    Then the user should have the following fields:
      | field | type |
      | id | int |
      | email | string |
      | password | string |
      | is_admin | boolean |
      | characters | []Character |
    And the user should have a one-to-many relationship with characters
    And the user can have up to three characters
    And the user can have god mode (unkillable)

  Scenario: Character entity structure
    Given I create a new character
    Then the character should have the following fields:
      | field | type |
      | id | int |
      | name | string |
      | gender | string |
      | description | string |
      | isNPC | boolean |
      | currentRoomId | int |
      | level | int |
      | experience | int |
      | userId | int |
    And the character should have stats:
      | stat |
      | Strength |
      | Intellect |
      | Wisdom |
      | Dexterity |
      | Fortitude |
    And the character can have equipment
    And the character has a class (Warrior, Magician, Thief, Charlatan, Vigilante)
    And the character has a race (Mutant, Human, Animal)

  Scenario: Room entity structure
    Given I create a new room
    Then the room should have the following fields:
      | field | type |
      | id | int |
      | name | string |
      | description | string |
      | exits | map[string]int |
      | isStartingRoom | boolean |
      | is_peerable | boolean |
    And the room can hold characters
    And the room can hold equipment
    And the room has atmosphere (air/water/wind)

  Scenario: Class system
    Given I create a character with a class
    Then the character can have one of the following classes:
      | class |
      | Warrior |
      | Magician |
      | Thief |
      | Charlatan |
      | Vigilante |

  Scenario: Race system
    Given I create a character with a race
    Then the character can have one of the following races:
      | race |
      | Mutant |
      | Human |
      | Animal |
    And if the race is Animal, it should have a species field

  Scenario: Gender system
    Given I create a character with a gender
    Then the character can have one of the following genders:
      | gender |
      | he/him |
      | she/her |
      | it/its |
      | they/them |
    And the gender field should be extensible

  Scenario: User entity validates required fields
    Given I attempt to create a user without an email
    Then the creation should fail
    And an error about required email should be returned
    When I attempt to create a user without a password
    Then the creation should fail
    And an error about required password should be returned

  Scenario: Character entity validates required fields
    Given I attempt to create a character without a name
    Then the creation should fail
    And an error about required name should be returned

  Scenario: Room entity has exits as key-value map
    Given I create a room with exits:
      | direction | destination |
      | north | 2 |
      | east | 3 |
    Then the room should have an exits map with 2 entries
    And exits["north"] should equal 2
    And exits["east"] should equal 3

  Scenario: NPC character has isNPC flag set
    Given I create an NPC character "Guard"
    Then the character should have isNPC: true
    And the character should not require a userId
    And the character should have _is_admin: true or admin flag

  Scenario: Character stats are integers
    Given I create a character
    Then Strength should be an integer
    And Intellect should be an integer
    And Wisdom should be an integer
    And Dexterity should be an integer
    And Fortitude should be an integer

  Scenario: Character experience accumulates
    Given a character "Hero" has 100 experience
    When the character gains 50 experience
    Then the character should have 150 experience

  Scenario: Character level increases at thresholds
    Given a character is at level 1
    When the character gains 1000 experience
    Then the character should be level 2
    And higher levels should require more experience