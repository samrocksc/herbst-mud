Feature: Readable Items
  As a player
  I want to read items that contain text
  So that I can discover lore, instructions, and hidden information

  Background:
    Given I am at a location with a "survival_manual" and "encrypted_terminal"

  Scenario: Read item with text content
    Given the "survival_manual" has content:
      """
      SURVIVAL IN THE NEW WORLD
      A Guide for the Post-Ooze Era

      Contents:
        1. Finding Water
        2. Identifying Safe Food
        3. Avoiding Ooze Pools...
      """
    When I type "read survival_manual"
    Then the full text content is displayed

  Scenario: Read item that is not readable
    Given "rusty_pipe" has no readable content
    When I type "read rusty_pipe"
    Then I see error: "the rusty_pipe is not readable"

  Scenario: Skill-gated item with sufficient skill
    Given the "encrypted_terminal" requires tech skill level 5
    And my character has tech skill level 7
    When I type "read encrypted_terminal"
    Then I see the decrypted content

  Scenario: Skill-gated item with insufficient skill
    Given the "encrypted_terminal" requires tech skill level 5
    And my character has tech skill level 3
    When I type "read encrypted_terminal"
    Then I see "[Requires tech skill level 5 to decode]"
    And I see "(You have tech skill level 3. Cannot read.)"

  Scenario: Skill-gated item at exact threshold
    Given the "encrypted_terminal" requires tech skill level 5
    And my character has tech skill level 5
    When I type "read encrypted_terminal"
    Then I can read the decrypted content

  Scenario: Read item not in room or inventory
    When I type "read nonexistent_item"
    Then I see error: "you don't see nonexistent_item here"

  Scenario: Read item in inventory
    Given "survival_manual" is in my inventory
    When I type "read survival_manual"
    Then the content is displayed
    (I can read items in inventory, not just in room)

  Scenario: Readable item shows skill requirement in description
    Given the "encrypted_terminal" requires tech skill level 5
    When I examine the terminal
    Then I see "[Encrypted - requires tech skill level 5]"
    And the skill requirement is visible before reading
