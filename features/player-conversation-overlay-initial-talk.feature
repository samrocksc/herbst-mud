Feature: Player Web Client — Conversation overlay opens on initial talk

  Background:
    Given I am authenticated as sma
    And I have a character named "smack" in Ooze Surfers
    And the NPC "Theodore Von Rad" is in the same room as "smack"
    And the NPC template for "Theodore Von Rad" has at least one dialog node

  Scenario: Initial talk opens the conversation overlay
    Given I am playing as "smack"
    When I type "talk Theodore Von Rad"
    Then the conversation overlay should show Theodore's greeting
    And the conversation overlay should show at least one numbered response

  Scenario: Talking to an NPC without dialog nodes keeps the text fallback
    Given I am playing as "smack"
    And the NPC "Some Text-Only NPC" is in the same room as "smack"
    And the NPC template for "Some Text-Only NPC" has zero dialog nodes
    When I type "talk Some Text-Only NPC"
    Then the scrollback should contain "Some Text-Only NPC says:"
    And the conversation overlay should not be visible

  Scenario: Server sends a screen payload with view_type=conversation on initial talk
    Given I am playing as "smack"
    When I type "talk Theodore Von Rad"
    Then the WebSocket should receive a message of type "screen"
    And the screen payload should have view_type "conversation"
    And the screen payload should include npc_name "Theodore Von Rad"
    And the screen payload should include current_node_id pointing to a known dialog node

  Scenario: Reproduce the regression (before fix)
    Given I am playing as "smack"
    When I type "talk Theodore Von Rad"
    Then the conversation overlay should NOT open
    And the scrollback should contain "Theodore Von Rad says:"
    And the WebSocket should NOT receive a "screen" message with view_type "conversation"

  Scenario: Dialog choice advances to the next node
    Given I am playing as "smack"
    And the conversation overlay is open with Theodore Von Rad
    When I click response "1" ("Tell me about this world")
    Then the conversation overlay should update with new NPC text
    And the conversation overlay should show response options for the next node

  Scenario: Dialog command routes to handleDialogChoice
    Given I am playing as "smack"
    And the conversation overlay is open with Theodore Von Rad
    When I send the command "dialog theodore_von_rad theodore_entry 1"
    Then the WebSocket should receive a "screen" message with view_type "conversation"
    And the screen payload should include current_node_id "theodore_about_world"

  Scenario: Dialog choice with empty next_node_id ends the conversation
    Given I am playing as "smack"
    And the conversation overlay is open with Theodore Von Rad
    When I choose a response whose next_node_id is empty
    Then the scrollback should contain "You end the conversation."
    And the conversation overlay should close

  Scenario: Leaving conversation then moving does not reopen the overlay
    Given I am playing as "smack"
    And the conversation overlay is open with Theodore Von Rad
    When I click "Leave"
    Then the conversation overlay should close
    And the conversation state should be cleared
    When I type "e" to move east
    Then the room screen should render for the destination room
    And the conversation overlay should NOT reopen
