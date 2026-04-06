Feature: Combat Action Costs
  As a combat system
  I need to define tick costs for all actions
  So that players understand how long their actions take to resolve

  Background:
    Given a combat session is active

  Scenario: Basic attack has 1 tick cost
    Given player is in combat
    When player uses "attack" action
    Then action should resolve in 1 tick
    And player can act again immediately after

  Scenario: Defend action applies buff immediately
    Given player is in combat
    When player uses "defend" action
    Then defend applies defense bonus immediately
    And buff lasts until player's next action

  Scenario: Flee attempt resolves in 1 tick
    Given player is in combat with enemy
    When player uses "flee" action
    Then flee attempt resolves in 1 tick
    And result is determined

  Scenario: Item use has 1 tick cost
    Given player has a healing potion
    When player uses "item" action
    Then item resolves in 1 tick

  Scenario: Wait action costs 1 tick
    Given player uses "wait" action
    Then action resolves in 1 tick

  Scenario: Slash talent has 1 tick cost
    Given player has "slash" talent equipped
    When player uses slash
    Then slash resolves in 1 tick
    And slash deals blade damage

  Scenario: Heavy strike has 2 tick channel time
    Given player has "heavy_strike" talent equipped
    When player starts channeling heavy_strike
    Then heavy_strike takes 2 ticks to resolve
    And player cannot act during channel

  Scenario: Heavy strike shows channeling message
    Given player started channeling heavy_strike
    And 1 tick has passed
    When player looks at combat screen
    Then message should show "channeling heavy_strike (1 tick left)"

  Scenario: Heavy strike resolves and deals damage
    Given player has been channeling heavy_strike for 2 ticks
    When heavy_strike resolves
    Then heavy_strike should deal more damage than basic attack
    And player can act again

  Scenario: Parry is a reaction (0 tick cost)
    Given player has "parry" talent equipped
    When player is attacked on the same tick
    And player uses parry
    Then parry resolves instantly (0 ticks)
    And parry can negate the incoming attack

  Scenario: Shield bash has 2 tick cost
    Given player has "shield_bash" talent equipped
    When player uses shield_bash
    Then shield_bash takes 2 ticks to complete
    And shield_bash can interrupt enemies

  Scenario: Battle cry has 1 tick cost
    Given player has "battle_cry" talent equipped
    When player uses battle_cry
    Then battle_cry resolves in 1 tick
    And battle_cry applies debuff to enemies

  Scenario: Second wind has 2 tick channel time
    Given player has "second_wind" talent equipped
    When player starts channeling second_wind
    Then second_wind takes 2 ticks to resolve
    And second_wind heals the player

  Scenario: Hail storm has 3 tick charge time
    Given player has "hail_storm" talent equipped
    When player starts channeling hail_storm
    Then hail_storm takes 3 ticks to complete
    And hail_storm performs double attacks

  Scenario: Instant actions resolve same tick
    Given player has an instant action queued
    When combat tick occurs
    Then instant action should resolve immediately
    And result should be shown in combat log

  Scenario: Channel actions lock player
    Given player is channeling an action
    When player tries to use another action
    Then player should be prevented from acting
    And message should show "You are channeling..."

  Scenario: Charge actions can be interrupted
    Given player is charging an action
    When enemy uses shield_bash
    Then player's charge action should be cancelled

  Scenario Outline: Action tick costs are correct
    Given player uses "<action>" action
    When action resolves
    Then action should take "<cost>" tick(s) to complete

    Examples:
      | action      | cost |
      | attack      | 1    |
      | defend      | 0    |
      | flee        | 1    |
      | item        | 1    |
      | wait        | 1    |
      | slash       | 1    |
      | heavy_strike| 2    |
      | parry       | 0    |
      | shield_bash | 2    |
      | battle_cry  | 1    |
      | second_wind | 2    |
      | hail_storm  | 3    |