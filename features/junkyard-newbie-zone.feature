Feature: Junkyard Newbie Zone
  As a new player
  I want a tutorial zone to learn game mechanics
  So that I can understand the game before venturing into harder areas

  Background:
    Given I am a new character
    And I have completed character creation
    And I am at the Fountain

  Scenario: Enter the Scrapyard zone
    Given I travel to Foggy Gate
    When I arrive
    Then I should see the entrance room description
    And I should see Guard Marco
    And I should see the rusty chain-link gate
    And I should see exits to New Venice (south) and Piles of Rust (north)

  Scenario: Read zone signage
    Given I am at Foggy Gate
    When I type "look sign" or "read sign"
    Then I should see "VISTA SALVAGE — EST. 1987"
    And I should understand this is an old salvage yard

  Scenario: Talk to quest NPC
    Given I am at Foggy Gate
    And Guard Marco is present
    When I type "talk marco"
    Then I should see Guard Marco's dialogue
    And he should warn me about the Crusher
    And he should offer the "Prove Yourself" quest

  Scenario: Accept starter quest
    Given I am at Foggy Gate
    And Guard Marco has offered "Prove Yourself" quest
    When I accept the quest
    Then the quest should appear in my quest log
    And I should see the objective: "Kill 3 Scrap Rats"
    And I should see the reward: 10 coins

  Scenario: Move between rooms
    Given I am at Foggy Gate
    When I type "north"
    Then I should move to Piles of Rust
    And I should see the room description
    And I should see available exits

  Scenario: Encounter creatures
    Given I am in Piles of Rust
    When I enter the room
    Then I should see Scrap Rats if they are spawned
    And I should see room items if present
    And I should be able to see creature descriptions

  Scenario: Pick up floor items
    Given I am in Piles of Rust
    And a rusty_pipe is on the floor
    When I type "get rusty_pipe"
    Then the rusty_pipe should be in my inventory
    And the item should be removed from the floor
    And I should see a pickup message

  Scenario: Combat with level 1 creature
    Given I am in Piles of Rust
    And a Scrap Rat is present
    When I type "attack rat"
    Then combat should begin
    And I should see combat rounds
    And the Scrap Rat should attack back

  Scenario: Creature flees at low HP
    Given I am fighting a Scrap Rat
    And the Scrap Rat HP is below 2
    When the combat round processes
    Then the Scrap Rat should attempt to flee
    And combat should end if the rat escapes

  Scenario: Complete "Prove Yourself" quest
    Given I have killed 3 Scrap Rats
    And I return to Guard Marco
    When I turn in the quest
    Then I should receive 10 coins
    And I should receive free entry to New Venice
    And the quest should be marked complete

  Scenario: Zone boss blocks progression
    Given I am at Crusher's Den
    And Old Scrap is alive
    When I try to enter Deep Scrap
    Then I should be blocked by Old Scrap
    And I should see a message about the path being guarded

  Scenario: Defeat zone boss
    Given I am fighting Old Scrap
    And Old Scrap has 25 HP
    When I reduce his HP to 0
    Then Old Scrap should die
    And he should drop scrap_machete
    And he should drop junk_crown
    And the path to Deep Scrap should unlock

  Scenario: Loot after boss kill
    Given I have defeated Old Scrap
    When items drop to the floor
    Then I should be able to type "get all" or "loot"
    And all items should go to my inventory

  Scenario: Environmental hazard - Ooze pools
    Given I am in Leaking Pipes
    And there are Ooze pools present
    When I stand in an Ooze pool
    Then I should take 1 damage per tick
    And I should see a warning message about the burning sensation

  Scenario: Secret area locked
    Given I am at Deep Scrap
    And I do not have the pre_ooze_key
    When I try to enter Buried Bunker
    Then I should see "This area is locked"
    And I should need the key or quest completion to enter

  Scenario: Unlock secret area with key
    Given I have found the pre_ooze_key
    And I am at Deep Scrap
    When I type "use key" or enter the bunker
    Then the Buried Bunker should unlock
    And I should be able to enter

  Scenario: Secret area contains rare loot
    Given I have entered the Buried Bunker
    When I look around
    Then I should see the vintage_terminal
    And I should see the pristine_laser_pistol
    And I should see the survival_manual

  Scenario: Interact with terminal for lore
    Given I am in the Buried Bunker
    And the vintage_terminal is present
    When I type "use terminal"
    Then I should see lore fragments about the Four Heroes
    And I should see logs from "Day 1" of the Ooze Incident
    And I should see coordinates to another location

  Scenario: Zone level range enforcement
    Given I am level 5
    And I try to enter the Scrapyard
    When I enter Foggy Gate
    Then I should be warned this is a newbie zone
    And I should still be allowed to enter
    But XP rewards should be reduced for overleveled players

  Scenario: Creature respawn timers
    Given I killed a Scrap Rat
    And 5 minutes have passed
    When I return to Piles of Rust
    Then the Scrap Rat should have respawned
    And I should be able to fight it again

  Scenario: Boss respawn timer
    Given I killed Old Scrap
    And 30 minutes have passed
    When I return to Crusher's Den
    Then Old Scrap should have respawned
    And his loot should be available again

  Scenario: Quest giver - Scavenger Jane
    Given I am in Deep Scrap
    And Scavenger Jane is present
    When I type "talk jane"
    Then she should offer the "Ooze Samples" quest
    And she should teach me the scavenge skill if I complete her quest

  Scenario: Trade with NPC
    Given I am in Deep Scrap
    And I have 5 glowing_goo
    When I trade with Scavenger Jane
    Then I should receive a repair_kit
    And my glowing_goo should be consumed

  Scenario: Pack hunter behavior
    Given I am in Metal Maze
    And there are 2 Junk Dogs
    When I attack one Junk Dog
    Then the other Junk Dog should join the fight
    And I should be fighting both creatures

  Scenario: Ooze spawn explosion
    Given I am fighting an Ooze Spawn
    When the Ooze Spawn dies
    Then it should explode
    And I should take 1 damage (AoE)
    And I should see an explosion message