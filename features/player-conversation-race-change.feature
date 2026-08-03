Feature: Player Web Client — Conversation race-change wizard

  Background:
    Given I am authenticated as sma
    And I have a character named "smack" in Ooze Surfers
    And the NPC "Theodore Von Rad" is in the same room as "smack"

  Scenario: Talk to Theodore and change race to Faerie
    Given I am playing as "smack"
    When I type "talk Theodore Von Rad"
    Then the conversation overlay should show Theodore's greeting
    When I choose "I want to change my race"
    Then I should see a list of races including "Faerie"
    When I choose "Faerie"
    Then I should see Faerie lore and options "Change my race", "Back", and "Leave"
    When I choose "Change my race"
    Then I should see "Done! You are now a Faerie."
    And the server should update "smack" race to "faerie"

  Scenario: Conversation branches do not duplicate screens
    Given I am playing as "smack"
    When I talk to Theodore and walk every branch back to the entry node
    Then the conversation overlay should never show two "Leave" buttons at once
    And leaving the conversation should return to the room screen

  Scenario: Race change is idempotent
    Given I am playing as "smack"
    And "smack" is already race "faerie"
    When I talk to Theodore and choose "Faerie" → "Change my race"
    Then I should see "Done! You are now a Faerie."
    And the server should still report "smack" race as "faerie"

  Scenario: Conversation keyboard shortcuts work while command input is focused
    Given I am playing as "smack"
    When I type "talk Theodore Von Rad"
    Then the conversation overlay should show Theodore's greeting
    When I press the "1" key without manually blurring the command input
    Then the conversation should advance to "Tell me about this world"
    When I press "Escape"
    Then the conversation overlay should close
