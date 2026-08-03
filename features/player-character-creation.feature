Feature: Player Web Client — Character Creation with Class Selection

  Background:
    Given I am authenticated as sma
    And I have selected the "Ooze Surfers" world
    And I have 0 characters in Ooze Surfers

  Scenario: Create character with class selection (the fix)
    Given I navigate to the character creation page
    When I fill the form with:
      | name   | class  | race  | gender |
      | Ninja  | ninja  | Ooze  | It     |
    And I click "Create Character"
    Then the request status should be 201
    And the created character has class "ninja"
    And the character is in world "2"

  Scenario: Default class is survivor when no class chosen (regression guard for backward compat)
    Given I navigate to the character creation page
    When I fill the form with:
      | name  | race  | gender |
      | TestA | Ooze  | It     |
    And I click "Create Character"
    Then the created character has class "survivor"

  Scenario: Class dropdown is populated from /classes endpoint
    Given I navigate to the character creation page
    Then the Class dropdown contains at least:
      | warrior |
      | chef    |
      | mystic  |
      | tinkerer|
      | survivor|
      | ninja   |

  Scenario: /classes endpoint returns the full class list
    Given I am authenticated as sma
    When I GET "/classes"
    Then the response status should be 200
    And the response is a JSON array of length >= 5

  Scenario: /classes/ninja/specialties returns 200 with specialties (or empty array)
    Given I am authenticated as sma
    When I GET "/classes/ninja/specialties"
    Then the response status should be 200
    And the response is a JSON array
