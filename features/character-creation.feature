Feature: Character Creation
  As a new player
  I want to create a character through an interactive form
  So that I can enter the game world with a unique identity

  Background:
    Given the MUD server is running
    And I am a logged-in user with no existing characters

  Scenario: New user sees character creation form
    Given I have logged in successfully
    And I have no existing characters
    When I reach the character selection screen
    Then I should be presented with a character creation form
    And I should be guided through the character creation process

  Scenario: Character creation - Fountain awakening intro
    Given I start character creation
    When the fountain awakening flow begins
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
    And character creation should proceed to name entry

  Scenario: Character creation - Enter name
    Given I have washed my face at the fountain
    And character creation has started
    When I am prompted for my name
    And I type "Michelangelo"
    Then my character name should be set to "Michelangelo"
    And I should proceed to the race selection

  Scenario: Character creation - Race selection
    Given I have entered my character name
    And I am on the race selection screen
    When I view the available races
    Then I should see "human" option with +1 all stats
    And I should see "turtle" option with +CON and Innate Block
    And I should see "rabbit" option
    And I should see "rat" option
    And I should see "rhino" option
    And each race should display its stat bonuses and penalties

  Scenario: Character creation - Select class
    Given I have selected my race
    And I am on the class selection screen
    When I see the available classes
    Then I should see "Warrior" option
    And I should see "Chef" option (locked until unlocked)
    And I should see "Mystic" option (locked until unlocked)
    And selecting a class should proceed to the next step

  Scenario: Character creation - Select gender
    Given I have selected my class
    And I am on the gender selection screen
    When I see the available genders
    Then I should see "he/him" option
    And I should see "she/her" option
    And I should see "it/its" option
    And I should see "they/them" option
    And I should be able to select one

  Scenario: Character creation - Select size
    Given I have selected my gender
    And I am on the size selection screen
    When I see the available sizes
    Then I should see size options that affect combat
    And I should see size options that affect weight for crash skill

  Scenario: Character creation - Complete and save
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

  Scenario: Character name must be unique
    Given I am at the character name prompt
    When I enter a name that is already taken
    Then I should see an error message
    And I should be prompted to enter a different name

  Scenario: Character name cannot be empty
    Given I am at the character name prompt
    When I enter a blank or empty name
    Then I should see an error message
    And I should be prompted to enter a valid name
