Feature: Examine-Triggered Events
  As a player
  I want examining objects to trigger world events
  So that the game world feels alive and reactive

  Background:
    Given the "fountain" has examine-triggered events
    And the "dusty_crate" has examine-triggered events
    And the "vintage_terminal" has examine-triggered events

  Scenario: Examine reveals hidden room
    Given the "fountain" triggers reveal_room at examine level 75
    And target room is "secret_passage_scrapyard"
    When I examine the fountain at level 75
    Then the message is shown: "You discover a hidden passage behind the fountain!"
    And the "secret_passage_scrapyard" room becomes accessible
    And the event is marked as fired (one-time)

  Scenario: Examine spawns hidden NPC
    Given the "dusty_crate" triggers spawn_npc at examine level 50
    And target NPC is "hidden_trader"
    When I examine dusty_crate at level 50
    Then the message is shown: "A figure steps out from behind the crates..."
    And "hidden_trader" NPC is now present in the room
    And the event cannot fire again

  Scenario: Examine unlocks lore entry
    Given the "vintage_terminal" triggers unlock_lore at examine level 60
    And target lore is "day_one_ooze_incident"
    When I examine vintage_terminal at level 60
    Then the message is shown: "The terminal flickers to life..."
    And lore entry is added to player's lore log
    And I can access it with "lore" command

  Scenario: Events fire only once per character
    Given I examined dusty_crate and spawned hidden_trader
    When I examine dusty_crate again
    Then no duplicate NPC spawns
    And no additional event fires

  Scenario: Different characters can trigger same event independently
    Given Character A triggered the event
    When Character B examines the same object
    Then Character B can trigger the event independently
    And Character A's event state is unchanged

  Scenario: Examine below threshold does not trigger event
    Given the fountain requires examine level 75
    When I examine at level 50
    Then no event triggers
    And the normal examine description is shown

  Scenario: Examine above threshold does trigger event
    Given the fountain requires examine level 75
    When I examine at level 80
    Then the event triggers
    And the room/lore/NPC is revealed

  Scenario: Multiple events on same object can exist
    Given the fountain has multiple event triggers
    When I examine at sufficient levels
    Then each applicable event fires
    And all are marked as fired

  Scenario: Event messages shown in examine output
    When I examine an object that triggers an event
    Then the event message appears in the examine output
    And the world state change is immediate
