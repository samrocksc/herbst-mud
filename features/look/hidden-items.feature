Feature: Hidden Items and Reveal Conditions
  As a player
  I want items that are hidden until I discover them
  So that exploration and skill checks feel rewarding

  Background:
    Given the "fountain" is in the current room
    And a hidden "old_key" is associated with the fountain

  Scenario: Hidden item not shown in room description
    Given "old_key" has is_visible: false
    When I look at the room
    Then "old_key" is not shown in the room listing

  Scenario: Examine fountain reveals hidden key
    Given "old_key" has reveal_condition: { type: "examine", target: "fountain", min_examine_level: 50 }
    And my examine level is 50
    When I examine the fountain
    Then the hidden key is revealed
    And I see "[A small brass key falls out from the crack!]"
    And "old_key" is now visible in the room

  Scenario: Examine below threshold does not reveal
    Given the key requires examine level 50
    And my examine level is 40
    When I examine the fountain
    Then the key is not revealed
    And the room does not show the key

  Scenario: Revealed item persists in room
    Given I revealed "old_key" by examining fountain
    When I look at the room again
    Then "old_key" is still visible
    And the reveal condition is not re-triggered

  Scenario: Hidden item becomes takeable after reveal
    Given "old_key" was hidden
    And I revealed it by examining fountain
    When I type "take old_key"
    Then the key is added to my inventory

  Scenario: Perception check reveals hidden items
    Given "hidden_treasure" has reveal_condition: { type: "perception_check" }
    And my WIS stat allows the check
    When the check succeeds on room entry or look
    Then the hidden item is revealed

  Scenario: Use item reveals hidden items
    Given "secret_door" has reveal_condition: { type: "use_item", item: "old_key" }
    When I use "old_key" on "secret_door"
    Then the hidden passage is revealed

  Scenario: Event reveals hidden items
    Given "hidden_chamber" has reveal_condition: { type: "event" }
    When the story event fires
    Then the hidden chamber becomes visible

  Scenario: Reveal conditions are per-item, not global
    Given Character A revealed "old_key"
    When Character B enters the room
    Then "old_key" is still hidden for Character B

  Scenario: Reveal condition with no min level reveals immediately on examine
    Given "easy_secret" has reveal_condition: { type: "examine" } with no min level
    When I examine the target object
    Then the secret is revealed immediately

  Scenario: Multiple reveal conditions on same item
    Given "complex_item" has multiple reveal conditions
    When any condition is met
    Then the item becomes visible

  Scenario: Hidden item not takeable before reveal
    Given "old_key" is hidden
    When I try to take it directly
    Then I see error: "you don't see old_key here"
