Feature: Examine-Triggered Quests
  As a player
  I want examining certain items to unlock hidden quests
  So that exploration is rewarded with story content

  Background:
    Given the "fountain" is in the current room
    And examining it can unlock the secret quest

  Scenario: Examine fountain with sufficient level unlocks quest
    Given my examine skill level is 75
    When I type "examine fountain"
    Then I discover a hidden compartment
    And the quest "The Fountain's Secret" is unlocked
    And it appears in my quest log

  Scenario: Examine fountain with low level does not unlock quest
    Given my examine skill level is 50
    When I type "examine fountain"
    Then I see the normal examine description
    And no quest is unlocked

  Scenario: Quest appears in quest log after unlock
    Given examining fountain unlocked "The Fountain's Secret"
    When I type "quests"
    Then I see "The Fountain's Secret" in my quest log

  Scenario: Examine skill grants XP on quest unlock
    Given examining fountain unlocks a quest
    When the quest is unlocked
    Then I gain examine skill XP
    And my examine level may increase

  Scenario: Secret quest reveals hidden path
    Given the "fountain" has a secret compartment
    And examining it at level 75 reveals it
    When the quest unlocks
    Then a new room or passage is revealed
    Or a hidden action becomes available

  Scenario: Examine path skips traditional quest giver
    Given a quest normally requires talking to an NPC
    When I examine the relevant object at high level
    Then the quest can be unlocked through examination alone
    And I bypass the NPC interaction

  Scenario: Multiple examine-triggered quests can exist
    Given multiple objects have examine-triggered quests
    When I examine each relevant object at correct level
    Then each respective quest unlocks independently

  Scenario: Quest unlock shows notification
    Given examining fountain unlocks a quest
    When the unlock occurs
    Then I see a notification: "Quest Unlocked: The Fountain's Secret"
    And the XP gain is shown

  Scenario: Cannot re-trigger already unlocked quest
    Given the "fountain" quest is already unlocked
    When I examine the fountain again
    Then no duplicate quest entry appears
    And no additional XP is granted
