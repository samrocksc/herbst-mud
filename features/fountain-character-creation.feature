Feature: Fountain Character Creation
  As a new player
  I want to create a character through the fountain awakening flow
  So that I can enter the game world with a unique identity

  Background:
    Given I am a new player starting the game
    And I have no existing characters

  Scenario: Player wakes at fountain covered in mud
    Given I connect to the game for the first time
    When I arrive at the starting area
    Then I should see "The Stone Fountain" room
    And I should see a description mentioning I am covered in mud
    And I should be unable to see my reflection clearly

  Scenario: Player washes face at fountain
    Given I am at the Stone Fountain
    And I am covered in mud
    When I type "wash face"
    Then I should see a message about washing my face
    And the mud should be cleared
    And I should see my reflection in the water
    And character creation should begin

  Scenario: Character creation - Enter name
    Given I have washed my face at the fountain
    And the character creation form has started
    When I am prompted for my name
    And I type "Michelangelo"
    Then my character name should be set to "Michelangelo"
    And I should proceed to the next field

  Scenario: Character creation - Select race
    Given I have entered my character name
    And I am on the race selection screen
    When I see the available races
    Then I should see "human" option
    And I should see "turtle" option
    And I should see "rabbit" option
    And I should see "rat" option
    And I should see "rhino" option
    And each race should display stat bonuses and penalties

  Scenario: Character creation - Turtle race bonus displayed
    Given I am on the race selection screen
    When I view the turtle race option
    Then I should see "+CON bonus"
    And I should see "Innate Block ability"
    And I should see "-INT penalty"

  Scenario: Character creation - Select gender
    Given I have selected my race
    And I am on the gender selection screen
    When I see the available genders
    Then I should see "male" option
    And I should see "female" option
    And I should see "other" option

  Scenario: Character creation - Select class
    Given I have selected my gender
    And I am on the class selection screen
    When I see the available classes
    Then I should see "Warrior" option
    And I should see "Chef" option (locked until unlocked)
    And I should see "Mystic" option (locked until unlocked)
    And selecting Warrior should show available specialties

  Scenario: Character creation - Select size
    Given I have selected my class
    And I am on the size selection screen
    When I see the available sizes
    Then I should see size options that affect combat
    And I should see size options that affect weight for crash skill

  Scenario: Character creation - Complete and enter game
    Given I have filled all character creation fields
    When I confirm my character
    Then my character should be saved to the database
    And I should enter the game world
    And I should see my character's stats displayed
    And I should be at the Stone Fountain location

  Scenario: Character creation uses full terminal width
    Given I start character creation
    When the form is displayed
    Then it should use the full terminal width
    And it should use the full terminal height
    And the I/O style should match the login screen